[![GoDoc](https://godoc.org/github.com/RidgeA/sonarqube-api-client-gen?status.svg)](https://godoc.org/github.com/RidgeA/sonarqube-api-client-gen)

# SonarQube WebApi Client Generator

This is a tool to generate client library for [SonarQube](https://www.sonarqube.org/) web-api 
It allows to generate client library based on your server version.

Tested with SQ v7.1

# __IMPORTANT__
__This is an alpha version.__

The API of generated code __will__ be changed when the SonarQube team implement api to [get response schema](https://community.sonarsource.com/t/is-it-possible-to-retrieve-response-schema/2184/3)

## Install

```
go get -u github.com/RidgeA/sonarqube-api-client-gen
```

## Usage of CLI tool

Make sure that your `$GOPATH/bin` directory listed in `$PATH` or run binary from `$GOPATH/bin`

Execute to generate client library with default settings
``` 
    sonarqube-api-client-gen
```

Available options:
```
  -deprecated
    	generate code for deprecated api methods (default: false)
  -help
    	show usage
  -host string
    	SonarQube server (default "http://localhost:9000")
  -internal
    	generate code for internal methods (default: false)
  -out string
    	output directory (default ".")
  -package string
    	package name, if not set will be sonarqube_client
  -target string
    	set target api version (default: server's version)
```

## Usage of generated code

Generated code has to external depends on two external dependencies:
* github.com/pkg/errors 
* github.com/google/go-querystring/query

You have to install them in your project manually (if you don't have them already), e.g.:

```
package main

import (
	"context"
	"encoding/json"
	"log"
 	sq "/sonarqube_client"
)

func main() {
	host := "http://localhost:9000"
	
	// see https://docs.sonarqube.org/display/DEV/Web+API#WebAPI-Authentication
	username := "3dc205e8974ab24fdde395c722655be39c2ac861"
	password := ""
	
	ctx := context.Background()
	c := sq.NewClient(nil, host, username, password)
 	params := &sq.ProjectsServiceCreateRequest{
		Name:    sq.String("project_name"),
		Project: sq.String("project_key"),
	}
	result, err := c.Projects().Create(ctx, params)
	if err != nil {
		log.Fatal(err)
	}
	
 	out := make(map[string]interface{})
	dec := json.NewDecoder(result.RawResponse.Body)
	if err := dec.Decode(&out); err != nil {
		log.Fatal(err)
	}
	outPretty, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("result:\n%s", outPretty)
 }

```

## TODO

 - [ ] Change response api
 - [ ] Add request params validation
