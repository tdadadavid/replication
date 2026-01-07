package cmd

import (
	"dbreplication/internal/api"
)

func Execute() {
	// register otel spans and other neccessary things for observability

	api.Start()
}
