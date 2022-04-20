package mysql

import (
	"gorm.io/gorm"
)

type UserDefinedModel interface {
	TableName() string
}

// MysqlTable 提供针对mysql数据表的sdk支持
type MysqlTable struct {
	tableName  string
	tableModel UserDefinedModel
	DB         *gorm.DB
}

func newMysqlTable(db *gorm.DB) *MysqlTable {
	return &MysqlTable{DB: db}
}

func (mt *MysqlTable) registerModel(model UserDefinedModel) error {
	mt.tableName = model.TableName()
	mt.tableModel = model
	return nil
}

// SelectWhere 两个参数: queryStr result
// 调用次函数，相当于针对当前table执行 select * from table where `queryStr`，例如`userId = ?`, args = "1"
// 将结果写入result， result类型只能是注册好的model数组，类型为 &[]UserDefinedModel{}
func (mt *MysqlTable) SelectWhere(queryStr string, result interface{}, args ...interface{}) error {
	// todo reflect 检查result类型,只能是 &[]model{}
	// TODO 既然被你看到了就由你来完善吧

	if err := mt.DB.Table(mt.tableModel.TableName()).Where(queryStr, args...).Find(result).Error; err != nil {
		return err
	}
	return nil
}

// Insert 一个参数 toInsertLines
// 调用次函数，相当于针对当前table，插入toInsertLines 对应的数据
// toInsertLines类型为 UserDefinedModel
func (mt *MysqlTable) Insert(toInsertLines UserDefinedModel) error {
	// todo reflect检查toInserLines类型，是数组则开多次插入
	// TODO 既然被你看到了就由你来完善吧

	if err := mt.DB.Table(mt.tableModel.TableName()).Create(toInsertLines).Error; err != nil {
		return err
	}
	return nil
}

// Update 三个参数：queryStr、 field、target
// 调用此函数，相当于针对queryStr 筛选出的数据条目(例如queryStr = 'userId = ?', args = "1") ，将筛选出的数据条目的field字段替换为target内容
func (mt *MysqlTable) Update(queryStr, field string, target interface{}, args ...interface{}) error {
	// todo reflect检查target 类型，是否与注册好的的field相符
	// TODO 既然被你看到了就由你来完善吧

	if err := mt.DB.Table(mt.tableModel.TableName()).Where(queryStr, args...).Update(field, target).Error; err != nil {
		return err
	}
	return nil
}

// Delete 一个参数：toDeleteTarget
// 传入一个UserDefinedModel类型，如果此对象的userId = 1,则删除掉数据库中userId= 1的字段
func (mt *MysqlTable) Delete(toDeleteTarget UserDefinedModel) error {
	// todo reflect检查toDeleteTarget 类型，确保所有字段不为空
	// TODO 既然被你看到了就由你来完善吧

	if err := mt.DB.Table(mt.tableModel.TableName()).Delete(toDeleteTarget).Error; err != nil {
		return err
	}
	return nil
}

// First 两个参数： queryStr、findTarget
// queryStr为筛选用的query，例如`userId = ?`, args = "1", findTarget为 UserDefinedModel 类型指针，为第一个找到的数据。
func (mt *MysqlTable) First(queryStr string, findTarget UserDefinedModel, args ...interface{}) error {
	// todo reflect 检查findTarget类型 确保是注册类型的指针相同
	// TODO 既然被你看到了就由你来完善吧

	if err := mt.DB.Table(mt.tableModel.TableName()).Where(queryStr, args...).Find(findTarget).Error; err != nil {
		return err
	}
	return nil
}
