package main

import (
	"regexp"
	"strings"
	"text/template"
)

var snakeToCamelRE = regexp.MustCompile("_([a-z])")

func snakeToCamel(str string) string {
	return snakeToCamelRE.ReplaceAllStringFunc(str, func(match string) string {
		return strings.ToUpper(strings.TrimPrefix(match, "_"))
	})
}

func makeExported(str string) string {
	first := string(str[0])
	return strings.Replace(str, first, strings.ToUpper(first), 1)
}

func makeUnexported(str string) string {
	first := string(str[0])
	return strings.Replace(str, first, strings.ToLower(first), 1)
}

func sanitizeItentifier(str string) string {
	return strings.Replace(str, "-", "_", -1)
}

func tick() string {
	return "`"
}

func replaceTags(str string) string {
	repl := strings.NewReplacer(
		"<br> ", "\n// ",
		"<br>", "\n// ",
		"<br/>", "\n// ",
		"<br />", "\n// ",
		"<ul>", "\n// ",
		"<li>", " * ",
		"</li>", "\n// ",
		"</ul>", "",
	)
	return repl.Replace(str)
}

var templateHelpers = template.FuncMap{
	"formatDescription": replaceTags,
	"tick":              tick,
}
