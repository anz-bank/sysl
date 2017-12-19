"""Confluence API"""

import base64
import getpass
import itertools
import json
import mimetypes
import operator
import os
import re
import urllib
import requests
import sys
import uuid

from sysl.util import cache


USER_ENV = 'CONFLUENCE_USER'

VERIFY_CERTS = bool(int(os.getenv('SYSL_VERIFY_CERTS', '1')))


class _LazySession(object):
    """Class to open urls with cached basic auth."""

    class _Proxy(object):
        def __init__(self, lazy_session, name):
            self.lazy_session = lazy_session
            self.name = name

        def __call__(self, url, *args, **kwargs):
            session = self.lazy_session.session
            if self.name not in (
                    'get', 'head') and not self.lazy_session._auth_added:
                print >>sys.stderr, '{} \033[1;36m{}\033[0m'.format(
                    self.name.upper(), url)
                user = os.getenv(USER_ENV)
                env_hint = ' (={} env var)'.format(USER_ENV)
                if user is None:
                    user = getpass.getuser()
                    env_hint = ' (=USER env var; set {} to override)'.format(
                        USER_ENV)
                print '\n\033[1mConfluence username:\033[0m ' + user + env_hint
                session.auth = (
                    user, getpass.getpass('\033[1mConfluence password:\033[0m '))
                self.lazy_session._auth_added = True
            return getattr(session, self.name)(url, *args, **kwargs)

    def __init__(self):
        self._session = None
        self._auth_added = False

    @property
    def session(self):
        if self._session is None:
            self._session = requests.Session()
        return self._session

    def __getattr__(self, name):
        assert name in ('get', 'head', 'put', 'post', 'delete', 'patch'), name
        return self._Proxy(self, name)


CONFLUENCE_SERVER = os.getenv('CONFLUENCE_SERVER')
LAZY_SESSION = _LazySession()


def _confluence_server():
    if not CONFLUENCE_SERVER:
        raise RuntimeError('Missing CONFLUENCE_SERVER environment variable')
    return CONFLUENCE_SERVER


def _confluence_to_http(confluence_url):
    """Convert confluence url to http url."""
    def api(*path, **kwargs):
        """Return an api url for a path."""
        assert set(kwargs) <= {'query'}

        url = _confluence_server() + '/rest/api/' + '/'.join(path)
        if 'query' in kwargs:
            url += '?' + '&'.join(k + '=' + urllib.quote(v)
                                  for (k, v) in kwargs['query'].iteritems())
        return url

    def content(content_id=None, *subpath, **kwargs):
        """Return an api url for content."""
        if content_id is not None:
            subpath = (content_id,) + subpath
        return api('content', *subpath, **kwargs)

    def attachment(content_id, att_id=None, *subpath):
        """Return an api url for an attachment."""
        if att_id is not None:
            subpath = (att_id,) + subpath
        return content(content_id, 'child', 'attachment', *subpath)

    (space, title, tail, att, query) = re.match(
        r'confluence://([^/]+)/(.+?)(/(?:@(.*?)|(?:[^@].*?)))?(\?.*)?$', confluence_url).groups()

    content(None, query={'space': space, 'title': title})
    response = LAZY_SESSION.get(
        content(None, query={'space': space, 'title': title}), verify=VERIFY_CERTS)
    response.raise_for_status()
    root = response.json()
    if att is None:
        return root['results'][0]['_links']['self'] + \
            (tail or '') + (query or '')

    children = [result['id'] for result in root['results']]

    for result in root['results']:
        child_id = result['id']
        response = LAZY_SESSION.get(
            content(child_id, 'child', 'attachment'), verify=VERIFY_CERTS)
        response.raise_for_status()
        root = response.json()
        for result in root['results']:
            if result['title'] == att:
                return (
                    attachment(child_id, result['id']).encode('utf-8'),
                    result['_links']['download'])
    raise RuntimeError('{} not found'.format(confluence_url))


def get_storage(confluence_url):
    response = LAZY_SESSION.get(_confluence_to_http(
        confluence_url) + '?expand=body.storage,ancestors,version', verify=VERIFY_CERTS)
    response.raise_for_status()
    root = response.json()
    body = root['body']
    return (body['storage']['value'], root['title'],
            root['ancestors'][-1:], root['version']['number'])


def put_storage(confluence_url, storage, title,
                ancestors, version, minorEdit=False):
    http_url = _confluence_to_http(confluence_url)
    page_id = http_url.rpartition('/')[-1]

    data = {
        'id': page_id,
        'title': title,
        'type': 'page',
        'space': {
            'key': 'PAD',  # TODO: extract from confluence_url
        },
        'body': {
            'storage': {
                'value': storage,
                'representation': 'storage',
            },
        },
        'ancestors': ancestors,
        'version': {
            'number': version,
            'minorEdit': minorEdit,
        },
    }
    LAZY_SESSION.put(
        _confluence_to_http(confluence_url),
        headers={
            'Content-Type': 'application/json',
            'X-Atlassian-Token': 'nocheck',
        },
        data=json.dumps(data),
        verify=VERIFY_CERTS).raise_for_status()


def upload_attachment(confluence_att_url, file_handle,
                      expire_cache=False, dry_run=False):
    """Upload an attachment to a confluence url."""
    if expire_cache:
        cache.expire(confluence_att_url)

    if isinstance(confluence_att_url, unicode):
        confluence_att_url = confluence_att_url.encode('utf-8')
    filename = confluence_att_url.rpartition('/@')[-1]

    def get_old():
        (url, download_path) = _confluence_to_http(confluence_att_url)
        url += '/data'
        old = LAZY_SESSION.get(
            _confluence_server() + download_path, verify=VERIFY_CERTS)
        old.raise_for_status()
        return (url, download_path, old.content)

    (url, download_path, old_content) = cache.get(confluence_att_url, get_old)

    new_content = file_handle.read()

    if old_content == new_content:
        return None

    if not dry_run:
        rsp = LAZY_SESSION.post(
            url,
            headers={'X-Atlassian-Token': 'nocheck'},
            data={'comment': 'sysl-generated image', 'minorEdit': 'true'},
            files={'file': (filename, new_content)},
            verify=VERIFY_CERTS)
        rsp.raise_for_status()

        cache.put(confluence_att_url, (url, download_path, new_content))

    return rsp.text


def upload_file(confluence_page_url, path):
    """Upload a file as an attachment to a confluence url."""
    filename = os.path.split(path)[1]
    upload_attachment(confluence_page_url + '/@' + filename, open(path))


def search_label(label):
    response = LAZY_SESSION.get(
        (_confluence_server() +
         '/rest/searchv3/1.0/search?queryString=labelText=' +
         label),
        verify=VERIFY_CERTS)
    response.raise_for_status()
    data = response.json()
    return data['results']
