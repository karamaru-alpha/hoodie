package core

import (
	"fmt"
	"go/format"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"golang.org/x/tools/imports"

	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/perrors"
)

type GenFile interface {
	Format() error
	CreateOrWrite(baseOutputDir string) error
	GetFilePath() string
	GetData() []byte
}

type genFile struct {
	filePath     string
	newData      []byte
	oldDataCache map[string][]byte
	mu           *sync.Mutex
}

func NewGenFile(filePath string, newData []byte) GenFile {
	return &genFile{
		filePath:     filePath,
		newData:      newData,
		oldDataCache: make(map[string][]byte),
		mu:           &sync.Mutex{},
	}
}

func (g *genFile) Format() error {
	startTime := time.Now()
	importsData, err := imports.Process("", g.newData, &imports.Options{
		Fragment:   true,
		AllErrors:  false,
		Comments:   true,
		TabIndent:  true,
		TabWidth:   8,
		FormatOnly: false,
	})
	if err != nil {
		return perrors.Wrap(err, "goimports").SetValues(map[string]any{
			"filePath": g.GetFilePath(),
			"fileData": string(g.newData),
		})
	}
	since := time.Since(startTime)
	if since > 1*time.Second {
		// NOTE: output log if it takes more than 1 second
		slog.Warn(fmt.Sprintf(" %s goimports end", g.GetFilePath()), "elapsed", since.String())
	}

	fmtData, err := format.Source(importsData)
	if err != nil {
		return perrors.Wrap(err, "gofmt").SetValues(map[string]any{
			"filePath": g.GetFilePath(),
			"fileData": string(g.newData),
		})
	}
	since = time.Since(startTime)
	if since > 1*time.Second {
		// NOTE: output log if it takes more than 1 second
		slog.Warn(fmt.Sprintf(" %s gofmt end", g.GetFilePath()), "elapsed", since.String())
	}

	g.newData = fmtData

	return nil
}

func (g *genFile) CreateOrWrite(baseOutputDir string) error {
	path := filepath.Join(baseOutputDir, g.filePath)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			slog.Warn("fail to close the file", "err", err, "filePath", path)
		}
	}()

	if _, err := file.Write(g.newData); err != nil {
		return err
	}

	return nil
}

func (g *genFile) GetFilePath() string {
	return g.filePath
}

func (g *genFile) GetData() []byte {
	return g.newData
}
