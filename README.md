# myapp
## 应用简介
myapp是一个基于**beego**框架的go语言web服务端程序，主要实现了如下功能:  
> * 与客户端通信，让客户端回显其消息  
> * 定时向客户端推送时间消息  
> * 在web页面上实现登录、注册、展示用户信息及浏览器相关信息  


## myapp页面功能介绍

### 1. 用户注册
 * 用户在/user页面点击注册  
![/user](https://github.com/sunlintong/myapp/blob/master/picture/1.PNG)  
 * 进入/user/signup页面注册  
![/user/signup](https://github.com/sunlintong/myapp/blob/master/picture/2.PNG)  
 * 注册过程中，密码转换为其**MD5编码**再保存至数据库，注册成功后可见数据库中插入了一条信息  
![database info](https://github.com/sunlintong/myapp/blob/master/picture/3.PNG)  
注册成功后跳转至登录界面  

### 2. 登录
 * 在登录界面输入用户名和密码即可登录  
![/user/login](https://github.com/sunlintong/myapp/blob/master/picture/4.PNG)  
 * 登录时，将用户输入的密码转换为MD5编码后再与数据库中的数据匹配，若相同则跳转至profile界面  
![/user/profile](https://github.com/sunlintong/myapp/blob/master/picture/5.PNG)  
profile界面展示：  
> * 登录用户的用户名
> * 登录用户的注册信息：introduction
> * UserAgent
> * 在访问myapp页面的主机名
> * 在访问myapp页面的主机ip

### 3. 与客户端通信
通过**grpc**实现：  
 * 定时向客户端推送含有当前时间的消息：  
![push time](https://github.com/sunlintong/myapp/blob/master/picture/6.PNG)  
 * echo：回显客户端的输入消息  
![echo](https://github.com/sunlintong/myapp/blob/master/picture/7.PNG)  

## 接口说明
* 软件名和版本号通过在Dockerfile里设置环境变量OEM和VER实现：  
````Dockerfile
ENV OEM myapp

ENV VER 5.1.3
````
* getLocalIP()获取IP  
````golang
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
````
* GetMd5String()函数，获取指定字符串的MD5字符串
用于在用户注册时用此函数将密码加密后存储和登录时将密码加密后与数据库中密码配对  
```golang
func GetMd5String(s string) string{
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
```
* Checkerr()检查并打印错误日志
```golang
func Checkerr(err error){
	if err != nil{
		glog.Errorln(err)
	}
}
```

## 在EKOS集群上部署myapp应用
### 1. 安装部署EKOS集群
* 在集群外的主机上下载安装包并`tar zxvf`解压
* `./`执行脚本安装
* `ekoslet`进入集群配置操作环境
* 配置集群的master、node，我的集群新建时只有一台master和一台node主机
`inventory init master:192.168.34.71:etcd:192.168.34.71:node:192.168.34.72-73`etcd应和master的ip一致
* 配置分布式文件系统，我选择将ceph安装在master上，但之后使用EKOS的时候，集群的存储时常出问题，需要`systemctl restart ceph.target`甚至重启集群才能解决，最好是装在node上
`ceph init rgw:192.168.34.72:mon:192.168.34.72:osd:192.168.34.72`
* 分发公钥`keygen +password`
* 安装ceph `ceph install`
* 再`install`
* 进入[我的EKOS](http://192.168.34.71:30000)  

### 2. 添加mysql
* 从本地将mysql官方镜像push至EKOS  
* 添加环境变量`MYSQL_ROOT_PASSWORD=<mypassword>`设置我的mysql密码  
* 添加负载均衡，转发规则，完成后可通过转发地址访问mysql  
![转发](/picture/8.PNG)  
![mysql](/picture/9.PNG)  

### 3. 添加myapp
* 用Makefile和Dockerfile编译并生成myapp镜像  
* push至EKOS  
* 添加启动命令 `./myapp -logtostderr=true`,所有日志均输出至标准控制台  
* 添加负载均衡，转发规则，完成后可通过转发地址访问myapp  
![转发](/picture/10.PNG)  
![myapp](/picture/11.PNG)  

### 4. 激活日志、监控
1. 在master上为node4打标签，准备用node4做日志功能  
`kubectl label nodes node4 ekos.ghostcloud.cn/label-role=logging`   
`kubectl taint nodes nodeX ekos.ghostcloud.cn/taint-role=logging:NoExecute`  
   
2. 在node4上： 设置系统参数  
`ulimit -n unlimited`  
`ulimit -l unlimited`  
`ulimit -s unlimited`  
3. 在master上为node5打标签，准备用node5做监控功能   
`kubectl label nodes node5 ekos.ghostcloud.cn/label-role=moniter`   
`kubectl taint nodes nodeX ekos.ghostcloud.cn/taint-role=moniter:NoExecute`  
4.进入[我的EKOS](http://192.168.34.71:30000)激活  
## 问题和解决方法

### 1. myapp-cli连接断开后，myapp也停止运行
经排查问题出现在下部分代码：  
````golang
for true{
		
	echorequest,err := stream.Recv()
				
//	if echorequest == nil || err !=nil {
			
//		break
		
//	}
		
	echoreply := &pb.EchoReply{}
		
	echoreply.Output = "echo:"+echorequest.Input
		
	stream.Send(echoreply)
		
	glog.Infoln("send a echo message:"+echoreply.Output)
	
}
````
因为如果流里面没消息进来，echoreply中两个参数都是空字符串，echoreply也就是空指针了，发送空指针导致发送错误，加上注释部分就成功了  
### 2. 不登陆，直接输入 /user/profile的链接时，页面会显示上一次登陆用户的用户名和个人介绍
我在`controllers`里用的全局变量保存用户名和个人介绍导致了上述问题    
````golang
var USERNAME,INTRODUCTION string
````
改用session缓存后，问题解决
### 3. ceph安装在了master上，存储容易死掉
### 4. mysql镜像要设置 MYSQL_ROOT_PASSWORD环境变量

### 5. myapp-cli从服务器断开后，服务器仍会一直打印推送时间消息的日志，表明服务器还在试图向客户端推送消息，此bug还未解决


## 几个kubectl命令及结果

### 1. 扩展副本`kubectl scale --replicas=4 deployment/myapp`  
扩展前：  
```
[root@node1 ~]# kubectl get pods
NAME                                    READY     STATUS    RESTARTS   AGE
default-http-backend-3138300093-cssvv   1/1       Running   0          1d
lb-myapp-1719601925-vcvs1               1/1       Running   0          21h
lb-mysql-3406302119-0jfm8               1/1       Running   0          1d
myapp-4112844737-vhmxd                  1/1       Running   0          21h
mysql-0                                 1/1       Running   0          1d
```  
扩展后：    
```
[root@node1 ~]# kubectl scale --replicas=4 deployment/myapp
deployment "myapp" scaled
[root@node1 ~]# kubectl get pods
NAME                                    READY     STATUS              RESTARTS   AGE
default-http-backend-3138300093-cssvv   1/1       Running             0          1d
lb-myapp-1719601925-vcvs1               1/1       Running             0          21h
lb-mysql-3406302119-0jfm8               1/1       Running             0          1d
myapp-4112844737-72zt5                  1/1       Running             0          4s
myapp-4112844737-nj8zh                  0/1       ContainerCreating   0          4s
myapp-4112844737-rdxn1                  1/1       Running             0          4s
myapp-4112844737-vhmxd                  1/1       Running             0          21h
mysql-0                                 1/1       Running             0          1d
[root@node1 ~]#
```
### 2. 查看pod的node
`kubectl get pod -o wide` 

```
[root@node1 ~]# kubectl get pods -o wide
NAME                                    READY     STATUS    RESTARTS   AGE       IP              NODE
default-http-backend-3138300093-cssvv   1/1       Running   0          1d        10.233.71.3     node3
lb-myapp-1719601925-vcvs1               1/1       Running   0          21h       192.168.34.72   node2
lb-mysql-3406302119-0jfm8               1/1       Running   0          1d        192.168.34.73   node3
myapp-4112844737-72zt5                  1/1       Running   0          1m        10.233.71.10    node3
myapp-4112844737-nj8zh                  1/1       Running   0          1m        10.233.75.23    node2
myapp-4112844737-rdxn1                  1/1       Running   0          1m        10.233.75.22    node2
myapp-4112844737-vhmxd                  1/1       Running   0          21h       10.233.71.4     node3
mysql-0                                 1/1       Running   0          1d        10.233.71.2     node3
[root@node1 ~]#

```
或`kubectl get pod myapp-4112844737-vhmxd -o yaml |grep node` 
  
```
[root@node1 ~]# kubectl get pod myapp-4112844737-vhmxd -o yaml |grep node
  nodeName: node3
[root@node1 ~]#
```
 
### 3. 实时查看日志
`kubectl logs -f <pod name>` 
由于一直在用的pod "myapp-4112844737-vhmxd" 已经吐了一天的日志，量太多找不到开始吐的时候，我就把这个pod删掉了。。。

```
[root@node1 ~]# kubectl get po myapp-4112844737-vhmxd -o yaml | kubectl delete -f -
pod "myapp-4112844737-vhmxd" deleted
```


```
[root@node1 ~]# kubectl logs -f myapp-4112844737-72zt5
2017/11/17 03:07:13 [I] [asm_amd64.s:2337] http server Running on http://:8080
2017/11/17 03:44:44 [D] [server.go:2619] |    10.233.75.0| 200 |    147.932µs|   match| GET      /     r:/
2017/11/17 03:44:49 [D] [server.go:2619] |    10.233.75.0| 200 |   1.514156ms|   match| GET      /user   r:/user
2017/11/17 03:44:54 [D] [server.go:2619] |    10.233.75.0| 200 |    851.102µs|   match| GET      /user/login   r:/user/login
I1117 03:45:00.487749       1 default.go:140] try to login
W1117 03:45:00.487856       1 default.go:172] dont have the user name
2017/11/17 03:45:00 [D] [server.go:2619] |    10.233.75.0| 200 |  19.241458ms|   match| POST     /user/login   r:/user/login
2017/11/17 03:45:06 [D] [server.go:2619] |    10.233.75.0| 200 |    831.562µs|   match| GET      /user/login   r:/user/login
I1117 03:45:13.398356       1 default.go:140] try to login
I1117 03:45:13.398398       1 default.go:157] a user log in
2017/11/17 03:45:13 [D] [server.go:2619] |    10.233.75.0| 302 |   3.816073ms|   match| POST     /user/login   r:/user/login
I1117 03:45:13.408510       1 default.go:65] UserAgent:Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36
I1117 03:45:13.409126       1 default.go:75] UserIP:10.233.71.10
2017/11/17 03:45:13 [D] [server.go:2619] |    10.233.75.0| 200 |    711.674µs|   match| GET      /user/profile   r:/user/profile
2017/11/17 03:58:27 [D] [server.go:2619] |    10.233.75.0| 404 |    621.173µs| nomatch| GET      /use
2017/11/17 03:58:33 [D] [server.go:2619] |    10.233.75.0| 200 |     70.554µs|   match| GET      /     r:/
2017/11/17 04:02:05 [D] [server.go:2619] |    10.233.75.0| 200 |    870.622µs|   match| GET      /user   r:/user
2017/11/17 04:02:17 [D] [server.go:2619] |    10.233.75.0| 200 |    682.862µs|   match| GET      /user/signup   r:/user/signup
I1117 04:02:42.478442       1 default.go:110] inserted id : 0
2017/11/17 04:02:42 [D] [server.go:2619] |    10.233.75.0| 200 |  27.910728ms|   match| POST     /user/signup   r:/user/signup
2017/11/17 04:02:43 [D] [server.go:2619] |    10.233.75.0| 200 |    769.526µs|   match| GET      /user/login   r:/user/login
2017/11/17 04:02:55 [D] [server.go:2619] |    10.233.75.0| 302 |  12.878449ms|   match| POST     /user/login   r:/user/login
I1117 04:02:55.470620       1 default.go:140] try to login
I1117 04:02:55.470669       1 default.go:157] a user log in
I1117 04:02:55.476231       1 default.go:65] UserAgent:Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36
2017/11/17 04:02:55 [D] [server.go:2619] |    10.233.75.0| 200 |    485.248µs|   match| GET      /user/profile   r:/user/profile
I1117 04:02:55.476654       1 default.go:75] UserIP:10.233.71.10

```