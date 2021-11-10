package cfg

import (
	"database/sql"
)

type User struct {
	ID   int            `db:"id"`
	Name string         `db:"names"`
	Age  sql.NullInt64 `db:"age"`
	Src  string         `db:"src"`
}

var DataChan = make(chan *[]*User, 64)

//对age列数据进行清洗
func ParseAge(u *User) {
	res := &(u.Age)
	// if res != nil {  
	// 	switch *res {
	// 	case "1":
	// 		*u.Age = "10"
	// 	case "2":
	// 		*u.Age = "20"
	// 	}
	// } else {

	// 	u.Age = new(string)
	// 	*u.Age = "0"
	// }
	if res.Valid {
		switch res.Int64 {
		case 1:
			res.Int64 = 100
		case 2:
			res.Int64 = 200
		}
	}else{
		res.Valid=true
		res.Int64=0
	}
	
}
