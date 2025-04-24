package main

import (
	"fmt"
	"os"
	"github.com/marcusolsson/tui-go"
)

func main() {
	ui, err := tui.New(tui.NewVBox(
		tui.NewLabel("Huffman Coding Visualization Tool"),
		tui.NewLabel("Press 'q' to quit."),
	))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create UI: %v\n", err)
		os.Exit(1)
	}

	ui.SetKeybinding("q", func() { ui.Quit() })

	if err := ui.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to run UI: %v\n", err)
		os.Exit(1)
	}
}
