package ghost

import (
	"fmt"
	"log"
	"os"
)

type logLv uint8

const (
	logPrefix = "[ghost]: "

	logI logLv = iota
	logW
	logE
	logC

	colf = "%c[%d;;%dm%s%c[0m"
)

func init() {
	log.SetPrefix(logPrefix)
}

// Info wrapper
func Info(v ...interface{}) {
	logger(logI, v...)
}

// Warn wrapper
func Warn(v ...interface{}) {
	logger(logW, v...)
}

// Error wrapper
func Error(v ...interface{}) {
	// log.Panicln(v...)
	panic(logger(logE, v...))
}

// ErrorInDefer without panic
func ErrorInDefer(v ...interface{}) {
	logger(logE, v...)
}

// Crash wrapper
func Crash(code int, v ...interface{}) {
	// log.Fatalln(v...)
	logger(logC, v...)
	os.Exit(code)
}

func logger(lv logLv, v ...interface{}) (s string) {
	s = fmt.Sprintln(v...)
	tag := ""
	switch lv {
	case logI:
		// tag = fmt.Sprintf(colf, 0x1B, 1, 34, "[debug]", 0x1B)
		tag = fmt.Sprintf(colf, 0x1B, 0, 32, " [info]", 0x1B)
	case logW:
		tag = fmt.Sprintf(colf, 0x1B, 1, 33, " [warn]", 0x1B)
	case logE:
		tag = fmt.Sprintf(colf, 0x1B, 1, 31, "[error]", 0x1B)
	case logC:
		tag = fmt.Sprintf(colf, 0x1B, 1, 31, "[crash]", 0x1B)
	default:
	}
	log.Print(tag, s)
	return
}
