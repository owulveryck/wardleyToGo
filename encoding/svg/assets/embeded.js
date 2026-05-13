var allLinks = new Array();

allLinks.push({{ range .AllLinks }}'{{.}}',{{ end }});

const max = 9
function buildTooltip(el) {
	var parts = [];
	el.querySelectorAll('desc').forEach(function(t) {
		var s = t.textContent.trim();
		if (s) {
			var sublines = s.split('\n');
			for (var k = 0; k < sublines.length; k++) {
				var trimmed = sublines[k].trim();
				if (trimmed) parts.push(trimmed);
			}
		}
	});
	if (parts.length === 0) return null;

	var maxWrap = 40;
	var allLines = [];
	for (var p = 0; p < parts.length; p++) {
		var part = parts[p];
		if (part.length <= maxWrap) {
			allLines.push({text: part, isLabel: false});
		} else {
			var wrapped = part.match(new RegExp('.{1,' + maxWrap + '}(\\s|$)', 'g')) || [part];
			for (var w = 0; w < wrapped.length; w++) {
				allLines.push({text: wrapped[w].trim(), isLabel: false});
			}
		}
		if (p < parts.length - 1) {
			allLines.push({text: '', isLabel: false});
		}
	}

	for (var i = 0; i < allLines.length; i++) {
		if (/^(Asset|Cost):\s/.test(allLines[i].text)) {
			allLines[i].isLabel = true;
		}
	}

	var ns = 'http://www.w3.org/2000/svg';
	var tip = document.createElementNS(ns, 'g');
	tip.setAttribute('class', 'tooltip-box');

	var fontSize = 13;
	var lineHeight = 18;
	var padX = 14;
	var padY = 10;
	var maxWidth = 0;
	for (var i = 0; i < allLines.length; i++) {
		if (allLines[i].text.length > maxWidth) maxWidth = allLines[i].text.length;
	}
	var boxW = maxWidth * 7.5 + padX * 2;
	var boxH = allLines.length * lineHeight + padY * 2;

	var rect = document.createElementNS(ns, 'rect');
	rect.setAttribute('x', '12');
	rect.setAttribute('y', String(-boxH / 2));
	rect.setAttribute('width', String(boxW));
	rect.setAttribute('height', String(boxH));
	rect.setAttribute('rx', '6');
	rect.setAttribute('ry', '6');
	rect.setAttribute('fill', 'rgba(25,25,30,0.94)');
	rect.setAttribute('stroke', 'rgba(255,255,255,0.15)');
	rect.setAttribute('stroke-width', '1');
	tip.appendChild(rect);

	for (var i = 0; i < allLines.length; i++) {
		var line = allLines[i];
		if (line.text === '') continue;

		var t = document.createElementNS(ns, 'text');
		t.setAttribute('x', String(12 + padX));
		t.setAttribute('y', String(-boxH / 2 + padY + fontSize + i * lineHeight));
		t.setAttribute('font-size', fontSize + 'px');
		t.setAttribute('font-family', "'Outfit', sans-serif");
		t.setAttribute('fill', 'white');

		if (line.isLabel) {
			var colonIdx = line.text.indexOf(': ');
			var labelPart = line.text.substring(0, colonIdx + 1);
			var valuePart = line.text.substring(colonIdx + 1);

			var boldSpan = document.createElementNS(ns, 'tspan');
			boldSpan.setAttribute('font-weight', 'bold');
			boldSpan.setAttribute('fill', 'rgba(255,255,255,0.95)');
			boldSpan.textContent = labelPart;
			t.appendChild(boldSpan);

			var valueSpan = document.createElementNS(ns, 'tspan');
			valueSpan.setAttribute('fill', 'rgba(255,255,255,0.75)');
			valueSpan.textContent = valuePart;
			t.appendChild(valueSpan);
		} else {
			t.textContent = line.text;
		}
		tip.appendChild(t);
	}

	return tip;
}

function getComponentTransform(el) {
	var g = el.querySelector('g[transform]');
	if (!g) return '';
	return g.getAttribute('transform');
}

function findTooltipFor(id) {
	var layer = document.getElementById('tooltip-layer');
	if (!layer) return null;
	return layer.querySelector('[data-for="' + id + '"]');
}

function showTooltip(id) {
	if (findTooltipFor(id)) return;

	var el = document.getElementById(id);
	if (!el) return;

	var tip = buildTooltip(el);
	if (!tip) return;

	tip.setAttribute('data-for', id);
	tip.setAttribute('transform', getComponentTransform(el));

	var layer = document.getElementById('tooltip-layer');
	if (layer) layer.appendChild(tip);
}

function hideTooltip(id) {
	var tip = findTooltipFor(id);
	if (tip && tip.getAttribute('data-pinned') !== 'true') {
		tip.remove();
	}
}

function pinTooltip(id) {
	var existing = findTooltipFor(id);
	if (existing) {
		if (existing.getAttribute('data-pinned') === 'true') {
			existing.remove();
		} else {
			existing.setAttribute('data-pinned', 'true');
		}
		return;
	}

	var el = document.getElementById(id);
	if (!el) return;

	var tip = buildTooltip(el);
	if (!tip) return;
	tip.setAttribute('data-for', id);
	tip.setAttribute('data-pinned', 'true');
	tip.setAttribute('transform', getComponentTransform(el));

	var layer = document.getElementById('tooltip-layer');
	if (layer) layer.appendChild(tip);
}

function toggleLink(id) { pinTooltip(id); }
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
