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
func (r *Resource) Show(id int) (resource iodsp.Resource) {
	r.db.Where(&iodsp.Resource{Id: id}).First(&resource)
	return resource
}

//find a Resource by url
func (r *Resource) ShowByUrl(url string)(resource iodsp.Resource) {
	r.db.Where(&iodsp.Resource{Url: url}).First(&resource)
	return resource
}

//find a Resource by name
func (r *Resource) ShowByName(name string)(resource iodsp.Resource) {
	r.db.Where(&iodsp.Resource{Name: name}).First(&resource)
	return resource
}

//Resource list
func (r *Resource) List() (resources []iodsp.Resource) {
	r.db.Model(&iodsp.Resource{}).Order("id desc").Find(&resources)
	return resources
}

//update Resource
func (r *Resource) Update(resource iodsp.Resource) error {
	updateErr := r.db.Save(&resource).Error
	return updateErr
}

//delete Resource
func (r *Resource) Delete(resource iodsp.Resource) error {
	deleteError := r.db.Delete(&resource).Error
	return deleteError
}

//find a Resource by url not id
func (r *Resource) ShowResourceByUrlNotId(url string, id int) (resource iodsp.Resource){
	r.db.Where(&iodsp.Resource{Url: url}).Not("id", id).First(&resource)
	return resource
}

//find a Resource by name not id
func (r *Resource) ShowResourceByNameNotId(name string, id int) (resource iodsp.Resource){
	r.db.Where(&iodsp.Resource{Name: name}).Not("id", id).First(&resource)
	return resource
}
