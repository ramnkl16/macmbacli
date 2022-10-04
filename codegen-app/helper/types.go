package helper

import (
	"database/sql"
	"encoding/xml"
)

type Module struct {
	ModuleName    xml.Name `xml:"module"`
	Name          string   `xml:"name,attr"`
	IncludeInPath string   `xml:"includeInPath,attr"`
	Models        Models   `xml:"models"`
}

type Models struct {
	XMLName xml.Name `xml:"models"`
	Name    string   `xml:"name,attr"`
	Models  []*Model `xml:"model"`
}

type Model struct {
	XMLName     xml.Name `xml:"model"`
	Name        string   `xml:"name,attr"`
	ModelName   string   `xml:"modelName,attr"`
	Type        string   `xml:"type,attr"`
	SqlName     string   `xml:"sql,attr"`
	Fields      []*Field `xml:"field"`
	Template    string   `xml:"template,attr"`
	FileName    string   `xml:"fileName,attr"`
	LayerFolder string   `xml:"layerFolder,attr"`
	ShortName   string   `xml:"shortName,attr"`
	EscapeChar  string   `xml:"escapeChar,attr"`
	Comments    string   `xml:"comments,attr"`
	FileNameExt string   `xml:"fileNameExt,attr"`
	RestApiPath string   `xml:"restApiPath,attr"`
	PackageName string
}

type Field struct {
	XMLName       xml.Name `xml:"field"`
	Type          string   `xml:"type,attr"`
	DartType      string   `xml:"dartType,attr"`
	Label         string   `xml:"label,attr"`
	Name          string   `xml:"name,attr"`
	TextInputType string   `xml:"textInputType,attr"`
	MaxCharLimit  string   `xml:"maxCharLimit,attr"`
	InputType     string   `xml:"inputType,attr"`

	Default     string `xml:"default,attr"`
	DbFieldName string `xml:"dbFieldName,attr"`
	FieldType   string `xml:"FieldType,attr"`
	IsNullable  string `xml:"isNullable,attr"`
}

type DbType interface {
	getSchema(scheme string) []ColumnSchema
	goType(col *ColumnSchema) (string, string, error)
	//getEnvMap() map[string]string
}
type TableEntity struct {
	TableName   string
	PackageName string
	EscapeChar  string
	Fields      []*ColumnSchema
}

type Configuration struct {
	DbUser     string `json:"db_user"`
	DbPassword string `json:"db_password"`
	DbName     string `json:"db_name"`
	DbHost     string `json:"db_host"`
	DbPort     int    `json:"db_port"`
}

type ColumnSchema struct {
	TableName              string
	ColumnName             string
	IsNullable             string
	DataType               string
	CharacterMaximumLength sql.NullInt64
	NumericPrecision       sql.NullInt64
	NumericScale           sql.NullInt64
	ColumnType             string
	ColumnKey              string
	GoType                 string
	TableComment           string
	JsonName               string
}
