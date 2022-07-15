import { jsonMember, jsonObject } from "typedjson";

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
}
