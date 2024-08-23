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

	note := yyp.NewNote("GoNote1", "test", proj.AsParent())

	err = proj.ImportResource(note)
	if err != nil {
		panic(err)
	}

	/*sprite, err := yyp.NewSprite("GoSprite1", proj.AsParent(), yyp.ProjectTextureGroupDefaultID(), 16, 16,
		[]yyp.SpriteLayer{
			{
				Name: yyp.SPRITELAYER_DEFAULTNAME,
				Frames: []*image.NRGBA{
					image.NewNRGBA(image.Rect(0, 0, 16, 16)),
				},
			},
			{
				Name: "Layer 1",
				Frames: []*image.NRGBA{
					image.NewNRGBA(image.Rect(0, 0, 16, 16)),
				},
			},
		})
	if err != nil {
		panic(err)
	}*/

	sprite, err := proj.SpriteLoad("GoSprite1")
	if err != nil {
		panic(err)
	}

	err = proj.ImportResource(sprite)
	if err != nil {
		panic(err)
	}

	/*script := yyp.NewScript("GoScript1", `
	function test() {
		// test script
	}
	`, proj.AsParent())*/

	script, err := proj.ScriptLoad("GoScript1")
	if err != nil {
		panic(err)
	}

	err = proj.ImportResource(script)
	if err != nil {
		panic(err)
	}

	//data := []byte("test")
	//includedFile := yyp.NewIncludedFile("GoIncludedFile.txt", yyp.INCLUDEDFILE_DEFAULTPATH, &data)

	includedFile, err := proj.IncludedFileLoad(yyp.INCLUDEDFILE_DEFAULTPATH, "GoIncludedFile.txt")
	if err != nil {
		panic(err)
	}

	err = proj.IncludedFileSave(includedFile)
	if err != nil {
		panic(err)
	}

	err = proj.DataSave()
	if err != nil {
		panic(err)
	}
}
