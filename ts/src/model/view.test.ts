import { realign } from "../common/format";
import { Element } from "./element";
import { Model } from "./model";

test.concurrent("FlatView", async () => {
    const view = (await Model.fromText(realign(`
        Company :: Backend:
            !table Customer:
                HomeAddress <: Company::App.Address
            !type Address:
                City <: string

        Company :: Frontend:
            /:
                GET:
                    Statement1:
                        Statement1_1:
                            Statement1_1_1
                            Statement1_1_2
                        Statement1_2
                    Statement2
                    Statement3:
                        if cond:
                            Statement3_if1_1
        `))).flat();

    const expectNames = (elements: readonly Element[]) => expect(elements.map(e => e.toString()));
    expectNames(view.apps).toEqual(["Company::Backend", "Company::Frontend"]);
    expectNames(view.types).toEqual(["!table Customer", "!type Address"]);
    expectNames(view.fields).toEqual(["HomeAddress <: Company::App.Address", "City <: string"]);
    expectNames(view.dataElements).toEqual([
        "Company::Backend",
        "Company::Frontend",
        "!table Customer",
        "!type Address",
        "HomeAddress <: Company::App.Address",
        "City <: string",
    ]);
    expectNames(view.endpoints).toEqual(["[REST] /"]);
    expectNames(view.statements).toEqual([
        "Statement1",
        "Statement2",
        "Statement3",
        "Statement1_1",
        "Statement1_2",
        "if cond",
        "Statement1_1_1",
        "Statement1_1_2",
        "Statement3_if1_1",
    ]);
    expectNames(view.behaviourElements).toEqual([
        "Company::Backend",
        "Company::Frontend",
        "[REST] /",
        "Statement1",
        "Statement2",
        "Statement3",
        "Statement1_1",
        "Statement1_2",
        "if cond",
        "Statement1_1_1",
        "Statement1_1_2",
        "Statement3_if1_1",
    ]);
    expectNames(view.allElements).toEqual([
        "Company::Backend",
        "Company::Frontend",
        "!table Customer",
        "!type Address",
        "HomeAddress <: Company::App.Address",
        "City <: string",
        "[REST] /",
        "Statement1",
        "Statement2",
        "Statement3",
        "Statement1_1",
        "Statement1_2",
        "if cond",
        "Statement1_1_1",
        "Statement1_1_2",
        "Statement3_if1_1",
    ]);
});
