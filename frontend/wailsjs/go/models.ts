export namespace main {
	
	export class Glyph {
	    name: string;
	    glyph: string;
	
	    static createFrom(source: any = {}) {
	        return new Glyph(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.glyph = source["glyph"];
	    }
	}

}

