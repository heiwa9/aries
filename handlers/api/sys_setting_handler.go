package api

import (
	"aries/config/setting"
	"aries/forms"
	"aries/models"
	"aries/utils"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-gomail/gomail"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type SysSettingHandler struct {
}

// @Summary 获取设置条目
// @Tags 系统设置
// @version 1.0
// @Accept application/json
// @Param name query string false "设置名称"
// @Success 100 object utils.Result 成功
// @Failure 103/104 object utils.Result 失败
// @Router /api/v1/sys_setting/items [get]
func (s *SysSettingHandler) GetSysSettingItem(ctx *gin.Context) {
	name := ctx.Query("name")
	result, _ := models.SysSettingItem{}.GetBySysSettingName(name)
	ctx.JSON(http.StatusOK, utils.Result{
		Code: utils.Success,
		Msg:  "查询成功",
		Data: result,
	})
}

// @Summary 保存网站配置信息
// @Tags 系统设置
// @version 1.0
// @Accept application/json
// @Param settingForm body forms.SiteSettingForm true "网站配置表单"
// @Success 100 object utils.Result 成功
// @Failure 103/104 object utils.Result 失败
// @Router /api/v1/sys_setting/site [post]
func (s *SysSettingHandler) SaveSiteSetting(ctx *gin.Context) {
	settingForm := forms.SiteSettingForm{}
	if err := ctx.ShouldBindJSON(&settingForm); err != nil {
		ctx.JSON(http.StatusOK, utils.Result{
			Code: utils.RequestError,
			Msg:  utils.GetFormError(err),
			Data: nil,
		})
		return
	}
	sysId, _ := strconv.ParseUint(settingForm.SysId, 10, 0)
	sysSetting := models.SysSetting{
		Model: gorm.Model{ID: uint(sysId)},
		Name:  settingForm.TypeName,
	}
	if sysId == 0 {
		if err := sysSetting.Create(); err != nil {
			log.Errorln("error: ", err.Error())
			ctx.JSON(http.StatusOK, utils.Result{
				Code: utils.ServerError,
				Msg:  "服务器端错误",
				Data: nil,
			})
			return
		}
	}
	t := reflect.TypeOf(settingForm)
	v := reflect.ValueOf(settingForm)
	var itemList []models.SysSettingItem
	for i := 0; i < t.NumField(); i++ {
		item := models.SysSettingItem{
			SysId: sysSetting.ID,
			Key:   t.Field(i).Tag.Get("json"),
			Val:   v.Field(i).Interface().(string),
		}
		itemList = append(itemList, item)
	}
	err := models.SysSettingItem{}.MultiCreateOrUpdate(itemList)
	if err != nil {
		log.Error("error: ", err.Error())
		ctx.JSON(http.StatusOK, utils.Result{
			Code: utils.ServerError,
			Msg:  "服务器端错误",
			Data: nil,
		})
		return
	}
	blogSetting, _ := models.SysSettingItem{}.GetBySysSettingName("网站设置")
	setting.BlogVars.InitBlogVars(blogSetting)
	ctx.JSON(http.StatusOK, utils.Result{
		Code: utils.Success,
		Msg:  "保存成功",
		Data: nil,
	})
}

// @Summary 保存 SMTP 服务配置信息
// @Tags 系统设置
// @version 1.0
// @Accept application/json
// @Param settingForm body forms.EmailSettingForm true "SMTP 配置表单"
// @Success 100 object utils.Result 成功
// @Failure 103/104 object utils.Result 失败
// @Router /api/v1/sys_setting/smtp [post]
func (s *SysSettingHandler) SaveSMTPSetting(ctx *gin.Context) {
	settingForm := forms.EmailSettingForm{}
	if err := ctx.ShouldBindJSON(&settingForm); err != nil {
		ctx.JSON(http.StatusOK, utils.Result{
			Code: utils.RequestError,
			Msg:  utils.GetFormError(err),
			Data: nil,
		})
		return
	}
	sysId, _ := strconv.ParseUint(settingForm.SysId, 10, 0)
	sysSetting := models.SysSetting{
		Model: gorm.Model{ID: uint(sysId)},
		Name:  settingForm.TypeName,
	}
	if sysId == 0 {
		if err := sysSetting.Create(); err != nil {
			log.Errorln("error: ", err.Error())
			ctx.JSON(http.StatusOK, utils.Result{
				Code: utils.ServerError,
				Msg:  "服务器端错误",
				Data: nil,
			})
			return
		}
	}
	settingForm.SysId = strconv.Itoa(int(sysSetting.ID))
	t := reflect.TypeOf(settingForm)
	v := reflect.ValueOf(settingForm)
	var itemList []models.SysSettingItem
	for i := 0; i < t.NumField(); i++ {
		item := models.SysSettingItem{
			SysId: sysSetting.ID,
			Key:   t.Field(i).Tag.Get("json"),
			Val:   v.Field(i).Interface().(string),
		}
		itemList = append(itemList, item)
	}
	err := models.SysSettingItem{}.MultiCreateOrUpdate(itemList)
	if err != nil {
		log.Error("error: ", err.Error())
		ctx.JSON(http.StatusOK, utils.Result{
			Code: utils.ServerError,
			Msg:  "服务器端错误",
			Data: nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, utils.Result{
		Code: utils.Success,
		Msg:  "保存成功",
		Data: nil,
	})
}

// @Summary 保存 sm.ms 配置信息
// @Tags 系统设置
// @version 1.0
// @Accept application/json
// @Param settingForm body forms.SmmsForm true "sm.ms 配置表单"
// @Success 100 object utils.Result 成功
// @Failure 103/104 object utils.Result 失败
// @Router /api/v1/sys_setting/pic_bed/smms [post]
func (s *SysSettingHandler) SaveSmmsSetting(ctx *gin.Context) {
	smmsForm := forms.SmmsForm{}
	if err := ctx.ShouldBindJSON(&smmsForm); err != nil {
		ctx.JSON(http.StatusOK, utils.Result{
			Code: utils.RequestError,
			Msg:  utils.GetFormError(err),
			Data: nil,
		})
		return
	}
	sysId, _ := strconv.ParseUint(smmsForm.SysId, 10, 0)
	smmsSetting := models.SysSetting{
		Model: gorm.Model{ID: uint(sysId)},
		Name:  smmsForm.StorageType,
	}
	picBedSetting := models.SysSetting{
		Name: "图床设置",
	}
	picBedSettingItems, _ := models.SysSettingItem{}.GetBySysSettingName("图床设置")
	if len(picBedSettingItems) == 0 {
		if err := picBedSetting.Create(); err != nil {
			log.Errorln("error: ", err.Error())
			ctx.JSON(http.StatusOK, utils.Result{
				Code: utils.ServerError,
				Msg:  "服务器端错误",
				Data: nil,
			})
			return
		}
	} else {
		sysId, _ := strconv.Atoi(picBedSettingItems["sys_id"])
		picBedSetting.ID = uint(sysId)
	}
	if sysId == 0 {
		if err := smmsSetting.Create(); err != nil {
			log.Errorln("error: ", err.Error())
			ctx.JSON(http.StatusOK, utils.Result{
				Code: utils.ServerError,
				Msg:  "服务器端错误",
				Data: nil,
			})
			return
		}
	}
	picBedForm := forms.PicBedSettingForm{
		SysId:       strconv.Itoa(int(picBedSetting.ID)),
		StorageType: "sm.ms",
	}
	t := reflect.TypeOf(picBedForm)
	v := reflect.ValueOf(picBedForm)
	var itemList []models.SysSettingItem
	for i := 0; i < t.NumField(); i++ {
		item := models.SysSettingItem{
			SysId: picBedSetting.ID,
			Key:   t.Field(i).Tag.Get("json"),
			Val:   v.Field(i).Interface().(string),
		}
		itemList = append(itemList, item)
	}
	err := models.SysSettingItem{}.MultiCreateOrUpdate(itemList)
	if err != nil {
		log.Error("error: ", err.Error())
		ctx.JSON(http.StatusOK, utils.Result{
			Code: utils.ServerError,
			Msg:  "服务器端错误",
			Data: nil,
		})
		return
	}
	smmsForm.SysId = strconv.Itoa(int(smmsSetting.ID))
	t = reflect.TypeOf(smmsForm)
	v = reflect.ValueOf(smmsForm)
	itemList = []models.SysSettingItem{}
	for i := 0; i < t.NumField(); i++ {
		item := models.SysSettingItem{
			SysId: smmsSetting.ID,
			Key:   t.Field(i).Tag.Get("json"),
			Val:   v.Field(i).Interface().(string),
		}
		itemList = append(itemList, item)
	}
	err = models.SysSettingItem{}.MultiCreateOrUpdate(itemList)
	if err != nil {
		log.Error("error: ", err.Error())
		ctx.JSON(http.StatusOK, utils.Result{
			Code: utils.ServerError,
			Msg:  "服务器端错误",
			Data: nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, utils.Result{
		Code: utils.Success,
		Msg:  "保存成功",
		Data: nil,
	})
}

// @Summary 保存腾讯云 COS 配置信息
// @Tags 系统设置
// @version 1.0
// @Accept application/json
// @Param cosForm body forms.TencentCosForm true "腾讯云 COS 配置表单"
// @Success 100 object utils.Result 成功
// @Failure 103/104 object utils.Result 失败
// @Router /api/v1/sys_setting/pic_bed/tencent_cos [post]
func (s *SysSettingHandler) SaveTencentCosSetting(ctx *gin.Context) {
	cosForm := forms.TencentCosForm{}
	if err := ctx.ShouldBindJSON(&cosForm); err != nil {
		ctx.JSON(http.StatusOK, utils.Result{
			Code: utils.RequestError,
			Msg:  utils.GetFormError(err),
			Data: nil,
		})
		return
	}
	sysId, _ := strconv.ParseUint(cosForm.SysId, 10, 0)
	cosSetting := models.SysSetting{
		Model: gorm.Model{ID: uint(sysId)},
		Name:  cosForm.StorageType,
	}
	picBedSetting := models.SysSetting{
		Name: "图床设置",
	}
	picBedSettingItems, _ := models.SysSettingItem{}.GetBySysSettingName("图床设置")
	if len(picBedSettingItems) == 0 {
		if err := picBedSetting.Create(); err != nil {
			log.Errorln("error: ", err.Error())
			ctx.JSON(http.StatusOK, utils.Result{
				Code: utils.ServerError,
				Msg:  "服务器端错误",
				Data: nil,
			})
			return
		}
	} else {
		sysId, _ := strconv.Atoi(picBedSettingItems["sys_id"])
		picBedSetting.ID = uint(sysId)
	}
	if sysId == 0 {
		if err := cosSetting.Create(); err != nil {
			log.Errorln("error: ", err.Error())
			ctx.JSON(http.StatusOK, utils.Result{
				Code: utils.ServerError,
				Msg:  "服务器端错误",
				Data: nil,
			})
			return
		}
	}
	picBedForm := forms.PicBedSettingForm{
		SysId:       strconv.Itoa(int(picBedSetting.ID)),
		StorageType: "cos",
	}
	t := reflect.TypeOf(picBedForm)
	v := reflect.ValueOf(picBedForm)
	var itemList []models.SysSettingItem
	for i := 0; i < t.NumField(); i++ {
		item := models.SysSettingItem{
			SysId: picBedSetting.ID,
			Key:   t.Field(i).Tag.Get("json"),
			Val:   v.Field(i).Interface().(string),
		}
		itemList = append(itemList, item)
	}
	err := models.SysSettingItem{}.MultiCreateOrUpdate(itemList)
	if err != nil {
		log.Error("error: ", err.Error())
		ctx.JSON(http.StatusOK, utils.Result{
			Code: utils.ServerError,
			Msg:  "服务器端错误",
			Data: nil,
		})
		return
	}
	cosForm.SysId = strconv.Itoa(int(cosSetting.ID))
	t = reflect.TypeOf(cosForm)
	v = reflect.ValueOf(cosForm)
	itemList = []models.SysSettingItem{}
	for i := 0; i < t.NumField(); i++ {
		item := models.SysSettingItem{
			SysId: cosSetting.ID,
			Key:   t.Field(i).Tag.Get("json"),
			Val:   v.Field(i).Interface().(string),
		}
		itemList = append(itemList, item)
	}
	err = models.SysSettingItem{}.MultiCreateOrUpdate(itemList)
	if err != nil {
		log.Error("error: ", err.Error())
		ctx.JSON(http.StatusOK, utils.Result{
			Code: utils.ServerError,
			Msg:  "服务器端错误",
			Data: nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, utils.Result{
		Code: utils.Success,
		Msg:  "保存成功",
		Data: nil,
	})
}

// @Summary 发送测试邮件
// @Tags 系统设置
// @version 1.0
// @Accept application/json
// Param sendForm body orm.EmailSendForm true "发送邮件表单"
// @Success 100 object utils.Result 成功
// @Failure 103/104 object utils.Result 失败
// @Router /api/v1/sys_setting/email/test [post]
func (s *SysSettingHandler) SendTestEmail(ctx *gin.Context) {
	sendForm := forms.EmailSendForm{}
	if err := ctx.ShouldBindJSON(&sendForm); err != nil {
		ctx.JSON(http.StatusOK, utils.Result{
			Code: utils.RequestError,
			Msg:  utils.GetFormError(err),
			Data: nil,
		})
		return
	}
	emailSetting, err := models.SysSettingItem{}.GetBySysSettingName("邮件设置")
	if err != nil {
		log.Error("error: ", err.Error())
		ctx.JSON(http.StatusOK, utils.Result{
			Code: utils.ServerError,
			Msg:  "服务器端错误",
			Data: nil,
		})
		return
	}
	if len(emailSetting) == 0 {
		ctx.JSON(http.StatusOK, utils.Result{
			Code: utils.RequestError,
			Msg:  "请先配置 SMTP，再进行邮件测试",
			Data: nil,
		})
		return
	}
	msg := gomail.NewMessage()
	// 设置收件人
	msg.SetHeader("To", sendForm.ReceiveEmail)
	// 设置发件人
	msg.SetAddressHeader("From", emailSetting["account"], sendForm.Sender)
	// 主题
	msg.SetHeader("Subject", sendForm.Title)
	// 正文
	msg.SetBody("text/html", utils.GetEmailHTML(sendForm.Title, sendForm.ReceiveEmail,
		sendForm.Content))
	port, _ := strconv.Atoi(emailSetting["port"])
	// 设置 SMTP 参数
	d := gomail.NewDialer(emailSetting["address"], port, emailSetting["account"], emailSetting["pwd"])
	// 发送
	err = d.DialAndSend(msg)
	if err != nil {
		log.Error("邮件发送失败：", err.Error())
		ctx.JSON(http.StatusOK, utils.Result{
			Code: utils.ServerError,
			Msg:  "邮件发送失败",
			Data: nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, utils.Result{
		Code: utils.Success,
		Msg:  "发送成功",
		Data: nil,
	})
}

// @Summary 获取后台首页数据
// @Tags 系统设置
// @version 1.0
// @Accept application/json
// @Success 100 object utils.Result 成功
// @Failure 103/104 object utils.Result 失败
// @Router /api/v1/sys_setting/index_info [get]
func (s *SysSettingHandler) GetAdminIndexData(ctx *gin.Context) {
	articleCount, err := models.Article{}.GetCount()
	if err != nil {
		log.Error("error: ", err.Error())
		ctx.JSON(http.StatusOK, utils.Result{
			Code: utils.ServerError,
			Msg:  "服务器端错误",
			Data: nil,
		})
		return
	}
	commentCount, err := models.Comment{}.GetCount()
	if err != nil {
		log.Error("error: ", err.Error())
		ctx.JSON(http.StatusOK, utils.Result{
			Code: utils.ServerError,
			Msg:  "服务器端错误",
			Data: nil,
		})
		return
	}
	latestArticles, err := models.Article{}.GetLatest(6)
	if err != nil {
		log.Error("error: ", err.Error())
		ctx.JSON(http.StatusOK, utils.Result{
			Code: utils.ServerError,
			Msg:  "服务器端错误",
			Data: nil,
		})
		return
	}
	latestComments, err := models.Comment{}.GetLatest(6)
	if err != nil {
		log.Error("error: ", err.Error())
		ctx.JSON(http.StatusOK, utils.Result{
			Code: utils.ServerError,
			Msg:  "服务器端错误",
			Data: nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, utils.Result{
		Code: utils.Success,
		Msg:  "查询成功",
		Data: gin.H{
			"article_count":   articleCount,
			"comment_count":   commentCount,
			"latest_articles": latestArticles,
			"latest_comments": latestComments,
		},
	})
}
