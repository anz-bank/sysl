import "jest-extended";
import { realign } from "../common/format";
import { Model } from "./model";
import { Element } from "./element";

test.concurrent("anno operations", async () => {
    const model = await Model.fromText(realign(`
        App:
            @anno1 = "A"
            @anno1 = "B"
    `));
    const app = model.getApp("App");
    const anno = app.annos[0];
    expect(() => app.getAnno("nonexistent")).toThrow();
    expect(app.findAnno("nonexistent")).toBeUndefined();
    expect(app.getAnno("anno1").value).toBe("A");
    expect(app.findAnno("anno1")?.value).toBe("A");

    expect(app.setAnno("anno1", "C")).toBe(anno);
    expect(anno.value).toBe("C");
    expect(app.getAnno("anno1").value).toBe("C");

    expect(app.setAnno("anno1", undefined)).toBe(anno);
    expect(app.findAnno("anno1")).toBeUndefined();

    expect(app.setAnno("anno2", "B")?.parent).toBe(app);
    expect(app.getAnno("anno2").value).toBe("B");

    expect(model.toSysl()).toEqual(realign(`
        App:
            @anno2 = "B"
    `));
});

test.concurrent("tag operations", async () => {
    const model = await Model.fromText(realign(`
        App [~tag1]:
            ...
    `));
    const app = model.getApp("App");
    const tag = app.tags[0];
    expect(() => app.getTag("nonexistent")).toThrow();
    expect(app.findTag("nonexistent")).toBeUndefined();
    expect(app.hasTag("nonexistent")).toBeFalse();

    expect(app.getTag("tag1")).toBe(tag);
    expect(app.findTag("tag1")).toBe(tag);
    expect(app.hasTag("tag1", "nonexistent")).toBeTrue();

    expect(app.setTag("tag1")).toBe(tag);
    expect(app.tags.length).toBe(1);

    expect(app.removeTag("tag1")).toBe(tag);
    expect(app.hasTag("tag1")).toBeFalse();

    expect(app.setTag("tag2")?.parent).toBe(app);
    expect(app.hasTag("tag2")).toBeTrue();

    expect(model.toSysl()).toEqual(realign(`
        App [~tag2]:
            ...
    `));
});

test.concurrent("fromDto", async () => {
    const model = await Model.fromText(realign(`
        App [~tag1, ~tag2]:
            @anno1 = "A"
            @anno2 = "B"
    `));

    const dto = model.apps[0].toDto();
    expect(Element.paramsFromDto(dto)).toMatchObject({
        locations: [{}],
        annos: [ { name: "anno1", locations: [{}] }, { name: "anno2", locations: [{}] }],
        tags: [ { name: "tag1", locations: [{}] }, { name: "tag2", locations: [{}] }]
    });
});
