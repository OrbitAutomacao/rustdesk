package model

const (
	GroupTypeDefault = 1 // 默认
	GroupTypeShare   = 2 // 共享
)

type Group struct {
	IdModel
	Name string `json:"name" gorm:"default:'';not null;"`
	Type int    `json:"type" gorm:"default:1;not null;"`
	Private bool `json:"private" gorm:"default:false;not null;index"`
	OwnerUserId uint `json:"owner_user_id" gorm:"default:0;not null;index"`
	TimeModel
}

type GroupList struct {
	Groups []*Group `json:"list"`
	Pagination
}

type DeviceGroup struct {
	IdModel
	Name string `json:"name" gorm:"default:'';not null;"`
	Private bool `json:"private" gorm:"default:false;not null;index"`
	OwnerUserId uint `json:"owner_user_id" gorm:"default:0;not null;index"`
	TimeModel
}

type DeviceGroupList struct {
	DeviceGroups []*DeviceGroup `json:"list"`
	Pagination
}
