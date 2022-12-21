function onTextChange() {
	var key = window.event.keyCode;

	// If the user has pressed enter
	if ((key == 10 || key == 13) && window.event.ctrlKey){
		displayImage();
		return false;
	}
	else {
		return true;
	}
}

function displayImage() {
	svg = generateSVG(document.getElementById("code").value);
	document.getElementById("svgContainer").innerHTML = svg;
	svgImage = document.getElementsByTagName("svg")[0];
	svgImage.style = ""
	svgImage.removeAttribute("width");
	svgImage.removeAttribute("height");
	svgImage.setAttribute("preserveAspectRatio", "xMidYMid meet")
	console.log(right.width)
	console.log(document.getElementById("colRight").getBoundingClientRect().width)
	console.log(svgSize)
	svgImage.getAttribute('viewBox')
	var box = svgImage.getAttribute('viewBox').split(/\s+|,/);
	if (box[2] < right.width) {
		svgImage.setAttribute('viewBox', `0 0 ${right.width} ${fullHeight}`);
	}
	svgSize = { w: svgImage.clientWidth, h: svgImage.clientHeight };

}

const svgContainer = document.getElementById("svgContainer");
const right = document.getElementById("svgContainer").getBoundingClientRect();
const left = document.getElementById("colLeft").getBoundingClientRect();
svgImage = document.getElementsByTagName("svg")[0];
svgImage.style = ""
svgImage.removeAttribute("width");
svgImage.removeAttribute("height");
console.log(window.innerHeight)


// Full height, including the scroll part
const fullHeight = Math.max(
	document.body.scrollHeight,
	document.documentElement.scrollHeight,
	document.body.offsetHeight,
	document.documentElement.offsetHeight,
	document.body.clientHeight,
	document.documentElement.clientHeight
);
const fullWidth = Math.max(
	document.body.scrollWidth,
	document.documentElement.scrollWidth,
	document.body.offsetWidth,
	document.documentElement.offsetWidth,
	document.body.clientWidth,
	document.documentElement.clientWidth
);
//svgImage.setAttribute('viewBox', `0 0 ${right.width} ${fullWidth}`);
svgImage.setAttribute('viewBox', `0 0 1200 ${fullHeight}`);
console.log(fullHeight)

var viewBox = { x: 0, y: 0, w: right.width, h: fullHeight };
//svgImage.setAttribute('viewBox', ` + "`" + `${viewBox.x} ${viewBox.y} ${viewBox.w} ${viewBox.h}` + "`" + `);
svgSize = { w: svgImage.clientWidth, h: svgImage.clientHeight };
var isPanning = false;
var startPoint = { x: 0, y: 0 };
var endPoint = { x: 0, y: 0 };;
var scale = 1;

svgContainer.onmousewheel = function (e) {
	e.preventDefault();
	var w = viewBox.w;
	var h = viewBox.h;
	var mx = e.offsetX;//mouse x  
	var my = e.offsetY;
	var dw = w * Math.sign(e.deltaY) * 0.05;
	var dh = h * Math.sign(e.deltaY) * 0.05;
	var dx = dw * mx / svgSize.w;
	var dy = dh * my / svgSize.h;
	viewBox = { x: viewBox.x + dx, y: viewBox.y + dy, w: viewBox.w - dw, h: viewBox.h - dh };
	scale = svgSize.w / viewBox.w;
	zoomValue.innerText = ` ${Math.round(scale * 100) / 100}`;
	svgImage.setAttribute('viewBox', `${viewBox.x} ${viewBox.y} ${viewBox.w} ${viewBox.h}`);
}


svgContainer.onmousedown = function (e) {
	isPanning = true;
	startPoint = { x: e.x, y: e.y };
}

svgContainer.onmousemove = function (e) {
	if (isPanning) {
		endPoint = { x: e.x, y: e.y };
		var dx = (startPoint.x - endPoint.x) / scale;
		var dy = (startPoint.y - endPoint.y) / scale;
		var movedViewBox = { x: viewBox.x + dx, y: viewBox.y + dy, w: viewBox.w, h: viewBox.h };
		svgImage.setAttribute('viewBox', `${movedViewBox.x} ${movedViewBox.y} ${movedViewBox.w} ${movedViewBox.h}`);
	}
}

svgContainer.onmouseup = function (e) {
	if (isPanning) {
		endPoint = { x: e.x, y: e.y };
		var dx = (startPoint.x - endPoint.x) / scale;
		var dy = (startPoint.y - endPoint.y) / scale;
		viewBox = { x: viewBox.x + dx, y: viewBox.y + dy, w: viewBox.w, h: viewBox.h };
		svgImage.setAttribute('viewBox', `${viewBox.x} ${viewBox.y} ${viewBox.w} ${viewBox.h}`);
		isPanning = false;
	}
}

svgContainer.onmouseleave = function (e) {
	isPanning = false;
}

