package parse

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"sync"
)

// Parse is responsible for handling all the logic for parsing
// go models
type Parse struct {
	sync.Mutex
	Files     []string
	Out       *os.File
	recursive bool

	// Mappings is a one-per-type map of each type
	Mappings     map[string][]Field
	BaseMappings map[string]Field
	Comments     map[string]string
}

// New returns a new parser
func New(r bool) *Parse {
	return &Parse{
		Comments:     make(map[string]string),
		Mappings:     make(map[string][]Field),
		BaseMappings: make(map[string]Field),
		Files:        []string{},
		recursive:    r,
	}
}

// Tag represents the Go struct tags. The original tags, the JSON specific tags, and the GoFlow (parse) tags.
// Parse tags have priority over the JSON tags
type Tag struct {
	Original string
	JSON     string
	Flow     FlowTag
}

type FlowTag struct {
	Name string
	Type string
}

func (t *FlowTag) String() string {
	return fmt.Sprintf(`name: "%s"    type: "%s"`, t.Name, t.Type)
}

// Field represents an un-ignored, exportable, json-tagged field within the Go type
type Field struct {
	// Name of the field
	Name string

	// Comment is the inline comment in Go that will be carried over to
	// flow
	Comment string

	// JSONTags represents all the Go json tags for which the fields
	// in the Flow type will be named
	JSONTags Tag

	// Go type, printed as string
	Type string

	Children []Field
}

func (p *Parse) ParseDir(d string) (e error) {
	files, err := ioutil.ReadDir(d)
	if err != nil {
		return err
	}

	for _, v := range files {
		name := v.Name()
		if v.IsDir() {
			if p.recursive {
				if err := p.ParseDir(d + "/" + v.Name()); err != nil {
					e = err
					return
				}
			}
		} else if strings.HasSuffix(name, "go") && !strings.Contains(name, "_test") {
			p.Lock()
			p.Files = append(p.Files, strings.Replace(fmt.Sprintf("%s/%s", d, name), "//", "/", -1))
			p.Unlock()

		}
	}
	return nil
}

func (p *Parse) ParseFiles() (e error) {
	var wg sync.WaitGroup
	for _, fname := range p.Files {
		wg.Add(1)
		go func(fname string) {
			defer wg.Done()
			fset := token.NewFileSet() // positions are relative to fset

			// Parse the file given in arguments
			f, err := parser.ParseFile(fset, fname, nil, parser.ParseComments)
			if err != nil {
				e = err
				return
			}
			bs, err := ioutil.ReadFile(fname)
			if err != nil {
				e = err
				return
			}
			structMap, baseMap := p.parseTypes(f)
			// Parse structs
			structKeys := make([]string, 0, len(structMap))
			for k := range structMap {
				structKeys = append(structKeys, k)
			}
			sort.Strings(structKeys)
			for _, structName := range structKeys {
				p.Lock()
				p.Mappings[structName] = p.parseStruct(structMap[structName], bs)
				p.Unlock()
			}

			baseKeys := make([]string, 0, len(baseMap))
			for k := range baseMap {
				baseKeys = append(baseKeys, k)
			}
			sort.Strings(baseKeys)
			for _, baseName := range baseKeys {
				p.Lock()
				p.BaseMappings[baseName] = Field{
					Type: baseMap[baseName],
					Name: baseName,
				}
				p.Unlock()
			}
		}(fname)
	}
	wg.Wait()
	return nil
}

func firstWord(value string) string {
	for i := range value {
		if value[i] == ' ' {
			return value[0:i]
		}
	}
	return value
}

func parseArray(ts *ast.ArrayType) string {
	var arr string
	ast.Inspect(ts, func(n ast.Node) bool {
		var s string
		switch x := n.(type) {
		case *ast.BasicLit:
			s = "?" + x.Value
		case *ast.Ident:
			s = x.Name
		}

		if s != "" {
			if s == "error" {
				s = "string"
			}
			arr = fmt.Sprintf("Array<%s>", s)
		}
		return true
	})
	return arr
}

func parseMap(ts *ast.MapType) string {
	var key, val string
	ast.Inspect(ts.Key, func(n ast.Node) bool {
		var s string
		switch y := n.(type) {
		case *ast.BasicLit:
			s = y.Value
		case *ast.Ident:
			s = y.Name
		case *ast.StarExpr:
			key = "?"
		}

		if s != "" {
			if s == "error" {
				s = "string"
			}
			key += s
		}
		return true
	})
	ast.Inspect(ts.Value, func(n ast.Node) bool {
		var s string
		switch y := n.(type) {
		case *ast.BasicLit:
			s = y.Value
		case *ast.Ident:
			s = y.Name
		case *ast.StarExpr:
			val = "?"
		}

		if s != "" {
			if s == "error" {
				s = "string"
			}
			val += s
		}
		return true
	})
	return fmt.Sprintf("{ [key: %s]: %s }", key, val)
}

func (p *Parse) parseTypes(f *ast.File) (map[string]*ast.FieldList, map[string]string) {
	structMap := map[string]*ast.FieldList{}
	baseMap := make(map[string]string)
	// range over the structs and fill struct map
	for _, d := range f.Scope.Objects {
		if f.Comments != nil {
			for _, v := range f.Comments {
				c := v.Text()
				p.Lock()
				p.Comments[firstWord(c)] = c
				p.Unlock()
			}
		}
		ts, ok := d.Decl.(*ast.TypeSpec)
		if !ok {
			continue
		}

		switch ts.Type.(type) {
		case *ast.StructType:
			x, ok := ts.Type.(*ast.StructType)
			if !ok {
				continue
			}
			structMap[ts.Name.String()] = x.Fields
		case *ast.InterfaceType:
			continue
		case *ast.MapType:
			x, ok := ts.Type.(*ast.MapType)
			if !ok {
				continue
			}
			baseMap[d.Name] = parseMap(x)
		case *ast.ArrayType:
			x, ok := ts.Type.(*ast.ArrayType)
			if !ok {
				continue
			}
			baseMap[d.Name] = parseArray(x)
		default:
			baseMap[d.Name] = fmt.Sprintf("%v", ts.Type)
		}
	}
	return structMap, baseMap
}

func (p *Parse) parseStruct(fs *ast.FieldList, bs []byte) []Field {
	out := []Field{}
	for _, field := range fs.List {
		newField := Field{}
		if len(field.Names) == 0 {

			// Somewhere in here lies the answer to getting the embedded structs
			//
			embededStructName := string(bs[field.Type.Pos()-1 : field.Type.End()-1])
			_ = embededStructName
			// Need logic to save types previously declared and declare the fields in those types in here.
			continue
		}

		if field.Tag == nil || strings.Contains(field.Tag.Value, "json:\"-\"") {
			continue
		}

		if strings.Contains(field.Tag.Value, "json") {
			newField.Name = field.Names[0].String()
			newField.JSONTags.Original = field.Tag.Value
			if field.Comment != nil {
				newField.Comment = field.Comment.Text()
			}

			switch field.Type.(type) {
			case *ast.InterfaceType:
				newField.Type = "Object"
			case *ast.MapType:
				x, ok := field.Type.(*ast.MapType)
				if !ok {
					continue
				}
				newField.Type = parseMap(x)
			case *ast.ArrayType:
				x, ok := field.Type.(*ast.ArrayType)
				if !ok {
					continue
				}
				newField.Type = parseArray(x)
			case *ast.StructType:
				// commented out until I can find out how Flow would support nested objects. They may not.
				// x, ok := field.Type.(*ast.StructType)
				// if !ok {
				// 	continue
				// }
				// if x.Fields.List == nil || x.Fields.NumFields() == 0 {
				// 	newField.Type = "object"
				// 	continue
				// }
				// newField.Children = p.parseStruct(x.Fields, bs)
				newField.Type = "Object"
			case *ast.StarExpr:
				x, ok := field.Type.(*ast.StarExpr)
				if !ok {
					continue
				}
				t := string(bs[x.Pos()-2 : x.End()-1])
				t = strings.Replace(t, "*", "?", -1)
				t = strings.Replace(t, "time.Duration", "string", -1)
				t = strings.Replace(t, "time.Time", "string", -1)

				if strings.Contains(t, "[]") {
					newField.Type = strings.Replace(t, "[]", "Array<", -1) + ">"
				} else {
					newField.Type = t
				}
			default:
				t := string(bs[field.Type.Pos()-1 : field.Type.End()-1])
				if strings.Contains(t, "time") {
					newField.Type = "string"
				} else {
					newField.Type = t
				}
			}
		}
		out = append(out, newField)

	}
	return out
}
