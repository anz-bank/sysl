import * as sysl from "@anz-bank/sysl";
import { Application } from "../../dist/model";

const model = new sysl.model.Model({
    apps: [new Application({ name: "Foo" })],
});
console.log(model.toSysl());
