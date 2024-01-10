package main

import (
	"context"
	"log"

	"go.uber.org/zap"

	"pokemon-rest-api/pkg/pokemons"
)

func main() {
	ctx := context.Background()
	logger := zap.Must(zap.NewDevelopment())

	if err := pokemons.Exec(ctx, logger); err != nil {
		log.Fatalf("ERROR: %v", err)
	}
}
