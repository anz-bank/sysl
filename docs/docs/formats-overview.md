---
id: formats-overview
title: Supported Formats
sidebar_label: Overview
---

import useBaseUrl from '@docusaurus/useBaseUrl';

Sysl is designed to be the universal modeling language for system interaction.
With this goal in mind, Sysl supports the importing of other formal specifications for the purpose of modelling.

# Example

Imagine building the `MyMusic` service that interacted with the following two systems:

<table>
<tr><td><b>OpenAPI</b></td><td><b>Spanner</b></td></tr>
<tr valign="top">
<td>

```yaml
openapi: "3.0.0"
paths:
  /artist/top_grossing:
    get:
      responses:
        "200":
          content:
            application/json:
              schema:
                type: array
                items:
                $ref: "#/components/schemas/Artist"
components:
  schemas:
    Artist:
      type: object
      properties:
        name:
          type: string
        earnings:
          type: int64
      required:
        - name
        - earnings
```

`searchr.yaml`

</td>
<td>

```sql
CREATE DATABASE Music
CREATE TABLE Artist (
  Id   INT64 NOT NULL,
  Name STRING(1024) NOT NULL,
  Age  INT64 NOT NULL
) PRIMARY KEY (Id);
```

`music.sql`

</td>
</tr>
</table>

The `MyMusic` service contains the `/artists` endpoint that returns information on the top grossing artists:

```sysl
import searchr.yaml as searchr.Searchr ~openapi
import music.sql ~spanner

MyMusic:

  /artists:
    GET
      Searchr <- GET /artist/top_grossing
      lookup artists in database
      transform artists
      return sequence of Artist

  !view Transform(earn <: Searchr.Artist, artist <: Music.Artist) -> Artist:
    ...

  !type Artist:
    id <: int64
    name <: string
    age <: int64
    earnings <: int64
```

The example above shows how we can import other formal specifications directly into Sysl.

For more information on importing, see the individual formats below.

# Supported Formats

Sysl supports the importing of the following formal specifications:

- [Avro](#formats-avro)
- [OpenAPI](#formats-openapi)
- [Spanner](#formats-spanner)
- [XSD](#formats-xsd)
