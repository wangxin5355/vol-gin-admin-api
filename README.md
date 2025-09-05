# vol-gin-admin-api
vol快速开发框架后端api  基于gin+gorm的实现
# 安装依赖
go generate
go mod tidy
# 编译
go build
# 允许
go run main.go
# swagger文档生成
swag init

---生成文件规划-----
api层
1、需要生成一个系统默认api，支持增删改查，默认api不可修改。
2、需要生成一个用户自定义api，用户可自行添加接口
自定义api
type AccountApi struct{
     DefaultAccountApi
}
默认api
type DefaultAccountApi struct{}
server层
1、需要生成一个默认server,支持默认增删改查，默认server不可修改
2、需要生成一个用户自定义server,用户可自行添加方法
router层
1、需要生成一个默认router,支持默认api接口，默认router不可修改
2、需要生成一个用户自定义router,用户可自行添加自定义api路由