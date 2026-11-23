package service

import (
	"testing"

	"github.com/lejianwen/rustdesk-api/v2/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestPrivateUserGroupVisibility(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:private_user_groups?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil { t.Fatal(err) }
	if err := db.AutoMigrate(&model.Group{}, &model.User{}, &model.Peer{}); err != nil { t.Fatal(err) }
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
	groups := &GroupService{}
	if !groups.CanAccessGroup(programming.Id, owner) || groups.CanAccessGroup(programming.Id, technician) { t.Fatal("private group ownership check failed") }
	for _, tc := range []struct { userId uint; want int64 }{{owner, 2}, {technician, 1}} {
		groupList := groups.List(1, 10, func(tx *gorm.DB) { VisibleGroups(tx, tc.userId) })
		userList := (&UserService{}).List(1, 10, func(tx *gorm.DB) { VisibleUsers(tx, tc.userId) })
		peerList := (&PeerService{}).List(1, 10, func(tx *gorm.DB) { VisibleGroupPeers(tx, tc.userId) })
		if groupList.Total != tc.want || userList.Total != tc.want || peerList.Total != tc.want {
			t.Fatalf("viewer %d saw groups=%d users=%d peers=%d; want %d", tc.userId, groupList.Total, userList.Total, peerList.Total, tc.want)
		}
	}
}
