package yyp

import (
	"fmt"
	"os"
	"path"
)

type Note struct {
	Name     string
	Text     string
	Resource *ResourceNote
}

type ResourceNote struct {
	ResourceType    string         `json:"resourceType"`
	ResourceVersion string         `json:"resourceVersion"`
	Name            string         `json:"name"`
	Parent          ResourceParent `json:"parent"`
	Tags            []string       `json:"tags,omitempty"`
}

func NewResourceNote(name string, parent ResourceParent) *ResourceNote {
	return &ResourceNote{
		ResourceType:    RESTYPE_NOTE,
		ResourceVersion: VERSION_NOTE,
		Name:            name,
		Parent:          parent,
	}
}

func NewNote(name, text string, parent ResourceParent) *Note {
	return &Note{
		Name:     name,
		Text:     text,
		Resource: NewResourceNote(name, parent),
	}
}

func (n *Note) Save(pdir string) (string, string, *ResourceParent, error) {
	d := path.Join(pdir, DIR_NOTE, n.Name)

	f, err := os.Stat(d)
	if err != nil {
		err := os.Mkdir(d, 0o666)
		if err != nil {
			return "", "", nil, err
		}
	} else {
		if !f.IsDir() {
			return "", "", nil, fmt.Errorf("path for note already exists, and it's not a directory")
		}
	}

	err = os.WriteFile(path.Join(d, n.Name+EXT_NOTE), []byte(n.Text), 0o666)
	if err != nil {
		return "", "", nil, err
	}

	err = saveJSON(path.Join(d, n.Name+EXT_RESOURCE), n.Resource)
	if err != nil {
		return "", "", nil, err
	}

	return n.Name, path.Join(DIR_NOTE, n.Name, n.Name+EXT_RESOURCE), &n.Resource.Parent, nil
}

func (p *Project) NoteLoad(name string) (*Note, error) {
	pth := path.Join(p.Path, DIR_NOTE, name)

	fs, err := os.Stat(pth)
	if err != nil {
		return nil, err
	}
	if !fs.IsDir() {
		return nil, fmt.Errorf("resource %s is not a directory", fs.Name())
	}

	data := &ResourceNote{}
	err = loadJSON(path.Join(pth, name+EXT_RESOURCE), data)
	if err != nil {
		return nil, fmt.Errorf("error parsing note json: %s", err)
	}

	if data.ResourceVersion != VERSION_NOTE {
		return nil, fmt.Errorf("resource version is unsupported: %s, expected %s", data.ResourceVersion, VERSION_NOTE)
	}
	if data.ResourceType != RESTYPE_NOTE {
		return nil, fmt.Errorf("resource is of incorrect type: %s, expected %s", data.ResourceType, RESTYPE_NOTE)
	}

	text, err := os.ReadFile(path.Join(pth, name+EXT_NOTE))
	if err != nil {
		return nil, err
	}

	note := &Note{
		Name:     name,
		Text:     string(text),
		Resource: data,
	}
	return note, nil
}
