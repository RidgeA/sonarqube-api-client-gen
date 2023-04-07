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
	auth          string
	packageName   string
)

var mainFlagsSet = flag.NewFlagSet("", flag.PanicOnError)

func parseFlags() {
	mainFlagsSet.StringVar(&host, "host", "http://localhost:9000", "SonarQube server")
	mainFlagsSet.BoolVar(&deprecated, "deprecated", false, "generate code for deprecated api methods (default: false)")
	mainFlagsSet.BoolVar(&internal, "internal", false, "generate code for internal methods (default: false)")
	mainFlagsSet.StringVar(&targetVersion, "target", "", "set target api version (default: server's version)")
	mainFlagsSet.BoolVar(&help, "help", false, "show usage")
	mainFlagsSet.StringVar(&out, "out", ".", "output directory")
	mainFlagsSet.StringVar(&auth, "auth", "", "the header Authorization value,example: Basic YWRtaW46YWRtaW4=")
	mainFlagsSet.StringVar(&packageName, "package", "", "package name, if not set will be sonarqube_client")
	mainFlagsSet.Parse(os.Args[1:])
	if help {
		mainFlagsSet.Usage()
		os.Exit(0)
	}
}

func main() {
	parseFlags()

	var err error
	var def *apiDefinition

	if def, err = loadAPI(nil, host, deprecated, internal, targetVersion, auth); err != nil {
		log.Fatal(err)
	}

	if err = generateCode(def, out); err != nil {
		log.Fatal(err)
	}
}
