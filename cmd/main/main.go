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

	//folder := yyp.NewFolder("GoFolder1")
	folder, err := proj.FolderLoad("GoFolder1")
	if err != nil {
		panic(err)
	}
	proj.FolderSave(folder)

	note := yyp.NewNote("GoNote1", "test", folder.AsParent())

	err = proj.ImportResource(note)
	if err != nil {
		panic(err)
	}

	err = proj.DataSave()
	if err != nil {
		panic(err)
	}
}
