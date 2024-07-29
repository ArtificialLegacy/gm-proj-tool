package main

import (
	"github.com/ArtificialLegacy/gm-proj-tool/pkg/yyp"
)

const TEST_DIR = "gm-proj-tool-testing"

func main() {
	proj, err := yyp.NewProject(TEST_DIR)
	if err != nil {
		panic(err)
	}

	note := yyp.NewNote("Note2", "test", proj.AsParent())

	err = proj.ImportResource(note)
	if err != nil {
		panic(err)
	}

	err = proj.DataSave()
	if err != nil {
		panic(err)
	}

	err = proj.OrderSave()
	if err != nil {
		panic(err)
	}
}
