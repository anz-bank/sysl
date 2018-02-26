def pytest_addoption(parser):
    parser.addoption('--syslexe', action='store', default='')
    parser.addoption('--reljamexe', action='store', default='')


def pytest_generate_tests(metafunc):
    if 'syslexe' in metafunc.fixturenames:
        metafunc.parametrize('syslexe', [metafunc.config.option.syslexe])
    if 'reljamexe' in metafunc.fixturenames:
        metafunc.parametrize('reljamexe', [metafunc.config.option.reljamexe])
