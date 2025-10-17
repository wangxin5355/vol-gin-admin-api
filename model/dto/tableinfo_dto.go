package dto

type TableInfo struct {
	Id       int    `json:"id"`
	PId      int    `json:"pId"` //父id
	ParentId int    `json:"parentId"`
	Name     string `json:"name"`
	OrderNo  int    `json:"orderNo"`
}

type TableTreeListData struct {
	Id       int    `json:"id"`
	PId      int    `json:"pId"` //父id
	ParentId int    `json:"parentId"`
	Name     string `json:"name"`
	IsParent bool   `json:"isParent"`
}
