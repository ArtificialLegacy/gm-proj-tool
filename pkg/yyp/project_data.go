package yyp

type ProjectData struct {
	ResourceType    string `json:"resourceType"`
	ResourceVersion string `json:"resourceVersion"`
	Name            string `json:"name"`

	MetaData ProjectMetaData `json:"MetaData"`

	DefaultScriptType int  `json:"defaultScriptType"`
	IsEcma            bool `json:"isEcma"`

	Folders         []string `json:"Folders"`
	IncludedFiles   []string `json:"IncludedFiles"`
	LibraryEmitters []string `json:"LibraryEmitters"`

	AudioGroups []ProjectAudioGroup `json:"AudioGroups"`
	Resources   []ProjectResource   `json:"Resources"`
	RoomOrder   []ProjectRoomOrder  `json:"RoomOrderNodes"`
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
	ID struct {
		Name string `json:"name"`
		Path string `json:"path"`
	} `json:"id"`
}

type ProjectRoomOrder struct {
	RoomID struct {
		Name string `json:"name"`
		Path string `json:"path"`
	} `json:"roomId"`
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
