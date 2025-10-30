package system

import (
	"log"

	"github.com/wangxin5355/vol-gin-admin-api/core/base"
	"github.com/wangxin5355/vol-gin-admin-api/core/initialize"
	"github.com/wangxin5355/vol-gin-admin-api/global"
	"github.com/wangxin5355/vol-gin-admin-api/model/common/request"
	"github.com/wangxin5355/vol-gin-admin-api/model/common/response"
	"github.com/wangxin5355/vol-gin-admin-api/model/dto"
	"github.com/wangxin5355/vol-gin-admin-api/model/system"
	"github.com/wangxin5355/vol-gin-admin-api/model/system/partial"
	"github.com/wangxin5355/vol-gin-admin-api/utils"
)

func InitDictionaryService() *DictionaryService {
	return &DictionaryService{
		BaseService: base.InitBaseService[partial.SysDictionaryEntity, system.SysDictionary](string(initialize.DbGin)),
	}
}

type DictionaryService struct {
	*base.BaseService[partial.SysDictionaryEntity, system.SysDictionary]
}

func (s *DictionaryService) GetPageData(options request.PageDataOptions) *response.PageGridData[partial.SysDictionaryEntity] {
	return s.BaseService.GetPageData(options)
}

func (dictionaryService *DictionaryService) GetBuilderDictionary() []string {
	var dicNos []string
	result := global.GVA_DB.Raw("SELECT DicNo FROM `sys_dictionary`").Scan(&dicNos)
	if result.Error != nil {
		log.Fatal(result.Error)
		return make([]string, 0)
	}
	return dicNos
}

// GetVueDictionar 获取vue页面需要的字典
func (dictionaryService *DictionaryService) GetVueDictionar(dicNos []string) []map[string]interface{} {
	if len(dicNos) == 0 {
		return []map[string]interface{}{}
	}

	type dicConfigRow struct {
		DicID    uint
		DicNo    string
		Config   string
		DbSql    string
		DBServer string
	}

	var dicConfigs []dicConfigRow
	result := global.GVA_DB.Raw("SELECT Dic_ID as DicID ,DicNo, Config, DbSql, DBServer FROM sys_dictionary WHERE DicNo IN (?)", dicNos).Scan(&dicConfigs)
	if result.Error != nil {
		return []map[string]interface{}{}
	}

	var out []map[string]interface{}
	for _, cfg := range dicConfigs {
		var list []dto.DictDetail

		if utils.IsNull(cfg.DbSql) {
			res := global.GVA_DB.Raw("SELECT DicValue as `Key`,DicName as `Value`  FROM `sys_dictionarylist` WHERE Dic_ID = ? and (Enable is null or Enable != 0) ORDER BY OrderNo", cfg.DicID).Scan(&list)
			if res.Error != nil {
				list = []dto.DictDetail{}
			}
		} else {
			//查询前做一下自定义sql参数替换
			sql := ReplaceGetCustomSqlParams(cfg.DicNo, cfg.DbSql)
			//执行自定义sql
			list = GetCustomDBSql(sql, cfg.DBServer)
		}
		item := map[string]interface{}{
			"dicNo":  cfg.DicNo,
			"config": cfg.Config,
			"data":   list,
		}
		out = append(out, item)
	}
	return out
}

// GetCustomDBSql 执行自定义sql
func GetCustomDBSql(sql string, dbServer string) []dto.DictDetail {
	var list []dto.DictDetail
	if utils.SqlInjectCheck(sql) {
		return list
	}

	//获取对应的db查询
	db := global.MustGetGlobalDBByDBName(dbServer)
	if db == nil {
		return list
	}

	result := db.Raw(sql).Scan(&list)
	if result.Error != nil {
		return list
	}
	return list
}

// ReplaceGetCustomSqlParams 自定义sql参数替换
func ReplaceGetCustomSqlParams(dicNo, sql string) string {
	if utils.IsNull(sql) {
		return sql
	}
	switch dicNo {
	case "roles":
	case "t_roles":
	case "tree_roles":
		sql = GetTree_RolesSql(sql)
		break
	}
	return sql
}
func GetTree_RolesSql(sql string) string {
	//TODO:获取当前用户的角色ID列表然后替换sql,超管应该要无视这个条件直接返回
	return sql
}
