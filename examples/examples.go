package main

import (
	"os"
	"time"

	"github.com/leonidasdeim/log"
)

func defaultLogger() {
	l := log.New()
	l.Info().Msg("I'm default logger. I log to os.Stdout in colourful text format")
	l.Info().Msgf("Formatted text: %d %t", 123, true)
}

func customLevel() {
	l := log.New(log.WithName("cust"), log.WithLevel(log.Warning))
	l.Warning().Prop("level", string(l.Level())).Msg("I only logging Warning level and above")
	l.Info().Msg("This will not be writter to the log")
}

func childLogger() {
	main := log.New(log.WithName("main"))
	main.Info().Msg("Main module logger")
	child1 := main.NewLocal(log.WithName("child1"), log.WithLevel(log.Debug))
	child1.Debug().Msg("Child module logger, I keep global settings from main logger, but can define my own log level and name")

	f, err := os.Create("example.log")
	if err != nil {
		main.Error().Msg(err.Error())
	}
	defer f.Close()

	child2 := main.NewLocal(log.WithName("child2"), log.WithLevel(log.Error), log.WithWriter(f, log.FormatJson))
	child2.Error().Msg("Another child logger which also can log to file")
}

func jsonLogger() {
	l := log.New(log.WithName("jsonMod"), log.WithWriter(os.Stdout, log.FormatJson))
	l.Warning().Msg("I can log in JSON format")
}

func simpleTextLogger() {
	l := log.New(log.WithName("text"), log.WithWriter(os.Stdout, log.FormatText))
	l.Warning().Msg("I can log in simple text format")
}

func nonBlocking() {
	l := log.New(log.WithName("nonblck"), log.WithMode(log.ModeNonBlocking))
	defer l.Close()
	l.Info().Msg("I'm working in Non Blocking mode, so I need Sync() before application exits")
	l.Info().Msg("All log lines are buffered")
	l.Info().Msg("And writted one by one")
	l.Info().Msg("1")
	l.Info().Msg("2")
	l.Info().Msg("3")
	l.Info().Msg("4")
	l.Info().Msg("5")
	l.Info().Msg("6")
	l.Info().Msg("7")
	l.Info().Msg("8")
	l.Info().Msg("9")
	l.Info().Msg("10")
	time.Sleep(100 * time.Millisecond)
}

func propsLogger() {
	l := log.New(log.WithName("props"), log.WithWriter(os.Stdout, log.FormatTextColor))
	l.Warning().Prop("prop1", "hello").Prop("prop2", "world").Msg("I can log in simple text format")
}

func propsJsonLogger() {
	l := log.New(log.WithName("propsJson"), log.WithWriter(os.Stdout, log.FormatJson))
	l.Warning().Prop("prop1", "hello").Prop("prop2", "world").Msg("I can log in simple json format")
}
