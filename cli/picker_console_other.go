//go:build !windows

package main

import "io"

func preparePickerOutput(_ io.Writer) (func() error, error) { return func() error { return nil }, nil }
