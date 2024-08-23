package yyp

import (
	"fmt"
	"os"
	"path"
)

type IncludedFile struct {
	Name     string
	Path     string
	Data     *[]byte
	Resource *ResourceIncludedFile
}

func NewIncludedFile(filepath, name string, data *[]byte) *IncludedFile {
	return &IncludedFile{
		Name: name,
		Path: filepath,
		Data: data,

		Resource: NewResourceIncludedFile(filepath, name),
	}
}

type ResourceIncludedFile struct {
	ResourceType    ResourceType `json:"resourceType"`
	ResourceVersion Version      `json:"resourceVersion"`
	Name            string       `json:"name"`
	CopyToMask      int          `json:"CopyToMask"`
	FilePath        string       `json:"filePath"`
}

const INCLUDEDFILE_DEFAULTPATH = DIR_DATAFILE

func NewResourceIncludedFile(filepath, name string) *ResourceIncludedFile {
	return &ResourceIncludedFile{
		ResourceType:    RESTYPE_INCLUDEDFILE,
		ResourceVersion: VERSION_INCLUDEDFILE,

		Name:     name,
		FilePath: filepath,

		CopyToMask: -1,
	}
}

func (p *Project) IncludedFileSave(file *IncludedFile) error {
	d := path.Join(p.Path, file.Path)

	_, err := os.Stat(d)
	if err != nil {
		err := os.Mkdir(d, 0o777)
		if err != nil {
			return fmt.Errorf("failed to find and create included files directory: %s, %s", d, err)
		}
	}

	err = os.WriteFile(path.Join(d, file.Name), *file.Data, 0o666)
	if err != nil {
		return fmt.Errorf("failed to write included file: %s", err)
	}

	found := false
	for i, f := range p.Data.IncludedFiles {
		if f.Name == file.Resource.Name && f.FilePath == file.Resource.FilePath {
			p.Data.IncludedFiles[i] = *file.Resource
			found = true
			break
		}
	}

	if !found {
		p.Data.IncludedFiles = append(p.Data.IncludedFiles, *file.Resource)
	}

	return nil
}

func (p *Project) IncludedFileLoad(filepath, name string) (*IncludedFile, error) {
	var resource *ResourceIncludedFile
	for _, f := range p.Data.IncludedFiles {
		if f.Name == name && f.FilePath == filepath {
			resource = &f
			break
		}
	}

	if resource == nil {
		return nil, fmt.Errorf("could not find included file: %s/%s", filepath, name)
	}
	if resource.ResourceVersion != VERSION_INCLUDEDFILE {
		return nil, fmt.Errorf("resource version is unsupported: %s, expected %s", resource.ResourceVersion, VERSION_INCLUDEDFILE)
	}
	if resource.ResourceType != RESTYPE_INCLUDEDFILE {
		return nil, fmt.Errorf("resource is of incorrect type: %s, expected %s", resource.ResourceType, RESTYPE_INCLUDEDFILE)
	}

	datapath := path.Join(p.Path, filepath, name)
	b, err := os.ReadFile(datapath)
	if err != nil {
		return nil, fmt.Errorf("failed to read included file: %s, %s", datapath, err)
	}

	return &IncludedFile{
		Name:     name,
		Path:     filepath,
		Data:     &b,
		Resource: resource,
	}, nil
}

func (p *Project) IncludedFileDelete(filepath, name string) error {
	pth := path.Join(p.Path, filepath, name)

	fs, err := os.Stat(pth)
	if err != nil {
		return nil // Trying to delete smth that doesn't exist shouldn't be an error
	}
	if fs.IsDir() {
		return fmt.Errorf("resource %s is a directory, included file cannot be a directory", pth)
	}

	err = os.Remove(pth)
	if err != nil {
		return fmt.Errorf("failed to delete included file: %s", err)
	}

	for i, f := range p.Data.IncludedFiles {
		if f.Name == name && f.FilePath == filepath {
			p.Data.IncludedFiles = append(p.Data.IncludedFiles[:i], p.Data.IncludedFiles[i+1:]...)
			break
		}
	}

	return nil
}
