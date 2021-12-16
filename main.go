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
                    // Create a context to place objects onto
                    graphical_ctx := layout.NewContext(&ops, event)

                    layout.Flex{
                        Axis: layout.Vertical,
                        Spacing: layout.SpaceStart,
                    }.Layout(graphical_ctx,
                        layout.Rigid(
                            // Create button with pointer from the clickable widget 'startButton'
                            func(graphical_ctx layout.Context) layout.Dimensions {
                                button := material.Button(theme, &startButton, "Start")
                                return button.Layout(graphical_ctx)
                            },
                        ),
                        layout.Rigid(
                            layout.Spacer{Height: unit.Dp(30)}.Layout,
                        ),
                    )

                    event.Frame(graphical_ctx.Ops)

            }
        }
    }()

    app.Main()
}
