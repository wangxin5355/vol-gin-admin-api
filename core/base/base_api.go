package base

type BaseApi[T any, S any] struct {
	ServiceGroup *S
}

// Service 子类必须重写此方法
func (b *BaseApi[T, S]) Service() *BaseApi[T, S] {
	return &BaseApi[T, S]{
		ServiceGroup: b.ServiceGroup,
	}
}

// GetPageData 通用分页接口
//func (b *BaseApi[T, S]) GetPageData(c *gin.Context) {
//	param, err := utils.BindJsonToPageDataOptions(c)
//	if err != nil {
//		response.FailWithMessage(err.Error(), c)
//		return
//	}
//	//data := b.Service().GetPageData(param)
//	data := b.ServiceGroup
//	response.OkWithData(data, c)
//}
