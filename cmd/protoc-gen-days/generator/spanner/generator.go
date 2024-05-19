package spanner

import (
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/core"
	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/generator/spanner/input"
	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/generator/spanner/output/db/ddl"
	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/generator/spanner/output/pkg/di"
	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/generator/spanner/output/pkg/domain/repository"
	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/generator/spanner/output/pkg/infra/spanner"
	repositoryimpl "github.com/karamaru-alpha/days/cmd/protoc-gen-days/generator/spanner/output/pkg/infra/spanner/repository"
	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/generator/spanner/output/pkg/infra/spanner/repository/base"
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
	messages, pkgMap, err := core.ConvertMessageFromProto(g.plugin.Files, nil, input.ConvertMessageFromProto)
	if err != nil {
		return nil, err
	}

	tpl := core.GetBaseTemplate("generator/spanner")

	eachCreators := []core.EachTemplateCreator[*input.Message]{
		base.NewBufferCreator(tpl),
		base.NewCacheCreator(tpl),
		base.NewDefinitionCreator(tpl),
		base.NewQueryBuilderCreator(tpl),
		repository.NewCreator(tpl),
		repositoryimpl.NewCreator(tpl),
	}
	genFiles, err := core.CreateEachTemplate(messages, eachCreators)
	if err != nil {
		return nil, err
	}

	bulkCreators := []core.BulkTemplateCreator[*input.Message]{
		di.NewCreator(tpl),
		ddl.NewCreator(tpl),
		spanner.NewCreator(tpl),
	}
	tmpGenFiles, err := core.CreateBulkTemplate(pkgMap, bulkCreators)
	if err != nil {
		return nil, err
	}
	genFiles = append(genFiles, tmpGenFiles...)

	return genFiles, nil
}
