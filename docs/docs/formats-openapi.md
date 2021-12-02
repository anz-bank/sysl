---
id: formats-openapi
title: OpenAPI
sidebar_label: OpenAPI
---

import useBaseUrl from '@docusaurus/useBaseUrl';

Sysl supports the importing of [OpenAPI3](https://swagger.io/specification/) schema.

### Extensions
OpenAPI 3 allows [extensions](https://swagger.io/docs/specification/openapi-extensions/) so that users are able to add custom data to their OpenAPI spec. To accomodate this, the Sysl importer will map these extensions into Sysl annotations and place them at the appropriate levels. Currently, only extensions at these specifications are handled:

- `info`: extensions here are mapped into application level annotation
- `paths`: extensions here are mapped into endpoint level annotations.
- `operation` (i.e. `get`, `post`, etc): extensions here are mapped into HTTP operation level annotations.
- `schema`: extensions here are mapped into type level annotations
- `responses`: extensions here are mapped into the return statement annotations
- `requestBody`: extensions here are mapped into the parameters type annotations.

Because Sysl and OpenAPI have different structures, extensions at different levels in OpenAPI may be reduced to annotations at the same level in Sysl, which may conflict. In these cases, the extensions lower in the OpenAPI hierarchy take precedence. For example:

```yaml
responses:
    200:
        content:
            plain/text:
                type: string
            x-anno: 2
        x-anno: 1
```

This specification will be mapped into the following Sysl specification:

```sysl
return ok <: string [mediatype="plain/text", anno="2"]
```

`anno` is `"2"` because the type is in the lower hierarchy and it has higher precedence. This applies to everywhere such as `paths` and `operation`, `responses`, `requestBody` and any `schema` inside them.

#### Special Annotations

To retain information, the importer would sometimes maps certain information into Sysl annotations. Due to the non-restrictive nature of OpenAPI 3 extensions, this creates a problem of possible clashes between annotations since the importer would trim the `x-` prefix. To handle this, Sysl importer will retain both annotations by not trimming the `x-` prefix on certain annotation. For example:

```yaml
paths:
    /path:
        get:
            description: "this is a description"
            x-description: "this is another description"
            responses:
                200:
                    content:
                        plain/text:
                            type: string
                            x-mediatype: "a mediatype information"
```

would be mapped into

```sysl
/path:
    GET:
        @description = "this is a description"
        @x-description = "this is another description"
        return ok <: string [mediatype="plain/text", x-mediatype="a mediatype information"]
```

Currently, the following annotations will have this behaviour:

- description
- mediatype
- openapi_type
- openapi_format
- package
- examples
- patterns
