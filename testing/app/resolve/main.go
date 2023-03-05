package main

import (
	"context"
	"flag"
	"kantoku/testing/app/base"
	"log"
)

var dep = flag.String("dep", "", "the dependency to resolve")

func main() {
	flag.Parse()

	kan, err := base.Generate(context.Background())
	if err != nil {
		log.Fatal(err)
		return
	}

	err = kan.Depot().Deps().Resolve(context.Background(), *dep)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("Done, enjoy!")
}
