package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"text/template"

	"go/format"

	"github.com/mitch000001/go-hbci/bankinfo"
)

func main() {
	flag.Parse()
	bankdataFiles := flag.Args()

	parser := bankinfo.Parser{}
	var bankInfos []bankinfo.BankInfo
	for _, bankdata := range bankdataFiles {
		file, err := os.Open(bankdata)
		if err != nil {
			log.Fatal("Cannot open file: %q", bankdata)
			os.Exit(1)
		}
		infos, err := parser.Parse(file)
		if err != nil {
			log.Fatal("Parse error: %q", err)
			os.Exit(1)
		}
		bankInfos = append(bankInfos, infos...)
	}
	data, err := writeDataToGoFile(bankInfos)
	if err != nil {
		log.Fatal("Error while parsing expression: %q", err)
		os.Exit(1)
	}
	goFile, err := os.Create("bankinfo/bank_data.go")
	if err != nil {
		log.Fatal("Cannot create file: %q", err)
		os.Exit(1)
	}
	_, err = io.Copy(goFile, data)
	if err != nil {
		log.Fatal("Error while writing file: %q", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func writeDataToGoFile(data []bankinfo.BankInfo) (io.Reader, error) {
	t, err := template.New("bank_data").Parse(dataTemplate)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing template: %v", err)
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	if err != nil {
		return nil, fmt.Errorf("Error while executing template: %v", err)
	}
	formattedBytes, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(formattedBytes), nil
}

const dataTemplate = `package bankinfo

var data = []BankInfo{
	{{range $element := .}}BankInfo{
		BankId: "{{.BankId}}",
		VersionNumber: "{{.VersionNumber}}",
		URL: "{{.URL}}",
		VersionName: "{{.VersionName}}",
	},
	{{end}}
}`