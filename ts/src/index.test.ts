import * as index from "./index";

it("exports each directory", () => {
    const dirs = ["model", "pbModel"];
    dirs.forEach(d => expect(index).toHaveProperty(d));

    expect(new index.model.Model()).not.toBeNull();
    expect(new index.pbModel.PbDocumentModel()).not.toBeNull();
});
