package helper

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/iancoleman/strcase"
)

var (
	TableColumns []*ColumnSchema
	SpColumns    []*ColumnSchema
)

type MySQL struct{}

func GetEntityFromSql() map[string]*Module {
	Db = MySQL{}
	columns := GetTableSchema(Config)
	modules := make(map[string]*Module)

	//build modules from table name
	for _, c := range columns {
		tableSplit := strings.Split(c.TableName, "_")
		module := modules[tableSplit[0]]
		if len(tableSplit) > 1 && module == nil {
			module.Models = Models{}
			module.Name = strings.ToLower(tableSplit[0])
			modules[tableSplit[0]] = module
		}
		model, isExist := findModel(module.Models.Models, c.TableName)
		if !isExist {
			model = &Model{}
			model.Fields = make([]*Field, 0)
			model.Name = strcase.ToCamel(c.TableName)
			model.SqlName = c.TableName
			model.Comments = c.TableComment
			module.Models.Models = append(module.Models.Models, model)
			model.Template = TemplateNames

		}
		field := Field{}
		field.Name = strcase.ToCamel(c.ColumnName)
		field.Type, _, _ = Db.goType(c)
		field.DbFieldName = c.ColumnName
		field.FieldType = c.ColumnType
		model.Fields = append(model.Fields, &field)
	}
	return modules
}

func findModel(models []*Model, tableName string) (*Model, bool) {
	var m Model
	for _, t := range models {
		if t.Name == tableName {
			return t, true
		}
	}
	return &m, false
}
func GetTableSchema(config Configuration) []*ColumnSchema {
	scheme := fmt.Sprintf("%s:%s@tcp(%s:%d)/information_schema", *FlagDbUser, *FlagPassword, config.DbHost, config.DbPort)
	log.Println(scheme)
	conn, err := sql.Open("mysql", scheme)
	if err != nil {
		sqlErr, ok := err.(*mysql.MySQLError)
		fmt.Println("Failed open sql", err, sqlErr, ok)

	}
	defer conn.Close()
	q := "SELECT c.TABLE_NAME, COLUMN_NAME, IS_NULLABLE, DATA_TYPE, COLUMN_TYPE, COLUMN_KEY, t.`TABLE_COMMENT` FROM COLUMNS c " +
		"INNER JOIN `tables` t ON t.table_name = c.table_name AND  t.TABLE_COMMENT NOT LIKE '%exclude%' " +
		"WHERE c.`TABLE_SCHEMA` = ? AND t.`TABLE_SCHEMA`=? ORDER BY TABLE_NAME, ORDINAL_POSITION"

	rows, err := conn.Query(q, config.DbName, config.DbName)
	log.Println(("sql executed for " + config.DbName))
	if err != nil {
		log.Fatal(err)
	}
	columns := []*ColumnSchema{}
	for rows.Next() {
		cs := ColumnSchema{}
		err := rows.Scan(&cs.TableName, &cs.ColumnName, &cs.IsNullable, &cs.DataType,
			&cs.ColumnType, &cs.ColumnKey, &cs.TableComment)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println(cs.ColumnName)
		columns = append(columns, &cs)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return columns
}

func GetProcedureParameters(
	config Configuration) []*ColumnSchema {
	scheme := fmt.Sprintf("%s:%s@tcp(%s:%d)/information_schema", *FlagDbUser, *FlagPassword, config.DbHost, config.DbPort)
	log.Println(scheme)
	conn, err := sql.Open("mysql", scheme)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	q := "SELECT  `SPECIFIC_NAME` SpName, `PARAMETER_NAME` parameterName, DATA_TYPE FROM PARAMETERS where SPECIFIC_SCHEMA = ? AND PARAMETER_NAME IS NOT NULL ORDER BY ORDINAL_POSITION"
	rows, err := conn.Query(q, config.DbName)
	log.Println(("sql executed for " + config.DbName))
	if err != nil {
		log.Fatal(err)
	}
	columns := []*ColumnSchema{}
	for rows.Next() {
		cs := ColumnSchema{}
		err := rows.Scan(&cs.TableName, &cs.ColumnName, &cs.DataType)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println(cs.ColumnName)
		columns = append(columns, &cs)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return columns
}

func (m MySQL) goType(col *ColumnSchema) (string, string, error) {
	requiredImport := ""
	// if col.IsNullable == "YES" {
	// 	requiredImport = "database/sql"
	// }
	var gt string = ""
	switch col.DataType {
	case "char", "varchar", "enum", "text", "longtext", "mediumtext", "tinytext":
		// if col.IsNullable == "YES" {
		// 	gt = "sql.NullString"
		// } else {
		gt = "string"
		//}
	case "blob", "mediumblob", "longblob", "varbinary", "binary":
		gt = "[]byte"
	case "date", "time", "datetime", "timestamp":
		// if col.IsNullable == "YES" {
		// 	gt = "sql.NullString"
		// } else {
		gt = "string"
		//}
	case "smallint", "int", "mediumint":
		// if col.IsNullable == "YES" {
		// 	gt = "sql.NullInt32"
		// } else {
		gt = "int32"
		//}
	case "bigint":
		// if col.IsNullable == "YES" {
		// 	gt = "sql.NullInt64"
		// } else {
		gt = "int64"
		//}

	case "float", "decimal", "double":
		// if col.IsNullable == "YES" {
		// 	gt = "sql.NullFloat64"
		// } else {
		gt = "float64"
		//}
	case "bit", "tinyint":
		gt = "int8"

	}
	if gt == "" {
		//fmt.Println("no compatible datatype :" + col.ColumnName)
		n := col.TableName + "." + col.ColumnName
		return "", "", errors.New("No compatible datatype (" + col.DataType + ") for " + n + " found")
	}
	return gt, requiredImport, nil
}

// func (m MySQL) getEnvMap() map[string]string {
// 	return map[string]string{
// 		EnvHostKey:     "MYSQL_HOST",
// 		EnvPortKey:     "MYSQL_PORT",
// 		EnvDataBaseKey: "MYSQL_DATABASE",
// 		EnvUserKey:     "MYSQL_USER",
// 		EnvPasswordKey: "MYSQL_PASSWORD",
// 	}
//}
