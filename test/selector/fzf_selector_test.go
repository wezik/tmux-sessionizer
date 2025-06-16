package selector_test

import (
	"bytes"
	"errors"
	"os/exec"
	"testing"
	"thop/internal/selector"
	"thop/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_SelectFrom(t *testing.T) {
	t.Run("selects from items", func(t *testing.T) {
		// given
		prompt := "foo prompt > "
		args := []string{"fzf", "--prompt", prompt}
		var cmdToExec *exec.Cmd
		cmdResult := "bar\n" // new line char due to input buffer separation needed by fzf

		execMock := new(test.MockExecutor)
		execMock.On("Execute", mock.Anything).Run(func(args mock.Arguments) {
			cmdToExec = args.Get(0).(*exec.Cmd)
		}).Return(cmdResult, 0, nil).Once()

		selector := selector.FzfSelector{E: execMock}

		// when
		s, err := selector.SelectFrom([]string{"foo", "bar", "baz"}, prompt)

		// then
		assert.Nil(t, err)
		assert.Equal(t, "bar", s)
		execMock.AssertExpectations(t)

		assert.Equal(t, cmdToExec.Args, args)

		stdin := cmdToExec.Stdin.(*bytes.Buffer)
		assert.Equal(t, "bar\nbaz\nfoo\n", stdin.String(), "stdin should be sorted")
	})

	t.Run("select maps 130 exit code to ErrSelectorCancelled", func(t *testing.T) {
		// given
		execMock := new(test.MockExecutor)
		execMock.On("Execute", mock.Anything).Return("foo", 130, nil).Once()

		s := selector.FzfSelector{E: execMock}

		// when
		_, err := s.SelectFrom([]string{"foo", "bar", "baz"}, "foo prompt > ")

		// then
		assert.True(t, selector.ErrSelectorCancelled.Equal(err))
		execMock.AssertExpectations(t)
	})

	t.Run("select propagates errors", func(t *testing.T) {
		// given
		execMock := new(test.MockExecutor)
		expectedErr := errors.New("unknown error")
		execMock.On("Execute", mock.Anything).Return("", 0, expectedErr).Once()

		s := selector.FzfSelector{E: execMock}

		// when
		_, err := s.SelectFrom([]string{"foo", "bar", "baz"}, "foo prompt > ")

		// then
		assert.True(t, selector.ErrSelectorFailed.Equal(err))
		execMock.AssertExpectations(t)
	})
}
