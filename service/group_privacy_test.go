package service

import (
	"testing"

	"github.com/lejianwen/rustdesk-api/v2/model"
	"github.com/lejianwen/rustdesk-api/v2/model/custom_types"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestPrivateUserGroupVisibility(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:private_user_groups?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil { t.Fatal(err) }
	if err := db.AutoMigrate(&model.Group{}, &model.User{}, &model.DeviceGroup{}, &model.Peer{}, &model.AddressBook{}); err != nil { t.Fatal(err) }
	previous := DB
	DB = db
	t.Cleanup(func() { DB = previous })

	owner, technician := uint(1), uint(9)
	programming := &model.Group{Name: "PROGRAMAÇÃO", Type: model.GroupTypeDefault, Private: true, OwnerUserId: owner}
	shared := &model.Group{Name: "ORBIT", Type: model.GroupTypeShare}
	for _, group := range []*model.Group{programming, shared} {
		if err := db.Create(group).Error; err != nil { t.Fatal(err) }
	}
	for _, user := range []*model.User{{Username: "dev", GroupId: programming.Id}, {Username: "support", GroupId: shared.Id}} {
		if err := db.Create(user).Error; err != nil { t.Fatal(err) }
		if err := db.Create(&model.Peer{Id: user.Username + "-pc", UserId: user.Id}).Error; err != nil { t.Fatal(err) }
	}
	deviceGroup := &model.DeviceGroup{Name: "DEV", Private: true, OwnerUserId: owner}
	if err := db.Create(deviceGroup).Error; err != nil { t.Fatal(err) }
	privatePeer := &model.Peer{Id: "private-device", UserId: 2, GroupId: deviceGroup.Id}
	if err := db.Create(privatePeer).Error; err != nil { t.Fatal(err) }
	if err := db.Create(&model.AddressBook{Id: privatePeer.Id, UserId: technician, Tags: custom_types.AutoJson("[]")}).Error; err != nil { t.Fatal(err) }
	groups := &GroupService{}
	if !groups.CanAccessGroup(programming.Id, owner) || groups.CanAccessGroup(programming.Id, technician) { t.Fatal("private group ownership check failed") }
	if !groups.CanAccessPeer(privatePeer, owner) || groups.CanAccessPeer(privatePeer, technician) { t.Fatal("private device group ownership check failed") }
	if !groups.CanAccessRemoteId(privatePeer.Id, owner) || groups.CanAccessRemoteId(privatePeer.Id, technician) { t.Fatal("private device ID must be hidden from other administrators") }
	for _, tc := range []struct { userId uint; want int64 }{{owner, 2}, {technician, 1}} {
		groupList := groups.List(1, 10, func(tx *gorm.DB) { VisibleGroups(tx, tc.userId) })
		userList := (&UserService{}).List(1, 10, func(tx *gorm.DB) { VisibleUsers(tx, tc.userId) })
		peerList := (&PeerService{}).List(1, 10, func(tx *gorm.DB) { VisibleGroupPeers(tx, tc.userId); VisibleDevicePeers(tx, tc.userId) })
		deviceGroups := groups.DeviceGroupList(1, 10, func(tx *gorm.DB) { VisibleDeviceGroups(tx, tc.userId) })
		wantPeers := tc.want
		wantDeviceGroups := int64(0)
		if tc.userId == owner { wantPeers = 3; wantDeviceGroups = 1 }
		if groupList.Total != tc.want || userList.Total != tc.want || peerList.Total != wantPeers || deviceGroups.Total != wantDeviceGroups {
			t.Fatalf("viewer %d saw groups=%d users=%d peers=%d device_groups=%d", tc.userId, groupList.Total, userList.Total, peerList.Total, deviceGroups.Total)
		}
	}
	if list := (&AddressBookService{}).ListByUserIdAndCollectionId(technician, 0, 1, 10); list.Total != 0 {
		t.Fatal("private device must not appear in another user's catalog")
	}
}
