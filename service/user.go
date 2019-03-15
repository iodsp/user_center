package service

import (
	"github.com/iodsp/user_center/context"
	"github.com/iodsp/user_center/models/iodsp"
	"github.com/iodsp/user_center/params"
	"github.com/jinzhu/gorm"
	"time"
)

//given an originals path and a db instance.
type User struct {
	db *gorm.DB
}

//returns a new User type with a given path and db instance
func NewUser(conf *context.Config) *User {
	db := conf.Db()
	if conf.Debug() {
		db.LogMode(true)
	}
	instance := &User{
		db: db,
	}
	return instance
}

// insert a new user record
func (r *User) StoreUser(params params.UserParams) error {
	userDb := r.db.NewScope(nil).DB()

	createTime := time.Now()
	updateTime := time.Now()
	insertErr := userDb.Create(&iodsp.User{
		Name:      params.Name,
		Phone:            params.Phone,
		RegisterPlatform: params.RegisterPlatform,
		DomainName:       params.DomainName,
		Type :            params.Type,
		IdentityId:       params.IdentityId,
		DomainId:         params.DomainId,
		CreatedAt: createTime,
		UpdatedAt: updateTime,
	}).Error

	return insertErr
}

//find a User by name
func (r *User) ShowByName(name string) (User iodsp.User){
	r.db.Where(&iodsp.User{Name: name}).First(&User)
	return User
}

//find a User by phone
func (r *User) ShowByPhone(phone string) (User iodsp.User){
	r.db.Where(&iodsp.User{Phone: phone}).First(&User)
	return User
}

//find a User by phone not id
func (r *User) UpdateShowByPhone(phone string, id int) (User iodsp.User){
	r.db.Where(&iodsp.User{Phone: phone}).Not("id", id).First(&User)
	return User
}

//find a User by name not id
func (r *User) UpdateShowByName(name string, id int) (User iodsp.User){
	r.db.Where(&iodsp.User{Name: name}).Not("id", id).First(&User)
	return User
}

//find a User by id
func (r *User) Show(id int) (User iodsp.User){
	r.db.Where(&iodsp.User{Id: id}).First(&User)
	return User
}

//User list
func (r *User) List() (Users []iodsp.User){
	r.db.Model(&iodsp.User{}).Order("id desc").Find(&Users)
	return Users
}

//update User
func (r *User) Update(User iodsp.User) error {
	updateErr := r.db.Save(User).Error
	return updateErr
}

//delete User
func (r *User) Delete(User iodsp.User) error {
	deleteError := r.db.Delete(&User).Error
	return deleteError
}
