# Data Dictionary

Generates a data dictionary as a CSV conforming to the Data Governance Metadata Template for `model.sysl`.

```bash
$ arrai run data_dictionary.arrai
Action, Platform Name, Technology Name, Physical Attribute Name, Physical Object Type, Physical Object Container Name, Physical Object Name, Physical Object Description, Business Term, Attribute Business Description
action, plat, tech, x, objType, Source, Foo, objDesc, term, attrDesc
action, plat, tech, y, objType, Source, Foo, objDesc, term, attrDesc
```
