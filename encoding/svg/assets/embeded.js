var g = new Map();
var allLinks = new Array(); 

allLinks.push({{ range .AllLinks }}'{{.}}',{{ end }});
{{range $key, $value := .G }}
g.set('{{$key}}',Array({{range $value}}'{{.}}',{{end}}));
{{- end}}

const max = 9
function toggleLink(clicked_id)
{
	if (g.has(clicked_id)) {
		g.get(clicked_id).forEach(element => {
			var style = document.getElementById(element).style.display;
			if(style === "none")
				document.getElementById(element).style.display = "block";
			else
				document.getElementById(element).style.display = "none";

		});
	}
}
function toggleLinks() {
	allLinks.forEach(element => {
		var style = document.getElementById(element).style.display;
		if(style === "none") {
			document.getElementById(element).style.display = "block";
		} else {
			document.getElementById(element).style.display = "none";
		}
	});
}

var allVisibilities = new Map();
var allInVisibilities = new Map();
function setVisibility() {
	{{ range $key, $value := .Visibility }}
	allVisibilities.set('in{{$value.Visibility}}',document.querySelectorAll('.{{$value.Visibility}}'))
	allInVisibilities.set('{{$value.Visibility}}',document.querySelectorAll('.in{{$value.Visibility}}'))
	{{- end}}
}

var visible = true

function toggleVisibility() {
	components = allVisibilities
	if (visible) {
		components = allInVisibilities
		visible = false
	} else {
		visible = true
	}
	components.forEach(function(value, key) {
		if (value.length === 0) {
			setVisibility();
		}
		value.forEach(element => {
			element.classList.toggle(key)	
		});
	})
}


setVisibility();
