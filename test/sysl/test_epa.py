# -*- coding: utf-8 -*-

from sysl.core import syslloader, syslints
from sysl.util import debug
import unittest
import re, os, sys
import traceback
import argparse as ap
import tempfile
from os import path 

class TestEpa(unittest.TestCase):
  def setUp(self):
    self.outpath  = tempfile.gettempdir()

  def integration_view_helper(self, modulename, d):
    (module, deps, _) = syslloader.load(modulename, True, '.')

    args = ap.Namespace(**d)

    if not args.exclude and args.project:
      args.exclude = {args.project}

    return syslints.integration_views(module, deps, args)

  def test_ints(self):

    try:
      d = {
        'project'   : 'Test EPA :: Integrations',
        'exclude'   : '',
        'output'    : path.join(self.outpath, 'test_ints-ints.png'),
        'plantuml'  : '',
        'clustered' : '',
        'title'     : 'Test EPA',
        'epa'       : False,
        'filter'    : None,
        'verbose'   : ''}

      out = self.integration_view_helper('/test/data/test_epa', d)

      self.assertTrue('_0 --> _1' in out[0])

    except (IOError, Exception) as e:
      self.fail(traceback.format_exc())

  def test_epa(self):

    try:

      d = {
        'project'   : 'Test EPA :: Integrations',
        'exclude'   : '',
        'output'    : path.join(self.outpath, 'test_epa-ints.png'),
        'plantuml'  : '',
        'clustered' : '',
        'epa'       : True,
        'filter'    : None,
        'title'     : 'Test EPA',
        'verbose'   : ''}

      out = self.integration_view_helper('/test/data/test_epa', d)

      #import pdb; pdb.set_trace()
      self.assertTrue(re.search('_0 -.*> _1', out[0]))
      self.assertTrue(re.search('_1 -.*> _2', out[0]))

    except (IOError, Exception) as e:
      self.fail(traceback.format_exc())

  def test_epa_repeated_calls(self):

    try:

      d = {
        'project'   : 'Test EPA :: Integrations',
        'exclude'   : '',
        'output'    : path.join(self.outpath, 'test_epa_repeated_calls-ints.png'),
        'plantuml'  : '',
        'clustered' : '',
        'epa'       : True,
        'filter'    : None,
        'title'     : 'Test EPA Repeated Calls',
        'verbose'   : ''}

      out = self.integration_view_helper('/test/data/test_epa_repeated_calls', d)

      #import pdb; pdb.set_trace()
      self.assertTrue('state "**App1 Input Method 1 client**" as _2' in out[0])
      self.assertTrue('state "**App1 Input Method 1**" as _3' in out[0])

      self.assertTrue(re.search('_1 -.*> _2', out[0]))
      self.assertTrue(re.search('_2 -.*> _3', out[0]))

    except (IOError, Exception) as e:
      self.fail(traceback.format_exc())

  def test_int_repeated_calls(self):

    try:

      d = {
        'project'   : 'Test EPA :: Integrations',
        'exclude'   : '',
        'output'    : path.join(self.outpath, 'test_int_repeated_calls-ints.png'),
        'plantuml'  : '',
        'clustered' : '',
        'epa'       : False,
        'filter'    : None,
        'title'     : 'Test EPA Repeated Calls',
        'verbose'   : ''}

      out = self.integration_view_helper('/test/data/test_epa_repeated_calls', d)

      self.assertTrue('_0 --> _1' in out[0])
      self.assertFalse('_1 --> _3' in out[0])
      self.assertFalse('_2 --> _3' in out[0])

    except (IOError, Exception) as e:
      self.fail(traceback.format_exc())

  def test_ignore_keyword(self):

    try:

      d = {
        'project'   : 'Test EPA :: Integrations',
        'exclude'   : '',
        'output'    : path.join(self.outpath, 'test_epa_ignore_keyword-ints.png'),
        'plantuml'  : '',
        'clustered' : '',
        'epa'       : True,
        'filter'    : None,
        'title'     : 'Test EPA Ignore Keyword',
        'verbose'   : ''}

      out = self.integration_view_helper('/test/data/test_epa_ignore_keyword', d)

      self.assertFalse('state "**.. * <- ***"' in out[0])

    except (IOError, Exception) as e:
      self.fail(traceback.format_exc())

  def test_labels(self):

    try:

      d = {
        'project'   : 'Test EPA :: Integrations',
        'exclude'   : '',
        'output'    : path.join(self.outpath, 'test_epa_labels-ints.png'),
        'plantuml'  : '',
        'clustered' : '',
        'epa'       : True,
        'filter'    : None,
        'title'     : 'Test EPA Labels',
        'verbose'   : ''}

      out = self.integration_view_helper('/test/data/test_epa_ignore_keyword', d)

      self.assertTrue('**«INT-001»**' in out[0])

    except (IOError, Exception) as e:
      self.fail(traceback.format_exc())

  def test_labels_for_events(self):

    try:

      d = {
        'project'   : 'Test EPA :: Events',
        'exclude'   : '',
        'output'    : path.join(self.outpath, 'test_epa_labels_for_events-ints.png'),
        'plantuml'  : '',
        'clustered' : '',
        'epa'       : True,
        'filter'    : None,
        'title'     : 'Test EPA Labels',
        'verbose'   : ''}

      #import pdb; pdb.set_trace()

      out = self.integration_view_helper('/test/data/test_epa_labels_for_events', d)

      self.assertTrue('**«INT-001»**' in out[0])

    except (IOError, Exception) as e:
      self.fail(traceback.format_exc())

  def test_patterns(self):

    try:

      d = {
        'project'   : 'Test EPA :: Integrations',
        'exclude'   : '',
        'output'    : path.join(self.outpath, 'test_epa_patterns-ints.png'),
        'plantuml'  : '',
        'clustered' : '',
        'epa'       : True,
        'filter'    : None,
        'title'     : 'Test EPA Patterns',
        'verbose'   : ''}

      out = self.integration_view_helper('/test/data/test_epa_patterns', d)
      print '==================', out[0]
      self.assertTrue('** <color green> → soap</color>**' in out[0])

    except (IOError, Exception) as e:
      self.fail(traceback.format_exc())

  def test_missing_patterns(self):

    try:

      d = {
        'project'   : 'Test EPA :: Integrations',
        'exclude'   : '',
        'output'    : path.join(self.outpath, 'test_epa_missing_patterns-ints.png'),
        'plantuml'  : '',
        'clustered' : '',
        'epa'       : True,
        'filter'    : None,
        'title'     : 'Test EPA Patterns',
        'verbose'   : ''}


      out = self.integration_view_helper('/test/data/test_epa_missing_patterns', d)

      #import pdb; pdb.set_trace()

      self.assertTrue('** <color red>pattern?</color>**' in out[0])

    except (IOError, Exception) as e:
      self.fail(traceback.format_exc())

  def test_missing_labels(self):

    try:

      d = {
        'project'   : 'Test EPA :: Integrations',
        'exclude'   : '',
        'output'    : path.join(self.outpath, 'test_epa_missing_labels-ints.png'),
        'plantuml'  : '',
        'clustered' : '',
        'epa'       : True,
        'filter'    : None,
        'title'     : 'Test EPA Patterns',
        'verbose'   : ''}

      # import pdb; pdb.set_trace()

      out = self.integration_view_helper('/test/data/test_epa_missing_labels', d)

      # import pdb; pdb.set_trace()

      self.assertTrue('<color red>(missing INT)</color>' in out[0])

    except (IOError, Exception) as e:
      self.fail(traceback.format_exc())

  # def test_events(self):
    
  #   try:

  #     d = {
  #       'project'   : 'Test EPA :: Events',
  #       'exclude'   : '',
  #       'output'    : path.join(self.outpath, 'test_epa_events-ints.png'),
  #       'plantuml'  : '',
  #       'clustered' : '',
  #       'epa'       : True,
  #       'title'     : 'Test EPA Events',
  #       'verbose'   : ''}

  #     out = self.integration_view_helper('/test/data/test_epa_events', d)

  #     #import pdb; pdb.set_trace()
  #     self.assertTrue('state "**Test EPA :: App1 -> App1 Event**" as _0' in out[0])
  #     self.assertTrue('state "**Trigger**" as _2' in out[0])
  #     self.assertTrue('state "**App1 Event**" as _1' in out[0])

  #     self.assertTrue(re.search('_2.*_1', out[0]))
  #     self.assertTrue(re.search('_1.*_0', out[0]))

  #     self.assertFalse(re.search('App1 Event client', out[0]))

  #   except (IOError, Exception) as e:
  #     self.fail(traceback.format_exc()) 

  def test_int_passthrough(self):

    try:
      
      d = {
        'project'   : 'Test EPA :: Passthrough',
        'exclude'   : '',
        'output'    : path.join(self.outpath, 'test_epa_passthrough-ints.png'),
        'plantuml'  : '',
        'clustered' : '',
        'epa'       : False,
        'title'     : 'Test EPA Passthrough',
        'verbose'   : '',
        'filter'    : ''}

      out = self.integration_view_helper('/test/data/test_epa_passthrough', d)
      
    except (IOError, Exception) as e:
      self.fail(traceback.format_exc())


if __name__ == '__main__':
  debug.init()
  unittest.main()
