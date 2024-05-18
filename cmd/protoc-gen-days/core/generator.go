package core

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/parallel"
)

type Generator interface {
	GeneratorBase
	Build() ([]GenFile, error)
}

type GeneratorBase interface {
	SetGenFiles(genFiles []GenFile)
	Format(ctx context.Context) error
	Generate(ctx context.Context) error
}

type generatorBase struct {
	plugin               *protogen.Plugin
	directGenerationMode bool
	baseOutputDir        string
	genFiles             []GenFile
}

func NewGeneratorBase(plugin *protogen.Plugin) GeneratorBase {
	params := strings.Split(plugin.Request.GetParameter(), ",")
	var directGenerationMode bool
	var baseOutputDir string
	for _, param := range params {
		s := strings.Split(param, "=")
		switch s[0] {
		case "direct_generation_mode":
			directGenerationMode = true
		case "output_dir":
			if len(s) == 2 {
				baseOutputDir = s[1]
			}
		}
	}
	return &generatorBase{
		plugin:               plugin,
		directGenerationMode: directGenerationMode,
		baseOutputDir:        baseOutputDir,
		genFiles:             make([]GenFile, 0),
	}
}

func (g *generatorBase) SetGenFiles(genFiles []GenFile) {
	g.genFiles = genFiles
}

func (g *generatorBase) Format(ctx context.Context) error {
	pg, ctx := parallel.NewGroupWithContext(ctx, parallel.DefaultSize)

	for _, file := range g.genFiles {
		if !strings.HasSuffix(file.GetFilePath(), ".gen.go") {
			continue
		}

		pg.Go(ctx, func(_ context.Context) error {
			return file.Format()
		})
	}

	return pg.Wait()
}

func (g *generatorBase) Generate(ctx context.Context) error {
	// NOTE: in direct generation mode, files are written directly. This method is faster, but it is not the original method of protoc-gen.
	if g.directGenerationMode {
		pg, ctx := parallel.NewGroupWithContext(ctx, parallel.DefaultSize)

		for _, file := range g.genFiles {
			pg.Go(ctx, func(_ context.Context) error {
				outputDir := filepath.Join(g.baseOutputDir, filepath.Dir(file.GetFilePath()))
				if err := os.MkdirAll(outputDir, 0777); err != nil {
					return err
				}

				return file.CreateOrWrite(g.baseOutputDir)
			})
		}

		return pg.Wait()
	}

	// NOTE: the original method of protoc-gen. Outputs the file contents written to stdout."
	for _, file := range g.genFiles {
		gf := g.plugin.NewGeneratedFile(file.GetFilePath(), "")
		if _, err := gf.Write(file.GetData()); err != nil {
			return err
		}
	}
	return nil
}
