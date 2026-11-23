package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/lejianwen/rustdesk-api/v2/global"
	"github.com/lejianwen/rustdesk-api/v2/http/request/admin"
	"github.com/lejianwen/rustdesk-api/v2/http/response"
	"github.com/lejianwen/rustdesk-api/v2/model"
	"github.com/lejianwen/rustdesk-api/v2/service"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type Peer struct {
}

// Detail 设备
// @Tags 设备
// @Summary 设备详情
// @Description 设备详情
// @Accept  json
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} response.Response{data=model.Peer}
// @Failure 500 {object} response.Response
// @Router /admin/peer/detail/{id} [get]
// @Security token
func (ct *Peer) Detail(c *gin.Context) {
	id := c.Param("id")
	iid, _ := strconv.Atoi(id)
	u := service.AllService.PeerService.InfoByRowId(uint(iid))
	current := service.AllService.UserService.CurUser(c)
	if u.RowId > 0 && service.AllService.GroupService.CanAccessPeer(u, current.Id) {
		response.Success(c, u)
		return
	}
	response.Fail(c, 101, response.TranslateMsg(c, "ItemNotFound"))
	return
}

// Create 创建设备
// @Tags 设备
// @Summary 创建设备
// @Description 创建设备
// @Accept  json
// @Produce  json
// @Param body body admin.PeerForm true "设备信息"
// @Success 200 {object} response.Response{data=model.Peer}
// @Failure 500 {object} response.Response
// @Router /admin/peer/create [post]
// @Security token
func (ct *Peer) Create(c *gin.Context) {
	f := &admin.PeerForm{}
	if err := c.ShouldBindJSON(f); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	errList := global.Validator.ValidStruct(c, f)
	if len(errList) > 0 {
		response.Fail(c, 101, errList[0])
		return
	}
	p := f.ToPeer()
	if !service.AllService.GroupService.CanAccessPeer(p, service.AllService.UserService.CurUser(c).Id) {
		response.Fail(c, 403, response.TranslateMsg(c, "NoAccess"))
		return
	}
	err := service.AllService.PeerService.Create(p)
	if err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "OperationFailed")+err.Error())
		return
	}
	response.Success(c, nil)
}

// List 列表
// @Tags 设备
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
// @Router /admin/peer/list [get]
// @Security token
func (ct *Peer) List(c *gin.Context) {
	query := &admin.PeerQuery{}
	if err := c.ShouldBindQuery(query); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	res := service.AllService.PeerService.List(query.Page, query.PageSize, func(tx *gorm.DB) {
		service.VisibleGroupPeers(tx, service.AllService.UserService.CurUser(c).Id)
		service.VisibleDevicePeers(tx, service.AllService.UserService.CurUser(c).Id)
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
		if query.Username != "" {
			tx.Where("username like ?", "%"+query.Username+"%")
		}
		if query.Ip != "" {
			tx.Where("last_online_ip like ?", "%"+query.Ip+"%")
		}
		if query.Alias != "" {
			tx.Where("alias like ?", "%"+query.Alias+"%")
		}
	})
	response.Success(c, res)
}

// Update 编辑
// @Tags 设备
// @Summary 设备编辑
// @Description 设备编辑
// @Accept  json
// @Produce  json
// @Param body body admin.PeerForm true "设备信息"
// @Success 200 {object} response.Response{data=model.Peer}
// @Failure 500 {object} response.Response
// @Router /admin/peer/update [post]
// @Security token
func (ct *Peer) Update(c *gin.Context) {
	f := &admin.PeerForm{}
	if err := c.ShouldBindJSON(f); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	if f.RowId == 0 {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError"))
		return
	}
	errList := global.Validator.ValidStruct(c, f)
	if len(errList) > 0 {
		response.Fail(c, 101, errList[0])
		return
	}
	u := f.ToPeer()
	current := service.AllService.UserService.CurUser(c)
	existing := service.AllService.PeerService.InfoByRowId(f.RowId)
	if existing.RowId == 0 || !service.AllService.GroupService.CanAccessPeer(existing, current.Id) || !service.AllService.GroupService.CanAccessPeer(u, current.Id) {
		response.Fail(c, 403, response.TranslateMsg(c, "NoAccess"))
		return
	}
	err := service.AllService.PeerService.Update(u)
	if err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "OperationFailed")+err.Error())
		return
	}
	response.Success(c, nil)
}

// Delete 删除
// @Tags 设备
// @Summary 设备删除
// @Description 设备删除
// @Accept  json
// @Produce  json
// @Param body body admin.PeerForm true "设备信息"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/peer/delete [post]
// @Security token
func (ct *Peer) Delete(c *gin.Context) {
	f := &admin.PeerForm{}
	if err := c.ShouldBindJSON(f); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	id := f.RowId
	errList := global.Validator.ValidVar(c, id, "required,gt=0")
	if len(errList) > 0 {
		response.Fail(c, 101, errList[0])
		return
	}
	u := service.AllService.PeerService.InfoByRowId(f.RowId)
	current := service.AllService.UserService.CurUser(c)
	if u.RowId > 0 && service.AllService.GroupService.CanAccessPeer(u, current.Id) {
		err := service.AllService.PeerService.Delete(u)
		if err == nil {
			response.Success(c, nil)
			return
		}
		response.Fail(c, 101, response.TranslateMsg(c, "OperationFailed")+err.Error())
		return
	}
	response.Fail(c, 101, response.TranslateMsg(c, "ItemNotFound"))
}

// BatchDelete 批量删除
// @Tags 设备
// @Summary 批量设备删除
// @Description 批量设备删除
// @Accept  json
// @Produce  json
// @Param body body admin.PeerBatchDeleteForm true "设备id"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/peer/batchDelete [post]
// @Security token
func (ct *Peer) BatchDelete(c *gin.Context) {
	f := &admin.PeerBatchDeleteForm{}
	if err := c.ShouldBindJSON(f); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	if len(f.RowIds) == 0 {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError"))
		return
	}
	var visible int64
	viewer := service.AllService.UserService.CurUser(c)
	service.DB.Model(&model.Peer{}).Where("row_id in ?", f.RowIds).Scopes(func(tx *gorm.DB) *gorm.DB { service.VisibleGroupPeers(tx, viewer.Id); service.VisibleDevicePeers(tx, viewer.Id); return tx }).Count(&visible)
	if visible != int64(len(f.RowIds)) {
		response.Fail(c, 403, response.TranslateMsg(c, "NoAccess"))
		return
	}
	err := service.AllService.PeerService.BatchDelete(f.RowIds)
	if err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "OperationFailed")+err.Error())
		return
	}
	response.Success(c, nil)
}

func (ct *Peer) SimpleData(c *gin.Context) {
	f := &admin.SimpleDataQuery{}
	if err := c.ShouldBindJSON(f); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	if len(f.Ids) == 0 {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError"))
		return
	}
	res := service.AllService.PeerService.List(1, 99999, func(tx *gorm.DB) {
		service.VisibleGroupPeers(tx, service.AllService.UserService.CurUser(c).Id)
		service.VisibleDevicePeers(tx, service.AllService.UserService.CurUser(c).Id)
		//可以公开的情报
		tx.Select("id,version")
		tx.Where("id in (?)", f.Ids)
	})
	response.Success(c, res)
}
