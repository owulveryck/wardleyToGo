package components

import (
	"image"

	"github.com/owulveryck/wardleyToGo/v2/internal/utils"
)

// CalcCoords calculates the coordinates wrt to the bounds.
// it scales accordingly
func CalcCoords(p image.Point, bounds image.Rectangle) image.Point {
	return utils.CalcCoords(p, bounds)
}
