# mysql-datasync
backup mysql database

#### 介绍
mysql数据库全量更新

#### 安装教程

1.  go get github.com/go-sql-driver/mysql
2.  go get github.com/jmoiron/sqlx


#### 使用说明

1.  配置cfg文件中的conf.ini文件
2.  更改cfg文件中table_struct文件中的User对应的数据库结构
3.  mian函数配置需要清洗的列
4.  go run main.go


#### 功能

1.  支持定时更新和及时更新
2.  从库更新前支持事务回滚功能
3.  log日志记录和自动分割功能
4.  支持对数据进行清洗
