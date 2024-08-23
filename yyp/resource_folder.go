package yyp

import (
	"fmt"
	"path"
)

type Folder struct {
	Name     string
	Resource *ResourceFolder
}

func NewFolder(name string) *Folder {
	return &Folder{
		Name:     name,
		Resource: NewResourceFolder(name),
	}
}

func (f *Folder) AsParent() ProjectResourceNode {
	return ProjectResourceNode{
		Name: f.Resource.Name,
		Path: f.Resource.FolderPath,
	}
}

func (p *Project) FolderSave(folder *Folder) error {
	found := false
	for i, f := range p.Data.Folders {
		if f.FolderPath == folder.Resource.FolderPath {
			p.Data.Folders[i] = *folder.Resource
			found = true
			break
		}
	}

	if !found {
		p.Data.Folders = append(p.Data.Folders, *folder.Resource)
	}

	return nil
}

func (p *Project) FolderLoad(name string) (*Folder, error) {
	var resource *ResourceFolder
	for _, f := range p.Data.Folders {
		if f.Name == name {
			resource = &f
			break
		}
	}

	if resource == nil {
		return nil, fmt.Errorf("could not find folder: %s", name)
	}
	if resource.ResourceVersion != VERSION_FOLDER {
		return nil, fmt.Errorf("resource version is unsupported: %s, expected %s", resource.ResourceVersion, VERSION_FOLDER)
	}
	if resource.ResourceType != RESTYPE_FOLDER {
		return nil, fmt.Errorf("resource is of incorrect type: %s, expected %s", resource.ResourceType, RESTYPE_FOLDER)
	}

	return &Folder{
		Name:     name,
		Resource: resource,
	}, nil
}

func (p *Project) FolderDelete(name string) error {
	for i, f := range p.Data.Folders {
		if f.Name == name {
			p.Data.Folders = append(p.Data.Folders[:i], p.Data.Folders[i+1:]...)
			break
		}
	}

	return nil
}

type ResourceFolder struct {
	ResourceType    ResourceType `json:"resourceType"`
	ResourceVersion Version      `json:"resourceVersion"`
	Name            string       `json:"name"`
	FolderPath      string       `json:"folderPath"`
	Tags            []string     `json:"tags,omitempty"`
}

func NewResourceFolder(name string) *ResourceFolder {
	return &ResourceFolder{
		ResourceType:    RESTYPE_FOLDER,
		ResourceVersion: VERSION_FOLDER,

		Name:       name,
		FolderPath: path.Join(DIR_FOLDER, name+EXT_RESOURCE),
	}
}
