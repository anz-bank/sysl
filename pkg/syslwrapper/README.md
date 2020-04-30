# syslwrapper

This package adds an abstraction layer around sysl proto types defined in the sysl package.
This is currently a work in progress, but is intended to make working with sysl types much easier and to solve common issues around working with sysl type references, as well as parsing return statements.

## Usage

```go
func main() {
	filename := "demo/grocerystore/GroceryStore.sysl"
	syslModules, err := parse.NewParser().Parse(filename, afero.NewOsFs())
	if err != nil {
		panic(err)
	}
   
	mapper := MakeAppMapper(syslModules)
    mapper.IndexTypes() // This indexes all the sysl.Type in each application in a map, which can be accessed using "AppName:TypeName"
    mapper.ConvertTypes() // This converts the indexed types into syslwrapper.Type
    simpleApps, err := mapper.Map() // Maps the applications in syslModules into simple syslwrapperApps
    
    getInventoryEndpoint := simpleApps.["GroceryStore"].Endpoints["GET /inventory"]
    params := getInventoryEndpoint.Params
    responses := getInventoryEndpoint.Responses
    
    prettyPrintParams, err := json.MarshalIndent(params, "", " ")
	prettyPrintResponses, err := json.MarshalIndent(responses, "", " ")
	prettyPrintInventoryResponse, err := json.MarshalIndent(mapper.SimpleTypes["GroceryStore:fooid"], "", " ")
    fmt.Println(string(prettyPrintParams))
	fmt.Println(string(prettyPrintResponses))
	fmt.Println(string(prettyPrintInventoryResponse))
}
```

### Response Statement Parsing

Responses are returned as
```json
{
 "200": {
  "In": "",
  "Description": "",
  "Name": "200",
  "Type": {
   "Description": "",
   "Optional": false,
   "Reference": "",
   "Type": "list",
   "Items": [
    {
     "Description": "",
     "Optional": false,
     "Reference": "GroceryStore:InventoryResponse",
     "Type": "ref",
     "Items": null,
     "Enum": null,
     "Properties": null
    }
   ],
   "Enum": null,
   "Properties": null
  }
 }
}
```
Let's compare this with what the `sysl.Application` type gives us if we called `syslModules.Apps["GroceryStore"].Endpoints["GET /inventory"].Stmt`
```
[
 {
  "Stmt": {
   "Call": {
    "target": {
     "part": [
      "Inventory"
     ]
    },
    "endpoint": "GET /inventory"
   }
  }
 },
 {
  "Stmt": {
   "Ret": {
    "payload": "sequence of InventoryResponse"
   }
  }
 }
]
```

### Parameter Type Conversion

Here's what `simpleApps["GroceryStore"].Endpoints["GET /inventory"].Params` gives us

```json
{
 "fooid": {
  "In": "header",
  "Description": "",
  "Name": "fooid",
  "Type": {
   "Description": "",
   "Optional": false,
   "Reference": "",
   "Type": "string",
   "Items": null,
   "Enum": null,
   "Properties": null
  }
 }
}
```

Let's compare this with what the `sysl.Application` type gives us if we called `syslModules.Apps["GroceryStore"].Endpoints["GET /inventory"].Param`
```json
[
 {
  "name": "fooid",
  "type": {
   "Type": {
    "Primitive": 6
   },
   "attrs": {
    "name": {
     "Attribute": {
      "S": "FooID"
     }
    },
    "patterns": {
     "Attribute": {
      "A": {
       "elt": [
        {
         "Attribute": {
          "S": "header"
         }
        },
        {
         "Attribute": {
          "S": "required"
         }
        }
       ]
      }
     }
    }
   }
  }
 }
]
```

### Type Lookups

If we wanted to lookup what InventoryResponse was, we can just call

```go
    mapper.SimpleTypes["GroceryStore:InventoryResponse"]
```

which returns

```json
{
 "Description": "",
 "Optional": false,
 "Reference": "",
 "Type": "tuple",
 "Items": null,
 "Enum": null,
 "Properties": {
  "item_id": {
   "Description": "",
   "Optional": false,
   "Reference": "",
   "Type": "string",
   "Items": null,
   "Enum": null,
   "Properties": null
  },
  "quantity": {
   "Description": "",
   "Optional": false,
   "Reference": "",
   "Type": "int",
   "Items": null,
   "Enum": null,
   "Properties": null
  }
 }
}
```

### Type Resolution (WIP)

Lets say you don't care for having to deal with looking up type references, you just want definite types for all of these

`mapper.ResolveTypes()`

Or if you want to just resolve one type, use `mapper.ResolveTypes()` Not Yet Fully Implemented

This function goes through each type reference and replaces it with the actual data type. This converts the data types from a graph representation to a tree representation.

This feature is still under construction, with some recursion that needs to be resolved.

Example
```go
func main() {
	filename := "demo/grocerystore/GroceryStore.sysl"
	syslModules, err := parse.NewParser().Parse(filename, afero.NewOsFs())
	if err != nil {
		panic(err)
	}

	mapper := syslwrapper.MakeAppMapper(syslModules)
	mapper.IndexTypes()
	mapper.ResolveTypes()
	mapper.ConvertTypes()
	simpleApps, err := mapper.Map()

	responses := simpleApps["GroceryStore"].Endpoints["GET /inventory"].Response
	prettyPrintResponses, err := json.MarshalIndent(responses, "", " ")
	fmt.Println(string(prettyPrintResponses))
}
```

```json
{
 "200": {
  "In": "",
  "Description": "",
  "Name": "200",
  "Type": {
   "Description": "",
   "Optional": false,
   "Reference": "",
   "Type": "list",
   "Items": [
    {
     "Description": "",
     "Optional": false,
     "Reference": "GroceryStore:InventoryResponse",
     "Type": "ref",
     "Items": null,
     "Enum": null,
     "Properties": null
    }
   ],
   "Enum": null,
   "Properties": null
  }
 }
}
```
