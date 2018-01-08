from os import path

import pytest

from sysl.core import syslparse
from sysl.core import syslloader
from sysl.proto import sysl_pb2


def test_annotations_simple():
	sysl_input = '''\
ThatApp:
	foo:
		...
MyApp:
    ep1:
        ThatApp <- foo
MyApp:
    .. * <- * :
        ThatApp <- foo [~abc]

'''
	module = sysl_pb2.Module()
	syslparse.Parser().parse(sysl_input.splitlines(), '', module)
	syslloader.infer(module)
	assert 2 == len(module.apps['MyApp'].endpoints)
	ep1 = module.apps['MyApp'].endpoints['ep1'].stmt
	assert 1 == len(ep1)
	assert 'foo' == ep1[0].call.endpoint
	assert 'ThatApp' == ep1[0].call.target.part[0]
	assert 1 == len(ep1[0].attrs['patterns'].a.elt)
	assert 'abc' == ep1[0].attrs['patterns'].a.elt[0].s


@pytest.mark.parametrize("endpoint_wildcard", [
    '.. * <- *',
    '.. *<-*',
    '..   *   <-  *  ',
    '*',
    ' * ',
])
def test_annotation_variants(endpoint_wildcard):
	SYSL_INPUT = '''\
ThatApp:
	foo:
		...
	bar:
		...
MyApp:
    ep1:
        ThatApp <- foo

    ep2:
        ThatApp <- foo

    ep3:
        ThatApp <- bar

MyApp:
    {ENDPOINT_WILDCARD}:
        ThatApp <- foo [~abc]

'''
	sysl_input = SYSL_INPUT.replace('{ENDPOINT_WILDCARD}', endpoint_wildcard)
	module = sysl_pb2.Module()
	syslparse.Parser().parse(sysl_input.splitlines(), '', module)
	syslloader.infer(module)

	assert 4 == len(module.apps['MyApp'].endpoints)

	ep1 = module.apps['MyApp'].endpoints['ep1'].stmt
	assert 1 == len(ep1)
	assert 'foo' == ep1[0].call.endpoint
	assert 'ThatApp' == ep1[0].call.target.part[0]
	assert 1 == len(ep1[0].attrs['patterns'].a.elt)
	assert 'abc' == ep1[0].attrs['patterns'].a.elt[0].s

	ep2 = module.apps['MyApp'].endpoints['ep2'].stmt
	assert 1 == len(ep2)
	assert 'foo' == ep2[0].call.endpoint
	assert 'ThatApp' == ep2[0].call.target.part[0]
	assert 1 == len(ep2[0].attrs['patterns'].a.elt)
	assert 'abc' == ep2[0].attrs['patterns'].a.elt[0].s

	ep3 = module.apps['MyApp'].endpoints['ep3'].stmt
	assert 1 == len(ep3)
	assert 'bar' == ep3[0].call.endpoint
	assert 'ThatApp' == ep3[0].call.target.part[0]
	assert 0 == len(ep3[0].attrs['patterns'].a.elt)

