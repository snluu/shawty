package main

import "fmt"
import "time"

// Log prints to the standard output with the timestamp
func log(level string, a ...interface{}) {
	fmt.Printf("[%v] %s: ", time.Now(), level)
	fmt.Println(a...)
}

// Logf prints to the standard output with the timestamp
func logf(level, format string, args ...interface{}) {
	fmt.Printf("[%v] %s: ", time.Now(), level)
	fmt.Printf(format, args...)
	fmt.Println()
}

func Linfo(a ...interface{}) {
	log("INFO", a...)
}

func Linfof(format string, args ...interface{}) {
	logf("INFO", format, args...)
}

func Lerror(a ...interface{}) {
	log("ERROR", a...)
}

func Lerrorf(format string, args ...interface{}) {
	logf("ERROR", format, args...)
}

func Ldebug(a ...interface{}) {
	log("DEBUG", a...)
}

func Ldebugf(format string, args ...interface{}) {
	logf("DEBUG", format, args...)
}
