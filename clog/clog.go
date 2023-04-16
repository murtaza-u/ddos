package clog

import (
	"log"

	"github.com/murtaza-u/ddos/conf"
)

// wrapper around log.Println
func Println(v ...any) {
	if !conf.Debug {
		return
	}

	log.Println(v...)
}

// wrapper around log.Print
func Print(v ...any) {
	if !conf.Debug {
		return
	}

	log.Print(v...)
}

// wrapper around log.Printf
func Printf(format string, v ...any) {
	if !conf.Debug {
		return
	}

	log.Printf(format, v...)
}

// wrapper around log.Fatal
func Fatal(v ...any) {
	if !conf.Debug {
		return
	}

	log.Fatal(v...)
}

// wrapper around log.Fatalf
func Fatalf(format string, v ...any) {
	if !conf.Debug {
		return
	}

	log.Fatalf(format, v...)
}
