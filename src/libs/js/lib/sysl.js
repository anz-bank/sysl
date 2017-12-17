((factory) => {
  if (typeof exports === 'object') {
    module.exports = factory();  // CommonJS
  } else if (typeof define === 'function' && define.amd) {
    define([], () => factory());  // AMD
  } else {
    window.sysl = factory();  // Browser
  }
})(() => {
  'use strict';

  const sysl_js_version = '1.0';

  const CamelCase = s => s[0].toUpperCase() + s.slice(1);

  class View {
    get isEmpty() {
      for (const e of this) {
        return false;
      }
      return true;
    }

    get size() {
      return this.count();
    }

    get singleOrNone() {
      let result;
      for (const e of this) {
        if (result !== undefined) {
          throw 'Too many items';
        }
        result = e;
      }
      return result;
    }

    get single() {
      const result = this.singleOrNone;
      if (result === undefined) {
        throw 'No items';
      }
      return result;
    }

    count(limit) {
      if (limit === undefined) {
        limit = Number.MAXINT;
      }
      let n = 0;
      for ({} of this) {
        if (++n >= limit) {
          break;
        }
      }
      return n;
    }

    map(func) {
      const self = this;
      return new (class extends View {
        * [Symbol.iterator]() {
          for (const e of self) {
            yield func(e);
          }
        }
      });
    }

    where(pred) {
      const self = this;
      return new (class extends View {
        * [Symbol.iterator]() {
          for (const e of self) {
            if (pred(e)) {
              yield e;
            }
          }
        }
      });
    }

    rank(name, comp) {
      const self = this;
      return new (class extends View {
        * [Symbol.iterator]() {
          let last_e, i = 0, r;
          for (const e of Array.from(self).sort(comp)) {
            if (last_e === undefined || comp(last_e, e)) {
              r = i;
            }
            yield Object.assign({[name]:r}, e._toObject());
            last_e = e;
            ++i;
          }
        }
      });
    }

    first(n, comp) {
      const self = this;
      return new (class extends View {
        * [Symbol.iterator]() {
          for (const e of Array.from(self).sort(comp).slice(0, n)) {
            yield e;
          }
        }
      });
    }

    // Not a view because views don't support ordering (and never will).
    * orderBy(comp) {
      for (const e of Array.from(this).sort(comp)) {
        yield e;
      }
    }
  }

  const validate_date = v => new Date(v).toString() != 'Invalid Date';

  const uuid_re = /[0-9A-Fa-f]{8}-(?:[0-9A-Fa-f]{4}-){3}[0-9A-Fa-f]{12}/;
  const validate_uuid = v => typeof v == 'string' && v.match(uuid_re);

  const TYPE_INFOS = {
    1 : ['EMPTY'    , v => false                  , v => v                ],
    2 : ['ANY'      , v => true                   , v => v                ],
    3 : ['BOOL'     , v => typeof v == 'boolean'  , v => v                ],
    4 : ['INT'      , v => typeof v == 'number'   , v => v                ],
    5 : ['FLOAT'    , v => typeof v == 'number'   , v => v                ],
    6 : ['STRING'   , v => typeof v == 'string'   , v => v                ],
    7 : ['BYTES'    , v => false                  , v => v                ],
    8 : ['STRING_8' , v => typeof v == 'string'   , v => v                ],
    9 : ['DATE'     , validate_date               , v => new Date(v)      ],
    10: ['DATETIME' , validate_date               , v => new Date(v)      ],
    11: ['XML'      , v => typeof v == 'string'   , v => v                ],
    12: ['DECIMAL'  , v => typeof v == 'string'   , v => v                ],
    13: ['UUID'     , validate_uuid               , v => v.toLowerCase()  ],
  };

  const TYPE_HELPERS = {};
  const TYPE_MAP = {};
  for (const i in TYPE_INFOS) {
    const tinfo = TYPE_INFOS[i];
    const [typename, valid, imp] = tinfo;
    TYPE_HELPERS[i] = {
      imp: v => {
        if (!valid(v)) {
          throw (
            `Invalid input: ${JSON.stringify(v)} is not ` +
            `${typename.toLowerCase()}-convertible`);
        }
        return imp(v);
      },
    };
    TYPE_MAP[typename] = i;
  }

  const privateProperties = (object, properties) => {
    for (const name in properties) {
      const value = properties[name];
      Object.defineProperty(object, name, {
        __proto__: null,
        value: value,
      })
    }
  };

  const createEntityClass = (tname, entityDef, Model, model) => {
    const Entity = class {
      constructor(table, data) {
        privateProperties(this, {
          _table: table,
          _fields: {},
          _def: entityDef,
        });
        this._fromObject(data);
        Object.seal(this);
      }

      _toObject() {
        const object = {};
        for (const fname in Entity.prototype) {
          const v = this[fname];
          if (v != null) {
            object[fname] = v;
          }
        }
        return object;
      }

      _fromObject(object) {
        for (const fname in object) {
          this._fields[fname] = object[fname];
        }

        return this;
      }

      _toJSON(replacer, space) {
        return JSON.stringify(this._toObject(), replacer, space);
      }

      _fromJSON(json, reviver) {
        return this._fromObject(JSON.parse(json, reviver));
      }

      _validate() {
        if (Entity._requiredFields.length) {
          const missingReqd = [];
          for (const f of Entity._requiredFields) {
            if (this[f] == null) {
              missingReqd.push(f);
            }
          }
          if (missingReqd.length) {
            throw 'Missing primary key fields: ' + missingReqd.join(', ');
          }
        }
        return this;
      }

      static _pkOf(object) {
        return JSON.stringify(Entity._pkFields.map(f => object[f]));
      }
    }

    Entity._allFields = {};
    Entity._assignableFields = [];
    Entity._requiredFields = [];
    Entity._pkFields = [];
    Entity._ftypes = {};

    for (const fspec in entityDef) {
      let m = fspec.match(/^([A-Za-z]\w+)([*!]?)$/);
      if (m) {
        const [, fname, reqd] = m;
        let ftype = entityDef[fspec];
        let fk;

        Entity._allFields[fname] = ftype;

        if (reqd) {
          Entity._requiredFields.push(fname);
        }
        if (reqd == '*') {
          Entity._pkFields.push(fname);
        }

        if (typeof ftype == 'string') {
          m = ftype.match(/^(\w+)(?:(=)|\.(\w+))?$/);
          let [, cls, equ, fld] = m;
          if (equ) {
            fld = fname;
          }
          const FkTable = model[cls];
          if (!FkTable) {
            throw new Error("No table class for " + cls);
          }
          const FkEntity = FkTable.constructor.Entity;
          if (!FkEntity) {
            throw new Error("No entity class for " + cls);
          }
          const fkType = FkEntity._ftypes[fld];
          ftype = fkType;
          fk = {tname: cls, fname: fld};
        }

        Entity._ftypes[fname] = ftype;

        const imp = TYPE_HELPERS[ftype].imp;

        if (fk) {
          Object.defineProperty(Entity.prototype, fname, {
            enumerable: true,
            get() {
              return this._fields[fname];
            },
          });
          const fkfname = fk.tname + (
            fk.fname == fname ? '' : 'Via' + CamelCase(fname));
          Object.defineProperty(Entity.prototype, fkfname, {
            enumerable: true,
            get() {
              return this._table.model[fk.tname].lookup(
                {[fk.fname]: this[fname]});
            },
            set(value) {
              this._fields[fname] = value[fk.fname];
            },
          });
          Entity._assignableFields.push(fkfname);
        } else {
          Object.defineProperty(Entity.prototype, fname, {
            enumerable: true,
            get() {
              return this._fields[fname];
            },
            set(value) {
              this._fields[fname] = imp(value);
            },
          });
          Entity._assignableFields.push(fname);
        }
      }
    }

    Entity._assignableFields.sort();
    Entity._requiredFields.sort();
    Entity._pkFields.sort();

    if (Entity._pkFields.length) {
      Object.defineProperties(Entity.prototype, {
        _pk: {
          get() {
            return Entity._pkOf(this);
          },
        },
      });
    }

    return Entity;
  };

  function createTableClass(tname, def, Model, model) {
    let Table = Model._tableClasses[tname];
    if (Table) {
      return Table;
    }

    Table = class extends View {
      constructor(model) {
        super();
        this.model = model;
        this._data = [];
        this._class = Table;
        if (Table.Entity._pkFields.length) {
          this._pkIndex = {};
        }
      }

      get length() {
        return this._data.length;
      }

      * [Symbol.iterator]() {
        for (const e of this._data) {
          yield e;
        }
      }

      insert(data) {
        const entity = new Table.Entity(this, data)._validate();

        if (this._pkIndex) {
          const pkey = entity._pk;
          if (pkey in this._pkIndex) {
            throw `Primary key violation inserting ${tname} with key = ${pkey}`;
          }
          this._pkIndex[pkey] = entity;
        }

        this._data.push(entity);
        return entity;
      }

      toObject() {
        return this._data.map(e => e._toObject());
      }

      fromObject(array) {
        for (const e of array) {
          this.insert(e)
        }
        return this;
      }

      toJSON(replacer, space) {
        return JSON.stringify(this.toObject(), replacer, space);
      }

      fromJSON(json, reviver) {
        return this.fromObject(JSON.parse(json, reviver));
      }
    }

    Model._tableClasses[tname] = Table;

    Table.Entity = createEntityClass(tname, def, Model, model);

    if (Table.Entity._pkFields.length) {
      Object.defineProperties(Table.prototype, {
        lookup: {
          value: function(fields) {
            return this._pkIndex[Table.Entity._pkOf(fields)];
          }
        },
      });
    }

    return Table;
  };

  const modelClasses = {};

  const createModelClass = def => {
    let Model = modelClasses[def.name];
    if (Model != null) {
      return Model;
    }
    Model = class {
      constructor(input) {
        this._def = def;
        this._tables = {};

        if (typeof input == 'string') {
          this.fromJSON(input);
        } else if (input != null) {
          this.fromObject(input);
        }
      }

      _has(tname) {
        return this._tables[tname] != null;
      }

      toObject() {
        const object = {};
        for (const k in Model.prototype) {
          if (this._has(k) && this[k].length) {
            object[k] = this[k].toObject();
          }
        }
        return object;
      };

      fromObject(object) {
        for (const k in Model.prototype) {
          if (k in object) {
            this[k].fromObject(object[k]);
          }
        }
        return this;
      };

      toJSON(replacer, space) {
        return JSON.stringify(this.toObject(), replacer, space);
      };

      fromJSON(json, reviver) {
        return this.fromObject(JSON.parse(json, reviver));
      };
    }

    Model._tableClasses = {};

    for (const tname in def.types) {
      const t = def.types[tname];
      if (t['_'].rel) {
        Object.defineProperty(Model.prototype, tname, {
          enumerable: true,
          get() {
            let table = this._tables[tname];
            if (table == null) {
              table = this._tables[tname] =
                new (createTableClass(tname, t, Model, this))(this);
            }
            return table;
          },
        });
      }
    }

    return modelClasses[def.name] = Model;
  };

  return {
    version: sysl_js_version,
    createModelClass(json) {
      return createModelClass(json);
    }
  };
});
