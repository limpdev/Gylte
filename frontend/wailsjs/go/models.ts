export namespace main {
	
	export class GlyphMatch {
	    id: number;
	    name: string;
	    glyph: string;
	    category?: string;
	    tags?: string;
	    score: number;
	    isFavorite: boolean;
	
	    static createFrom(source: any = {}) {
	        return new GlyphMatch(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.glyph = source["glyph"];
	        this.category = source["category"];
	        this.tags = source["tags"];
	        this.score = source["score"];
	        this.isFavorite = source["isFavorite"];
	    }
	}
	export class SearchResult {
	    glyphs: GlyphMatch[];
	    total: number;
	    searchTime: number;
	    hasMore: boolean;
	    categories?: string[];
	
	    static createFrom(source: any = {}) {
	        return new SearchResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.glyphs = this.convertValues(source["glyphs"], GlyphMatch);
	        this.total = source["total"];
	        this.searchTime = source["searchTime"];
	        this.hasMore = source["hasMore"];
	        this.categories = source["categories"];
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

