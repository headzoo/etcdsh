package io

import (
	"fmt"
	"os"
)

var (
	Stdin, Stdout, Stderr = os.Stdin, os.Stdout, os.Stderr
)

// Print formats using the default formats for its operands and writes to standard output.
// Spaces are added between operands when neither is a string.
// It returns the number of bytes written and any write error encountered.
func Print(line string) (int, error) {
	return fmt.Fprint(Stdout, line)
}

// Println formats using the default formats for its operands and writes to standard output.
// Spaces are always added between operands and a newline is appended.
// It returns the number of bytes written and any write error encountered.
func Println(line string) (int, error) {
	return fmt.Fprintln(Stdout, line)
}

// Printf formats according to a format specifier and writes to standard output.
// It returns the number of bytes written and any write error encountered.
func Printf(format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(Stdout, format, a...)
}

// Printfln formats according to a format specifier and writes to standard output.
// Spaces are always added between operands and a newline is appended.
// It returns the number of bytes written and any write error encountered.
func Printfln(format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(Stdout, format+"\n", a...)
}

// Print formats using the default formats for its operands and writes to standard output.
// Spaces are added between operands when neither is a string.
// It returns the number of bytes written and any write error encountered.
func PrintErr(line string) (int, error) {
	return fmt.Fprint(Stderr, line)
}

// Println formats using the default formats for its operands and writes to standard error.
// Spaces are always added between operands and a newline is appended.
// It returns the number of bytes written and any write error encountered.
func PrintErrln(line string) (int, error) {
	return fmt.Fprintln(Stderr, line)
}

// Printf formats according to a format specifier and writes to standard error.
// It returns the number of bytes written and any write error encountered.
func PrintErrf(format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(Stderr, format, a...)
}

// Printfln formats according to a format specifier and writes to standard error.
// Spaces are always added between operands and a newline is appended.
// It returns the number of bytes written and any write error encountered.
func PrintErrfln(format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(Stderr, format+"\n", a...)
}
