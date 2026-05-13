var allLinks = new Array();

allLinks.push({{ range .AllLinks }}'{{.}}',{{ end }});

const max = 9
var activeTooltip = null;
function toggleLink(clicked_id)
{
	var el = document.getElementById(clicked_id);
	if (!el) return;

	var existing = el.querySelector('.tooltip-box');
	if (existing) {
		existing.remove();
		if (activeTooltip === clicked_id) activeTooltip = null;
		return;
	}

	if (activeTooltip) {
		var prev = document.querySelector('.tooltip-box');
		if (prev) prev.remove();
		activeTooltip = null;
	}

	var parts = [];
	el.querySelectorAll('title').forEach(function(t) {
		var s = t.textContent.trim();
		if (s) parts.push(s);
	});
	var text = parts.join(' — ');
	if (!text) return;

	var target = el.querySelector('g[transform]') || el;

	var ns = 'http://www.w3.org/2000/svg';
	var tip = document.createElementNS(ns, 'g');
	tip.setAttribute('class', 'tooltip-box');

	var lines = text.match(/.{1,30}(\s|$)/g) || [text];
	var lineHeight = 16;
	var maxWidth = 0;
	lines.forEach(function(l) { if (l.length > maxWidth) maxWidth = l.length; });
	var boxW = maxWidth * 7 + 16;
	var boxH = lines.length * lineHeight + 12;

	var rect = document.createElementNS(ns, 'rect');
	rect.setAttribute('x', '12');
	rect.setAttribute('y', String(-boxH / 2));
	rect.setAttribute('width', String(boxW));
	rect.setAttribute('height', String(boxH));
	rect.setAttribute('rx', '4');
	rect.setAttribute('ry', '4');
	rect.setAttribute('fill', 'rgba(30,30,30,0.92)');
	rect.setAttribute('stroke', 'none');
	tip.appendChild(rect);

	for (var i = 0; i < lines.length; i++) {
		var t = document.createElementNS(ns, 'text');
		t.setAttribute('x', '20');
		t.setAttribute('y', String(-boxH / 2 + 16 + i * lineHeight));
		t.setAttribute('font-size', '12px');
		t.setAttribute('font-family', "'Outfit', sans-serif");
		t.setAttribute('fill', 'white');
		t.textContent = lines[i].trim();
		tip.appendChild(t);
	}

	target.appendChild(tip);
	activeTooltip = clicked_id;
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
