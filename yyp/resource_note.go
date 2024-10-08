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

func NewResourceNote(name string, parent ProjectResourceNode) *ResourceNote {
	return &ResourceNote{
		ResourceType:    RESTYPE_NOTE,
		ResourceVersion: VERSION_NOTE,

		Name:   name,
		Parent: parent,
	}
}

func NewNote(name, text string, parent ProjectResourceNode) *Note {
	return &Note{
		Name:     name,
		Text:     text,
		Resource: NewResourceNote(name, parent),
	}
}

func (n *Note) Save(pdir string) (string, string, *ProjectResourceNode, error) {
	d := path.Join(pdir, DIR_NOTE, n.Name)

	f, err := os.Stat(d)
	if err != nil {
		err := os.Mkdir(d, 0o777)
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

func (p *Project) NoteExists(name string) bool {
	pth := path.Join(p.Path, DIR_NOTE, name)

	fs, err := os.Stat(pth)
	if err != nil {
		return false
	}
	if !fs.IsDir() {
		return false
	}

	return true
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

func (p *Project) NoteDelete(name string) error {
	pth := path.Join(p.Path, DIR_NOTE, name)

	fs, err := os.Stat(pth)
	if err != nil {
		return nil // Trying to delete smth that doesn't exist shouldn't be an error
	}
	if !fs.IsDir() {
		return fmt.Errorf("resource %s is not a directory", fs.Name())
	}

	err = os.RemoveAll(pth)
	if err != nil {
		return fmt.Errorf("failed to delete resource directory: %s", err)
	}

	idPath := path.Join(DIR_NOTE, name, name+EXT_RESOURCE)
	for i, r := range p.Data.Resources {
		if r.ID.Name == name && r.ID.Path == idPath {
			p.Data.Resources = append(p.Data.Resources[:i], p.Data.Resources[i+1:]...)
			break
		}
	}

	return nil
}

type ResourceNote struct {
	ResourceType    ResourceType        `json:"resourceType"`
	ResourceVersion Version             `json:"resourceVersion"`
	Name            string              `json:"name"`
	Parent          ProjectResourceNode `json:"parent"`
	Tags            []string            `json:"tags,omitempty"`
}
