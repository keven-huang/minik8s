# minik8s验收文档 G15    
 
minik8s是一个类似于kubernetes的迷你容器编排工具，能够在多机上对满足CRI接口的容器进行管理。支持Pod抽象管理容器生命周期，提供统一ClusterIP的Service，支持自动重启Pod的Replicaset，支持DNS服务，支持自动扩缩容，支持HPA，并基于交我算平台提供了对GPU应用的支持。最后实现了自选功能Serverless和控制面容错。

答辩演示视频（19min 三倍速）：
https://jbox.sjtu.edu.cn/l/519wZ8

各功能点视频合集：
https://jbox.sjtu.edu.cn/l/K1B746

小组成员：
- 韩金铂，组长，@[Chasingdreams6](https://github.com/Chasingdreams6)
- 陆浩旗，@[luhaoqi](https://github.com/luhaoqi)
- 黄嘉敏，@[keven-huang](https://github.com/keven-huang)

## 总体架构
![总体架构图](https://gitee.com/jinbohan/minik8s/raw/feature_doc/assets/arch.png)

minik8s的总体架构参考了kubernetes的架构，集群中存在唯一的master节点和若干worker节点。在master节点上运行api-server和etcd数据库，其它组件都通过api-server和数据库交互。

master的控制面包括api-server, scheduler, kube-proxy和kube-controller-manager；scheduler负责把pod调度到集群中的不同机器上，kube-proxy负责配置iptable与dns服务，controller-manager负责管理各个资源的生命周期，如replicaset, hpa, job, workflow等；kubelet负责监听api-server，创建对应的pod，同时通过心跳机制保证集群内节点的可用性；集群中的不同节点通过flannel进行通信；



## 软件栈
项目总体采用go语言开发，使用python语言完成serverless模块

使用flannel完成跨node通信，使用docker-go-client完成与docker的交互,
dockerClient调用了 https://github.com/fsouza/go-dockerclient

Pod的pause容器使用了google Pause3.6；etcd使用docker镜像部署。

使用iptable完成service，使用coredns+nginx完成dns；coredns和nginx都是docker镜像部署。

iptable的系统调用引用了一个开源的iptable for go，项目地址：https://github.com/coreos/go-iptables


![](https://notes.sjtu.edu.cn/uploads/upload_3ac674626951d4b23944c7d7e8127791.png)



## 分工和贡献度


| 成员   | 分工 | 贡献度 |
| ------ | -------- | -------- |
| 韩金铂 | dockerclient, 多机网络，service, dns, worker 心跳机制，scheduler随机策略，部分kubectl| 33.3%     |
| 陆浩旗 |kubectl，Replicaset，HPA，serverless中function相关，部分dockerclient中的工具，kubelet         |      33.3%    |
| 黄嘉敏|ci/cd,apiserver,list&watch机制,scheduler,hpa部分逻辑,gpu,workflow          |     33.3%     |



## 开发流程

gitee仓库地址：https://gitee.com/jinbohan/minik8s
ci/cd仓库地址：https://ipads.se.sjtu.edu.cn:2020/520021911392/minik8s/-/pipelines

### 分支介绍
- master 主分支
- develop 主开发分支
- feature_xxx 各个特性对应的分支
- apiserver, kubectl, kubelet 一些基本组件对应的分支

### git-flow

项目使用git-flow工作流开发。小组成员在开发新的feature时，从develop分支拉取对应的feature分支。在开发完毕后，提交P/R请求到develop分支。

下图是开发serverless的部分git工作流：
完整工作流：https://gitee.com/jinbohan/minik8s/graph/feature_serverless
![](https://notes.sjtu.edu.cn/uploads/upload_a60795f137b1ef04480564802e21101a.png)

下图是将feature分支合并到develop分支的一些P/R，P/R会指定代码的审查者与测试者。审查者查看历次commit记录，测试者运行`makefile`的一些测试脚本和`go test`，确保新的特性不会对原有的特性产生影响。

![](https://notes.sjtu.edu.cn/uploads/upload_3a84d974d5bf1adf0710bde87380a43b.png)

在项目开发基本结束时，提出将`develop`分支合并到`master`分支的P/R，发布了`minik8sv1.0`版本。

### 开发实践

在开发流程中，我们遵循极限编程（XP）的部分工作实践。

- 代码集体所有，任何人都可以自由修改代码，没有唯一对代码有所有权的人。
- 每日例会，在项目后期，每天通过腾讯会议or微信群交流进度，分配任务。
- 结对编程，通过腾讯会议的共享屏幕结对编程，轮流执行键入和质询，确保代码质量。
- 持续集成，我们使用了gitlab上的CI/CD。
- 测试先行，我们编写了一些自动化测试脚本和样例yaml。在项目开发的早期，我们还编写了一些针对Pod的example用于测试功能。


### 审查与测试

使用了代码走查的策略。当一位成员完成了部分工作，其它两位成员会对功能的实现方法和效果提出质询，当效果存在偏差时，会督促进行返工修改。

使用了测试脚本辅助测试。编写了Makefile，在其中对`service`, `dns`, `scheduling`等功能进行了测试，模拟输入kubectl命令的格式。

使用了`go test`来测试。对于`replicaset`, `hpa`, `service`等组件，编写了`test`文件测试，在`CI/CD`中会执行`go test`，确保最新的代码不会产生副作用。

### CI/CD

![](https://notes.sjtu.edu.cn/uploads/upload_71dfee844bf5eb15e69666357fe1c795.png)

在ci/cd配置的过程中，我发现如果按照教程中使用容器运行ci/cd很难做到环境的一致性，首先是容器运行时的测试依赖于flannel等插件的配置，并且还有go包下载过慢需要换源的问题出现。如果只是简单的单元测试固然是可以做到容器的运行，但我们希望一个真实的环境来模拟更为复杂的测试，比如创建replicaset命令等，因此我选择了shell作为gitlab-runner的执行器。

ci/cd分为两个步骤，第一个是build编译，通过编译来查看代码会有语法上的错误出现编译出错，第二个是test，我们编写了`Makefile`来辅助`go test`。`Makefile`脚本中包含`sudo make stop`,`sudo make clean`以及`sudo make run`。

`make stop` : 停止可能还没停止的kube-apiserver等进程

`make clean`: 清空可能出现的`etcd`多余数据残留

`make run`: 运行控制面的组件 
![](https://notes.sjtu.edu.cn/uploads/upload_f186d88f96998d6091d52212b37d56b5.png)

## 目录结构
- /assets 项目用到的一些图片

- /build 项目编译出的可执行程序


- /cmd 项目对外可以运行的.go程序
  - /client 测试与api-server交互的代码
  - /example 一些可运行的例子
  - /gpu 创建pod所用的go程序server
    - /server
        - server.go gpu上传下载逻辑server
        - ssh_test.go ssh的test文件
        - ssh.go ssh以及scp相关接口封装
    - main.go 启动server入口
  - /kube-apiserver api-server的运行入口
    - /app api-server的接口定义与handler
    - /test api-server的测试文件
    - /kube-apiserver.go 启动api-server
  - /kube-controller-manager
    - /kube-controller-manager.go 启动controller-manager
  - /kube-proxy 
    - /kube-proxy.go 启动kube-proxy
  - /kube-scheduler
    - /kube-scheduler.go 启动scheduler
  - /kube-service
    - /kubeservice.go 启动service-controller
  - /kubectl
    - /dns-example 演示dns的一些yaml文件
    - /gpu-yaml 演示gpu的一些yaml文件
    - /pod-example 演示pod的一些yaml
    - /sche-example 演示scheduler的一些yaml
    - /serverless 演示serverless的一些yaml
    - /serverless_math 演示serverless的数学功能
    - /service-example 演示service的一些yaml
    - /workflow-example 演示serverless-workflow的yaml
    - /kubectl.go 启动kubectl
    - *.yaml 其余的一些演示yaml文件，不一一赘述
  - /kubelet
    - /kubelet.go 启动kubelet
    - /dockerClient 一些与docker相关的示例
        
- /configs 一些全局配置
  - /dns 配置coredns相关模板
  - /gateway 配置nginx相关模板
  - /netconfig.go dns有关的一些全局常量

- /docs 文档
  - /assets 一些图片
  - *.pdf, *.md, *.pptx 项目需求验收指南，开题结题报告等

- /log 各个组件输出的log
- /pkg 源代码主文件夹
  - /api/core
     - /gputype.go gpu数据结构
     - /types.go 一些核心数据结构
     - /workflow.go workflow相关
  - /api/meta/v1  metaData使用的数据结构，从kubernetes中选择而来
    - /micro_time.go
    - /time.go
    - /types.go
  - /client
    - /informer informer实现
    - /tool
      - /httpclient.go 一些与api-server交互的函数
      - /listwatcher.go list&watch机制
  - /cmd
    - /apply kubectl的apply命令实现
    - /create kubectl的create命令实现
    - /delete kubectl的delete命令实现
    - /get kubectl的get命令实现
    - /invoke kubectl的invoke命令实现
  - /iptables go client for iptable, 是github上的开源组件
  - /kube-apiserver/etcd
    - /etcd.go 封装的etcd接口
    - /watcher.go watch实现
  - /kube-controller
    - /hpa-controller 
    - /job-controller
    - /replicaset-controller
    - /workflow-controller
  - /kube-proxy
    - /DnsManager.go DNS主体
    - /chain.go 构建svchain, podchain和DNAT rule规则
    - /kubeproxy.go kube-proxy主体
    - /types.go kube-proxy需要的数据结构
  - /kubelet
    - /config 一些kubelet全局配置
    - /dockerClient 封装的与docker engine相关的接口
    - /monitor 监控container status, 为了自动重启
    - /kubelet.go 启动kubelet
  - /runtime
    - /kube-service 
      - / manager.go Service-controller的主体逻辑
      - / service.go Service-controller的剩余逻辑
      - / singletons.go coredns和nginx的数据结构
      - / types.go RuntimeService的类型
    - /schema 仿照k8s的运行时
    - /interfaces.go 仿照k8s的运行时
    - /types.go  仿照k8s的运行时
  - /scheduler
    - /scheduler.go 调度器
    - /strategy.go 调度器使用的策略
  - /service
    - / service_test.go service测试
    - / types.go service需要的数据结构
  - /types 
    - /uid.go metadata需要的数据结构
  - /util 一些工具函数
- /scripts 脚本
  - /linux linux的sh脚本
  - /win windows的bat脚本
- /.gitlab-ci.yml ci/cd配置文件
- /Makefile make scripts

## 项目实现的功能

### 多机部署
我们使用了三台机器部署minik8s，其中m1作为master主节点，其余两台作为worker节点加入master。

我们使用 HTTP协议 进行网络通信来实现 worker 节点与 master 节点之间的消息传递，主要是发送给在master上的api-server， 使用 Flannel + etcd 来实现多台主机上分配给 Pod 的 ip 地址是全局唯一的。

启动master和worker我们采用编写Makefile脚本的方式，以下为具体命令：

``` shell
run:
	mkdir -p /root/nginx
	go run ./cmd/kube-apiserver/kube-apiserver.go > log/apiserver.log &
	go run ./cmd/kube-scheduler/kube-scheduler.go --strategy=RRStrategy > log/scheduler.log &
	go run ./cmd/kubelet/kubelet.go --nodename=node1 --nodeip=192.168.1.7 --masterip=http://192.168.1.7:8080 > log/kubelet.log &
	go run ./cmd/kube-controller-manager/kube-controller-manager.go > log/controller-manager.log &
	go run ./cmd/kube-proxy/kubeproxy.go --masterip=http://192.168.1.7:8080 > log/kubeproxy.log &
	go run ./cmd/kube-service/kubeservice.go > log/kubeservice.log &
m3:
	go run ./cmd/kubelet/kubelet.go --nodename=node3 --nodeip=192.168.1.11 --masterip=http://192.168.1.7:8080 > log/kubelet-m3.log &
	go run ./cmd/kube-proxy/kubeproxy.go --masterip=http://192.168.1.7:8080 > log/kubeproxy.log &
m2:
	go run ./cmd/kubelet/kubelet.go --nodename=node2 --nodeip=192.168.1.8 --masterip=http://192.168.1.7:8080 > log/kubelet-m2.log &
	go run ./cmd/kube-proxy/kubeproxy.go --masterip=http://192.168.1.7:8080 > log/kubeproxy.log &
```
重要的参数是--nodename指定节点名字，--masterip指定master节点ip以及--nodeip指定自己的ip

运行的时候在主节点运行`make run`, m2,m3作为worker节点运行`make m2`和`make m3`

使用命令 `kubectl get node -a`即可查看所有运行节点
![](https://notes.sjtu.edu.cn/uploads/upload_bafb0c8dc766ab942052375051626032.jpg)


### Pod抽象
Pod是minik8s的最小调度的单位，它包含多个container，同一个pod之间的container共享pause容器的网络空间，通过localhost进行通信。

不同Pod之间通过PodIP通信，使用flannel修改了各个节点docker0网桥对应的网段，m1的容器网段是10.0.13.1/24, m2的网段是10.0.12.1/24， m3是10.0.8.1/24。

> Pod的配置文件

以下是一个示例pod，它包含kind, name, image, containers, limit, ports， volumeMounts等字段。

```yaml=
apiVersion: v1
kind: Pod
metadata:
  name: test5
  labels:
    app: my-app
spec:
  nodeName: "node1"
  containers:
    - name: my-container
      image: nginx:latest
      command: ["nginx", "-g", "daemon off;"]
      limitResource:
          cpu: "1"
          memory: "128M"
      ports:
        - containerPort: 80
      volumeMounts:
        - name: my-volume
          mountPath: /usr/share/nginx/html
  volumes:
    - name: my-volume
      hostPath: /home/nginx
```

> 多容器Pod的创建，与Pod内部容器的通信

创建一个有两个容器的pod，配置文件如下：
```yaml=
apiVersion: v1
kind: Pod
metadata:
  name: two-con
  labels:
    app: two-con
spec:
  containers:
    - name: server
      image: chasingdreams/tinyweb:3
      ports:
        - name: http
          containerPort: 80
    - name: user
      image: chasingdreams/minor_ubuntu:v3
      tty: true
      command: [ "/bin/sh"]
```
演示Pod内部的user容器访问server容器的nginx服务。
![](https://notes.sjtu.edu.cn/uploads/upload_428d4fa99ff339e62b84c1843a950825.png)
创建的pod名为two-con，它有两个容器。
![](https://notes.sjtu.edu.cn/uploads/upload_74c292322a6983f02ec5f8d1803b2d4a.png)
在user容器中`curl localhost:80`，可以看到成功访问了server的服务。

> Pod的多机调度（支持RR和Random策略）

在3台机器的集群中，创建6个Pod，使用RRStrategy调度；
![](https://notes.sjtu.edu.cn/uploads/upload_f7c5c5decd848ea77c2e65b06735d44b.png)
可以看到，6个Pod按照node2, node3, node1的顺序轮询调度到了3个节点。

调度策略可以通过修改`make run`中scheduler的strategy来指定，目前支持`RRStrategy`和`RandomStrategy`。

调度器实现了调度算法和控制器解耦，添加调度策略只需要在`strategy.go`中添加新的逻辑即可。

> 使用volume完成同一Pod内多个容器、与宿主机的文件共享

示例Pod文件:
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: mnt
  labels:
    app: mnt
spec:
  containers:
    - name: user2
      image: chasingdreams/minor_ubuntu:v3
      tty: true
      command: [ "/bin/sh"]
      volumeMounts:
        - name: volume-007
          mountPath: /tmp/inner_path1
    - name: user1
      image: chasingdreams/minor_ubuntu:v3
      tty: true
      command: [ "/bin/sh"]
      volumeMounts:
        - name: volume-007
          mountPath: /tmp/inner_path2
  volumes:
    - name: volume-007
      hostPath: /tmp/host_path
```
上述文件创建了一个`volume-007`，绑定到了宿主机的`/tmp/host_path`上。
Pod内有两个user容器，user1把`volume-007`绑定到了`/tmp/inner_path1`上，user2把`volume-007`绑定到了`/tmp/inner_path2`上.

![](https://notes.sjtu.edu.cn/uploads/upload_4bd8a4a608d97fd8b03262d82337ced5.png)

上述演示，首先在宿主机的`/tmp/host_path`中创建了`1.txt`，进入user1的`/tmp/inner_path1`发现`1.txt`存在。这可以说明宿主机创建的文件可以通过host_mount传递给Pod内部的容器。

其次在user1的`/tmp/inner_path1`创建`2.txt`，再进入user2，发现在`/tmp/inner_path2`中发现`1.txt`和`2.txt`。这可以说明共享同一个volume的文件是共享的，同一个Pod的不同container之间可以共享文件。

最后在host的`/tmp/host_path`下发现了`1.txt`和`2.txt`，这可以说明Pod也可以通过host_mount，向宿主机传递文件。

> 对Pod的容器进行资源限制

示例Pod文件：
```yaml=
apiVersion: v1
kind: Pod
metadata:
  name: res
  labels:
    app: res
spec:
  containers:
    - name: server
      image: chasingdreams/tinyweb:3
      ports:
        - name: http
          containerPort: 80
      limitResource:
        memory: 50M
        cpu: 1
```
上述yaml创建了一个web，限制资源为1个CPU，50M内存。

创建这个pod，并在node2上用docker stats观察：
![](https://notes.sjtu.edu.cn/uploads/upload_eaea9000fe0e0e08bb39f47376cfb46b.png)

发现res-server的MEM_LIM为50mb，说明成功限制了资源占用。
![](https://notes.sjtu.edu.cn/uploads/upload_005b80cc3df868e555cdd6010ec679e4.png)



### Service抽象
我们实现了ClusterIP类型的Service，即Service对外提供一个统一的ClusterIP，用户可以通过ClusterIP访问到对应的服务。Service和Pod之间通过labels进行匹配。用户对Pod透明，Service通过一定的负载均衡策略把对ClusterIp的请求分配到对应的Pod。

>  演⽰利⽤配置⽂件创建Service的配置⽂件及运⾏情况

以下是我们service的yaml文件示例：
```yaml=
apiVersion: v1
kind: Service
metadata:
  name: webservice
spec:
  name: webservice
  clusterIP: "11.1.1.1"
  selector:
    app: "tinyserver"
  ports:
    - port: "80"
      protocol: tcp
      targetPort: "80"
```
service通过seletors匹配对应的pod，上述service就匹配label["app"]="tinyserver"的所有pod；

匹配的pod示例：
```yaml=
apiVersion: v1
kind: Pod
metadata:
  name: tinyserver1
  labels:
    app: tinyserver
spec:
  containers:
    - name: server
      image: chasingdreams/tinyweb:3
      ports:
        - name: http
          containerPort: 80
```

> 演⽰Service的单机可访问性
> 演示Service对多个Pod的映射

![](https://notes.sjtu.edu.cn/uploads/upload_319edc8546e625c48c90078ee8fcfbd9.png)

如上图所示，创建了一个名为webservice的服务，clusterIP为11.1.1.1, 对应两个服务pod，ip分别为10.0.12.2， 10.0.8.3，它们分别位于node2和node3上。		

`chasingdreams/tinyweb:3`镜像是一个特意修改的镜像，它可以返回自己的podIp，用于演示负载均衡。


在m1上使用`curl 11.1.1.1`连续两次，可以发现返回的结果分别是10.0.12.2和10.0.8.3，说明service通过负载均衡的策略对应到了不同pod。

> service可以动态监控pod的状态，更新当前的iptable

![](https://notes.sjtu.edu.cn/uploads/upload_d6c0a51be142606916113b8d75ccaa01.png)

如上图所示，增加了pod3，发现它被调度到了node1，ip为10.0.13.3，使用
`get service`命令，发现service中的pod变成了三个，说明service可以动态依据label添加pod。

在m1上使用`curl 11.1.1.1`访问，可以发现新增的pod也依照了负载均衡策略。

> 演⽰Service的多机可访问性

![](https://notes.sjtu.edu.cn/uploads/upload_73f117408dd6b422cc11c745a28cc535.png)

启动一个pod，在pod内部`curl 11.1.1.1`访问，发现也可以按照负载均衡策略正常访问到service。

#### 实现思路

service的功能需要service-controller和kube-proxy联合完成。service-controller中的informer负责管理runtimeService这个数据结构。这个数据结构是运行时的数据结构，不存储在etcd中。它维护了PodNameAndIP这个数组，保存符合条件的所有pods。

service需要实现动态筛选pod。这是通过service-controller中的一个定时线程ticker，它会每10秒钟向主线程发一次信号，主线程接受到这个信号之后，就会调用findPods函数，这个函数中会执行筛选pods和更新etcd。findPods中有一个小优化，就是调用samePods函数比较两次筛选的pods是否相同，如果相同就不更新etcd。

kube-proxy部分是监听etcd，当service-controller完成podNameAndIP的分配之后，更新iptable。iptable的转发逻辑是，在初始化时有一个名为SERVICE的链，它被添加到PREROUTING和OUTPUT链中的第一条，这样所有进入or发出该节点的流量都会被转到SERVICE链。SERVICE链中根据不同的serviceName, 将对应的service分配到svcchain。svcchain依照pod的个数，将流量依照负载均衡的策略导入到podchain。 每个podchain中写入DNAT规则，将导向clusterIP的流量改成导向这个pod的流量。注意要在每个链后添加RETURN规则，这样可以在当无法匹配时，不会影响原有的网络功能。

多机的flannel配置参考了这篇博客：
https://www.cnblogs.com/breezey/p/9419612.html#33%E5%AE%89%E8%A3%85flannel



### Replicaset抽象

> 演示利用配置⽂件创建ReplicaSet以及Pod失败后重启

以下是我们的replicaset-restart.yaml文件
``` yaml=
apiVersion: apps/v1
kind: ReplicaSet
metadata:
  name: my-replicaset
spec:
  replicas: 3
  selector:
    matchLabels:
      app: my-app
  template:
    metadata:
      labels:
        app: my-app
    spec:
      containers:
        - name: my-container
          image: chasingdreams/tinyweb:30sfail
#          command: ["timeout", "30s", "/bin/sh", "-c", "while true; do :; done"]
          resources:
            limits:
              cpu: "1"
              memory: "128Mi"
          ports:
            - containerPort: 80
          volumeMounts:
            - name: my-volume
              mountPath: /usr/share/nginx/html
      volumes:
        - name: my-volume
          configMap:
            name: my-configmap
```
重要的参数是Replicas指定了Replicaset的副本数量，selector指定了匹配的labels，支持后续动态手动创建的Pod加入Replicaset，只要Label符合即可，以及template中的Pod模板

这里为了演示自动重启我们使用的`tinyweb:30sfail`镜像是一个自制镜像，30s内会作为一个正常的Nginx镜像，访问80端口即会返回自己的IP地址，用来演示Replicaset与Service协同下的负载均衡策略，同时在30s后会自动Exit一个非零值，表示进入Failed状态，此时我们的kubelet会自动监听本地的Container状态，发现Failed或Succeed之后回去更新etcd中对应Pod的状态

Replicaset发现其中的Pod Fail之后会自动清理Pod并重启一个新的Pod，如果是Succeed则默认不删除，当做还在Running处理（这里的处理方式由自己决定，都可以）

创建命令为：

`kubectl create -f replicaset.yaml`

可以配合Get命令查看启动后的状态：

`kubectl get replicaset -a`和`kubectl get pod -a`

![](https://notes.sjtu.edu.cn/uploads/upload_dcbad3ec1a4aef9df67fb55f0d9ee7a5.jpg)

> 演示将Deployment绑定至Service的配置文件和运行情况

其实和Service的绑定并没有特别之处，因为Service支持Label匹配，只需要我们Replicaset中的Pod的Label和Service需要的一致就可以自动匹配使用ClusterIP了，以下是一个对应的Serviceyaml，多余不再赘述，可以再Service部分查看Service的实现方式

``` yaml=
apiVersion: v1
kind: Service
metadata:
  name: replicaset-service
spec:
  name: replicaset-service
  clusterIP: "11.1.1.2"
  selector:
    app: "my-app"
  ports:
    - port: "80"
      protocol: tcp
      targetPort: "80"
```
下图可以看到，可以通过serviceIP去访问这个repliaset里的pod，并且随着RC中pod的更换，service也可以动态更新。

![](https://notes.sjtu.edu.cn/uploads/upload_16aab2989ab821a54d134a30fc5fcd5a.png)



#### 实现方式
Replicaset实现主要需要写一个Replicaset Controller，重点需要2个Informer来监听Replicaset和Pod以及一个worker queue来处理副本数量不匹配的情况

当创建一个Replicaset的时候我们需要遍历Pod来看哪些已有的Pod属于这个Replicaset，同时修改对应的Status中的Replicas属性，如果和Spec中预期的Replicas数量不一致则加入Worker queue

Worker线程会轮询queue来看是否有需要处理的Replicaset，从queue中取出之后观察是需要增加副本还是减少副本，然后再去创建或删除对应的Pod即可

其中的难点在于如何考虑对etcd中的状态的获取即由此引发的并发问题，因为我们没有使用ResourceVersion，因此需要考虑好修改的逻辑链条的前后顺序，当然如果有并发支持的设计写起来应该会简单很多

另外我们目前的处理是在Informer监听或者自己修改对应的Replicaset属性后判断是否加入queue，还可以由一种实现方式是worker隔一段时间轮询所有的Replicaset检查是否需要处理，好处在于不用担心在哪里忘记加入队列


### Auto-scaling
hpa的yaml文件如下：

```yaml=
apiVersion: app/v1
kind: HorizontalPodAutoscaler
metadata:
  name: my-hpa
spec:
  scaleTargetRef:
    kind: ReplicaSet    	#管理的资源名
    name: my-replicaset
  minReplicas: 2	  # 最小副本数
  maxReplicas: 10	  # 最大副本数
  periodSeconds: 10   # 扩缩容速度
  metrics: 
  -  name: cpu 				# 对应指标
     target:
      type: Utilization
      value: 30
```

我们有两种指标支持，一个是cpu利用率，另一个是内存利用率，内存利用率的写法为

```yaml=
  metrics: 
  -  name: memory 				# 对应指标
     target:
      type: Utilization
      value: 30
```
#### 使用方法

创建完对应的replicaset后我们就可以创建对应的hpa来管理这个replicaset

```shell
$  kubectl create -f hpa.yaml
```

![](https://notes.sjtu.edu.cn/uploads/upload_526f2cd7ad9c4a7383c41af64c7f737d.png)

![](https://notes.sjtu.edu.cn/uploads/upload_aeff6f582af4c1005f5852b3f480dfd7.png)

上图演示了replicaset的扩缩容过程。

#### 实现方式

对于节点的容器资源监控，我们尝试了两种方式，一种直接用docker的stats监控得到docker cpu占用资源，另一种用Cadvisor和Prometheus监控docker资源，通过Prometheus Client编写query查询对应的资源占用。

在尝试了Cadvisor和Prometheus监控docker资源后发现不太准，延迟较高，我们写了个死循环的container从cadvisor数据来看大概滞后了1分钟左右，所以我们放弃了使用cadvisor转而使用docker stats的方法监控容器资源。

由于已经有kubelet组件，我们决定让他同时作为一个server能够被动的接受http请求从而获取节点容器的内存使用率和cpu使用率，这样获取容器资源利用率的问题就解决了。

有了获取资源的方法，那么我们就只需要一个`hpacontroller`来管理对应的pod副本就可以了，在获取创建hpa的信息后，`hpacontroller`获取对应的pod资源信息，要想支持auto-scaling我们只需要复用ReplicaSet提供的管理副本的功能，简单来说，HPA根据指标计算出ReplicaSet的新的replica数量，修改Replicaset对应的Replicas属性，从而让Replicaset管理对应的replica。

对于扩容/缩容速度的问题，我们使用了比较简单的实现方式，即使用periodSeconds属性来表示扩容/缩容一个Pod之后停顿的时间，来间接表示扩缩容的速度，这样我们实现的时候只需要在扩/缩容创建/删除一个Pod之后进行一个简单的Sleep就可以实现对扩缩容速度的控制。

### DNS

DNS提供了hostname+path能访问对应service的功能；

> 演⽰利⽤配置⽂件对DNS与转发规则进⾏配置的配置⽂件与运⾏状况

DNS配置文件：
```yaml=
apiVersion: v1
kind: DNS
metadata:
  name: testDns
spec:
  host: hanjinbo.com
  gatewayIp: 13.1.1.1
  paths:
    - name: /path1
      service: service1
      port: "80"
    - name: /path2
      service: service2
      port: "80"
```
host表示这个dns的域名，gatewayIP是这个dns对应的nginx服务的svcip（这个nginx服务被称为gatewayService，作用是将/path1, /path2导向service1, service2的clusterIP）。 paths是各个path规则。

> 集群中的宿主机上使用DNS

![](https://notes.sjtu.edu.cn/uploads/upload_ee039742e4dcdeabab2229db34b379cd.png)

在m1上执行`make testdns`命令，它会创建一些pod和service。上图所示，创建了两个service，ip分别是11.1.1.2和11.1.1.3，它们各自下面有一个pod；
再创建一个dns；最后效果是在宿主机上使用`curl hanjinbo.com/path1`能访问到service1，使用`curl hanjinbo.com/path2`能访问到service2。

> 集群中的Pod使用DNS

![](https://notes.sjtu.edu.cn/uploads/upload_b74c57a85c2df4fef79c5775b16d2cf1.png)

在集群中创建一个test5Pod，在pod内部同样curl两次，依然可以返回正确结果。

#### 实现思路

使用coredns+nginx的方式完成DNS。coredns是一个固定的dns服务，我使用了一个在master节点上的pod，在其中跑一个coredns；再创建了一个coredns服务，服务的IP固定为12.1.1.1. 那么12.1.1.1就是整个集群的dns服务的ip。这个coredns是在控制面启动的时候创建的。

当创建一个dns资源时，需要在coredns的hostfile中添加一项。按照上面的例子，就添加一个`13.1.1.1 hanjinbo.com`的映射。这在coredns服务器中增加了一个entry，coredns服务器就会把对于`hanjinbo.com`的dns请求指向13.1.1.1。13.1.1.1是一个nginx-service。当这个dns启动时，会在master节点下的某个路径创建一个名为DNSName的nginx文件夹，将这个nginx文件夹映射入一个nginx的pod里。master再创建一个nginx.conf，在其中写path1, path2映射到svc1, svc2的规则。当nginx-service启动完成之后，重启nginx-pod和coredns，应用当前更改。

DNS同样需要service-controller和kube-proxy的配合。dns服务有5个状态，分别是初始状态(s1)、文件夹创立状态(s2)、Nginx服务已创立状态(s3)，运行状态(s4)，被删除状态(s5)。

当kube-proxy监听到dns资源的创建之后，它会创建nginx文件夹，之后这个dns就变成文件夹创立状态(s2)。

当service-contoller监听到有处于s2状态的dns时，他会找到singletons.go里的gateway模板，根据当前的dns名字创建正确的nginxPod和nginxService，将dns的path对应的serviceIP更新，并将dns的状态变更为s3。

当kube-proxy监听到s3状态的dns时。会根据dns中的serviceIP信息，重写对应文件夹下的nginx.conf文件，并重启nginx和coredns,dns状态变更为s4，即可用状态。

当kube-proxy, service-controller监听到deletedns请求时，分别会删除nginx文件夹, coredns项目和nginx服务。dns状态变更为s5, 随即在etcd中被删除，生命周期结束。


### 容错

####  List & Watch的容错
 由于我们minik8s的list&watch实现采用了http分块编码传输的实现方式，因此apiserver进程崩溃意味着http连接断开，而controller和kubelet中的list&watch都依赖于http连接，因此我们在两边增加了http连接断开后的定时重试机制，能在apiserver重启或是组件重启后重新连接。

#### 动态增删Node
我们在worker节点的kubelet上运行了心跳发送器，它会每5秒向master发送心跳。master如果10秒没有接受到对应worker的心跳，说明worker节点已经fail，master会将其从数据库中删除，并且调度器也不会调度到一个fail的worker。

首先，启动master和两个worker，使用get node可以发现有三个节点。

在m3节点上使用`make stop`停止m3，再在master上get node，发现只有两个节点。说明master可以动态删node。

再将m3节点加回集群。

![](https://notes.sjtu.edu.cn/uploads/upload_f3c4573e02cfa975b0917ea39ecebfc7.png)


#### 控制面容错
> minik8s的控制面重启之后，pod, service均能正常访问

首先创建一个pod和service。
![](https://notes.sjtu.edu.cn/uploads/upload_7aacb2fea97e30fd487fa7263e69bc3b.png)

在master节点上运行`make stop`，停止所有控制面。使用`docker ps -a`发现pod依然存在。
![](https://notes.sjtu.edu.cn/uploads/upload_b824bef00537c8e9a12c454825ae47ff.png)



再重启控制面，使用`get pod`和`get service`命令，可以发现刚才创建的pod和service还能查询到。并且通过`curl 11.1.1.1`命令依然能够访问service。说明控制面的重启并不会影响已有pod和service的正常运行。
![](https://notes.sjtu.edu.cn/uploads/upload_ab23c771e45d46f61b2090d511da4a92.png)


这些容错并不需要特殊的实现，因为数据本来就存在etcd，控制面重启不会影响，只需要对应的Controller有相应的处理有初始值的逻辑即可


### GPU

首先是我们`gpu`的yaml文件介绍，yaml文件如下：

```yaml=
apiVersion: app/v1
kind: Job
metadata:
  name: cuda
spec:
  task:
    name: cuda
    partition: dgx2 # 分区（账号限制只能为dgx2）
    nodes: 1 # 计算节点数
    ntasks-per-node: 1 # 每个节点上任务数
    cpus-per-task: 1 # 每个任务使用cpu核心数
    gpu: 1 # gpu卡数（由于账号限制只能为1）
    output: cuda.out # 结果输出文件
    error: cuda.err # 错误输出文件
    mail-type: end # 邮件通知（任务结束后发送邮件通知）
    mail-user: jmhuang_2020@sjtu.edu.cn # 邮箱地址
    program: ./cmd/kubectl/gpu-yaml/cuda.cu # 运行的cuda程序
```
#### 使用方法
使用下述命令即可执行job文件

```shell
$ kubectl create -f job-example.yaml
```

通过下述命令即可获取job任务的当前状态为pending或succeeded或failed

```shell
$ kubectl get job <jobname> or kubectl get job -a
```

任务结果存放于`/home/job/<jobname>`的文件夹下，可以查看对应的`.out`以及`.err`文件

#### 实现方式

实现gpu的主要方式为创建一个go server，对于每一个任务都创建一个蕴含server的pod中去执行上传cuda程序获取结果的过程。

在实现的过程中有几个逻辑是需要考虑的，一个是多机过程中文件上传和结果获取的形式，二是如何上传文件到交大集群并获取结果。

由于创建pod节点的是从节点，而上传文件的是主节点，这时候就会出现上传文件以及获得文件的问题。如何将主节点上传的文件如何打包到从节点以及如何将从节点创建pod后得到的job结果传回主节点是完成job功能的两个关键点。

文件的上传有两个方法，一个是让kubelet去监听主节点存放file的地方，一个是在上传的时候将文件打包上传给apiserver，当kubelet创建对应pod的时候让kubelet下载对应的cuda程序和生成的slurm脚本。我们minik8s的实现采用了第二种方式，在通过`kubectl create -f`上传cuda程序和对应yaml文件的时候，apiserver将其转化为slurm脚本，将其保存在etcd中。kubelet创建pod的时候设定一个标志`gpu`，当kubelet创建pod时观察到`gpu`标志后从apiserver处获取对应的job文件保存到文件夹中。

而当pod运行完毕后如何将结果的out文件和err文件放回到主节点，让主节点能够访问是十分关键的，我采用了共享文件系统的形式，master节点将`/home/job`文件夹共享给两个从节点，而从节点的server将自己目录的`/home/job`文件夹挂载到主节点处，实现文件系统的共享，从而实现文件实现结果的共享。

而上传文件到交大集群并获取结果实现方式采用写一个go语言的server，并在pod执行这个server来解决。创建的pod挂载了`/home/job`文件夹，因此我们可以拿到对应的job文件信息，再此基础上我们可以写一个基于go语言的server，启动后将文件上传到交我算运算平台，并且通过`ssh`远程连接的形式，输入`sbatch <jobname>.slurm`上传工作脚本文件等待结果，并且上传后能获取`jobid`来为将结果文件下载做准备。而获取结果采用轮询的方式，交我算提供了`saact -j <jobid>`的命令来获取程序的运行状态，通过输入命令`"sacct -j " + s.jobID + "| tail -n +3 | awk '{print $6}'"`来获取上传的gpu任务的状态，每5s轮询一次，如果发现结果是`COMPLETED`那么就将远程的文件下载至本地创建的任务文件夹内，实现结果的保存和共享。

`kubectl get job`的功能则是基于`pod`运行状态来管理的，由于当pod运行的是一个go语言写的server程序，当获取完结果并成功下载过后对应的server程序将退出，容器状态为exited，当发现pod状态为exited时，如果是异常退出说明任务提交出异常错误，job状态为failed，如果正常退出则job状态为succeed

此外我们也实现了job的delete功能，能够删除对应的pod和job信息。
### Serverless 中的 Function部分

> 演示Function的定义和运行程

还是使用`kubectl create -f function.yaml -d buildcontext`命令来进行创建

注意到这里我们不止有-f来制定对应的yaml文件，还有一个-d来指定一个目录，这是因为我们在实现过程中一个function的创建就会创建一个对应的image，因此需要一个build context来作为创建镜像的上下文，里面需要有Dockerfile，web server的python代码，需要额外import的库的requirements.txt文件以及用户自定义的my_module.py函数文件，我们的Dockerfile会把这些文件都放到image中，启动的时候自行运行serverless_server.py文件

每次在master的kubectl创建函数都会build image & push image，push到DockerHub是为了在多台机器之间进行镜像的同步，因为也会受限于和DockerHub之间的网络通信情况，不过在我们的测试中多数情况下还是可以的，create一个全新的function用时不会超过一分钟，多数在半分钟以内

为了防止worker使用之前有的image，即使是同一个函数的不同版本的image name也需要不同

以下是对应的演示：
创建function：
![](https://notes.sjtu.edu.cn/uploads/upload_9dd1a34f6ba5bf83615c44cda22f6043.jpg)

调用函数：
`kubectl invoke -f <function_name> -p <param>`
![](https://notes.sjtu.edu.cn/uploads/upload_b5a49516bfe21199b8aadf6ba6593f0a.jpg)

#### 实现过程
Serverless部分的处理相对来说不是特别复杂，实现的过程中参考了@[WilliamX1学长](https://github.com/WilliamX1/minik8s#DNS-%E4%B8%8E%E8%BD%AC%E5%8F%91)的实现

简单来说我们需要用python来写一个web server来处理HTTP请求，serverless_server.py文件，核心逻辑是：
```python=
@app.route('/function/<string:module_name>/<string:function_name>', methods=['POST'])
def execute_function(module_name: str, function_name: str):

    module = importlib.import_module(module_name)
    event = {"method": "http"}

    if request.headers.get('Content-Type') == 'application/json':
        context = request.get_json()
    else:
        context = request.form.to_dict()

    try:
        result = getattr(module, function_name)(event, context)
        return result, 200
    except Exception as e:
        logger.info("An error occurred during function execution:", e)
        return "Error during function execution", 500
```
就是调用放在同目录下的my_module.py文件中的对应函数即可，这个函数的名称可以有用户自定义，但需要和function.yaml中的name相对应，同时传入的参数必须是指定格式，下面是一个例子：
``` python=
# my_module.py
def calc(event: dict, context: dict)->dict:
    operator = context.get('operator')
    x = context.get('x')
    y = context.get('y')

    if operator and x is not None and y is not None:
        expression = f"{x} {operator} {y}"
        try:
            result = eval(expression)
            return {'result': result}
        except Exception as e:
            return {'error': str(e)}

    return {'error': 'Invalid input'}

```

**如何实现Scale-up和scale to zero?**

说来也是比较简单：我们在上面的python写的web server中加入两个计时器，一个用来计时离上次调用多少时间，时间达到我们设定的阈值（例如30s）之后就会通过`os._exit`来退出整个程序，这样对应的Container就会Failed，我们的kubelet监听到之后就会自动回收，由此实现scale to zero

另外一个计时器用来计算一段时间内的RPS，也就是请求数，当超过我们的阈值（为了方便演示我们定为20s内invoke超过四次就进行scale up）就会进行扩容，主要的逻辑就是向apiserver对应的地址发送一条HTTP请求，里面包含自己的function name即可，然后由api-server来完成对应的扩容请求

**如何实现热启动？**
其实就是在invoke的时候先检查一遍Pod是否有我们需要的Pod，如果有的话就列出来按照一定的策略，例如Round-robin来进行选择，如果没有的话则需要创建一个新的Pod也就是冷启动，我们冷启动一次因为需要等待Container的创建，包括网络IP的分配大概需要10s左右的时间，热启动相当于就是api-server找到对应的PodIP之后把HTTP请求进行一个转发，相对简单

### Workflow
workflow的yaml文件如下:

```yaml=
apiVersion: apps/v1
kind: Workflow
metadata:
  name: my-workflow
spec:
  InputData : ./cmd/kubectl/workflow-example/input.yaml
  states:
    - Name: StartState 
      Type: Input
      Next: branchstate
    - Name: branchstate
      Type: Choice
      Choices:
        - Condition:
            Type: Numeric
            Variable: x
            Operator: ">"
            Value: 10
          Next: Function1
        - Condition:
            Type: Numeric
            Variable: x
            Operator: "<="
            Value: 10
          Next: Function2
    - Name: Function1
      Type: Task
      Resource : function1
      Next: Function3
    - Name: Function2
      Type: Task
      Resource: function2
      Next: Function3
    - Name: Function3
      Type: Task
      Resource: function3
      End : true
```

yaml文件模仿的是AWS Step Function的写法。state的类型分为Input、choice和task。在workflow的开头需要指明input文件的path地址，每个state的next标签指代下一个要跳转的状态，如果是结束状态则需要标志符End。调用函数需要state的类型为Task并且定义Resource为创建的对应函数。需要创建分支则需要Choice类型的State并且定义比较的变量，运算符和值来跳转到对应的状态。

#### 使用方法

首先需要创建workflow涉及到的resource function。如例子中的yaml就需要创建函数`function1`、`function2`和`function3`

创建完对应的function后输入下述命令创建并执行工作流

```shell
$  kubectl create -f workflow.yaml
```

同时也可以通过`kubectl get workflow`来获取工作流的结果。如果显示`Pending`则意味着函数仍在计算当中，否则则会显示工作流的返回结果。

![](https://notes.sjtu.edu.cn/uploads/upload_680a912e9cc48bafc8f2faa62141114b.png)
上图显示了通过`get workflow` 获取工作流的结果。

#### 实现方式

而workflow的实现方式即为通过apiserver将workflow的数据结构转化为dag图的结构，图的结构包含node和edge。

同时创建`workflow_controller`，list&watch对应workflow资源，当监听到新的workflow任务便开启额外的`goroutine`来完成一个工作流。dag图中存储了开始节点startnode以及已经通过json解析好的input文件，因此从图开始遍历调用相关节点即可。如果是分支节点就判断跳转到下一个节点，如果是task节点就调用apiserver的接口来调用对应函数。遇到end节点就返回并通知apiserver工作流完成，更新工作流状态。





