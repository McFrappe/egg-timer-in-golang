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

  "gioui.org/app"
  "gioui.org/unit"
  "gioui.org/font/gofont"
  "gioui.org/io/system"
  "gioui.org/layout"
  "gioui.org/op"
  "gioui.org/widget"
  "gioui.org/widget/material"
)

type C = layout.Context
type D = layout.Dimensions

// rgba(11, 232, 129,1.0)
var emeraldColor = color.NRGBA {
  R: 255,
  G: 232,
  B: 129,
  A: 1.0,
}

// Progressvariables, a channel and a variable
var progress_incrementer chan float32
var progress_value float32
var is_boiling bool

func draw(window * app.Window) error {
  // ops are for operations from the UI
  var ops op.Ops
  var startButton widget.Clickable
  theme := material.NewTheme(gofont.Collection())

  // listen for events in the window
  // for event := range window.Events() {
  for {
    select {
      case event := <-window.Events():
          switch event := event.(type) {
            case system.FrameEvent:
              // Create a context to place objects onto
              graphical_ctx := layout.NewContext(&ops, event)

              if startButton.Clicked() {
                is_boiling = !is_boiling
              }

              layout.Flex {
                Axis: layout.Vertical,
                Spacing: layout.SpaceStart,
              }.Layout(graphical_ctx,
                layout.Rigid(
                    func(graphical_ctx C) D {
                        circle := clip.Circle{
                            Center: f32.Point{X: float32(graphical_ctx.Constraints.Max.X) / 2, Y: 200},
                            Radius: 120,
                        }.Op(graphical_ctx.Ops)
                        color := color.NRGBA{R: 200, A: 255}
                        paint.FillShape(graphical_ctx.Ops, color, circle)
                        d := image.Point{Y: 500}
                        return layout.Dimensions{Size: d}
                    },
                ),
                layout.Rigid(
                  func(graphical_ctx C) D {
                    return material.ProgressBar(
                      theme,
                      progress_value).Layout(
                      graphical_ctx)
                  },
                ),
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
                    } else {
                      button_text = "Boiling"
                    }

                    return margins.Layout(
                      graphical_ctx,
                      func(graphical_ctx C) D {
                        return material.Button(theme, &startButton, button_text).Layout(graphical_ctx)}, 
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
