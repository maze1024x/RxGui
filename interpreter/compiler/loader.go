package compiler

import (
	"fmt"
	"strings"
	"rxgui/standalone/util/richtext"
	"rxgui/lang/source"
	"rxgui/lang/textual/ast"
	"rxgui/lang/textual/syntax"
	"rxgui/lang/textual/parser"
	"rxgui/lang/textual/transformer"
)


const SourceFilenameSuffix = ".km"

type Loader struct {
	Suffix  string
	Load    func(bytes ([] byte), key source.FileKey) (*ast.Root, richtext.Error)
}
var loaders = [] Loader {{
	Suffix: SourceFilenameSuffix,
	Load:   func(content ([] byte), key source.FileKey) (*ast.Root, richtext.Error) {
		var code = source.DecodeUtf8ToCode(content)
		var name = key.String()
		var root = syntax.DefaultRoot()
		var cst_, err = parser.Parse(code, root, name)
		if err != nil { return nil, err }
		var ast_ = transformer.Transform(cst_, key).(ast.Root)
		return &ast_, nil
	},
}}
func Load(file string, fs FileSystem) (*ast.Root, richtext.Error) {
	var key = source.FileKey {
		Context: fs.Key(),
		Path:    file,
	}
	for _, l := range loaders {
		if strings.HasSuffix(file, l.Suffix) {
			var content, err = ReadFile(file, fs)
			if err != nil { return nil, richtext.ErrorFrom(err) }
			return l.Load(content, key)
		}
	}
	var err = fmt.Errorf("no loader available for file \"%s\"", key)
	return nil, richtext.ErrorFrom(err)
}


