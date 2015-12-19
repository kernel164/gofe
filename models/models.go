package models

type GenericReq struct {
	Params GenericParams `json:"params" binding:"Required"`
}
type GenericParams struct {
	Mode             string `json:"mode" binding:"Required"`
	OnlyFoldreferers bool   `json:"onlyFolders"` // list
	Path             string `json:"path"`        // ALL
	NewPath          string `json:"newPath"`     // move/rename, copy
	Content          string `json:"content"`     // edit
	Name             string `json:"name"`        // addfolder
	Perms            string `json:"perms"`       // changepermissions
	PermsCode        string `json:"permsCode"`   // changepermissions
	Recursive        bool   `json:"recursive"`   // changepermissions
	Destination      string `json:"destination"` // compress, extract
	SourceFile       string `json:"sourceFile"`  // extract
	Preview          bool   `json:"preview"`     // download
}
type ListDirResp struct {
	Result []ListDirEntry `json:"result" binding:"Required"`
}
type ListDirEntry struct {
	Name   string `json:"name" binding:"Required"`
	Rights string `json:"rights" binding:"Required"`
	Size   string `json:"size" binding:"Required"`
	Date   string `json:"date" binding:"Required"`
	Type   string `json:"type" binding:"Required"`
}
type GenericResp struct {
	Result GenericRespBody `json:"result" binding:"Required"`
}
type GenericRespBody struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}
type GetContentResp struct {
	Result string `json:"result" binding:"Required"`
}
