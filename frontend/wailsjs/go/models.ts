export namespace main {
	
	export class CommandInfo {
	    category: string;
	    icon: string;
	    name: string;
	    description: string;
	    command: string;
	    sudo: boolean;
	    isUser: boolean;
	
	    static createFrom(source: any = {}) {
	        return new CommandInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.category = source["category"];
	        this.icon = source["icon"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.command = source["command"];
	        this.sudo = source["sudo"];
	        this.isUser = source["isUser"];
	    }
	}
	export class DocCategory {
	    name: string;
	    count: number;
	
	    static createFrom(source: any = {}) {
	        return new DocCategory(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.count = source["count"];
	    }
	}
	export class DocFile {
	    name: string;
	    path: string;
	    category: string;
	
	    static createFrom(source: any = {}) {
	        return new DocFile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.category = source["category"];
	    }
	}

}

