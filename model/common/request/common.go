package request

import (
	"gorm.io/gorm"
)

// PageInfo Paging common input parameter structure
type PageInfo struct {
	Page     int    `json:"page" form:"page"`         // 页码
	PageSize int    `json:"pageSize" form:"pageSize"` // 每页大小
	Keyword  string `json:"keyword" form:"keyword"`   // 关键字
}

func (r *PageInfo) Paginate() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if r.Page <= 0 {
			r.Page = 1
		}
		switch {
		case r.PageSize > 100:
			r.PageSize = 100
		case r.PageSize <= 0:
			r.PageSize = 10
		}
		offset := (r.Page - 1) * r.PageSize
		return db.Offset(offset).Limit(r.PageSize)
	}
}

// GetById Find by id structure
type GetById struct {
	ID int `json:"id" form:"id"` // 主键ID
}

func (r *GetById) Uint() uint {
	return uint(r.ID)
}

type IdsReq struct {
	Ids []int `json:"ids" form:"ids"`
}

// GetAuthorityId Get role by id structure
type GetAuthorityId struct {
	AuthorityId uint `json:"authorityId" form:"authorityId"` // 角色ID
}

type Empty struct{}

// 分页参数
type PageDataOptions struct {
	Sort   string             `json:"sort"`
	Order  string             `json:"order"`
	Page   int                `json:"page"`
	Rows   int                `json:"rows"`
	Export bool               `json:"-"`
	Filter []SearchParameters `json:"-"`
	Wheres string             `json:"wheres"`
	Value  any                `json:"value"`
}

type SearchParameters struct {
	Name        string
	Value       string
	DisplayType string
}

// SaveModel 示例
type SaveModel struct {
	MainData   map[string]any   `json:"mainData"`   // 主表数据
	DetailData []map[string]any `json:"detailData"` // 明细数据列表
	DelKeys    []any            `json:"delKeys"`    // 删除的明细数据主键列表
}
