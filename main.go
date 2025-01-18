package main

import (
	"fmt"

	"go.uber.org/fx"
)

type SampleA struct {
	Id string
}

func NewA(s string) SampleA {
	return SampleA{Id: s}
}

func CallSampleA(i SampleA) {
	fmt.Println(i.Id)
}

func main() {
	fx.New(
		fx.Supply("example"),
		fx.Provide(NewA),
		fx.Invoke(CallSampleA),
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
