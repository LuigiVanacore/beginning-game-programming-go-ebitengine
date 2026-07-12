package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/basicfont"
)

// UILabel draws text with optional containment within a parent's bounds.
type UILabel struct {
	Text               string
	X, Y               float64 // X: left edge or center if CenterHorizontally; Y: top of text
	Color              color.Color
	ContainInParent    bool    // when true, text is truncated to fit maxWidth
	CenterHorizontally bool    // when true, X is the horizontal center of the text
	CenterVertically   bool    // when true, Y is the vertical center (text height inferred from content)
	MaxWidth           float64
}

// truncateTextToFit returns text truncated with "..." if it exceeds maxWidth.
func truncateTextToFit(s string, maxWidth float64) string {
	if maxWidth <= 0 {
		return s
	}
	face := text.NewGoXFace(basicfont.Face7x13)
	w, _ := text.Measure(s, face, 0)
	if w <= maxWidth {
		return s
	}
	ellipsis := "..."
	ellipsisW, _ := text.Measure(ellipsis, face, 0)
	available := maxWidth - ellipsisW
	for i := len(s) - 1; i >= 0; i-- {
		sub := s[:i]
		sw, _ := text.Measure(sub, face, 0)
		if sw <= available {
			return sub + ellipsis
		}
	}
	return ellipsis
}

// DrawLabel draws text, optionally truncated to fit within container bounds.
// If label.ContainInParent and maxWidth > 0, text is truncated with "...".
// If label.CenterHorizontally, X is the horizontal center of the text.
// If lineHeight > 0, text is vertically centered in a cell of height lineHeight starting at Y.
func DrawLabel(target *ebiten.Image, label UILabel, maxWidth float64, lineHeight float64) {
	face := text.NewGoXFace(basicfont.Face7x13)
	textToDraw := label.Text
	if label.ContainInParent && maxWidth > 0 {
		textToDraw = truncateTextToFit(label.Text, maxWidth)
	}
	w, h := text.Measure(textToDraw, face, 0)
	x := label.X
	if label.CenterHorizontally {
		x = label.X - w/2
	}
	y := label.Y
	if lineHeight > 0 && h > 0 {
		y = label.Y + (lineHeight-h)/2
	}
	opts := &text.DrawOptions{}
	opts.GeoM.Translate(x, y)
	opts.ColorScale.Reset()
	if label.Color != nil {
		opts.ColorScale.ScaleWithColor(label.Color)
	} else {
		opts.ColorScale.ScaleWithColor(color.White)
	}
	text.Draw(target, textToDraw, face, opts)
}
