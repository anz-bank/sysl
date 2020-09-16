// nolint
package importer

const importSpannerScript = `
### ------------------------------------------------------------------------ ###
###  ../../arrai/util.arrai                                                  ###
### ------------------------------------------------------------------------ ###

# A collection of helper functions for arr.ai.
#
# If generally useful, these should gradually migrate to a more standard library.

# Invokes a macro on a string as if it were source code at parsing time.
let invokeMacro = \macro \s
    macro -> (//dict(.@transform) >>> \rule \fn
        fn(//grammar.parse(.@grammar, rule, s))).@value
;

# Transforms an AST into a simple tuple of its values.
# Useful for the @transform of a flat grammar.
let rec simpleTransform = \ast
    cond ast {
        (...):
            let d = //dict(ast) >> \term cond term {
                ('':value): value,
                (...): simpleTransform(term),
                [...]: term >> simpleTransform(.) where .@item,
                _: {},
            };
            //tuple(d where .@value),
        _: {}
    }
;

# Filters the nodes of a hierarchical data structure based on a (key, value) predicate.
# Key-value pairs for which the predicate returns false will be removed from the result.
let rec filterTree = \pred \ast
    cond ast {
        {(@:..., @value:...), ...}: ast where pred(.@, .@value) >> filterTree(pred, .),
        [...]: ast >> filterTree(pred, .),
        {...}: ast => filterTree(pred, .),
        (...): safetuple(//dict(ast) where pred(.@, .@value) >> filterTree(pred, .)),
        _: ast,
    }
;

# Sequentially applies ` + "`" + `fn(accumulator, i)` + "`" + ` for each ` + "`" + `i` + "`" + ` in ` + "`" + `arr` + "`" + `. The ` + "`" + `accumulator` + "`" + ` is initialised
# to ` + "`" + `val` + "`" + `, and updated to the result of ` + "`" + `fn` + "`" + ` after each invocation.
# Returns the final accumulated value.
let rec reduce = \arr \fn \val
    cond arr {
        [head, ...]:
            let tail = -1\(arr without (@:0, @item:head));
            reduce(tail, fn, fn(val, head)),
        _: val,
    }
;

# Sequentially applies ` + "`" + `fn(accumulator, k, v)` + "`" + ` for each ` + "`" + `{k: v}` + "`" + ` pair in ` + "`" + `obj` + "`" + `.
# The ` + "`" + `accumulator` + "`" + ` is initialised to ` + "`" + `val` + "`" + `, and updated to the result of ` + "`" + `fn` + "`" + `
# after each invocation.
# Returns the final accumulated value.
let rec reduceObj = \obj \fn \val
    cond obj {
        {}: val,
        [(@:k, @value:v), ...tail]: reduceObj(tail rank (:.@), fn, fn(val, k, v)),
        [(@:k, @item:v), ...tail]:  reduceObj(tail rank (:.@), fn, fn(val, k, v)),
        (...): reduceObj(//dict(obj) orderby .@, fn, val),
        _:     reduceObj(obj orderby .@, fn, val),
    }
;

# Performs ` + "`" + `reduce` + "`" + ` once on ` + "`" + `arr` + "`" + `, and once for each array output of ` + "`" + `fn` + "`" + `. Accumulates to the same
# value across all invocations.
let reduceFlat = \arr \fn \val
    reduce(arr, \z \i reduce(i, fn, z), val)
;

# Returns a sequence with any offset and holes removed.
let ranked = \s s rank (:.@);
# Explore constructs a dependency graph by starting at source and calling step
# to find adjacent nodes. Deps is the graph constructed so far.
# Self-edges are ignored.
let rec _explore = \source \step \deps
    cond {
        {source} & (deps => .@): deps,
        _:
            let next = step(source) where . != source;
            let deps = deps | {(@:source, @value: next)};
            reduce(next orderby ., \v \i _explore(i, step, v), deps)
    };
let explore = \source \step _explore(source, step, {});

# Unimported returns the set of nodes with no in-edges.
let unimported = \g (g => .@) where !({.} & //rel.union(g => .@value));

# Topsort returns an array of nodes in graph in dependency order.
let rec _topsort = \graph \sorted \sources
    cond sources orderby . {
        []: sorted,
        [..., tail]:
            let adjs = graph(tail);
            let graph = graph where .@ != tail;
            let sources = (sources &~ {tail}) | (adjs & unimported(graph));
            _topsort(graph, sorted ++ [tail], sources)
    };
let topsort = \graph _topsort(graph, [], unimported(graph));

# TODO: this should be part of stdlib
let rec trimWhitespace = \str
    let prefix = //seq.trim_prefix(' ');
    let suffix = //seq.trim_suffix(' ');
    let trimmed = prefix(suffix(str));
    cond trimmed {
        (str): str,
        _: trimWhitespace(trimmed)
    }
;


# TODO: Handle type without app reference
let typeGrammar = {:
        //grammar.lang.wbnf[grammar]:
        types -> (app=([^\.]+) ".")? type=([^\.]+):".";
    :};
let unpackType = \type (
    cond type {
        (''): (app: '', type: '', field: ''),
        _: (//grammar -> .parse(typeGrammar, 'types', type))
            ->
            # TODO: remove once .field?: is fixed
            let t = .;
            let app = trimWhitespace(ranked(t.app?.''?:''));
            let typeCount = .type count;
            cond (typeCount) {
                (1): (
                    :app,
                    type : .type >> trimWhitespace(ranked(.'')),
                    field: ''
                ),
                _: (
                    :app,
                    type : .type where .@ != typeCount - 1 >> trimWhitespace(ranked(.'')),
                    field: trimWhitespace(ranked(.type(typeCount-1).''))
                )
            }
    }
);

let packType = \(app: appName, type: typeSeq, field: fieldName) (
    cond fieldName {
        (''): //seq.join('.', [appName] ++ typeSeq),
        _: //seq.join('.', [appName] ++ typeSeq ++ [fieldName]),
    }
)
;

# isValidIdentifier checks whether the identifier name is valid.
let isValidIdentifier = \identifier
    # InvalidIdentifiers that would be appended by underscore('_') when used as an identifier in the ingested SQL
    # for instance  a column "Int64 INT64" becomes _Int64 <: int [name="Int64"].
    # List is still fairly limited but more keywords could be added as we go.
    let invalidIdentifiers = { "any", "as", "bool", "bytes", "date", "datetime", "decimal",
    "else", "float", "float64", "if", "int", "int64", "string" };

    # sysl is largely case insensitive so lowercase the identifier before comparison
    # taken from pkg/grammar/SyslLexer.g4
    let regex = //re.compile("('%'[0-9a-fA-F][0-9a-fA-F])*[a-zA-Z_]([-a-zA-Z0-9_]|('%'[0-9a-fA-F][0-9a-fA-F]))*");
    !(//str.lower(identifier) <: invalidIdentifiers) && regex.match(identifier);

# resolveValidIdentifier resolves the invalid identifier name.
let resolveValidIdentifier = \identifier
    cond {
        !isValidIdentifier(identifier): '_' ++ identifier,
        _: identifier
    };

let util_arrai =
(
    :explore,
    :filterTree,
    :invokeMacro,
    :simpleTransform,
    :ranked,
    :reduce,
    :reduceFlat,
    :reduceObj,
    :ranked,
    :simpleTransform,
    :topsort,
    :unimported,
    :unpackType,
    :packType,
    :trimWhitespace,
    :isValidIdentifier,
    :resolveValidIdentifier,
);

### ------------------------------------------------------------------------ ###
###  spanner.arrai                                                           ###
### ------------------------------------------------------------------------ ###

# spanner ddl grammar
# CREATE DATABASE statements are parsed to avoid errors but ignore by the modelling functions
let grammar = {://grammar.lang.wbnf[grammar]:
ddl                 -> stmt=(create_database|create_table|create_index|alter_table|drop_table|drop_index):";" ";" \s*;

create_database     -> "CREATE DATABASE" database=([a-z][0-9a-z_]*[0-9a-z]);

create_table        -> "CREATE TABLE" table_name "("(
                              attr=(column_name attr_type not_null? options_def?)
                            | table_constraint
                        ):"," "," ")"
                       pk=("PRIMARY KEY" "(" key_part ")")
                       ("," cluster)*;
not_null            -> "NOT NULL";
table_constraint    -> ("CONSTRAINT" constraint_name)? "FOREIGN KEY" "(" column_name:"," ")" "REFERENCES" foreign=(table_name "(" column_name:"," ")");
cluster             -> "INTERLEAVE IN PARENT" table_name ("ON DELETE" on_delete)?;
on_delete           -> CASCADE   = "CASCADE"
                     | NO_ACTION = "NO ACTION";

create_index        -> "CREATE" unique=("UNIQUE")? nullfiltered=("NULL_FILTERED")? "INDEX" index_name "ON" table_name "(" key_part ")" storing_clause? interleaving_clause?;
storing_clause      -> "STORING" "(" column_name:",",? ")";
interleaving_clause -> "INTERLEAVE IN" table_name;

alter_table         -> "ALTER TABLE" table_name (alter=(table_alteration|attr_alteration)):",";
table_alteration    -> ADD_COLUMN      = ("ADD COLUMN" column_name (attr_type | options_def))
                     | DROP_COLUMN     = ("DROP COLUMN" column_name)
                     | SET_COLUMN      = ("SET ON DELETE" on_delete)
                     | ADD_CONSTRAINT  = ("ADD" table_constraint)
                     | DROP_CONSTRAINT = ("DROP CONSTRAINT" constraint_name);
attr_alteration     -> "ALTER COLUMN" column_name (attr_type | "SET" options_def);

drop_table          -> "DROP TABLE" table_name;

drop_index          -> "DROP INDEX" index_name;

options_def         -> "OPTIONS" "("(option):",",? ")";
option              -> "allow_commit_timestamp" "=" ("true"|"false");
key_part            -> column_def=(column_name sort_by=("ASC"|"DESC")?):",",?;
attr_type           -> (SCALAR_TYPE|ARRAY_TYPE);

ARRAY_TYPE          -> "ARRAY<" SCALAR_TYPE ">";
SCALAR_TYPE         -> BOOL      = "BOOL"
                     | INT64     = "INT64"
                     | FLOAT64   = "FLOAT64"
                     | DATE      = "DATE"
                     | TIMESTAMP = "TIMESTAMP"
                     | STRING    = "STRING(" length ")"
                     | BYTES     = "BYTES(" length ")";

length              -> (int64_value|"MAX");

table_name          -> [$@A-Za-z_][0-9$@A-Za-z_]*;
constraint_name     -> [$@A-Za-z_][0-9$@A-Za-z_]*;
column_name         -> /{` + "`" + `[^` + "`" + `]*` + "`" + `|[$@A-Za-z_][0-9$@A-Za-z_]*};
index_name          -> [$@A-Za-z_][0-9$@A-Za-z_]*;

int64_value         -> hex_value|decimal_value;
decimal_value       -> [-]?\d+;
hex_value           -> /{-?0x[[:xdigit:]]+};

.wrapRE -> /{(?i)\s*()};
:};

#################### PARSE TREE EVALUATORS ####################
# These functions turn a spanner ddl parse tree into ddl statements

# evalInt64 turns an int64 parse node into an integer
let evalInt64 = \parsed cond parsed {
    (decimal_value: ('': val), ...): //eval.value(val),
    (hex_value: ('': val), ...): val,
    _: "MAX"
};

# evalType turns an attr_type parse node into an attribute type
let evalType = \parsed
    let data = cond parsed {
        (SCALAR_TYPE: scalar, ...):                    (scalar: scalar, array: false),
        (ARRAY_TYPE: (SCALAR_TYPE: scalar, ...), ...): (scalar: scalar, array: true),
    };
    let type = cond data.scalar {
        (BOOL: _, ...):                   (type: "bool", length: 0),
        (INT64: _, ...):                  (type: "int", length: 0),
        (FLOAT64: _, ...):                (type: "float", length: 0),
        (DATE: _, ...):                   (type: "date", length: 0),
        (TIMESTAMP: _, ...):              (type: "datetime", length: 0),
        (STRING: _, length: length, ...): (type: "string", length: evalInt64(length.int64_value?:"MAX")),
        (BYTES: _, length: length, ...):  (type: "bytes", length: evalInt64(length.int64_value?:"MAX")),
    };
    (
        type: type.type,
        length: type.length,
        array: data.array,
    );

# evalAttribute turns an attr node into an attribute of a relation
let evalAttribute = \parsed
    let type = evalType(parsed.attr_type);
    (
        name:     parsed.column_name.'' rank (:.@),
        type:     type.type,
        length:   type.length,
        array:    type.array,
        options: cond {
            (parsed.options_def?:false): //seq.join('',parsed.options_def.option.@item.''),
            _: {}
        },
        nullable: cond parsed {(not_null: _, ...): false, _: true},
    );

let evalForeignKeyConstraint = \parsed
    let reference = parsed.foreign -> (
        table_name: .table_name.'' rank (@: .@),
        attributes: .column_name >> (.'' rank (@: .@)),
    );
    (
        constraint_name: cond {
            (parsed.constraint_name?:false): parsed.constraint_name.'' rank (@: .@),
            _: {}
        },
        foreign_keys: parsed.column_name => (
            attribute: .@item.'' rank (@: .@),
            reference_table: reference.table_name,
            reference_attribute: reference.attributes(.@),
        ),
    );

let evalTableAlteration = \parsed cond parsed {
    (table_alteration: (ADD_COLUMN: data, ...), ...): (
        type: "add_column",
        alteration: (
            name: data.column_name.'' rank (@: .@),
            type: evalType(data.attr_type),
        ),
    ),

    (table_alteration: (DROP_COLUMN: data, ...), ...): (
        type: "drop_column",
        alteration: (
            name: data.column_name.'' rank (@: .@),
        ),
    ),

    (table_alteration: (SET_COLUMN: data, ...), ...): (
        type: "on_delete",
        alteration: data
    ),

    (table_alteration: (ADD_CONSTRAINT: data, ...), ...): (
        type: "add_constraint",
        alteration: evalForeignKeyConstraint(data.table_constraint),
    ),

    (table_alteration: (DROP_CONSTRAINT: data, ...), ...): (
        type: "drop_constraint",
        alteration: data,
    ),

    (attr_alteration: data, ...): (type: "alter_column", alteration: data),
};

# concatOffset appends the two strings preserving string offsets
let concatOffset = \str1 \str2
    (str1 => .@ orderby .)(0)\$` + "`" + `${str1}${str2}` + "`" + `;

# parseKeyPart parses the primary_key generated from spanner sql
let parseKeyPart = \t
    t.column_def >> \def cond {
        (def.sort_by?:false): concatOffset(def.column_name.'', $` + "`" + `(${//str.lower(def.sort_by.'')})` + "`" + `),
        _: def.column_name.'',
    };

# evalDdl turns a ddl parse tree into a list of ddl statements ready to be applied to a model
# Use applyStmt to apply these statements to a spanner model
let evalDdl = \parsed parsed.stmt >> cond . {
    (create_table: val, ...): (
        stmt: "create_table",
        data: (
            name: val.table_name.'' rank (:.@),
            attributes: val.attr >> evalAttribute(.),
            foreign_keys: cond val {
                (table_constraint: [...constraints], ...): constraints => evalForeignKeyConstraint(.@item),
            },
            primary_key: parseKeyPart(val.pk.key_part),
            cluster: cond {
                (val.cluster?:false): val.cluster >> (
                    interleaved_in: (.table_name.'' rank (@: .@)),
                    on_delete: (.on_delete.CASCADE.'' rank (@: .@))
                ),
                _: {}
            },
        ),
    ),

    (create_index: val, ...): (
        stmt: "create_index",
        data: (
            unique: cond {
                (val.unique?:false): true,
            },
            nullfiltered: cond {
                (val.nullfiltered?:false): true,
            },
            name: val.index_name.'' rank (@: .@),
            table_name: val.table_name.'' rank (@: .@),
            key_part: parseKeyPart(val.key_part),
            storing_col: cond {
                (val.storing_clause?:false): val.storing_clause.column_name >> (.'' rank (@: .@)),
                _: {}
            },
            interleaved_table: (val.interleaving_clause?.table_name.'':'') rank (@: .@),
        ),
    ),

    (create_database: val, ...): (
        stmt: "create_database",
        data: (
            name: val.database.'' rank (@: .@),
        ),
    ),

    (alter_table: val, ...): (
        stmt: "alter_table",
        data: (
            table_name: val.table_name.'' rank (@: .@),
            alterations: val.alter >> evalTableAlteration(.),
        ),
    ),

    (drop_table: val, ...): (
        stmt: "drop_table",
        data: val.table_name.'' rank (:.@),
    ),
};

################## PARSERS ##################

# parses a byte array against the ddl grammar and hands it to eval
let parseDdl = \bytes evalDdl(//grammar.parse(grammar, "ddl", bytes));

# parses a list of schema files. reads each file and hands to parseDdl
let parseSchema = \files //seq.concat(files >> parseDdl(//os.file(.)));

################## STATEMENTS #################

# applies a create table statement
let applyCreateTable = \relation \model
    let relations = cond model.entities where .name=relation.name {
        false: model.entities | {relation},
        true: false, # a match means the ddl is trying to create a table that already exists
    };
    (
        entities: relations,
        indexes: model.indexes,
        schema: model.schema,
    );

# applies an alter table stamement
# NOT IMPLEMENTED
let applyAlterTable = \alteration \model
    model;

# applies a create index statement
let applyCreateIndex = \index \model
    let indxs = cond model.indexes where .name=index.name {
        false: model.indexes | {index},
        true: false, # a match means the ddl is trying to create a index that already exists
    };
    (
        entities: model.entities,
        indexes: indxs,
        schema: model.schema,
    );

let applyCreateDatabase = \database \model
    let dbschema = cond model.schema where .name=database.name {
        false: model.schema | {database},
        true: false, # a match means the ddl is trying to create a database that already exists
    };
    (
        entities: model.entities,
        indexes: model.indexes,
        schema: dbschema,
    );

# applies a drop table statement
let applyDropTable = \name \model
    let relations = model.entities where .name != relation;
    (
        entities: relations,
        indexes: model.indexes,
        schema: model.schema,
    );

# applies a drop_index statement
# NOT IMPLEMENTED
let applyDropIndex = \name \model
    model;

# applies either a single staement or a list of statements in the given order
let rec applyStmt = \stmt \model
    cond stmt {
        # match against the single statement types
        (stmt: "create_database", data: schema):  applyCreateDatabase(schema, model),
        (stmt: "create_table", data: relation):  applyCreateTable(relation, model),
        (stmt: "create_index", data: index):     applyCreateIndex(index, model),
        (stmt: "alter_table", data: alteration): applyAlterTable(alteration, model),
        (stmt: "drop_table", data: name):        applyDropTable(name, model),
        (stmt: "drop_index", data: name):        applyDropIndex(name, model),
        (...): model, # guard against unrecognised statements

        # match against an arrai of statements and recursively apply them in order
        [first, ...rem]: applyStmt(rem, applyStmt(first, model)),
        []:              model,
    };

let spanner_arrai =
################# EXPOSE ################
(
    # empty model, use this as the base of an applyStmt call to create a model from a ddl statement
    emptyModel:        (entities: {}, indexes: {}, schema: {}),

    # parses a single byte array representing a spanner schema
    parseDdl:          parseDdl,

    # parses a list of schema files. opens the files and calls parseDdl on them in the given order
    parseSchema: parseSchema,

    # applies a ddl stmt to a model. Use emptyModel to get a model from scratch
    applyStmt:         applyStmt,
);

### ------------------------------------------------------------------------ ###
###  sysl.arrai                                                              ###
### ------------------------------------------------------------------------ ###

# Transforms that generate sysl from an arr.ai-based SQL model.

# import sysl lib.
let util = util_arrai;

# size returns the size of an attribute.
let size = \length
    cond {
        length = 'MAX': '',
        //seq.has_prefix('0x', length): $` + "`" + `(${//seq.trim_prefix('0x', length)})` + "`" + `,
        length > 0: $` + "`" + `(${length})` + "`" + `,
    };

# sortingOrder determines and appends sorting order.
let sortingOrder = \e \attr
    let re = //re.compile($` + "`" + `${attr.name}\((asc|desc)\)` + "`" + `);
    let keyOrder = e.primary_key >> (re.match(.)(0)?(1)?:{} rank (:.@)) where .;
    (keyOrder where .@item) rank (@: .@);

# compareColumnOrder compares the col order between primary key and table.
let compareColumnOrder = \entity
    let pk = entity.primary_key >> //seq.split('(', .)(0);
    pk !(<=) (entity.attributes >> .name);

# matchingFKs returns the foreign keys matching attribute name.
let matchingFKs = \entity \attr
    entity.foreign_keys => (.foreign_keys where .attribute = attr.name) where .;

# attributePatterns generates the patterns for an attribute.
let attributePatterns = \entity \attr
    let attrName = cond {
        !util.isValidIdentifier(attr.name): $` + "`" + `name="${attr.name}"` + "`" + `
    };
    let options = cond {
        attr.options:
            let [k, v, ...] = //seq.split('=', attr.options);
            $` + "`" + `${k}="${v}"` + "`" + `,
        };
    let pk = cond {
        (entity.primary_key where //seq.contains(attr.name, .@item)):
            $` + "`" + `~pk${cond {sortingOrder(entity, attr):$` + "`" + `, ~${sortingOrder(entity, attr) ::}` + "`" + `}}` + "`" + `
    };
    let fk = cond { matchingFKs(entity, attr): '~fk' };
    let length = cond { attr.length = 'MAX': '~max' };
    let hexPrefix = cond { //seq.has_prefix('0x', attr.length): '~hex' };
    let byteLength = cond {
        attr.type = 'bytes' && attr.length > 0 && attr.length != 'MAX': $` + "`" + `length="${attr.length}"` + "`" + `
    };
    [attrName, options, byteLength, pk, fk, length, hexPrefix] where .@item;

# entityPatterns generates the patterns for an entity.
let entityPatterns = \entity \model
    let pk = cond {
        entity.primary_key count > 1 && compareColumnOrder(entity):
            $` + "`" + `primary_key="${entity.primary_key ::,}"` + "`" + `,
    };
    let cluster = cond {
        entity.cluster:
            //seq.join(', ', entity.cluster >>
            $` + "`" + `interleave_in_parent="${.interleaved_in}", interleave_on_delete="${//str.lower(.on_delete)}"` + "`" + `),
    };
    let fk = cond {
        entity.foreign_keys: $` + "`" + `
            foreign_keys=[${(entity.foreign_keys => \keys $` + "`" + `
                [${cond {keys.constraint_name: $` + "`" + `"constraint:${keys.constraint_name}",` + "`" + `}}"columns:${keys.foreign_keys => .attribute orderby .::,}"]
            ` + "`" + `) orderby .::,}]
        ` + "`" + `,
    };
    let indx = model.indexes where .table_name = entity.name => ([
        $` + "`" + `"name:${.name}"` + "`" + `,
        cond {.unique: $` + "`" + `"unique:${.unique}"` + "`" + `},
        cond {.nullfiltered: $` + "`" + `"null_filtered:${.nullfiltered}"` + "`" + `},
        $` + "`" + `"key_parts:${.key_part ::,}"` + "`" + `,
        cond {.storing_col: $` + "`" + `"storing:${.storing_col::,}"` + "`" + `},
        cond {.interleaved_table: $` + "`" + `"interleave_in:${.interleaved_table}"` + "`" + `},
    ] where .@item) => '[' ++ //seq.join(',', .) ++ ']';
    [pk, cluster, fk, cond { indx: $` + "`" + `indexes=[${indx orderby . ::,}]` + "`" + ` }] where .@item;

# entityPatternsString returns the annotation for an entity's patterns.
let entityPatternsString = \entity \model cond entityPatterns(entity, model) {[]: '', ePats: $` + "`" + `[${ePats ::, }]` + "`" + `};

# attributePatternsString returns the annotation for an attribute's patterns.
let attributePatternsString = \entity \attr cond attributePatterns(entity, attr) {[]: '', aPats: $` + "`" + `[${aPats ::, }]` + "`" + `};

# typeInfo generates the type info for an attribute.
let typeInfo = \entity \attr \type \isArray
    let fks = matchingFKs(entity, attr);
    let fkAttr = cond {
        fks:  fks => .reference_attribute
    };
    let fkTable = cond {
        fks:  fks => .reference_table
    };
    cond {
        fkAttr && fkTable: $` + "`" + `${fkTable orderby . ::}.${fkAttr orderby . ::}` + "`" + `,
        isArray: $` + "`" + `sequence of ${type}` + "`" + `,
        _: type,
    };

# transformModel translates the empty model into sysl file.
let transformModel = \model \package
    # sysl specification
    # https://github.com/anz-bank/sysl/blob/master/pkg/sysl/sysl.proto
    $` + "`" + `
        ##########################################
        ##                                      ##
        ##  AUTOGENERATED CODE -- DO NOT EDIT!  ##
        ##                                      ##
        ##########################################

        ${model.schema => $` + "`" + `
            ${.name} ${cond {package: $` + "`" + `[spanner_spec="1.0", package="${package}"]` + "`" + `} }:
                ${(model.entities => \entity
                let eps = //seq.sub(", , ", ", ", entityPatternsString(entity, model)); $` + "`" + `
                    !table ${entity.name}${cond{eps: ' ' ++ eps}}:
                        ${entity.attributes >> let aps = attributePatternsString(entity, .); $` + "`" + `
                            ${util.resolveValidIdentifier(//seq.sub("` + "`" + `", "", .name))} <: ${typeInfo(entity, ., .type ++ cond {.type != 'bytes': size( .length)}, .array)}${cond {.nullable:'?'}} ${aps}` + "`" + ` ::\i:
                        }` + "`" + `
                    ) orderby .::\i:}
                ` + "`" + ` orderby .::}` + "`" + `;

let sysl_arrai =
(
    :transformModel,
);

### ------------------------------------------------------------------------ ###
###  import.arrai                                                            ###
### ------------------------------------------------------------------------ ###

let import = \importSql \appName \syslPackage
    let spanner = spanner_arrai;
    let sysl = sysl_arrai;
    let stmts = spanner.parseSchema([importSql]);
    let model = spanner.applyStmt(stmts, spanner.emptyModel);
    sysl.transformModel(
        cond {
            (model.schema): model,
            _: (
                entities: model.entities,
                indexes: model.indexes,
                schema: {(name: appName)},
            )
        }, syslPackage);

(
    :import,
)
`
