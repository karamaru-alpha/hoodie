package core

import (
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/huandu/xstrings"
	"google.golang.org/protobuf/compiler/protogen"
)

const (
	module        = "github.com/karamaru-alpha/days"
	captureLength = 4
	uuidGoName    = "UUID"
)

var (
	idRegExp   = regexp.MustCompile(`^(.*)(Id|id)(s?|\d+|s\d+)$`)
	uuidRegExp = regexp.MustCompile(`^(.*)(Uuid|uuid)(s?|\d+|s\d+)$`)
)

type mapString struct {
	value map[string]string
	mutex *sync.RWMutex
}

func newMapString() *mapString {
	return &mapString{
		value: make(map[string]string, 10000),
		mutex: &sync.RWMutex{},
	}
}

func (m *mapString) Load(key string) (string, bool) {
	m.mutex.Lock()
	v, ok := m.value[key]
	m.mutex.Unlock()
	return v, ok
}

func (m *mapString) Store(key, value string) {
	m.mutex.Lock()
	m.value[key] = value
	m.mutex.Unlock()
}

var toSnakeCaseCache = newMapString()

// ToSnakeCase
//
//	userIDs -> user_ids
func ToSnakeCase(str string) string {
	if v, ok := toSnakeCaseCache.Load(str); ok {
		return v
	}

	// NOTE: ensure 'IDs' should not become 'i_ds'.
	result := xstrings.ToSnakeCase(strings.ReplaceAll(str, "IDs", "Ids"))
	toSnakeCaseCache.Store(str, result)
	return result
}

var toUpperSnakeCaseCache = newMapString()

// ToUpperSnakeCase
//
//	userIDs -> USER_IDS
func ToUpperSnakeCase(str string) string {
	if v, ok := toUpperSnakeCaseCache.Load(str); ok {
		return v
	}

	result := strings.ToUpper(ToSnakeCase(str))

	toUpperSnakeCaseCache.Store(str, result)
	return result
}

var toGolangCamelCaseCache = newMapString()

// ToGolangCamelCase
//
//	user_id -> userID
//	user_uuid -> userUUID
func ToGolangCamelCase(str string) string {
	if v, ok := toGolangCamelCaseCache.Load(str); ok {
		return v
	}

	result := ToGolangPascalCase(str)
	switch {
	case strings.HasPrefix(result, "ID"):
		result = strings.Replace(result, "ID", "id", 1)
	case result == uuidGoName:
		result = "uuid"
	default:
		result = xstrings.FirstRuneToLower(result)
	}

	toGolangCamelCaseCache.Store(str, result)
	return result
}

var toPascalCaseCache = newMapString()

// ToPascalCase
//
//	user_id -> UserId
func ToPascalCase(str string) string {
	if v, ok := toPascalCaseCache.Load(str); ok {
		return v
	}

	snakeStr := ToSnakeCase(str)
	result := xstrings.ToCamelCase(snakeStr)

	toPascalCaseCache.Store(str, result)
	return result
}

var toGolangPascalCaseCache = newMapString()

// ToGolangPascalCase
//
//	user_id -> UserID
func ToGolangPascalCase(str string) string {
	if v, ok := toGolangPascalCaseCache.Load(str); ok {
		return v
	}

	result := ToPascalCase(str)
	if captures := uuidRegExp.FindStringSubmatch(result); len(captures) == captureLength {
		result = captures[1] + uuidGoName + captures[3]
	} else if captures := idRegExp.FindStringSubmatch(result); len(captures) == captureLength {
		result = captures[1] + "ID" + captures[3]
	}

	toGolangPascalCaseCache.Store(str, result)
	return result
}

var toPkgNameCache = newMapString()

var toLocalNameCache = newMapString()

// ToLocalName
//
//	user_id -> userID
//	type -> typ
func ToLocalName(str string) string {
	if v, ok := toLocalNameCache.Load(str); ok {
		return v
	}

	result := ToGolangCamelCase(str)
	if result == "type" {
		result = "typ"
	}

	toLocalNameCache.Store(str, result)
	return result
}

// ToPkgName
//
//	pkg/domain/entity/admin -> admin
//	github.com/karamaru-alpha/days/pkg/domain/entity/admin -> admin
func ToPkgName(str string) string {
	if v, ok := toPkgNameCache.Load(str); ok {
		return v
	}

	fileDirs := strings.Split(GoPackageNameToFileDirName(str), "/")
	var result string
	if len(fileDirs) > 0 {
		result = fileDirs[len(fileDirs)-1]
	}

	toPkgNameCache.Store(str, result)
	return result
}

var goPackageNameToFileDirNameCache = newMapString()

// GoPackageNameToFileDirName
//
//	github.com/karamaru-alpha/days/pkg/domain/entity/admin -> pkg/domain/entity/admin
func GoPackageNameToFileDirName(str string) string {
	if v, ok := goPackageNameToFileDirNameCache.Load(str); ok {
		return v
	}
	return strings.TrimPrefix(str, module+"/")
}

type BaseInput interface {
	GetPkgName() string
	GetSnakeName() string
	Nil() bool
}

func ConvertMessageFromProto[T BaseInput](files []*protogen.File, flagKindSet FlagKindSet, converter func(*protogen.File, FlagKindSet) (T, error)) (messages []T, pkgMap map[string][]T, err error) {
	messages = make([]T, 0, len(files))
	pkgMap = make(map[string][]T, 0)
	for _, file := range files {
		if !file.Generate {
			continue
		}

		msg, err := converter(file, flagKindSet)
		if err != nil {
			return nil, nil, err
		}
		if msg.Nil() {
			continue
		}
		messages = append(messages, msg)
		pkgMap[msg.GetPkgName()] = append(pkgMap[msg.GetPkgName()], msg)
	}

	sort.SliceStable(messages, func(i, j int) bool {
		return messages[i].GetSnakeName() < messages[j].GetSnakeName()
	})
	for _, msgs := range pkgMap {
		sort.SliceStable(msgs, func(i, j int) bool {
			return msgs[i].GetSnakeName() < msgs[j].GetSnakeName()
		})
	}
	return messages, pkgMap, nil
}
