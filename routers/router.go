package routers

import (
	"myapp/controllers"
	"github.com/astaxie/beego"
)


func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/user",&controllers.UserController{})
	beego.Router("/user/profile",&controllers.User_ProfileController{})
	beego.Router("/user/signup",&controllers.User_SignupController{})
	beego.Router("/user/login",&controllers.User_LoginController{})
	
}
