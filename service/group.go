package service

import (
	"github.com/lejianwen/rustdesk-api/v2/model"
	"gorm.io/gorm"
)

type GroupService struct {
}

// InfoById 根据用户id取用户信息
func (us *GroupService) InfoById(id uint) *model.Group {
	u := &model.Group{}
	DB.Where("id = ?", id).First(u)
	return u
}

func (us *GroupService) List(page, pageSize uint, where func(tx *gorm.DB)) (res *model.GroupList) {
	res = &model.GroupList{}
	res.Page = int64(page)
	res.PageSize = int64(pageSize)
	tx := DB.Model(&model.Group{})
	if where != nil {
		where(tx)
	}
	tx.Count(&res.Total)
	tx.Scopes(Paginate(page, pageSize))
	tx.Find(&res.Groups)
	return
}

// Create 创建
func (us *GroupService) Create(u *model.Group) error {
	res := DB.Create(u).Error
	return res
}
func (us *GroupService) Delete(u *model.Group) error {
	return DB.Delete(u).Error
}

// Update 更新
func (us *GroupService) Update(u *model.Group) error {
	return DB.Model(u).Updates(map[string]interface{}{"name": u.Name, "type": u.Type, "private": u.Private, "owner_user_id": u.OwnerUserId}).Error
}

func VisibleGroups(tx *gorm.DB, userId uint) {
	tx.Where("private = ? OR owner_user_id = ?", false, userId)
}

func VisibleUsers(tx *gorm.DB, userId uint) {
	tx.Where("NOT EXISTS (SELECT 1 FROM groups g WHERE g.id = users.group_id AND g.private = ? AND g.owner_user_id <> ?)", true, userId)
}

func VisibleGroupPeers(tx *gorm.DB, userId uint) {
	tx.Where("NOT EXISTS (SELECT 1 FROM users u JOIN groups g ON g.id = u.group_id WHERE u.id = peers.user_id AND g.private = ? AND g.owner_user_id <> ?)", true, userId)
}

func VisibleDevicePeers(tx *gorm.DB, userId uint) {
	tx.Where("NOT EXISTS (SELECT 1 FROM device_groups dg WHERE dg.id = peers.group_id AND dg.private = ? AND dg.owner_user_id <> ?)", true, userId)
}

func VisibleAddressBooks(tx *gorm.DB, userId uint) {
	tx.Where("NOT EXISTS (SELECT 1 FROM peers p JOIN device_groups dg ON dg.id = p.group_id WHERE p.id = address_books.id AND dg.private = ? AND dg.owner_user_id <> ?)", true, userId)
	tx.Where("NOT EXISTS (SELECT 1 FROM peers p JOIN users u ON u.id = p.user_id JOIN groups g ON g.id = u.group_id WHERE p.id = address_books.id AND g.private = ? AND g.owner_user_id <> ?)", true, userId)
}

func VisibleDeviceGroups(tx *gorm.DB, userId uint) {
	tx.Where("private = ? OR owner_user_id = ?", false, userId)
}

func (us *GroupService) CanAccessGroup(id, userId uint) bool {
	group := us.InfoById(id)
	return group.Id != 0 && (!group.Private || group.OwnerUserId == userId)
}

func (us *GroupService) CanAccessUser(targetId, viewerId uint) bool {
	if targetId == 0 { return true }
	target := (&UserService{}).InfoById(targetId)
	return target.Id != 0 && us.CanAccessGroup(target.GroupId, viewerId)
}

func (us *GroupService) CanAccessDeviceGroup(id, userId uint) bool {
	if id == 0 { return true }
	group := us.DeviceGroupInfoById(id)
	return group.Id != 0 && (!group.Private || group.OwnerUserId == userId)
}

func (us *GroupService) CanAccessPeer(peer *model.Peer, viewerId uint) bool {
	return us.CanAccessUser(peer.UserId, viewerId) && us.CanAccessDeviceGroup(peer.GroupId, viewerId)
}

func (us *GroupService) CanAccessRemoteId(id string, viewerId uint) bool {
	var peers []*model.Peer
	DB.Where("id = ?", id).Find(&peers)
	for _, peer := range peers {
		if !us.CanAccessPeer(peer, viewerId) { return false }
	}
	return true
}

// DeviceGroupInfoById 根据用户id取用户信息
func (us *GroupService) DeviceGroupInfoById(id uint) *model.DeviceGroup {
	u := &model.DeviceGroup{}
	DB.Where("id = ?", id).First(u)
	return u
}

func (us *GroupService) DeviceGroupList(page, pageSize uint, where func(tx *gorm.DB)) (res *model.DeviceGroupList) {
	res = &model.DeviceGroupList{}
	res.Page = int64(page)
	res.PageSize = int64(pageSize)
	tx := DB.Model(&model.DeviceGroup{})
	if where != nil {
		where(tx)
	}
	tx.Count(&res.Total)
	tx.Scopes(Paginate(page, pageSize))
	tx.Find(&res.DeviceGroups)
	return
}

func (us *GroupService) DeviceGroupCreate(u *model.DeviceGroup) error {
	res := DB.Create(u).Error
	return res
}
func (us *GroupService) DeviceGroupDelete(u *model.DeviceGroup) error {
	return DB.Delete(u).Error
}

func (us *GroupService) DeviceGroupUpdate(u *model.DeviceGroup) error {
	return DB.Model(u).Updates(map[string]interface{}{"name": u.Name, "private": u.Private, "owner_user_id": u.OwnerUserId}).Error
}
