#!/usr/bin/env python

import argparse
import os

import lxml.etree as ET

from sysl.util import writer


BASE_TYPE_MAP = {
    'xs:decimal': ('decimal', True),
    'xs:string': ('string', False),
    'xs:integer': ('int', True),
    'xs:date': ('date', False),
    'xs:time': ('string', False),
}

NS = {
    'xs': 'http://www.w3.org/2001/XMLSchema'
}


def find(node, expr):
    return node.find(expr, namespaces=NS)


def findall(node, expr):
    return node.findall(expr, namespaces=NS)


def get_attr(node, expr, attr, mapf=None, default=None):
    if mapf is None:
        def mapf(x):
            return x

    e = find(node, expr)
    return mapf(e.get(attr)) if e is not None else default


def syslForField(field):
    setof = 'set of ' * (int(field.get('maxOccurs') or 1) > 1)
    optional = '?' * (not setof and field.get('minOccurs') == '0')

    type_ = field.get('type')
    if type_:
        return setof + type_ + optional

    fr = find(field, './xs:simpleType/xs:restriction')
    if fr is not None:
        (base, numeric) = BASE_TYPE_MAP[fr.get('base')]

        type_specs = []

        if numeric:
            dig = get_attr(fr, './xs:totalDigits', 'value', int)
            fdig = get_attr(fr, './xs:fractionDigits', 'value', None, '')
            if dig is not None:
                type_specs.append('{}{}'.format(dig, fdig and ('.' + fdig)))

            mini = get_attr(fr, './xs:minInclusive', 'value')
            maxi = get_attr(fr, './xs:maxInclusive', 'value')
            if mini is not None and maxi is not None:
                type_specs.append('{}..{}'.format(mini, maxi))
            elif mini is not None or maxi is not None:
                raise RuntimeError('xs:minInclusive xor xs:maxInclusive')
        else:
            maxl = get_attr(fr, './xs:maxLength', 'value')
            if maxl is not None:
                type_specs.append(maxl)

        # TODO: syslparse.py to support multiple type-specs
        type_qual = '(' + ', '.join(type_specs[:1]) + ')' if type_specs else ''

        return '{}{}{}{}'.format(setof, base, type_qual, optional)

    raise RuntimeError('unexpected simpleType spec: {}'.format(
        ET.tostring(field)))


def main():
    argp = argparse.ArgumentParser(
        description='import xsd to sysl')

    argp.add_argument('--appname', help='output appname')
    argp.add_argument('--package', help='output package')

    argp.add_argument('input', help='xsd input file')
    argp.add_argument('output', help='sysl output file')

    args = argp.parse_args()

    root = ET.XML(open(args.input).read())

    w = writer.Writer('sysl')

    w('{}{}:',
      args.appname or os.path.splitext(os.path.basename(args.input))[0],
      ' [package="{}"]'.format(args.package) if args.package else '')
    with w.indent():
        toplevel = [
            (e.get('name'), list(e)[0])
            for e in findall(root, './xs:element[@name]')]
        secondary = [
            (ct.get('name'), ct)
            for ct in findall(root, './xs:complexType[@name]')]

        for (i, (tname, t)) in enumerate(toplevel + secondary):
            w('\n!type {}{}:'[not i:],
              tname,
              (' [xml_order=["index"]]'
               if any(
                   attr.get('name') == 'index'
                   for attr in findall(t, './xs:attribute'))
               else ''))

            with w.indent():
                for attr in findall(t, './xs:attribute'):
                    (type_, _) = BASE_TYPE_MAP[attr.get('type')]
                    w('{} <: {}{} [~xml_attribute]',
                      attr.get('name'),
                      type_,
                      '' if attr.get('use') == 'required' else '?')

                for field in (
                        findall(t, './xs:sequence/xs:element') +
                        findall(t, './xs:all/xs:element')):
                    w('{} <: {}', field.get('name'), syslForField(field))

        open(args.output, 'w').write(str(w))


if __name__ == '__main__':
    main()
