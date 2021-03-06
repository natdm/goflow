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

	"io"

	log "github.com/Sirupsen/logrus"
)

// Parse is responsible for handling all the logic for parsing
// go models
type Parse struct {
	sync.Mutex
	Files     []string
	recursive bool

	// Mappings is a one-per-type map of each type
	mappings     map[string][]field
	baseMappings map[string]field
	comments     map[string]string
	embeds       map[string][]string
	outfile      io.Writer
}

// New returns a new parser
func New(r bool, f *os.File) *Parse {
	return &Parse{
		comments:     make(map[string]string),
		mappings:     make(map[string][]field),
		embeds:       make(map[string][]string),
		baseMappings: make(map[string]field),
		Files:        []string{},
		recursive:    r,
		outfile:      f,
	}
}

// Tag represents the Go struct tags. The original tags, the JSON specific tags, and the GoFlow (parse) tags.
// Parse tags have priority over the JSON tags
type tag struct {
	original string
	json     string
	flow     flowTag
}

type flowTag struct {
	name string
	typ  string
}

func (t *flowTag) String() string {
	return fmt.Sprintf(`name: "%s"    type: "%s"`, t.name, t.typ)
}

// Field represents an un-ignored, exportable, json-tagged field within the Go type
type field struct {
	// Name of the field
	name string

	// Comment is the inline comment in Go that will be carried over to
	// flow
	comment string

	// Tags represents all the Go json tags for which the fields
	// in the Flow type will be named
	tags tag

	// Go type, printed as string
	typ string

	// Children are used to gain access to nested (not embedded) structs.
	// Not currently supported by flow, but keeping the logic in case it is soon.
	children []field
}

// ParseDir parses a directory for all go files
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

// ParseFiles parses all files in p.Files to get all go types
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
			p.Lock()
			for _, structName := range structKeys {
				p.mappings[structName] = p.parseStruct(structMap[structName], structName, bs)
			}
			p.Unlock()
			baseKeys := make([]string, 0, len(baseMap))
			for k := range baseMap {
				baseKeys = append(baseKeys, k)
			}
			sort.Strings(baseKeys)
			p.Lock()
			for _, baseName := range baseKeys {
				p.baseMappings[baseName] = field{
					typ:  baseMap[baseName],
					name: baseName,
				}
			}
			p.Unlock()
		}(fname)
	}
	wg.Wait()
	return nil
}

func (p *Parse) parseTypes(f *ast.File) (map[string]*ast.FieldList, map[string]string) {
	structMap := map[string]*ast.FieldList{}
	baseMap := make(map[string]string)
	// range over the structs and fill struct map
	for _, d := range f.Scope.Objects {
		if f.Comments != nil {
			p.Lock()
			for _, v := range f.Comments {
				c := v.Text()
				p.comments[firstWord(c)] = c
			}
			p.Unlock()
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

func (p *Parse) parseStruct(fs *ast.FieldList, name string, bs []byte) []field {
	out := []field{}
	if fs == nil {
		return out
	}

	for _, f := range fs.List {
		newField := field{}

		if len(f.Names) == 0 {
			// // Need logic to save types previously declared and declare the fields in those types in here.
			switch f.Type.(type) {
			case *ast.Ident:
				// Treat as embedded struct
				t := string(bs[f.Type.Pos()-1 : f.Type.End()-1])
				// p.Embeds[name] = append(p.Embeds[name], t)
				if f.Comment != nil {
					newField.comment = f.Comment.Text()
				}
				newField.typ = "embedded"
				newField.tags = tag{
					flow: flowTag{
						name: t,
					},
				}
			default:
				log.Printf("unknown type, %s\n", f.Type)
			}
			out = append(out, newField)
			continue
		}

		// Parse tags after looking for embedded types.
		// Ignore any json-ignore tags
		if f.Tag == nil || strings.Contains(f.Tag.Value, "json:\"-\"") {
			continue
		}

		// If there are JSON tags,  parse it.
		if strings.Contains(f.Tag.Value, "json:") {
			newField.name = f.Names[0].String()
			newField.tags.original = f.Tag.Value
			if f.Comment != nil {
				newField.comment = f.Comment.Text()
			}
			switch f.Type.(type) {
			case *ast.InterfaceType:
				newField.typ = "Object"
			case *ast.MapType:
				x := f.Type.(*ast.MapType)
				newField.typ = parseMap(x)
			case *ast.ArrayType:
				x := f.Type.(*ast.ArrayType)
				newField.typ = parseArray(x)
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
				newField.typ = "Object"
			case *ast.StarExpr:
				x := f.Type.(*ast.StarExpr)
				t := string(bs[x.Pos()-2 : x.End()-1])
				t = strings.Replace(t, "*", "?", -1)
				t = strings.Replace(t, "time.Duration", "string", -1)
				t = strings.Replace(t, "time.Time", "string", -1)

				if strings.Contains(t, "[]") {
					newField.typ = strings.Replace(t, "[]", "Array<", -1) + ">"
				} else {
					newField.typ = t
				}
			default:
				t := string(bs[f.Type.Pos()-1 : f.Type.End()-1])
				if strings.Contains(t, "time") {
					newField.typ = "string"
				} else {
					newField.typ = t
				}
			}
			out = append(out, newField)
		}
	}
	return out
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
	return fmt.Sprintf("{ [key: %s]: %s }", inspectNode(ts.Key), inspectNode(ts.Value))
}

// inspectNode checks what is determined to be the value of a node based on a type-assertion of the node.
func inspectNode(node ast.Node) string {
	var out string
	ast.Inspect(node, func(n ast.Node) bool {
		var s string
		switch y := n.(type) {
		case *ast.BasicLit:
			s = y.Value
		case *ast.Ident:
			s = y.Name
		case *ast.StarExpr:
			out = "?"
		}

		if s != "" {
			if s == "error" {
				s = "string"
			}
			out += s
		}
		return true
	})
	return out
}

func removeDuplicates(s []string) []string {
	found := make(map[string]bool)
	j := 0
	for i, x := range s {
		if !found[x] {
			found[x] = true
			(s)[j] = (s)[i]
			j++
		}
	}
	s = (s)[:j]
	cp := []string{}
	for _, v := range s {
		cp = append(cp, v)
	}
	return cp
}

// updateTags updates tags in place. If a comma is before an ending quote, it stops at the comma
func updateTags(f []field) {
	for i := range f {
		if strings.Contains(f[i].comment, "a duration") {
			log.Println(f[i])
		}
		// If the type is not exported, ignore the type and all fields
		// Set it to blank to ignore later
		if !isExported(f[i].name) {
			f = []field{}
			return
		}
		flowTags := parseFlowTag(getTag("flow", f[i].tags.original))
		f[i].tags.json = getTag("json", f[i].tags.original)
		f[i].tags.flow = flowTags
		if flowTags.typ != "" {
			f[i].typ = flowTags.typ
		} else if f[i].typ == "struct" {
			updateTags(f[i].children)
		}
	}
}

func parseFlowTag(tag string) flowTag {
	sp := strings.Split(tag, ".")
	switch len(sp) {
	case 2:
		return flowTag{
			name: sp[0],
			typ:  sp[1],
		}
	case 1:
		return flowTag{
			name: sp[0],
		}
	default:
		return flowTag{}
	}
}

func getTag(tag string, tags string) string {
	loc := strings.Index(tags, fmt.Sprintf("%s:\"", tag))
	if loc > -1 {
		bs := []byte(tags)
		bs = bs[loc+len(tag)+2:]
		loc = strings.Index(string(bs), "\"")
		commaLoc := strings.Index(string(bs), ",")
		if commaLoc > -1 && commaLoc < loc {
			return string(bs[:commaLoc])
		}
		if loc == -1 {
			return ""
		}
		return string(bs[:loc])
	}
	return ""
}

func updateTypes(name string, fields []field) {
	for i := range fields {
		fields[i].typ = updateType(fields[i].typ)
		if len(fields[i].children) > 0 {
			updateTypes(name, fields[i].children)
		}
	}
}

func updateType(t string) string {
	// conversions is used in strings.Replacer to replace each key/val pair. Left is Go, right is Flow.
	conversions := []string{
		"*", "?",
		"int64", "number",
		"int32", "number",
		"int16", "number",
		"int8", "number",
		"int", "number",
		"uint64", "number",
		"uint32", "number",
		"uint16", "number",
		"uint8", "number",
		"uint", "number",
		"uintptr", "number",
		"byte", "number",
		"rune", "number",
		"float32", "number",
		"float64", "number",
		"complex64", "number",
		"complex128", "number",
		"bool", "boolean"}

	replacer := strings.NewReplacer(conversions...)
	return replacer.Replace(t)
}

// isExported returns true if the first character in a string is already capital
func isExported(s string) bool {
	if len(s) == 0 {
		return false
	}
	l := string(byte(s[0]))
	return strings.ToUpper(l) == l
}

// RemoveUnexported removes anything that is not exported from the map used to reference types before writing
func removeUnexported(m map[string][]field) {
	for k := range m {
		if !isExported(k) {
			delete(m, k)
		}
	}
}
