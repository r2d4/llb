package main

import (
	"context"
	"log"

	"github.com/moby/buildkit/frontend/gateway/grpcclient"
	"github.com/r2d4/llb/pkg/build"
)

func main() {
	if err := grpcclient.RunFromEnvironment(context.Background(), build.BuildFunc); err != nil {
		log.Fatal(err)
	}
}
