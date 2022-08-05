import { ILocational } from "../model/common";

function getStart<T extends ILocational>(item: T): number {
    if (!item.locations?.length) {
        return 0;
    }
    const firstLoc = item.locations[0];
    return firstLoc.start ? firstLoc.start.line : firstLoc.end.line;
}

export function sortLocationalArray<T extends ILocational>(array: T[]): T[] {
    return array.sort((i1, i2) => {
        return getStart(i1) - getStart(i2);
    });
}
