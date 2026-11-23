package my

import (
	"github.com/gin-gonic/gin"
	"github.com/lejianwen/rustdesk-api/v2/http/request/admin"
	"github.com/lejianwen/rustdesk-api/v2/http/response"
	"github.com/lejianwen/rustdesk-api/v2/model"
	"github.com/lejianwen/rustdesk-api/v2/service"
	"gorm.io/gorm"
	"time"
)

type Peer struct {
}

// List 列表
// @Tags 我的设备
// @Summary 设备列表
// @Description 设备列表
// @Accept  json
// @Produce  json
// @Param page query int false "页码"
// @Param page_size query int false "页大小"
// @Param time_ago query int false "时间"
// @Param id query string false "ID"
// @Param hostname query string false "主机名"
// @Param uuids query string false "uuids 用逗号分隔"
// @Success 200 {object} response.Response{data=model.PeerList}
// @Failure 500 {object} response.Response
// @Router /admin/my/peer/list [get]
// @Security token
func (ct *Peer) List(c *gin.Context) {
	query := &admin.PeerQuery{}
	if err := c.ShouldBindQuery(query); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	u := service.AllService.UserService.CurUser(c)
	res := service.AllService.PeerService.List(query.Page, query.PageSize, func(tx *gorm.DB) {
		if service.AllService.UserService.IsAdmin(u) {
			tx.Where("(user_id = ? OR user_id = 0)", u.Id)
		} else {
			tx.Where("user_id = ?", u.Id)
		}
		if query.TimeAgo > 0 {
			lt := time.Now().Unix() - int64(query.TimeAgo)
			tx.Where("last_online_time < ?", lt)
		}
		if query.TimeAgo < 0 {
			lt := time.Now().Unix() + int64(query.TimeAgo)
			tx.Where("last_online_time > ?", lt)
		}
		if query.Id != "" {
			tx.Where("id like ?", "%"+query.Id+"%")
		}
		if query.Hostname != "" {
			tx.Where("hostname like ?", "%"+query.Hostname+"%")
		}
		if query.Uuids != "" {
			tx.Where("uuid in (?)", query.Uuids)
		}
	})
	response.Success(c, res)
}

// Claim associa um dispositivo ainda sem proprietário ao administrador atual.
func (ct *Peer) Claim(c *gin.Context) {
	u := service.AllService.UserService.CurUser(c)
	if !service.AllService.UserService.IsAdmin(u) {
		response.Fail(c, 403, "Apenas administradores podem atribuir dispositivos pendentes")
		return
	}
	var request struct {
		RowId uint `json:"row_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Fail(c, 101, "ID do dispositivo inválido")
		return
	}
	result := service.DB.Model(&model.Peer{}).Where("row_id = ? AND user_id = 0", request.RowId).Update("user_id", u.Id)
	if result.Error != nil || result.RowsAffected != 1 {
		response.Fail(c, 101, "Dispositivo indisponível para atribuição")
		return
	}
	response.Success(c, nil)
}
