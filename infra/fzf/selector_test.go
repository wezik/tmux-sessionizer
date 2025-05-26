package fzf_test

import (
	"errors"
	"os/exec"
	. "phopper/dom/model"
	. "phopper/dom/utils"
	. "phopper/infra/fzf"
	"slices"
	"testing"
)

type MockCommandExecutor struct {
	ExecuteParam1   *exec.Cmd
	ExecuteCalls    int
	ExecuteReturn   string
	ExecuteErr      error
	ExecuteExitCode int
}

func (m *MockCommandExecutor) Execute(cmd *exec.Cmd) (string, int, error) {
	m.ExecuteParam1 = cmd
	m.ExecuteCalls++
	return m.ExecuteReturn, m.ExecuteExitCode, m.ExecuteErr
}

func Test_FzfSelector(t *testing.T) {
	t.Run("select from items", func(t *testing.T) {
		// given
		executor := &MockCommandExecutor{}
		executor.ExecuteReturn = "bar\n" // should contain new line char due to input buffer
		selector := NewFzfSelector(executor)

		prompt := "foo prompt > "
		expectedCommand := exec.Command("fzf", "--prompt", prompt)

		// when
		selected, err := selector.SelectFrom([]string{"foo", "bar", "baz"}, prompt)

		// then
		Assert(t, err == nil, "Error should be nil")
		Assert(t, executor.ExecuteCalls == 1, "Execute should be called once")

		cmdParam := executor.ExecuteParam1
		Assert(t, slices.Equal(expectedCommand.Args, cmdParam.Args), "Execute param should be %s is %s", expectedCommand, cmdParam)
		Assert(t, selected == "bar", "Selected item should be %s is %s", "bar", selected)
	})

	t.Run("select maps exit code 130 to ErrSelectorCancelled", func(t *testing.T) {
		// given
		executor := &MockCommandExecutor{}
		executor.ExecuteReturn = "foo\n"
		executor.ExecuteExitCode = 130

		selector := NewFzfSelector(executor)

		prompt := "foo prompt > "

		// when
		_, err := selector.SelectFrom([]string{"foo", "bar", "baz"}, prompt)

		// then
		Assert(t, err == ErrSelectorCancelled, "Error should be %s is %s", ErrSelectorCancelled, err)
		Assert(t, executor.ExecuteCalls == 1, "Execute should be called once")
	})

	t.Run("select propagates errors", func(t *testing.T) {
		// given
		executor := &MockCommandExecutor{}
		expectedErr := errors.New("unknown error")
		executor.ExecuteErr = expectedErr

		selector := NewFzfSelector(executor)

		prompt := "foo prompt > "

		// when
		_, err := selector.SelectFrom([]string{"foo", "bar", "baz"}, prompt)

		// then
		Assert(t, expectedErr == err, "Error should be %s is %s", expectedErr, err)
		Assert(t, executor.ExecuteCalls == 1, "Execute should be called once")
	})
}
