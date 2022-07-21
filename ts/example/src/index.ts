// Import the TypeScript model directly using the `sysl/model` export:
import {
    Action,
    Application,
    AppName,
    Endpoint,
    Model,
    Param,
    Primitive,
    Statement,
    Struct,
    Type,
    TypePrimitive,
} from "@anz-bank/sysl/model";

const model1 = new Model({
    apps: [
        new Application({
            name: new AppName(["App1"]),
            types: [
                new Type({
                    discriminator: "!type",
                    name: "Type1",
                    value: new Struct([
                        new Type({
                            name: "field1",
                            value: new Primitive(TypePrimitive.INT),
                        }),
                    ]),
                }),
            ],
            endpoints: [
                new Endpoint({
                    name: "Endpoint1",
                    params: [new Param("param1")],
                    statements: [
                        new Statement({ value: new Action("statement1") }),
                    ],
                }),
            ],
        }),
    ],
});
console.log(model1.toSysl());

// ... or import the entire library using `sysl` and take what you need:
import * as sysl from "@anz-bank/sysl";

const model2 = new sysl.model.Model({
    apps: [
        new sysl.model.Application({ name: new sysl.model.AppName(["App2"]) }),
        // etc.
    ],
});
console.log(model2.toSysl());
