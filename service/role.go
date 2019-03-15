package service

import (
	"github.com/iodsp/user_center/context"
	"github.com/iodsp/user_center/models/iodsp"
	"github.com/iodsp/user_center/params"
	"github.com/jinzhu/gorm"
	"time"
)

//given an originals path and a db instance.
type Role struct {
	db *gorm.DB
}

//returns a new Role type with a given path and db instance
func NewRole(conf *context.Config) *Role {
	db := conf.Db()
	if conf.Debug() {
		db.LogMode(true)
	}
	instance := &Role{
		db: db,
	}
	return instance
}

// insert a new role record
func (r *Role) Store(params params.RoleParams) error {
	roleDb := r.db.NewScope(nil).DB()

	createTime := time.Now()
	updateTime := time.Now()
	insertErr := roleDb.Create(&iodsp.Role{
		Name:      params.Name,
		CreatedAt: createTime,
		UpdatedAt: updateTime,
	}).Error

	return insertErr
}

//find a role by name
func (r *Role) ShowByName(name string) (role iodsp.Role){
	r.db.Where(&iodsp.Role{Name: name}).First(&role)
	return role
}

//find a role by name not id
func (r *Role) UpdateShowByName(name string, id int) (role iodsp.Role){
	r.db.Where(&iodsp.Role{Name: name}).Not("id", id).First(&role)
	return role
}

//find a role by id
func (r *Role) Show(id int) (role iodsp.Role){
	r.db.Where(&iodsp.Role{Id: id}).First(&role)
	return role
}

//role list
func (r *Role) List() (roles []iodsp.Role){
	r.db.Model(&iodsp.Role{}).Order("id desc").Find(&roles)
	return roles
}

//update role
func (r *Role) Update(role iodsp.Role) error {
	updateErr := r.db.Save(role).Error
	return updateErr
}

//delete role
func (r *Role) Delete(role iodsp.Role) error {
	deleteError := r.db.Delete(&role).Error
	return deleteError
}


