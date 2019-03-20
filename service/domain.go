package service

import (
	"github.com/iodsp/user_center/context"
	"github.com/iodsp/user_center/models/iodsp"
	"github.com/iodsp/user_center/params"
	"github.com/jinzhu/gorm"
	"time"
)


//given an originals path and a db instance.
type Domain struct {
	db            *gorm.DB
}

//returns a new domain type with a given path and db instance
func NewDomain(conf *context.Config) *Domain {
	db := conf.Db()
	if conf.Debug() {
		db.LogMode(true)
	}
	instance := &Domain{
		db:            db,
	}
	return instance
}

// insert a new domain record
func (r *Domain) StoreDomain(params params.DomainParams) error {
	domainDb := r.db.NewScope(nil).DB()
	createTime := time.Now()
	updateTime := time.Now()
	insertErr := domainDb.Create(&iodsp.Domain{
		Name:      params.Name,
		Type:      params.Type,
		CreatedAt: createTime,
		UpdatedAt: updateTime,
	}).Error
	return insertErr
}

//find a Domain by name
func (r *Domain) ShowDomainByName(name string) (domain iodsp.Domain){
	r.db.Where(&iodsp.Domain{Name: name}).First(&domain)
	return domain
}

//find a Domain by name not id
func (r *Domain) ShowDomainByNameNotId(name string, id int) (domain iodsp.Domain){
	r.db.Where(&iodsp.Domain{Name: name}).Not("id", id).First(&domain)
	return domain
}

//find a Domain by id
func (r *Domain) ShowDomain(id int) (domain iodsp.Domain){
	r.db.Where(&iodsp.Domain{Id: id}).First(&domain)
	return domain
}

//Domain list
func (r *Domain) DomainList() (domains []iodsp.Domain){
	r.db.Model(&iodsp.Domain{}).Order("id desc").Find(&domains)
	return domains
}

//update Domain
func (r *Domain) UpdateDomain(domain iodsp.Domain) error {
	updateErr := r.db.Save(&domain).Error
	return updateErr
}

//delete Domain
func (r *Domain) DeleteDomain(domain iodsp.Domain) error {
	deleteError := r.db.Delete(&domain).Error
	return deleteError
}
