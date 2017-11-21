package main

import (
	"time"

	"github.com/JonShard/CloudComputingAssignment2/cmd/api_hooks"
)

func main() {
	for {
		delay := time.Minute * 15
		api_hooks.StoreCurrencies("EUR")
		api_hooks.InvokeAllHooks(false)

		time.Sleep(delay)
	}
}
