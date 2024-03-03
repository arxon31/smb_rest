package entity

type DirGetRequest struct {
	Dirs map[string]FileNode `json:"path"`
}

func (d DirGetRequest) Validate() error {
	return nil
}

type DirGetResponse struct {
	DirPath string `json:"path"`
}

func (d DirGetResponse) Validate() error {

	return nil
}
