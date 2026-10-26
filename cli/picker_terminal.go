package main

import (
	"errors"
	"io"
	"os"

	"golang.org/x/term"
)

func preparePickerTerminal(input io.Reader, out io.Writer) (func() error, error) {
	restoreOutput, err := preparePickerOutput(out)
	if err != nil {
		return nil, err
	}
	restoreInput := func() error { return nil }
	if file, ok := input.(*os.File); ok && term.IsTerminal(int(file.Fd())) {
		state, err := term.MakeRaw(int(file.Fd()))
		if err != nil {
			_ = restoreOutput()
			return nil, err
		}
		restoreInput = func() error { return term.Restore(int(file.Fd()), state) }
	}
	return func() error { return errors.Join(restoreInput(), restoreOutput()) }, nil
}
func pickerSize(out io.Writer) (int, int) {
	if file, ok := out.(*os.File); ok {
		if width, height, err := term.GetSize(int(file.Fd())); err == nil && width > 0 && height > 0 {
			return width, height
		}
	}
	return 80, 24
}
