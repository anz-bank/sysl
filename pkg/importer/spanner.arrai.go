package importer

const importSpannerScript = `
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
table_constraint    -> "CONSTRAINT" constraint_name "FOREIGN KEY" "(" column_name:"," ")" "REFERENCES" foreign=(table_name "(" column_name:"," ")");
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

options_def         -> "OPTIONS ("(option):",",? ")";
option              -> "allow_commit_timestamp="("true"|"false");
key_part            -> (column_name sort_by=("ASC"|"DESC")?):",",?;
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
column_name         -> [$@A-Za-z_][0-9$@A-Za-z_]*;
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
        constraint_name: parsed.constraint_name.'' rank (@: .@),
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

let parseKeyPart = \t
     t.column_name => \(@:i, @item:n)
        cond {
            (t.sort_by?(i)?:false): $'${n.''}(${//str.lower(t.sort_by(i).'')})',
            _:  $'${n.''}',
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
                (val.unique?:false): 'true',
                _: 'false'
            },
            nullfiltered: cond {
                (val.nullfiltered?:false): 'true',
                _: 'false'
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

# Transforms that generate sysl from arrai based sql model.
# size returns the size of an attribute.
let size = \length
    cond {
        length = 'MAX': '',
        //seq.has_prefix('0x', length): $`+"`"+`(${//seq.trim_prefix('0x', length)})`+"`"+`,
        length > 0: $`+"`"+`(${length})`+"`"+`,
    };

# hasAttributePattern checks whether an attribute has associated patterns.
let hasAttributePattern = \entity \attr
    attr.options?:{} || (entity.primary_key where //seq.contains(attr.name, .)) ||
    attr.length = 'MAX' || //seq.has_prefix('0x', attr.length) ||
    ((entity.foreign_keys => ((.foreign_keys where .attribute = attr.name)
        => .reference_attribute orderby .)) where .);

# hasEntityPattern checks whether an entity has associated patterns.
let hasEntityPattern = \entity \model
    entity.primary_key | entity.foreign_keys | entity.cluster | model.indexes;

# sortingOrder determines and appends sorting order.
let sortingOrder = \e \attr
    let re = //re.compile($`+"`"+`${attr.name}\((asc|desc)\)`+"`"+`);
    let keyOrder = e.primary_key => (re.match(.)(0)?(1)?:{} rank (:.@)) where .;
    cond keyOrder {{x}: $`+"`"+`, ~${x}`+"`"+`};

# compareColumnOrder compares the col order between primary key and table.
let compareColumnOrder = \entity
    let pk = (entity.primary_key orderby .) >> //seq.split('(', .) >> .(0);
    pk !(<=) (entity.attributes >> .name);

# attributePatterns generates the patterns for an attribute.
let attributePatterns = \entity \attr
    let options = cond {
        attr.options:
            let opt = //seq.split('=',attr.options);
            $`+"`"+`${opt(0)}='${opt(1)}'`+"`"+`,
        };
    let pk = cond {
        (entity.primary_key where //seq.contains(attr.name, .)):
            '~pk' ++ sortingOrder(entity, attr)
    };
    let fk = cond {
        (entity.foreign_keys => ((.foreign_keys where .attribute = attr.name) =>
        .reference_attribute orderby .)) where .:
            '~fk'
    };
    let length = cond {
        attr.length = 'MAX': '~max',
    };
    let hexPrefix = cond {
        //seq.has_prefix('0x', attr.length): '~hex',
    };
    //seq.join(', ', ([options, pk, fk, length, hexPrefix] where .@item != ''));

# entityPatterns generates the patterns for an entity.
let entityPatterns = \entity \model
    let pk = cond {
        entity.primary_key count > 1 && compareColumnOrder(entity):
            $`+"`"+`primary_key="${entity.primary_key orderby .::, }"`+"`"+`,
    };
    let cluster = cond {
        entity.cluster:
            //seq.join('', entity.cluster >>
            $`+"`"+`interleave_in_parent="${.interleaved_in}", interleave_on_delete="${.on_delete}"`+"`"+`),
    };
    let fk = cond {
        entity.foreign_keys:
            $`+"`"+`foreign_keys=[${//seq.join(', ',(entity.foreign_keys =>
                \keys $`+"`"+`["constraint:${keys.constraint_name}","columns:${//seq.join(', ', (keys.foreign_keys => .attribute) orderby .)}"]`+"`"+`) orderby .)}]`+"`"+`,
    };
    let indices = cond {
        (model.indexes where .table_name = entity.name):
            $`+"`"+`indexes=[${//seq.join(', ', ((model.indexes orderby .) where .@item.table_name = entity.name) >>
                $`+"`"+`"name:${.table_name}","unique:${.unique}","null_filtered:${.nullfiltered}","key_parts:${//seq.join(', ', .key_part orderby .)}${cond {.storing_col: $`+"`"+`","storing:${//seq.join(', ',.storing_col)}`+"`"+`}}${cond {.interleaved_table: $`+"`"+`","interleave_in:${.interleaved_table}`+"`"+`}}`+"`"+`)}"]`+"`"+`,
    };
    //seq.join(', ', ([pk, cluster, fk, indices] where .@item != ''));

# appendEntityPatterns appends the entity patterns together.
let appendEntityPatterns = \entity \model
    cond {
        hasEntityPattern(entity, model): $`+"`"+`[${entityPatterns(entity, model)}]`+"`"+`
    };

# appendAttributePatterns appends the attribute patterns together.
let appendAttributePatterns = \entity \attr
    cond {
        hasAttributePattern(entity,attr): $`+"`"+`[${attributePatterns(entity, attr)}]`+"`"+`
    };

# typeInfo generates the type info for an attribute.
let typeInfo = \entity \attr \type \isArray
    let fkAttr = entity.foreign_keys => ((.foreign_keys where .attribute = attr.name) =>
        .reference_attribute orderby .) where .;
    let fkTable = entity.foreign_keys => ((.foreign_keys where .attribute = attr.name) =>
        .reference_table orderby .) where .;
    cond {
        fkAttr | fkTable:  //seq.join('', (fkTable => //seq.join('', .) orderby .)) ++ '.'
            ++ //seq.join('', (fkAttr => //seq.join('', .) orderby .)),
        isArray: $`+"`"+`sequence of ${type}`+"`"+`,
        _: type,
    };

# transformModel translates the empty model into sysl file.
let transformModel = \model \package
    # sysl specification
    # https://github.com/anz-bank/sysl/blob/master/pkg/sysl/sysl.proto
    $`+"`"+`
        ########   THIS IS AUTOGENERATED BY sysl   ########
        ${model.schema => $`+"`"+`
            ${.name}${cond {package: $`+"`"+`[package="${package}"]`+"`"+`} }:
            ${//seq.sub('#>>>\n', '')($`+"`"+`
            ${(model.entities => \entity $`+"`"+`
            #>>>
                !table ${entity.name} ${appendEntityPatterns(entity, model)}:
                    ${entity.attributes >>
                        $`+"`"+`
                            ${.name} <: ${typeInfo(entity, ., .type ++ size(.length), .array)}${cond {.nullable:'?'}} ${appendAttributePatterns(entity, .)}`+"`"+` ::\i:
                        }`+"`"+`
                    )
                    orderby .::\i:
            }
            `+"`"+`)}
        `+"`"+` orderby .::}
    `+"`"+`;

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
