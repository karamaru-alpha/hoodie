package enum

import (
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/core"
	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/generator/enum/input"
	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/generator/enum/output/pkg/domain/enum"
)

type generator struct {
	core.GeneratorBase
	plugin *protogen.Plugin
}

func NewGenerator(plugin *protogen.Plugin) core.Generator {
	return &generator{
		GeneratorBase: core.NewGeneratorBase(plugin),
		plugin:        plugin,
	}
}

func (g *generator) Build() ([]core.GenFile, error) {
	enums, _, err := core.ConvertMessageFromProto(g.plugin.Files, nil, input.ConvertMessageFromProto)
	if err != nil {
		return nil, err
	}

	tpl := core.GetBaseTemplate("generator/enum")

	eachCreators := []core.EachTemplateCreator[*input.Enum]{
		enum.NewCreator(tpl),
	}

	return core.CreateEachTemplate(enums, eachCreators)
}
