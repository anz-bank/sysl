import "reflect-metadata";
import { jsonArrayMember, jsonMapMember, jsonMember, jsonObject } from "typedjson";
import { Location } from "../location";
import { Endpoint, Param, RestParams } from "../model/statement";
import { serializerFor } from "../util";
import { PbAttribute } from "./attribute";
import { PbAppName } from "./appname";
import { PbTypeDef, PbValue } from "./type";

@jsonObject
export class PbParam {
    @jsonMember name!: string;
    @jsonMember(() => PbTypeDef) type!: PbTypeDef;

    toModel(): Param {
        let param = new Param(this.name, this.type.noType ? undefined : this.type.toModel())
        return param;
    }
}

@jsonObject
export class PbAction {
    @jsonMember action!: string;

    // toModel(): Action {
    //     return new Action(this.action)
    // }
}

@jsonObject
export class PbCallArg {
    @jsonMember(() => PbValue) value?: PbValue;
    @jsonMember name!: string;

    // toModel(): CallArg {
    //     return new CallArg(this.value?.toModel(), this.name)
    // }
}

@jsonObject
export class PbCall {
    @jsonMember endpoint!: string;
    @jsonArrayMember(PbCallArg) arg?: PbCallArg[];
    @jsonMember target!: PbAppName;

    // toModel(appName: string[]): Call {
    //     return new Call(this.endpoint,
    //         this.arg?.select(a => a.toModel()).toArray() ?? [],
    //         this.target.part, appName);
    // }
}

@jsonObject
export class PbLoopN {
    @jsonMember count!: number;
    @jsonArrayMember(() => PbStatement) stmt!: PbStatement[];

    // toModel(appName: string[]): LoopN {
    //     return new LoopN(this.count, this.stmt.select(s => s.toModel(appName)).toArray())
    // }
}

@jsonObject
export class PbForeach {
    @jsonMember collection!: string;
    @jsonArrayMember(() => PbStatement) stmt!: PbStatement[];

    // toModel(appName: string[]): Foreach {
    //     return new Foreach(this.collection, this.stmt.select(s => s.toModel(appName)).toArray())
    // }

}

@jsonObject
export class PbAltChoice {
    @jsonMember cond!: string;
    @jsonArrayMember(() => PbStatement) stmt!: PbStatement[];

    // toModel(appName: string[]): AltChoice {
    //     return new AltChoice(this.cond, this.stmt.select(s => s.toModel(appName)).toArray())
    // }
}

@jsonObject
export class PbAlt {
    @jsonArrayMember(PbAltChoice) choice!: PbAltChoice[];

    // toModel(appName: string[]): Alt {
    //     return new Alt(this.choice.select(c => c.toModel(appName)).toArray())
    // }
}

@jsonObject
export class PbGroup {
    @jsonMember title!: string;
    @jsonArrayMember(() => PbStatement) stmt!: PbStatement[];

    // toModel(appName: string[]): Group {
    //     return new Group(this.title, this.stmt.select(s => s.toModel(appName)).toArray())
    // }
}

@jsonObject
export class PbReturn {
    @jsonMember payload!: string;

    // toModel(): Return {
    //     return new Return(this.payload)
    // }

}

@jsonObject
export class PbCond {
    @jsonMember test!: string;
    @jsonArrayMember(() => PbStatement) stmt!: PbStatement[];

    // toModel(appName: string[]): Cond {
    //     return new Cond(this.test, this.stmt.select(s => s.toModel(appName)).toArray())
    // }

}

@jsonObject
export class PbLoop {
    @jsonMember(String) mode!: LoopMode;
    @jsonMember criterion?: string;
    @jsonArrayMember(() => PbStatement) stmt!: PbStatement[];

    // toModel(appName: string[]): Loop {
    //     return new Loop(this.mode, this.criterion, this.stmt.select(s => s.toModel(appName)).toArray());
    // }

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

    // toModel(appName: string[]): Statement {
    //     return new Statement(this.action?.toModel(),
    //         this.call?.toModel(appName),
    //         this.cond?.toModel(appName),
    //         this.loop?.toModel(appName),
    //         this.loopN?.toModel(appName),
    //         this.foreach?.toModel(appName),
    //         this.alt?.toModel(appName),
    //         this.group?.toModel(appName),
    //         this.ret?.toModel(),
    //         this.sourceContexts,
    //         getTags(this.attrs),
    //         getAnnos(this.attrs))
    // }
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
            this.urlParam?.select(p => p.toModel()).toArray() ?? [])
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
            this.param?.select(p => p.toModel()).toArray() ?? [],
            [],
            // this.stmt?.select(s => s.toModel(appName)).toArray(),
            this.restParams?.toModel(),
            this.sourceContexts,
            this.source?.part ?? [])

    }
}
