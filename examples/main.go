package main

import (
	"fmt"

	"github.com/leonidasdeim/log"
)

func main() {
	defaultLogger()
	fmt.Println("-")

	customLevel()
	fmt.Println("-")

	childLogger()
	fmt.Println("-")

	jsonLogger()
	fmt.Println("-")

	nonBlocking()
	fmt.Println("-")
}

func defaultLogger() {
	l := log.New()
	l.Info("I'm default logger. I log to os.Stdout in text format")
}

func customLevel() {
	l := log.New(log.WithName("cust"), log.WithLevel(log.Warning))
	l.Warning("I only logging Warning level and above")
	l.Info("This will not be writter to the log")
}

func childLogger() {
	main := log.New(log.WithName("main"))
	main.Info("Main module logger")
	child1 := main.Local(log.WithName("child1"), log.WithLevel(log.Debug))
	child1.Debug("Child module logger, I keep global settings from main logger, but can define my own log level and name")
	child2 := main.Local(log.WithName("child2"), log.WithLevel(log.Error))
	child2.Error("Another child logger")
}

func jsonLogger() {
	l := log.New(log.WithName("cust"), log.WithFormat(log.FormatJson))
	l.Warning("I can log in JSON format")
}

func nonBlocking() {
	l := log.New(log.WithName("nonblck"), log.WithWriteMode(log.ModeNonBlocking))
	defer l.Sync()
	l.Info("I'm working in Non Blocking mode, so I need Sync() before application exits")
	l.Info("All logs are written in goroutines")
}
