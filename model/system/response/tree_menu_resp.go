package response

// 首页菜单树响应结构体
type TreeMenuResp struct {
	ID         int      `json:"id"` //菜单id
	Name       string   `json:"name" example:"菜单名称"`
	Url        string   `json:"url" example:"/FTQOperationStatistics"`
	ParentId   int      `json:"parentId"` //父菜单id
	Icon       string   `json:"icon" example:"图标icon"`
	Enable     int      `json:"enable"`
	TableName  int      `json:"tableName" example:"表/视图名称"`
	Permission []string `json:"permission" example:"Search,Export"`
}
