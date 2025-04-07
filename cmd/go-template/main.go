package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/harness/go-template/internal"
)

var (
	versionMajor string = "0"
	versionMinor string = "1"
	versionPatch string = "1"
	gitref       string = "dev-default"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var valueFlags arrayFlags
var setFlags arrayFlags

func printUsage() {
	fmt.Println("Usage:")
	fmt.Printf("Template: %s <-t templateFilePath> [-f valuesFilePath] [-s variableOverride] [-o outputFolder]\n", os.Args[0])
	fmt.Printf("Version : %s -v\n", os.Args[0])
}

func main() {

	versionPtr := flag.Bool("v", false, "version")
	templatePtr := flag.String("t", "", "Path to template File/Folder")
	putPtr := flag.String("o", "", "Path to output Folder")
	flag.Var(&valueFlags, "f", "Path to Values file.")
	flag.Var(&setFlags, "s", "Set variable override.")
	flag.Parse()

	if *versionPtr {
		versionString := versionMajor
		if versionMinor != "" {
			versionString += "." + versionMinor
		}
		if versionPatch != "" {
			versionString += "." + versionPatch
		}
		if gitref != "" {
			versionString += " (" + gitref + ")"
		}
		fmt.Println(versionString)
		return
	}

	if *templatePtr == "" {
		printUsage()
		os.Exit(1)
	}

	internal.Render(*templatePtr, *putPtr, valueFlags)
}
