package wardley

import (
	"image"
	"image/color"

	"github.com/owulveryck/wardleyToGo"
)

type Element interface {
	wardleyToGo.Component
	wardleyToGo.Chainer
	wardleyToGo.Positioner
	Colorer
	Labeler
}

type Colorer interface {
	SetColor(color.Color)
}

type Labeler interface {
	SetLabelPlacement(image.Point)
	SetLabelAnchor(int)
	GetLabel() string
	SetLabel(string)
}

type Configurer interface {
	SetConfiguredFlag()
	GetConfiguredFlag() bool
}
