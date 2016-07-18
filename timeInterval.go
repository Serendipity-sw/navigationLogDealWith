package main

import (
	"github.com/guotie/deferinit"
	"github.com/smtc/glog"
	"sync"
	"time"
)

var (
	jsTmr *time.Timer
)

func init() {
	deferinit.AddRoutine(watchFilesProcess)
}

/**
定时执行程序,每日定点处理日志文件
创建人:邵炜
创建时间:2016年7月18日11:19:15
*/
func watchFilesProcess(ch chan struct{}, wg *sync.WaitGroup) {
	go func() {
		<-ch

		jsTmr.Stop()
		wg.Done()
	}()

	jsTmr = time.NewTimer(getMyEveyDayFourTime())
	glog.Info(" watch userBlackList is waiting! \n")
	<-jsTmr.C
	glog.Info(" watch userBlackList is loading! \n")
	for {
		logProcessFunc()
		jsTmr.Reset(getMyEveyDayFourTime())
		<-jsTmr.C
	}
}

/**
获取每天凌晨4点到现在时间是多少小时
创建人:邵炜
创建时间:2016年3月28日11:33:05
*/
func getMyEveyDayFourTime() time.Duration {
	return time.Duration(28-time.Now().Hour()) * time.Hour
}
