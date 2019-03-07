package user_center

import (
	"github.com/iodsp/user_center/models/fionaUserCenter"
	"github.com/iodsp/user_center/params"
	"github.com/jinzhu/gorm"
	"time"
)

type Params struct {
	Id   int    `json:"primary_key"`
	Name string `json:"name" binding:"required"`
}

//given an originals path and a db instance.
type Role struct {
	db *gorm.DB
}

//returns a new Role type with a given path and db instance
func NewRole(db *gorm.DB) *Role {
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
	insertErr := roleDb.Create(&fionaUserCenter.Role{
		Name:      params.Name,
		CreatedAt: createTime,
		UpdatedAt: updateTime,
	}).Error

	return insertErr
}

//find a role by name
func (r *Role) ShowByName(id int) (role fionaUserCenter.Role){
	r.db.Where(&fionaUserCenter.Role{Id: id}).First(&role)
	return role
}

//find a role by id
func (r *Role) Show(id int) (role fionaUserCenter.Role){
	r.db.Where(&fionaUserCenter.Role{Id: id}).First(&role)
	return role
}

//role list
func (r *Role) List() (roles []fionaUserCenter.Role){
	r.db.Model(&fionaUserCenter.Role{}).Find(&roles)
	return roles
}

//update role
func (r *Role) Update(role fionaUserCenter.Role) error {
	updateErr := r.db.Save(role).Error
	return updateErr
}

//delete role
func (r *Role) Delete(role fionaUserCenter.Role) error {
	deleteError := r.db.Delete(&role).Error
	return deleteError
}


