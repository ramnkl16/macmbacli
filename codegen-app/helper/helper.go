package helper

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
)

var (
	Tmpl                  *template.Template
	Wd                    string
	Entity                Module
	FlagDbUser            = flag.String("db_user", "root", "user name")
	FlagPassword          = flag.String("db_password", "Kish123", "db password")
	FlagDbName            = flag.String("db_name", "booknplay", "database name")
	FlagDbHost            = flag.String("db_host", "localhost", "host name  localhost")
	FlagDbPort            = flag.Int("db_port", 3306, "db port")
	FlagIgnore            = flag.String("ignore", "", "field ignore list")
	FlagEnableModule      bool
	FlagExtension         string
	FlagMysql             bool
	Db                    MySQL
	Config                Configuration
	ExcludeDtoColNames    string
	ExcludeResultColNames string
	SpecialEscapeChar     = "`"
	TemplateNames         = "dto|dao|service|controller|cacheloader|cachemodel"
	PackageNames          = "dto-models|dao-daos|service-services|controller-controllers|cacheloader-cacheloader|cachemodel-cachemodels"
)

func UpdateTableColumns(module *Module) {
	for _, m := range module.Models.Models {

		if m.Type == "table" && len(m.Fields) == 0 { //needs column defintion from informat_schema
			cols := getTableFields(m.SqlName)
			m.Fields = make([]*Field, 0)
			wordSplit := strings.Split(m.SqlName, "_")
			if len(wordSplit) == 1 {
				m.ShortName = string(m.SqlName[0:2])
			} else {
				for _, w := range wordSplit {
					m.ShortName = m.ShortName + string(w[0])
				}
			}
			if len(m.Name) == 0 {
				m.Name = strcase.ToCamel(m.SqlName)
			}
			m.Comments = cols[0].TableComment

			if len(m.Template) == 0 {
				m.Template = TemplateNames
			}
			for _, c := range cols {
				f := Field{}
				f.Name = strcase.ToCamel(c.ColumnName)
				f.Type, _, _ = Db.goType(c)
				f.DbFieldName = c.ColumnName
				f.FieldType = c.ColumnType
				m.Fields = append(m.Fields, &f)
			}

		}
		fmt.Println("Inside update model", m)
	}

}
func updateSpColumns(module *Module) {
	for _, m := range module.Models.Models {

		if m.Type == "sp" && len(m.Fields) == 0 { //needs column defintion from informat_schema
			cols := getTableFields(m.SqlName)
			m.Fields = make([]*Field, 0)
			wordSplit := strings.Split(m.SqlName, "_")
			if len(wordSplit) == 1 {
				m.ShortName = string(m.SqlName[0:2])
			} else {
				for _, w := range wordSplit {
					m.ShortName = m.ShortName + string(w[0])
				}
			}

			m.Name = strcase.ToCamel(m.SqlName)
			m.Comments = cols[0].TableComment

			if len(m.Template) == 0 {
				m.Template = TemplateNames
			}
			for _, c := range cols {
				f := Field{}
				f.Name = strcase.ToCamel(c.ColumnName)
				f.Type, _, _ = Db.goType(c)
				f.DbFieldName = c.ColumnName
				f.FieldType = c.ColumnType
				m.Fields = append(m.Fields, &f)
			}

		}
		fmt.Println("Inside update model", m)
	}

}

func getTableFields(tableName string) []*ColumnSchema {
	fields := make([]*ColumnSchema, 0)
	for _, c := range TableColumns {
		if strings.EqualFold(tableName, c.TableName) {
			fields = append(fields, c)
		}
	}
	return fields
}
func getProcedureFields(sp string) []*ColumnSchema {
	fields := make([]*ColumnSchema, 0)
	for _, c := range SpColumns {
		if strings.EqualFold(sp, c.TableName) {
			fields = append(fields, c)
		}
	}
	return fields
}

func GetEntity(fileName string) Module {
	fileName = fmt.Sprintf("%s/codegendef/%s", Wd, fileName)
	xmlFile, err := os.Open(fileName)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(fileName)
	// defer the closing of our xmlFile so that we can parse it later on
	defer xmlFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(xmlFile)

	// we initialize our Users array
	// we unmarshal our byteArray which contains our
	// xmlFiles content into 'users' which we defined above
	Entity.Models.Models = make([]*Model, 0)
	xml.Unmarshal(byteValue, &Entity)
	fmt.Println("entity unmarshal", &Entity)
	for _, m := range Entity.Models.Models {
		fmt.Println("entity", m.Name, m.SqlName)
		for _, f := range m.Fields {
			fmt.Println(m.Name, f.Name, f.DbFieldName)
		}

	}
	return Entity
}

func GenerateCode(module *Module) {
	pkgNames := make(map[string]string)
	for _, v := range strings.Split(PackageNames, "|") {
		vsplit := strings.Split(v, "-")
		pkgNames[vsplit[0]] = vsplit[1]
	}
	for _, m := range module.Models.Models {
		//fileNames := strings.Split(m.FileName, "|")
		templates := strings.Split(m.Template, "|")
		fmt.Println("models", m)
		//folders := strings.Split(m.LayerFolder, "|")
		m.EscapeChar = SpecialEscapeChar
		for index, fn := range templates {
			folder := pkgNames[fn]
			m.PackageName = pkgNames[fn]
			if FlagEnableModule {
				folder = Entity.Name
				m.PackageName = Entity.Name
			}
			// entityName := entity.Name
			// if entity.IncludeInPath == "false" {
			// 	entityName = ""
			// }
			p := path.Join(Wd, folder)
			fmt.Println(p)
			if _, err := os.Stat(p); os.IsNotExist(err) {
				fmt.Println(err)
				if err := os.MkdirAll(p, os.ModePerm); err != nil {
					fmt.Println(err)
					panic(err)
				}
			}
			var fileName string
			fmt.Println("extension", FlagExtension)
			if len(FlagExtension) == 0 {
				fileName = fmt.Sprintf("%s.%s.go", m.Name, fn)
			} else {
				fileName = fmt.Sprintf("%s.%s.go.%s", m.Name, fn, FlagExtension)
			}
			file, err := os.Create(path.Join(p, fileName))

			if err != nil {
				fmt.Println(err)
			}
			if err := Tmpl.ExecuteTemplate(file, templates[index], *m); err != nil {
				fmt.Println(err)
			}
		}
	}

}

func GenerateCodeForFlutter(module *Module) {
	pkgNames := make(map[string]string)
	for _, v := range strings.Split(PackageNames, "|") {
		vsplit := strings.Split(v, "-")
		pkgNames[vsplit[0]] = vsplit[1]
	}
	for _, m := range module.Models.Models {
		//fileNames := strings.Split(m.FileName, "|")
		templates := strings.Split(m.Template, "|")
		fmt.Println("models", m)
		//folders := strings.Split(m.LayerFolder, "|")
		m.EscapeChar = SpecialEscapeChar
		folder := strings.ToLower(m.Name)
		for index, fn := range templates {
			//folder := pkgNames[fn]
			m.PackageName = pkgNames[fn]
			// if FlagEnableModule {
			// 	folder = Entity.Name
			// 	m.PackageName = Entity.Name
			// }
			// entityName := entity.Name
			// if entity.IncludeInPath == "false" {
			// 	entityName = ""
			// }
			p := path.Join(Wd, folder)
			fmt.Println(p)
			if _, err := os.Stat(p); os.IsNotExist(err) {
				fmt.Println(err)
				if err := os.MkdirAll(p, os.ModePerm); err != nil {
					fmt.Println(err)
					panic(err)
				}
			}
			var fileName string
			// fmt.Println("file ext", m.FileNameExt)
			// 	fileName = fmt.Sprintf("%s.%s.go", m.Name, fn)
			// } else {
			fileName = fmt.Sprintf("%s.%s.%s", strings.ToLower(m.Name), fn, m.FileNameExt)
			//}
			file, err := os.Create(path.Join(p, fileName))

			if err != nil {
				fmt.Println(err)
			}
			if err := Tmpl.ExecuteTemplate(file, templates[index], *m); err != nil {
				fmt.Println(err)
			}
		}
	}

}
