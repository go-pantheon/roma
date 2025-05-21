package pool

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/emicklei/proto"
	"github.com/go-pantheon/fabrica-util/camelcase"
	"github.com/go-pantheon/fabrica-util/errors"
	tmpl "github.com/go-pantheon/roma/vulcan/app/api/db/internal/template"
)

func ParseData(path string, goModuleBase string) (*tmpl.Data, error) {
	allMsgs := make(map[string]struct{})
	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(d.Name(), ".proto") {
			return nil
		}

		parseAllMessages(path, allMsgs)
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "error during second pass of proto scanning")
	}

	ret := &tmpl.Data{
		Files: make([]*tmpl.File, 0, len(allMsgs)),
	}

	err = filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(d.Name(), ".proto") {
			return nil
		}

		file, err := genFile(path, goModuleBase, allMsgs)
		if err != nil {
			return err
		}

		ret.Files = append(ret.Files, file)
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "error during second pass of proto scanning")
	}

	return ret, nil
}

func parseAllMessages(path string, messages map[string]struct{}) error {
	reader, err := os.Open(path)
	if err != nil {
		return errors.Wrapf(err, "error opening proto file:%s", path)
	}
	defer reader.Close()

	parser := proto.NewParser(reader)
	definition, err := parser.Parse()
	if err != nil {
		return errors.Wrapf(err, "error parsing proto file:%s", path)
	}

	for _, each := range definition.Elements {
		if msg, ok := each.(*proto.Message); ok {
			messages[msg.Name] = struct{}{}
		}
	}

	return nil
}

func genFile(path string, goModuleBase string, allMessages map[string]struct{}) (*tmpl.File, error) {
	ret := &tmpl.File{
		Dir:      filepath.Dir(path),
		FileName: strings.TrimSuffix(filepath.Base(path), ".proto"),
	}

	reader, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "error opening proto file:%s", path)
	}
	defer reader.Close()

	parser := proto.NewParser(reader)
	definition, err := parser.Parse()
	if err != nil {
		return nil, errors.Wrapf(err, "error parsing proto file:%s", path)
	}

	importsMap := make(map[string]*tmpl.Import)

	goPkgAlias, goPkg, _ := extractGoPackage(definition.Elements, goModuleBase)
	if goPkg != "" && goPkgAlias != "" {
		ret.Package = goPkgAlias
		importsMap[goPkg] = &tmpl.Import{Alias: goPkgAlias, Path: goPkg}
	}

	for _, impl := range importsMap {
		ret.Imports = append(ret.Imports, impl)
	}

	for _, each := range definition.Elements {
		if parsedMsg, ok := each.(*proto.Message); ok {
			msg := &tmpl.Message{
				Name:      parsedMsg.Name,
				GoPackage: goPkgAlias,
				Fields:    []*tmpl.Field{},
			}

			var fieldData *tmpl.Field

			for _, element := range parsedMsg.Elements {
				switch field := element.(type) {
				case *proto.MapField:
					fieldData = parseMapField(field, allMessages)
				case *proto.NormalField:
					if field.Repeated {
						fieldData = parseRepeatedField(field, allMessages)
					} else {
						fieldData = parseNormalField(field, allMessages)
					}
				default:
					continue
				}

				msg.Fields = append(msg.Fields, fieldData)
			}

			ret.Messages = append(ret.Messages, msg)
		}
	}

	return ret, nil
}

func parseMapField(protoField *proto.MapField, allDefinedMessageTypes map[string]struct{}) *tmpl.Field {
	ret := &tmpl.Field{
		Name:      camelcase.ToUnderScore(protoField.Name),
		IsMap:     true,
		KeyType:   protoField.KeyType,
		ValueType: protoField.Type,
	}

	if _, isMsg := allDefinedMessageTypes[protoField.Type]; isMsg {
		ret.ValueIsMessage = true
	}

	return ret
}

func parseRepeatedField(protoField *proto.NormalField, allDefinedMessageTypes map[string]struct{}) *tmpl.Field {
	ret := &tmpl.Field{
		Name:       camelcase.ToUnderScore(protoField.Name),
		IsRepeated: true,
		ValueType:  protoField.Type,
	}

	if _, isMsg := allDefinedMessageTypes[protoField.Type]; isMsg {
		ret.ValueIsMessage = true
	}

	return ret
}

func parseNormalField(protoField *proto.NormalField, allDefinedMessageTypes map[string]struct{}) *tmpl.Field {
	ret := &tmpl.Field{
		Name: camelcase.ToUnderScore(protoField.Name),
		Type: protoField.Type,
	}

	if _, isMsg := allDefinedMessageTypes[protoField.Type]; isMsg {
		ret.IsMessage = true
	}

	return ret
}

func extractGoPackage(elements []proto.Visitee, goModuleBase string) (alias, path string, err error) {
	for _, each := range elements {
		opt, ok := each.(*proto.Option)
		if !ok || opt.Name != "go_package" {
			continue
		}

		parts := strings.Split(opt.Constant.Source, ";")
		pkgPath := strings.Trim(parts[0], "\"")
		fullImportPath := pkgPath

		if goModuleBase != "" && !strings.HasPrefix(pkgPath, goModuleBase) {
			fullImportPath = filepath.Join(goModuleBase, pkgPath)
		}

		if len(parts) > 1 {
			alias = strings.TrimSpace(parts[1])
		} else {
			alias = filepath.Base(pkgPath)
		}

		return alias, fullImportPath, nil
	}

	return "", "", errors.New("go_package option not found")
}
