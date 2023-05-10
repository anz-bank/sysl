import { CloneContext } from "./clone";

export class Range {
    constructor(public min?: number, public max?: number) {}
    public toString = () => `${this.min ?? ""}..${this.max ?? ""}`;
    public clone = () => new Range(this.min, this.max);
}

/** e.g.: 3 decimal places = {base = 10, index = -3} */
export class DecimalResolution {
    constructor(public base?: number, public index?: number) {}
    public toString = () => `${this.base}.${this.index}`;
    public clone = () => new DecimalResolution(this.base, this.index);
}

export class TypeConstraint {
    constructor(
        public range?: Range,
        public length?: Range,
        public resolution?: DecimalResolution,
        public precision?: number,
        public scale?: number,
        public bitWidth?: number
    ) {}

    toString(): string {
        return [
            this.range ? `range: ${this.range}` : undefined,
            this.length ? `length: ${this.length}` : undefined,
            this.resolution ? `resolution: ${this.resolution}` : undefined,
            this.precision ? `precision: ${this.precision}` : undefined,
            this.scale ? `scale: ${this.scale}` : undefined,
            this.bitWidth ? `bitWidth: ${this.bitWidth}` : undefined,
        ].filter(x => x).join(", ");
    }

    clone(_context?: CloneContext): TypeConstraint {
        return new TypeConstraint(
            this.range?.clone(),
            this.length?.clone(),
            this.resolution?.clone(),
            this.precision,
            this.scale,
            this.bitWidth
        );
    }
}
