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
	).Run()
}

func Submodule() fx.Option {
	return fx.Module(
		"sampleMod",
		fx.Invoke(CallSampleA), // Gets it from parent module
	)
}
