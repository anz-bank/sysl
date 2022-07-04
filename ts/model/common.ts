import { Location } from "../location";

export abstract class BaseType {
    discriminator: string;

    constructor(discriminator: string) {
        this.discriminator = discriminator;
    }

    toSysl(): string { return ""; }
}

export abstract class ComplexType extends BaseType {
    locations: Location[];
    name: string;

    constructor(discriminator: string, locations: Location[], name: string) {
        super(discriminator);
        this.locations = locations;
        this.name = name ?? undefined;
    }
}

export abstract class SimpleType extends BaseType {
    constructor() {
        super("");
    }
}
