import "jest-extended";
import { Union } from "./union";

test.concurrent("UnionWithUndefinedChildren", async () => {
    // @ts-ignore
    const u = new Union("x", undefined, {});
    const u2 = new Union("x", [], {});
    const uDto = u.toDto();
    expect(uDto).toEqual(u2.toDto());
});
