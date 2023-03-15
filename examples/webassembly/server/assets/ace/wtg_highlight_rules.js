// WTG_highlight_rules.js

define(function(require, exports, module) {
	"use strict";

	var oop = require("../lib/oop");
	var TextHighlightRules = require("ace/mode/text_highlight_rules").TextHighlightRules;

	var WTGHighlightRules = function () {
		var keywordMapper = this.createKeywordMapper({
			"keyword": "evolution|type|color|title|stage1|stage2|stage3|stage4"
		}, "identifier", true);

		this.$rules = {
			"start": [
				{
					token: keywordMapper,
					regex: "\\b(" + keywordMapper + ")\\b"
				}
			]
		};
	};

	oop.inherits(WTGHighlightRules, TextHighlightRules);

	exports.WTGHighlightRules = WTGHighlightRules;

	oop.inherits(WTGHighlightRules, TextHighlightRules);
	exports.WTGHighlightRules = WTGHighlightRules;
});
/*
	ace.define("ace/mode/WTG_highlight_rules", ["require", "exports", "module", "ace/lib/oop", "ace/mode/text_highlight_rules"], function (require, exports, module) {
		var oop = require("../lib/oop");
		var TextHighlightRules = require("./text_highlight_rules").TextHighlightRules;

		var WTGHighlightRules = function () {
			var keywordMapper = this.createKeywordMapper({
				"keyword": "evolution|type|color|title|stage1|stage2|stage3|stage4"
			}, "identifier", true);

			this.$rules = {
				"start": [
					{
						token: keywordMapper,
						regex: "\\b(" + keywordMapper + ")\\b"
					}
				]
			};
		};

		oop.inherits(WTGHighlightRules, TextHighlightRules);

		exports.WTGHighlightRules = WTGHighlightRules;
	});
*/
