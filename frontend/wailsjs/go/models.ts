export namespace api {
	
	export class HistoricalTransaction {
	    t: string;
	    c: number;
	    zdf: number;
	    jlrl: number;
	    hsl: number;
	    qbjlr: number;
	    cddlr: number;
	    cddjlr: number;
	    ddlr: number;
	    ddjlr: number;
	    xdlr: number;
	    xdjlr: number;
	    sdlr: number;
	    sdjlr: number;
	
	    static createFrom(source: any = {}) {
	        return new HistoricalTransaction(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.t = source["t"];
	        this.c = source["c"];
	        this.zdf = source["zdf"];
	        this.jlrl = source["jlrl"];
	        this.hsl = source["hsl"];
	        this.qbjlr = source["qbjlr"];
	        this.cddlr = source["cddlr"];
	        this.cddjlr = source["cddjlr"];
	        this.ddlr = source["ddlr"];
	        this.ddjlr = source["ddjlr"];
	        this.xdlr = source["xdlr"];
	        this.xdjlr = source["xdjlr"];
	        this.sdlr = source["sdlr"];
	        this.sdjlr = source["sdjlr"];
	    }
	}
	export class Index {
	    dm: string;
	    mc: string;
	    jys: string;
	
	    static createFrom(source: any = {}) {
	        return new Index(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.dm = source["dm"];
	        this.mc = source["mc"];
	        this.jys = source["jys"];
	    }
	}
	export class KLineData {
	    d: string;
	    o: number;
	    h: number;
	    l: number;
	    c: number;
	    v: number;
	    e: number;
	    zf: number;
	    hs: number;
	    zd: number;
	    zde: number;
	
	    static createFrom(source: any = {}) {
	        return new KLineData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.d = source["d"];
	        this.o = source["o"];
	        this.h = source["h"];
	        this.l = source["l"];
	        this.c = source["c"];
	        this.v = source["v"];
	        this.e = source["e"];
	        this.zf = source["zf"];
	        this.hs = source["hs"];
	        this.zd = source["zd"];
	        this.zde = source["zde"];
	    }
	}
	export class RealtimeData {
	    fm: number;
	    h: number;
	    hs: number;
	    lb: number;
	    l: number;
	    lt: number;
	    o: number;
	    pe: number;
	    pc: number;
	    p: number;
	    sz: number;
	    cje: number;
	    ud: number;
	    v: number;
	    yc: number;
	    zf: number;
	    zs: number;
	    sjl: number;
	    zdf60: number;
	    zdfnc: number;
	    t: string;
	
	    static createFrom(source: any = {}) {
	        return new RealtimeData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.fm = source["fm"];
	        this.h = source["h"];
	        this.hs = source["hs"];
	        this.lb = source["lb"];
	        this.l = source["l"];
	        this.lt = source["lt"];
	        this.o = source["o"];
	        this.pe = source["pe"];
	        this.pc = source["pc"];
	        this.p = source["p"];
	        this.sz = source["sz"];
	        this.cje = source["cje"];
	        this.ud = source["ud"];
	        this.v = source["v"];
	        this.yc = source["yc"];
	        this.zf = source["zf"];
	        this.zs = source["zs"];
	        this.sjl = source["sjl"];
	        this.zdf60 = source["zdf60"];
	        this.zdfnc = source["zdfnc"];
	        this.t = source["t"];
	    }
	}

}

export namespace data {
	
	export class UpdateStatus {
	    isUpdating: boolean;
	    // Go type: time
	    startTime: any;
	    total: number;
	    completed: number;
	    current: string;
	    progress: number;
	    speed: number;
	    estimateTime: number;
	    errorCount: number;
	    lastError: string;
	    lastUpdateStr: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.isUpdating = source["isUpdating"];
	        this.startTime = this.convertValues(source["startTime"], null);
	        this.total = source["total"];
	        this.completed = source["completed"];
	        this.current = source["current"];
	        this.progress = source["progress"];
	        this.speed = source["speed"];
	        this.estimateTime = source["estimateTime"];
	        this.errorCount = source["errorCount"];
	        this.lastError = source["lastError"];
	        this.lastUpdateStr = source["lastUpdateStr"];
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

export namespace engine {
	
	export class ExecutionRecord {
	    fileName: string;
	    strategyId: number;
	    strategyName: string;
	    // Go type: time
	    executionTime: any;
	    signalCount: number;
	    processedCount: number;
	    totalStocks: number;
	
	    static createFrom(source: any = {}) {
	        return new ExecutionRecord(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.fileName = source["fileName"];
	        this.strategyId = source["strategyId"];
	        this.strategyName = source["strategyName"];
	        this.executionTime = this.convertValues(source["executionTime"], null);
	        this.signalCount = source["signalCount"];
	        this.processedCount = source["processedCount"];
	        this.totalStocks = source["totalStocks"];
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
	export class ExecutionResult {
	    strategyId: number;
	    strategyName: string;
	    // Go type: time
	    executionTime: any;
	    // Go type: time
	    completionTime: any;
	    totalStocks: number;
	    processedStocks: number;
	    signals: types.StockSignal[];
	
	    static createFrom(source: any = {}) {
	        return new ExecutionResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.strategyId = source["strategyId"];
	        this.strategyName = source["strategyName"];
	        this.executionTime = this.convertValues(source["executionTime"], null);
	        this.completionTime = this.convertValues(source["completionTime"], null);
	        this.totalStocks = source["totalStocks"];
	        this.processedStocks = source["processedStocks"];
	        this.signals = this.convertValues(source["signals"], types.StockSignal);
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
	export class ExecutionStatus {
	    status: string;
	    // Go type: time
	    startTime: any;
	    totalStocks: number;
	    processedCount: number;
	    currentStock: string;
	    progress: number;
	    speed: number;
	    estimateTime: number;
	    error: string;
	    strategyId: number;
	
	    static createFrom(source: any = {}) {
	        return new ExecutionStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.startTime = this.convertValues(source["startTime"], null);
	        this.totalStocks = source["totalStocks"];
	        this.processedCount = source["processedCount"];
	        this.currentStock = source["currentStock"];
	        this.progress = source["progress"];
	        this.speed = source["speed"];
	        this.estimateTime = source["estimateTime"];
	        this.error = source["error"];
	        this.strategyId = source["strategyId"];
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
	export class Strategy {
	    id: number;
	    name: string;
	    description: string;
	    filePath: string;
	
	    static createFrom(source: any = {}) {
	        return new Strategy(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.filePath = source["filePath"];
	    }
	}

}

export namespace indicators {
	
	export class MACDResult {
	    DIF: number[];
	    DEA: number[];
	    MACD: number[];
	
	    static createFrom(source: any = {}) {
	        return new MACDResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.DIF = source["DIF"];
	        this.DEA = source["DEA"];
	        this.MACD = source["MACD"];
	    }
	}

}

export namespace main {
	
	export class ExecutionResults {
	    signals: types.StockSignal[];
	    totalStocks: number;
	    status: engine.ExecutionStatus;
	
	    static createFrom(source: any = {}) {
	        return new ExecutionResults(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.signals = this.convertValues(source["signals"], types.StockSignal);
	        this.totalStocks = source["totalStocks"];
	        this.status = this.convertValues(source["status"], engine.ExecutionStatus);
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

export namespace types {
	
	export class StockSignal {
	    code: string;
	    name: string;
	    price: number;
	    change: number;
	    turnover: number;
	    reason: string;
	
	    static createFrom(source: any = {}) {
	        return new StockSignal(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.name = source["name"];
	        this.price = source["price"];
	        this.change = source["change"];
	        this.turnover = source["turnover"];
	        this.reason = source["reason"];
	    }
	}

}

