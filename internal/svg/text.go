package svg

import (
	"encoding/xml"
	"image"
	"strconv"
	"sync"
)

const (
	TextAnchorUndefined int = iota
	TextAnchorStart
	TextAnchorMiddle
	TextAnchorEnd
)

var (
	defaultFontMu sync.RWMutex
	defaultFont   string
)

func UpdateDefaultFont(fontName string) {
	defaultFontMu.Lock()
	defer defaultFontMu.Unlock()
	defaultFont = fontName
}

func getDefaultFont() string {
	defaultFontMu.RLock()
	defer defaultFontMu.RUnlock()
	return defaultFont
}

type TextArea struct {
	P          image.Point
	Text       []byte
	Fill       Color
	TextAnchor int
	FontWeight string
	FontSize   string
	FontFamily string
	TextAdjust bool
}

func (t TextArea) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	foreignObject := xml.StartElement{
		Name: xml.Name{Local: "foreignObject"},
	}
	foreignObject.Attr = []xml.Attr{
		{
			Name:  xml.Name{Local: "transform"},
			Value: "translate(12,-10)",
		},
		{
			Name:  xml.Name{Local: "width"},
			Value: "90",
		},
		{
			Name:  xml.Name{Local: "height"},
			Value: "60",
		},
	}
	textArea := xml.StartElement{
		Name: xml.Name{Local: "textArea"},
	}
	style := ""
	if t.FontFamily != "" {
		style = style + "font-family:" + t.FontFamily + ";"
	} else if f := getDefaultFont(); f != "" {
		style = style + "font-family:" + f + ";"
	}
	atr, _ := t.Fill.MarshalXMLAttr(xml.Name{Local: "fill"})
	style = style + "color:" + atr.Value + ";"

	textArea.Attr = []xml.Attr{
		{
			Name:  xml.Name{Local: "xmlns"},
			Value: "http://www.w3.org/1999/xhtml",
		},
		{
			Name:  xml.Name{Local: "style"},
			Value: style,
		},
		{
			Name:  xml.Name{Local: "class"},
			Value: "componentText",
		},
	}

	if err := e.EncodeToken(foreignObject); err != nil {
		return err
	}
	if err := e.EncodeToken(textArea); err != nil {
		return err
	}
	if err := e.EncodeToken(xml.CharData(t.Text)); err != nil {
		return err
	}
	if err := e.EncodeToken(textArea.End()); err != nil {
		return err
	}
	return e.EncodeToken(foreignObject.End())
}

type Text struct {
	P          image.Point
	Text       []byte
	Fill       Color
	TextAnchor int
	FontWeight string
	FontSize   string
	FontFamily string
	TextAdjust bool
	MaxChars   int
}

func (t Text) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	maxChars := t.MaxChars
	if maxChars == 0 {
		maxChars = 8
	}
	element := xml.StartElement{
		Name: xml.Name{Local: "text"},
	}
	element.Attr = []xml.Attr{
		/*
			{
				Name:  xml.Name{Local: "x"},
				Value: strconv.Itoa(t.P.X),
			},
			{
				Name:  xml.Name{Local: "y"},
				Value: strconv.Itoa(t.P.Y),
			},
		*/
		must(t.Fill.MarshalXMLAttr(xml.Name{Local: "fill"})),
		must(t.Fill.MarshalXMLAttr(xml.Name{Local: "fill-opacity"})),
	}
	if t.FontWeight != "" {
		element.Attr = append(element.Attr, xml.Attr{
			Name:  xml.Name{Local: "font-weight"},
			Value: t.FontWeight,
		})
	}
	if t.FontFamily != "" {
		element.Attr = append(element.Attr, xml.Attr{
			Name:  xml.Name{Local: "font-family"},
			Value: t.FontFamily,
		})
	} else if f := getDefaultFont(); f != "" {
		element.Attr = append(element.Attr, xml.Attr{
			Name:  xml.Name{Local: "font-family"},
			Value: f,
		})
	}
	if t.FontSize != "" {
		element.Attr = append(element.Attr, xml.Attr{
			Name:  xml.Name{Local: "font-size"},
			Value: t.FontSize,
		})
	}
	switch t.TextAnchor {
	case TextAnchorEnd:
		element.Attr = append(element.Attr, xml.Attr{
			Name:  xml.Name{Local: "text-anchor"},
			Value: "end",
		})
	case TextAnchorStart:
		element.Attr = append(element.Attr, xml.Attr{
			Name:  xml.Name{Local: "text-anchor"},
			Value: "start",
		})
	case TextAnchorMiddle:
		element.Attr = append(element.Attr, xml.Attr{
			Name:  xml.Name{Local: "text-anchor"},
			Value: "middle",
		})
	}
	if err := e.EncodeToken(element); err != nil {
		return err
	}
	words := []string{string(t.Text)}
	if t.TextAdjust {
		words = splitString(string(t.Text), maxChars)
	}
	for i, word := range words {
		dy := t.P.Y
		switch {
		case dy < 0:
			if i == 0 {
				dy = dy * len(words)
			} else {
				dy = 18
			}
		case dy == 0:
			dy = 18
			if i == 0 {
				dy = ((len(words))/2)*-6 + 6
			}
		case dy > 0:
			if i > 0 && dy < 10 {
				dy = 18
			}
		}
		tspan := xml.StartElement{
			Name: xml.Name{Local: "tspan"},
		}
		tspan.Attr = []xml.Attr{
			{
				Name: xml.Name{Local: "x"},
				//Value: "10",
				Value: strconv.Itoa(t.P.X),
			},
			{
				Name:  xml.Name{Local: "dy"},
				Value: strconv.Itoa(dy),
				//Value: "20",
			},
		}
		if err := e.EncodeToken(tspan); err != nil {
			return err
		}
		if err := e.EncodeToken(xml.CharData(word)); err != nil {
			return err
		}
		if err := e.EncodeToken(tspan.End()); err != nil {
			return err
		}
	}
	return e.EncodeToken(element.End())
}
