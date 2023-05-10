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
