export WTG2SVG_WITHCONTROLS=false
export WTG2SVG_WIDTH=1100
export WTG2SVG_HEIGHT=900
export WTG2SVG_WITHSPACE=true
export WTG2SVG_WITHVALUECHAIN=true

WTG2SVG := wtg2svg 
INKSCAPE := inkscape

SRCS := $(wildcard *.wtg)
SVGS := $(SRCS:.wtg=.svg)
PNGS := $(SVGS:.svg=.png)

.PHONY: all clean

all: $(PNGS) $(SVGS)


%.svg: %.wtg
	cat $< | $(WTG2SVG) > $@

%.png: %.svg %.wtg
	cat $< | inkscape -p --export-type=png > $@

clean:
	rm $(SVGS) $(PNGS)
