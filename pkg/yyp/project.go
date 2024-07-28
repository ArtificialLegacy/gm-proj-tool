package yyp

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/tailscale/hujson"
)

type Project struct {
	Path string
	Data *ProjectData
}

func NewProject(path string) (*Project, error) {
	fs, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !fs.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", path)
	}

	return &Project{
		Path: path,
		Data: nil,
	}, nil
}

func (p *Project) DataLoad() error {
	fs, err := os.ReadDir(p.Path)
	if err != nil {
		return err
	}

	var projFile string
	for _, f := range fs {
		if path.Ext(f.Name()) == EXT_PROJ {
			projFile = f.Name()
			break
		}
	}

	if projFile == "" {
		return fmt.Errorf("%s file for project was not found", EXT_PROJ)
	}

	b, err := os.ReadFile(path.Join(p.Path, projFile))
	if err != nil {
		return err
	}

	data := ProjectData{}
	b, err = hujson.Standardize(b)
	if err != nil {
		return fmt.Errorf("failed to standardize project data json: %s", err)
	}

	err = json.Unmarshal(b, &data)
	if err != nil {
		return fmt.Errorf("error parsing project data json: %s", err)
	}

	if data.ResourceVersion != RES_VERSION {
		return fmt.Errorf("project resource version is unsupported: %s, expected %s", data.ResourceVersion, RES_VERSION)
	}
	if data.ResourceType != RESTYPE_PROJ {
		return fmt.Errorf("project resource is of wrong type %s, expected %s", data.ResourceType, RESTYPE_PROJ)
	}

	p.Data = &data
	return nil
}

func (p *Project) DataSave() error {
	return nil
}
