# Terraform Google Data Catalog

Generates a Terraform file describing the databases in a Sysl model as resources representing [Google Data Catalog entries](https://www.terraform.io/docs/providers/google/r/data_catalog_entry.html#example-usage-data-catalog-entry-full).

```bash
$ arrai run terraform_gdc.arrai
resource "google_data_catalog_entry" "foo_entry" {
  entry_group = google_data_catalog_entry_group.entry_group.id
  entry_id = "Source.Foo"

  user_specified_type = "Foo"
  user_specified_system = "Source"
  linked_resource = "model.sysl:6"

  display_name = "Foo"
  description  = "A database.
 Stores data."

  schema = <<EOF
{
  "columns": [
    {
      "column": "y",
      "description": "A Foo.
 Represents foo things.",
      "mode": "REQUIRED",
      "type": "INT"
    },
    {
      "column": "x",
      "description": "The x value.",
      "mode": "REQUIRED",
      "type": "INT"
    }
  ]
}
EOF
}

[...]
```
