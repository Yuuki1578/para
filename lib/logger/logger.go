package logger

import (
	"log"
	"os"
)

var (
	stdout *log.Logger = log.New(os.Stdout, "[PARA INFO]: ", log.Ltime)
	stderr *log.Logger = log.New(os.Stderr, "[PARA WARN]: ", log.Ltime)
)

func Print(v ...any) {
	stdout.Print(v...)
}

func Printf(format string, v ...any) {
	stdout.Printf(format, v...)
}

func Println(v ...any) {
	stdout.Println(v...)
}

func Eprint(v ...any) {
	stderr.Print(v...)
}

func Eprintf(format string, v ...any) {
	stderr.Printf(format, v...)
}

func Eprintln(v ...any) {
	stderr.Println(v...)
}
