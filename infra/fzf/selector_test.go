package fzf_test

import (
	"bytes"
	"errors"
	"os/exec"
	"testing"
	. "thop/dom/model"
	. "thop/infra/fzf"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCommandExecutor struct {
	mock.Mock
}

func (m *MockCommandExecutor) Execute(cmd *exec.Cmd) (string, int, error) {
	args := m.Called(cmd)
	return args.String(0), args.Int(1), args.Error(2)
}

func (m *MockCommandExecutor) ExecuteInteractive(cmd *exec.Cmd) (int, error) {
	args := m.Called(cmd)
	return args.Int(0), args.Error(1)
}

func Test_SelectFrom(t *testing.T) {
	t.Run("selects from items", func(t *testing.T) {
		// given
		prompt := "foo prompt > "
		args := []string{"fzf", "--prompt", prompt}
		var cmdToExec *exec.Cmd
		cmdResult := "bar\n" // new line char due to input buffer separation needed by fzf

		execMock := new(MockCommandExecutor)
		execMock.On("Execute", mock.Anything).Run(func(args mock.Arguments) {
			cmdToExec = args.Get(0).(*exec.Cmd)
		}).Return(cmdResult, 0, nil).Once()

		selector := NewFzfSelector(execMock)

		// when
		selected, err := selector.SelectFrom([]string{"foo", "bar", "baz"}, prompt)

		// then
		assert.Nil(t, err)
		assert.Equal(t, "bar", selected)
		execMock.AssertExpectations(t)

		assert.Equal(t, cmdToExec.Args, args)

		stdin := cmdToExec.Stdin.(*bytes.Buffer)
		assert.Equal(t, "bar\nbaz\nfoo\n", stdin.String(), "stdin should be sorted")
	})

	t.Run("select maps 130 exit code to ErrSelectorCancelled", func(t *testing.T) {
		// given
		execMock := new(MockCommandExecutor)
		execMock.On("Execute", mock.Anything).Return("foo", 130, nil).Once()

		selector := NewFzfSelector(execMock)

		// when
		_, err := selector.SelectFrom([]string{"foo", "bar", "baz"}, "foo prompt > ")

		// then
		assert.Equal(t, ErrSelectorCancelled, err)
		execMock.AssertExpectations(t)
	})

	t.Run("select propagates errors", func(t *testing.T) {
		// given
		execMock := new(MockCommandExecutor)
		expectedErr := errors.New("unknown error")
		execMock.On("Execute", mock.Anything).Return("", 0, expectedErr).Once()

		selector := NewFzfSelector(execMock)

		// when
		_, err := selector.SelectFrom([]string{"foo", "bar", "baz"}, "foo prompt > ")

		// then
		assert.Equal(t, expectedErr, err)
		execMock.AssertExpectations(t)
	})
}
