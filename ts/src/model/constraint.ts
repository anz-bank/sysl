export type TypeConstraintRangeParams = {
    min?: number | undefined;
    max?: number | undefined;
};

export class TypeConstraintRange {
    min: number | undefined;
    max: number | undefined;

    constructor({ min, max }: TypeConstraintRangeParams) {
        this.min = min;
        this.max = max;
    }
}

export type TypeConstraintLengthParams = {
    min?: number | undefined;
    max?: number | undefined;
};

export class TypeConstraintLength {
    min: number | undefined;
    max: number | undefined;

    constructor({ min, max }: TypeConstraintLengthParams) {
        this.min = min;
        this.max = max;
    }
}

export type TypeConstraintResolutionParams = {
    base?: number | undefined;
    index?: number | undefined;
};

/** e.g.: 3 decimal places = {base = 10, index = -3} */
export class TypeConstraintResolution {
    base: number | undefined;
    index: number | undefined;

    constructor({ base, index }: TypeConstraintResolutionParams) {
        this.base = base;
        this.index = index;
    }
}

export type TypeConstraintParams = {
    range?: TypeConstraintRange;
    length?: TypeConstraintLength;
    resolution?: TypeConstraintResolution;
    precision?: number;
    scale?: number;
    bitWidth?: number;
};

export class TypeConstraint {
    range: TypeConstraintRange | undefined;
    length: TypeConstraintLength | undefined;
    resolution: TypeConstraintResolution | undefined;
    precision?: number;
    scale?: number;
    bitWidth?: number;

    constructor({ range, length, resolution, precision, scale, bitWidth }: TypeConstraintParams) {
        this.range = range;
        this.length = length;
        this.resolution = resolution;
        this.precision = precision;
        this.scale = scale;
        this.bitWidth = bitWidth;
    }
}
