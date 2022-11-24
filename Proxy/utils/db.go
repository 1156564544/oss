package utils

import (
	sql "database/sql"
	_ "github.com/go-sql-driver/mysql"
	"errors"
)

type Users struct{
	Name string
	Password string
	Isroot int
	Isread int
	Iswrite int
}

func SelectFromDB(name string) (*Users,error){
	DB, _ := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/oss")
	if err := DB.Ping(); err != nil {
		return nil,errors.New("open database fail")
	}
	stmt, e := DB.Prepare("select * from users where name=?")
	if e != nil {
		return nil,errors.New("prepare fail")
	}
	query, e := stmt.Query(name)
	if e != nil {
		return nil,errors.New("query fail")
	}
	defer query.Close()
	user:=new(Users)
	if query.Next(){
		e:=query.Scan(&user.Name,&user.Password,&user.Isroot,&user.Isread,&user.Iswrite)
		if e!=nil{
			return nil,errors.New("scan fail")
		}
	}else{
		return nil,errors.New("no data")
	}
	return user,nil
}

func InsertIntoDB(user *Users) error{
	DB, _ := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/oss")
	if err := DB.Ping(); err != nil {
		return errors.New("open database fail")
	}
	stmt, e := DB.Prepare("insert into users values(?,?,?,?,?)")
	if e != nil {
		return errors.New("prepare fail")
	}
	_, e = stmt.Exec(user.Name,user.Password,user.Isroot,user.Isread,user.Iswrite)
	if e != nil {
		return errors.New("exec fail")
	}
	return nil
}

func DeleteFromDB(username string) error{
	DB, _ := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/oss")
	if err := DB.Ping(); err != nil {
		return errors.New("open database fail")
	}
	stmt, e := DB.Prepare("delete from users where name=?")
	if e != nil {
		return errors.New("prepare fail")
	}
	_, e = stmt.Exec(username)
	if e != nil {
		return errors.New("exec fail")
	}
	return nil
}

func UpdateDB(username string ,user *Users) error{
	DB, _ := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/oss")
	if err := DB.Ping(); err != nil {
		return errors.New("open database fail")
	}
	stmt, e := DB.Prepare("update users set password=?,isroot=?,isread=?,iswrite=? where name=?")
	if e != nil {
		return errors.New("prepare fail")
	}
	_, e = stmt.Exec(user.Password,user.Isroot,user.Isread,user.Iswrite,username)
	if e != nil {
		return errors.New("exec fail")
	}
	return nil
}