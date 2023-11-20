import "reflect-metadata";
import { jsonArrayMember, jsonMapMember, jsonMember, jsonObject } from "typedjson";
import { Location } from "../common/location";
import {
    ActionStatement,
    CallArg,
    CallStatement,
    CondStatement,
    ForEachStatement,
    GroupStatement,
    LoopStatement,
    OneOfStatement,
    RestParams,
    ReturnStatement,
    Statement,
    ellipsis,
} from "../model/statement";
import { Endpoint } from "../model/endpoint";
import { PbAppName } from "./appname";
import { getAnnos, getTags, PbAttribute } from "./attribute";
import { serializerFor } from "./serialize";
import { PbTypeDef, PbValue } from "./type";
import { ElementRef, Field, Primitive, TypePrimitive } from "../model";
import * as util from "util";

@jsonObject
export class PbParam {
    @jsonMember name!: string;
    @jsonMember(() => PbTypeDef) type!: PbTypeDef;

    toModel(parentRef: ElementRef): Field {
        if (this.type.hasValue()) return this.type.toField(parentRef, this.name);
        return new Field(this.name, new Primitive(TypePrimitive.ANY), false, { locations: this.type.sourceContexts });
    }
}

@jsonObject
export class PbCallArg {
    @jsonMember(() => PbValue) value?: PbValue;
    @jsonMember name!: string;

    toModel(): CallArg {
        return new CallArg(this.name, this.value?.toModel());
    }
}

@jsonObject
export class PbLoopN {
    @jsonMember count!: number;
    @jsonArrayMember(() => PbStatement) stmt!: PbStatement[];
}

@jsonObject
export class PbForeach {
    @jsonMember collection!: string;
    @jsonArrayMember(() => PbStatement) stmt!: PbStatement[];
}

@jsonObject
export class PbAltChoice {
    @jsonMember cond?: string; // TODO: Remove optionality when Sysl parser rejects empty names for choices.
    @jsonArrayMember(() => PbStatement) stmt!: PbStatement[];
}

@jsonObject
export class PbGroup {
    @jsonMember title!: string;
    @jsonArrayMember(() => PbStatement) stmt!: PbStatement[];
}


@jsonObject
export class PbCond {
    @jsonMember test!: string;
    @jsonArrayMember(() => PbStatement) stmt!: PbStatement[];
}

@jsonObject
export class PbLoop {
    @jsonMember mode!: string;
    @jsonMember criterion?: string;
    @jsonArrayMember(() => PbStatement) stmt!: PbStatement[];
}

@jsonObject
export class PbCall {
    @jsonMember endpoint!: string;
    @jsonArrayMember(PbCallArg) arg?: PbCallArg[];
    @jsonMember target!: PbAppName;
}

@jsonObject export class PbAction { @jsonMember action!: string; }
@jsonObject export class PbReturn { @jsonMember payload!: string; }
@jsonObject export class PbAlt    { @jsonArrayMember(PbAltChoice) choice!: PbAltChoice[]; }

@jsonObject
export class PbStatement {
    @jsonMember action?: PbAction;
    @jsonMember call?: PbCall;
    @jsonMember cond?: PbCond;
    @jsonMember loop?: PbLoop;
    @jsonMember loopN?: PbLoopN;
    @jsonMember foreach?: PbForeach;
    @jsonMember alt?: PbAlt;
    @jsonMember group?: PbGroup;
    @jsonMember ret?: PbReturn;
    @jsonArrayMember(Location) sourceContexts!: Location[];
    @jsonMapMember(String, () => PbAttribute, serializerFor(PbAttribute))
    attrs!: Map<string, PbAttribute>;

    toModel(appRef: ElementRef): Statement | undefined {
        const params = {
            locations: this.sourceContexts,
            tags: getTags(this.attrs),
            annos: getAnnos(this.attrs),
        };
        if (this.action) {
            if (this.action.action == ellipsis) return undefined
            return new ActionStatement(this.action.action, params);
        }
        if (this.ret) return new ReturnStatement(this.ret.payload, params);
        if (this.call)
            return new CallStatement(
                ElementRef.fromAppParts(this.call.target.part).with({ endpointName: this.call.endpoint }),
                this.call.arg?.map(c => c.toModel()) ?? [],
                appRef,
                params
            );

        const children = PbStatement.getBlock(this.group ?? this.cond ?? this.loop ?? this.foreach, appRef);
        const blockParams = { ...params, children };
        if (this.group) return new GroupStatement(this.group.title, blockParams);
        if (this.cond) return new CondStatement(this.cond.test, blockParams);
        if (this.foreach) return new ForEachStatement(this.foreach.collection, blockParams);
        if (this.alt) {
            const toGroup = (choice: PbAltChoice) => 
                // TODO: Remove fallback to empty string when Sysl parser rejects empty names for choices.
                new GroupStatement(choice.cond ?? "", { ...blockParams, children: PbStatement.getBlock(choice, appRef) });
            
            return new OneOfStatement(this.alt.choice.map(c => toGroup(c)), blockParams);
        }
        if (this.loop) {
            if (this.loop.mode == "WHILE" || this.loop.mode == "UNTIL")
                return new LoopStatement(this.loop.criterion ?? "", this.loop.mode.toLowerCase(), blockParams);
            throw new Error(`Unrecognized loop mode: ${this.loop.mode.toString()}`);
        }

        if (this.loopN) throw new Error("LoopN statement is not supported.");
        throw new Error("Encountered unrecognized statement type in PB:\n" + util.inspect(this));
    }

    static getBlock(container: { stmt?: PbStatement[] } | undefined, appRef: ElementRef): Statement[] {
        if (!container?.stmt?.length) return [];
        const block = container.stmt.map(s => s.toModel(appRef)).filter(s => s) as Statement[];
        return block.sort(Location.compareFirst);
    }
}

export enum RestMethod {
    NOMethod = "NOMethod",
    GET = "GET",
    PUT = "PUT",
    POST = "POST",
    DELETE = "DELETE",
    PATCH = "PATCH",
}

@jsonObject
export class PbRestParams {
    @jsonMember method!: RestMethod;
    @jsonMember path!: string;
    @jsonArrayMember(PbParam) queryParam?: PbParam[];
    @jsonArrayMember(PbParam) urlParam?: PbParam[];

    toModel(parentRef: ElementRef): RestParams {
        return new RestParams({
            method: this.method,
            path: this.path,
            queryParams: this.queryParam?.map(p => p.toModel(parentRef)) ?? [],
            urlParams: this.urlParam?.map(p => p.toModel(parentRef)) ?? [],
        });
    }
}

@jsonObject
export class PbEndpoint {
    @jsonMember name!: string;
    @jsonMember longName?: string;
    @jsonMember docstring?: string;
    @jsonMember isPubsub?: boolean;
    @jsonArrayMember(PbParam) param?: PbParam[];
    @jsonArrayMember(PbStatement) stmt!: PbStatement[];
    @jsonMember restParams?: PbRestParams;
    @jsonArrayMember(Location) sourceContexts!: Location[];
    @jsonMapMember(String, () => PbAttribute, serializerFor(PbAttribute))
    attrs!: Map<string, PbAttribute>;
    @jsonMember source?: PbAppName;

    toModel(appRef: ElementRef): Endpoint {
        return new Endpoint(this.name, {
            longName: this.longName,
            isPubsub: this.isPubsub ?? false,
            params: (this.param ?? []).filter(p => p.name || p.type).map(p => p.toModel(appRef)),
            children: PbStatement.getBlock(this, appRef),
            restParams: this.restParams?.toModel(appRef),
            locations: this.sourceContexts,
            tags: getTags(this.attrs),
            annos: getAnnos(this.attrs),
        });
    }
}
