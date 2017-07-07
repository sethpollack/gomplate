package main

import (
	"io/ioutil"
	"log"
	"net/url"
	"text/template"

	yaml "gopkg.in/yaml.v2"

	"github.com/hairyhenderson/gomplate/funcs"
)

func ReadFile(path string) string {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}

func ReadDir(dir string) string {
	fileMap := make(map[string]string)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if !file.IsDir() {
			fileMap[file.Name()] = ReadFile(dir + "/" + file.Name())
		}
	}

	b, err := yaml.Marshal(fileMap)
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}

// initFuncs - The function mappings are defined here!
func initFuncs(data *Data) template.FuncMap {
	env := &Env{}
	typeconv := &TypeConv{}

	f := template.FuncMap{
		"readDir":          ReadDir,
		"readFile":         ReadFile,
		"getenv":           env.Getenv,
		"bool":             typeconv.Bool,
		"has":              typeconv.Has,
		"json":             typeconv.JSON,
		"jsonArray":        typeconv.JSONArray,
		"yaml":             typeconv.YAML,
		"yamlArray":        typeconv.YAMLArray,
		"toml":             typeconv.TOML,
		"csv":              typeconv.CSV,
		"csvByRow":         typeconv.CSVByRow,
		"csvByColumn":      typeconv.CSVByColumn,
		"slice":            typeconv.Slice,
		"join":             typeconv.Join,
		"toJSON":           typeconv.ToJSON,
		"toJSONPretty":     typeconv.toJSONPretty,
		"toYAML":           typeconv.ToYAML,
		"toTOML":           typeconv.ToTOML,
		"toCSV":            typeconv.ToCSV,
		"urlParse":         url.Parse,
		"datasource":       data.Datasource,
		"ds":               data.Datasource,
		"datasourceExists": data.DatasourceExists,
		"include":          data.include,
	}
	funcs.AWSFuncs(f)
	funcs.AddBase64Funcs(f)
	funcs.AddNetFuncs(f)
	funcs.AddReFuncs(f)
	funcs.AddStringFuncs(f)

	return f
}
