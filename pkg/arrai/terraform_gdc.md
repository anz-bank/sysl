# Terraform Google Data Catalog

Generates a Terraform file describing the databases in a Sysl model as resources representing [Google Data Catalog entries](https://www.terraform.io/docs/providers/google/r/data_catalog_entry.html#example-usage-data-catalog-entry-full).

```bash
$ arrai run terraform_gdc.arrai
resource "google_data_catalog_entry" "bar_entry" {
  entry_group = google_data_catalog_entry_group.entry_group.id
  entry_id = "Source.Bar"

  user_specified_type = "Bar"
  user_specified_system = "Source"
  linked_resource = "model.sysl:14"

  display_name = "Bar"
  description  = "A database. Stores data."

  schema = <<EOF
{
  "columns": [
    {
      "column": "a",
      "description": "A bar table.",
      "mode": "REQUIRED",
      "type": "STRING"
    },
    {
      "column": "b",
      "description": "An optional int",
      "mode": "NULLABLE",
      "type": "INT"
    }
  ]
}
EOF
}

resource "google_data_catalog_entry_group" "entry_group" {
  entry_group_id = "Source_group"
}
```
