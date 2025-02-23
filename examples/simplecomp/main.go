package main

import (
	"log"

	"github.com/VanThen60hz/service-context"
)

func main() {
	const compId = "foo"

	serviceCtx := sctx.NewServiceContext(
		sctx.WithName("simple-component"),
		sctx.WithComponent(NewSimpleComponent(compId)),
	)

	if err := serviceCtx.Load(); err != nil {
		log.Fatal(err)
	}

	type CanGetValue interface {
		GetValue() string
	}

	comp := serviceCtx.MustGet(compId).(CanGetValue)

	log.Println(comp.GetValue())

	_ = serviceCtx.Stop()
}
