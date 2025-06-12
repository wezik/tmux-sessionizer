package editor

import (
	"os"
	"os/exec"
	. "thop/dom/service"

	"github.com/dsnet/try"
)

type ShellEditorLauncher struct {
	editor string
	exec   CommandExecutor
}

func NewShellEditorLauncher(editor string, exec CommandExecutor) *ShellEditorLauncher {
	return &ShellEditorLauncher{editor: editor, exec: exec}
}

func (e *ShellEditorLauncher) Open(path string) (err error) {
	defer try.Handle(&err)

	cmd := exec.Command(e.editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	try.E1(e.exec.ExecuteInteractive(cmd))

	return nil
}
