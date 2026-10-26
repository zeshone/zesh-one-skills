package main

import (
	"io"
	"os"

	"golang.org/x/sys/windows"
	"golang.org/x/term"
)

func preparePickerOutput(out io.Writer) (func() error, error) {
	file, ok := out.(*os.File)
	if !ok || !term.IsTerminal(int(file.Fd())) {
		return func() error { return nil }, nil
	}
	handle := windows.Handle(file.Fd())
	var original uint32
	if err := windows.GetConsoleMode(handle, &original); err != nil {
		return nil, err
	}
	if err := windows.SetConsoleMode(handle, original|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING); err != nil {
		return nil, err
	}
	return func() error { return windows.SetConsoleMode(handle, original) }, nil
}
