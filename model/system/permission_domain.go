package system

// 菜单权限
type MenuAction struct {
	MenuId  int      `json:"menuId"`                                                                // 菜单ID
	Actions []string `json:"actions" example:"Search,Add,Delete,Update,Import,Export,Upload,Audit"` // 菜单权限
}

type Permission struct {
	Menu_Id     int      `json:"menu_Id"`
	ParentId    int      `json:"parentId"` //父菜单id
	TableName   string   `json:"tableName"`
	MenuAuth    string   `json:"menuAuth"`
	UserAuth    string   `json:"userAuth"`
	UserAuthArr []string `json:"userAuthArr" example:"Add,Search"`
	MenuType    int      `json:"menuType"`
}

type Action struct {
	Action_Id int    `json:"action_Id"`
	Menu_Id   int    `json:"menu_Id"` //父菜单id
	Text      string `json:"text"`
	Value     string `json:"value"`
}

type TreeMenu struct {
	ID         int      `json:"id"` //菜单id
	Name       string   `json:"name" example:"菜单名称"`
	Url        string   `json:"url" example:"/FTQOperationStatistics"`
	ParentId   int      `json:"parentId"` //父菜单id
	Icon       string   `json:"icon" example:"图标icon"`
	Enable     int      `json:"enable"`
	TableName  string   `json:"tableName" example:"表/视图名称"`
	Permission []string `json:"permission" example:"Search,Export"`
}
