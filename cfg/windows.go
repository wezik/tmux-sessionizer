//go:build windows

package cfg

import (
	"os"
	"path/filepath"
)

var (
	dll32         = filepath.Join(os.Getenv("SYSTEMROOT"), "System32", "rundll32.exe")
	defaultEditor = dll32 + " url.dll,FileProtocolHandler"
)
