package main

import (
	"fmt"
	"github.com/smtc/glog"
	"net/smtp"
	"strings"
	"time"
)

var (
	javaRequestCountNumber         int            //java共返回导航数
	javaRequestCountNumberByTaskId int            //java共返回某个特定任务的导航树
	feedBackNumber                 int            //所有任务共曝光次数
	feedBackNumberByTaskId         int            //某个特定任务的曝光次数
	portal11                       int            //某个探针设备编号共发送请求书
	portal12                       int            //某个探针设备编号共发送请求书
	portal13                       int            //某个探针设备编号共发送请求书
	portal14                       int            //某个探针设备编号共发送请求书
	portal209                      int            //某个探针设备编号共发送请求书
	portal210                      int            //某个探针设备编号共发送请求书
	portal6                        int            //某个探针设备编号共发送请求书
	portal7                        int            //某个探针设备编号共发送请求书
	javaRequestPhoneNumber         map[string]int //java共返回需要导航的用户数
	javaRequestPhoneNumberByTaskId map[string]int //java共返回某个特定任务的导航用户数
	feedBackPhoneNumber            map[string]int //所有任务的导航用户数
	feedBackPhoneNumberByTaskId    map[string]int //某个特定任务的导航用户数
)

/**
日志处理主方法
创建人:邵炜
创建时间:2016年7月18日11:36:26
*/
func logProcessFunc() {
	timeDates := time.Now().Add(time.Duration(day) * 24 * time.Hour)
	timeDatef = timeDates.Format("2006/01/02")
	timeDate = timeDates.Format("20060102")
	sendMailContent = sendMailContent[:0]
	getMyReceiveNumber()
	javaRequestCountNumber = 0
	javaRequestCountNumberByTaskId = 0
	feedBackNumber = 0
	feedBackNumberByTaskId = 0
	portal11 = 0
	portal12 = 0
	portal13 = 0
	portal14 = 0
	portal209 = 0
	portal210 = 0
	portal6 = 0
	portal7 = 0
	javaRequestPhoneNumber = make(map[string]int)
	javaRequestPhoneNumberByTaskId = make(map[string]int)
	feedBackPhoneNumber = make(map[string]int)
	feedBackPhoneNumberByTaskId = make(map[string]int)
	infoFileProcess()
	sendMail()
}

func sendMail() {
	auth := smtp.PlainAuth("", account, password, "smtp.qiye.163.com")
	to := sendMailS
	msg := []byte("To: bain@axon.com.cn;shaow@axon.com.cn\r\n" +
		"Subject: Navigation Of Task Daily Analysis-" + timeDatef + "\r\n" +
		"Content-Type: text/plain; charset=UTF-8" + "\r\n\r\n")
	for _, value := range sendMailContent {
		msg = append(msg, []byte(value)...)
	}
	err := smtp.SendMail("smtp.qiye.163.com:25", auth, account, to, msg)
	if err != nil {
		glog.Error("邮件发送失败! err: %s \n", err.Error())
		return
	}
	glog.Info("邮件发送成功! \n")
}

/**
info日志分析主方法
创建人:邵炜
创建时间:2016年7月18日16:48:56
*/
func infoFileProcess() {
	var where []func(string) bool
	where = append(where, portalCount)
	where = append(where, feedBackNumberByTaskIdFunc)
	where = append(where, feedBackNumberProcessFunc)
	where = append(where, javaRequestSuccessFunc)
	where = append(where, javaRequestSuccessByTaskIdFunc)
	infoFilePath := fmt.Sprintf("%s/logs/INFO-%s.log", processFilePath, timeDate)
	readFile(infoFilePath, where, nil)
	sendMailContent = append(sendMailContent, fmt.Sprintf("java共返回导航数: %d \r\n", javaRequestCountNumber))
	sendMailContent = append(sendMailContent, fmt.Sprintf("java共返回导航数用户数: %d \r\n", len(javaRequestPhoneNumber)))
	sendMailContent = append(sendMailContent, fmt.Sprintf("java共返回任务%s: %d \r\n", taskId, javaRequestCountNumberByTaskId))
	sendMailContent = append(sendMailContent, fmt.Sprintf("java共返回任务%s用户数: %d \r\n", taskId, len(javaRequestPhoneNumberByTaskId)))
	sendMailContent = append(sendMailContent, fmt.Sprintf("曝光总量: %d \r\n", feedBackNumber))
	sendMailContent = append(sendMailContent, fmt.Sprintf("曝光总量用户数: %d \r\n", len(feedBackPhoneNumber)))
	sendMailContent = append(sendMailContent, fmt.Sprintf("曝光总量任务%s: %d \r\n", taskId, feedBackNumberByTaskId))
	sendMailContent = append(sendMailContent, fmt.Sprintf("曝光总量任务%s用户数: %d \r\n", taskId, len(feedBackPhoneNumberByTaskId)))
	sendMailContent = append(sendMailContent, fmt.Sprintf("探针编号11: %d \r\n", portal11))
	sendMailContent = append(sendMailContent, fmt.Sprintf("探针编号12: %d \r\n", portal12))
	sendMailContent = append(sendMailContent, fmt.Sprintf("探针编号13: %d \r\n", portal13))
	sendMailContent = append(sendMailContent, fmt.Sprintf("探针编号14: %d \r\n", portal14))
	sendMailContent = append(sendMailContent, fmt.Sprintf("探针编号209: %d \r\n", portal209))
	sendMailContent = append(sendMailContent, fmt.Sprintf("探针编号210: %d \r\n", portal210))
	sendMailContent = append(sendMailContent, fmt.Sprintf("探针编号6: %d \r\n", portal6))
	sendMailContent = append(sendMailContent, fmt.Sprintf("探针编号7: %d", portal7))
}

/**
探针各设备编号收集
创建人:邵炜
创建时间:2016年7月18日16:15:00
*/
func portalCount(content string) bool {
	bo := strings.Index(content, "parseParams:") >= 0 && strings.Index(content, "pid=") >= 0
	if bo {
		if strings.Index(content, "pid=11") >= 0 {
			portal11++
		} else if strings.Index(content, "pid=12") >= 0 {
			portal12++
		} else if strings.Index(content, "pid=13") >= 0 {
			portal13++
		} else if strings.Index(content, "pid=14") >= 0 {
			portal14++
		} else if strings.Index(content, "pid=209") >= 0 {
			portal209++
		} else if strings.Index(content, "pid=210") >= 0 {
			portal210++
		} else if strings.Index(content, "pid=6") >= 0 {
			portal6++
		} else if strings.Index(content, "pid=7") >= 0 {
			portal7++
		}
	}
	return bo
}

/**
根据ID查询该任务的曝光量
创建人:邵炜
创建时间:2016年7月18日15:37:00
*/
func feedBackNumberByTaskIdFunc(content string) bool {
	bo := strings.Index(content, "feedback:") >= 0 && strings.Index(content, "activeid=0") >= 0 && strings.Index(content, fmt.Sprintf("task=%s", taskId)) >= 0
	if bo {
		feedBackNumberByTaskId++
		feedBackPhoneNumberByTaskId[strings.Split(content, " ")[5]] = 0
	}
	return bo
}

/**
曝光量总量
创建人:邵炜
创建时间:2016年7月18日15:35:18
*/
func feedBackNumberProcessFunc(content string) bool {
	bo := strings.Index(content, "feedback:") >= 0 && strings.Index(content, "activeid=0") >= 0
	if bo {
		feedBackNumber++
		feedBackPhoneNumber[strings.Split(content, " ")[5]] = 0
	}
	return bo
}

/**
java请求返回成功数
创建人:邵炜
创建时间:2016年7月18日15:21:20
*/
func javaRequestSuccessFunc(content string) bool {
	bo := strings.Index(content, "getScene: success:") >= 0
	if bo {
		javaRequestCountNumber++
		javaRequestPhoneNumber[strings.Split(content, " ")[4]] = 0
	}
	return bo
}

/**
判断java请求返回匹配为某个任务ID数
创建人:邵炜
创建时间:2016年7月18日15:25:48
*/
func javaRequestSuccessByTaskIdFunc(content string) bool {
	bo := strings.Index(content, "getScene: success:") >= 0 && strings.Index(content, fmt.Sprintf("task=%s", taskId)) >= 0
	if bo {
		javaRequestCountNumberByTaskId++
		javaRequestPhoneNumberByTaskId[strings.Split(content, " ")[4]] = 0
	}
	return bo
}

/**
获取nohup文件共共接收多少数据量
创建人:邵炜
创建时间:2016年7月18日14:48:42
*/
func getMyReceiveNumber() {
	nohupFilePath := fmt.Sprintf("%s/nohup.out", processFilePath)
	var where []func(string) bool
	where = append(where, nohupJudgment)
	matchNumber, _ := readFile(nohupFilePath, where, nil)
	sendMailContent = append(sendMailContent, fmt.Sprintf("导航共接收探针%d次请求", matchNumber))
}

/**
nohup文件筛选判断条件
创建人:邵炜
创建时间:2016年7月18日14:49:12
*/
func nohupJudgment(content string) bool {
	return strings.Index(content, "mainjs/analysis") >= 0 && strings.Index(content, timeDatef) >= 0
}
