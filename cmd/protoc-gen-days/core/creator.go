package core

type TemplateInfo struct {
	Data     []byte
	FilePath string
}

type EachTemplateCreator[T any] interface {
	Create(message T) (*TemplateInfo, error)
}

type BulkTemplateCreator[T any] interface {
	Create(pkgName string, messages []T) (*TemplateInfo, error)
}

func CreateEachTemplate[T any](messages []T, creators []EachTemplateCreator[T]) ([]GenFile, error) {
	genFiles := make([]GenFile, 0, len(messages)*len(creators))
	for _, msg := range messages {
		for _, creator := range creators {
			info, err := creator.Create(msg)
			if err != nil {
				return nil, err
			}
			if info == nil {
				continue
			}
			genFiles = append(genFiles, NewGenFile(info.FilePath, info.Data))
		}
	}
	return genFiles, nil
}
