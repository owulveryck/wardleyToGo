// wardleyToGo: sketch with SVGo, (derived from the old misc/goplay), except:
// (1) only listen on localhost, (default port 1999)
// (2) always render html,
// (3) SVGo default code,
package main

import (
	"flag"
	"log"
	"net/http"
	"text/template"

	svgmap "github.com/owulveryck/wardleyToGo/internal/encoding/svg"
	"github.com/owulveryck/wardleyToGo/internal/parser"
)

var (
	httpListen = flag.String("port", "1999", "port to listen on")
)

var (
	// a source of numbers, for naming temporary files
	uniq = make(chan int)
)

func main() {
	flag.Parse()

	// source of unique numbers
	go func() {
		for i := 0; ; i++ {
			uniq <- i
		}
	}()
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/render", render)
	log.Printf("☠ ☠ ☠ Warning: this server allows a client connecting to 127.0.0.1:%s to execute code on this computer ☠ ☠ ☠", *httpListen)
	log.Fatal(http.ListenAndServe("127.0.0.1:"+*httpListen, nil))
}

// FrontPage is an HTTP handler that renders the svgoplay interface.
// If a filename is supplied in the path component of the URI,
// its contents will be put in the interface's text area.
// Otherwise, the default "hello, world" program is displayed.
func FrontPage(w http.ResponseWriter, req *http.Request) {
	frontPage.Execute(w, sampleMap)
}

// render is an HTTP handler that reads Go source code from the request,
// runs the program (returning any errors),
// and sends the program's output as the HTTP response.
func render(w http.ResponseWriter, r *http.Request) {
	// target is the base name for .go, object, executable files
	// write body to target.go
	width := 1100
	height := 800
	padLeft := 25
	padBottom := 30

	p := parser.NewParser(r.Body)
	defer r.Body.Close()
	m, err := p.Parse() // the map
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	svgmap.Encode(m, w, width, height, padLeft, padBottom)
}

var frontPage = template.Must(template.New("frontPage").Parse(frontPageText)) // HTML template

var frontPageText = `<!doctype html>
<html>
<head>
<title>wardleyToGo: SVGo sketching</title>
<style>
.svg-container { 
	display: inline-block;
	position: relative;
	width: 100%;
	padding-bottom: 100%; 
	vertical-align: middle; 
	overflow: hidden; 
}
.wardley-map { 
	display: inline-block;
	position: absolute;
	top: 0;
	left: 0;
}
body {
	font-size: 12pt;
}
pre, textarea {
	font-family: Optima, Calibri, 'DejaVu Sans', sans-serif;
	font-size: 100%;
	line-height: 15pt;
}
.hints {
	font-size: 80%;
	font-family: sans-serif;
	color: #333333;
	font-style: italic;
	text-align: right;
}
#edit, #output, #errors { width: 100%; text-align: left; }
#edit { height: 560px; }
#errors { color: #c00; }
</style>
<script>

function insertTabs(n) {
	// find the selection start and end
	var cont  = document.getElementById("edit");
	var start = cont.selectionStart;
	var end   = cont.selectionEnd;
	// split the textarea content into two, and insert n tabs
	var v = cont.value;
	var u = v.substr(0, start);
	for (var i=0; i<n; i++) {
		u += "\t";
	}
	u += v.substr(end);
	// set revised content
	cont.value = u;
	// reset caret position after inserted tabs
	cont.selectionStart = start+n;
	cont.selectionEnd = start+n;
}

function autoindent(el) {
	var curpos = el.selectionStart;
	var tabs = 0;
	while (curpos > 0) {
		curpos--;
		if (el.value[curpos] == "\t") {
			tabs++;
		} else if (tabs > 0 || el.value[curpos] == "\n") {
			break;
		}
	}
	setTimeout(function() {
		insertTabs(tabs);
	}, 1);
}

function keyHandler(event) {
	var e = window.event || event;
	if (e.keyCode == 9) { // tab
		insertTabs(1);
		e.preventDefault();
		return false;
	}
	if (e.keyCode == 13) { // enter
		if (e.shiftKey) { // +shift
			render(e.target);
			e.preventDefault();
			return false;
		} else {
			autoindent(e.target);
		}
	}
	return true;
}

var xmlreq;

function autorender() {
	if(!document.getElementById("autorender").checked) {
		return;
	}
	render();
}

function render() {
	var prog = document.getElementById("edit").value;
	var req = new XMLHttpRequest();
	xmlreq = req;
	req.onreadystatechange = renderUpdate;
	req.open("POST", "/render", true);
	req.setRequestHeader("Content-Type", "text/plain; charset=utf-8");
	req.send(prog);	
}

function renderUpdate() {
	var req = xmlreq;
	if(!req || req.readyState != 4) {
		return;
	}
	if(req.status == 200) {
		document.getElementById("output").innerHTML = req.responseText;
		document.getElementById("errors").innerHTML = "";
	} else {
		document.getElementById("errors").innerHTML = req.responseText;
		document.getElementById("output").innerHTML = "";
	}
}
</script>
</head>
<body>
<textarea autofocus="true" id="edit" spellcheck="false" onkeydown="keyHandler(event);" onkeyup="autorender();">{{printf "%s" . |html}}</textarea>
<div class="hints">
(Shift-Enter to render and run.)&nbsp;&nbsp;&nbsp;&nbsp;
<input type="checkbox" id="autorender" value="checked" /> Compile and run after each keystroke
</div>
<div id="output"></div>
<div id="errors"></div>
</body>
</html>
`

var sampleMap = []byte(`
title Tea Shop
anchor Business [0.95, 0.63]
anchor Public [0.95, 0.78]
component Cup of Tea [0.79, 0.61] label [19, -4]
component Cup [0.73, 0.78] label [19,-4] (dataProduct)
component Tea [0.63, 0.81]
component Hot Water [0.52, 0.80]
component Water [0.38, 0.82]
component Kettle [0.43, 0.35] label [-73, 4] (build)
evolve Kettle 0.62 label [22, 9] (buy)
component Power [0.1, 0.7] label [-29, 30] (outsource)
evolve Power 0.89 label [-12, 21]
Business->Cup of Tea
Public->Cup of Tea
Cup of Tea-collaboration>Cup
Cup of Tea-collaboration>Tea
Cup of Tea-collaboration>Hot Water
Hot Water->Water
Hot Water-facilitating>Kettle 
Kettle-xAsAService>Power
build Kettle


annotation 1 [[0.43,0.49],[0.08,0.79]] Standardising power allows Kettles to evolve faster
annotation 2 [0.48, 0.85] Hot water is obvious and well known
annotations [0.60, 0.02]

note +a generic note appeared [0.16, 0.36]

style wardley
streamAlignedTeam team A [0.84, 0.58, 0.74, 0.68]
enablingTeam team B [0.52, 0.23, 0.32, 0.43]
platformTeam team C [0.18, 0.61, 0.02, 0.94]
complicatedSubsystemTeam team D [0.83, 0.73, 0.45, 0.90]
`)
