package yyp

type ProjectOrder struct {
	FolderOrderSettings []ProjectFolderOrderSettings `json:"FolderOrderSettings"`

	ResourceOrderSettings []ProjectResourceOrderSettings `json:"ResourceOrderSettings"`
}

type ProjectFolderOrderSettings struct {
	Name  string `json:"name"`
	Order int    `json:"order"`
	Path  string `json:"path"`
}

type ProjectResourceOrderSettings struct {
	Name  string `json:"name"`
	Order int    `json:"order"`
	Path  string `json:"path"`
}
