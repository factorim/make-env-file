package main

import (
	"log"

	"github.com/factorim/make-env-file/makeenvfile"
)

func init() {
	log.SetPrefix("make-env ")
}

func main() {
	// Parse command line options
	sourceFile, destFile, overwrite := makeenvfile.GetFlags()
	log.Printf("create %s from %s", destFile, sourceFile)

	// Check env
	report, err := makeenvfile.CheckEnv(sourceFile, destFile)
	if err != nil {
		log.Fatal(err)
	}
	// Make env
	err = makeenvfile.MakeEnv(sourceFile, destFile, overwrite, report)
	if err != nil {
		log.Fatal(err)
	}
}
