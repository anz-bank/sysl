import { CollectionDecorator } from "./decorator";
import { ElementRef } from "./elementRef";
import { Primitive } from "./primitive";

export type FieldValue = Primitive | ElementRef | CollectionDecorator;

export namespace FieldValue {
    export function toDto(value: FieldValue) {
        let collectionType: "set" | "sequence" | undefined;
        if (value instanceof CollectionDecorator) {
            collectionType = value.isSet ? "set" : "sequence";
            value = value.innerType;
        }

        return {
            collectionType,
            ref: value instanceof ElementRef ? value.toString() : undefined,
            primitive: value instanceof Primitive ? value.toString() : undefined,
            constraint: value instanceof Primitive ? value.constraint?.toString() : undefined,
        };
    }

    export function fromDto(dto: ReturnType<typeof toDto>): FieldValue {
        let value: FieldValue = dto.ref
            ? ElementRef.parse(dto.ref)
            : Primitive.fromParts(dto.primitive!, dto.constraint);
        if (dto.collectionType) value = new CollectionDecorator(value, dto.collectionType == "set");
        return value;
    }
}
