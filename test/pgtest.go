package main

/*
this is about xorm use postgresql
 */
import (
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/golang/glog"
	_ "github.com/lib/pq"
	"log"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123"
	dbName   = "pg_test"
)

var engine *xorm.Engine

func main() {

	selectAll()

	DeleteTest(21)

	t := &test{20, 20}
	err := Insertto(t)
	if err != true {
		glog.Infof("insert is fail with action %s", err)
	}

	selectTest(23)

	t.Num = 21
	UpdateTest(t)

	selectAll()
}

func init() {
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)
	engine, err = xorm.NewEngine("postgres", psqlInfo)
	if err != nil {
		glog.Infof("this is a err %s", err)
	}

	////同步结构体与数据表
	//if err=engine.Sync(new(test));err!=nil{
	//	glog.Infof("Fail to sync database: %v\n", err)
	//}
}


//自检机制
func getDBEngine() *xorm.Engine {
	err := engine.Ping()
	if err != nil {
		glog.Infof("connect postgresql fail %s", err)
		return nil
	}
	fmt.Println("connect postgresql success")
	return engine
}

type test struct {
	Id  int64
	Num int
}

func dataInit() {
	if err := engine.Sync2(new(test)); err != nil {
		glog.Infof("this is init database is fail %s", err)
	}
}

//查询所有条件
func selectAll() {
	var t []test
	engine := getDBEngine()
	engine.SQL("select * from test").Find(&t)
	fmt.Println(t)
}

//条件查询
func selectTest(id int) {
	var t []test
	engine := getDBEngine()
	engine.Where("test.id=?", id).Find(&t)
	fmt.Println(t)
}

//插入测试
func Insertto(t *test) bool {
	engine := getDBEngine()
	rows, err := engine.Insert(t)
	if err != nil {
		glog.Infof("err abdou insert %s", err)
		return false
	}

	if rows == 0 {
		return false
	}
	return true
}

//删除测试
func DeleteTest(num int) bool {
	t := test{
		Num: num,
	}
	engine := getDBEngine()
	rows, err := engine.Delete(&t)
	if err != nil {
		glog.Infof("delete is fail %s", err)
		return false
	}
	if rows == 0 {
		log.Print("delect not found ")
		return false
	}
	return true
}

func UpdateTest(t *test) bool {
	engine := getDBEngine()
	rows, err := engine.Update(t, test{Id: t.Id})
	if err != nil {
		glog.Infof("updata is fail about %s", err)
		return false
	}
	if rows == 0 {
		glog.Infof("no once is fail")
		return false
	}
	return true
}
