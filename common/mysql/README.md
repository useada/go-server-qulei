# 可参照uploader服务
```
import (
    "a.com/go-server/common/mysql"
)

// 全局初始化调用一次
func init() {
    if err := mysql.Init(Conf.Mysql); err != nil {
        panic(err)
    }
}

// 使用示例
func DemoGetMysqlData(state int) ([]TmpData, error) {
    items := []TmpData{}
    handle := func(orm *gorm.DB) error {
        return orm.Where("state=?", state).Find(&items).Error
    }
    //下面的语句没有指定表名, test是库名，可以在TmpData处指定表名， 参见GORM
    return items, mysql.Master("test").Doit(handle)
}
```
