package cfg

import (
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
)

type MysqlconfigHost struct {
	Address  string `ini:"address"`
	Port     int    `ini:"port"`
	Username string `ini:"username"`
	Password string `ini:"password"`
	DB       string `ini:"db"`
	Table    string `ini:"table"`
}

type MysqlconfigSlave struct {
	Address  string `ini:"address"`
	Port     int    `ini:"port"`
	Username string `ini:"username"`
	Password string `ini:"password"`
	DB       string `ini:"db"`
	Table    string `ini:"table"`
}

type SqlSyntax struct {
	HostSelcet  string `ini:"host_selcet"`
	SlaveUpdate string `ini:"slave_update"`
	UpdateTime  int    `ini:"update_time"`
	HostCount   string `ini:"host_count"`
	SendCount   int    `ini:"send_count"`
}

type Config struct {
	MysqlconfigHost  `ini:"mysql_host"`
	MysqlconfigSlave `ini:"mysql_slave"`
	SqlSyntax       `ini:"sql_syntax"`
}

//解析配置文件
func Loadini(filename string, data interface{}) (err error) {
	t := reflect.TypeOf(data)
	if t.Kind() != reflect.Ptr {
		err = errors.New("data shoud be a pointer")
		return
	}
	if t.Elem().Kind() != reflect.Struct {
		err = errors.New("data shoud be a struct pointer")
		return
	}
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("read file failed, err:", err)
		return
	}
	linestr := strings.Split(string(content), "\r\n")
	var structName string
	for idx, line := range linestr {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "[") {
			if line[0] != '[' && line[len(line)-1] != ']' {
				err = fmt.Errorf("line:%d syntax error", idx+1)
				return
			}
			sectionName := strings.TrimSpace(line[1 : len(line)-1])
			if len(sectionName) == 0 {
				err = fmt.Errorf("line:%d syntax error", idx+1)
				return
			}
			for i := 0; i < t.Elem().NumField(); i++ {
				field := t.Elem().Field(i)
				if field.Tag.Get("ini") == sectionName {
					structName = field.Name
				}
			}

		} else {
			if !strings.Contains(line, "=") || strings.HasPrefix(line, "=") {
				err = fmt.Errorf("line:%d syntax error", idx+1)
				return
			}

			kv := strings.Split(line, "=")
			v := reflect.ValueOf(data)

			svalue := v.Elem().FieldByName(structName)
			stype := svalue.Type()

			if stype.Kind() != reflect.Struct {
				err = fmt.Errorf("%s shoud be a struct", structName)
				return
			}
			var fieldname string

			for i := 0; i < stype.NumField(); i++ {
				field := stype.Field(i)
				if strings.TrimSpace(kv[0]) == field.Tag.Get("ini") {
					fieldname = field.Name
					break
				}
			}
			confobj := svalue.FieldByName(fieldname)
			switch confobj.Type().Kind() {
			case reflect.String:
				confobj.SetString(strings.TrimSpace(kv[1]))
			case reflect.Int:
				ins, ok := strconv.Atoi(strings.TrimSpace(kv[1]))
				if ok != nil {
					err = fmt.Errorf("line:%d value type error", idx+1)
				}
				confobj.SetInt(int64(ins))
			}

		}

	}
	return
}
