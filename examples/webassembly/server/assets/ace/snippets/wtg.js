define("ace/snippets/wtg",["require","exports","module"],function(e,t,n){"use strict";t.snippetText="# Evolution\nsnippet evolution\n	|${1:..}|${2:..}|${3:..}|${4:..}|\n# Configure\nsnippet configure\n	$1: {\n		evolution: |${2:..}|${3:..}|${4:..}|${5:..}|\n	}\n",t.scope="wtg"});
                (function() {
                    window.require(["ace/snippets/wtg"], function(m) {
                        if (typeof module == "object" && typeof exports == "object" && module) {
                            module.exports = m;
                        }
                    });
                })();
            
