package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"

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

// preparePickerSession captures termination before changing terminal state.
func preparePickerSession(input io.Reader, out io.Writer) (func() error, error) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	stop := func() {
		signal.Stop(signals)
		select {
		case sig := <-signals:
			finishPickerSignal(sig)
		default:
		}
	}
	restore, err := preparePickerTerminal(input, out)
	if err != nil {
		stop()
		return nil, err
	}
	fmt.Fprint(out, "\x1b[?1049h\x1b[?25l")
	done := make(chan struct{})
	exited := make(chan struct{})
	var once sync.Once
	var restoreErr error
	cleanup := func() error {
		once.Do(func() {
			fmt.Fprint(out, "\x1b[?25h\x1b[?1049l")
			restoreErr = restore()
		})
		return restoreErr
	}
	go func() {
		defer close(exited)
		select {
		case sig := <-signals:
			_ = cleanup()
			signal.Stop(signals)
			finishPickerSignal(sig)
		case <-done:
		}
	}()
	return func() error {
		err := cleanup()
		stop()
		close(done)
		<-exited
		return err
	}, nil
}

// Reemit after restoration so Unix callers still observe signal termination.
func finishPickerSignal(sig os.Signal) {
	signal.Reset(sig)
	if process, err := os.FindProcess(os.Getpid()); err == nil {
		if err := process.Signal(sig); err == nil {
			select {} // Wait for the default signal handler to terminate the process.
		}
	}
	os.Exit(128 + int(sig.(syscall.Signal)))
}

func pickerSize(out io.Writer) (int, int) {
	if file, ok := out.(*os.File); ok {
		if width, height, err := term.GetSize(int(file.Fd())); err == nil && width > 0 && height > 0 {
			return width, height
		}
	}
	return 80, 24
}
