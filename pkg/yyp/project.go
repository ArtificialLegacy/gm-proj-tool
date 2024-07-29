package yyp

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/tailscale/hujson"
)

type Project struct {
	Path          string
	Data          *ProjectData
	ResourceOrder *ProjectOrder
}

type ResourceParent struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func NewProject(path string) (*Project, error) {
	fs, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !fs.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", path)
	}

	proj := &Project{
		Path:          path,
		Data:          nil,
		ResourceOrder: nil,
	}

	err = proj.DataLoad()
	if err != nil {
		return nil, err
	}
	err = proj.OrderLoad()
	if err != nil {
		return nil, err
	}

	return proj, nil
}

func (p *Project) AsParent() ResourceParent {
	if p.Data == nil {
		return ResourceParent{}
	}

	return ResourceParent{
		Name: p.Data.Name,
		Path: p.Data.Name + EXT_PROJ,
	}
}

func (p *Project) DataLoad() error {
	projFile := findFile(p.Path, EXT_PROJ)
	if projFile == "" {
		return fmt.Errorf("%s file for project was not found", EXT_PROJ)
	}

	data := &ProjectData{}
	err := loadJSON(path.Join(p.Path, projFile), data)
	if err != nil {
		return fmt.Errorf("error parsing project data json: %s", err)
	}

	if data.ResourceVersion != VERSION_PROJ {
		return fmt.Errorf("project resource version is unsupported: %s, expected %s", data.ResourceVersion, VERSION_PROJ)
	}
	if data.ResourceType != RESTYPE_PROJ {
		return fmt.Errorf("project resource is of wrong type %s, expected %s", data.ResourceType, RESTYPE_PROJ)
	}

	p.Data = data
	return nil
}

func (p *Project) DataSave() error {
	if p.Data == nil {
		return fmt.Errorf("project data is nil, nothing to save")
	}

	err := saveJSON(path.Join(p.Path, p.Data.Name+EXT_PROJ), p.Data)
	if err != nil {
		return err
	}

	return nil
}

func (p *Project) OrderLoad() error {
	orderFile := findFile(p.Path, EXT_ORDER)
	if orderFile == "" {
		return fmt.Errorf("%s file for project was not found", EXT_ORDER)
	}

	order := &ProjectOrder{}
	err := loadJSON(path.Join(p.Path, orderFile), order)
	if err != nil {
		return fmt.Errorf("error parsing project order json: %s", err)
	}

	p.ResourceOrder = order
	return nil
}

func (p *Project) OrderSave() error {
	if p.ResourceOrder == nil {
		return fmt.Errorf("project order data is nil, nothing to save")
	}

	err := saveJSON(path.Join(p.Path, p.Data.Name+EXT_ORDER), p.ResourceOrder)
	if err != nil {
		return err
	}

	return nil
}

func (p *Project) ResourceList(dir string) []string {
	list := []string{}

	fs, err := os.ReadDir(path.Join(p.Path, dir))
	if err != nil {
		return list
	}

	for _, f := range fs {
		if f.IsDir() {
			list = append(list, f.Name())
		}
	}

	return list
}

type ImportableResource interface {
	Save(pdir string) (string, string, *ResourceParent, error)
}

func (p *Project) ImportResource(res ImportableResource) error {
	if p.Data == nil {
		return fmt.Errorf("project data must not be nil")
	}
	if p.ResourceOrder == nil {
		return fmt.Errorf("project resource order must not be nil")
	}

	name, d, _, err := res.Save(p.Path)
	if err != nil {
		return fmt.Errorf("failed to save resource: %s", err)
	}

	p.Data.Resources = append(p.Data.Resources, ProjectResource{
		ID: ProjectResourceNode{
			Name: name,
			Path: d,
		},
	})

	p.ResourceOrder.ResourceOrderSettings = append(p.ResourceOrder.ResourceOrderSettings, ProjectResourceOrderSettings{
		Name:  name,
		Order: 0,
		Path:  d,
	})

	return nil
}

func findFile(p, ext string) string {
	fs, err := os.ReadDir(p)
	if err != nil {
		return ""
	}

	var file string
	for _, f := range fs {
		if path.Ext(f.Name()) == ext {
			file = f.Name()
			break
		}
	}

	return file
}

func loadJSON(p string, v any) error {
	b, err := os.ReadFile(p)
	if err != nil {
		return err
	}

	b, err = hujson.Standardize(b) // GM uses non-standard json
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	return nil
}

func saveJSON(p string, v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(p, b, 0o666)
	if err != nil {
		return err
	}

	return nil
}
