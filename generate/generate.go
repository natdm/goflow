package generate

import (
	"fmt"
	"io"

	"strings"

	"github.com/natdm/goflow/parse"
)

type Generate struct {
	out io.WriteCloser
}

var conversions = []string{
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

// UpdateTags updates tags in place. If a comma is before an ending quote, it stops at the comma
func UpdateTags(f []parse.Field) {
	for i := range f {

		// If the type is not exported, ignore the type and all fields
		// Set it to blank to ignore later
		if !IsExported(f[i].Name) {
			f = []parse.Field{}
			return
		}
		flowTags := parseFlowTag(getTag("flow", f[i].JSONTags.Original))

		f[i].JSONTags.JSON = getTag("json", f[i].JSONTags.Original)
		f[i].JSONTags.Flow = flowTags
		if flowTags.Type != "" {
			f[i].Type = flowTags.Type
		} else if f[i].Type == "struct" {
			UpdateTags(f[i].Children)
		}
	}
}

func parseFlowTag(tag string) parse.FlowTag {
	sp := strings.Split(tag, ".")
	switch len(sp) {
	case 2:
		return parse.FlowTag{
			Name: sp[0],
			Type: sp[1],
		}
	case 1:
		return parse.FlowTag{
			Name: sp[0],
		}
	default:
		return parse.FlowTag{}
	}
}

func UpdateTypes(name string, fields []parse.Field) {
	for i := range fields {
		fields[i].Type = UpdateType(fields[i].Type)
		if len(fields[i].Children) > 0 {
			UpdateTypes(name, fields[i].Children)
		}
	}
}

func UpdateType(t string) string {
	replacer := strings.NewReplacer(conversions...)
	return replacer.Replace(t)
}

// If the first character in a string is already capital, it is exported
func IsExported(s string) bool {
	if len(s) == 0 {
		return false
	}
	l := string(byte(s[0]))
	return strings.ToUpper(l) == l
}

func RemoveUnexported(m map[string][]parse.Field) {
	for k := range m {
		if !IsExported(k) {
			delete(m, k)
		}
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
