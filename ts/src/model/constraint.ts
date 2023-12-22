export class Range {
    constructor(public readonly min?: number, public readonly max?: number) {
        if (isNaN(min ?? 0) || isNaN(max ?? 0)) throw new Error("Range min and max must not be NaN.");
        if (min == undefined && max == undefined) throw new Error("Range must have at least one value.");
        if (min != undefined && max != undefined && min > max) throw new Error("Range min cannot be greater than max.");
    }

    toString(): string {
        return `(${this.min ? this.min + ".." : ""}${this.max ?? ""})`;
    }
}

/** Describes constraints on primitive types, such as string length or number of digits. */
export class TypeConstraint {
    constructor(
        /** Optional. The minimum and maximum number of characters in a string or bytes in a blob. */
        public readonly length?: Range,
        /** Optional. The total number of digits in a decimal number, on both sides of the floating point. If specified,
         *  {@link length} is ignored. */
        public readonly precision?: number,
        /** Optional. The number of digits to the right of the floating point, defaults to zero if unspecified. Must be
         * equal to or less than {@link precision}, if specified. */
        public readonly scale?: number,
        /** Optional. The number of bits in a scalar primitive: 32 or 64. Mustn't be specified if {@link precision} is
         *  specified. */
        public readonly bitWidth?: 32 | 64
    ) {
        if (precision && bitWidth) throw new Error("Bit width constraint is not compatible with precision/scale.");
        if (bitWidth != undefined && bitWidth != 32 && bitWidth != 64) throw new Error("Bit width must be 32 or 64.");
        if (scale != undefined && precision == undefined) throw new Error(`Scale (${scale}) requires precision.`);
        if (precision != undefined && scale == undefined) scale = 0;
        if (scale != undefined && precision != undefined) {
            if (precision < 1 || isNaN(precision)) throw new Error(`Precision must be positive: ${precision}`);
            if (scale < 0 || isNaN(scale)) throw new Error(`Scale must be positive or zero: ${scale}`);
            if (scale > precision) throw new Error(`Scale (${scale}) cannot be greater than precision (${precision}).`);
            this.length = undefined;
        }
        this.scale = scale;
    }

    toString(): string {
        if (this.precision) return `(${this.precision}.${this.scale})`;
        return `${this.bitWidth ?? ""}${this.length ?? ""}`;
    }

    static parse(constraintStr: string) {
        const openParenParts = constraintStr.split("(");
        if (openParenParts.length == 1)
            return new TypeConstraint(undefined, undefined, undefined, Number(constraintStr) as 32 | 64);
        if (openParenParts.length != 2 || !openParenParts[1])
            throw new Error(`Invalid constraint string: ${constraintStr}`);
        const closedParenParts = openParenParts[1].split(")");
        if (closedParenParts.length != 2 || !closedParenParts[0] || closedParenParts[1])
            throw new Error(`Invalid constraint string: ${constraintStr}`);

        let length: Range | undefined;
        let precision: number | undefined;
        let scale: number | undefined;
        const bitWidth = openParenParts[0] ? (Number(openParenParts[0]) as 32 | 64) : undefined;

        if (closedParenParts[0]) {
            const constraintParts = closedParenParts[0].split(".").map(s => Number(s));

            switch (constraintParts.length) {
                case 1:
                    length = new Range(undefined, constraintParts[0]);
                    break;
                case 2:
                    precision = constraintParts[0];
                    scale = constraintParts[1];
                    break;
                case 3:
                    if (constraintParts[1] != 0) throw new Error(`Invalid constraint string: ${constraintStr}`);
                    length = new Range(
                        constraintParts[0] ? constraintParts[0] : undefined,
                        constraintParts[2] ? constraintParts[2] : undefined
                    );
                    break;
                default:
                    throw new Error(`Invalid constraint string: ${constraintStr}`);
            }
        }

        return new TypeConstraint(length, precision, scale, bitWidth);
    }
}
