package main

import (
	"go.uber.org/fx"

	"github.com/karamaru-alpha/hoodie/pkg/cmd/api/di"
)

func main() {
	fx.New(di.Initialize()).Run()
}
