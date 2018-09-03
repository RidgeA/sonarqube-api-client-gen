/*
This is a tool to generate client library for SonarQube (https://www.sonarqube.org/) web-api
It allows to generate client library based on your server version.
*/
package main

import (
	"flag"
	"log"
	"os"
)

// flags
var (
	host          string
	deprecated    bool
	internal      bool
	targetVersion string
	help          bool
	out           string
	packageName   string
)

func parseFlags() {
	flag.StringVar(&host, "host", "http://localhost:9000", "SonarQube server")
	flag.BoolVar(&deprecated, "deprecated", false, "generate code for deprecated api methods (default: false)")
	flag.BoolVar(&internal, "internal", false, "generate code for internal methods (default: false)")
	flag.StringVar(&targetVersion, "target", "", "set target api version (default: server's version)")
	flag.BoolVar(&help, "help", false, "show usage")
	flag.StringVar(&out, "out", ".", "output directory")
	flag.StringVar(&packageName, "package", "", "package name, if not set will be sonarqube_client")
	flag.Parse()

	if help {
		flag.PrintDefaults()
		os.Exit(0)
	}
}

func main() {
	parseFlags()

	var err error
	var def *apiDefinition

	if def, err = loadAPI(nil, host, deprecated, internal, targetVersion); err != nil {
		log.Fatal(err)
	}

	if err = generateCode(def, out); err != nil {
		log.Fatal(err)
	}
}
