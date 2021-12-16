package main

import (
    "log"
    "os"
    "image/color"

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
var emeraldColor = color.NRGBA{
    R: 255,
    G: 232, 
    B: 129, 
    A: 1.0,
}

func draw(window *app.Window) error {
    // ops are for operations from the UI
    var ops op.Ops
    var startButton widget.Clickable
    theme := material.NewTheme(gofont.Collection())

    // listen for events in the window
    for event := range window.Events() {
        switch event := event.(type) {
            case system.FrameEvent:
                // Create a context to place objects onto
                graphical_ctx := layout.NewContext(&ops, event)

                layout.Flex{
                    Axis: layout.Vertical,
                    Spacing: layout.SpaceStart,
                }.Layout(graphical_ctx,
                    layout.Rigid(
                    func(graphical_ctx C) D {
                        margins := layout.Inset{
                            Top: unit.Dp(25),
                            Bottom: unit.Dp(25),
                            Right: unit.Dp(35),
                            Left: unit.Dp(35),
                        }

                        return margins.Layout(graphical_ctx, func(graphical_ctx C) D {
                            button := material.Button(theme, &startButton, "Start")
                                // button.Background = emeraldColor
                                return button.Layout(graphical_ctx)
                        },)
                    },
                ),
            )

            event.Frame(graphical_ctx.Ops)

            case system.DestroyEvent:
                return event.Err
        }
    }

    return nil
}

func main() {
    go func() {
        // create a new window
        window := app.NewWindow(
            app.Title("Egg Timer"),
            app.MaxSize(unit.Dp(400), unit.Dp(600)),
            app.MinSize(unit.Dp(400), unit.Dp(600)),
        )

        if err := draw(window); err != nil {
            log.Fatal(err)
        }

        os.Exit(0)
    }()

    app.Main()
}
