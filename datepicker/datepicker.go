package datepicker

import (
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"strconv"
	"time"

	"gioui.org/font"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/hd-buddy/GioCalendarPicker/util"
)

type DatePicker struct {
	Date      time.Time
	IsOpen    bool
	Openbtn   widget.Clickable
	PrevBtn   widget.Clickable
	NextBtn   widget.Clickable
	Days      [31]widget.Clickable
	MonthBtn  widget.Clickable
	YearBtn   widget.Clickable
	Months    [12]widget.Clickable
	Years     [20]widget.Clickable
	ViewMode  string // "date", "month", or "year"
	YearRange int    // Starting year for the year picker
	TodayBtn  widget.Clickable
	Editor    *widget.Editor
}
type (
	C  = layout.Context
	D  = layout.Dimensions
	mt = material.Theme
)

// / Date Picker

//go:embed assets/date.svg
var DateIcon []byte
var FlagDpDate bool

func (dp *DatePicker) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if dp.Openbtn.Clicked(gtx) {
		dp.IsOpen = !dp.IsOpen
		dp.ViewMode = "date"
	}
	if dp.TodayBtn.Clicked(gtx) {
		dp.Date = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, dp.Date.Location())
		// dp.IsOpen = false
	}
	DateIcon := util.LoadSvg(DateIcon)
	dp.Editor.SetText(dp.Date.Format("02-Jan-2006"))
	return layout.Stack{Alignment: layout.N}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			// gtx.Constraints.Min.X = gtx.Constraints.Max.X
			inputBox, _ := util.LayoutInputBoxWithIcon(gtx, th, util.InputBox{
				Editor:       dp.Editor,
				BorderColor:  BlackColor,
				CornerRadius: 6,
				FontWeight:   font.SemiBold,
				Hint:         dp.Date.Format("02-Jan-2006"),
				TextColor:    BlackColor,
				Size:         16,
				Width:        unit.Dp(gtx.Constraints.Max.X),
				Height:       42,
				InInset:      layout.Inset{Left: 10, Top: 10, Bottom: 10, Right: 30},
			}, util.Icon{
				Height:     20,
				Width:      20,
				Icon1:      DateIcon,
				IconButton: &dp.Openbtn,
				Inset:      layout.Inset{Right: 10},
			}, layout.E)

			return inputBox
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			if dp.IsOpen {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Spacer{Height: 50}.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return widget.Border{
							Color:        LGrayColor,
							Width:        unit.Dp(1),
							CornerRadius: 6,
						}.Layout(gtx, func(gtx C) D {
							return widget.Border{
								Color:        MGrayColor,
								Width:        unit.Dp(1),
								CornerRadius: 6,
							}.Layout(gtx, func(gtx C) D {
								return widget.Border{
									Color:        GrayColor,
									Width:        unit.Dp(1),
									CornerRadius: 6,
								}.Layout(gtx, func(gtx C) D {
									return layout.UniformInset(unit.Dp(1)).Layout(gtx, func(gtx C) D {
										return dp.calendarLayout(gtx, th)
									})
								})
							})
						})
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return util.LayoutText(gtx, th, util.Text{
									Text:       "Created by Dhruv Hingu",
									Size:       12,
									TextColor:  BlackColor,
									FontWeight: font.Bold,
									Inset:      layout.UniformInset(unit.Dp(10)),
								})
							}),
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								return layout.Spacer{}.Layout(gtx)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return util.LayoutButton(gtx, th, util.Button{
									Text:            "Go to Today",
									TextColor:       color.NRGBA{R: 0, G: 0, B: 255, A: 255},
									Size:            12,
									FontWeight:      font.Bold,
									BackgroundColor: Transparent,
									BorderColor:     GrayColor,
									CornerRadius:    6,
									Button:          &dp.TodayBtn,
									InInset:         layout.UniformInset(unit.Dp(10)),
									OutInset:        layout.UniformInset(unit.Dp(10)),
								})
							}),
						)
					}),
				)
			}
			return layout.Dimensions{}
		}),
	)
}

var btnText string

func (dp *DatePicker) calendarLayout(gtx layout.Context, th *material.Theme) layout.Dimensions {

	if dp.PrevBtn.Clicked(gtx) {
		switch dp.ViewMode {
		case "date":
			dp.Date = dp.Date.AddDate(0, -1, 0)
		case "month":
			dp.Date = dp.Date.AddDate(-1, 0, 0)
		case "year":
			dp.YearRange -= 20
		}
	}
	if dp.NextBtn.Clicked(gtx) {
		switch dp.ViewMode {
		case "date":
			dp.Date = dp.Date.AddDate(0, 1, 0)
		case "month":
			dp.Date = dp.Date.AddDate(1, 0, 0)
		case "year":
			dp.YearRange += 20
		}
	}
	if dp.MonthBtn.Clicked(gtx) {
		dp.ViewMode = "month"
	}
	if dp.YearBtn.Clicked(gtx) {
		dp.ViewMode = "year"
		dp.YearRange = dp.Date.Year() - 10
	}

	gtx.Constraints.Min.Y = 300
	return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return widget.Border{
			Color:        BlackColor,
			CornerRadius: unit.Dp(6),
			// Width:        unit.Dp(1),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return util.LayoutButton(gtx, th, util.Button{
								Text:            "<",
								TextColor:       BlackColor,
								Size:            20,
								FontWeight:      font.Bold,
								BackgroundColor: Transparent,
								BorderColor:     GrayColor,
								CornerRadius:    4,
								Button:          &dp.PrevBtn,
								InInset:         layout.Inset{Left: 15, Right: 15, Top: 5, Bottom: 5},
							})

						}),
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							return layout.Spacer{}.Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							btnText = ""
							switch dp.ViewMode {
							case "date":
								btnText = dp.Date.Format("January 2006")
							case "month":
								btnText = dp.Date.Format("2006")
							case "year":
								btnText = fmt.Sprintf("%d-%d", dp.YearRange, dp.YearRange+19)
							}
							return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return util.LayoutButton(gtx, th, util.Button{
										Text:            dp.Date.Format("Jan"),
										TextColor:       BlackColor,
										Size:            20,
										FontWeight:      font.Bold,
										BackgroundColor: Transparent,
										BorderColor:     GrayColor,
										CornerRadius:    4,
										Button:          &dp.MonthBtn,
										InInset:         layout.Inset{Left: 15, Right: 15, Top: 5, Bottom: 5},
									})
								}),
								layout.Rigid(layout.Spacer{Width: unit.Dp(5)}.Layout),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return util.LayoutButton(gtx, th, util.Button{
										Text:            dp.Date.Format("2006"),
										TextColor:       BlackColor,
										Size:            20,
										FontWeight:      font.Bold,
										BackgroundColor: Transparent,
										BorderColor:     GrayColor,
										CornerRadius:    4,
										Button:          &dp.YearBtn,
										InInset:         layout.Inset{Left: 15, Right: 15, Top: 5, Bottom: 5},
									})
								}),
							)
						}),
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							return layout.Spacer{}.Layout(gtx)
						}),

						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return util.LayoutButton(gtx, th, util.Button{
								Text:            ">",
								TextColor:       BlackColor,
								Size:            20,
								FontWeight:      font.Bold,
								BackgroundColor: Transparent,
								BorderColor:     GrayColor,
								CornerRadius:    4,
								Button:          &dp.NextBtn,
								InInset:         layout.Inset{Left: 15, Right: 15, Top: 5, Bottom: 5},
							})

						}),
					)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					switch dp.ViewMode {
					case "date":
						return dp.daysGrid(gtx, th)
					case "month":
						return dp.monthGrid(gtx, th)
					case "year":
						return dp.yearGrid(gtx, th)
					default:
						return layout.Dimensions{}
					}
				}),
			)
		})
	})

}

var hoverbg color.NRGBA
var textcolor color.NRGBA
var bordercolor color.NRGBA

func (dp *DatePicker) daysGrid(gtx layout.Context, th *material.Theme) layout.Dimensions {
	firstDay := time.Date(dp.Date.Year(), dp.Date.Month(), 1, 0, 0, 0, 0, dp.Date.Location())
	daysInMonth := 32 - time.Date(dp.Date.Year(), dp.Date.Month(), 32, 0, 0, 0, 0, dp.Date.Location()).Day()
	totalspaceX := gtx.Constraints.Max.X
	totalspaceY := gtx.Constraints.Max.Y - 100
	startOffset := int(firstDay.Weekday())
	if startOffset == 0 {
		startOffset = 7
	}
	startOffset--

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return DateSpacedLayout(gtx, th, "Mon", totalspaceX, totalspaceY)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return DateSpacedLayout(gtx, th, "Tue", totalspaceX, totalspaceY)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return DateSpacedLayout(gtx, th, "Wed", totalspaceX, totalspaceY)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return DateSpacedLayout(gtx, th, "Thu", totalspaceX, totalspaceY)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return DateSpacedLayout(gtx, th, "Fri", totalspaceX, totalspaceY)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return DateSpacedLayout(gtx, th, "Sat", totalspaceX, totalspaceY)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return DateSpacedLayout(gtx, th, "Sun", totalspaceX, totalspaceY)
				}),
			)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			var children []layout.FlexChild
			day := 1
			for week := 0; week < 6; week++ {
				children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							var weekDays []layout.FlexChild
							for weekday := 0; weekday < 7; weekday++ {
								if (week == 0 && weekday < startOffset) || day > daysInMonth {
									gtx.Constraints.Min.X = totalspaceX / 7
									gtx.Constraints.Max.X = totalspaceX / 7
									gtx.Constraints.Min.Y = totalspaceY / 7
									gtx.Constraints.Max.Y = totalspaceY / 7
									weekDays = append(weekDays, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										return D{Size: image.Point{X: totalspaceX / 7, Y: totalspaceY / 7}}
									}))
								} else {
									currentDay := day
									currentDate := time.Date(dp.Date.Year(), dp.Date.Month(), currentDay, 0, 0, 0, 0, dp.Date.Location())

									weekDays = append(weekDays, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										if dp.Days[currentDay-1].Clicked(gtx) {
											dp.Date = currentDate
											dp.IsOpen = false
										}

										textcolor = BlackColor

										if dp.Days[currentDay-1].Hovered() {
											hoverbg = MGrayColor
											pointer.CursorPointer.Add(gtx.Ops)
										} else {
											hoverbg = Transparent
											pointer.CursorDefault.Add(gtx.Ops)
										}

										bordercolor = Transparent
										if currentDate.Equal(dp.Date) {
											bordercolor = RedColor
										}

										gtx.Constraints.Max.X = totalspaceX / 7
										gtx.Constraints.Min.X = totalspaceX / 7
										gtx.Constraints.Min.Y = totalspaceY / 7
										gtx.Constraints.Max.Y = totalspaceY / 7
										return util.LayoutButton(gtx, th, util.Button{
											Text:            fmt.Sprintf("%d", currentDay),
											Button:          &dp.Days[currentDay-1],
											TextColor:       textcolor,
											Size:            12,
											FontWeight:      font.Bold,
											BackgroundColor: hoverbg,
											CornerRadius:    4,
											BorderColor:     bordercolor,
										})
									}))
									day++
								}
							}
							return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx, weekDays...)
						}),
					)
				}))
				if day > daysInMonth {
					break
				}
			}
			return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceEvenly}.Layout(gtx, children...)
		}),
	)
}

func DateSpacedLayout(gtx C, th *mt, name string, totalspaceX int, totalspaceY int) D {
	gtx.Constraints.Max.X = totalspaceX / 7
	gtx.Constraints.Min.X = totalspaceX / 7
	gtx.Constraints.Min.Y = totalspaceY / 7
	gtx.Constraints.Max.Y = totalspaceY / 7
	return layout.Center.Layout(gtx, func(gtx C) D {
		return util.LayoutText(gtx, th, util.Text{
			Text:       name,
			Size:       14,
			FontWeight: font.Bold,
			TextColor:  BlackColor,
		})
	})
}

func (dp *DatePicker) monthGrid(gtx layout.Context, th *material.Theme) layout.Dimensions {
	totalspaceX := gtx.Constraints.Max.X
	totalspaceY := gtx.Constraints.Max.Y - 100
	months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			var rows []layout.FlexChild
			for row := 0; row < 3; row++ {
				rows = append(rows, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					var monthButtons []layout.FlexChild
					for col := 0; col < 4; col++ {
						i := row*4 + col
						monthButtons = append(monthButtons, layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							if dp.Months[i].Clicked(gtx) {
								dp.Date = time.Date(dp.Date.Year(), time.Month(i+1), 1, 0, 0, 0, 0, dp.Date.Location())
								dp.ViewMode = "date"
							}
							if dp.Months[i].Hovered() {
								hoverbg = MGrayColor
								pointer.CursorPointer.Add(gtx.Ops)
							} else {
								hoverbg = Transparent
								pointer.CursorDefault.Add(gtx.Ops)
							}
							gtx.Constraints.Min.X = totalspaceX / 4
							gtx.Constraints.Max.X = totalspaceX / 4
							gtx.Constraints.Min.Y = totalspaceY / 3
							gtx.Constraints.Max.Y = totalspaceY / 3
							return util.LayoutButton(gtx, th, util.Button{
								Text:            months[i],
								Button:          &dp.Months[i],
								TextColor:       BlackColor,
								Size:            14,
								FontWeight:      font.Bold,
								BackgroundColor: hoverbg,
								CornerRadius:    4,
							})
						}))
					}
					return layout.Flex{Alignment: layout.Middle, Spacing: layout.SpaceEvenly}.Layout(gtx, monthButtons...)
				}))
			}
			return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceEvenly}.Layout(gtx, rows...)
		}),
	)
}

func (dp *DatePicker) yearGrid(gtx layout.Context, th *material.Theme) layout.Dimensions {
	totalSpaceX := gtx.Constraints.Max.X
	totalSpaceY := gtx.Constraints.Max.Y - 100

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			var rows []layout.FlexChild
			for row := 0; row < 5; row++ {
				rows = append(rows, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					var yearButtons []layout.FlexChild
					for col := 0; col < 4; col++ {
						i := row*4 + col
						year := dp.YearRange + i
						yearButtons = append(yearButtons, layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							if dp.Years[i].Clicked(gtx) {
								dp.Date = time.Date(year, dp.Date.Month(), 1, 0, 0, 0, 0, dp.Date.Location())
								dp.ViewMode = "date"
							}
							if dp.Years[i].Hovered() {
								hoverbg = MGrayColor
								pointer.CursorPointer.Add(gtx.Ops)
							} else {
								hoverbg = Transparent
								pointer.CursorDefault.Add(gtx.Ops)
							}
							gtx.Constraints.Min.X = totalSpaceX / 4
							gtx.Constraints.Max.X = totalSpaceX / 4
							gtx.Constraints.Min.Y = totalSpaceY / 5
							gtx.Constraints.Max.Y = totalSpaceY / 5
							return util.LayoutButton(gtx, th, util.Button{
								Text:            strconv.Itoa(year),
								Button:          &dp.Years[i],
								TextColor:       BlackColor,
								Size:            14,
								FontWeight:      font.Bold,
								BackgroundColor: hoverbg,
								CornerRadius:    4,
							})
						}))
					}
					return layout.Flex{Alignment: layout.Middle, Spacing: layout.SpaceEvenly}.Layout(gtx, yearButtons...)
				}))
			}
			return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceEvenly}.Layout(gtx, rows...)
		}),
	)
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
