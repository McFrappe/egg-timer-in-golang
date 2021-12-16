package main

import (
    "gioui.org/app"
    "gioui.org/unit"
    "gioui.org/font/gofont"
    "gioui.org/io/system"
    "gioui.org/layout"
    "gioui.org/op"
    "gioui.org/widget"
    "gioui.org/widget/material"
)

func main() {
    go func() {
        // create a new window
        window := app.NewWindow(
            app.Title("Egg Timer"),
            app.MaxSize(unit.Dp(400), unit.Dp(600)),
            app.MinSize(unit.Dp(400), unit.Dp(600)),
        )

        // ops are for operations from the UI
        var ops op.Ops
        var startButton widget.Clickable
        theme := material.NewTheme(gofont.Collection())

        // listen for events in the window
        for event := range window.Events() {
            switch event := event.(type) {
                case system.FrameEvent:
                    graphical_ctx := layout.NewContext(&ops, event)
                    button := material.Button(theme, &startButton, "Start")
                    button.Layout(graphical_ctx)
                    event.Frame(graphical_ctx.Ops)
            }
        }
    }()

    app.Main()
}
