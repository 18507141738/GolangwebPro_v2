package controllers

import (
	"Artifice_V2.0.0/models"
	"Artifice_V2.0.0/util"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"log"
	"strconv"
	"strings"
	"time"
)

//展示首页
type PlatHomePage struct {
	beego.Controller
}

func (this *PlatHomePage) Get() {
	this.Data["SystemTitle"] = SelSystemTitleUtil()
	this.Data["bmk"] = "homemenu"
	this.Data["Time"] = GetTimeSE()
	admin_id := this.Ctx.GetCookie("user_id")
	o := O
	var user models.User
	err := o.QueryTable(new(models.User)).Filter("admin_id", admin_id).RelatedSel().One(&user)
	if err != nil {
		log.Println(err)
		this.TplName = "login.html"
	} else {
		this.Data["admin_name"] = user.Name
		this.Data["admin_user"] = user.Acount
		this.Data["organize_name"] = user.Organize.Name
		this.Data["jurisdiction"] = user.Jurisdiction
		this.TplName = "platform/index.html"
	}

}

//展示功能页面
type PlatFuncPage struct {
	beego.Controller
}

func (this *PlatFuncPage) Get() {
	this.Data["SystemTitle"] = SelSystemTitleUtil()
	t := this.GetString("type")
	this.Data["bmk"] = titleTypeByParamType(t)
	this.Data["Time"] = GetTimeSE()
	admin_id := this.Ctx.GetCookie("user_id")
	o := O
	var user models.User
	err := o.QueryTable(new(models.User)).Filter("admin_id", admin_id).RelatedSel().One(&user)
	if err != nil {
		log.Println(err)
		this.TplName = "platform/index.html"
	} else {
		this.Data["admin_name"] = user.Name
		this.Data["admin_user"] = user.Acount
		this.Data["organize_name"] = user.Organize.Name
		this.Data["jurisdiction"] = user.Jurisdiction
		this.TplName = "platform/funcPage.html"
	}

}

//展示功能页面
type PlatFuncPageIE struct {
	beego.Controller
}

func (this *PlatFuncPageIE) Get() {
	this.Data["SystemTitle"] = SelSystemTitleUtil()
	t := this.GetString("type")
	this.Data["bmk"] = titleTypeByParamType(t)
	this.Data["Time"] = GetTimeSE()
	admin_id := this.Ctx.GetCookie("user_id")
	o := O
	var user models.User
	err := o.QueryTable(new(models.User)).Filter("admin_id", admin_id).RelatedSel().One(&user)
	if err != nil {
		log.Println(err)
		this.TplName = "platform/index.html"
	} else {
		this.Data["admin_name"] = user.Name
		this.Data["admin_user"] = user.Acount
		this.Data["organize_name"] = user.Organize.Name
		this.Data["jurisdiction"] = user.Jurisdiction
		this.TplName = "platform/funcPageIE.html"
	}

}

type PlatHomeController struct {
	beego.Controller
}

func (this *PlatHomeController) Post() {
	userid := this.Ctx.GetCookie("user_id")
	var startTime = time.Now().AddDate(0, 0, -30).Format("2006-01-02") + " 00:00:00"
	var endTime = time.Now().Format("2006-01-02") + " 23:59:59"

	o := O
	var user models.User

	//获取用户账户信息
	err := o.QueryTable(new(models.User)).Filter("admin_id", userid).RelatedSel().One(&user)
	if err != nil {
		log.Println(err)
		this.Data["json"] = map[string]interface{}{"Code": 0, "Reason": "查询用户信息异常"}
		this.ServeJSON()
		return
	}

	if user.Organize.ID == "" {
		this.Data["json"] = map[string]interface{}{"Code": 0, "Reason": "用户组织架构异常"}
		this.ServeJSON()
		return
	}

	if user.Jurisdiction == "" {
		this.Data["json"] = map[string]interface{}{"Code": 0, "Reason": "用户账号无权限"}
		this.ServeJSON()
		return
	}

	jurs := strings.Split(user.Jurisdiction, ",")

	if len(jurs) == 0 {
		this.Data["json"] = map[string]interface{}{"Code": 0, "Reason": "用户账号无权限"}
		this.ServeJSON()
		return
	}

	orgIds, err := SelOrgIDBelongUserOrg(user.Organize)
	if err != nil {
		LogsError("查询用户下属组织异常:", err)

	}
	var ids = []string{user.Organize.ID}
	for _, v := range orgIds {
		ids = append(ids, v["id"].(string))
	}

	countInfo := make(orm.Params)
	jurc := [...]string{"cloths", "smoke", "fire", "boundary", "queue_count", "sleep_count", "leakage"}
	var z_num = []string{}
	var z_key = []string{}
	var c_num = []string{}
	for _, v := range jurs {
		if v == "7" || v == "8" || v == "9" {
			continue
		}
		n, err := strconv.Atoi(v)
		if err != nil {
			log.Println(err)
			this.Data["json"] = map[string]interface{}{"Code": 0, "Reason": "用户权限异常"}
			this.ServeJSON()
			return
		}
		//查询用户下的主机特定区域的告警统计
		cnt, err := UserAlarmClassifyCount(ids, jurc[n])
		if err != nil {
			this.Data["json"] = map[string]interface{}{"Code": 0, "Reason": "分类统计异常"}
			this.ServeJSON()
			return
		}
		log.Println(jurc[n], ":", cnt)

		switch v {
		case "0":
			countInfo["zh"] = cnt
			z_num = append(z_num, "IFNULL(a.num_zh,0) as num_zh")
			c_num = append(c_num, "COUNT(alarm_type = 'cloths' OR NULL) AS num_zh")
			z_key = append(z_key, "num_zh")
		case "1":
			countInfo["yw"] = cnt
			z_num = append(z_num, "IFNULL(a.num_yw,0) as num_yw")
			c_num = append(c_num, "COUNT(alarm_type = 'smoke' OR NULL) AS num_yw")
			z_key = append(z_key, "num_yw")
		case "2":
			countInfo["hm"] = cnt
			z_num = append(z_num, "IFNULL(a.num_hm,0) as num_hm")
			c_num = append(c_num, "COUNT(alarm_type = 'fire' OR NULL) AS num_hm")
			z_key = append(z_key, "num_hm")
		case "3":
			countInfo["qy"] = cnt
			z_num = append(z_num, "IFNULL(a.num_qy,0) as num_qy")
			c_num = append(c_num, "COUNT(alarm_type = 'boundary' OR NULL) AS num_qy")
			z_key = append(z_key, "num_qy")
		case "4":
			countInfo["lg"] = cnt
			z_num = append(z_num, "IFNULL(a.num_lg,0) as num_lg")
			c_num = append(c_num, "COUNT(alarm_type = 'queue_count' OR NULL) AS num_lg")
			z_key = append(z_key, "num_lg")
		case "5":
			countInfo["sg"] = cnt
			z_num = append(z_num, "IFNULL(a.num_sg,0) as num_sg")
			c_num = append(c_num, "COUNT(alarm_type = 'sleep_count' OR NULL) AS num_sg")
			z_key = append(z_key, "num_sg")
			//case "6":
			//	countInfo["xl"] = cnt
			//	z_num = append(z_num, "IFNULL(a.num_xl,0) as num_xl")
			//	c_num = append(c_num, "COUNT(alarm_type = 'leakage' OR NULL) AS num_xl")
			//	z_key = append(z_key, "num_xl")
		}

	}

	var z_numstr = ""
	for i, v := range z_num {
		if i == 0 {
			z_numstr = ","
		}
		z_numstr += v
		if i < len(z_num)-1 {
			z_numstr += ","
		}
	}
	var z_keystr = ""
	for i, v := range z_key {
		z_keystr += v
		if i < len(z_key)-1 {
			z_keystr += "+"
		}
	}
	if len(z_keystr) > 0 {
		z_keystr = ",IFNULL(" + z_keystr + ",0) as num "
	}

	var c_numstr = ""
	for i, v := range c_num {
		if i == 0 {
			c_numstr = ","
		}
		c_numstr += v
		if i < len(c_num)-1 {
			c_numstr += ","
		}
	}

	var sql = "SELECT mm.aa as brief_time " +
		z_numstr + " " +
		z_keystr +
		" FROM  " +
		" ( SELECT DATE_FORMAT(DATE_SUB(CONVERT(?,datetime),INTERVAL xc-1 DAY),'%m.%d') as aa," + //获取xc-1天的日期“月日”
		" DATE_FORMAT( DATE_SUB(CONVERT(?,datetime), INTERVAL xc - 1 DAY), '%Y-%m-%d' ) AS bb " + //获取xc-1天的日期“年月日”
		" FROM" +
		" ( SELECT @xi:=@xi+1 as xc " + //设置xi自加1 赋值给xc
		" FROM " +
		" (SELECT 1 UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5 UNION SELECT 6) xc1, " + //查询6次
		" (SELECT 1 UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5 UNION SELECT 6) xc2, " + //查询6次 在上6次中，共6*6=36次
		" (SELECT @xi:=0) xc0 ) xcxc " + //xi 默认0
		" WHERE " +
		" xc<=DATEDIFF(CONVERT(?,datetime),CONVERT(?,datetime))+1 )mm " + //判断条件 计算两个时间的相差天数
		" LEFT OUTER JOIN " + //向左加入
		"( SELECT " +
		" DATE_FORMAT(alarm_time, '%m.%d') as brief_time " +
		c_numstr +
		" FROM ss_alarm ss " +
		" INNER JOIN ss_camera sc ON sc.camera_id = ss.camera_id " + //关联摄像头 查询对应告警下的摄像头
		" INNER JOIN ss_place sp ON sp.place_id = sc.place_id " +
		" INNER JOIN ss_organize so ON so.id = sp.organize_id " +
		" WHERE " +
		"so.id in ("
	var sqlc = ""
	for i := 0; i < len(ids); i++ {
		sqlc += ",?"
	}
	sqlc = sqlc[1:len(sqlc)]
	//var sqlc = "'" + user.Organize.ID + "'"
	//for i, v := range orgIds {
	//	if i < len(orgIds) {
	//		sqlc += ","
	//	}
	//	sqlc += "'" + v["id"].(string) + "'"
	//}
	//if len(sqlc) == 0 || sqlc == "" {
	//	sqlc = "''"
	//}
	sql += sqlc

	sql += ") AND alarm_time>=? and alarm_time<=? " +
		"GROUP BY brief_time)a " +
		"ON a.brief_time = mm.aa ORDER BY mm.bb asc"

	var maps []orm.Params
	_, err = o.Raw(sql, endTime, endTime, endTime, startTime, ids, startTime, endTime).Values(&maps)
	if err != nil {
		log.Println(err)
		this.Data["json"] = map[string]interface{}{"Code": 0, "Reason": "查询统计信息异常"}
		this.ServeJSON()
		return
	}

	this.Data["json"] = map[string]interface{}{"Code": 1, "Jur": jurs, "AC": countInfo, "Chart": maps}
	this.ServeJSON()
}

//首页算法统计-折线
type PlatFuncLineController struct {
	beego.Controller
}

func (this *PlatFuncLineController) Post() {
	zId, funcType := this.GetString("zId"), this.GetString("type")

	if funcType == "" || len(funcType) == 0 {
		this.Data["json"] = &util.ResultB{0, 0, nil}
		this.ServeJSON()
		return
	}

	userid := this.Ctx.GetCookie("user_id")
	var startTime = time.Now().AddDate(0, 0, -30).Format("2006-01-02") + " 00:00:00"
	var endTime = time.Now().Format("2006-01-02") + " 23:59:59"

	o := O
	var organize models.Organize

	var user models.User
	err := o.QueryTable(new(models.User)).Filter("admin_id", userid).RelatedSel().One(&user)
	if err != nil {
		log.Print("用户查询err:", err)
		this.Data["json"] = &util.ResultB{0, 0, nil}
		this.ServeJSON()
		return
	}

	if user.Organize.ID == "" {
		this.Data["json"] = map[string]interface{}{"Code": 0, "Num": 0, "Reason": "用户组织架构异常"}
		this.ServeJSON()
		return
	}

	var org = user.Organize

	if len(zId) > 0 && user.Organize.ID != zId {
		err = o.QueryTable(new(models.Organize)).Filter("id", zId).One(&organize)
		if err != nil {
			LogsError("查询组织架构信息异常:", err)
			this.Data["json"] = map[string]interface{}{"Code": 0, "Reason": "组织信息异常"}
			this.ServeJSON()
			return
		}

		if InOrganize(&organize, user.Organize) { //user组织是否包含传入org.id的组织
			org = &organize
		}
	}

	orgIds, err := SelOrgIDBelongUserOrg(org)
	if err != nil {
		LogsError("查询用户下属组织异常:", err)

	}

	var ids = []string{user.Organize.ID}
	for _, v := range orgIds {
		ids = append(ids, v["id"].(string))
	}

	var sql = "SELECT mm.aa as brief_time,IFNULL(a.num,0) as num " +
		"FROM  " +
		"( SELECT DATE_FORMAT(DATE_SUB(CONVERT(?,datetime),INTERVAL xc-1 DAY),'%m%d') as aa," +
		" DATE_FORMAT( DATE_SUB(CONVERT(?,datetime), INTERVAL xc - 1 DAY), '%Y-%m-%d' ) AS bb " +
		" FROM" +
		"( SELECT @xi:=@xi+1 as xc " +
		"FROM " +
		"(SELECT 1 UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5 UNION SELECT 6) xc1, " +
		"(SELECT 1 UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5 UNION SELECT 6) xc2, " +
		"(SELECT @xi:=0) xc0 ) xcxc " +
		"WHERE " +
		"xc<=DATEDIFF(CONVERT(?,datetime),CONVERT(?,datetime))+1 )mm " +
		"LEFT OUTER JOIN " +
		"(SELECT " +
		"DATE_FORMAT(alarm_time, '%m%d') as brief_time," +
		"COUNT(*) AS num " +
		"FROM ss_alarm ss " +
		"INNER JOIN ss_camera sc ON sc.camera_id = ss.camera_id " +
		" INNER JOIN ss_place sp ON sp.place_id = sc.place_id " +
		" INNER JOIN ss_organize so ON so.id = sp.organize_id " +
		"WHERE " +
		"so.id in ("
	var sqlc = ""
	for i := 0; i < len(ids); i++ {
		sqlc += ",?"
	}
	sqlc = sqlc[1:len(sqlc)]
	//var sqlc = "'" + org.ID + "'"
	//for i, v := range orgIds {
	//	if i < len(orgIds) {
	//		sqlc += ","
	//	}
	//	sqlc += "'" + v["id"].(string) + "'"
	//
	//}
	//if len(sqlc) == 0 || sqlc == "" {
	//	sqlc = "''"
	//}
	sql += sqlc

	sql += ") AND alarm_time>=? and alarm_time<=? " +
		"AND alarm_type = ? " +
		"GROUP BY brief_time)a " +
		"ON a.brief_time = mm.aa ORDER BY mm.bb asc"
	var maps []orm.Params
	_, err = o.Raw(sql, endTime, endTime, endTime, startTime, ids, startTime, endTime, funcType).Values(&maps)
	if err != nil {
		LogsError("查询统计折线图数据异常:", err)
		this.Data["json"] = map[string]interface{}{"Code": 0, "Reason": "查询信息异常"}
		this.ServeJSON()
		return
	}
	this.Data["json"] = &util.ResultB{1, 0, &maps}
	this.ServeJSON()
}

//首页算法告警列表
type PlatFuncAlarmsController struct {
	beego.Controller
}

func (this *PlatFuncAlarmsController) Post() {

	limits, pages, dateTime, zId, funcType, cameraID :=
		this.GetString("limit"),
		this.GetString("page"),
		this.GetString("dateTime"),
		this.GetString("zId"),
		this.GetString("type"),
		this.GetString("cameraID")
	if funcType == "" || len(funcType) == 0 {
		log.Println("算法类型为空")
		this.Data["json"] = &util.ResultB{0, 0, nil}
		this.ServeJSON()
		return
	}
	start_page, err := strconv.Atoi(pages)
	limitss, err := strconv.Atoi(limits)
	start_pages := (start_page - 1) * limitss
	var startTime = ""
	var endTime = ""

	times := strings.Split(dateTime, " - ")
	if len(dateTime) > 0 || len(times) > 1 {
		startTime = times[0] + " 00:00:00"
		endTime = times[1] + " 23:59:59"

		day := TimeDifferDay(times[0], times[1])

		//超过31天或比对异常
		if day < 0 || day > 31 {
			loc, _ := time.LoadLocation("Local") //获取本地时区
			eTime, err := time.ParseInLocation("2006-01-02 15:04:05", endTime, loc)
			if err != nil {
				//时间转换异常从当前时间查询
				startTime = time.Now().AddDate(0, 0, -30).Format("2006-01-02") + " 00:00:00"
				endTime = time.Now().Format("2006-01-02") + " 23:59:59"
			} else {
				//超过的时间不显示
				startTime = eTime.AddDate(0, 0, -30).Format("2006-01-02") + " 00:00:00"
			}
		}
	} else { // 默认查询当天记录
		startTime = time.Now().AddDate(0, 0, -30).Format("2006-01-02") + " 00:00:00"
		endTime = time.Now().Format("2006-01-02") + " 23:59:59"
	}

	userid := this.Ctx.GetCookie("user_id")
	var maps []models.Alarm
	o := O
	var user models.User
	var organize models.Organize
	err = o.QueryTable(new(models.User)).Filter("admin_id", userid).RelatedSel().One(&user)
	if err != nil {
		log.Print("用户查询err:", err)
		this.Data["json"] = &util.ResultB{0, 0, nil}
		this.ServeJSON()
		return
	}

	if user.Organize.ID == "" {
		this.Data["json"] = map[string]interface{}{"Code": 0, "Reason": "用户组织架构异常"}
		this.ServeJSON()
		return
	}

	var org = user.Organize

	if len(zId) > 0 && user.Organize.ID != zId {
		err = o.QueryTable(new(models.Organize)).Filter("id", zId).One(&organize)
		if err != nil {
			LogsError("查询组织架构信息异常:", err)
			this.Data["json"] = map[string]interface{}{"Code": 0, "Reason": "组织信息异常"}
			this.ServeJSON()
			return
		}

		if InOrganize(&organize, user.Organize) { //user组织是否包含传入org.id的组织
			org = &organize
		}

	}

	orgIds, err := SelOrgIDBelongUserOrg(org)
	if err != nil {
		LogsError("查询用户下属组织异常:", err)

	}
	var ids = []string{org.ID}
	for _, v := range orgIds {
		ids = append(ids, v["id"].(string))
	}

	qs := o.QueryTable(new(models.Alarm)).Filter("alarm_type", funcType).Filter("Camera__Place__Organize__ID__in", ids).Filter("alarm_time__gte", startTime).Filter("alarm_time__lte", endTime).Limit(limitss, start_pages).OrderBy("-alarm_time").RelatedSel(3)
	if len(cameraID) > 0 && cameraID != "0" {
		qs = qs.Filter("Camera__ID", cameraID)
	}
	_, err = qs.All(&maps)
	if err != nil {
		LogsError("查询告警信息异常:", err)
		this.Data["json"] = &util.ResultB{0, 0, nil}
		this.ServeJSON()
		return
	}

	var nums []models.Alarm
	qs2 := o.QueryTable(new(models.Alarm)).Filter("alarm_type", funcType).Filter("Camera__Place__Organize__ID__in", ids).Filter("alarm_time__gte", startTime).Filter("alarm_time__lte", endTime).OrderBy("-alarm_time").RelatedSel(3)
	if len(cameraID) > 0 && cameraID != "0" {
		qs2 = qs2.Filter("Camera__ID", cameraID)
	}
	_, err = qs2.All(&nums)
	if err != nil {
		LogsError("查询告警总数异常:", err)
		this.Data["json"] = &util.ResultB{0, 0, nil}
		this.ServeJSON()
		return
	}

	this.Data["json"] = map[string]interface{}{"Code": 0, "Num": len(nums), "Reason": &maps}
	this.ServeJSON()
}

// 新告警查询
type SelNewAlarmHandler struct {
	beego.Controller
}

func (this *SelNewAlarmHandler) Post() {
	userid := this.Ctx.GetCookie("user_id")
	o := O
	var user models.User
	err := o.QueryTable(new(models.User)).Filter("admin_id", userid).One(&user)
	if err == nil {
		stationlist, err := SelOrgIDBelongUserOrg(user.Organize)
		if err != nil {
			mystruct := map[string]interface{}{"Code": 1, "Num": 0, "Reason": "处理失败"}
			this.Data["json"] = mystruct
			this.ServeJSON()
			return
		}

		ids := []string{}
		for i := 0; i < len(stationlist); i++ {
			ids = append(ids, stationlist[i]["id"].(string))
		}

		var maps []models.Alarm
		num, err := o.QueryTable(new(models.Alarm)).Exclude("alarm_type", "flush_count").Filter("page_read_type", "0").Filter("Camera__Place__Organize__ID__in", ids).RelatedSel("Camera__Place__Organize").All(&maps)
		if err != nil || num <= 0 {
			mystruct := map[string]interface{}{"Code": 1, "Num": 0, "Reason": "处理失败"}
			this.Data["json"] = mystruct
			this.ServeJSON()
			return
		}

		o.QueryTable(new(models.Alarm)).Exclude("alarm_type", "flush_count").Filter("page_read_type", "0").Filter("Camera__Place__Organize__ID__in", ids).RelatedSel("Camera__Place__Organize").Update(orm.Params{
			"page_read_type": "1",
		})

		mystruct := map[string]interface{}{"Code": 0, "Num": num, "Reason": &maps}
		this.Data["json"] = mystruct

	} else {
		mystruct := map[string]interface{}{"Code": 1, "Num": 0, "Reason": "处理失败"}
		this.Data["json"] = mystruct
	}
	this.ServeJSON()
}

// 告警处理
type AlarmHandle struct {
	beego.Controller
}

func (this *AlarmHandle) Post() {
	alarmid := this.GetString("alarmid")
	o := O
	_, err := o.QueryTable(new(models.Alarm)).Filter("ID", alarmid).Update(orm.Params{
		"alarmhandler": "1",
	})
	if err != nil {
		mystruct := map[string]interface{}{"Code": 1, "Reason": "操作失败"}
		this.Data["json"] = mystruct
		this.ServeJSON()
		return
	}
	mystruct := map[string]interface{}{"Code": 0, "Reason": "操作成功"}
	this.Data["json"] = mystruct
	this.ServeJSON()
}

//用户组织下所有摄像头
func UserCameras(organize *models.Organize) []*models.Camera {
	if organize == nil {
		return nil
	}

	places := placeByOrgainze(organize)

	var placeids []string
	for _, v := range places {
		placeids = append(placeids, v.ID)
	}

	o := O
	var maps []*models.Camera
	_, err := o.QueryTable(new(models.Camera)).Filter("place_id__in", placeids).RelatedSel().All(&maps)
	if err != nil {
		LogsError("通过区域ids查询区域信息异常:", err)
	}
	return maps
}

/*
	用户告警分类查询
*/
func UserAlarmClassifyCount(orgids []string, atype string) (int64, error) {

	o := O

	var nowTime = time.Now().Format("2006-01-02") + " 23:59:59"
	var startTime = time.Now().AddDate(0, 0, -30).Format("2006-01-02") + " 00:00:00"

	qs := o.QueryTable(new(models.Alarm)).Filter("Camera__Place__Organize__ID__in", orgids).Filter("alarm_type", atype).Filter("alarm_time__gte", startTime).Filter("alarm_time__lte", nowTime).RelatedSel()
	cnt, err := qs.Count()
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return cnt, nil

}

//计算相差天数
func TimeDifferDay(startTime string, endTime string) int {
	loc, _ := time.LoadLocation("Local") //获取本地时区
	sTime, err := time.ParseInLocation("2006-01-02 15:04:05", startTime+" 00:00:00", loc)
	if err != nil {
		return -1
	}
	eTime, err := time.ParseInLocation("2006-01-02 15:04:05", endTime+" 00:00:00", loc)
	if err != nil {
		return -1
	}
	subM := eTime.Sub(sTime)
	log.Println(subM.Hours() / 24)
	return int(subM.Hours()/24) + 1
}

func titleTypeByParamType(t string) string {
	switch t {
	case "cloths":
		return "zhjc"
	case "fire":
		return "hmjc"
	case "boundary":
		return "qyjc"
	case "smoke":
		return "ywjc"
	case "queue_count":
		return "lgjc"
	case "leakage":
		return "xljc"
	case "sleep_count":
		return "sgjc"

	}
	return ""
}
