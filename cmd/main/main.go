package main

import "github.com/ArtificialLegacy/gm-proj-tool/pkg/yyp"

const TEST_DIR = "gm-proj-tool-testing"

func main() {
	proj, err := yyp.NewProject(TEST_DIR)
	if err != nil {
		panic(err)
	}

	err = proj.DataLoad()
	if err != nil {
		panic(err)
	}

	println(proj.Data.Name)
}
