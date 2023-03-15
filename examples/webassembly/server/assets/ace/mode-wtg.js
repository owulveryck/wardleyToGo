define("ace/mode/matching_brace_outdent", ["require", "exports", "module", "ace/range"], function(e, t, n) {
	"use strict";
	var r = e("../range").Range,
		i = function() {};
	(function() {
		this.checkOutdent = function(e, t) {
			return /^\s+$/.test(e) ? /^\s*\}/.test(t) : !1
	}, this.autoOutdent = function(e, t) {
		var n = e.getLine(t),
			i = n.match(/^(\s*\})/);
		if (!i) return 0;
		var s = i[1].length,
			o = e.findMatchingBracket({
				row: t,
				column: s
			});
	if (!o || o.row == t) return 0;
	var u = this.$getIndent(e.getLine(o.row));
	e.replace(new r(t, 0, t, s - 1), u)
}, this.$getIndent = function(e) {
	return e.match(/^\s*/)[0]
}
}).call(i.prototype), t.MatchingBraceOutdent = i
}), define("ace/mode/doc_comment_highlight_rules", ["require", "exports", "module", "ace/lib/oop", "ace/mode/text_highlight_rules"], function(e, t, n) {
	"use strict";
	var r = e("../lib/oop"),
		i = e("./text_highlight_rules").TextHighlightRules,
		s = function() {
			this.$rules = {
				start: [{
					token: "comment.doc.tag",
					regex: "@[\\w\\d_]+"
				}, s.getTagRule(), {
					defaultToken: "comment.doc",
					caseInsensitive: !0
				}]
			}
		};
	r.inherits(s, i), s.getTagRule = function(e) {
		return {
			token: "comment.doc.tag.storage.type",
			regex: "\\b(?:TODO|FIXME|XXX|HACK)\\b"
		}
	}, s.getStartRule = function(e) {
		return {
			token: "comment.doc",
			regex: "\\/\\*(?=\\*)",
			next: e
		}
	}, s.getEndRule = function(e) {
		return {
			token: "comment.doc",
			regex: "\\*\\/",
			next: e
		}
	}, t.DocCommentHighlightRules = s
}), define("ace/mode/wtg_highlight_rules", ["require", "exports", "module", "ace/lib/oop", "ace/lib/lang", "ace/mode/text_highlight_rules", "ace/mode/doc_comment_highlight_rules"], function(e, t, n) {
	"use strict";
	var r = e("../lib/oop"),
		i = e("../lib/lang"),
		s = e("./text_highlight_rules").TextHighlightRules,
		o = e("./doc_comment_highlight_rules").DocCommentHighlightRules,
		u = function() {
			var e = i.arrayToMap("evolution|title|stage1|stage2|stage3|stage4|color|type".split("|")),
				t = i.arrayToMap("pipeline|build|buy|outsource".split("|"));
			this.$rules = {
				start: [{
					token: "comment",
					regex: /\/\/.*$/
				}, {
					token: "comment",
					merge: !0,
					regex: /\/\*/,
					next: "comment"
				}, {
					token: "keyword.operator",
					regex: /:/
				}, {
					token: "paren.lparen",
					regex: /[\[{]/
					}, {
						token: "paren.rparen",
						regex: /[\]}]/
				}, {
					token: "constant",
					regex: /\|\.*x?>?\.*\|\.*x?>?\.*\|\.*x?>?\.*\|\.*x?>?\.*\|/
				}, {
					token: "invalid",
					regex: /\|[\.]*\|[\.]*\|[\.]*\|[\.]*\|/
				}, {
					token: "comment",
					regex: /\./
				}, {
					token: "constant",
					regex: /\-/
				}, {
					token: function(n) {
						return e.hasOwnProperty(n.toLowerCase()) ? "keyword" : t.hasOwnProperty(n.toLowerCase()) ? "support" : "text"
					},
					regex: "\\-?[a-zA-Z_][a-zA-Z0-9_\\-]*"
				}],
				comment: [{
					token: "comment",
					regex: "\\*\\/",
					next: "start"
				}, {
					defaultToken: "comment"
				}],
				qqstring: [{
					token: "string",
					regex: '[^"\\\\]+',
					merge: !0
				}, {
					token: "string",
					regex: "\\\\$",
					next: "qqstring",
					merge: !0
				}, {
					token: "string",
					regex: '"|$',
					next: "start",
					merge: !0
				}],
				qstring: [{
					token: "string",
					regex: "[^'\\\\]+",
					merge: !0
				}, {
					token: "string",
					regex: "\\\\$",
					next: "qstring",
					merge: !0
				}, {
					token: "string",
					regex: "'|$",
					next: "start",
					merge: !0
				}]
			}
		};
	r.inherits(u, s), t.DotHighlightRules = u
}), define("ace/mode/folding/cstyle", ["require", "exports", "module", "ace/lib/oop", "ace/range", "ace/mode/folding/fold_mode"], function(e, t, n) {
	"use strict";
	var r = e("../../lib/oop"),
		i = e("../../range").Range,
		s = e("./fold_mode").FoldMode,
		o = t.FoldMode = function(e) {
			e && (this.foldingStartMarker = new RegExp(this.foldingStartMarker.source.replace(/\|[^|]*?$/, "|" + e.start)), this.foldingStopMarker = new RegExp(this.foldingStopMarker.source.replace(/\|[^|]*?$/, "|" + e.end)))
		};
	r.inherits(o, s),
		function() {
			this.foldingStartMarker = /([\{\[\(])[^\}\]\)]*$|^\s*(\/\*)/, this.foldingStopMarker = /^[^\[\{\(]*([\}\]\)])|^[\s\*]*(\*\/)/, this.singleLineBlockCommentRe = /^\s*(\/\*).*\*\/\s*$/, this.tripleStarBlockCommentRe = /^\s*(\/\*\*\*).*\*\/\s*$/, this.startRegionRe = /^\s*(\/\*|\/\/)#?region\b/, this._getFoldWidgetBase = this.getFoldWidget, this.getFoldWidget = function(e, t, n) {
				var r = e.getLine(n);
				if (this.singleLineBlockCommentRe.test(r) && !this.startRegionRe.test(r) && !this.tripleStarBlockCommentRe.test(r)) return "";
				var i = this._getFoldWidgetBase(e, t, n);
				return !i && this.startRegionRe.test(r) ? "start" : i
			}, this.getFoldWidgetRange = function(e, t, n, r) {
				var i = e.getLine(n);
				if (this.startRegionRe.test(i)) return this.getCommentRegionBlock(e, i, n);
				var s = i.match(this.foldingStartMarker);
				if (s) {
					var o = s.index;
					if (s[1]) return this.openingBracketBlock(e, s[1], n, o);
					var u = e.getCommentFoldRange(n, o + s[0].length, 1);
					return u && !u.isMultiLine() && (r ? u = this.getSectionRange(e, n) : t != "all" && (u = null)), u
				}
				if (t === "markbegin") return;
				var s = i.match(this.foldingStopMarker);
				if (s) {
					var o = s.index + s[0].length;
					return s[1] ? this.closingBracketBlock(e, s[1], n, o) : e.getCommentFoldRange(n, o, -1)
				}
			}, this.getSectionRange = function(e, t) {
				var n = e.getLine(t),
					r = n.search(/\S/),
					s = t,
					o = n.length;
				t += 1;
				var u = t,
					a = e.getLength();
				while (++t < a) {
					n = e.getLine(t);
					var f = n.search(/\S/);
					if (f === -1) continue;
					if (r > f) break;
					var l = this.getFoldWidgetRange(e, "all", t);
					if (l) {
						if (l.start.row <= s) break;
						if (l.isMultiLine()) t = l.end.row;
						else if (r == f) break
					}
					u = t
				}
				return new i(s, o, u, e.getLine(u).length)
			}, this.getCommentRegionBlock = function(e, t, n) {
				var r = t.search(/\s*$/),
					s = e.getLength(),
					o = n,
					u = /^\s*(?:\/\*|\/\/|--)#?(end)?region\b/,
					a = 1;
				while (++n < s) {
					t = e.getLine(n);
					var f = u.exec(t);
					if (!f) continue;
					f[1] ? a-- : a++;
					if (!a) break
				}
				var l = n;
				if (l > o) return new i(o, r, l, t.length)
			}
		}.call(o.prototype)
}), define("ace/mode/wtg", ["require", "exports", "module", "ace/lib/oop", "ace/mode/text", "ace/mode/matching_brace_outdent", "ace/mode/wtg_highlight_rules", "ace/mode/folding/cstyle"], function(e, t, n) {
	"use strict";
	var r = e("../lib/oop"),
		i = e("./text").Mode,
		s = e("./matching_brace_outdent").MatchingBraceOutdent,
		o = e("./wtg_highlight_rules").DotHighlightRules,
		u = e("./folding/cstyle").FoldMode,
		a = function() {
			this.HighlightRules = o, this.$outdent = new s, this.foldingRules = new u, this.$behaviour = this.$defaultBehaviour
		};
	r.inherits(a, i),
		function() {
			this.lineCommentStart = ["//", "#"], this.blockComment = {
				start: "/*",
				end: "*/"
			}, this.getNextLineIndent = function(e, t, n) {
				var r = this.$getIndent(t),
					i = this.getTokenizer().getLineTokens(t, e),
					s = i.tokens,
					o = i.state;
				if (s.length && s[s.length - 1].type == "comment") return r;
				if (e == "start") {
					var u = t.match(/^.*(?:\bcase\b.*:|[\{\(\[])\s*$/);
						u && (r += n)
					}
						return r
					}, this.checkOutdent = function(e, t, n) {
						return this.$outdent.checkOutdent(t, n)
					}, this.autoOutdent = function(e, t, n) {
						this.$outdent.autoOutdent(t, n)
					}, this.$id = "ace/mode/wtg"
			}.call(a.prototype), t.Mode = a
		});
					(function() {
						window.require(["ace/mode/wtg"], function(m) {
							if (typeof module == "object" && typeof exports == "object" && module) {
								module.exports = m;
							}
						});
					})();
