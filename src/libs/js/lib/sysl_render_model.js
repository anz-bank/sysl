((root, factory) => {
  if (typeof exports === 'object') {
    module.exports = root.sysl_render_model = factory();  // CommonJS
  } else if (typeof define === 'function' && define.amd) {
    define([], () => root.sysl_render_model = factory());  // AMD
  } else {
    root.sysl_render_model = factory();  // Browser
  }
})(this, () => {
  'use strict';

  const __ = Object.assign(React.createElement, React.DOM);

  const EntityViewer = props => {
    const e = props.entity;
    const fields = e._fields;
    const key = e.constructor._pkFields;
    const fieldDefs = e.constructor._allFields;

    return (
      __.table({
            style: {
              border: '1px solid #bde',
              borderCollapse: 'separate',
              borderRadius: '5px',
              fontSize: '10pt',
              margin: '0.5em',
              padding: '0.5em',
            },
          },
        __.tbody({},
          Object.keys(e._fields).filter(f => (''+e._fields[f]).trim()).sort().map(fname => {
            const f = fields[fname];
            const isBool = fieldDefs[fname] == 3;  // sysl.proto/Type.Primitive
            return __.tr({key: fname},
              __.th({
                    style: {
                      paddingLeft: '0.5em',
                      paddingRight: '0.5em',
                      borderRight: '1px solid #cef',
                    },
                  },
                fname,
                isBool ? '?' : ''),
              __.td({
                    style: {
                      paddingLeft: '0.5em',
                      paddingRight: '0.5em',
                    },
                  },
                isBool ? __.input({type: 'checkbox', checked: f, disabled: true}) : ''+f
              )
            )
          })
        )
      )
    )
  };

  const TableViewer = props => {
    // TODO: Proper keys
    let i = 0;

    return (
      __.div({className: "panel panel-default"},
        __.div({className: "panel-heading"},
          __.h3({className: "panel-title"},
            props.name,
            ' x ',
            props.data.length
          )
        ),
        __.div({className: "panel-body"},
          props.data.map(entity =>
            __(EntityViewer, {key: ++i, entity})
          )
        )
      )
    )
  };

  const ModelViewer = props => {
    const m = props.model;
    const types = m._def.types;
    return __.div({},
      __.h1({}, m._def.model),
      Object.keys(types).filter(tname => m[tname].length).sort().map(tname =>
        __(TableViewer, {
          key: tname,
          name: tname,
          data: m[tname],
          def: types[tname],
        })
      )
    )
  };

  return {__, ModelViewer};
});
