package ui

import (
	. "book/code/ch11/internal/core"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	upgradePanelWidth  = 160.0 // panel width in pixels
	upgradePanelHeight = 180.0 // panel height in pixels
	upgradePanelGap    = 20.0  // gap between panels
)

// UpgradePanel draws a panel with weapon name, icon, upgrade value, and select button.
// ButtonHovered is set by the caller before Draw for hover effect.
// icon40: icon pre-rendered once at a fixed 40x40 size; drawn without scaling each frame (avoids shimmer)
type UpgradePanel struct {
	X, Y          float64
	Width         float64
	Height        float64
	WeaponName    string
	UpgradeValue  string
	Icon          *ebiten.Image // nil = no icon (source)
	icon40        *ebiten.Image // pre-rendered 40x40 icon; used for stable drawing
	OnSelect      func()
	ButtonHovered bool // true when mouse is over the select button
}

const (
	panelPadding     = 12.0
	panelBorderW     = 2.0
	buttonHeight     = 28.0
	iconSize         = 40.0
	nameRowHeight    = 18.0 // height for the name
	iconRowHeight    = 48.0 // icon + gap
	upgradeRowHeight = 18.0 // height for the upgrade value
	rowGap           = 4.0  // gap between rows
)

// NewUpgradePanel creates a panel for the given upgrade choice.
// The icon is pre-rendered once at 40x40: each frame draws the same image without scaling (no shimmer).
func NewUpgradePanel(x, y, width, height float64, opt UpgradeChoice, rm *ResourceManager, onSelect func()) *UpgradePanel {
	icon, _ := rm.GetTexture(opt.IconKey)
	if icon == nil && opt.IconKey == "holy_shield" {
		icon, _ = rm.GetTexture("sacred_book") // fallback
	}
	p := &UpgradePanel{
		X:            x,
		Y:            y,
		Width:        width,
		Height:       height,
		WeaponName:   opt.WeaponName,
		UpgradeValue: opt.UpgradeDesc,
		Icon:         icon,
		OnSelect:     onSelect,
	}
	if icon != nil {
		p.icon40 = newScaledIcon(icon)
	}
	return p
}

// newScaledIcon creates a 40x40 image with the icon scaled to fit. Used only once.
func newScaledIcon(src *ebiten.Image) *ebiten.Image {
	const size = 40
	dst := ebiten.NewImage(size, size)
	b := src.Bounds()
	iw, ih := float64(b.Dx()), float64(b.Dy())
	scale := size / iw
	if ih*scale > size {
		scale = size / ih
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate((size-scale*iw)/2, (size-scale*ih)/2)
	dst.DrawImage(src, op)
	return dst
}

// ContainsPoint returns true if (sx, sy) is inside the panel (for click detection).
func (p *UpgradePanel) ContainsPoint(sx, sy float64) bool {
	return sx >= p.X && sx <= p.X+p.Width && sy >= p.Y && sy <= p.Y+p.Height
}

// ButtonBounds returns the button rect (x, y, w, h) for hit testing in screen coords.
func (p *UpgradePanel) ButtonBounds() (x, y, w, h float64) {
	btnX := p.X + panelPadding
	btnY := p.Y + p.Height - panelPadding - buttonHeight
	btnW := p.Width - panelPadding*2
	btnH := buttonHeight
	return btnX, btnY, btnW, btnH
}

// ButtonContainsPoint returns true if (sx, sy) is inside the select button.
func (p *UpgradePanel) ButtonContainsPoint(sx, sy float64) bool {
	bx, by, bw, bh := p.ButtonBounds()
	return sx >= bx && sx <= bx+bw && sy >= by && sy <= by+bh
}

// Draw renders the panel on the target. Uses UIContainer to clip content to panel bounds.
func (p *UpgradePanel) Draw(target *ebiten.Image) {
	cont := NewUIContainer(p.X, p.Y, p.Width, p.Height)
	cont.BorderWidth = panelBorderW
	cont.Draw(target, func(canvas *ebiten.Image) {
		maxW := p.Width - panelPadding*2
		y := panelPadding

		// Row 1: weapon name
		DrawLabel(canvas, UILabel{
			Text:               p.WeaponName,
			X:                  panelPadding + maxW/2,
			Y:                  y,
			Color:              color.RGBA{255, 255, 255, 255},
			ContainInParent:    true,
			CenterHorizontally: true,
		}, maxW, nameRowHeight)
		y += nameRowHeight + rowGap

		// Row 2: 40x40 icon (Round = same position each frame)
		if p.icon40 != nil {
			iconX := math.Round((p.Width - iconSize) / 2)
			iconY := math.Round(y)
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(iconX, iconY)
			canvas.DrawImage(p.icon40, op)
		}
		y += iconRowHeight + rowGap

		// Row 3: upgrade value
		DrawLabel(canvas, UILabel{
			Text:               p.UpgradeValue,
			X:                  panelPadding + maxW/2,
			Y:                  y,
			Color:              color.RGBA{180, 220, 255, 255},
			ContainInParent:    true,
			CenterHorizontally: true,
		}, maxW, upgradeRowHeight)

		// Button (Round = same position each frame)
		btnX := math.Round(panelPadding)
		btnY := math.Round(p.Height - panelPadding - buttonHeight)
		btnW := math.Round(p.Width - panelPadding*2)
		btnH := math.Round(buttonHeight)

		btnFill := color.RGBA{60, 100, 160, 255}
		btnBorder := color.RGBA{100, 150, 220, 255}
		if p.ButtonHovered {
			btnFill = color.RGBA{85, 135, 210, 255}
			btnBorder = color.RGBA{140, 190, 255, 255}
		}
		vector.DrawFilledRect(canvas, float32(btnX), float32(btnY), float32(btnW), float32(btnH),
			btnFill, true)
		vector.StrokeRect(canvas, float32(btnX), float32(btnY), float32(btnW), float32(btnH),
			1, btnBorder, true)

		DrawLabel(canvas, UILabel{
			Text:               "SELECT",
			X:                  btnX + btnW/2,
			Y:                  btnY,
			Color:              color.RGBA{255, 255, 255, 255},
			CenterHorizontally: true,
		}, btnW, btnH)
	})
}
