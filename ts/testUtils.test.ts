import { initializeLinq } from "linq-to-typescript";
import { realign } from "./testUtils";

initializeLinq();

test("realign", () => {
    expect(realign(`
    App "display name" [~abstract]:
        ...
    `)).toBe(
`App "display name" [~abstract]:
    ...
`);
});
