package main

import (
	"go.uber.org/fx"

	"github.com/karamaru-alpha/days/pkg/cmd/api/di"
)

func main() {
	fx.New(di.Initialize()).Run()
}
