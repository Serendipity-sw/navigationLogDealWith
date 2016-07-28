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
	portalArray                    map[string]int //某个探针设备编号共发送请求数
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
	portalArray = map[string]int{}
	getMyReceiveNumber()
	javaRequestCountNumber = 0
	javaRequestCountNumberByTaskId = 0
	feedBackNumber = 0
	feedBackNumberByTaskId = 0
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
	if strings.TrimSpace(taskId) != "" {
		where = append(where, feedBackNumberByTaskIdFunc)
		where = append(where, javaRequestSuccessByTaskIdFunc)
	}
	where = append(where, feedBackNumberProcessFunc)
	where = append(where, javaRequestSuccessFunc)
	infoFilePath := fmt.Sprintf("%s/logs/INFO-%s.log", processFilePath, timeDate)
	readFile(infoFilePath, where, nil)
	sendMailContent = append(sendMailContent, fmt.Sprintf("java共返回导航数: %d \r\n", javaRequestCountNumber))
	sendMailContent = append(sendMailContent, fmt.Sprintf("java共返回导航数用户数: %d \r\n", len(javaRequestPhoneNumber)))
	if strings.TrimSpace(taskId) != "" {
		sendMailContent = append(sendMailContent, fmt.Sprintf("java共返回任务%s: %d \r\n", taskId, javaRequestCountNumberByTaskId))
		sendMailContent = append(sendMailContent, fmt.Sprintf("java共返回任务%s用户数: %d \r\n", taskId, len(javaRequestPhoneNumberByTaskId)))
		sendMailContent = append(sendMailContent, fmt.Sprintf("曝光总量任务%s: %d \r\n", taskId, feedBackNumberByTaskId))
		sendMailContent = append(sendMailContent, fmt.Sprintf("曝光总量任务%s用户数: %d \r\n", taskId, len(feedBackPhoneNumberByTaskId)))
	}
	sendMailContent = append(sendMailContent, fmt.Sprintf("曝光总量: %d \r\n", feedBackNumber))
	sendMailContent = append(sendMailContent, fmt.Sprintf("曝光总量用户数: %d \r\n", len(feedBackPhoneNumber)))
	for key, value := range portalArray {
		sendMailContent = append(sendMailContent, fmt.Sprintf("探针编号%s: %d \r\n", key, value))
	}
}

/**
探针各设备编号收集
创建人:邵炜
创建时间:2016年7月18日16:15:00
*/
func portalCount(content string) bool {
	bo := strings.Index(content, "parseParams:") >= 0 && strings.Index(content, "pid=") >= 0
	if bo {
		pidStr := strings.Split(content, "&")[4]
		pidNumber := strings.Split(pidStr, "=")[1]
		_, ok := portalArray[pidNumber]
		if ok {
			portalArray[pidNumber]++
		} else {
			portalArray[pidNumber] = 0
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
	sendMailContent = append(sendMailContent, fmt.Sprintf("省份: %s \r\n", province))
	sendMailContent = append(sendMailContent, fmt.Sprintf("导航共接收探针%d次请求 \r\n", matchNumber))
}

/**
nohup文件筛选判断条件
创建人:邵炜
创建时间:2016年7月18日14:49:12
*/
func nohupJudgment(content string) bool {
	return strings.Index(content, "mainjs/analysis") >= 0 && strings.Index(content, timeDatef) >= 0
}
