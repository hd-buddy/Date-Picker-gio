// This file is specially created to provide a custom button style for the application.
// This custom function removes the gray focused background color on clicking the elements
// This only gives gray color to hovered clicable buttons.
package util

import (
	"image"
	"image/color"
	"math"

	"gioui.org/font"
	"gioui.org/io/pointer"
	"gioui.org/io/semantic"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// CustomButtonStyle extends material.ButtonStyle.
type CustomButtonStyle struct {
	Text         string
	Color        color.NRGBA
	Font         font.Font
	TextSize     unit.Sp
	Background   color.NRGBA
	CornerRadius unit.Dp
	Inset        layout.Inset
	Button       *widget.Clickable
	shaper       *text.Shaper
}

// CustomButtonLayoutStyle extends material.ButtonLayoutStyle.
type CustomButtonLayoutStyle struct {
	Background   color.NRGBA
	CornerRadius unit.Dp
	Button       *widget.Clickable
}

// CustomButton initializes a CustomButtonStyle.
func CustomButton(th *material.Theme, button *widget.Clickable, txt string) CustomButtonStyle {
	b := material.Button(th, button, txt)
	return CustomButtonStyle{
		Text:         b.Text,
		Color:        b.Color,
		Font:         b.Font,
		TextSize:     b.TextSize,
		Background:   b.Background,
		CornerRadius: b.CornerRadius,
		Inset:        b.Inset,
		Button:       b.Button,
		shaper:       th.Shaper,
	}
}

// CustomButtonLayout initializes a CustomButtonLayoutStyle.
func CustomButtonLayout(th *material.Theme, button *widget.Clickable) CustomButtonLayoutStyle {
	return CustomButtonLayoutStyle{
		Button:       button,
		Background:   th.Palette.ContrastBg,
		CornerRadius: 4,
	}
}

// Layout method for CustomButtonStyle.
func (b CustomButtonStyle) Layout(gtx layout.Context) layout.Dimensions {
	return CustomButtonLayoutStyle{
		Background:   b.Background,
		CornerRadius: b.CornerRadius,
		Button:       b.Button,
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return b.Inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			colMacro := op.Record(gtx.Ops)
			paint.ColorOp{Color: b.Color}.Add(gtx.Ops)
			return widget.Label{Alignment: text.Middle}.Layout(gtx, b.shaper, b.Font, b.TextSize, b.Text, colMacro.Stop())
		})
	})
}

// Layout method for CustomButtonLayoutStyle.
func (b CustomButtonLayoutStyle) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	min := gtx.Constraints.Min
	return b.Button.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Background{}.Layout(gtx,
			func(gtx layout.Context) layout.Dimensions {
				rr := gtx.Dp(b.CornerRadius)
				defer clip.UniformRRect(image.Rectangle{Max: gtx.Constraints.Min}, rr).Push(gtx.Ops).Pop()
				background := b.Background
				if !gtx.Enabled() {
					background = DisabledColor(b.Background)
				} else if b.Button.Hovered() || gtx.Focused(b.Button) {
					pointer.CursorPointer.Add(gtx.Ops)
				} else {
					pointer.CursorDefault.Add(gtx.Ops)
				}
				paint.Fill(gtx.Ops, background)
				for _, c := range b.Button.History() {
					drawInk(gtx, c)
				}
				return layout.Dimensions{Size: gtx.Constraints.Min}
			},
			func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min = min
				return layout.Center.Layout(gtx, w)
			},
		)
	})
}

func Clickable(gtx layout.Context, button *widget.Clickable, w layout.Widget) layout.Dimensions {

	return button.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		semantic.Button.Add(gtx.Ops)
		return layout.Background{}.Layout(gtx,
			func(gtx layout.Context) layout.Dimensions {
				defer clip.Rect{Max: gtx.Constraints.Min}.Push(gtx.Ops).Pop()
				if button.Hovered() || gtx.Focused(button) {
					paint.Fill(gtx.Ops, Transparent)
				}
				//
				if button.Hovered() {
					pointer.CursorPointer.Add(gtx.Ops)
				} else {
					pointer.CursorDefault.Add(gtx.Ops)
				}
				for _, c := range button.History() {
					drawInk(gtx, c)
				}
				return layout.Dimensions{Size: gtx.Constraints.Min}
			},
			w,
		)
	})
}

func Hovered(c color.NRGBA) (h color.NRGBA) {
	if c.A == 0 {
		// Provide a reasonable default for transparent widgets.
		return color.NRGBA{A: 0x44, R: 0x88, G: 0x88, B: 0x88}
	}
	const ratio = 0x20
	m := color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: c.A}
	if approxLuminance(c) > 128 {
		m = color.NRGBA{A: c.A}
	}
	return mix(m, c, ratio)
}

func approxLuminance(c color.NRGBA) byte {
	const (
		r = 13933 // 0.2126 * 256 * 256
		g = 46871 // 0.7152 * 256 * 256
		b = 4732  // 0.0722 * 256 * 256
		t = r + g + b
	)
	return byte((r*int(c.R) + g*int(c.G) + b*int(c.B)) / t)
}

func mix(c1, c2 color.NRGBA, a uint8) color.NRGBA {
	ai := int(a)
	return color.NRGBA{
		R: byte((int(c1.R)*ai + int(c2.R)*(256-ai)) / 256),
		G: byte((int(c1.G)*ai + int(c2.G)*(256-ai)) / 256),
		B: byte((int(c1.B)*ai + int(c2.B)*(256-ai)) / 256),
		A: byte((int(c1.A)*ai + int(c2.A)*(256-ai)) / 256),
	}
}

// DisabledColor returns the color to be used when the button is disabled.
func DisabledColor(c color.NRGBA) color.NRGBA {
	const disableAlpha = 0.5
	c.R = byte(float32(c.R) * disableAlpha)
	c.G = byte(float32(c.G) * disableAlpha)
	c.B = byte(float32(c.B) * disableAlpha)
	c.A = byte(float32(c.A) * disableAlpha)
	return c
}

func drawInk(gtx layout.Context, c widget.Press) {
	// duration is the number of seconds for the
	// completed animation: expand while fading in, then
	// out.
	const (
		expandDuration = float32(0)
		fadeDuration   = float32(0)
	)

	now := gtx.Now

	t := float32(now.Sub(c.Start).Seconds())

	end := c.End
	if end.IsZero() {
		// If the press hasn't ended, don't fade-out.
		end = now
	}

	endt := float32(end.Sub(c.Start).Seconds())

	// Compute the fade-in/out position in [0;1].
	// var alphat float32
	{
		var haste float32
		if c.Cancelled {
			// If the press was cancelled before the inkwell
			// was fully faded in, fast forward the animation
			// to match the fade-out.
			if h := 0.5 - endt/fadeDuration; h > 0 {
				haste = h
			}
		}
		// Fade in.
		half1 := t/fadeDuration + haste
		if half1 > 0.5 {
			half1 = 0.5
		}

		// Fade out.
		half2 := float32(now.Sub(end).Seconds())
		half2 /= fadeDuration
		half2 += haste
		if half2 > 0.5 {
			// Too old.
			return
		}

		// alphat = half1 + half2
	}

	// Compute the expand position in [0;1].
	sizet := t
	if c.Cancelled {
		// Freeze expansion of cancelled presses.
		sizet = endt
	}
	sizet /= expandDuration

	// Animate only ended presses, and presses that are fading in.
	if !c.End.IsZero() || sizet <= 1.0 {
		gtx.Execute(op.InvalidateCmd{})
	}

	if sizet > 1.0 {
		sizet = 1.0
	}

	// if alphat > .5 {
	// 	// Start fadeout after half the animation.
	// 	// alphat = 1.0 - alphat
	// }
	// Twice the speed to attain fully faded in at 0.5.
	// t2 := alphat * 2
	// Beziér ease-in curve.
	// alphaBezier := t2 * t2 * (3.0 - 2.0*t2)
	sizeBezier := sizet * sizet * (3.0 - 2.0*sizet)
	size := gtx.Constraints.Min.X
	if h := gtx.Constraints.Min.Y; h > size {
		size = h
	}
	// Cover the entire constraints min rectangle and
	// apply curve values to size and color.
	size = int(float32(size) * 2 * float32(math.Sqrt(2)) * sizeBezier)
	// alpha := 0.7 * alphaBezier
	// const col = 0.8
	// ba, bc := byte(alpha*0xff), byte(col*0xff)
	rgba := Transparent
	ink := paint.ColorOp{Color: rgba}
	ink.Add(gtx.Ops)
	rr := size / 2
	defer op.Offset(c.Position.Add(image.Point{
		X: -rr,
		Y: -rr,
	})).Push(gtx.Ops).Pop()
	defer clip.UniformRRect(image.Rectangle{Max: image.Pt(size, size)}, rr).Push(gtx.Ops).Pop()
	paint.PaintOp{}.Add(gtx.Ops)
}
