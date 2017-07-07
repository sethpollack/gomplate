package main

import (
	"io"
	"log"
	"text/template"

	"github.com/Masterminds/sprig"
)

func (g *Gomplate) createTemplate() *template.Template {
	return template.New("template").Funcs(sprig.TxtFuncMap()).Funcs(g.funcMap).Option("missingkey=error")
}

// Gomplate -
type Gomplate struct {
	funcMap    template.FuncMap
	leftDelim  string
	rightDelim string
}

// RunTemplate -
func (g *Gomplate) RunTemplate(text string, out io.Writer) {
	context := &Context{}
	tmpl, err := g.createTemplate().Delims(g.leftDelim, g.rightDelim).Parse(text)
	if err != nil {
		log.Fatalf("Line %q: %v\n", text, err)
	}

	if err := tmpl.Execute(out, context); err != nil {
		panic(err)
	}
}

// NewGomplate -
func NewGomplate(data *Data, o *GomplateOpts) *Gomplate {
	return &Gomplate{
		leftDelim:  o.lDelim,
		rightDelim: o.rDelim,
		funcMap:    initFuncs(data, o),
	}
}

func runTemplate(o *GomplateOpts) error {
	defer runCleanupHooks()
	data := NewData(o.dataSources, o.dataSourceHeaders)

	g := NewGomplate(data, o)

	if o.inputDir != "" {
		return processInputDir(o.inputDir, o.outputDir, g)
	}

	return processInputFiles(o.input, o.inputFiles, o.outputFiles, g)
}

// Called from process.go ...
func renderTemplate(g *Gomplate, inString string, outPath string) error {
	outFile, err := openOutFile(outPath)
	if err != nil {
		return err
	}
	// nolint: errcheck
	defer outFile.Close()
	g.RunTemplate(inString, outFile)
	return nil
}
