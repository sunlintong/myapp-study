# kubernetes命令
## kubectl get
* `kubectl get po`获取当前运行的所有pods的信息  
![kubectl get po](/kubectl_picture/1.PNG)  
可以看到集群中一共跑着5个pod，第一个不清楚是什么，第二、三个是mysql和myapp的负载均衡，最后两个是他们的服务，如果要获取pod运行在哪个节点上的信息，可用`kubectl get po -o wide`  
![kubectl get po -o wide](/kubectl_picture/2.PNG)  
* `kubuctl get namespace`获取namespace信息
* `kubectl get rc`, `kubectl get svc`, `kubectl get nodes`等获取其他resource信息  
* 可以这样组合使用：`kubectl get po --namespace=kube-system`  
![kubectl get ](/kubectl_picture/3.PNG)  

## kubectl describe
`kubectl describe`类似于get，同样用于获取resource的相关信息。不同的是，get获得的是更详细的resource个性的详细信息，**describe获得的是resource集群相关的信息**如`kubectl describe po mysql-server-0`

## kubectl create
* `kubectl create`用于根据文件或输入创建集群resource,如果已经定义了相应resource的yaml或son文件，直接`kubectl create -f filename`即可创建文件内定义的resource  
* 直接使用`kubectl create`则可以基于yaml文件创建出ReplicationController（rc），rc会创建两个副本  

## kubectl replace
`kubectl create`命令用于对已有资源进行更新、替换,当我们需要更新resource的一些属性的时候，如果修改副本数量，增加、修改label，更改image版本，修改端口等。都可以直接修改原yaml文件，然后执行replace命令。 

## kubectl patch
`kubectl patch`用于修改正在运行的容器的属性  

## kubectl edit
`kubectl edit`用于灵活的在一个common的resource基础上，发展出更新过的significant resource。

## kubectl delete
`kubectl delete`根据resource名或label删除resource

## kubectl apply
`kubectl reply	直接在原有resource的基础上进行更新,同时还会resource中添加一条注释，标记当前的apply  

## kubectl logs
`kubectl logs <pod-name>`命令用于显示pod运行中，容器内程序输出到标准输出的内容  
`kubectl logs -f <pod-name>`可*跟踪查看容器日志*  
![kubectl logs](/kubectl_picture/4.PNG)

## kubectl exec
`kubectl exec`类似于docker的exec命令，为在一个已经运行的容器中执行一条shell命令,如果一个pod容器中，有多个容器，需要使用-c选项指定容器   
![kubectl exec](/kubectl_picture/5.PNG)  
可以看出已经进入该容器

## kubectl attach
`kubectl attach`命令类似于docker的attach命令，可以直接查看容器中以daemon形式运行的进程的输出，效果类似于logs -f，退出查看使用ctrl-c


## 问题诊断
### 问题：新建负载均衡一直成功不了
* 先用`kubectl get po -o wide`查看pod状态  
![kubectl get po -o wide](/kubectl_picture/6.PNG)   
可以看出`lb-myapp-2937719659-lg838`的STATUS为`ImagePullBackOff`  
* 再去EKOS看事件：  
![事件](/kubectl_picture/7.PNG)  
逐步查看可知，`lb-myapp-2937719659-lg838`被分配至node2，但是在启动其容器时失败了，因为pull其需要的image  `registry.ekos.local/ekos/ingress-controller`发生了错误，因次查看EKOS中的docker images  
![images](/kubectl_picture/8.PNG)  
没有发现该镜像  
* 去kubernets中看  
![docker images](/kubectl_picture/9.PNG)   
![docker images](/kubectl_picture/10.PNG)   
在node2上找到了该镜像，推测是ceph的原因导致该镜像不可用  
* 重启ceph：`systemctl restart ceph.target`  
* 重启ceph后依然没有，只好重启所有虚拟机,重启后成功






