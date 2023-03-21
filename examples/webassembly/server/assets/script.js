var gvar = this;
//(function (document) {
	function onceLoaded(document) {
		//http://stackoverflow.com/a/10372280/398634
		window.URL = window.URL || window.webkitURL;
		var el_stetus = document.getElementById("status"),
			t_stetus = -1,
			reviewer = document.getElementById("review"),
			scale = window.devicePixelRatio || 1,
			downloadBtn = document.getElementById("download"),
			owmBtn = document.getElementById("toowm"),
			sampleBtn = document.getElementById("sample"),
			editor = ace.edit("editor"),
			lastHD = -1,
			worker = null,
			parser = new DOMParser(),
			showError = null,
			shareEl = document.querySelector("#share"),
			applyEl = document.querySelector("#apply"),
			shareURLEl = document.querySelector("#shareurl"),
			errorEl = document.querySelector("#error");

		function show_status(text, hide) {
			hide = hide || 0;
			clearTimeout(t_stetus);
			el_stetus.innerHTML = text;
			if (hide) {
				t_stetus = setTimeout(function () {
					el_stetus.innerHTML = "";
				}, hide);
			}
		}

		function show_error(e) {
			show_status("error", 500);
			reviewer.classList.remove("working");
			reviewer.classList.add("error");

			var message = e.message === undefined ? "An error occurred while processing the graph input." : e.message;
			while (errorEl.firstChild) {
				errorEl.removeChild(errorEl.firstChild);
			}
			errorEl.appendChild(document.createTextNode(message));
		}

		function svgXmlToImage(svgXml, callback) {
			var pngImage = new Image(), svgImage = new Image();

			svgImage.onload = function () {
				var canvas = document.createElement("canvas");
				canvas.width = svgImage.width * scale;
				canvas.height = svgImage.height * scale;

				var context = canvas.getContext("2d");
				context.drawImage(svgImage, 0, 0, canvas.width, canvas.height);

				pngImage.src = canvas.toDataURL("image/png");
				pngImage.width = svgImage.width;
				pngImage.height = svgImage.height;

				if (callback !== undefined) {
					callback(null, pngImage);
				}
			}

			svgImage.onerror = function (e) {
				if (callback !== undefined) {
					callback(e);
				}
			}
			svgImage.src = svgXml;
		}

		function restoreContent() {
			var previousContent = window.localStorage.getItem("wtg");
			editor.getSession().setValue(previousContent);
			renderGraph();
		}
		function copyShareURL(e) {
			var compressedFlow = compress(editor.getSession().getDocument().getValue(),"gzip");

			compressedFlow.then(function(result) {
				// do something with result
				var param = arrayBufferToBase64(result);
				var url = document.URL;
				let params = new URLSearchParams(url.search);
				params.set('wtg', param);
				console.log(params.toString())
				window.history.replaceState({}, '', `${location.pathname}?${params}`);
				shareURLEl.style.display = "inline";
				shareURLEl.value = document.URL;
			});
		}

		function copyToClipboard(str) {
			const el = document.createElement('textarea');
			el.value = str;
			el.setAttribute('readonly', '');
			el.style.position = 'absolute';
			el.style.left = '-9999px';
			document.body.appendChild(el);
			const selected =
				document.getSelection().rangeCount > 0
				? document.getSelection().getRangeAt(0)
				: false;
			el.select();
			var result = document.execCommand('copy')
			document.body.removeChild(el);
			if (selected) {
				document.getSelection().removeAllRanges();
				document.getSelection().addRange(selected);
			}
			return result;
		};

		function renderGraph() {
			var w = parseInt(document.getElementById("width").value);
			var h = parseInt(document.getElementById("height").value);

			svg = generateSVG(editor.getSession().getDocument().getValue(),w,h);
			updateOutput(svg);
			// include script from the SVG
			gvar.eval(document.getElementById('SVGScript').textContent);
		}

		function updateState() {
			var content = encodeURIComponent(editor.getSession().getDocument().getValue());
			history.pushState({"content": content}, "", "#" + content)
		}

		function updateOutput(result) {

			var text = reviewer.querySelector("#text");
			if (text) {
				reviewer.removeChild(text);
			}

			var a = reviewer.querySelector("a");
			if (a) {
				reviewer.removeChild(a);
			}

			if (!result) {
				return;
			}

			reviewer.classList.remove("working");
			reviewer.classList.remove("error");


			//a.appendChild(svgEl);
			var a = document.createElement("a");
			a.innerHTML = result;
			reviewer.appendChild(a);
			// TODO
			var svgEl = document.getElementsByTagName("svg")[0];

			//const url = "data:image/svg+xml;charset=utf-8,"+encodeURIComponent(result);
			const url = 'data:image/svg+xml;base64,'+window.btoa(unescape(encodeURIComponent(a.innerHTML)));
			downloadBtn.href = url;
			downloadBtn.download = "wtg.svg";
			svgPanZoom(svgEl, {
				zoomEnabled: true,
				controlIconsEnabled: true,
				fit: true,
				center: true,
			});

			//updateState()
		}

		editor.setTheme("ace/theme/twilight");
		   editor.setOptions({
        enableBasicAutocompletion: true,
        enableSnippets: true,
        enableLiveAutocompletion: true
    });

		editor.getSession().setMode("ace/mode/wtg");

		editor.getSession().on("change", function () {
			clearTimeout(lastHD);
			lastHD = setTimeout(function(){
				renderGraph();
				// save the content locally for future restore
				window.localStorage.setItem("wtg",editor.getSession().getDocument().getValue());
			}, 1500);
		});

		window.onpopstate = function(event) {
			if (event.state != null && event.state.content != undefined) {
				editor.getSession().setValue(decodeURIComponent(event.state.content));
			}
		};

		share.addEventListener("click", copyShareURL);
		sampleBtn.addEventListener("click", setSample);
		owmBtn.addEventListener("click",function(){
				var w = window.open("");
                w.document.write("<pre>");
                w.document.write(toOWM(editor.getSession().getDocument().getValue()));
                w.document.write("</pre>");
		});
		apply.addEventListener('click', function(){
			console.log("rendering")
			renderGraph()
		})

		// Since apparently HTMLCollection does not implement the oh so convenient array functions
		HTMLOptionsCollection.prototype.indexOf = function(name) {
			for (let i = 0; i < this.length; i++) {
				if (this[i].value == name) {
					return i;
				}
			}

			return -1;
		};

		restoreContent();
		/* come from sharing */
			const params = new URLSearchParams(location.search.substring(1));

		if (params.has('raw')) {
			editor.getSession().setValue(params.get('raw'));
			renderGraph();
		} else if (params.has('wtg')) {
			let wtgText = params.get('wtg'); // "some_value"export 
			console.log(wtgText);
			var content = base64ToArrayBuffer(wtgText);
			if (content!=null) {
				decompress(content,"gzip").then(function(result) {
					console.log(result);
					editor.getSession().setValue(result);
					renderGraph();

				});
			}
		} else if (params.has('compressed')) {
			const compressed = params.get('compressed');
		} else if (params.has('url')) {
			const url = params.get('url');
			let ok = false;
			fetch(url)
				.then(res => {
					ok = res.ok;
					return res.text();
				})
				.then(res => {
					if (!ok) {
						throw { message: res };
					}

					editor.getSession().setValue(res);
					renderGraph();
				}).catch(e => {
					show_error(e);
				});
		} else if (location.hash.length > 1) {
			editor.getSession().setValue(decodeURIComponent(location.hash.substring(1)));
		} else if (editor.getValue()) { // Init
			renderGraph();
		}

		function compress(string, encoding) {
			const byteArray = new TextEncoder().encode(string);
			const cs = new CompressionStream(encoding);
			const writer = cs.writable.getWriter();
			writer.write(byteArray);
			writer.close();
			return new Response(cs.readable).arrayBuffer();
		}

		function decompress(byteArray, encoding) {
			const cs = new DecompressionStream(encoding);
			const writer = cs.writable.getWriter();
			writer.write(byteArray);
			writer.close();
			return new Response(cs.readable).arrayBuffer().then(function (arrayBuffer) {
				return new TextDecoder().decode(arrayBuffer);
			});
		}

		function arrayBufferToBase64( buffer ) {
			var binary = '';
			var bytes = new Uint8Array( buffer );
			var len = bytes.byteLength;
			for (var i = 0; i < len; i++) {
				binary += String.fromCharCode( bytes[ i ] );
			}
			return window.btoa( binary );
		}

		function base64ToArrayBuffer(base64) {
			var binary_string =  window.atob(base64);
			var len = binary_string.length;
			var bytes = new Uint8Array( len );
			for (var i = 0; i < len; i++)        {
				bytes[i] = binary_string.charCodeAt(i);
			}
			return bytes.buffer;
		}
		function setSample() {
			editor.getSession().setValue(`
title: sample map // title is optional
/***************
  value chain 
****************/

business - cup of tea
public - cup of tea
cup of tea - cup
cup of tea -- tea
cup of tea --- hot water
hot water - water
hot water -- kettle
kettle - power

/***************
  definitions 
****************/

// you can inline the evolution
business: |....|....|...x.|.........|

public: |....|....|....|.x...|

// or create blocks
cup of tea: {
  evolution: |....|....|..x..|........|
  color: Green // you can set colors
}
cup: {
  type: buy
  evolution: |....|....|....|....x....|
}
tea: {
  type: buy
  evolution: |....|....|....|.....x....|
}
hot water: {
  evolution: |....|....|....|....x....|
  color: Blue
}
water: {
  type: outsource
  evolution: |....|....|....|.....x....|
}

// you can set the evolution with a >
kettle: {
  type: build
  evolution: |...|...x.|..>.|.......|
}
power: {
  type: pipeline
  evolution: |...|...|....x|.....>..|
}

stage1: genesis / concept
stage2: custom / emerging
stage3: product / converging
stage4: commodity / accepted
				`);
		}

	};
	//})(document);
