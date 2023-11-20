import "reflect-metadata";
import { jsonMember, jsonObject } from "typedjson";
import path from "path";
import { ILocational, Model } from "../model";

/** Describes the offset location within a file by zero-based line and column. */
@jsonObject
export class Offset {
    /** The zero-based line number of the {@link Offset}. */
    @jsonMember line: number;
    /** The zero-based column number of the {@link Offset}. */
    @jsonMember col: number;

    /**
     * Create a new {@link Offset} object with the specified zero-based line and column.
     * @param line Optional. The zero-based line number of the {@link Offset}. If not specified, defaults to `0`.
     * @param col Optional. The zero-based column number of the {@link Offset}. If not specified, defaults to `0`.
     */
    constructor(line: number = 0, col: number = 0) {
        this.line = line;
        this.col = col;
    }
}

/** Describes the source location of a Sysl document fragment. */
@jsonObject
export class Location {
    /** The Sysl path of the file where the fragment is located, relative to {@link Model.syslRoot}. */
    @jsonMember file: string;
    /** The {@link Offset} inside the {@link file} where the fragment begins. */
    @jsonMember start: Offset;
    /** The {@link Offset} inside the {@link file} where the fragment ends. */
    @jsonMember end: Offset;

    /**
     * Creates a new {@link Location} object with the specified file and offset range.
     * @param file The file name where the fragment is located.
     * @param start The {@link Offset} inside the {@link file} where the fragment begins.
     * @param end The {@link Offset} inside the {@link file} where the fragment ends.
     */
    constructor(file: string, start: Offset, end: Offset) {
        this.file = file;
        this.start = start;
        this.end = end;
    }

    /**
     * Gets the file extension part of {@link file} path. This is the part of the file name beginning at the first
     * period.
     * @example
     * If the {@link file} is
     * ```//github.com/someone/myproject/program.edg.sysl@master```
     * the returned value will be `.edg.sysl`
     */
    public get extension(): string {
        return "." + path.basename(this.file).split("@")[0].split(".").slice(1).join(".");
    }

    /**
     * Gets the git reference (branch, tag of hash) of the {@link file} path. This is the part after the '@' symbol,
     * or an empty string if none is specified.
     */
    public get gitRef(): string {
        return path.basename(this.file).split("@")[1] ?? "";
    }

    /**
     * True if the {@link file} path point to a remote file (starts with `//`), otherwise false.
     */
    public get isRemote(): boolean {
        return this.file.startsWith("//");
    }

    /**
     * Formats the {@link Location} object as a string with five parts joined by the colon character:
     * - File name
     * - 1-based start line
     * - 1-based start column (omitted if 0 and end line/column are omitted)
     * - 1-based end line (omitted if identical to start line)
     * - 1-based end column (omitted if not greater than start column and end line is omitted)
     * This format is backwards-compatible with the file linking in VS Code (for file, and start line/col).
     * @returns A string representation of the location.
     */
    public toString(startLineOnly: boolean = false): string {
        const parts = [this.file, this.start.line + 1];

        if (!startLineOnly) {
            const hasEndLine = this.end.line > this.start.line;
            const hasEnd = hasEndLine || this.end.col > this.start.col;
            if (this.start.col > 0 || hasEnd) {
                parts.push(this.start.col + 1);
                if (hasEnd) {
                    parts.push(hasEndLine ? this.end.line + 1 : "");
                    parts.push(this.end.col + 1);
                }
            }
        }

        return parts.join(":");
    }

    public clone(): Location {
        return new Location(
            this.file,
            new Offset(this.start.line, this.start.col),
            new Offset(this.end.line, this.end.col)
        );
    }

    public static compare(a: Location, b: Location): number {
        return (
            a?.file.localeCompare(b?.file) ||
            a?.start.line - b?.start.line ||
            a?.start.col - b?.start.col ||
            a?.end.line - b?.end.line ||
            a?.end.col - b?.end.col
        );
    }

    public static compareFirst(a: ILocational, b: ILocational): number {
        return Location.compare(a.locations[0], b.locations[0]);
    }
}
