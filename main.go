package main

import (
	"context"
	"fmt"
	"strconv"

	"go.uber.org/fx"
)

type Config struct {
	Id int
}

type SampleA struct {
	Id string
}

type SampleB struct {
	Id string
}

func NewB(s string) SampleB {
	return SampleB{Id: s}
}

func NewAFromConfig(c Config) SampleA {
	return SampleA{Id: strconv.Itoa(c.Id)}
}

func CallSampleA(i SampleA) {
	fmt.Println(i.Id)
}

func CallSampleB(i SampleB) {
	fmt.Println(i.Id)
}

func NewConfig(lc fx.Lifecycle) Config {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			fmt.Println("start hook")
			return nil
		},
		OnStop: func(context.Context) error {
			fmt.Println("stop hook")
			return nil
		},
	})
	return Config{Id: 1}
}

func main() {
	fx.New(
		fx.Supply("example B"), // Provides bare string which will be used for B constructor
		fx.Provide(
			NewB,
			NewConfig,
			NewAFromConfig,
		),
		fx.Invoke(
			CallSampleA,
			CallSampleB,
		),
		Submodule(),
		SecondSubmodule(),
	).Run()
}

func Submodule() fx.Option {
	return fx.Module(
		"sampleMod",
		fx.Decorate(func() SampleA {
			return SampleA{Id: "2"} // Just updates the dependency
		}),
		fx.Invoke(CallSampleA), // Gets it from parent module, and decorates
		SubSubmodule(),         // Gets decorated dependencies
	)
}

func SubSubmodule() fx.Option {
	return fx.Module(
		"sampleMod",
		fx.Decorate(func(s SampleA) SampleA {
			return SampleA{Id: "submodule_call: " + s.Id}
		}),
		fx.Invoke(CallSampleA),
	)
}

func SecondSubmodule() fx.Option {
	return fx.Module(
		"sampleMod2",
		fx.Invoke(CallSampleA), // Gets it from parent module
		SubSubmodule(),
	)
}
