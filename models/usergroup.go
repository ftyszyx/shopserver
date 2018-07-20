package models

type UserGroup struct {
	Model
}

type UserGroupData struct {
	Name        string `empty:"组名不能为空"`
	Module_ids  string
	expire_time int
	Group_type  int
}

func (self *UserGroup) GetModelStruct() interface{} {
	return UserGroupData{}
}
