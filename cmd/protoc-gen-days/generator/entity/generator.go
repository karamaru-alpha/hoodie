package entity

import (
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/core"
	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/generator/entity/input"
	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/generator/entity/output/pkg/domain/entity"
)

type generator struct {
	core.GeneratorBase
	plugin      *protogen.Plugin
	flagKindSet core.FlagKindSet
}

func NewGenerator(plugin *protogen.Plugin, flagKindSet core.FlagKindSet) core.Generator {
	return &generator{
		GeneratorBase: core.NewGeneratorBase(plugin),
		plugin:        plugin,
		flagKindSet:   flagKindSet,
	}
}

func (g *generator) Build() ([]core.GenFile, error) {
	messages, _, err := core.ConvertMessageFromProto(g.plugin.Files, g.flagKindSet, input.ConvertMessageFromProto)
	if err != nil {
		return nil, err
	}
	tpl := core.GetBaseTemplate("generator/entity")
	eachCreators := []core.EachTemplateCreator[*input.Message]{
		entity.NewTypesCreator(tpl),
	}
	genFiles, err := core.CreateEachTemplate(messages, eachCreators)
	if err != nil {
		return nil, err
	}

	return genFiles, nil
}
