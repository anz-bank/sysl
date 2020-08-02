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

resource "google_data_catalog_entry_group" "entry_group" {
  entry_group_id = "Source_group"
}

resource "google_data_catalog_entry" "bar_entry" {
  entry_group = google_data_catalog_entry_group.entry_group.id
  entry_id = "Source.Bar"

  user_specified_type = "Bar"
  user_specified_system = "Source"
  linked_resource = "model.sysl:14"

  display_name = "Bar"
  description  = "A database.
 Stores data."

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
      "column": "x",
      "description": "A foreign key",
      "mode": "REQUIRED",
      "type": "?"
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
