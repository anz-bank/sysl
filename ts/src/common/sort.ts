import { ILocational } from "../model/common";

export function sortLocationalArray<T extends ILocational>(array: T[]): T[] {
    return array.sort((i1, i2) => {
        return i1.locations[0]?.file.localeCompare(i2.locations[0]?.file)
        ||     i1.locations[0]?.start.line - i2.locations[0]?.start.line
        ||     i1.locations[0]?.start.col  - i2.locations[0]?.start.col
        ||     i1.locations[0]?.end.line   - i2.locations[0]?.end.line
        ||     i1.locations[0]?.end.col    - i2.locations[0]?.end.col
    });
}
