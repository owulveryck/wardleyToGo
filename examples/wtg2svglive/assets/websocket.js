var gvar = this;
let socket = new WebSocket("ws://localhost:8080/ws");

socket.onopen = function(e) {
	//alert("[open] Connection established");
};

socket.onmessage = function(event) {
	//alert(`[message] Data received from server: ${event.data}`);
	console.log("received message");
	//console.log(`received ${event.data}`);
	document.getElementById("svgContainer").innerHTML = event.data;
	gvar.eval(document.getElementById('SVGScript').textContent);

};

socket.onclose = function(event) {
	if (event.wasClean) {
		alert(`[close] Connection closed cleanly, code=${event.code} reason=${event.reason}`);
	} else {
		// e.g. server process killed or network down
		// event.code is usually 1006 in this case
		alert('[close] Connection died');
	}
};

socket.onerror = function(error) {
	alert(`[error]`);
};
