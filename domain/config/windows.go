//go:build windows

package config

import (
	"os"
	"path/filepath"
)

var (
	defaultEditor = filepath.Join(os.Getenv("SYSTEMROOT"), "System32", "rundll32.exe") + " " + "url.dll,FileProtocolHandler"
)
