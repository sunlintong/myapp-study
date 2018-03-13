package controllers

import (
	"encoding/hex"
	"crypto/md5"
//	"fmt"
//	"log"
	"net"
	"github.com/golang/glog"
	"github.com/astaxie/beego"
	"os"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)
var USERNAME,INTRODUCTION string

const MYSQL_USER string = "root"
const MYSQL_PASSWORD string = "hbuckwpzg"
const MYSQL_IP string = "192.168.34.181"
const DATABASENAME string = "myapp"

type MainController struct {
	beego.Controller
}

type UserController struct{
	beego.Controller
}

type User_ProfileController struct{
	beego.Controller
}

type User_SignupController struct{
	beego.Controller
}

type User_LoginController struct{
	beego.Controller
}

//beego.Controller 拥有很多方法，其中包括
//Init、Prepare、Post、Get、Delete、Head 等方法。
//通过重写的方式来实现这些


func(this *MainController) Get(){
	this.Ctx.WriteString("this is /"+"\n")
}


func(this *UserController) Get(){
	this.Data["content"] = os.Getenv("OEM")+os.Getenv("VER")
	this.TplName = "user.tpl"
}

//func (this *UserController) Profile(){	
//}

func(this *User_ProfileController) Get(){
	this.Ctx.WriteString("username:"+USERNAME+"\n"+"introduction:"+INTRODUCTION+"\n")

	userAgentString := this.Ctx.Request.UserAgent()
	this.Ctx.WriteString("UserAgent:"+userAgentString+"\n")
	glog.Infoln("UserAgent:"+userAgentString)
	hostname,err1:=os.Hostname()
	if err1!=nil{
		this.Ctx.WriteString("hostname出错")
	}else{
		this.Ctx.WriteString("主机名："+hostname+"\n")
	}
//	ipString := this.Ctx.Input.IP()
	ipString := getLocalIP()
	this.Ctx.WriteString("IP:"+ipString)
	glog.Infoln("UserIP:"+ipString)
}


//注册页面
func(this *User_SignupController) Get(){
	this.Data["content"] = os.Getenv("OEM")+os.Getenv("VER")
	this.TplName = "signup.tpl"
}

func(this *User_SignupController) Post(){
	var name,inputPassword,introduction string
	inputs := this.Input()
	name = inputs.Get("username")
	//记录password时用MD5加密
	inputPassword = inputs.Get("password")
	password := GetMd5String(inputPassword)
	introduction = inputs.Get("introduction")

	db,err := sql.Open("mysql",MYSQL_USER+":"+MYSQL_PASSWORD+"@tcp("+MYSQL_IP+":3306)/"+DATABASENAME+"?charset=utf8")
	Checkerr(err)
	rows,err1 := db.Query("SELECT name FROM info WHERE name=?",name)
	Checkerr(err1)
	
	//输入不完整
	if name == ""||inputPassword == "" {		
		this.TplName = "inputInvalid.tpl"
	//用户名重复
	}else if rows.Next() == true{
		this.TplName = "nameRepeat.tpl"
	//注册信息合法，往数据库添加
	}else{
		result,err2 := db.Exec("INSERT INTO info(name,password,introduction) VALUES(?,?,?)",name,password,introduction)
		Checkerr(err2)
		ids,_ := result.LastInsertId()
		glog.Infoln("inserted id :",ids)
		this.TplName = "signupSuccess.tpl"
	}	
}


//登录页面
func(this *User_LoginController) Get(){	
	this.Data["content"] = os.Getenv("OEM")+os.Getenv("VER")
	this.TplName = "login.tpl"
}

func(this *User_LoginController) Post(){
	var name,password string
	inputs := this.Input()
	name = inputs.Get("username")
	//记录用户输入时变作MD5转换
	password = GetMd5String(inputs.Get("password"))

	db,err := sql.Open("mysql",MYSQL_USER+":"+MYSQL_PASSWORD+"@tcp("+MYSQL_IP+":3306)/"+DATABASENAME+"?charset=utf8")
	Checkerr(err)
	rows,err1 := db.Query("SELECT name,password,introduction FROM info WHERE name=?",name)
	Checkerr(err1)
	
//	for rows.Next() {
//		var dbname,dbpassword string
//		err2 := rows.Scan(&dbname,&dbpassword)
//		Checkerr(err2)
//		fmt.Println(dbname,dbpassword)
//	}
	glog.Infoln("try to login")

	if name == ""||password == "" {		
		this.TplName = "inputInvalid.tpl"
	//无论是没有该用户名或者密码错误，都只提示登录失败
	}else{
		//rows.Next()读取指向的下一行，一开始指向的是第一行之前，并返回bool，判断其中是否还有元组
		//如果rows中有元组
		if rows.Next() == true{
			var dbname,dbpassword,dbintroduction string
			err2 := rows.Scan(&dbname,&dbpassword,&dbintroduction)
			Checkerr(err2)
			//登录成功
			if dbname == name && dbpassword == password {
		
//				userController := UserController{}

				glog.Infoln("a user log in")
				USERNAME = dbname
				INTRODUCTION = dbintroduction
//				beego.Router("/user/profile",userController,"*:Profile")
//				url := userController.URLFor("userController.Profile","username",dbname,"introduction",dbintroduction)
//				log.Println(url)
				//可以用Redirect方法直接跳转页面
				this.Redirect("/user/profile",302)
//				this.TplName = "loginSuccess.tpl"
			}else{
				glog.Warningln("password error")
				this.TplName = "loginFailed.tpl"
			}
		//rows中没有元组	
		}else{
			glog.Warningln("dont have the user name")
			this.TplName = "loginFailed.tpl"
		}
	}	
}

func getLocalIP() string{
	addrSlice,err := net.InterfaceAddrs()
	if err != nil{
		glog.Errorln("get local ip failed!")
		return "localhost"
	}
	for _,addr := range addrSlice {
		ipnet,ok := addr.(*net.IPNet);
		if ok && !ipnet.IP.IsLoopback(){
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}		
	}
	return "localhost"
}


func Checkerr(err error){
	if err != nil{
		glog.Errorln(err)
	}
}

//为指定字符串生成MD5字符串
func GetMd5String(s string) string{
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
