package main

import (
	"log"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/hd-buddy/GioCalendarPicker/datepicker"
	"github.com/hd-buddy/GioCalendarPicker/util"
)

func main() {

	go func() {
		w := new(app.Window)

		w.Option(app.Title("Date Picker"))
		if err := Display(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

var dp = datepicker.DatePicker{
	Date:      time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Time.Location(time.Now())),
	Openbtn:   widget.Clickable{},
	IsOpen:    false,
	PrevBtn:   widget.Clickable{},
	NextBtn:   widget.Clickable{},
	Days:      [31]widget.Clickable{},
	MonthBtn:  widget.Clickable{},
	YearBtn:   widget.Clickable{},
	Months:    [12]widget.Clickable{},
	Years:     [20]widget.Clickable{},
	ViewMode:  "date", // "date", "month", or "year"
	YearRange: time.Now().Year() - 10,
	TodayBtn: widget.Clickable{},
	Editor:    &widget.Editor{},
}

type (
	C = layout.Context
	D = layout.Dimensions
)

func Display(w *app.Window) error {
	var ops op.Ops
	th := material.NewTheme()
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			DisplayDatePicker(gtx, th)
			e.Frame(gtx.Ops)
		}
	}
}

func DisplayDatePicker(gtx C, th *material.Theme) D {
	return layout.Flex{
		Axis:      layout.Vertical,
		Alignment: layout.Middle,
		Spacing:   layout.SpaceEnd,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return util.LayoutText(gtx, th, util.Text{
				Text:       "Date Picker in Gio",
				Size:       20,
				TextColor:  util.BlackColor,
				FontWeight: font.Bold,
				Inset:      layout.Inset{Top: 10, Bottom: 10},
				Alignment:  text.Middle,
			})
		}),
		layout.Rigid(func(gtx C) D {
			return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx C) D {
				return dp.Layout(gtx, th)
			})
		}),
	)
}
