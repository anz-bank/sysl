import "reflect-metadata";
import { jsonMember, jsonObject } from "typedjson";
import path from "path";

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
    /** The file name where the fragment is located. */
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
     * If the file is
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
}
