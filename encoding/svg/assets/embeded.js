var g = new Map();
var allLinks = new Array(); 

allLinks.push({{ range .AllLinks }}'{{.}}',{{ end }});
{{range $key, $value := .G }}
g.set('{{$key}}',Array({{range $value}}'{{.}}',{{end}}));
{{- end}}

const max = 9
function replyClick(clicked_id)
{
	console.log(clicked_id);
	var rx = /element_(.*)/;
	var arr = rx.exec(clicked_id);
	var id = arr[1];
	hideID(id)
}
function hideID(id) {
	for (let i = 0; i < max; i++) {
		var myEle = document.getElementById("edge_"+id+"_"+i);
		if(!myEle){
			continue;
		}
		var style = document.getElementById("edge_"+id+"_"+i).style.display;
		if(style === "none")
			document.getElementById("edge_"+id+"_"+i).style.display = "block";
		else
			document.getElementById("edge_"+id+"_"+i).style.display = "none";
		if (id < max) {
			hideID(i)
		}
	}
}

