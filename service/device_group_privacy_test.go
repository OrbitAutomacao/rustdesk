package service

import (
	"testing"

	"github.com/lejianwen/rustdesk-api/v2/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestPrivateDeviceGroupVisibility(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:private_device_groups?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.DeviceGroup{}, &model.Peer{}); err != nil {
		t.Fatal(err)
	}
	previous := DB
	DB = db
	t.Cleanup(func() { DB = previous })

	owner, other := uint(7), uint(8)
	private := &model.DeviceGroup{Name: "Programação", Private: true, OwnerUserId: owner}
	public := &model.DeviceGroup{Name: "Compartilhado", OwnerUserId: owner}
	for _, group := range []*model.DeviceGroup{private, public} {
		if err := db.Create(group).Error; err != nil {
			t.Fatal(err)
		}
	}
	for _, peer := range []*model.Peer{{Id: "private-pc", GroupId: private.Id}, {Id: "public-pc", GroupId: public.Id}} {
		if err := db.Create(peer).Error; err != nil {
			t.Fatal(err)
		}
	}
	groups := &GroupService{}
	if !groups.CanAccessDeviceGroup(private.Id, owner) || groups.CanAccessDeviceGroup(private.Id, other) {
		t.Fatal("private group must be accessible only to its owner")
	}
	for _, tc := range []struct {
		userId uint
		want   int64
	}{{owner, 2}, {other, 1}} {
		visibleGroups := groups.DeviceGroupList(1, 10, func(tx *gorm.DB) { VisibleDeviceGroups(tx, tc.userId) })
		visiblePeers := (&PeerService{}).List(1, 10, func(tx *gorm.DB) { VisiblePeers(tx, tc.userId) })
		if visibleGroups.Total != tc.want || visiblePeers.Total != tc.want {
			t.Fatalf("user %d sees %d groups and %d peers; want %d each", tc.userId, visibleGroups.Total, visiblePeers.Total, tc.want)
		}
	}
}
