package main

import (
	"embed"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	os.Setenv("GTK_THEME", "Adwaita:dark")

	app := NewApp()

	err := wails.Run(&options.App{
		Title:            "Popeye",
		Width:            1200,
		Height:           800,
		MinWidth:         800,
		MinHeight:        600,
		DisableResize:    false,
		Frameless:        false,
		StartHidden:      false,
		AlwaysOnTop:      false,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 13, G: 17, B: 23, A: 255},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		Linux: &linux.Options{
			ProgramName: "Popeye",
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
