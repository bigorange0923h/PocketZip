export namespace app {
	
	export class BatchExtractResult {
	    archivePath: string;
	    outputDir: string;
	    success: boolean;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new BatchExtractResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.archivePath = source["archivePath"];
	        this.outputDir = source["outputDir"];
	        this.success = source["success"];
	        this.error = source["error"];
	    }
	}
	export class ExtractStrategy {
	    name: string;
	    outputDir: string;
	    autoRetry: boolean;
	    maxRetries: number;
	    autoOpen: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ExtractStrategy(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.outputDir = source["outputDir"];
	        this.autoRetry = source["autoRetry"];
	        this.maxRetries = source["maxRetries"];
	        this.autoOpen = source["autoOpen"];
	    }
	}

}

export namespace archive {
	
	export class ArchiveEntry {
	    path: string;
	    size: number;
	    isDir: boolean;
	    modified: string;
	
	    static createFrom(source: any = {}) {
	        return new ArchiveEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.size = source["size"];
	        this.isDir = source["isDir"];
	        this.modified = source["modified"];
	    }
	}

}

export namespace history {
	
	export class ExtractHistory {
	    ID: number;
	    ArchivePath: string;
	    OutputDir: string;
	    Success: boolean;
	    UsedPassword: boolean;
	    ErrorMessage: string;
	    // Go type: time
	    CreatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new ExtractHistory(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.ArchivePath = source["ArchivePath"];
	        this.OutputDir = source["OutputDir"];
	        this.Success = source["Success"];
	        this.UsedPassword = source["UsedPassword"];
	        this.ErrorMessage = source["ErrorMessage"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
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

export namespace password {
	
	export class PasswordRecord {
	    ID: number;
	    ArchivePath: string;
	    ArchiveName: string;
	    ArchiveHash: string;
	    EncryptedPassword: number[];
	    SuccessCount: number;
	    // Go type: time
	    LastUsedAt: any;
	    // Go type: time
	    CreatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new PasswordRecord(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.ArchivePath = source["ArchivePath"];
	        this.ArchiveName = source["ArchiveName"];
	        this.ArchiveHash = source["ArchiveHash"];
	        this.EncryptedPassword = source["EncryptedPassword"];
	        this.SuccessCount = source["SuccessCount"];
	        this.LastUsedAt = this.convertValues(source["LastUsedAt"], null);
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
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
	export class PasswordStats {
	    totalRecords: number;
	    totalUsed: number;
	
	    static createFrom(source: any = {}) {
	        return new PasswordStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.totalRecords = source["totalRecords"];
	        this.totalUsed = source["totalUsed"];
	    }
	}

}

