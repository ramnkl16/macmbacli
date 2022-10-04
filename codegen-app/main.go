package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"

	"dev.azure.com/RAMcoretech/booknplay/_git/codegen-app/helper"
)

func main() {
	flag.BoolVar(&helper.FlagMysql, `mysql`, false, ``)
	flag.BoolVar(&helper.FlagEnableModule, `enableModule`, false, ``)
	flag.StringVar(&helper.FlagExtension, "ext", "", "extension for code generation file")
	flag.StringVar(&helper.ExcludeDtoColNames, "edto", "", "extension for code generation file")
	flag.StringVar(&helper.ExcludeResultColNames, "erc", "active_flag|updated_by|updated_at", "extension for code generation file")

	flag.Parse()
	fmt.Println("dbhost", *helper.FlagDbHost)

	helper.Tmpl = template.Must(template.New("types").Funcs(helper.NewTemplateFuncs()).ParseGlob("./codegendef/*.tmpl"))
	helper.Wd, _ = os.Getwd()
	dir, err := os.Open(helper.Wd + "/codegendef")
	if err != nil {
		fmt.Println("No codegendef folder. Pleae try after create code defintion file")
		return
	}
	helper.Config.DbPassword = *helper.FlagPassword
	helper.Config.DbHost = *helper.FlagDbHost
	helper.Config.DbName = *helper.FlagDbName
	helper.Config.DbPort = *helper.FlagDbPort
	helper.Config.DbUser = *helper.FlagDbUser

	helper.IgnoreNames = make(map[string]bool)
	for _, v := range strings.Split(*helper.FlagIgnore, ",") {
		helper.IgnoreNames[v] = true
	}

	helper.TableColumns = helper.GetTableSchema(helper.Config)
	helper.SpColumns = helper.GetProcedureParameters(helper.Config)
	//read defintion files and generated code
	list, _ := dir.Readdirnames(0) // 0 to read all files and folders
	for _, name := range list {
		filePart := strings.Split(name, ".")
		fmt.Println(name)
		if strings.ToLower(filePart[len(filePart)-1]) == "xml" {
			module := helper.GetEntity(name)
			fmt.Print(module)
			//helper.UpdateTableColumns(&module)
			//helper.GenerateCode(&module)
			helper.GenerateCodeForFlutter(&module)
		}
		// moduels := GetEntityFromSql()
		// for _, v := range moduels {
		// 	GenerateCode(v)
		// }

	}
}
