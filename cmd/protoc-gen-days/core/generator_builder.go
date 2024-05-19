package core

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/parallel"
)

type GeneratorBuilder interface {
	AppendGenerator(generator Generator) GeneratorBuilder
	Generate(ctx context.Context) error
}

type generatorBuilder struct {
	generators []Generator
}

func NewGeneratorBuilder() GeneratorBuilder {
	return &generatorBuilder{
		generators: make([]Generator, 0),
	}
}

func (g *generatorBuilder) AppendGenerator(generator Generator) GeneratorBuilder {
	g.generators = append(g.generators, generator)

	return g
}

func (g *generatorBuilder) Generate(ctx context.Context) error {
	pg, ctx := parallel.NewGroupWithContext(ctx, parallel.DefaultSize)

	for _, generator := range g.generators {
		startTime := time.Now()
		slog.InfoContext(ctx, fmt.Sprintf("%T build start", generator))

		pg.Go(ctx, func(ctx context.Context) error {
			genFiles, err := generator.Build()
			if err != nil {
				return err
			}
			generator.SetGenFiles(genFiles)

			slog.InfoContext(ctx, fmt.Sprintf(" %T build end, elapsed: %s", generator, time.Since(startTime).String()))
			if err := generator.Format(ctx); err != nil {
				return err
			}
			slog.InfoContext(ctx, fmt.Sprintf(" %T format end, elapsed: %s", generator, time.Since(startTime).String()))

			if err := generator.Generate(ctx); err != nil {
				return err
			}
			slog.InfoContext(ctx, fmt.Sprintf(" %T generate end, elapsed: %s", generator, time.Since(startTime).String()))

			return nil
		})
	}
	return pg.Wait()
}
