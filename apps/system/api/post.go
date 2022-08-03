package api

import (
	"errors"
	"fmt"
	"github.com/XM-GO/PandaKit/biz"
	"github.com/XM-GO/PandaKit/restfulx"
	"github.com/XM-GO/PandaKit/utils"
	"pandax/apps/system/entity"
	"pandax/apps/system/services"
	"pandax/pkg/global"
)

type PostApi struct {
	PostApp services.SysPostModel
	UserApp services.SysUserModel
	RoleApp services.SysRoleModel
}

// @Summary 职位列表数据
// @Description 获取JSON
// @Tags 职位
// @Param postName query string false "postName"
// @Param postCode query string false "postCode"
// @Param status query string false "status"
// @Param pageSize query int false "页条数"
// @Param pageNum query int false "页码"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Router /system/post [get]
// @Security
func (p *PostApi) GetPostList(rc *restfulx.ReqCtx) {

	pageNum := restfulx.QueryInt(rc, "pageNum", 1)
	pageSize := restfulx.QueryInt(rc, "pageSize", 10)
	status := restfulx.QueryParam(rc, "status")
	postName := restfulx.QueryParam(rc, "postName")
	postCode := restfulx.QueryParam(rc, "postCode")
	post := entity.SysPost{Status: status, PostName: postName, PostCode: postCode}

	if !IsTenantAdmin(rc.LoginAccount.TenantId) {
		post.TenantId = rc.LoginAccount.TenantId
	}

	list, total := p.PostApp.FindListPage(pageNum, pageSize, post)

	rc.ResData = map[string]any{
		"data":     list,
		"total":    total,
		"pageNum":  pageNum,
		"pageSize": pageSize,
	}
}

// @Summary 获取职位
// @Description 获取JSON
// @Tags 职位
// @Param postId path int true "postId"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Router /system/post/{postId} [get]
// @Security
func (p *PostApi) GetPost(rc *restfulx.ReqCtx) {
	postId := restfulx.PathParamInt(rc, "postId")
	p.PostApp.FindOne(int64(postId))
}

// @Summary 添加职位
// @Description 获取JSON
// @Tags 职位
// @Accept  application/json
// @Product application/json
// @Param data body entity.SysPost true "data"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": 400, "message": "添加失败"}"
// @Router /system/post [post]
// @Security X-TOKEN
func (p *PostApi) InsertPost(rc *restfulx.ReqCtx) {
	var post entity.SysPost
	restfulx.BindQuery(rc, &post)
	post.TenantId = rc.LoginAccount.TenantId
	post.CreateBy = rc.LoginAccount.UserName
	p.PostApp.Insert(post)
}

// @Summary 修改职位
// @Description 获取JSON
// @Tags 职位
// @Accept  application/json
// @Product application/json
// @Param data body entity.SysPost true "body"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": 400, "message": "添加失败"}"
// @Router /system/post [put]
// @Security X-TOKEN
func (p *PostApi) UpdatePost(rc *restfulx.ReqCtx) {
	var post entity.SysPost
	restfulx.BindQuery(rc, &post)

	post.CreateBy = rc.LoginAccount.UserName
	p.PostApp.Update(post)
}

// @Summary 删除职位
// @Description 删除数据
// @Tags 职位
// @Param postId path string true "postId "
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": 400, "message": "删除失败"}"
// @Router /system/post/{postId} [delete]
func (p *PostApi) DeletePost(rc *restfulx.ReqCtx) {
	postId := restfulx.PathParam(rc, "postId")
	postIds := utils.IdsStrToIdsIntGroup(postId)

	deList := make([]int64, 0)
	for _, id := range postIds {
		user := entity.SysUser{}
		user.PostId = id
		list := p.UserApp.FindList(user)
		if len(*list) == 0 {
			deList = append(deList, id)
		} else {
			global.Log.Info(fmt.Sprintf("dictId: %d 存在岗位绑定用户无法删除", id))
		}
	}
	if len(deList) == 0 {
		biz.ErrIsNil(errors.New("所有岗位都已绑定用户，无法删除"), "所有岗位都已绑定用户，无法删除")
	}
	p.PostApp.Delete(deList)
}
