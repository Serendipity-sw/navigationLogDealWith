package main

import (
	"flag"
	"github.com/guotie/config"
	"github.com/guotie/deferinit"
	"github.com/smtc/glog"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
)

var (
	configFn        = flag.String("config", "./config.json", "config file path")
	debugFlag       = flag.Bool("d", false, "debug mode")
	sendMailS       []string //需要发送的email地址
	account         string   //发送方email账号
	password        string   //发送方email密码
	processFilePath string   //需要分析的目录
	day             int64    //当前日期减几天
	timeDatef       string   //程序运行需要处理的日期 时间格式为2006/01/02
	timeDate        string   //程序运行需要处理的日期 时间格式为20060102
	sendMailContent []string
	taskId          string //需要检查的任务ID
)

/**
服务启动预处理方法
创建人:邵炜
创建时间:2016年7月18日10:58:13
*/
func serverRun(cfn string, debug bool) {
	config.ReadCfg(cfn)
	logInit(debug)
	sendMails := config.GetString("sendMails")
	if len(sendMails) > 0 {
		sendMailS = strings.Split(sendMails, ",")
	}
	account = config.GetString("account")
	password = config.GetString("password")
	processFilePath = config.GetString("processFilePath")
	taskId = config.GetString("taskId")
	day = config.GetInt64Default("day", -1)
	// 初始化
	deferinit.InitAll()
	glog.Info("init all module successfully.\n")
	// 设置多cpu运行
	runtime.GOMAXPROCS(runtime.NumCPU())
	deferinit.RunRoutines()
	glog.Info("run routines successfully.\n")
}

/**
服务退出
创建人:邵炜
创建时间:2016年7月18日11:11:30
*/
func serverExit() {
	// 结束所有go routine
	deferinit.StopRoutines()
	glog.Info("stop routine successfully.\n")

	deferinit.FiniAll()
	glog.Info("fini all modules successfully.\n")
}

func main() {
	//判断进程是否存在
	if checkPid() {
		return
	}
	flag.Parse()
	serverRun(*configFn, *debugFlag)
	c := make(chan os.Signal, 1)
	writePid()
	// 信号处理
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	// 等待信号
	<-c
	serverExit()
	rmPidFile()
	glog.Close()
	os.Exit(0)
}
