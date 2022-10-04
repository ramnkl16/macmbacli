package helper

import (
	"fmt"
	"strings"
	"time"
	"unicode"

	"text/template"

	"github.com/gedex/inflector"
	"github.com/google/uuid"
	"github.com/iancoleman/strcase"
	"github.com/kenshaw/snaker"
)

var (
	IgnoreNames      map[string]bool
	ShortNameTypeMap map[string]string
)

func NewTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"removeidsuffix":      removeIdSuffix,
		"formatName":          formatName,
		"tolowercase":         strings.ToLower,
		"touppercase":         strings.ToUpper,
		"pluralize":           inflector.Pluralize,
		"timenow":             func() string { return time.Now().Local().String() },
		"tolowercamel":        strcase.ToLowerCamel,
		"colvals":             colvals,
		"sqlColnames":         sqlColnames,
		"fieldnames":          fieldnames,
		"fieldnamesmulti":     fieldnamesmulti,
		"GetLocalFields":      GetLocalFields,
		"shortname":           shortname,
		"fieldnamesforUpdate": fieldnamesforUpdate,
		"NewGUID":             uuid.New,
		"getField":            getField,
		"getFilteredFields":   getFilteredFields,
	}
}

func getField(fields []*Field, fieldName string) *Field {
	for _, f := range fields {
		//fmt.Println("getfield:" + f.ColumnName)
		if strings.EqualFold(f.Name, fieldName) {
			fmt.Println("return getfield:", f.Name)
			return f
		}
	}
	return nil
}

func getFilteredFields(fields []*Field, ignores ...string) []*Field {
	additionalIgnore := map[string]bool{}
	for _, fn := range ignores {
		additionalIgnore[fn] = true
	}
	list := make([]*Field, 0)

	for _, f := range fields {
		if IgnoreNames[f.DbFieldName] || (additionalIgnore[f.DbFieldName]) {
			continue
		} else {
			list = append(list, f)
		}
	}
	return list
}

// fieldnames creates a list of field names from fields of the adding the
// provided prefix, and excluding any Field with Name contained in ignoreNames.
//
// Used to present a comma separated list of field names, ie in a Go statement
// (ie, "t.Field1, t.Field2, t.Field3 ...")
func fieldnamesforUpdate(fields []*Field, ignores ...string) string {

	additionalIgnore := map[string]bool{}
	for _, n := range ignores {
		additionalIgnore[n] = true
	}
	str := ""
	i := 0
	for _, f := range fields {
		if IgnoreNames[f.DbFieldName] || (additionalIgnore[f.DbFieldName]) {
			continue
		}
		//assuming id is primay for all table
		if f.DbFieldName == "id" {
			continue
		}

		if i != 0 {
			str = str + ", "
		}
		str = str + f.DbFieldName + " = ?"
		i++
	}

	return str
}

// fieldnames creates a list of field names from fields of the adding the
// provided prefix, and excluding any Field with Name contained in ignoreNames.
//
// Used to present a comma separated list of field names, ie in a Go statement
// (ie, "t.Field1, t.Field2, t.Field3 ...")
func fieldnames(fields []*Field, prefix string, ignores ...string) string {

	additionalIgnore := map[string]bool{}
	for _, n := range ignores {
		additionalIgnore[strings.ToLower(n)] = true
	}
	str := ""
	i := 0
	for _, f := range fields {
		if IgnoreNames[strings.ToLower(f.Name)] || (additionalIgnore[strings.ToLower(f.Name)]) {
			continue
		}

		if i != 0 {
			str = str + ", "
		}
		str = str + prefix + "." + snaker.ForceCamelIdentifier(f.Name)
		i++
	}

	return str
}
func fieldnamesmulti(fields []*Field, prefix string) string {

	str := ""
	i := 0
	for _, f := range fields {
		// if ignoreNames[f.LocalDbFieldName] {
		// 	continue
		// }

		if i != 0 {
			str = str + ", "
		}
		str = str + prefix + "." + formatName(f.DbFieldName)
		i++
	}
	return str
}

func sqlColnames(fields []*Field, ignores ...string) string {
	additionalIgnore := map[string]bool{}
	for _, n := range ignores {
		additionalIgnore[strings.ToLower(n)] = true
		fmt.Println(n)
	}
	str := ""
	i := 0
	for _, f := range fields {
		if additionalIgnore[strings.ToLower(f.DbFieldName)] || len(strings.TrimSpace(strings.ToLower(f.DbFieldName))) == 0 {

			continue
		}
		fmt.Println("local", f.DbFieldName)
		if i != 0 {
			str = str + ", "
		}
		str = str + (f.DbFieldName)
		i++
	}

	return str
}
func GetLocalFields(fields []*Field) []*Field {
	filteredField := make([]*Field, 0)
	for _, f := range fields {
		if len(strings.TrimSpace(f.DbFieldName)) == 0 {
			continue
		}
		filteredField = append(filteredField, f)
	}
	return filteredField
}

// var (
// 	formatName  = template.FuncMap{"formatName": formatName}
// 	colValsFunc = template.FuncMap{"colvalues": colvals}
// 	colnames    = template.FuncMap{"colnames": columns}
func colvals(fields []*Field, ignores ...string) string {
	additionalIgnore := map[string]bool{}
	for _, n := range ignores {
		additionalIgnore[n] = true
	}
	str := ""
	i := 0
	for _, f := range fields {
		if additionalIgnore[f.DbFieldName] || len(strings.TrimSpace(f.DbFieldName)) == 0 {
			continue
		}
		if i != 0 {
			str = str + ", "
		}
		str = str + "?"
		i++
	}
	return str
}
func formatName(name string) string {
	newName := lintName(strings.Title(name))
	// If a first charactor of the table is number, add "A" to the top
	if unicode.IsNumber(rune(newName[0])) {
		newName = "A" + newName
	}

	return newName
}

func removeIdSuffix(name string) string {
	if strings.HasSuffix(name, "Id") || strings.HasSuffix(name, "id") || strings.HasSuffix(name, "ID") {
		v := strcase.ToCamel(name)
		return string(v[0 : len(v)-2])
	}
	return name
}

func lintName(name string) (should string) {
	// Fast path for simple cases: "_" and all lowercase.
	if name == "_" {
		return name
	}
	allLower := true
	for _, r := range name {
		if !unicode.IsLower(r) {
			allLower = false
			break
		}
	}
	if allLower {
		return name
	}
	// Split camelCase at any lower->upper transition, and split on underscores.
	// Check each word for common initialisms.
	runes := []rune(name)
	w, i := 0, 0 // index of start of word, scan
	for i+1 <= len(runes) {
		eow := false // whether we hit the end of a word
		if i+1 == len(runes) {
			eow = true
		} else if runes[i+1] == '_' {
			// underscore; shift the remainder forward over any run of underscores
			eow = true
			n := 1
			for i+n+1 < len(runes) && runes[i+n+1] == '_' {
				n++
			}

			// Leave at most one underscore if the underscore is between two digits
			if i+n+1 < len(runes) && unicode.IsDigit(runes[i]) && unicode.IsDigit(runes[i+n+1]) {
				n--
			}

			copy(runes[i+1:], runes[i+n+1:])
			runes = runes[:len(runes)-n]
		} else if unicode.IsLower(runes[i]) && !unicode.IsLower(runes[i+1]) {
			// lower->non-lower
			eow = true
		}
		i++
		if !eow {
			continue
		}

		// [w,i) is a word.
		word := string(runes[w:i])
		if u := strings.ToUpper(word); commonInitialisms[u] {
			// Keep consistent case, which is lowercase only at the start.
			if w == 0 && unicode.IsLower(runes[w]) {
				u = strings.ToLower(u)
			}
			// All the common initialisms are ASCII,
			// so we can replace the bytes exactly.
			copy(runes[w:], []rune(u))
		} else if w > 0 && strings.ToLower(word) == word {
			// already all lowercase, and not the first word, so uppercase the first character.
			runes[w] = unicode.ToUpper(runes[w])
		}
		w = i
	}
	return string(runes)
}
func shortname(typ string) string {

	runes := []rune(strings.Trim(typ, " "))
	return strings.ToLower(string(runes[0:2]))
}

// commonInitialisms is a set of common initialisms.
// Only add entries that are highly unlikely to be non-initialisms.
// For instance, "ID" is fine (Freudian code is rare), but "AND" is not.
var commonInitialisms = map[string]bool{
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SQL":   true,
	"SSH":   true,
	"TCP":   true,
	"TLS":   true,
	"TTL":   true,
	"UDP":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XSRF":  true,
	"XSS":   true,
}

// goReservedNames is a map of of go reserved names to "safe" names.
var goReservedNames = map[string]string{
	"break":       "brk",
	"case":        "cs",
	"chan":        "chn",
	"const":       "cnst",
	"continue":    "cnt",
	"default":     "def",
	"defer":       "dfr",
	"else":        "els",
	"fallthrough": "flthrough",
	"for":         "fr",
	"func":        "fn",
	"go":          "goVal",
	"goto":        "gt",
	"if":          "ifVal",
	"import":      "imp",
	"interface":   "iface",
	"map":         "mp",
	"package":     "pkg",
	"range":       "rnge",
	"return":      "ret",
	"select":      "slct",
	"struct":      "strct",
	"switch":      "swtch",
	"type":        "typ",
	"var":         "vr",

	// go types
	"error":      "e",
	"bool":       "b",
	"string":     "str",
	"byte":       "byt",
	"rune":       "r",
	"uintptr":    "uptr",
	"int":        "i",
	"int8":       "i8",
	"int16":      "i16",
	"int32":      "i32",
	"int64":      "i64",
	"uint":       "u",
	"uint8":      "u8",
	"uint16":     "u16",
	"uint32":     "u32",
	"uint64":     "u64",
	"float32":    "z",
	"float64":    "f",
	"complex64":  "c",
	"complex128": "c128",
}
