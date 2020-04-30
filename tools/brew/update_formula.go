package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/pkg/errors"
)

type formula struct {
	Version string
	Mac     formulaFile
	Linux   formulaFile
}
type formulaFile struct {
	Checksum string
	URL      string
}

func main() {
	var formula formula

	templatePath := flag.String("template", "", "template file to use")
	outputPath := flag.String("outputPath", "", "path to the file with the processed template")
	flag.StringVar(&formula.Version, "version", "", "eksctl version to publish")
	flag.StringVar(&formula.Linux.URL, "linux-url", "", "URL to eksctl binary file for Linux")
	flag.StringVar(&formula.Mac.URL, "mac-url", "", "URL to eksctl binary file for MacOs")

	flag.Parse()

	// Validate flags
	if *templatePath == "" {
		log.Fatal("missing templatePath")
	}
	if _, err := os.Stat(*templatePath); err != nil {
		log.Fatalf("unable to open template file. Check if it exists %q", *templatePath)
	}

	if *outputPath == "" {
		log.Fatal("missing outputPath")
	}

	if formula.Version == "" {
		log.Fatal("missing version")
	}

	if formula.Linux.URL == "" {
		log.Fatal("missing linuxUrl")
	}

	if formula.Mac.URL == "" {
		log.Fatal("missing macUrl")
	}

	// Calculate checksums and write template
	templateFile, err := readFile(*templatePath)
	if err != nil {
		log.Fatal(err.Error())
	}

	formula.Linux.Checksum, err = calculateSHA256(formula.Linux.URL)
	if err != nil {
		log.Fatalf("unable to download linux file %q: %s", formula.Linux.URL, err)
	}

	formula.Mac.Checksum, err = calculateSHA256(formula.Mac.URL)
	if err != nil {
		log.Fatalf("unable to download linux file %q: %s", formula.Mac.URL, err)
	}

	parsedTemplate, err := template.New(*templatePath).Parse(templateFile)
	if err != nil {
		log.Fatalf("could not parse template %q: %s", *templatePath, err.Error())
	}

	executedTemplate := bytes.NewBuffer(nil)
	if err = parsedTemplate.Execute(executedTemplate, formula); err != nil {
		log.Fatalf("could not apply values to template: %s", err.Error())
	}
	log.Println("Successfully executed template")

	log.Println("Writing file")
	if err := ioutil.WriteFile(*outputPath, executedTemplate.Bytes(), 0644); err != nil {
		log.Fatalf("error while writing executed template file %q: %s", *outputPath, err.Error())
	}

}
func calculateSHA256(url string) (string, error) {
	contents, err := downloadFile(url)
	if err != nil {
		return "", err
	}
	checksum := sha256.Sum256(contents)
	return hex.EncodeToString(checksum[:]), nil

}

func downloadFile(url string) ([]byte, error) {
	log.Printf("downloading file %q", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return contents, nil
}

func readFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", errors.Wrapf(err, "could not open file")
	}
	defer file.Close()
	contents, err := ioutil.ReadAll(file)
	if err != nil {
		return "", errors.Wrapf(err, "unable to read file %q", path)
	}
	return string(contents), err
}
