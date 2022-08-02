import "reflect-metadata";
import { jsonArrayMember, jsonObject } from "typedjson";

@jsonObject
export class PbAppName {
    @jsonArrayMember(String) part!: string[];
}
