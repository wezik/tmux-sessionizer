package editor

import (
	"os"
	"os/exec"
	. "phopper/dom/service"
)

type ShellEditorLauncher struct {
	editor string
	exec   CommandExecutor
}

func NewShellEditorLauncher(editor string, exec CommandExecutor) *ShellEditorLauncher {
	return &ShellEditorLauncher{editor: editor, exec: exec}
}

func (e *ShellEditorLauncher) Open(path string) error {
	cmd := exec.Command(e.editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	_, err := e.exec.ExecuteInteractive(cmd)
	return err
}
