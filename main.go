package main

import (
    "gioui.org/app"
)

func main() {
    go func() {
        // create a new window
        w := app.NewWindow()

        // listen for events in the window
        for range w.Events() {
        }
    }()

    app.Main()
}
