package check

import (
	"encoding/json"
	"log"
	"os"

	resource "github.com/pivotal-cf-experimental/new_version_resource"
)

func Main() {
	var request resource.CheckRequest
	inputRequest(&request)

	command := resource.NewCheckCommand()
	response, err := command.Run(request)
	if err != nil {
		log.Fatal("running command", err)
	}

	outputResponse(response)
}

func inputRequest(request *resource.CheckRequest) {
	if err := json.NewDecoder(os.Stdin).Decode(request); err != nil {
		log.Fatal("reading request from stdin", err)
	}
}

func outputResponse(response []resource.Version) {
	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		log.Fatal("writing response to stdout", err)
	}
}
