package api

import (
	"fmt"
	"github.com/XM-GO/PandaKit/biz"
	"github.com/XM-GO/PandaKit/oss"
	"github.com/XM-GO/PandaKit/restfulx"
	"github.com/XM-GO/PandaKit/utils"
	"pandax/apps/resource/entity"
	"pandax/apps/resource/services"
	"time"
)

/**
 * @Description 添加qq群467890197 交流学习
 * @Author 熊猫
 * @Date 2022/1/13 15:19
 **/

type ResOssesApi struct {
	ResOssesApp services.ResOssesModel
}

// @Summary ResOsses列表数据
// @Tags ResOsses
// @Param pageSize query int false "页条数"
// @Param pageNum query int false "页码"
// @Param status query string false "状态"
// @Param category query string false "分类"
// @Param ossCode query string false "编号"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Router /resource/oss/list [get]
// @Security
func (p *ResOssesApi) GetResOssesList(rc *restfulx.ReqCtx) {

	pageNum := restfulx.QueryInt(rc, "pageNum", 1)
	pageSize := restfulx.QueryInt(rc, "pageSize", 10)
	status := restfulx.QueryParam(rc, "status")
	category := restfulx.QueryParam(rc, "category")
	ossCode := restfulx.QueryParam(rc, "ossCode")

	data := entity.ResOss{Status: status, Category: category, OssCode: ossCode}
	list, total := p.ResOssesApp.FindListPage(pageNum, pageSize, data)
	li := *list
	for i, data := range li {
		data.AccessKey = utils.DdmKey(data.AccessKey)
		data.SecretKey = utils.DdmKey(data.SecretKey)
		li[i] = data
	}
	rc.ResData = map[string]any{
		"data":     list,
		"total":    total,
		"pageNum":  pageNum,
		"pageSize": pageSize,
	}
}

// @Summary 获取ResOsses
// @Description 获取JSON
// @Tags ResOsses
// @Param ossId path int true "ossId"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Router /resource/oss/{ossId} [get]
// @Security
func (p *ResOssesApi) GetResOsses(rc *restfulx.ReqCtx) {
	ossId := restfulx.PathParamInt(rc, "ossId")
	p.ResOssesApp.FindOne(int64(ossId))
}

// @Summary 添加ResOsses
// @Description 获取JSON
// @Tags ResOsses
// @Accept  application/json
// @Product application/json
// @Param data body entity.ResOss true "data"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": 400, "message": "添加失败"}"
// @Router /resource/oss [post]
// @Security X-TOKEN
func (p *ResOssesApi) InsertResOsses(rc *restfulx.ReqCtx) {
	var data entity.ResOss
	restfulx.BindQuery(rc, &data)

	p.ResOssesApp.Insert(data)
}

// @Summary 修改ResOsses
// @Description 获取JSON
// @Tags ResOsses
// @Accept  application/json
// @Product application/json
// @Param data body entity.ResOss true "body"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": 400, "message": "添加失败"}"
// @Router /resource/oss [put]
// @Security X-TOKEN
func (p *ResOssesApi) UpdateResOsses(rc *restfulx.ReqCtx) {
	var data entity.ResOss
	restfulx.BindQuery(rc, &data)
	if utils.IsDdmKey(data.AccessKey) {
		data.AccessKey = ""
	}
	if utils.IsDdmKey(data.SecretKey) {
		data.SecretKey = ""
	}
	p.ResOssesApp.Update(data)
}

// @Summary 删除ResOsses
// @Description 删除数据
// @Tags ResOsses
// @Param ossId path string true "ossId"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": 400, "message": "删除失败"}"
// @Router /resource/oss/{ossId} [delete]
func (p *ResOssesApi) DeleteResOsses(rc *restfulx.ReqCtx) {
	ossId := restfulx.PathParam(rc, "ossId")
	ossIds := utils.IdsStrToIdsIntGroup(ossId)
	p.ResOssesApp.Delete(ossIds)
}

// @Summary 上传文件ResOsses
// @Description 上传文件
// @Tags ResOsses
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": 400, "message": "删除失败"}"
// @Router /resource/oss/uploadFile [post]
func (p *ResOssesApi) UplaodResOsses(rc *restfulx.ReqCtx) {
	_, handler, _ := rc.Request.Request.FormFile("file")
	ossCode := restfulx.QueryParam(rc, "ossCode")
	list := p.ResOssesApp.FindList(entity.ResOss{OssCode: ossCode})
	li := *list
	yunFileTmpPath := "uploads/" + time.Now().Format("2006-01-02") + "/" + handler.Filename
	// 读取本地文件。
	f, openError := handler.Open()
	biz.ErrIsNil(openError, "function file.Open() Failed")
	biz.ErrIsNil(NewOss(li[0]).PutObj(yunFileTmpPath, f), "上传OSS失败")

	rc.ResData = fmt.Sprintf("https://%s.%s/%s", li[0].BucketName, li[0].Endpoint, yunFileTmpPath)
}

// @Summary 修改ResOsses状态
// @Description 获取JSON
// @Tags ResOsses
// @Accept  application/json
// @Product application/json
// @Param data body entity.ResOss true "body"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": 400, "message": "添加失败"}"
// @Router /resource/oss [put]
// @Security X-TOKEN
func (p *ResOssesApi) UpdateOssStatus(rc *restfulx.ReqCtx) {
	var data entity.ResOss
	restfulx.BindQuery(rc, &data)

	p.ResOssesApp.Update(entity.ResOss{OssId: data.OssId, Status: data.Status})
}

func NewOss(ens entity.ResOss) oss.Driver {
	switch ens.Category {
	case "0":
		return oss.NewAliOss(oss.AliConfig{
			AccessKey: ens.AccessKey,
			SecretKey: ens.SecretKey,
			Bucket:    ens.BucketName,
			Endpoint:  ens.Endpoint,
		})
	case "1":
		//return oss.NewQnOss()
	case "2":
		//return oss.NewTencentOss()
	}
	return nil
}
