
function insertTabs(n) {
    // find the selection start and end
    var cont = document.getElementById("edit");
    var start = cont.selectionStart;
    var end = cont.selectionEnd;
    // split the textarea content into two, and insert n tabs
    var v = cont.value;
    var u = v.substr(0, start);
    for (var i = 0; i < n; i++) {
        u += "\t";
    }
    u += v.substr(end);
    // set revised content
    cont.value = u;
    // reset caret position after inserted tabs
    cont.selectionStart = start + n;
    cont.selectionEnd = start + n;
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
    setTimeout(function () {
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
    if (!document.getElementById("autorender").checked) {
        return;
    }
    render();
}

function render() {
    var prog = document.getElementById("edit").value;
    var svg = generateSVG(prog);
    document.getElementById("output").innerHTML = svg;
    /*
    var req = new XMLHttpRequest();
    xmlreq = req;
    req.onreadystatechange = renderUpdate;
    req.open("POST", "/render", true);
    req.setRequestHeader("Content-Type", "text/plain; charset=utf-8");
    req.send(prog);
    */
    save();

}

function save() {
    var content = document.getElementById("edit").value;
    //localStorage["user"] = user ;
    localStorage.setItem("content", content);

}

function load() {
    var content = localStorage.getItem("content");
    document.getElementById("edit").value = content;
}

function start() {
    render();
}
window.onload = start();

function showTooltip(evt, text) {
    let tooltip = document.getElementById("tooltip");
    tooltip.innerHTML = text;
    tooltip.style.display = "block";
    tooltip.style.left = evt.pageX + 10 + 'px';
    tooltip.style.top = evt.pageY - 10 + 'px';
}

function hideTooltip() {
    var tooltip = document.getElementById("tooltip");
    tooltip.style.display = "none";
}

function insertAtCursor(myField, myValue) {
    //IE support
    if (document.selection) {
        myField.focus();
        var sel = document.selection.createRange();
        sel.text = myValue;
    }
    //MOZILLA and others
    else if (myField.selectionStart || myField.selectionStart == '0') {
        var startPos = myField.selectionStart;
        var endPos = myField.selectionEnd;
        myField.value = myField.value.substring(0, startPos)
            + myValue
            + myField.value.substring(endPos, myField.value.length);
    } else {
        myField.value += myValue;
    }
}

function renderUpdate() {
    var req = xmlreq;
    if (!req || req.readyState != 4) {
        return;
    }
    if (req.status == 200) {
        document.getElementById("output").innerHTML = req.responseText;
        document.getElementById("errors").innerHTML = "";
    } else {
        document.getElementById("errors").innerHTML = req.responseText;
        document.getElementById("output").innerHTML = "";
    }
}