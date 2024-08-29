package yyp

import (
	"fmt"
	"os"
	"path"
)

type Script struct {
	Name     string
	Code     string
	Resource *ResourceScript
}

func NewScript(name, code string, parent ProjectResourceNode) *Script {
	return &Script{
		Name:     name,
		Code:     code,
		Resource: NewResourceScript(name, parent),
	}
}

func (n *Script) Save(pdir string) (string, string, *ProjectResourceNode, error) {
	d := path.Join(pdir, DIR_SCRIPT, n.Name)

	f, err := os.Stat(d)
	if err != nil {
		err := os.Mkdir(d, 0o777)
		if err != nil {
			return "", "", nil, err
		}
	} else {
		if !f.IsDir() {
			return "", "", nil, fmt.Errorf("path for script already exists, and it's not a directory")
		}
	}

	err = os.WriteFile(path.Join(d, n.Name+EXT_SCRIPT), []byte(n.Code), 0o666)
	if err != nil {
		return "", "", nil, err
	}

	err = saveJSON(path.Join(d, n.Name+EXT_RESOURCE), n.Resource)
	if err != nil {
		return "", "", nil, err
	}

	return n.Name, path.Join(DIR_SCRIPT, n.Name, n.Name+EXT_RESOURCE), &n.Resource.Parent, nil
}

func (p *Project) ScriptExists(name string) bool {
	pth := path.Join(p.Path, DIR_SCRIPT, name)

	fs, err := os.Stat(pth)
	if err != nil {
		return false
	}
	if !fs.IsDir() {
		return false
	}

	return true
}

func (p *Project) ScriptLoad(name string) (*Script, error) {
	pth := path.Join(p.Path, DIR_SCRIPT, name)

	fs, err := os.Stat(pth)
	if err != nil {
		return nil, err
	}
	if !fs.IsDir() {
		return nil, fmt.Errorf("resource %s is not a directory", fs.Name())
	}

	data := &ResourceScript{}
	err = loadJSON(path.Join(pth, name+EXT_RESOURCE), data)
	if err != nil {
		return nil, fmt.Errorf("error parsing script json: %s", err)
	}

	if data.ResourceVersion != VERSION_SCRIPT {
		return nil, fmt.Errorf("resource version is unsupported: %s, expected %s", data.ResourceVersion, VERSION_SCRIPT)
	}
	if data.ResourceType != RESTYPE_SCRIPT {
		return nil, fmt.Errorf("resource is of incorrect type: %s, expected %s", data.ResourceType, RESTYPE_SCRIPT)
	}

	if data.IsDND {
		return nil, fmt.Errorf("DnD scripts are not supported")
	}

	code, err := os.ReadFile(path.Join(pth, name+EXT_SCRIPT))
	if err != nil {
		return nil, err
	}

	script := &Script{
		Name:     name,
		Code:     string(code),
		Resource: data,
	}
	return script, nil
}

func (p *Project) ScriptDelete(name string) error {
	pth := path.Join(p.Path, DIR_SCRIPT, name)

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

	idPath := path.Join(DIR_SCRIPT, name, name+EXT_RESOURCE)
	for i, r := range p.Data.Resources {
		if r.ID.Name == name && r.ID.Path == idPath {
			p.Data.Resources = append(p.Data.Resources[:i], p.Data.Resources[i+1:]...)
			break
		}
	}

	return nil
}

type ResourceScript struct {
	ResourceType    ResourceType `json:"resourceType"`
	ResourceVersion Version      `json:"resourceVersion"`

	Name   string              `json:"name"`
	Parent ProjectResourceNode `json:"parent"`

	IsCompatibility bool `json:"isCompatibility"`
	IsDND           bool `json:"isDnD"`

	Tags []string `json:"tags,omitempty"`
}

func NewResourceScript(name string, parent ProjectResourceNode) *ResourceScript {
	return &ResourceScript{
		ResourceType:    RESTYPE_SCRIPT,
		ResourceVersion: VERSION_SCRIPT,

		Name:   name,
		Parent: parent,
	}
}
