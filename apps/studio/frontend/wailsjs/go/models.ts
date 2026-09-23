export namespace agentadapter {
	
	export class Event {
	    type: string;
	    data: string;
	
	    static createFrom(source: any = {}) {
	        return new Event(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.data = source["data"];
	    }
	}
	export class Manifest {
	    id: string;
	    version: string;
	    transport: string;
	    capabilities: string[];
	
	    static createFrom(source: any = {}) {
	        return new Manifest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.version = source["version"];
	        this.transport = source["transport"];
	        this.capabilities = source["capabilities"];
	    }
	}

}

export namespace execution {
	
	export class RunResult {
	    cardId: string;
	    worktree: string;
	    state: string;
	    events: agentadapter.Event[];
	    changedFiles: string[];
	
	    static createFrom(source: any = {}) {
	        return new RunResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cardId = source["cardId"];
	        this.worktree = source["worktree"];
	        this.state = source["state"];
	        this.events = this.convertValues(source["events"], agentadapter.Event);
	        this.changedFiles = source["changedFiles"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace productgraph {
	
	export class Provenance {
	    derivedFrom?: string[];
	    generatedBy?: string[];
	    attributedTo?: string[];
	    confidence: number;
	    // Go type: time
	    generatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Provenance(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.derivedFrom = source["derivedFrom"];
	        this.generatedBy = source["generatedBy"];
	        this.attributedTo = source["attributedTo"];
	        this.confidence = source["confidence"];
	        this.generatedAt = this.convertValues(source["generatedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SourceRef {
	    kind: string;
	    ref: string;
	
	    static createFrom(source: any = {}) {
	        return new SourceRef(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.kind = source["kind"];
	        this.ref = source["ref"];
	    }
	}
	export class NodeMetadata {
	    source?: SourceRef[];
	    status: string;
	    impact?: string;
	    provenance?: Provenance;
	    revision: number;
	
	    static createFrom(source: any = {}) {
	        return new NodeMetadata(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.source = this.convertValues(source["source"], SourceRef);
	        this.status = source["status"];
	        this.impact = source["impact"];
	        this.provenance = this.convertValues(source["provenance"], Provenance);
	        this.revision = source["revision"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class AcceptanceCriterion {
	    id: string;
	    requirementId: string;
	    statement: string;
	    mandatory: boolean;
	    verificationType: string;
	    metadata: NodeMetadata;
	
	    static createFrom(source: any = {}) {
	        return new AcceptanceCriterion(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.requirementId = source["requirementId"];
	        this.statement = source["statement"];
	        this.mandatory = source["mandatory"];
	        this.verificationType = source["verificationType"];
	        this.metadata = this.convertValues(source["metadata"], NodeMetadata);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class AdapterManifest {
	    id: string;
	    version: string;
	    transport: string;
	    capabilities: string[];
	
	    static createFrom(source: any = {}) {
	        return new AdapterManifest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.version = source["version"];
	        this.transport = source["transport"];
	        this.capabilities = source["capabilities"];
	    }
	}
	export class CardScheduleDecision {
	    cardId: string;
	    state: string;
	    reasons?: string[];
	
	    static createFrom(source: any = {}) {
	        return new CardScheduleDecision(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cardId = source["cardId"];
	        this.state = source["state"];
	        this.reasons = source["reasons"];
	    }
	}
	export class CardSpec {
	    id: string;
	    goal: string;
	    scopePaths: string[];
	    forbiddenPaths: string[];
	    dependencies?: string[];
	    semanticLocks?: string[];
	    sharedResourceLocks?: string[];
	    acceptanceCriteria: string[];
	    verificationCommands: string[];
	    requiredAdapterCapabilities: string[];
	    state: string;
	
	    static createFrom(source: any = {}) {
	        return new CardSpec(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.goal = source["goal"];
	        this.scopePaths = source["scopePaths"];
	        this.forbiddenPaths = source["forbiddenPaths"];
	        this.dependencies = source["dependencies"];
	        this.semanticLocks = source["semanticLocks"];
	        this.sharedResourceLocks = source["sharedResourceLocks"];
	        this.acceptanceCriteria = source["acceptanceCriteria"];
	        this.verificationCommands = source["verificationCommands"];
	        this.requiredAdapterCapabilities = source["requiredAdapterCapabilities"];
	        this.state = source["state"];
	    }
	}
	export class Gap {
	    id: string;
	    type: string;
	    severity: string;
	    impact: string;
	    status: string;
	    metadata: NodeMetadata;
	
	    static createFrom(source: any = {}) {
	        return new Gap(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.type = source["type"];
	        this.severity = source["severity"];
	        this.impact = source["impact"];
	        this.status = source["status"];
	        this.metadata = this.convertValues(source["metadata"], NodeMetadata);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Requirement {
	    id: string;
	    statement: string;
	    metadata: NodeMetadata;
	
	    static createFrom(source: any = {}) {
	        return new Requirement(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.statement = source["statement"];
	        this.metadata = this.convertValues(source["metadata"], NodeMetadata);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class CompiledPlan {
	    revisionId: string;
	    requirements: Requirement[];
	    acceptanceCriteria: AcceptanceCriterion[];
	    cards: CardSpec[];
	    blockingGaps?: Gap[];
	
	    static createFrom(source: any = {}) {
	        return new CompiledPlan(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.revisionId = source["revisionId"];
	        this.requirements = this.convertValues(source["requirements"], Requirement);
	        this.acceptanceCriteria = this.convertValues(source["acceptanceCriteria"], AcceptanceCriterion);
	        this.cards = this.convertValues(source["cards"], CardSpec);
	        this.blockingGaps = this.convertValues(source["blockingGaps"], Gap);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	export class Project {
	    id: string;
	    name: string;
	    brief?: string;
	    revision: number;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Project(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.brief = source["brief"];
	        this.revision = source["revision"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	export class Revision {
	    id: string;
	    projectId: string;
	    number: number;
	    status: string;
	    snapshot: string;
	    // Go type: time
	    createdAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Revision(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.projectId = source["projectId"];
	        this.number = source["number"];
	        this.status = source["status"];
	        this.snapshot = source["snapshot"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ScheduleResult {
	    adapter: AdapterManifest;
	    decisions: CardScheduleDecision[];
	
	    static createFrom(source: any = {}) {
	        return new ScheduleResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.adapter = this.convertValues(source["adapter"], AdapterManifest);
	        this.decisions = this.convertValues(source["decisions"], CardScheduleDecision);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace terminal {
	
	export class Event {
	    type: string;
	    data: string;
	
	    static createFrom(source: any = {}) {
	        return new Event(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.data = source["data"];
	    }
	}
	export class Info {
	    id: string;
	    command: string;
	    directory: string;
	    state: string;
	
	    static createFrom(source: any = {}) {
	        return new Info(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.command = source["command"];
	        this.directory = source["directory"];
	        this.state = source["state"];
	    }
	}

}

