def export_plantuml(model, with_fields):
  w = writer.Writer('plantuml')

  with w.uml():
    for (tname, ft, t) in wrapped_types(model):
      if t.HasField(u'relation'):
        pkey = primary_key_params(t)
        pkey_fields = {f for (_, f, _) in pkey}
        fkey_fields = {f for (f, _, _) in foreign_keys(t)}

        # Yes, I know; PlantUML isn't Java. But this fits the bill.
        with java.Block(w, 'class {}', tname):
          if with_fields:
            for (fname, f) in t.relation.attr_defs.iteritems():
              if fname in pkey_fields:
                w('{{static}} +{}', java.name(fname))
            for (fname, f) in t.relation.attr_defs.iteritems():
              (_, _, type_) = typeref(f)
              if fname not in pkey_fields:
                if fname in fkey_fields:
                  w('+{}', fname)

    for (tname, t) in sorted(app.types.iteritems()):
      if t.HasField(u'relation'):
        for (fname, _, type_info) in foreign_keys(t):
          w('{} -- {}{}',
            type_info.parent_path,
            '' if {fname} == pkey_fields else '"0..*" ',
            tname)

  diagutil.output_plantuml(
    diagutil.OutputArgs(
      output=root_class + '.png',
      plantuml=None,
      verbose=True,
      expire_cache=False,
      dry_run=False),
    str(w))
