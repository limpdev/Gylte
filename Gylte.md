## `app.go`

```go
package main

import (
	"context"
	"database/sql"
	"log"
	"sort"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	_ "modernc.org/sqlite"
)

// App struct
type App struct {
	ctx context.Context
	db  *sql.DB
}

// Glyph struct for database results
type Glyph struct {
	Name  string `json:"name"`
	Glyph string `json:"glyph"`
}

// GlyphMatch represents a glyph with its fuzzy match score
type GlyphMatch struct {
	Glyph
	Score int
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// and the database connection is initialized.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	var err error
	a.db, err = sql.Open("sqlite", "./gylte.db")
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
}

// fuzzyMatch implements fzf-style fuzzy matching
// Returns a score (higher is better) and whether it matches
func fuzzyMatch(pattern, text string) (int, bool) {
	pattern = strings.ToLower(pattern)
	text = strings.ToLower(text)

	if pattern == "" {
		return 0, true
	}

	// Exact substring match gets highest priority
	if strings.Contains(text, pattern) {
		score := 1000
		// Bonus for matches at the beginning
		if strings.HasPrefix(text, pattern) {
			score += 500
		}
		// Bonus for shorter strings (more relevant)
		score -= len(text) - len(pattern)
		return score, true
	}

	// Fuzzy matching: all characters of pattern must appear in order in text
	patternChars := []rune(pattern)
	textChars := []rune(text)

	patternIdx := 0
	score := 0
	consecutiveMatches := 0

	for textIdx, char := range textChars {
		if patternIdx < len(patternChars) && char == patternChars[patternIdx] {
			// Character matches
			patternIdx++
			consecutiveMatches++

			// Bonus for consecutive matches
			score += consecutiveMatches * 2

			// Bonus for matches at word boundaries
			if textIdx == 0 || textChars[textIdx-1] == ' ' || textChars[textIdx-1] == '-' || textChars[textIdx-1] == '_' {
				score += 10
			}
		} else {
			consecutiveMatches = 0
		}
	}

	// All pattern characters must be matched
	if patternIdx == len(patternChars) {
		// Penalty for longer strings
		score -= len(text) - len(pattern)
		return score, true
	}

	return 0, false
}

// GetGlyphs retrieves all glyphs or filters them by a search term with fuzzy matching
func (a *App) GetGlyphs(searchTerm string) ([]Glyph, error) {
	// Always get all glyphs from database
	rows, err := a.db.Query("SELECT name, glyph FROM glyphs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allGlyphs []Glyph
	for rows.Next() {
		var glyph Glyph
		if err := rows.Scan(&glyph.Name, &glyph.Glyph); err != nil {
			return nil, err
		}
		allGlyphs = append(allGlyphs, glyph)
	}

	// If no search term, return all glyphs
	if strings.TrimSpace(searchTerm) == "" {
		return allGlyphs, nil
	}

	// Apply fuzzy matching
	var matches []GlyphMatch
	for _, glyph := range allGlyphs {
		if score, isMatch := fuzzyMatch(searchTerm, glyph.Name); isMatch {
			matches = append(matches, GlyphMatch{
				Glyph: glyph,
				Score: score,
			})
		}
	}

	// Sort by score (highest first)
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Score > matches[j].Score
	})

	// Convert back to []Glyph
	var result []Glyph
	for _, match := range matches {
		result = append(result, match.Glyph)
	}

	return result, nil
}

// CopyToClipboard takes a string and copies it to the user's clipboard.
func (a *App) CopyToClipboard(text string) {
	runtime.ClipboardSetText(a.ctx, text)
}
```

## `db_generator\main.go`

```go
package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	_ "modernc.org/sqlite" // Pure Go SQLite driver
)

// Glyph struct to match the JSON structure
type Glyph struct {
	Name  string `json:"name"`
	Glyph string `json:"glyph"`
}

func main() {
	// Open and read the JSON file
	jsonFile, err := os.Open("../frontend/src/glyphs.json")
	if err != nil {
		log.Fatal("Error opening JSON file:", err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var glyphs []Glyph
	json.Unmarshal(byteValue, &glyphs)

	// Create a new SQLite database file
	// Use "sqlite" instead of "sqlite3" for modernc.org/sqlite
	db, err := sql.Open("sqlite", "./gylte.db")
	if err != nil {
		log.Fatal("Error creating database:", err)
	}
	defer db.Close()

	// Create the glyphs table
	sqlStmt := `
	CREATE TABLE glyphs (id INTEGER NOT NULL PRIMARY KEY, name TEXT, glyph TEXT);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	// Insert unique glyphs into the database
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("INSERT INTO glyphs(name, glyph) VALUES(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	uniqueGlyphs := make(map[string]bool)
	for _, glyph := range glyphs {
		if _, ok := uniqueGlyphs[glyph.Name]; !ok {
			_, err = stmt.Exec(glyph.Name, glyph.Glyph)
			if err != nil {
				log.Fatal(err)
			}
			uniqueGlyphs[glyph.Name] = true
		}
	}
	tx.Commit()

	log.Println("Database 'gylte.db' created and populated successfully.")
}
```

## `frontend\dist\assets\index-CRzIgH-E.js`

```javascript
var Be=Object.defineProperty;var Xe=(t,n,i)=>n in t?Be(t,n,{enumerable:!0,configurable:!0,writable:!0,value:i}):t[n]=i;var ge=(t,n,i)=>Xe(t,typeof n!="symbol"?n+"":n,i);(function(){const n=document.createElement("link").relList;if(n&&n.supports&&n.supports("modulepreload"))return;for(const s of document.querySelectorAll('link[rel="modulepreload"]'))l(s);new MutationObserver(s=>{for(const r of s)if(r.type==="childList")for(const c of r.addedNodes)c.tagName==="LINK"&&c.rel==="modulepreload"&&l(c)}).observe(document,{childList:!0,subtree:!0});function i(s){const r={};return s.integrity&&(r.integrity=s.integrity),s.referrerPolicy&&(r.referrerPolicy=s.referrerPolicy),s.crossOrigin==="use-credentials"?r.credentials="include":s.crossOrigin==="anonymous"?r.credentials="omit":r.credentials="same-origin",r}function l(s){if(s.ep)return;s.ep=!0;const r=i(s);fetch(s.href,r)}})();function ue(){}function pe(t,n){for(const i in n)t[i]=n[i];return t}function Ee(t){return t()}function Se(){return Object.create(null)}function M(t){t.forEach(Ee)}function Oe(t){return typeof t=="function"}function Re(t,n){return t!=t?n==n:t!==n||t&&typeof t=="object"||typeof t=="function"}function Fe(t){return Object.keys(t).length===0}function Ne(t){const n={};for(const i in t)i[0]!=="$"&&(n[i]=t[i]);return n}function a(t,n){t.appendChild(n)}function he(t,n,i){t.insertBefore(n,i||null)}function fe(t){t.parentNode&&t.parentNode.removeChild(t)}function z(t){return document.createElement(t)}function o(t){return document.createElementNS("http://www.w3.org/2000/svg",t)}function _e(t){return document.createTextNode(t)}function x(){return _e(" ")}function se(t,n,i,l){return t.addEventListener(n,i,l),()=>t.removeEventListener(n,i,l)}function e(t,n,i){i==null?t.removeAttribute(n):t.getAttribute(n)!==i&&t.setAttribute(n,i)}function Ge(t,n){for(const i in n)e(t,i,n[i])}function Qe(t){return Array.from(t.childNodes)}function Ie(t,n){n=""+n,t.data!==n&&(t.data=n)}function Le(t,n){t.value=n??""}let oe;function re(t){oe=t}function We(){if(!oe)throw new Error("Function called outside component initialization");return oe}function qe(t){We().$$.on_mount.push(t)}const ie=[],Ve=[];let ae=[];const Pe=[],Ke=Promise.resolve();let we=!1;function De(){we||(we=!0,Ke.then(xe))}function ve(t){ae.push(t)}const be=new Set;let ne=0;function xe(){if(ne!==0)return;const t=oe;do{try{for(;ne<ie.length;){const n=ie[ne];ne++,re(n),He(n.$$)}}catch(n){throw ie.length=0,ne=0,n}for(re(null),ie.length=0,ne=0;Ve.length;)Ve.pop()();for(let n=0;n<ae.length;n+=1){const i=ae[n];be.has(i)||(be.add(i),i())}ae.length=0}while(ie.length);for(;Pe.length;)Pe.pop()();we=!1,be.clear(),re(t)}function He(t){if(t.fragment!==null){t.update(),M(t.before_update);const n=t.dirty;t.dirty=[-1],t.fragment&&t.fragment.p(t.ctx,n),t.after_update.forEach(ve)}}function Je(t){const n=[],i=[];ae.forEach(l=>t.indexOf(l)===-1?n.push(l):i.push(l)),i.forEach(l=>l()),ae=n}const de=new Set;let $;function Ye(){$={r:0,c:[],p:$}}function Ze(){$.r||M($.c),$=$.p}function le(t,n){t&&t.i&&(de.delete(t),t.i(n))}function ye(t,n,i,l){if(t&&t.o){if(de.has(t))return;de.add(t),$.c.push(()=>{de.delete(t),l&&(i&&t.d(1),l())}),t.o(n)}else l&&l()}function je(t){return(t==null?void 0:t.length)!==void 0?t:Array.from(t)}function et(t,n){t.d(1),n.delete(t.key)}function tt(t,n,i,l,s,r,c,u,d,g,S,b){let h=t.length,p=r.length,N=h;const V={};for(;N--;)V[t[N].key]=N;const _=[],v=new Map,k=new Map,P=[];for(N=p;N--;){const w=b(s,r,N),m=i(w);let f=c.get(m);f?P.push(()=>f.p(w,n)):(f=g(m,w),f.c()),v.set(m,_[N]=f),m in V&&k.set(m,Math.abs(N-V[m]))}const C=new Set,A=new Set;function G(w){le(w,1),w.m(u,S),c.set(w.key,w),S=w.first,p--}for(;h&&p;){const w=_[p-1],m=t[h-1],f=w.key,y=m.key;w===m?(S=w.first,h--,p--):v.has(y)?!c.has(f)||C.has(f)?G(w):A.has(y)?h--:k.get(f)>k.get(y)?(A.add(f),G(w)):(C.add(y),h--):(d(m,c),h--)}for(;h--;){const w=t[h];v.has(w.key)||d(w,c)}for(;p;)G(_[p-1]);return M(P),_}function nt(t,n){const i={},l={},s={$$scope:1};let r=t.length;for(;r--;){const c=t[r],u=n[r];if(u){for(const d in c)d in u||(l[d]=1);for(const d in u)s[d]||(i[d]=u[d],s[d]=1);t[r]=u}else for(const d in c)s[d]=1}for(const c in l)c in i||(i[c]=void 0);return i}function it(t){t&&t.c()}function Me(t,n,i){const{fragment:l,after_update:s}=t.$$;l&&l.m(n,i),ve(()=>{const r=t.$$.on_mount.map(Ee).filter(Oe);t.$$.on_destroy?t.$$.on_destroy.push(...r):M(r),t.$$.on_mount=[]}),s.forEach(ve)}function Te(t,n){const i=t.$$;i.fragment!==null&&(Je(i.after_update),M(i.on_destroy),i.fragment&&i.fragment.d(n),i.on_destroy=i.fragment=null,i.ctx=[])}function at(t,n){t.$$.dirty[0]===-1&&(ie.push(t),De(),t.$$.dirty.fill(0)),t.$$.dirty[n/31|0]|=1<<n%31}function $e(t,n,i,l,s,r,c=null,u=[-1]){const d=oe;re(t);const g=t.$$={fragment:null,ctx:[],props:r,update:ue,not_equal:s,bound:Se(),on_mount:[],on_destroy:[],on_disconnect:[],before_update:[],after_update:[],context:new Map(n.context||(d?d.$$.context:[])),callbacks:Se(),dirty:u,skip_bound:!1,root:n.target||d.$$.root};c&&c(g.root);let S=!1;if(g.ctx=i?i(t,n.props||{},(b,h,...p)=>{const N=p.length?p[0]:h;return g.ctx&&s(g.ctx[b],g.ctx[b]=N)&&(!g.skip_bound&&g.bound[b]&&g.bound[b](N),S&&at(t,b)),h}):[],g.update(),S=!0,M(g.before_update),g.fragment=l?l(g.ctx):!1,n.target){if(n.hydrate){const b=Qe(n.target);g.fragment&&g.fragment.l(b),b.forEach(fe)}else g.fragment&&g.fragment.c();n.intro&&le(t.$$.fragment),Me(t,n.target,n.anchor),xe()}re(d)}class Ue{constructor(){ge(this,"$$");ge(this,"$$set")}$destroy(){Te(this,1),this.$destroy=ue}$on(n,i){if(!Oe(i))return ue;const l=this.$$.callbacks[n]||(this.$$.callbacks[n]=[]);return l.push(i),()=>{const s=l.indexOf(i);s!==-1&&l.splice(s,1)}}$set(n){this.$$set&&!Fe(n)&&(this.$$.skip_bound=!0,this.$$set(n),this.$$.skip_bound=!1)}}const lt="4";typeof window<"u"&&(window.__svelte||(window.__svelte={v:new Set})).v.add(lt);function st(t){let n,i,l,s,r,c,u,d,g,S,b,h,p,N,V,_,v,k,P,C,A,G,w,m,f,y,I,U,B,X,F,E,Q,W,q,K,O,D,H,J,Y,R,T,Z,ee,te,me=[{xmlns:"http://www.w3.org/2000/svg"},{width:t[0]},{height:t[1]},{viewBox:"0 0 24 24"},t[3]],ce={};for(let L=0;L<me.length;L+=1)ce=pe(ce,me[L]);return{c(){n=o("svg"),i=o("rect"),l=o("animate"),s=o("animate"),r=o("animate"),c=o("animate"),u=o("rect"),d=o("animate"),g=o("animate"),S=o("animate"),b=o("animate"),h=o("rect"),p=o("animate"),N=o("animate"),V=o("animate"),_=o("animate"),v=o("rect"),k=o("animate"),P=o("animate"),C=o("animate"),A=o("animate"),G=o("rect"),w=o("animate"),m=o("animate"),f=o("animate"),y=o("animate"),I=o("rect"),U=o("animate"),B=o("animate"),X=o("animate"),F=o("animate"),E=o("rect"),Q=o("animate"),W=o("animate"),q=o("animate"),K=o("animate"),O=o("rect"),D=o("animate"),H=o("animate"),J=o("animate"),Y=o("animate"),R=o("rect"),T=o("animate"),Z=o("animate"),ee=o("animate"),te=o("animate"),e(l,"id","SVGzjrPLenI"),e(l,"attributeName","x"),e(l,"begin","0;SVGXAURnSRI.end+0.25s"),e(l,"dur","0.75s"),e(l,"values","1;4;1"),e(s,"attributeName","y"),e(s,"begin","0;SVGXAURnSRI.end+0.25s"),e(s,"dur","0.75s"),e(s,"values","1;4;1"),e(r,"attributeName","width"),e(r,"begin","0;SVGXAURnSRI.end+0.25s"),e(r,"dur","0.75s"),e(r,"values","7.33;1.33;7.33"),e(c,"attributeName","height"),e(c,"begin","0;SVGXAURnSRI.end+0.25s"),e(c,"dur","0.75s"),e(c,"values","7.33;1.33;7.33"),e(i,"width","7.33"),e(i,"height","7.33"),e(i,"x","1"),e(i,"y","1"),e(i,"fill",t[2]),e(d,"attributeName","x"),e(d,"begin","SVGzjrPLenI.begin+0.125s"),e(d,"dur","0.75s"),e(d,"values","8.33;11.33;8.33"),e(g,"attributeName","y"),e(g,"begin","SVGzjrPLenI.begin+0.125s"),e(g,"dur","0.75s"),e(g,"values","1;4;1"),e(S,"attributeName","width"),e(S,"begin","SVGzjrPLenI.begin+0.125s"),e(S,"dur","0.75s"),e(S,"values","7.33;1.33;7.33"),e(b,"attributeName","height"),e(b,"begin","SVGzjrPLenI.begin+0.125s"),e(b,"dur","0.75s"),e(b,"values","7.33;1.33;7.33"),e(u,"width","7.33"),e(u,"height","7.33"),e(u,"x","8.33"),e(u,"y","1"),e(u,"fill",t[2]),e(p,"attributeName","x"),e(p,"begin","SVGzjrPLenI.begin+0.125s"),e(p,"dur","0.75s"),e(p,"values","1;4;1"),e(N,"attributeName","y"),e(N,"begin","SVGzjrPLenI.begin+0.125s"),e(N,"dur","0.75s"),e(N,"values","8.33;11.33;8.33"),e(V,"attributeName","width"),e(V,"begin","SVGzjrPLenI.begin+0.125s"),e(V,"dur","0.75s"),e(V,"values","7.33;1.33;7.33"),e(_,"attributeName","height"),e(_,"begin","SVGzjrPLenI.begin+0.125s"),e(_,"dur","0.75s"),e(_,"values","7.33;1.33;7.33"),e(h,"width","7.33"),e(h,"height","7.33"),e(h,"x","1"),e(h,"y","8.33"),e(h,"fill",t[2]),e(k,"attributeName","x"),e(k,"begin","SVGzjrPLenI.begin+0.25s"),e(k,"dur","0.75s"),e(k,"values","15.66;18.66;15.66"),e(P,"attributeName","y"),e(P,"begin","SVGzjrPLenI.begin+0.25s"),e(P,"dur","0.75s"),e(P,"values","1;4;1"),e(C,"attributeName","width"),e(C,"begin","SVGzjrPLenI.begin+0.25s"),e(C,"dur","0.75s"),e(C,"values","7.33;1.33;7.33"),e(A,"attributeName","height"),e(A,"begin","SVGzjrPLenI.begin+0.25s"),e(A,"dur","0.75s"),e(A,"values","7.33;1.33;7.33"),e(v,"width","7.33"),e(v,"height","7.33"),e(v,"x","15.66"),e(v,"y","1"),e(v,"fill",t[2]),e(w,"attributeName","x"),e(w,"begin","SVGzjrPLenI.begin+0.25s"),e(w,"dur","0.75s"),e(w,"values","8.33;11.33;8.33"),e(m,"attributeName","y"),e(m,"begin","SVGzjrPLenI.begin+0.25s"),e(m,"dur","0.75s"),e(m,"values","8.33;11.33;8.33"),e(f,"attributeName","width"),e(f,"begin","SVGzjrPLenI.begin+0.25s"),e(f,"dur","0.75s"),e(f,"values","7.33;1.33;7.33"),e(y,"attributeName","height"),e(y,"begin","SVGzjrPLenI.begin+0.25s"),e(y,"dur","0.75s"),e(y,"values","7.33;1.33;7.33"),e(G,"width","7.33"),e(G,"height","7.33"),e(G,"x","8.33"),e(G,"y","8.33"),e(G,"fill",t[2]),e(U,"attributeName","x"),e(U,"begin","SVGzjrPLenI.begin+0.25s"),e(U,"dur","0.75s"),e(U,"values","1;4;1"),e(B,"attributeName","y"),e(B,"begin","SVGzjrPLenI.begin+0.25s"),e(B,"dur","0.75s"),e(B,"values","15.66;18.66;15.66"),e(X,"attributeName","width"),e(X,"begin","SVGzjrPLenI.begin+0.25s"),e(X,"dur","0.75s"),e(X,"values","7.33;1.33;7.33"),e(F,"attributeName","height"),e(F,"begin","SVGzjrPLenI.begin+0.25s"),e(F,"dur","0.75s"),e(F,"values","7.33;1.33;7.33"),e(I,"width","7.33"),e(I,"height","7.33"),e(I,"x","1"),e(I,"y","15.66"),e(I,"fill",t[2]),e(Q,"attributeName","x"),e(Q,"begin","SVGzjrPLenI.begin+0.375s"),e(Q,"dur","0.75s"),e(Q,"values","15.66;18.66;15.66"),e(W,"attributeName","y"),e(W,"begin","SVGzjrPLenI.begin+0.375s"),e(W,"dur","0.75s"),e(W,"values","8.33;11.33;8.33"),e(q,"attributeName","width"),e(q,"begin","SVGzjrPLenI.begin+0.375s"),e(q,"dur","0.75s"),e(q,"values","7.33;1.33;7.33"),e(K,"attributeName","height"),e(K,"begin","SVGzjrPLenI.begin+0.375s"),e(K,"dur","0.75s"),e(K,"values","7.33;1.33;7.33"),e(E,"width","7.33"),e(E,"height","7.33"),e(E,"x","15.66"),e(E,"y","8.33"),e(E,"fill",t[2]),e(D,"attributeName","x"),e(D,"begin","SVGzjrPLenI.begin+0.375s"),e(D,"dur","0.75s"),e(D,"values","8.33;11.33;8.33"),e(H,"attributeName","y"),e(H,"begin","SVGzjrPLenI.begin+0.375s"),e(H,"dur","0.75s"),e(H,"values","15.66;18.66;15.66"),e(J,"attributeName","width"),e(J,"begin","SVGzjrPLenI.begin+0.375s"),e(J,"dur","0.75s"),e(J,"values","7.33;1.33;7.33"),e(Y,"attributeName","height"),e(Y,"begin","SVGzjrPLenI.begin+0.375s"),e(Y,"dur","0.75s"),e(Y,"values","7.33;1.33;7.33"),e(O,"width","7.33"),e(O,"height","7.33"),e(O,"x","8.33"),e(O,"y","15.66"),e(O,"fill",t[2]),e(T,"id","SVGXAURnSRI"),e(T,"attributeName","x"),e(T,"begin","SVGzjrPLenI.begin+0.5s"),e(T,"dur","0.75s"),e(T,"values","15.66;18.66;15.66"),e(Z,"attributeName","y"),e(Z,"begin","SVGzjrPLenI.begin+0.5s"),e(Z,"dur","0.75s"),e(Z,"values","15.66;18.66;15.66"),e(ee,"attributeName","width"),e(ee,"begin","SVGzjrPLenI.begin+0.5s"),e(ee,"dur","0.75s"),e(ee,"values","7.33;1.33;7.33"),e(te,"attributeName","height"),e(te,"begin","SVGzjrPLenI.begin+0.5s"),e(te,"dur","0.75s"),e(te,"values","7.33;1.33;7.33"),e(R,"width","7.33"),e(R,"height","7.33"),e(R,"x","15.66"),e(R,"y","15.66"),e(R,"fill",t[2]),Ge(n,ce)},m(L,j){he(L,n,j),a(n,i),a(i,l),a(i,s),a(i,r),a(i,c),a(n,u),a(u,d),a(u,g),a(u,S),a(u,b),a(n,h),a(h,p),a(h,N),a(h,V),a(h,_),a(n,v),a(v,k),a(v,P),a(v,C),a(v,A),a(n,G),a(G,w),a(G,m),a(G,f),a(G,y),a(n,I),a(I,U),a(I,B),a(I,X),a(I,F),a(n,E),a(E,Q),a(E,W),a(E,q),a(E,K),a(n,O),a(O,D),a(O,H),a(O,J),a(O,Y),a(n,R),a(R,T),a(R,Z),a(R,ee),a(R,te)},p(L,[j]){j&4&&e(i,"fill",L[2]),j&4&&e(u,"fill",L[2]),j&4&&e(h,"fill",L[2]),j&4&&e(v,"fill",L[2]),j&4&&e(G,"fill",L[2]),j&4&&e(I,"fill",L[2]),j&4&&e(E,"fill",L[2]),j&4&&e(O,"fill",L[2]),j&4&&e(R,"fill",L[2]),Ge(n,ce=nt(me,[{xmlns:"http://www.w3.org/2000/svg"},j&1&&{width:L[0]},j&2&&{height:L[1]},{viewBox:"0 0 24 24"},j&8&&L[3]]))},i:ue,o:ue,d(L){L&&fe(n)}}}function rt(t,n,i){let{width:l="24"}=n,{height:s="24"}=n,{fill:r="{fill}"}=n,{customClass:c="animatedToast"}=n;return t.$$set=u=>{i(3,n=pe(pe({},n),Ne(u))),"width"in u&&i(0,l=u.width),"height"in u&&i(1,s=u.height),"fill"in u&&i(2,r=u.fill),"customClass"in u&&i(4,c=u.customClass)},n=Ne(n),[l,s,r,n,c]}class ut extends Ue{constructor(n){super(),$e(this,n,rt,st,Re,{width:0,height:1,fill:2,customClass:4})}}function ot(t){return window.go.main.App.CopyToClipboard(t)}function ze(t){return window.go.main.App.GetGlyphs(t)}function ft(){window.runtime.WindowMinimise()}function ct(){window.runtime.Quit()}function ke(t,n,i){const l=t.slice();return l[8]=n[i],l[10]=i,l}function Ce(t,n){let i,l,s=n[8].glyph+"",r,c,u,d=n[8].name+"",g,S,b,h,p;function N(){return n[5](n[8])}function V(..._){return n[6](n[8],..._)}return{key:t,first:null,c(){i=z("div"),l=z("span"),r=_e(s),c=x(),u=z("span"),g=_e(d),S=x(),e(l,"class","glyph-icon"),e(u,"class","glyph-name"),e(i,"class","glyph-card"),e(i,"title",b=`Click to copy "${n[8].glyph}"`),e(i,"role","button"),e(i,"tabindex","0"),this.first=i},m(_,v){he(_,i,v),a(i,l),a(l,r),a(i,c),a(i,u),a(u,g),a(i,S),h||(p=[se(i,"click",N),se(i,"keydown",V)],h=!0)},p(_,v){n=_,v&2&&s!==(s=n[8].glyph+"")&&Ie(r,s),v&2&&d!==(d=n[8].name+"")&&Ie(g,d),v&2&&b!==(b=`Click to copy "${n[8].glyph}"`)&&e(i,"title",b)},d(_){_&&fe(i),h=!1,M(p)}}}function Ae(t){let n,i,l;return i=new ut({props:{width:"12",height:"12",fill:"#45a847",class:"animatedToast"}}),{c(){n=z("div"),it(i.$$.fragment),e(n,"class","toast")},m(s,r){he(s,n,r),Me(i,n,null),l=!0},i(s){l||(le(i.$$.fragment,s),l=!0)},o(s){ye(i.$$.fragment,s),l=!1},d(s){s&&fe(n),Te(i)}}}function dt(t){let n,i,l,s,r,c,u,d,g,S,b,h,p,N,V,_=[],v=new Map,k,P,C,A,G=je(t[1]);const w=f=>f[10];for(let f=0;f<G.length;f+=1){let y=ke(t,G,f),I=w(y);v.set(I,_[f]=Ce(I,y))}let m=t[2]&&Ae();return{c(){n=z("div"),i=z("div"),l=z("div"),l.textContent="",s=x(),r=z("div"),c=x(),u=z("div"),d=z("button"),d.textContent="-",g=x(),S=z("button"),S.textContent="×",b=x(),h=z("div"),p=z("input"),N=x(),V=z("div");for(let f=0;f<_.length;f+=1)_[f].c();k=x(),m&&m.c(),e(l,"class","title"),e(r,"class","spacer"),e(u,"class","window-controls"),e(i,"class","title-bar draggable"),e(p,"type","text"),e(p,"autocomplete","on"),e(p,"class","search-input"),e(p,"placeholder",""),e(h,"class","search-container draggable"),e(V,"class","glyph-grid draggable"),e(n,"id","app")},m(f,y){he(f,n,y),a(n,i),a(i,l),a(i,s),a(i,r),a(i,c),a(i,u),a(u,d),a(u,g),a(u,S),a(n,b),a(n,h),a(h,p),Le(p,t[0]),a(n,N),a(n,V);for(let I=0;I<_.length;I+=1)_[I]&&_[I].m(V,null);a(n,k),m&&m.m(n,null),P=!0,C||(A=[se(d,"click",ft),se(S,"click",ct),se(p,"input",t[4])],C=!0)},p(f,[y]){y&1&&p.value!==f[0]&&Le(p,f[0]),y&10&&(G=je(f[1]),_=tt(_,y,w,1,f,G,v,V,et,Ce,null,ke)),f[2]?m?y&4&&le(m,1):(m=Ae(),m.c(),le(m,1),m.m(n,null)):m&&(Ye(),ye(m,1,1,()=>{m=null}),Ze())},i(f){P||(le(m),P=!0)},o(f){ye(m),P=!1},d(f){f&&fe(n);for(let y=0;y<_.length;y+=1)_[y].d();m&&m.d(),C=!1,M(A)}}}function ht(t,n,i){let l="",s=[],r=!1;qe(async()=>{i(1,s=await ze(""))});const c=async()=>{i(1,s=await ze(l))},u=b=>{ot(b),i(2,r=!0),setTimeout(()=>{i(2,r=!1)},1900)};function d(){l=this.value,i(0,l)}const g=b=>u(b.glyph),S=(b,h)=>h.key==="Enter"&&u(b.glyph);return setTimeout(c,300),[l,s,r,u,d,g,S]}class mt extends Ue{constructor(n){super(),$e(this,n,ht,dt,Re,{})}}new mt({target:document.getElementById("app")});
```

## `frontend\src\App.svelte`

```
<script lang="ts">
    import { onMount } from "svelte";
    import AniToast from "./comps/aniToast.svelte";
    import { CopyToClipboard, SearchGlyphs, GetAllGlyphs } from "../wailsjs/go/main/App";
    import { WindowMinimise, Quit } from "../wailsjs/runtime";

    // Define the type for a single glyph object
    interface Glyph {
        name: string;
        glyph: string;
    }

    let searchTerm = "";
    let glyphs: Glyph[] = [];
    let toastVisible = false;
    let searchTimeout: NodeJS.Timeout;
    let isLoading = true;

    // Load all glyphs on initial render
    onMount(async () => {
        try {
            // The Go backend automatically loads the embedded JSON on startup
            // So we just need to get all glyphs for the initial display
            filteredGlyphs = await GetAllGlyphs();
            isLoading = false;
        } catch (error) {
            console.error("Failed to load glyphs from Go backend:", error);
            isLoading = false;
        }
    });

    // Debounced search function for better performance
    const performSearch = async (query: string) => {
        try {
            filteredGlyphs = await SearchGlyphs(query);
        } catch (error) {
            console.error("Search failed:", error);
        }
    };

    // Reactive statement with debouncing for search
    $: {
        // Clear previous timeout
        if (searchTimeout) {
            clearTimeout(searchTimeout);
        }
        
        // Debounce search by 150ms to avoid excessive backend calls
        if (!isLoading) {
            searchTimeout = setTimeout(() => {
                performSearch(searchTerm);
            }, 150);
        }
    }

    const handleGlyphClick = (glyph: string) => {
        CopyToClipboard(glyph);
        // Show "Copied!" toast message
        toastVisible = true;
        setTimeout(() => {
            toastVisible = false;
        }, 1900);
    };
</script>

<div id="app">
    <!-- Custom Title Bar -->
    <div class="title-bar draggable">
        <div class="title"></div>
        <div class="spacer"></div>
        <div class="window-controls">
            <button on:click={WindowMinimise}>-</button>
            <button on:click={Quit}>Ã—</button>
        </div>
    </div>

    <!-- Search Input -->
    <div class="search-container draggable">
        <input 
            type="text" 
            autocomplete="on" 
            class="search-input" 
            placeholder="" 
            bind:value={searchTerm}
            disabled={isLoading}
        />
    </div>

    <!-- Loading State -->
    {#if isLoading}
        <div class="loading">Loading glyphs...</div>
    {:else}
        <!-- Glyph Grid -->
        <div class="glyph-grid draggable">
            {#each filteredGlyphs as item (item.name)}
                <div
                    class="glyph-card"
                    title={`Click to copy "${item.glyph}"`}
                    on:click={() => handleGlyphClick(item.glyph)}
                    on:keydown={(e) => e.key === "Enter" && handleGlyphClick(item.glyph)}
                    role="button"
                    tabindex="0"
                >
                    <span class="glyph-icon">{item.glyph}</span>
                    <span class="glyph-name">{item.name}</span>
                </div>
            {/each}
        </div>
    {/if}

    <!-- "Copied!" Toast Notification -->
    {#if toastVisible}
        <div class="toast">
            <AniToast width="12" height="12" fill="#45a847" class="animatedToast" />
        </div>
    {/if}
</div>

<style>
    .loading {
        display: flex;
        justify-content: center;
        align-items: center;
        height: 200px;
        color: #888;
        font-size: 14px;
        background-color: #12121290;
    }
    
    /* Your existing CSS styles go here */
</style>
```

## `frontend\src\comps\aniToast.svelte`

```
<script>
    export let width = "24";
    export let height = "24";
    export let fill = "{fill}";
    export let customClass = "animatedToast"; // Use a different name to avoid conflict with the class attribute
</script>

<svg
    xmlns="http://www.w3.org/2000/svg"
    width={width}
    height={height}
    viewBox="0 0 24 24"
    {...$$props}
>
    <rect width="7.33" height="7.33" x="1" y="1" {fill}>
        <animate
            id="SVGzjrPLenI"
            attributeName="x"
            begin="0;SVGXAURnSRI.end+0.25s"
            dur="0.75s"
            values="1;4;1"
        />
        <animate
            attributeName="y"
            begin="0;SVGXAURnSRI.end+0.25s"
            dur="0.75s"
            values="1;4;1"
        />
        <animate
            attributeName="width"
            begin="0;SVGXAURnSRI.end+0.25s"
            dur="0.75s"
            values="7.33;1.33;7.33"
        />
        <animate
            attributeName="height"
            begin="0;SVGXAURnSRI.end+0.25s"
            dur="0.75s"
            values="7.33;1.33;7.33"
        />
    </rect>
    <rect width="7.33" height="7.33" x="8.33" y="1" {fill}>
        <animate
            attributeName="x"
            begin="SVGzjrPLenI.begin+0.125s"
            dur="0.75s"
            values="8.33;11.33;8.33"
        />
        <animate
            attributeName="y"
            begin="SVGzjrPLenI.begin+0.125s"
            dur="0.75s"
            values="1;4;1"
        />
        <animate
            attributeName="width"
            begin="SVGzjrPLenI.begin+0.125s"
            dur="0.75s"
            values="7.33;1.33;7.33"
        />
        <animate
            attributeName="height"
            begin="SVGzjrPLenI.begin+0.125s"
            dur="0.75s"
            values="7.33;1.33;7.33"
        />
    </rect>
    <rect width="7.33" height="7.33" x="1" y="8.33" {fill}>
        <animate
            attributeName="x"
            begin="SVGzjrPLenI.begin+0.125s"
            dur="0.75s"
            values="1;4;1"
        />
        <animate
            attributeName="y"
            begin="SVGzjrPLenI.begin+0.125s"
            dur="0.75s"
            values="8.33;11.33;8.33"
        />
        <animate
            attributeName="width"
            begin="SVGzjrPLenI.begin+0.125s"
            dur="0.75s"
            values="7.33;1.33;7.33"
        />
        <animate
            attributeName="height"
            begin="SVGzjrPLenI.begin+0.125s"
            dur="0.75s"
            values="7.33;1.33;7.33"
        />
    </rect>
    <rect width="7.33" height="7.33" x="15.66" y="1" {fill}>
        <animate
            attributeName="x"
            begin="SVGzjrPLenI.begin+0.25s"
            dur="0.75s"
            values="15.66;18.66;15.66"
        />
        <animate
            attributeName="y"
            begin="SVGzjrPLenI.begin+0.25s"
            dur="0.75s"
            values="1;4;1"
        />
        <animate
            attributeName="width"
            begin="SVGzjrPLenI.begin+0.25s"
            dur="0.75s"
            values="7.33;1.33;7.33"
        />
        <animate
            attributeName="height"
            begin="SVGzjrPLenI.begin+0.25s"
            dur="0.75s"
            values="7.33;1.33;7.33"
        />
    </rect>
    <rect width="7.33" height="7.33" x="8.33" y="8.33" {fill}>
        <animate
            attributeName="x"
            begin="SVGzjrPLenI.begin+0.25s"
            dur="0.75s"
            values="8.33;11.33;8.33"
        />
        <animate
            attributeName="y"
            begin="SVGzjrPLenI.begin+0.25s"
            dur="0.75s"
            values="8.33;11.33;8.33"
        />
        <animate
            attributeName="width"
            begin="SVGzjrPLenI.begin+0.25s"
            dur="0.75s"
            values="7.33;1.33;7.33"
        />
        <animate
            attributeName="height"
            begin="SVGzjrPLenI.begin+0.25s"
            dur="0.75s"
            values="7.33;1.33;7.33"
        />
    </rect>
    <rect width="7.33" height="7.33" x="1" y="15.66" {fill}>
        <animate
            attributeName="x"
            begin="SVGzjrPLenI.begin+0.25s"
            dur="0.75s"
            values="1;4;1"
        />
        <animate
            attributeName="y"
            begin="SVGzjrPLenI.begin+0.25s"
            dur="0.75s"
            values="15.66;18.66;15.66"
        />
        <animate
            attributeName="width"
            begin="SVGzjrPLenI.begin+0.25s"
            dur="0.75s"
            values="7.33;1.33;7.33"
        />
        <animate
            attributeName="height"
            begin="SVGzjrPLenI.begin+0.25s"
            dur="0.75s"
            values="7.33;1.33;7.33"
        />
    </rect>
    <rect width="7.33" height="7.33" x="15.66" y="8.33" {fill}>
        <animate
            attributeName="x"
            begin="SVGzjrPLenI.begin+0.375s"
            dur="0.75s"
            values="15.66;18.66;15.66"
        />
        <animate
            attributeName="y"
            begin="SVGzjrPLenI.begin+0.375s"
            dur="0.75s"
            values="8.33;11.33;8.33"
        />
        <animate
            attributeName="width"
            begin="SVGzjrPLenI.begin+0.375s"
            dur="0.75s"
            values="7.33;1.33;7.33"
        />
        <animate
            attributeName="height"
            begin="SVGzjrPLenI.begin+0.375s"
            dur="0.75s"
            values="7.33;1.33;7.33"
        />
    </rect>
    <rect width="7.33" height="7.33" x="8.33" y="15.66" {fill}>
        <animate
            attributeName="x"
            begin="SVGzjrPLenI.begin+0.375s"
            dur="0.75s"
            values="8.33;11.33;8.33"
        />
        <animate
            attributeName="y"
            begin="SVGzjrPLenI.begin+0.375s"
            dur="0.75s"
            values="15.66;18.66;15.66"
        />
        <animate
            attributeName="width"
            begin="SVGzjrPLenI.begin+0.375s"
            dur="0.75s"
            values="7.33;1.33;7.33"
        />
        <animate
            attributeName="height"
            begin="SVGzjrPLenI.begin+0.375s"
            dur="0.75s"
            values="7.33;1.33;7.33"
        />
    </rect>
    <rect width="7.33" height="7.33" x="15.66" y="15.66" {fill}>
        <animate
            id="SVGXAURnSRI"
            attributeName="x"
            begin="SVGzjrPLenI.begin+0.5s"
            dur="0.75s"
            values="15.66;18.66;15.66"
        />
        <animate
            attributeName="y"
            begin="SVGzjrPLenI.begin+0.5s"
            dur="0.75s"
            values="15.66;18.66;15.66"
        />
        <animate
            attributeName="width"
            begin="SVGzjrPLenI.begin+0.5s"
            dur="0.75s"
            values="7.33;1.33;7.33"
        />
        <animate
            attributeName="height"
            begin="SVGzjrPLenI.begin+0.5s"
            dur="0.75s"
            values="7.33;1.33;7.33"
        />
    </rect>
</svg>
```

## `frontend\src\main.ts`

```typescript
import "./style.css"; // Keep your existing CSS file name
import App from "./App.svelte";

const app = new App({
    target: document.getElementById("app")!,
});

export default app;
```

## `frontend\src\vite-env.d.ts`

```typescript
/// <reference types="svelte" />
/// <reference types="vite/client" />
```

## `frontend\svelte.config.js`

```javascript
import { vitePreprocess } from "@sveltejs/vite-plugin-svelte";

export default {
    // Consult https://svelte.dev/docs#compile-time-svelte-preprocess
    // for more information about preprocessors
    preprocess: vitePreprocess(),
};
```

## `frontend\vite.config.js`

```javascript
import {defineConfig} from 'vite'
import {svelte} from '@sveltejs/vite-plugin-svelte'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [svelte()]
})
```

## `frontend\wailsjs\go\main\App.d.ts`

```typescript
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {main} from '../models';

export function CopyToClipboard(arg1:string):Promise<void>;

export function GetGlyphs(arg1:string):Promise<Array<main.Glyph>>;
```

## `frontend\wailsjs\go\main\App.js`

```javascript
// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

export function CopyToClipboard(arg1) {
  return window['go']['main']['App']['CopyToClipboard'](arg1);
}

export function GetGlyphs(arg1) {
  return window['go']['main']['App']['GetGlyphs'](arg1);
}
```

## `frontend\wailsjs\go\models.ts`

```typescript
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

```

## `frontend\wailsjs\runtime\runtime.d.ts`

```typescript
/*
 _       __      _ __
| |     / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/

export interface Position {
    x: number;
    y: number;
}

export interface Size {
    w: number;
    h: number;
}

export interface Screen {
    isCurrent: boolean;
    isPrimary: boolean;
    width : number
    height : number
}

// Environment information such as platform, buildtype, ...
export interface EnvironmentInfo {
    buildType: string;
    platform: string;
    arch: string;
}

// [EventsEmit](https://wails.io/docs/reference/runtime/events#eventsemit)
// emits the given event. Optional data may be passed with the event.
// This will trigger any event listeners.
export function EventsEmit(eventName: string, ...data: any): void;

// [EventsOn](https://wails.io/docs/reference/runtime/events#eventson) sets up a listener for the given event name.
export function EventsOn(eventName: string, callback: (...data: any) => void): () => void;

// [EventsOnMultiple](https://wails.io/docs/reference/runtime/events#eventsonmultiple)
// sets up a listener for the given event name, but will only trigger a given number times.
export function EventsOnMultiple(eventName: string, callback: (...data: any) => void, maxCallbacks: number): () => void;

// [EventsOnce](https://wails.io/docs/reference/runtime/events#eventsonce)
// sets up a listener for the given event name, but will only trigger once.
export function EventsOnce(eventName: string, callback: (...data: any) => void): () => void;

// [EventsOff](https://wails.io/docs/reference/runtime/events#eventsoff)
// unregisters the listener for the given event name.
export function EventsOff(eventName: string, ...additionalEventNames: string[]): void;

// [EventsOffAll](https://wails.io/docs/reference/runtime/events#eventsoffall)
// unregisters all listeners.
export function EventsOffAll(): void;

// [LogPrint](https://wails.io/docs/reference/runtime/log#logprint)
// logs the given message as a raw message
export function LogPrint(message: string): void;

// [LogTrace](https://wails.io/docs/reference/runtime/log#logtrace)
// logs the given message at the `trace` log level.
export function LogTrace(message: string): void;

// [LogDebug](https://wails.io/docs/reference/runtime/log#logdebug)
// logs the given message at the `debug` log level.
export function LogDebug(message: string): void;

// [LogError](https://wails.io/docs/reference/runtime/log#logerror)
// logs the given message at the `error` log level.
export function LogError(message: string): void;

// [LogFatal](https://wails.io/docs/reference/runtime/log#logfatal)
// logs the given message at the `fatal` log level.
// The application will quit after calling this method.
export function LogFatal(message: string): void;

// [LogInfo](https://wails.io/docs/reference/runtime/log#loginfo)
// logs the given message at the `info` log level.
export function LogInfo(message: string): void;

// [LogWarning](https://wails.io/docs/reference/runtime/log#logwarning)
// logs the given message at the `warning` log level.
export function LogWarning(message: string): void;

// [WindowReload](https://wails.io/docs/reference/runtime/window#windowreload)
// Forces a reload by the main application as well as connected browsers.
export function WindowReload(): void;

// [WindowReloadApp](https://wails.io/docs/reference/runtime/window#windowreloadapp)
// Reloads the application frontend.
export function WindowReloadApp(): void;

// [WindowSetAlwaysOnTop](https://wails.io/docs/reference/runtime/window#windowsetalwaysontop)
// Sets the window AlwaysOnTop or not on top.
export function WindowSetAlwaysOnTop(b: boolean): void;

// [WindowSetSystemDefaultTheme](https://wails.io/docs/next/reference/runtime/window#windowsetsystemdefaulttheme)
// *Windows only*
// Sets window theme to system default (dark/light).
export function WindowSetSystemDefaultTheme(): void;

// [WindowSetLightTheme](https://wails.io/docs/next/reference/runtime/window#windowsetlighttheme)
// *Windows only*
// Sets window to light theme.
export function WindowSetLightTheme(): void;

// [WindowSetDarkTheme](https://wails.io/docs/next/reference/runtime/window#windowsetdarktheme)
// *Windows only*
// Sets window to dark theme.
export function WindowSetDarkTheme(): void;

// [WindowCenter](https://wails.io/docs/reference/runtime/window#windowcenter)
// Centers the window on the monitor the window is currently on.
export function WindowCenter(): void;

// [WindowSetTitle](https://wails.io/docs/reference/runtime/window#windowsettitle)
// Sets the text in the window title bar.
export function WindowSetTitle(title: string): void;

// [WindowFullscreen](https://wails.io/docs/reference/runtime/window#windowfullscreen)
// Makes the window full screen.
export function WindowFullscreen(): void;

// [WindowUnfullscreen](https://wails.io/docs/reference/runtime/window#windowunfullscreen)
// Restores the previous window dimensions and position prior to full screen.
export function WindowUnfullscreen(): void;

// [WindowIsFullscreen](https://wails.io/docs/reference/runtime/window#windowisfullscreen)
// Returns the state of the window, i.e. whether the window is in full screen mode or not.
export function WindowIsFullscreen(): Promise<boolean>;

// [WindowSetSize](https://wails.io/docs/reference/runtime/window#windowsetsize)
// Sets the width and height of the window.
export function WindowSetSize(width: number, height: number): void;

// [WindowGetSize](https://wails.io/docs/reference/runtime/window#windowgetsize)
// Gets the width and height of the window.
export function WindowGetSize(): Promise<Size>;

// [WindowSetMaxSize](https://wails.io/docs/reference/runtime/window#windowsetmaxsize)
// Sets the maximum window size. Will resize the window if the window is currently larger than the given dimensions.
// Setting a size of 0,0 will disable this constraint.
export function WindowSetMaxSize(width: number, height: number): void;

// [WindowSetMinSize](https://wails.io/docs/reference/runtime/window#windowsetminsize)
// Sets the minimum window size. Will resize the window if the window is currently smaller than the given dimensions.
// Setting a size of 0,0 will disable this constraint.
export function WindowSetMinSize(width: number, height: number): void;

// [WindowSetPosition](https://wails.io/docs/reference/runtime/window#windowsetposition)
// Sets the window position relative to the monitor the window is currently on.
export function WindowSetPosition(x: number, y: number): void;

// [WindowGetPosition](https://wails.io/docs/reference/runtime/window#windowgetposition)
// Gets the window position relative to the monitor the window is currently on.
export function WindowGetPosition(): Promise<Position>;

// [WindowHide](https://wails.io/docs/reference/runtime/window#windowhide)
// Hides the window.
export function WindowHide(): void;

// [WindowShow](https://wails.io/docs/reference/runtime/window#windowshow)
// Shows the window, if it is currently hidden.
export function WindowShow(): void;

// [WindowMaximise](https://wails.io/docs/reference/runtime/window#windowmaximise)
// Maximises the window to fill the screen.
export function WindowMaximise(): void;

// [WindowToggleMaximise](https://wails.io/docs/reference/runtime/window#windowtogglemaximise)
// Toggles between Maximised and UnMaximised.
export function WindowToggleMaximise(): void;

// [WindowUnmaximise](https://wails.io/docs/reference/runtime/window#windowunmaximise)
// Restores the window to the dimensions and position prior to maximising.
export function WindowUnmaximise(): void;

// [WindowIsMaximised](https://wails.io/docs/reference/runtime/window#windowismaximised)
// Returns the state of the window, i.e. whether the window is maximised or not.
export function WindowIsMaximised(): Promise<boolean>;

// [WindowMinimise](https://wails.io/docs/reference/runtime/window#windowminimise)
// Minimises the window.
export function WindowMinimise(): void;

// [WindowUnminimise](https://wails.io/docs/reference/runtime/window#windowunminimise)
// Restores the window to the dimensions and position prior to minimising.
export function WindowUnminimise(): void;

// [WindowIsMinimised](https://wails.io/docs/reference/runtime/window#windowisminimised)
// Returns the state of the window, i.e. whether the window is minimised or not.
export function WindowIsMinimised(): Promise<boolean>;

// [WindowIsNormal](https://wails.io/docs/reference/runtime/window#windowisnormal)
// Returns the state of the window, i.e. whether the window is normal or not.
export function WindowIsNormal(): Promise<boolean>;

// [WindowSetBackgroundColour](https://wails.io/docs/reference/runtime/window#windowsetbackgroundcolour)
// Sets the background colour of the window to the given RGBA colour definition. This colour will show through for all transparent pixels.
export function WindowSetBackgroundColour(R: number, G: number, B: number, A: number): void;

// [ScreenGetAll](https://wails.io/docs/reference/runtime/window#screengetall)
// Gets the all screens. Call this anew each time you want to refresh data from the underlying windowing system.
export function ScreenGetAll(): Promise<Screen[]>;

// [BrowserOpenURL](https://wails.io/docs/reference/runtime/browser#browseropenurl)
// Opens the given URL in the system browser.
export function BrowserOpenURL(url: string): void;

// [Environment](https://wails.io/docs/reference/runtime/intro#environment)
// Returns information about the environment
export function Environment(): Promise<EnvironmentInfo>;

// [Quit](https://wails.io/docs/reference/runtime/intro#quit)
// Quits the application.
export function Quit(): void;

// [Hide](https://wails.io/docs/reference/runtime/intro#hide)
// Hides the application.
export function Hide(): void;

// [Show](https://wails.io/docs/reference/runtime/intro#show)
// Shows the application.
export function Show(): void;

// [ClipboardGetText](https://wails.io/docs/reference/runtime/clipboard#clipboardgettext)
// Returns the current text stored on clipboard
export function ClipboardGetText(): Promise<string>;

// [ClipboardSetText](https://wails.io/docs/reference/runtime/clipboard#clipboardsettext)
// Sets a text on the clipboard
export function ClipboardSetText(text: string): Promise<boolean>;

// [OnFileDrop](https://wails.io/docs/reference/runtime/draganddrop#onfiledrop)
// OnFileDrop listens to drag and drop events and calls the callback with the coordinates of the drop and an array of path strings.
export function OnFileDrop(callback: (x: number, y: number ,paths: string[]) => void, useDropTarget: boolean) :void

// [OnFileDropOff](https://wails.io/docs/reference/runtime/draganddrop#dragandddropoff)
// OnFileDropOff removes the drag and drop listeners and handlers.
export function OnFileDropOff() :void

// Check if the file path resolver is available
export function CanResolveFilePaths(): boolean;

// Resolves file paths for an array of files
export function ResolveFilePaths(files: File[]): void
```

## `frontend\wailsjs\runtime\runtime.js`

```javascript
/*
 _       __      _ __
| |     / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/

export function LogPrint(message) {
    window.runtime.LogPrint(message);
}

export function LogTrace(message) {
    window.runtime.LogTrace(message);
}

export function LogDebug(message) {
    window.runtime.LogDebug(message);
}

export function LogInfo(message) {
    window.runtime.LogInfo(message);
}

export function LogWarning(message) {
    window.runtime.LogWarning(message);
}

export function LogError(message) {
    window.runtime.LogError(message);
}

export function LogFatal(message) {
    window.runtime.LogFatal(message);
}

export function EventsOnMultiple(eventName, callback, maxCallbacks) {
    return window.runtime.EventsOnMultiple(eventName, callback, maxCallbacks);
}

export function EventsOn(eventName, callback) {
    return EventsOnMultiple(eventName, callback, -1);
}

export function EventsOff(eventName, ...additionalEventNames) {
    return window.runtime.EventsOff(eventName, ...additionalEventNames);
}

export function EventsOnce(eventName, callback) {
    return EventsOnMultiple(eventName, callback, 1);
}

export function EventsEmit(eventName) {
    let args = [eventName].slice.call(arguments);
    return window.runtime.EventsEmit.apply(null, args);
}

export function WindowReload() {
    window.runtime.WindowReload();
}

export function WindowReloadApp() {
    window.runtime.WindowReloadApp();
}

export function WindowSetAlwaysOnTop(b) {
    window.runtime.WindowSetAlwaysOnTop(b);
}

export function WindowSetSystemDefaultTheme() {
    window.runtime.WindowSetSystemDefaultTheme();
}

export function WindowSetLightTheme() {
    window.runtime.WindowSetLightTheme();
}

export function WindowSetDarkTheme() {
    window.runtime.WindowSetDarkTheme();
}

export function WindowCenter() {
    window.runtime.WindowCenter();
}

export function WindowSetTitle(title) {
    window.runtime.WindowSetTitle(title);
}

export function WindowFullscreen() {
    window.runtime.WindowFullscreen();
}

export function WindowUnfullscreen() {
    window.runtime.WindowUnfullscreen();
}

export function WindowIsFullscreen() {
    return window.runtime.WindowIsFullscreen();
}

export function WindowGetSize() {
    return window.runtime.WindowGetSize();
}

export function WindowSetSize(width, height) {
    window.runtime.WindowSetSize(width, height);
}

export function WindowSetMaxSize(width, height) {
    window.runtime.WindowSetMaxSize(width, height);
}

export function WindowSetMinSize(width, height) {
    window.runtime.WindowSetMinSize(width, height);
}

export function WindowSetPosition(x, y) {
    window.runtime.WindowSetPosition(x, y);
}

export function WindowGetPosition() {
    return window.runtime.WindowGetPosition();
}

export function WindowHide() {
    window.runtime.WindowHide();
}

export function WindowShow() {
    window.runtime.WindowShow();
}

export function WindowMaximise() {
    window.runtime.WindowMaximise();
}

export function WindowToggleMaximise() {
    window.runtime.WindowToggleMaximise();
}

export function WindowUnmaximise() {
    window.runtime.WindowUnmaximise();
}

export function WindowIsMaximised() {
    return window.runtime.WindowIsMaximised();
}

export function WindowMinimise() {
    window.runtime.WindowMinimise();
}

export function WindowUnminimise() {
    window.runtime.WindowUnminimise();
}

export function WindowSetBackgroundColour(R, G, B, A) {
    window.runtime.WindowSetBackgroundColour(R, G, B, A);
}

export function ScreenGetAll() {
    return window.runtime.ScreenGetAll();
}

export function WindowIsMinimised() {
    return window.runtime.WindowIsMinimised();
}

export function WindowIsNormal() {
    return window.runtime.WindowIsNormal();
}

export function BrowserOpenURL(url) {
    window.runtime.BrowserOpenURL(url);
}

export function Environment() {
    return window.runtime.Environment();
}

export function Quit() {
    window.runtime.Quit();
}

export function Hide() {
    window.runtime.Hide();
}

export function Show() {
    window.runtime.Show();
}

export function ClipboardGetText() {
    return window.runtime.ClipboardGetText();
}

export function ClipboardSetText(text) {
    return window.runtime.ClipboardSetText(text);
}

/**
 * Callback for OnFileDrop returns a slice of file path strings when a drop is finished.
 *
 * @export
 * @callback OnFileDropCallback
 * @param {number} x - x coordinate of the drop
 * @param {number} y - y coordinate of the drop
 * @param {string[]} paths - A list of file paths.
 */

/**
 * OnFileDrop listens to drag and drop events and calls the callback with the coordinates of the drop and an array of path strings.
 *
 * @export
 * @param {OnFileDropCallback} callback - Callback for OnFileDrop returns a slice of file path strings when a drop is finished.
 * @param {boolean} [useDropTarget=true] - Only call the callback when the drop finished on an element that has the drop target style. (--wails-drop-target)
 */
export function OnFileDrop(callback, useDropTarget) {
    return window.runtime.OnFileDrop(callback, useDropTarget);
}

/**
 * OnFileDropOff removes the drag and drop listeners and handlers.
 */
export function OnFileDropOff() {
    return window.runtime.OnFileDropOff();
}

export function CanResolveFilePaths() {
    return window.runtime.CanResolveFilePaths();
}

export function ResolveFilePaths(files) {
    return window.runtime.ResolveFilePaths(files);
}
```

## `main.go`

```go
package main

import (
	"context"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "Gylte",
		Width:  500, // Taller than wide
		Height: 720,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Frameless:        true, // Key for custom title bar
		BackgroundColour: &options.RGBA{R: 19, G: 19, B: 19, A: 1},
		OnStartup:        app.startup,
		OnDomReady:       app.onDomReady,
		OnBeforeClose:    app.onBeforeClose,
		OnShutdown:       app.onShutdown,
		Bind: []interface{}{
			app, // This automatically binds all public methods of App
		},
		// --- CSS properties for dragging the window ---
		CSSDragProperty: "--wails-draggable",
		CSSDragValue:    "drag",
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

// Optional lifecycle methods for the App
func (a *App) onDomReady(ctx context.Context) {
	// Called when the frontend Dom is ready
}

func (a *App) onBeforeClose(ctx context.Context) (prevent bool) {
	// Called when the application is about to quit,
	// either by clicking the window close button or calling runtime.Quit.
	// Returning true will cause the application to continue running.
	return false
}

func (a *App) onShutdown(ctx context.Context) {
	// Called during shutdown after OnBeforeClose
}
```

