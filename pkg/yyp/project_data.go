package yyp

type ProjectData struct {
	ResourceType    string `json:"resourceType"`
	ResourceVersion string `json:"resourceVersion"`
	Name            string `json:"name"`

	AudioGroups []ProjectAudioGroup `json:"AudioGroups"`

	Configs           ProjectConfig `json:"configs"`
	DefaultScriptType int           `json:"defaultScriptType"`

	Folders       []ProjectFolder       `json:"Folders"`
	IncludedFiles []ProjectIncludedFile `json:"IncludedFiles"`

	IsEcma bool `json:"isEcma"`

	LibraryEmitters []ProjectLibraryEmitters `json:"LibraryEmitters"`

	MetaData ProjectMetaData `json:"MetaData"`

	Resources []ProjectResource  `json:"resources"`
	RoomOrder []ProjectRoomOrder `json:"RoomOrderNodes"`

	TemplateType string `json:"templateType"`

	TextureGroups []ProjectTextureGroups `json:"TextureGroups"`
}

type ProjectMetaData struct {
	IDEVersion string `json:"IDEVersion"`
}

type ProjectAudioGroup struct {
	ResourceType    string `json:"resourceType"`
	ResourceVersion string `json:"resourceVersion"`
	Name            string `json:"name"`
	Targets         int    `json:"targets"`
}

type ProjectResource struct {
	ID ProjectResourceNode `json:"id"`
}

type ProjectResourceNode struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type ProjectRoomOrder struct {
	RoomID ProjectRoomOrderNode `json:"roomId"`
}

type ProjectRoomOrderNode struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type ProjectTextureGroups struct {
	ResourceType    string             `json:"resourceType"`
	ResourceVersion string             `json:"resourceVersion"`
	Name            string             `json:"name"`
	Autocrop        bool               `json:"autocrop"`
	Border          int                `json:"border"`
	CompressFormat  string             `json:"compressFormat"`
	Directory       string             `json:"directory"`
	GroupParent     TextureGroupParent `json:"groupParent"`
	IsScaled        bool               `json:"isScaled"`
	LoadType        string             `json:"loadType"`
	MipsToGenerate  int                `json:"mipsToGenerate"`
	Targets         int                `json:"targets"`
}

type TextureGroupParent struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type ProjectFolder struct {
	ResourceType    string   `json:"resourceType"`
	ResourceVersion string   `json:"resourceVersion"`
	Name            string   `json:"name"`
	FolderPath      string   `json:"folderPath"`
	Tags            []string `json:"tags,omitempty"`
}

type ProjectIncludedFile struct {
	ResourceType    string `json:"resourceType"`
	ResourceVersion string `json:"resourceVersion"`
	Name            string `json:"name"`
	CopyToMask      int    `json:"CopyToMask"`
	FilePath        string `json:"filePath"`
}

type ProjectConfig struct {
	Children []ProjectConfig `json:"children"`
	Name     string          `json:"name"`
}

type ProjectLibraryEmitters struct {
	Name string `json:"name"`
	Path string `json:"path"`
}
