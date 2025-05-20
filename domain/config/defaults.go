package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type Config struct {
	Editor string
}

// TODO this should be loaded from some sort of config file
const settingsEditor = ""

func GetDefaults() Config {
	editor := settingsEditor
	if editor == "" {
		editor = prepareLaunchCmd(runtime.GOOS)
	}

	return Config{
		Editor: editor,
	}
}

func prepareLaunchCmd(current_os string) string {
	switch current_os {
	case "darwin":
		return "open"
	case "linux":
		return "xdg-open"
	case "windows":
		cmd := "url.dll,FileProtocolHandler"
		runDll32 := filepath.Join(os.Getenv("SYSTEMROOT"), "System32", "rundll32.exe")
		return runDll32 + " " + cmd
	default:
		fmt.Println("Unknown OS:", current_os, "falling back to vi as an editor")
		return "vi"
	}
}
