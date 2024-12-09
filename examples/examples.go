package main

import (
	"errors"
	"os"
	"time"

	"github.com/leonidasdeim/log"
)

func defaultLogger() {
	l := log.New()
	l.Info("I'm default logger. I log to os.Stdout in colourful text format")
	l.Infof("Formatted text: %d %t", 123, true)
}

func customLevel() {
	l := log.New(log.WithName("cust"), log.WithLevel(log.Warning))
	l.Warning("I only logging Warning level and above", "level", l.Level())
	l.Info("This will not be writter to the log")
}

func childLogger() {
	main := log.New(log.WithName("main"))
	main.Info("Main module logger")
	child1 := main.NewLocal(log.WithName("child1"), log.WithLevel(log.Debug))
	child1.Debug("Child module logger, I keep global settings from main logger, but can define my own log level and name")

	f, err := os.Create("example.log")
	if err != nil {
		main.Err("create file", err)
	}
	defer f.Close()

	child2 := main.NewLocal(log.WithName("child2"), log.WithLevel(log.Error), log.WithWriter(f, log.FormatJson))
	child2.Error("Another child logger which also can log to file")
}

func jsonLogger() {
	l := log.New(log.WithName("jsonMod"), log.WithWriter(os.Stdout, log.FormatJson))
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
	l.Info("All log lines are buffered")
	l.Info("And writted one by one")
	l.Info("1")
	l.Info("2")
	l.Info("3")
	l.Info("4")
	l.Info("5")
	l.Info("6")
	l.Info("7")
	l.Info("8")
	l.Info("9")
	l.Info("10")
	time.Sleep(100 * time.Millisecond)
}

func propsLogger() {
	l := log.New(log.WithName("props"), log.WithWriter(os.Stdout, log.FormatTextColor))
	l.Warning("I can log in text format with props", "prop1", 123, "prop2", "hello")
}

func errorLogger() {
	l := log.New(log.WithName("props"), log.WithWriter(os.Stdout, log.FormatTextColor))
	err := errors.New("something failed")
	l.Warn("I can log error as a special prop", err)
}

func propsJsonLogger() {
	l := log.New(log.WithName("propsJson"), log.WithWriter(os.Stdout, log.FormatJson))
	l.Warning("I can log in json format with props", "prop1", struct{ Name string }{"Leo"}, "prop2", 1.23)
}
