package main

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/core"
	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/generator/entity"
	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/generator/enum"
)

func main() {
	locationName := "Asia/Tokyo"
	location, err := time.LoadLocation(locationName)
	if err != nil {
		location = time.FixedZone(locationName, 9*60*60)
	}
	time.Local = location

	startTime := time.Now()
	slog.Info("protoc-gen-days start")

	generatorBuilder := core.NewGeneratorBuilder()

	protogen.Options{}.Run(func(plugin *protogen.Plugin) error {
		generatorMap := createGeneratorMap(plugin)
		kinds := make([]string, 0, len(generatorMap))
		for kind, generator := range generatorMap {
			kinds = append(kinds, string(kind))
			generatorBuilder.AppendGenerator(generator)
		}
		slog.Info(fmt.Sprintf("flag: %s", strings.Join(kinds, ",")))

		return generatorBuilder.Generate(context.Background())
	})

	endTime := time.Now()
	slog.Info(fmt.Sprintf("protoc-gen-days end, elapsed: %s", endTime.Sub(startTime).String()))
}

func createGeneratorMap(plugin *protogen.Plugin) map[core.FlagKind]core.Generator {
	params := strings.Split(plugin.Request.GetParameter(), ",")
	flagKindSet := make(core.FlagKindSet, len(params))

	for _, param := range params {
		switch core.FlagKind(strings.Split(param, "=")[0]) {
		case core.FlagKindGenEntity:
			flagKindSet.Add(core.FlagKindGenEntity)
		case core.FlagKindGenEnum:
			flagKindSet.Add(core.FlagKindGenEnum)
		default:
			// do nothing
		}
	}
	generatorMap := make(map[core.FlagKind]core.Generator, flagKindSet.Size())
	if flagKindSet.Has(core.FlagKindGenEntity) {
		generatorMap[core.FlagKindGenEntity] = entity.NewGenerator(plugin, flagKindSet)
	}
	if flagKindSet.Has(core.FlagKindGenEnum) {
		generatorMap[core.FlagKindGenEnum] = enum.NewGenerator(plugin)
	}

	return generatorMap
}
