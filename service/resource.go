package service

import (
	"github.com/iodsp/user_center/context"
	"github.com/iodsp/user_center/models/iodsp"
	"github.com/iodsp/user_center/params"
	"github.com/jinzhu/gorm"
	"time"
)

//given an originals path and a db instance.
type Resource struct {
	db *gorm.DB
}

//returns a new Resource type with a given path and db instance
func NewResource(conf *context.Config) *Resource {
	db := conf.Db()
	if conf.Debug() {
		db.LogMode(true)
	}
	instance := &Resource{
		db: db,
	}
	return instance
}

// insert a new Resource record
func (r *Resource) StoreResource(params params.ResourceParams) error {
	ResourceDb := r.db.NewScope(nil).DB()

	createTime := time.Now()
	updateTime := time.Now()
	insertErr := ResourceDb.Create(&iodsp.Resource{
		DomainId:   params.DomainId,
		DomainName: params.DomainName,
		Name:       params.Name,
		Url:        params.Url,
		Desc:       params.Desc,
		CreatedAt:  createTime,
		UpdatedAt:  updateTime,
	}).Error

	return insertErr
}

//find a Resource by id
func (r *Resource) Show(id int) (Resource iodsp.Resource) {
	r.db.Where(&iodsp.Resource{Id: id}).First(&Resource)
	return Resource
}

//Resource list
func (r *Resource) List() (Resources []iodsp.Resource) {
	r.db.Model(&iodsp.Resource{}).Order("id desc").Find(&Resources)
	return Resources
}

//update Resource
func (r *Resource) Update(Resource iodsp.Resource) error {
	updateErr := r.db.Save(Resource).Error
	return updateErr
}

//delete Resource
func (r *Resource) Delete(Resource iodsp.Resource) error {
	deleteError := r.db.Delete(&Resource).Error
	return deleteError
}
