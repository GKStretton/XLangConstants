package main

import (
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

var (
	input  = flag.String("input", "/input", "specify input directory containing yaml files")
	output = flag.String("output", "/output", "specify output directory for generated code files")
)

type TemplateData struct {
	Constants map[string]string
	Package   string
}

func main() {
	flag.Parse()

	fmt.Printf("running XLangConstants for input: %s, output: %s\n", *input, *output)

	err := os.MkdirAll(*output, os.ModePerm)
	if err != nil {
		panic(err)
	}

	err = filepath.WalkDir(*input, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			processFile(path)
		}

		return nil
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("done")
}

func processFile(inputFile string) {
	fmt.Printf("processing file: %s\n", inputFile)

	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}

	var yamlData map[string]map[string]string
	err = yaml.Unmarshal(data, &yamlData)
	if err != nil {
		panic(err)
	}

	name := basenameWithoutExt(inputFile)
	templateData := TemplateData{
		Constants: yamlData["constants"],
		Package:   name,
	}

	// Generate code files for each language
	// todo: iterate template files instead
	generateCodeFile("./templates/python.tmpl", filepath.Join(*output, "python", name, fmt.Sprintf("%s.py", name)), templateData)
	generateCodeFile("./templates/c.tmpl", filepath.Join(*output, "c", name, fmt.Sprintf("%s.h", name)), templateData)
	generateCodeFile("./templates/go.tmpl", filepath.Join(*output, "go", name, fmt.Sprintf("%s.go", name)), templateData)
	generateCodeFile("./templates/typescript.tmpl", filepath.Join(*output, "typescript", name, fmt.Sprintf("%s.ts", name)), templateData)
}

func generateCodeFile(templateFile, outputFile string, templateData TemplateData) {
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		panic(err)
	}

	fmt.Printf("creating dir: %s\n", filepath.Dir(outputFile))
	err = os.MkdirAll(filepath.Dir(outputFile), os.ModePerm)
	if err != nil {
		panic(err)
	}

	fmt.Printf("creating file: %s\n", outputFile)
	outFile, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	err = tmpl.Execute(outFile, templateData)
	if err != nil {
		panic(err)
	}
}

func basenameWithoutExt(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	return strings.TrimSuffix(base, ext)
}
