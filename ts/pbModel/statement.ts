import "reflect-metadata";
import {
    jsonArrayMember,
    jsonMapMember,
    jsonMember,
    jsonObject,
} from "typedjson";
import { Location } from "../location";
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
import { getAnnos, getTags, sortLocationalArray } from "../sort";
import { PbAppName } from "./appname";
import { PbAttribute } from "./attribute";
import { serializerFor } from "./serialize";
import { PbTypeDef, PbValue } from "./type";

@jsonObject
export class PbParam {
    @jsonMember name!: string;
    @jsonMember(() => PbTypeDef) type!: PbTypeDef;

    toModel(): Param {
        let param = new Param(
            this.name,
            this.type.hasValue() ? this.type.toModel() : undefined
        );
        return param;
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
        return new CallArg(this.value?.toModel(), this.name);
    }
}

@jsonObject
export class PbCall {
    @jsonMember endpoint!: string;
    @jsonArrayMember(PbCallArg) arg?: PbCallArg[];
    @jsonMember target!: PbAppName;

    toModel(appName: string[]): Call {
        return new Call(
            this.endpoint,
            this.arg?.select(a => a.toModel()).toArray() ?? [],
            this.target.part,
            appName
        );
    }
}

@jsonObject
export class PbLoopN {
    @jsonMember count!: number;
    @jsonArrayMember(() => PbStatement) stmt!: PbStatement[];

    toModel(appName: string[]): LoopN {
        return new LoopN(
            this.count,
            this.stmt.select(s => s.toModel(appName)).toArray()
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
            this.stmt.select(s => s.toModel(appName)).toArray()
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
            this.stmt.select(s => s.toModel(appName)).toArray()
        );
    }
}

@jsonObject
export class PbAlt {
    @jsonArrayMember(PbAltChoice) choice!: PbAltChoice[];

    toModel(appName: string[]): Alt {
        return new Alt(this.choice.select(c => c.toModel(appName)).toArray());
    }
}

@jsonObject
export class PbGroup {
    @jsonMember title!: string;
    @jsonArrayMember(() => PbStatement) stmt!: PbStatement[];

    toModel(appName: string[]): Group {
        return new Group(
            this.title,
            this.stmt.select(s => s.toModel(appName)).toArray()
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
            this.stmt.select(s => s.toModel(appName)).toArray()
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
            this.stmt.select(s => s.toModel(appName)).toArray()
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

    getValue():
        | PbAction
        | PbCall
        | PbCond
        | PbLoop
        | PbLoopN
        | PbForeach
        | PbAlt
        | PbGroup
        | PbReturn
        | undefined {
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
        return new Statement(
            this.getValue()?.toModel(appName),
            this.sourceContexts,
            getTags(this.attrs),
            getAnnos(this.attrs)
        );
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

    toModel(): RestParams {
        return new RestParams(
            this.method,
            this.path,
            this.queryParam?.select(p => p.toModel()).toArray() ?? [],
            this.urlParam?.select(p => p.toModel()).toArray() ?? []
        );
    }
}

@jsonObject
export class PbEndpoint {
    @jsonMember name!: string;
    @jsonMember longName?: string;
    @jsonMember docstring?: string;
    @jsonArrayMember(String) flag?: string[];
    @jsonMember isPubsub?: boolean;
    @jsonArrayMember(PbParam) param?: PbParam[];
    @jsonArrayMember(PbStatement) stmt!: PbStatement[];
    @jsonMember restParams?: PbRestParams;
    @jsonArrayMember(Location) sourceContexts!: Location[];
    @jsonMapMember(String, () => PbAttribute, serializerFor(PbAttribute))
    attrs!: Map<string, PbAttribute>;
    @jsonMember source?: PbAppName;

    toModel(appName: string[]): Endpoint {
        return new Endpoint(
            this.name,
            this.longName,
            this.docstring,
            this.flag,
            this.isPubsub ?? false,
            this.param
                ?.where(p => p.name != undefined || p.type != undefined)
                .select(p => p.toModel())
                .toArray() ?? [],
            sortLocationalArray(
                this.stmt?.select(s => s.toModel(appName)).toArray() ?? []
            ),
            this.restParams?.toModel(),
            this.source?.part ?? [],
            this.sourceContexts,
            getTags(this.attrs),
            getAnnos(this.attrs)
        );
    }
}
