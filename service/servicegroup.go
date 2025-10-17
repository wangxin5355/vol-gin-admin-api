package service

import (
	"fmt"
	"github.com/wangxin5355/vol-gin-admin-api/service/example"
	"github.com/wangxin5355/vol-gin-admin-api/service/system"
	"github.com/wangxin5355/vol-gin-admin-api/service/test"
	"sync"
)

// 注意，如果在一个服务中想应用另外一个服务，使用ServiceInstance，会导致循环依赖，
// 想要实现在一个服务中引用另一个服务，请直接在服务中获取或者创建对应服务的实例 请参照MenuService-》PermissionService
type ServiceInstance struct {
	ExampleTestService     *example.ExampleTestService
	TestService            *test.TestService
	OperationRecordService *system.OperationRecordService
	UserService            *system.UserService
	SystemConfigService    *system.SystemConfigService
	JwtService             *system.JwtService
	CasbinService          *system.CasbinService
	MenuService            *system.MenuService
	PermissionService      *system.PermissionService //依赖实现
	//PermissionService system.IPermissions //依赖接口，两种方式都可以,如果是大型项目，建议使用接口依赖
	TableInfoService  *system.TableInfoService
	DictionaryService *system.DictionaryService
}

var (
	ServiceInstances *ServiceInstance
	once             sync.Once // 确保初始化代码只执行一次
)

// 初始化服务层实例

func InitServiceInstance() {
	once.Do(func() {
		ServiceInstances = &ServiceInstance{
			TestService:        test.NewTestService(),         //提供方法
			ExampleTestService: &example.ExampleTestService{}, //直接实例化方式，
			PermissionService:  system.GetPermissionService(), //返回单例
			//正常来说全局就一个service就行了。全部直接实例化方式初始化就行
			OperationRecordService: &system.OperationRecordService{},
			UserService:            &system.UserService{},
			SystemConfigService:    &system.SystemConfigService{},
			JwtService:             &system.JwtService{},
			MenuService:            &system.MenuService{},
			TableInfoService:       &system.TableInfoService{},
			DictionaryService:      &system.DictionaryService{},
		}
		fmt.Println("ServiceInstances 单例已初始化")
	})
}
