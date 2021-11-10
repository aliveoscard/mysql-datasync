package cfg

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func (c *Config) InitDB() (dbHost, dbSlave *sqlx.DB, err error) {

	dsnHost := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", c.MysqlconfigHost.Username, c.MysqlconfigHost.Password, c.MysqlconfigHost.Address, c.MysqlconfigHost.Port, c.MysqlconfigHost.DB)
	dbHost, err = sqlx.Connect("mysql", dsnHost)
	if err != nil {
		return
	}
	fmt.Printf("主数据库连接成功(%s)...\n", c.MysqlconfigHost.Address)

	dsnSlave := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", c.MysqlconfigSlave.Username, c.MysqlconfigSlave.Password, c.MysqlconfigSlave.Address, c.MysqlconfigSlave.Port, c.MysqlconfigSlave.DB)
	dbSlave, err = sqlx.Connect("mysql", dsnSlave)
	if err != nil {
		return
	}
	fmt.Printf("从数据库连接成功(%s)...\n", c.MysqlconfigSlave.Address)

	return
}

//查询主库
func (c *Config) HostData(dbHost *sqlx.DB, dataChanel chan<- *[]*User,parse ...func(*User)) (err error) {
	var totalCount int
	err = dbHost.Get(&totalCount, c.HostCount) //查询总条数
	if err != nil {
		return
	}
	fmt.Println("总条数", totalCount)
	percount := c.SendCount
	ss := totalCount / percount
	if totalCount%percount != 0 {
		ss++
	}
	for i := 0; i < ss; i++ {
		res := []*User{}
		err = dbHost.Select(&res, fmt.Sprintf(c.HostSelcet, i*percount, percount))
		if err != nil {
			fmt.Println(err)
			return
		}
		
		for _,dataClean := range res {
			for _,p := range parse {
				p(dataClean)
			}
		}
		dataChanel <- &res
	}

	fmt.Println("主库查询完成")
	close(dataChanel)
	return
}

func (c *Config) SlaveBak(dbSlave *sqlx.DB, dataChanel <-chan *[]*User) (err error) {
	//开启事务
	tx, err := dbSlave.Beginx()
	if err != nil {
		fmt.Printf("begin trans failed, err:%v\n", err)
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			fmt.Println("回滚成功")
			tx.Rollback()
		} else {
			err = tx.Commit()
			fmt.Println("事务提交成功")
		}
	}()

	//清空从库的表
	_, err = tx.Exec(fmt.Sprintf("delete from %s", c.MysqlconfigSlave.Table)) //truncate table %s 不会回滚
	if err != nil {
		fmt.Printf("delete failed, err:%v\n", err)
		return
	}
	fmt.Println("正在清空从库成功")
	//插入语句
	for v := range dataChanel {
		_, err = tx.NamedExec(c.SlaveUpdate, *v)
		if err != nil {
			return
		}
	}
	fmt.Println("正在更新从库")
	return
}
