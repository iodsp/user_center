package service

import (
	"github.com/iodsp/user_center/context"
	"github.com/iodsp/user_center/models/iodsp"
	"github.com/iodsp/user_center/params"
	"github.com/jinzhu/gorm"
	"time"
)

//given an originals path and a db instance.
type Grant struct {
	db *gorm.DB
}

//returns a new grant type with a given path and db instance
func NewGrant(conf *context.Config) *Grant {
	db := conf.Db()
	if conf.Debug() {
		db.LogMode(true)
	}
	instance := &Grant{
		db: db,
	}
	return instance
}

func (r *Grant) AddUserRole(params params.UserRoleParam) error {
	domainDb := r.db.NewScope(nil).DB()
	createTime := time.Now()
	updateTime := time.Now()
	insertErr := domainDb.Create(&iodsp.UserRole{
		UserId:     params.UserId,
		RoleId:     params.RoleId,
		RoleName:   params.RoleName,
		DomainId:   params.DomainId,
		DomainName: params.DomainName,
		CreatedAt:  createTime,
		UpdatedAt:  updateTime,
	}).Error
	return insertErr
}

func (r *Grant) DeleteUserRole(userRole iodsp.UserRole) error {
	deleteError := r.db.Delete(&userRole).Error
	return deleteError
}

func (r *Grant) AddRoleResource(params params.RoleResource) error {
	domainDb := r.db.NewScope(nil).DB()
	createTime := time.Now()
	updateTime := time.Now()
	insertErr := domainDb.Create(&iodsp.RoleResource{
		ResourceId:  params.ResourceId,
		RoleId:      params.RoleId,
		ResourceUrl: params.ResourceUrl,
		CreatedAt:   createTime,
		UpdatedAt:   updateTime,
	}).Error
	return insertErr
}

func (r *Grant) DeleteRoleResource(roleResource iodsp.RoleResource) error {
	deleteError := r.db.Delete(&roleResource).Error
	return deleteError
}

//roleIds has resource
func (r *Grant) HasRoleResource(roleId string, resourceUrl string) (roleResource iodsp.RoleResource) {
	r.db.Where(&iodsp.RoleResource{ResourceUrl: resourceUrl}).Where("roleId in (?)", roleId).First(&roleResource)
	return roleResource
}

//select roleIds by userId
func (r *Grant) ShowRoleIdsByUserId(userId int) (userRoles []iodsp.UserRole) {
	r.db.Select([]string{"roleId"}).Where(map[string]interface{}{"UserId": userId}).Find(&userRoles)
	return userRoles
}

//select userRole by userId roleId
func (r *Grant) ShowByUserIdRoleId(userId int, roleId int) (userRole iodsp.UserRole) {
	r.db.Where(&iodsp.UserRole{UserId: userId, RoleId: roleId}).Find(&userRole)
	return userRole
}

//select roleResource by roleId resourceId
func (r *Grant) ShowByRoleIdResourceId(roleId int, resourceId int) (roleResource iodsp.RoleResource) {
	r.db.Where(&iodsp.RoleResource{ResourceId: resourceId, RoleId: roleId}).Find(&roleResource)
	return roleResource
}

//update userRole
func (r *Grant) UpdateUserRole(updateUserRole iodsp.UserRole) error {
	updateErr := r.db.Unscoped().Save(&updateUserRole).Error
	return updateErr
}

//select a userRole by userId roleId include those having been softly deleted
func (r *Grant) ShowByUserIdRoleIdIncDel(userId int, roleId int) (userRole iodsp.UserRole) {
	r.db.Unscoped().Where(&iodsp.UserRole{UserId: userId, RoleId: roleId}).Find(&userRole)
	return userRole
}

//select a resourceRole by resourceId roleId include those having been softly deleted
func (r *Grant) ShowByResourceIdRoleIdIncDel(resourceId int, roleId int) (roleResource iodsp.RoleResource) {
	r.db.Unscoped().Where(&iodsp.RoleResource{ResourceId: resourceId, RoleId: roleId}).Find(&roleResource)
	return roleResource
}

//update roleResource
func (r *Grant) UpdateResourceRole(updateResourceRole iodsp.RoleResource) error {
	updateErr := r.db.Unscoped().Save(&updateResourceRole).Error
	return updateErr
}
