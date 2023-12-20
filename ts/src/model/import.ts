import { Location } from "../common/location";
import { CloneContext, ICloneable } from "./clone";

export type ImportParams = {
    filePath: string;
    locations?: Location[];
    appAlias?: string;
};

export class Import implements ICloneable {
    filePath: string;
    locations: Location[];
    appAlias?: string;

    constructor({ filePath, locations, appAlias }: ImportParams) {
        this.filePath = filePath;
        this.locations = locations ?? [];
        this.appAlias = appAlias ? appAlias : undefined;
    }

    toSysl(): string {
        return `import ${this.filePath}${this.appAlias ? ` as ${this.appAlias}` : ""}`;
    }

    toDto() {
        return {
            filePath: this.filePath,
            locations: this.locations.map(l => l.toString()),
            appAlias: this.appAlias,
        };
    }

    static fromDto(dto: ReturnType<Import["toDto"]>): Import {
        return new Import({
            filePath: dto.filePath,
            locations: dto.locations.map(Location.parse),
            appAlias: dto.appAlias,
        });
    }

    toString(): string {
        return this.filePath;
    }

    clone(_context?: CloneContext): Import {
        return new Import({ filePath: this.filePath, appAlias: this.appAlias });
    }
}
