import "jest-extended";
import * as index from "./index";

test("exports each directory", () => {
    const dirs = ["common", "model", "pbModel"];
    dirs.forEach(d => expect(index).toHaveProperty(d));

    expect(new index.model.Model()).not.toBeNull();
    expect(new index.pbModel.PbDocumentModel()).not.toBeNull();
});
