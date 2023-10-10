package main

import (
	"os"

	"github.com/leonidasdeim/log"
)

func defaultLogger() {
	l := log.New()
	l.Info("I'm default logger. I log to os.Stdout in colourful text format")
}

func customLevel() {
	l := log.New(log.WithName("cust"), log.WithLevel(log.Warning))
	l.Warning("I only logging Warning level and above")
	l.Info("This will not be writter to the log")
}

func childLogger() {
	main := log.New(log.WithName("main"))
	main.Info("Main module logger")
	child1 := main.NewLocal(log.WithName("child1"), log.WithLevel(log.Debug))
	child1.Debug("Child module logger, I keep global settings from main logger, but can define my own log level and name")

	f, err := os.Create("example.log")
	if err != nil {
		main.Error(err.Error())
	}
	defer f.Close()

	child2 := main.NewLocal(log.WithName("child2"), log.WithLevel(log.Error), log.WithWriter(f, log.FormatJson))
	child2.Error("Another child logger which also can log to file")
}

func jsonLogger() {
	l := log.New(log.WithName("json"), log.WithWriter(os.Stdout, log.FormatJson))
	l.Warning("I can log in JSON format")
}

func simpleTextLogger() {
	l := log.New(log.WithName("text"), log.WithWriter(os.Stdout, log.FormatText))
	l.Warning("I can log in simple text format")
}

func nonBlocking() {
	l := log.New(log.WithName("nonblck"), log.WithMode(log.ModeNonBlocking))
	defer l.Close()
	l.Info("I'm working in Non Blocking mode, so I need Sync() before application exits")
	l.Info("All logs are written in goroutines")
}
