package yyp

type ProjectData struct {
	ResourceType    ResourceType `json:"resourceType"`
	ResourceVersion Version      `json:"resourceVersion"`
	Name            string       `json:"name"`

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

	TemplateType TemplateType `json:"templateType"`

	TextureGroups []ProjectTextureGroups `json:"TextureGroups"`
}

type ProjectMetaData struct {
	IDEVersion string `json:"IDEVersion"`
}

type ProjectAudioGroup struct {
	ResourceType    ResourceType `json:"resourceType"`
	ResourceVersion Version      `json:"resourceVersion"`
	Name            string       `json:"name"`
	Targets         int          `json:"targets"`
}

type ProjectResource struct {
	ID ProjectResourceNode `json:"id"`
}

type ProjectRoomOrder struct {
	RoomID ProjectResourceNode `json:"roomId"`
}

type ProjectTextureGroups struct {
	ResourceType    ResourceType            `json:"resourceType"`
	ResourceVersion Version                 `json:"resourceVersion"`
	Name            string                  `json:"name"`
	Autocrop        bool                    `json:"autocrop"`
	Border          int                     `json:"border"`
	CompressFormat  TextureGroupCompression `json:"compressFormat"`
	Directory       string                  `json:"directory"`
	GroupParent     ProjectResourceNode     `json:"groupParent"`
	IsScaled        bool                    `json:"isScaled"`
	LoadType        TextureGroupType        `json:"loadType"`
	MipsToGenerate  int                     `json:"mipsToGenerate"`
	Targets         int                     `json:"targets"`
}

func ProjectTextureGroupDefaultID() ProjectResourceNode {
	return ProjectResourceNode{
		Name: "Default",
		Path: "texturegroups/Default",
	}
}

type ProjectFolder struct {
	ResourceType    ResourceType `json:"resourceType"`
	ResourceVersion Version      `json:"resourceVersion"`
	Name            string       `json:"name"`
	FolderPath      string       `json:"folderPath"`
	Tags            []string     `json:"tags,omitempty"`
}

type ProjectIncludedFile struct {
	ResourceType    ResourceType `json:"resourceType"`
	ResourceVersion Version      `json:"resourceVersion"`
	Name            string       `json:"name"`
	CopyToMask      int          `json:"CopyToMask"`
	FilePath        string       `json:"filePath"`
}

type ProjectConfig struct {
	Children []ProjectConfig `json:"children"`
	Name     string          `json:"name"`
}

type ProjectLibraryEmitters struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type ProjectEmpty struct{}