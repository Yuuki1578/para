package logger

import (
	"log"
	"os"
)

var (
	Stdout *log.Logger = log.New(os.Stdout, "[PARA INFO]: ", log.Ltime)
	Stderr *log.Logger = log.New(os.Stderr, "[PARA WARN]: ", log.Ltime)
)

func Print(v ...any) {
	Stdout.Print(v...)
}

func Printf(format string, v ...any) {
	Stdout.Printf(format, v...)
}

func Println(v ...any) {
	Stdout.Println(v...)
}

func Eprint(v ...any) {
	Stderr.Print(v...)
}

func Eprintf(format string, v ...any) {
	Stderr.Printf(format, v...)
}

func Eprintln(v ...any) {
	Stderr.Println(v...)
}
