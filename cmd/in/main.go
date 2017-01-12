package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	resource "github.com/pivotal-cf-experimental/new_version_resource"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage: %s <sources directory>\n", os.Args[0])
		os.Exit(1)
	}

	var request resource.InRequest
	inputRequest(&request)

	destDir := os.Args[1]
	os.MkdirAll(destDir, 0755)

	urlFile := filepath.Join(destDir, "url")
	var url string

	if request.Source.Type == "http" {
		url = request.Source.HTTP.URL
	} else if request.Source.Type == "git" {
		url = "https://github.com/" + request.Source.Git.Organization + "/" + request.Source.Git.Repo
	}

	err := ioutil.WriteFile(urlFile, []byte(url), 0644)
	if err != nil {
		fmt.Println(err)
	}

	err = ioutil.WriteFile(filepath.Join(destDir, "version"), []byte(request.Version.Version), 0644)
	if err != nil {
		fmt.Println(err)
	}

	fp, err := os.Create(filepath.Join(destDir, "input.json"))
	if err != nil {
		log.Fatalf("Unable to create %v. Err: %v.", "input.json", err)
	}
	defer fp.Close()

	encoder := json.NewEncoder(fp)
	if err = encoder.Encode(request); err != nil {
		log.Fatalf("Unable to encode Json file. Err: %v.", err)
	}

	response := resource.InResponse{
		Version: request.Version,
		Metadata: []resource.MetadataPair{
			{Name: "url", Value: url},
			{Name: "version", Value: request.Version.Version},
		},
	}

	outputResponse(response)
}

func inputRequest(request *resource.InRequest) {
	if err := json.NewDecoder(os.Stdin).Decode(request); err != nil {
		log.Fatal("reading request from stdin", err)
	}
}

func outputResponse(response resource.InResponse) {
	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		log.Fatal("writing response to stdout", err)
	}
}
