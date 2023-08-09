import "reflect-metadata";
import { jsonArrayMember, jsonMapMember, jsonMember, jsonObject } from "typedjson";
import { Location } from "../common/location";
import { sortByLocation } from "../common/sort";
import {
    Action,
    Alt,
    AltChoice,
    Call,
    CallArg,
    Cond,
    Endpoint,
    Foreach,
    Group,
    Loop,
    LoopN,
    Param,
    RestParams,
    Return,
    Statement,
} from "../model/statement";
import { PbAppName } from "./appname";
import { getAnnos, getTags, PbAttribute } from "./attribute";
import { serializerFor } from "./serialize";
import { PbTypeDef, PbValue } from "./type";
import { ElementRef } from "../model";

@jsonObject
export class PbParam {
    @jsonMember name!: string;
    @jsonMember(() => PbTypeDef) type!: PbTypeDef;

    toModel(parentRef: ElementRef): Param {
        return new Param(
            this.name,
            this.type.sourceContexts ?? [],
            this.type.hasValue() ? this.type.toModel(undefined, undefined, parentRef) : undefined
        );
    }
}

@jsonObject
export class PbAction {
    @jsonMember action!: string;

    toModel(): Action {
        return new Action(this.action);
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
export class PbCall {
    @jsonMember endpoint!: string;
    @jsonArrayMember(PbCallArg) arg?: PbCallArg[];
    @jsonMember target!: PbAppName;

    toModel(appName: string[]): Call {
        return new Call({
            endpoint: this.endpoint,
            args: this.arg?.map(a => a.toModel()) ?? [],
            targetApp: this.target.part,
            originApp: appName,
        });
    }
}

@jsonObject
export class PbLoopN {
    @jsonMember count!: number;
    @jsonArrayMember(() => PbStatement) stmt!: PbStatement[];

    toModel(appName: string[]): LoopN {
        return new LoopN(
            this.count,
            this.stmt.map(s => s.toModel(appName))
        );
    }
}

@jsonObject
export class PbForeach {
    @jsonMember collection!: string;
    @jsonArrayMember(() => PbStatement) stmt!: PbStatement[];

    toModel(appName: string[]): Foreach {
        return new Foreach(
            this.collection,
            this.stmt.map(s => s.toModel(appName))
        );
    }
}

@jsonObject
export class PbAltChoice {
    @jsonMember cond!: string;
    @jsonArrayMember(() => PbStatement) stmt!: PbStatement[];

    toModel(appName: string[]): AltChoice {
        return new AltChoice(
            this.cond,
            this.stmt.map(s => s.toModel(appName))
        );
    }
}

@jsonObject
export class PbAlt {
    @jsonArrayMember(PbAltChoice) choice!: PbAltChoice[];

    toModel(appName: string[]): Alt {
        return new Alt(this.choice.map(c => c.toModel(appName)));
    }
}

@jsonObject
export class PbGroup {
    @jsonMember title!: string;
    @jsonArrayMember(() => PbStatement) stmt!: PbStatement[];

    toModel(appName: string[]): Group {
        return new Group(
            this.title,
            this.stmt.map(s => s.toModel(appName))
        );
    }
}

@jsonObject
export class PbReturn {
    @jsonMember payload!: string;

    toModel(): Return {
        return new Return(this.payload);
    }
}

@jsonObject
export class PbCond {
    @jsonMember test!: string;
    @jsonArrayMember(() => PbStatement) stmt!: PbStatement[];

    toModel(appName: string[]): Cond {
        return new Cond(
            this.test,
            this.stmt.map(s => s.toModel(appName))
        );
    }
}

@jsonObject
export class PbLoop {
    @jsonMember(String) mode!: LoopMode;
    @jsonMember criterion?: string;
    @jsonArrayMember(() => PbStatement) stmt!: PbStatement[];

    toModel(appName: string[]): Loop {
        return new Loop(
            this.mode,
            this.criterion,
            this.stmt.map(s => s.toModel(appName))
        );
    }
}

export enum LoopMode {
    NOMode = 0,
    WHILE = 1,
    UNTIL = 2,
    UNRECOGNIZED = -1,
}

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

    getValue(): PbAction | PbCall | PbCond | PbLoop | PbLoopN | PbForeach | PbAlt | PbGroup | PbReturn | undefined {
        return (
            this.action ||
            this.call ||
            this.cond ||
            this.loop ||
            this.loopN ||
            this.foreach ||
            this.alt ||
            this.group ||
            this.ret
        );
    }

    toModel(appName: string[]): Statement {
        return new Statement(this.getValue()?.toModel(appName), {
            locations: this.sourceContexts,
            tags: getTags(this.attrs),
            annos: getAnnos(this.attrs),
        });
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

    toModel(appName: string[]): Endpoint {
        const appRef = ElementRef.fromAppParts(appName);
        return new Endpoint(this.name, {
            longName: this.longName,
            docstring: this.docstring,
            isPubsub: this.isPubsub ?? false,
            params: (this.param ?? []).filter(p => p.name || p.type).map(p => p.toModel(appRef)),
            statements: sortByLocation(this.stmt?.map(s => s.toModel(appName)) ?? []),
            restParams: this.restParams?.toModel(appRef),
            pubsubSource: this.source?.part ?? [],
            locations: this.sourceContexts,
            tags: getTags(this.attrs),
            annos: getAnnos(this.attrs),
        });
    }
}
