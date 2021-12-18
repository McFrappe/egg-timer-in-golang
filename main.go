package main

import (
  "log"
  "os"
  "image/color"
  "image"
  "gioui.org/f32"
  "gioui.org/op/clip"
  "gioui.org/op/paint"
  "time"
  "math"
  "fmt"
  "strconv"
  "strings"

  "gioui.org/app"
  "gioui.org/unit"
  "gioui.org/font/gofont"
  "gioui.org/io/system"
  "gioui.org/layout"
  "gioui.org/op"
  "gioui.org/widget"
  "gioui.org/widget/material"
  "gioui.org/text"
)

type C = layout.Context
type D = layout.Dimensions

// rgba(11, 232, 129,1.0)
var emerald_color = color.NRGBA {
  R: 255,
  G: 232,
  B: 129,
  A: 1.0,
}

// Progressvariables, a channel and a variable
var progress_incrementer chan float32
var progress_value float32

// textfield to input boil duration
var boil_duration_input widget.Editor
var is_boiling bool
var boil_duration float32

func egg_drawing(gtx C) D {
    // Draw a custom path, shaped like an egg
    var egg_path clip.Path
    op.Offset(f32.Pt(200, 150)).Add(gtx.Ops)
    egg_path.Begin(gtx.Ops)
    // Rotate from 0 to 360 degrees
    for deg := 0.0; deg <= 360; deg++ {

      rad := deg / 360 * 2 * math.Pi
      // Trig gives the distance in X and Y direction
      cosT := math.Cos(rad)
      sinT := math.Sin(rad)
      // Constants to define the eggshape
      a := 110.0
      b := 150.0
      d := 20.0
      // The x/y coordinates
      x := a * cosT
      y := -(math.Sqrt(b*b - d*d*cosT*cosT) + d*sinT) * sinT

      // the point on the outline
      p := f32.Pt(float32(x), float32(y))
      // Draw the line to this point
      egg_path.LineTo(p)
    }
    // Close the path
    egg_path.Close()

    // Get hold of the actual clip
    egg_area := clip.Outline{Path: egg_path.End()}.Op()

    // Fill the shape
    // color := color.NRGBA{R: 255, G: 239, B: 174, A: 255}
    color := color.NRGBA{R: 255, G: uint8(239 * (1 - progress_value)), B: uint8(174 * (1 - progress_value)), A: 255}
    paint.FillShape(gtx.Ops, color, egg_area)

    d := image.Point{Y: 375}
    return layout.Dimensions{Size: d}
}

func draw(window * app.Window) error {
  // ops are for operations from the UI
  var ops op.Ops
  var start_button widget.Clickable
  theme := material.NewTheme(gofont.Collection())

  // listen for events in the window
  for {
    select {
      case event := <-window.Events():
          switch event := event.(type) {
            case system.FrameEvent:
              // Create a context to place objects onto
              graphical_ctx := layout.NewContext(&ops, event)

              if start_button.Clicked() {
                is_boiling = !is_boiling

                if progress_value >= 1 {
                    progress_value = 0
                }

                input_string := boil_duration_input.Text()
                input_string = strings.TrimSpace(input_string)
                input_float, _ := strconv.ParseFloat(input_string, 32)

                boil_duration = float32(input_float)
                boil_duration = boil_duration / (1 - progress_value)
              }

              layout.Flex {
                Axis: layout.Vertical,
                Spacing: layout.SpaceStart,
              }.Layout(graphical_ctx,
                // The egg drawing
                layout.Rigid(
                    func(graphical_ctx C) D {
                        return egg_drawing(graphical_ctx)
                    },
                ),
                // The progressbar
                layout.Rigid(
                  func(graphical_ctx C) D {
                    return material.ProgressBar(
                      theme,
                      progress_value).Layout(
                      graphical_ctx)
                  },
                ),
                // The input field
                layout.Rigid(
                    func(graphical_ctx C) D {
                        editor := material.Editor(theme, &boil_duration_input, "seconds")

                        boil_duration_input.SingleLine = true
                        boil_duration_input.Alignment = text.Middle
                        
                        margins := layout.Inset{
                            Top: unit.Dp(0),
                            Right: unit.Dp(170),
                            Bottom: unit.Dp(40),
                            Left: unit.Dp(170),
                        }

                        border := widget.Border{
                            Color: color.NRGBA{R: 204, G: 204, B: 204, A: 255},
                            CornerRadius: unit.Dp(3),
                            Width: unit.Dp(2),
                        }

                        return margins.Layout(graphical_ctx, 
                            func(graphical_ctx C) D {
                                return border.Layout(graphical_ctx, editor.Layout)
                            },
                        )
                    },
                ),
                // The button
                layout.Rigid(
                  func(graphical_ctx C) D {
                    margins := layout.Inset {
                      Top: unit.Dp(25),
                      Bottom: unit.Dp(25),
                      Right: unit.Dp(35),
                      Left: unit.Dp(35),
                    }

                    var button_text string
                    if !is_boiling {
                      button_text = "Start"
                    } else if is_boiling && progress_value < 1 {
                      button_text = "Stop"
                    } else if is_boiling && progress_value >= 1 {
                        button_text = "Finished"
                    }

                    return margins.Layout(
                      graphical_ctx,
                      func(graphical_ctx C) D {
                        return material.Button(theme, &start_button, button_text).Layout(graphical_ctx)
                    }, 
                    )
                  },
                ),
              )
              event.Frame(graphical_ctx.Ops)

            case system.DestroyEvent:
              return event.Err
          }
      case p := <-progress_incrementer:
          if is_boiling && progress_value < 1 {
              progress_value += p

              boil_remain := (1 - progress_value) * boil_duration
              input_str := fmt.Sprintf("%.1f", math.Round(float64(boil_remain) * 10) / 10)

              boil_duration_input.SetText(input_str)
              window.Invalidate()
          }
    }
  }
}

func main() {
  progress_incrementer = make(chan float32)

  go func() {
    for {
      time.Sleep(time.Second / 25)
      progress_incrementer <-0.004
    }
  }()

  go func() {
    // create a new window
    window := app.NewWindow(
      app.Title("Egg Timer"),
      app.MaxSize(unit.Dp(400), unit.Dp(600)),
      app.MinSize(unit.Dp(400), unit.Dp(600)),
    )

    if err := draw(window);err != nil {
      log.Fatal(err)
    }

    os.Exit(0)
  }()

  app.Main()
}
