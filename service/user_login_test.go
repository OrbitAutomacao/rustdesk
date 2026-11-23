package service

import (
	"testing"

	"github.com/lejianwen/rustdesk-api/v2/config"
	"github.com/lejianwen/rustdesk-api/v2/model"
	"github.com/lejianwen/rustdesk-api/v2/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestInfoByUsernamePasswordIgnoresUsernameCase(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:login_case_test?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.User{}); err != nil {
		t.Fatal(err)
	}
	oldDB, oldConfig := DB, Config
	DB, Config = db, &config.Config{}
	t.Cleanup(func() { DB, Config = oldDB, oldConfig })

	hash, err := utils.EncryptPassword("valid-test-password")
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&model.User{Username: "edeymerson", Password: hash}).Error; err != nil {
		t.Fatal(err)
	}

	users := &UserService{}
	if got := users.InfoByUsernamePassword("EDEYMERSON", "valid-test-password"); got.Id == 0 {
		t.Fatal("uppercase username should match the stored lowercase account")
	}
	if got := users.InfoByUsernamePassword("EDEYMERSON", "wrong-password"); got.Id != 0 {
		t.Fatal("wrong password must still be rejected")
	}
}
