package rpc

import (
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/core"
	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/generator/rpc/input"
	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/generator/rpc/output/pkg/cmd/di"
	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/generator/rpc/output/pkg/pb/rpc/connect"
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
	_, pkgMap, err := core.ConvertMessageFromProto(g.plugin.Files, nil, input.ConvertMessageFromProto)
	if err != nil {
		return nil, err
	}

	tpl := core.GetBaseTemplate("generator/rpc")

	bulkCreators := []core.BulkTemplateCreator[*input.Service]{
		di.NewCreator(tpl),
		connect.NewCreator(tpl),
	}
	genFiles, err := core.CreateBulkTemplate(pkgMap, bulkCreators)
	if err != nil {
		return nil, err
	}

	return genFiles, nil
}
