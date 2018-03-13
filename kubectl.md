# kubernetes����
## kubectl get
* `kubectl get po`��ȡ��ǰ���е�����pods����Ϣ  
![kubectl get po](/kubectl_picture/1.PNG)  
���Կ�����Ⱥ��һ������5��pod����һ���������ʲô���ڶ���������mysql��myapp�ĸ��ؾ��⣬������������ǵķ������Ҫ��ȡpod�������ĸ��ڵ��ϵ���Ϣ������`kubectl get po -o wide`  
![kubectl get po -o wide](/kubectl_picture/2.PNG)  
* `kubuctl get namespace`��ȡnamespace��Ϣ
* `kubectl get rc`, `kubectl get svc`, `kubectl get nodes`�Ȼ�ȡ����resource��Ϣ  
* �����������ʹ�ã�`kubectl get po --namespace=kube-system`  
![kubectl get ](/kubectl_picture/3.PNG)  

## kubectl describe
`kubectl describe`������get��ͬ�����ڻ�ȡresource�������Ϣ����ͬ���ǣ�get��õ��Ǹ���ϸ��resource���Ե���ϸ��Ϣ��**describe��õ���resource��Ⱥ��ص���Ϣ**��`kubectl describe po mysql-server-0`

## kubectl create
* `kubectl create`���ڸ����ļ������봴����Ⱥresource,����Ѿ���������Ӧresource��yaml��son�ļ���ֱ��`kubectl create -f filename`���ɴ����ļ��ڶ����resource  
* ֱ��ʹ��`kubectl create`����Ի���yaml�ļ�������ReplicationController��rc����rc�ᴴ����������  

## kubectl replace
`kubectl create`�������ڶ�������Դ���и��¡��滻,��������Ҫ����resource��һЩ���Ե�ʱ������޸ĸ������������ӡ��޸�label������image�汾���޸Ķ˿ڵȡ�������ֱ���޸�ԭyaml�ļ���Ȼ��ִ��replace��� 

## kubectl patch
`kubectl patch`�����޸��������е�����������  

## kubectl edit
`kubectl edit`����������һ��common��resource�����ϣ���չ�����¹���significant resource��

## kubectl delete
`kubectl delete`����resource����labelɾ��resource

## kubectl apply
`kubectl reply	ֱ����ԭ��resource�Ļ����Ͻ��и���,ͬʱ����resource�����һ��ע�ͣ���ǵ�ǰ��apply  

## kubectl logs
`kubectl logs <pod-name>`����������ʾpod�����У������ڳ����������׼���������  
`kubectl logs -f <pod-name>`��*���ٲ鿴������־*  
![kubectl logs](/kubectl_picture/4.PNG)

## kubectl exec
`kubectl exec`������docker��exec���Ϊ��һ���Ѿ����е�������ִ��һ��shell����,���һ��pod�����У��ж����������Ҫʹ��-cѡ��ָ������   
![kubectl exec](/kubectl_picture/5.PNG)  
���Կ����Ѿ����������

## kubectl attach
`kubectl attach`����������docker��attach�������ֱ�Ӳ鿴��������daemon��ʽ���еĽ��̵������Ч��������logs -f���˳��鿴ʹ��ctrl-c


## �������
### ���⣺�½����ؾ���һֱ�ɹ�����
* ����`kubectl get po -o wide`�鿴pod״̬  
![kubectl get po -o wide](/kubectl_picture/6.PNG)   
���Կ���`lb-myapp-2937719659-lg838`��STATUSΪ`ImagePullBackOff`  
* ��ȥEKOS���¼���  
![�¼�](/kubectl_picture/7.PNG)  
�𲽲鿴��֪��`lb-myapp-2937719659-lg838`��������node2������������������ʱʧ���ˣ���Ϊpull����Ҫ��image  `registry.ekos.local/ekos/ingress-controller`�����˴�����β鿴EKOS�е�docker images  
![images](/kubectl_picture/8.PNG)  
û�з��ָþ���  
* ȥkubernets�п�  
![docker images](/kubectl_picture/9.PNG)   
![docker images](/kubectl_picture/10.PNG)   
��node2���ҵ��˸þ����Ʋ���ceph��ԭ���¸þ��񲻿���  
* ����ceph��`systemctl restart ceph.target`  
* ����ceph����Ȼû�У�ֻ���������������,������ɹ�






