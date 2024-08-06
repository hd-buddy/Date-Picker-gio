package util

import (
	"image"
	"image/color"
	"unicode"

	"gioui.org/font"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/inkeliz/giosvg"
)

type (
	C  = layout.Context
	D  = layout.Dimensions
	mt = material.Theme
)

func LoadSvg(Icon []byte) *giosvg.Icon {
	Vector, _ := giosvg.NewVector(Icon)
	svgIcon := giosvg.NewIcon(Vector)
	return svgIcon
}
func LayoutText(gtx C, th *mt, text Text) D {
	return layout.Inset{
		Top:    unit.Dp(text.Inset.Top),
		Right:  unit.Dp(text.Inset.Right),
		Bottom: unit.Dp(text.Inset.Bottom),
		Left:   unit.Dp(text.Inset.Left),
	}.Layout(gtx, func(gtx C) D {
		element := material.Label(th, unit.Sp(text.Size), text.Text)
		element.Color = text.TextColor
		element.Font = font.Font{Typeface: "Nunito", Weight: text.FontWeight}
		element.Alignment = text.Alignment

		return element.Layout(gtx)
	})
}

func LayoutInputBox(gtx C, th *material.Theme, inputBox InputBox) (D, *widget.SubmitEvent) {
	var submitEvent *widget.SubmitEvent

	dims := layout.Inset{
		Top:    unit.Dp(inputBox.OutInset.Top),
		Right:  unit.Dp(inputBox.OutInset.Right),
		Bottom: unit.Dp(inputBox.OutInset.Bottom),
		Left:   unit.Dp(inputBox.OutInset.Left),
	}.Layout(gtx, func(gtx C) D {
		return widget.Border{
			Color:        inputBox.BorderColor,
			Width:        unit.Dp(1),
			CornerRadius: unit.Dp(inputBox.CornerRadius),
		}.Layout(gtx, func(gtx C) D {
			element := material.Editor(th, inputBox.Editor, inputBox.Hint)
			element.Font = font.Font{Typeface: "Nunito", Weight: inputBox.FontWeight}
			element.Color = inputBox.TextColor
			element.TextSize = unit.Sp(inputBox.Size)
			element.Editor.SingleLine = inputBox.SingleLine
			element.Editor.Submit = inputBox.Submit

			for {
				ev, ok := element.Editor.Update(gtx)
				if !ok {
					break
				}
				switch ev := ev.(type) {
				case widget.SubmitEvent:
					submitEvent = &ev
				}
			}
			if inputBox.HideText {
				inputBox.Editor.Mask = '●'
			} else {
				inputBox.Editor.Mask = 0
			}
			inputWidth := int(float32(inputBox.Width))
			inputHeight := int(float32(inputBox.Height))
			gtx.Constraints = layout.Exact(image.Pt(inputWidth, inputHeight))
			return layout.Inset{
				Top:    unit.Dp(inputBox.InInset.Top),
				Right:  unit.Dp(inputBox.InInset.Right),
				Bottom: unit.Dp(inputBox.InInset.Bottom),
				Left:   unit.Dp(inputBox.InInset.Left),
			}.Layout(gtx, element.Layout)
		})
	})

	return dims, submitEvent
}
func LayoutNumberInputBox(gtx C, th *material.Theme, inputBox InputBox) (D, *widget.SubmitEvent) {
	var submitEvent *widget.SubmitEvent

	dims := layout.Inset{
		Top:    unit.Dp(inputBox.OutInset.Top),
		Right:  unit.Dp(inputBox.OutInset.Right),
		Bottom: unit.Dp(inputBox.OutInset.Bottom),
		Left:   unit.Dp(inputBox.OutInset.Left),
	}.Layout(gtx, func(gtx C) D {
		return widget.Border{
			Color:        inputBox.BorderColor,
			Width:        unit.Dp(1),
			CornerRadius: unit.Dp(inputBox.CornerRadius),
		}.Layout(gtx, func(gtx C) D {
			element := material.Editor(th, inputBox.Editor, inputBox.Hint)
			element.Font = font.Font{Typeface: "Nunito", Weight: inputBox.FontWeight}
			element.Color = inputBox.TextColor
			element.TextSize = unit.Sp(inputBox.Size)
			element.Editor.SingleLine = inputBox.SingleLine
			element.Editor.Submit = inputBox.Submit
			if inputBox.IsRightAlign {
				element.Editor.Alignment = text.End
			} else {
				element.Editor.Alignment = text.Start
			}
			// Track if we need to update the text and cursor position
			var updateText bool
			var filteredText string
			var dotFound bool

			for _, r := range inputBox.Editor.Text() {
				if unicode.IsDigit(r) {
					filteredText += string(r)
				} else if r == '.' && !dotFound {
					filteredText += string(r)
					dotFound = true
				} else {
					updateText = true
				}
			}

			if updateText {
				// Set the filtered text and move the cursor to the end
				inputBox.Editor.SetText(filteredText)
				inputBox.Editor.MoveCaret(len(filteredText), len(filteredText))
			}

			for {
				ev, ok := element.Editor.Update(gtx)
				if !ok {
					break
				}
				switch ev := ev.(type) {
				case widget.SubmitEvent:
					submitEvent = &ev
				}
			}

			if inputBox.HideText {
				inputBox.Editor.Mask = '●'
			} else {
				inputBox.Editor.Mask = 0
			}

			inputWidth := int(float32(inputBox.Width))
			inputHeight := int(float32(inputBox.Height))
			gtx.Constraints = layout.Exact(image.Pt(inputWidth, inputHeight))
			return layout.Inset{
				Top:    unit.Dp(inputBox.InInset.Top),
				Right:  unit.Dp(inputBox.InInset.Right),
				Bottom: unit.Dp(inputBox.InInset.Bottom),
				Left:   unit.Dp(inputBox.InInset.Left),
			}.Layout(gtx, element.Layout)
		})
	})

	return dims, submitEvent
}

func LayoutInputBoxWithIcon(gtx C, th *mt, inputBox InputBox, icon Icon, iconPostion layout.Direction) (D, *widget.SubmitEvent) {
	var submitEvent *widget.SubmitEvent

	dims := layout.Stack{Alignment: iconPostion}.Layout(gtx,
		layout.Stacked(func(gtx C) D {
			var inputBox1 D
			inputBox1, submitEvent = LayoutInputBox(gtx, th, InputBox{
				Editor:       inputBox.Editor,
				Hint:         inputBox.Hint,
				TextColor:    inputBox.TextColor,
				BorderColor:  inputBox.BorderColor,
				FontWeight:   inputBox.FontWeight,
				CornerRadius: inputBox.CornerRadius,
				Width:        inputBox.Width,
				Height:       inputBox.Height,
				Size:         inputBox.Size,
				SingleLine:   inputBox.SingleLine,
				Submit:       inputBox.Submit,
				HideText:     inputBox.HideText,

				OutInset: layout.Inset{
					Top:    unit.Dp(inputBox.OutInset.Top),
					Right:  unit.Dp(inputBox.OutInset.Right),
					Bottom: unit.Dp(inputBox.OutInset.Bottom),
					Left:   unit.Dp(inputBox.OutInset.Left),
				},
				InInset: layout.Inset{
					Top:    unit.Dp(inputBox.InInset.Top),
					Right:  unit.Dp(inputBox.InInset.Right),
					Bottom: unit.Dp(inputBox.InInset.Bottom),
					Left:   unit.Dp(inputBox.InInset.Left),
				},
			})
			return inputBox1
		}),
		layout.Stacked(func(gtx C) D {
			return layout.Inset{Top: icon.Inset.Top,
				Right:  icon.Inset.Right,
				Bottom: icon.Inset.Bottom,
				Left:   icon.Inset.Left}.Layout(gtx, func(gtx C) D {
				Width := int(float32(icon.Width))
				Height := int(float32(icon.Height))
				gtx.Constraints = layout.Exact(image.Pt(Width, Height))
				if icon.ToggleIcon {
					return Clickable(gtx, icon.IconButton, func(gtx C) D {
						if !icon.ToggleCondition {
							return icon.Icon1.Layout(gtx)
						} else {
							return icon.Icon2.Layout(gtx)
						}
					})
				} else {
					return Clickable(gtx, icon.IconButton, func(gtx C) D {
						return icon.Icon1.Layout(gtx)
					})
				}

			})
		}),
	)
	return dims, submitEvent
}

func LayoutButton(gtx C, th *mt, button Button) D {
	return layout.Inset{
		Top:    unit.Dp(button.OutInset.Top),
		Right:  unit.Dp(button.OutInset.Right),
		Bottom: unit.Dp(button.OutInset.Bottom),
		Left:   unit.Dp(button.OutInset.Left),
	}.Layout(gtx, func(gtx C) D {
		element := CustomButton(th, button.Button, button.Text)
		element.TextSize = unit.Sp(button.Size)
		element.Color = button.TextColor
		element.Font = font.Font{Typeface: "Nunito", Weight: button.FontWeight}
		if button.Button.Hovered() {
			pointer.CursorPointer.Add(gtx.Ops)
		} else {
			pointer.CursorDefault.Add(gtx.Ops)
		}

		element.Background = button.BackgroundColor
		element.CornerRadius = button.CornerRadius
		border := widget.Border{Color: button.BorderColor, CornerRadius: button.CornerRadius, Width: unit.Dp(1)}
		element.Inset = layout.Inset{
			Top:    unit.Dp(button.InInset.Top),
			Right:  unit.Dp(button.InInset.Right),
			Bottom: unit.Dp(button.InInset.Bottom),
			Left:   unit.Dp(button.InInset.Left),
		}
		return border.Layout(gtx, func(gtx C) D {
			return element.Layout(gtx)
		})
	})
}

func LayoutIconButton(gtx C, th *mt, button Button, icon Icon, alignment layout.Direction) D {
	return layout.Stack{Alignment: alignment}.Layout(gtx,
		layout.Stacked(func(gtx C) D {
			return LayoutButton(gtx, th, Button{
				BackgroundColor: button.BackgroundColor,
				Button:          button.Button,
				BorderColor:     button.BorderColor,
				CornerRadius:    button.CornerRadius,
				FontWeight:      button.FontWeight,
				InInset: layout.Inset{
					Top:    unit.Dp(icon.Height)/2 + icon.Inset.Top,
					Bottom: unit.Dp(icon.Height)/2 + icon.Inset.Bottom,
					Left:   unit.Dp(icon.Width)/2 + icon.Inset.Left,
					Right:  unit.Dp(icon.Width)/2 + icon.Inset.Right,
				},
				OutInset:  button.OutInset,
				Size:      button.Size,
				Text:      button.Text,
				TextColor: button.TextColor,
			})
		}),
		layout.Stacked(func(gtx C) D {
			return layout.Inset{
				Top:    icon.Inset.Top,
				Right:  icon.Inset.Right,
				Bottom: icon.Inset.Bottom,
				Left:   icon.Inset.Left,
			}.Layout(gtx, func(gtx C) D {
				IconWidth := int(float32(icon.Width))
				IconHeight := int(float32(icon.Height))
				gtx.Constraints = layout.Exact(image.Pt(IconWidth, IconHeight))
				if icon.ToggleIcon {
					if !icon.ToggleCondition {
						return icon.Icon1.Layout(gtx)
					} else {
						return icon.Icon2.Layout(gtx)
					}
				} else {
					return icon.Icon1.Layout(gtx)
				}

			})
		}),
	)
}
func DrawALine(gtx C, linecolor color.NRGBA) D {
	line := widget.Border{Color: linecolor, CornerRadius: unit.Dp(0), Width: unit.Dp(1)}
	return line.Layout(gtx, func(gtx C) D {
		return D{Size: image.Point{X: gtx.Constraints.Max.X, Y: 1}}
	})
}
func DrawALineThicknes(gtx C, linecolor color.NRGBA, thickness int) D {
	line := widget.Border{Color: linecolor, CornerRadius: unit.Dp(0), Width: unit.Dp(1)}
	return line.Layout(gtx, func(gtx C) D {
		return D{Size: image.Point{X: gtx.Constraints.Max.X, Y: thickness}}
	})
}

func DrawALineVertical(gtx C, linecolor color.NRGBA) D {
	line := widget.Border{Color: linecolor, CornerRadius: unit.Dp(0), Width: unit.Dp(0.5)}
	return line.Layout(gtx, func(gtx C) D {
		return D{Size: image.Point{X: 1, Y: gtx.Constraints.Max.Y}}
	})

}

func ModalHeader(gtx C, th *material.Theme, TitleText string, buttonClose *widget.Clickable) D {
	return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx C) D {
		return layout.Flex{
			Axis:      layout.Horizontal,
			Alignment: layout.Middle,
			Spacing:   layout.SpaceBetween,
		}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return LayoutText(gtx, th, Text{
					Text:       TitleText,
					Size:       20,
					FontWeight: font.SemiBold,
					TextColor:  BlackColor,
				})
			}),
			layout.Flexed(1, func(gtx C) D {
				return layout.Spacer{}.Layout(gtx) // Take up remaining space
			}),
			layout.Rigid(func(gtx C) D {
				var closeColor color.NRGBA
				if buttonClose.Hovered() {
					closeColor = RedColor
				} else {
					closeColor = GrayColor
				}
				return LayoutButton(gtx, th, Button{
					Text:            "✕",
					TextColor:       closeColor,
					Size:            20,
					FontWeight:      font.Bold,
					BackgroundColor: Transparent,
					Button:          buttonClose,
					InInset:         layout.Inset{Right: 10},
				})
			}),
		)
	})
}

type Button struct {
	BackgroundColor color.NRGBA
	Button          *widget.Clickable
	BorderColor     color.NRGBA
	CornerRadius    unit.Dp
	FontWeight      font.Weight
	InInset         layout.Inset
	OutInset        layout.Inset
	Size            int
	Text            string
	TextColor       color.NRGBA
}

type Icon struct {
	Height          int
	Icon1           *giosvg.Icon
	Icon2           *giosvg.Icon
	IconButton      *widget.Clickable
	Inset           layout.Inset
	ToggleCondition bool
	ToggleIcon      bool
	Width           int
}

type InputBox struct {
	BorderColor  color.NRGBA
	CornerRadius int
	Editor       *widget.Editor
	EditorMask   string
	FontWeight   font.Weight
	Height       int
	HideText     bool
	Hint         string
	InInset      layout.Inset
	OutInset     layout.Inset
	Size         int
	SingleLine   bool
	Submit       bool
	TextColor    color.NRGBA
	Width        unit.Dp
	Alignment    text.Alignment
	IsRightAlign bool
}

type Text struct {
	FontWeight font.Weight
	Inset      layout.Inset
	Size       int
	Text       string
	TextColor  color.NRGBA
	Alignment  text.Alignment
}

type Border struct {
	BorderColor  color.NRGBA
	Width        int
	CornerRadius int
}

var GreenColor = color.NRGBA{R: 0, G: 255, B: 0, A: 180}
var GrayColor = color.NRGBA{R: 128, G: 128, B: 128, A: 255}
var WhiteColor = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
var BlackColor = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
var RedColor = color.NRGBA{R: 255, G: 0, B: 0, A: 255}
var DarkGrayColor = color.NRGBA{R: 34, G: 40, B: 50, A: 255}
var Transparent = color.NRGBA{R: 0, G: 0, B: 0, A: 0}
var DarkRedColor = color.NRGBA{R: 171, G: 36, B: 33, A: 255}
var LGrayColor = color.NRGBA{R: 238, G: 238, B: 238, A: 255}
var BlackTextColor = color.NRGBA{R: 28, G: 28, B: 28, A: 255}
var GrayBackground = color.NRGBA{R: 230, G: 230, B: 235, A: 255}
var OverlayGray = color.NRGBA{R: 238, G: 238, B: 238, A: 50}
var GrayTextColor = color.NRGBA{R: 84, G: 100, B: 125, A: 255}
var GrayShadow_high = color.NRGBA{R: 0, G: 0, B: 0, A: 50}
var GrayShadow_med = color.NRGBA{R: 0, G: 0, B: 0, A: 30}
var GrayShadow_low = color.NRGBA{R: 0, G: 0, B: 0, A: 10}
var MGrayColor = color.NRGBA{R: 220, G: 220, B: 220, A: 255}
var GrayCardBackground = color.NRGBA{R: 240, G: 240, B: 240, A: 255}
