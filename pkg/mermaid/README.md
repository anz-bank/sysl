**Mermaid diagram generation design doc**

**Objective**

- Explore alternatives for PlantUML for diagram generation.

**Background**

Generating diagram using PlantUML requires making HTTP requests to the PlantUML server. In order to support systems generating several diagrams at once, it is not advisable to rely on an external service which could also affect performance and reliability. This issue also impedes our ability to integrate diagram generation as a part of CI.

**Summary**

Mermaid diagram generator offers users a way to generate diagrams similar to that of what PlantUML diagram generator offers. The four types of diagrams that can be generated are *sequence diagrams*, *integration diagrams*, *endpoint analysis diagrams* and *data model diagrams*. The sysl modules produced from parsing a sysl file contains all the information to generate a mermaid file. The mermaid diagram generator iterates through and prints the appropriate mermaid code. The data model diagram generator makes use of the sysl wrapper that simplifies the sysl modules to maps of simplified types for easy access to defined types.

**Detailed Design**

**Sequence Diagram** (https://github.com/anz-bank/sysl/pull/756)

GenerateSequenceDiagram function takes in a sysl module, an app name and an endpoint as input and returns a string that contains the mermaid code for the sequence diagram. The sysl statements for that app and endpoint is retrieved and iterated through in the function printSequenceDiagramStatements. If the statement is a group or a condition, this function is recursively called. If the statement is a call to another application and endpoint, the function GenerateSequenceDiagram is called recursively and the appropriate mermaid code is generated and stored. For returns, actions and loops, the appropriate mermaid output is generated and stored.

**Integration Diagram** (https://github.com/anz-bank/sysl/pull/787)

Integration diagrams can be generated via three different functions.

- `GenerateFullIntegrationDiagram` takes in the sysl module, iterates through all its applications and stores all the calls that each of the applications makes to other applications and generates and stores the appropriate mermaid output for that relationship.
- `GenerateIntegrationDiagram` takes in an application and does a similar process, but finds all the relationships for that application only.
- `GenerateMultipleAppIntegrationDiagram` takes in a list of applications and follows the same process.

**Endpoint Analysis Diagram** (https://github.com/anz-bank/sysl/pull/803)

The endpoint analysis diagram generator is an extension of integration diagram generator, but it adds the concept of a subgraph from mermaid in order to display the endpoints of each of the applications in the generated diagram.

**Data-model Diagram** (https://github.com/anz-bank/sysl/pull/838)

The data-model diagram generator uses the sysl wrapper library to simplify the process of iterating through types defined in different applications. There are two functions for generating diagrams here.

- `GenerateFullDataDiagram` takes a sysl module and uses the sysl wrapper library to generating a map of strings to Type. The types in these maps can be the either of the following – **tuple, relation, enum, map or empty**. Tuples and relations consist of primitive types, lists of other types or references to other types. If a type refers to or is a list of another type, this relationship is stored and printed in the diagram. Maps take a similar approach as well. Enums do not contain any references, thus appropriate mermaid output is printed directly.
- `GetDataDiagramWithAppAndType` takes a sysl module, along with an app and a type. This runs in a similar fashion to the previous function but it runs through it twice. The first iteration captures all the relationships that exists in that module and this result is used to find all the relationship that the given app and types has to other types. The resulting types are run through in the second iteration and the appropriate mermaid code is generated along with the relationships between the printed types. This uses the class diagram concept in mermaid.

**Command Line Parameters**

-     sysl diagram --help

Generate mermaid diagrams

Optional flags:

-     --help                  Show context-sensitive help (also try --help-long and --help-man).
-     --version               Show application version.
-     --log="warn"            log level: [off,debug,info,warn,trace]
-     -v, --verbose           enable verbose logging
-     --root=ROOT             sysl root directory for input model file. If root is not found, the module directory becomes the root, but the module can not import with absolute paths (or imports must be
                              relative).
-     -i, --integration           Generate an integration diagram (Optional- specify --app)
-     -s, --sequence              Generate a sequence diagram (Specify --app and --endpoint)
-     -p, --endpointanalysis      Generate an integration diagram with its endpoints (Optional- specify --app)
-     -d, --data                  Generate a Data model diagram
-     -a, --app=APP               Application
-     -e, --endpoint=ENDPOINT     Endpoint
-     -o, --output="diagram.svg"  Output file (Default: diagram.svg)


**Alternatives Considered**

- The current approach in sequence, integration and endpoint analysis uses recursion which makes the resulting string non deterministic which makes it hard to assert the resulting mermaid code. Though it doesn't change the resulting diagram, it would be nice to have a map that hold all the information and iterate through it so that the resulting code is in order.

**Enhancements/Future works**

- Enhancement on this would be to add a logger to detect if there are any nil references for bad inputs.
- Escape special characters instead of replacing them in mermaid code. (https://github.com/mermaid-js/mermaid/issues/170)
