package system

import (
	"github.com/wangxin5355/vol-gin-admin-api/global"
	"log"
)

type DictionaryService struct {
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
