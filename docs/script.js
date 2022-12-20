const right = document.getElementById("content").getBoundingClientRect();

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

function onTextChange() {
	var key = window.event.keyCode;

	// If the user has pressed enter
	if ((key == 10 || key == 13) && window.event.ctrlKey){
		svg = generateSVG(document.getElementById("code").value);
		document.getElementById("content").innerHTML = svg;
		svgImage = document.getElementsByTagName("svg")[0];

		var box = svgImage.getAttribute('viewBox').split(/\s+|,/);
		if (box[2] < right.width) {
			svgImage.setAttribute('viewBox', `0 0 ${right.width-200} ${fullHeight}`);
		}
		return false;
	}
	else {
		return true;
	}
}

