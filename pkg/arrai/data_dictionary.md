# Data Dictionary

Generates a data dictionary as a CSV conforming to the Data Governance Metadata Template for `model.sysl`.

```bash
$ arrai run data_dictionary.arrai
Action, Platform Name, Technology Name, Physical Attribute Name, Physical Object Type, Physical Object Container Name, Physical Object Name, Physical Object Description, Business Term, Attribute Business Description
?, ?, ?, x, Table, Source, Bar, "A foreign key", ?, "?"
?, ?, ?, b, Table, Source, Bar, "An optional int", ?, "?"
?, ?, ?, x, Table, Source, Foo, "The x value.", ?, "?"
?, ?, ?, a, Table, Source, Bar, "A bar table.", ?, "?"
?, ?, ?, y, Table, Source, Foo, "A Foo.
 Represents foo things.", ?, "?"
```
