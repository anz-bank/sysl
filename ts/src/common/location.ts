import "reflect-metadata";
import { jsonMember, jsonObject } from "typedjson";
import path from "path";

@jsonObject
export class Offset {
    @jsonMember col!: number;
    @jsonMember line!: number;
}

@jsonObject
export class Location {
    @jsonMember file!: string;
    @jsonMember start!: Offset;
    @jsonMember end!: Offset;

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
