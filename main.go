package main

import (
	"mysqlDatasync/cfg"
	"mysqlDatasync/mylogger"
	"fmt"
	"github.com/jmoiron/sqlx"
	"sync"
	"time"
)

var log mylogger.LoggerIO
var mysqlConfig cfg.Config

func main() {
	defer time.Sleep(time.Second)
	log = mylogger.NewFileLog("info", "./mylogger", "mysql.log", 1024*10*1024) //1K=1024
	//加载配置文件

	err := cfg.Loadini("./cfg/conf.ini", &mysqlConfig)
	if err != nil {
		log.Error("配置文件加载失败:%s", err)
		fmt.Println("配置文件加载失败")
		return
	}
	log.Info("配置文件加载成功:%v", mysqlConfig)
	fmt.Printf("配置文件加载成功:%#v\n", mysqlConfig)

	//链接数据库
	dbHost, dbSlave, err := mysqlConfig.InitDB()
	if err != nil {
		log.Error("连接数据库失败:%s", err)
		fmt.Println("连接数据库失败...")
		return
	}
	log.Info("连接数据库成功")
	fmt.Println("连接数据库成功")
	defer dbHost.Close()
	defer dbSlave.Close()



	//更新一次
	onceupdate(dbHost, dbSlave)

	//定时更新
	// timerupdate(dbHost, dbSlave)

}




func onceupdate(dbHost, dbSlave *sqlx.DB) (err error) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = mysqlConfig.SlaveBak(dbSlave, cfg.DataChan)
		if err != nil {
			log.Error("从库更新失败%s", err)
			fmt.Println("从库更新失败")
			return
		}
		log.Info("从库更新成功")
		fmt.Println("从库更新成功")
	}()

	err = mysqlConfig.HostData(dbHost, cfg.DataChan,cfg.ParseAge)  //配置要清洗的列
	if err != nil {
		log.Error("主库查询失败:%s", err)
		return
	}
	log.Info("主库查询成功")
	wg.Wait()

	return
}

func timerupdate(dbHost, dbSlave *sqlx.DB) (err error) {
	ticker := time.Tick(time.Hour*time.Duration(mysqlConfig.UpdateTime))
	for range ticker {
		datachanl :=make(chan *[]*cfg.User,64)
		err = mysqlConfig.HostData(dbHost, datachanl,cfg.ParseAge)
		if err != nil {
			log.Error("主库查询失败:%s", err)
			return
		}
		log.Info("主库查询成功")
		err = mysqlConfig.SlaveBak(dbSlave, datachanl)
		if err != nil {
			log.Error("从库更新失败%s", err)
			fmt.Println("从库更新失败")
			return
		}
		log.Info("从库更新成功")
		fmt.Println("从库更新成功")
	}
	return
}
