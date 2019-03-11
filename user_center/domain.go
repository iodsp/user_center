package user_center

import (
	"github.com/iodsp/user_center/models/fionaUserCenter"
	"github.com/iodsp/user_center/params"
	"github.com/jinzhu/gorm"
	"time"
)


//given an originals path and a db instance.
type Domain struct {
	db            *gorm.DB
}

//returns a new domain type with a given path and db instance
func NewDomain(db *gorm.DB, isDebug bool) *Domain {
	if isDebug {
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
	insertErr := domainDb.Create(&fionaUserCenter.Domain{
		Name:      params.Name,
		Type:      params.Type,
		CreatedAt: createTime,
		UpdatedAt: updateTime,
	}).Error
	return insertErr
}

//find a Domain by name
func (r *Domain) ShowDomainByName(name string) (Domain fionaUserCenter.Domain){
	r.db.Where(&fionaUserCenter.Domain{Name: name}).First(&Domain)
	return Domain
}

//find a Domain by name not id
func (r *Domain) ShowDomainByNameNotId(name string, id int) (Domain fionaUserCenter.Domain){
	r.db.Where(&fionaUserCenter.Domain{Name: name}).Not("id", id).First(&Domain)
	return Domain
}

//find a Domain by id
func (r *Domain) ShowDomain(id int) (Domain fionaUserCenter.Domain){
	r.db.Where(&fionaUserCenter.Domain{Id: id}).First(&Domain)
	return Domain
}

//Domain list
func (r *Domain) DomainList() (Domains []fionaUserCenter.Domain){
	r.db.Model(&fionaUserCenter.Domain{}).Order("id desc").Find(&Domains)
	return Domains
}

//update Domain
func (r *Domain) UpdateDomain(Domain fionaUserCenter.Domain) error {
	updateErr := r.db.Save(Domain).Error
	return updateErr
}

//delete Domain
func (r *Domain) DeleteDomain(Domain fionaUserCenter.Domain) error {
	deleteError := r.db.Delete(&Domain).Error
	return deleteError
}
