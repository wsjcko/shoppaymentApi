go install go-micro.dev/v4/cmd/micro@master

micro new service github.com/wsjcko/shopcart

mkdir -p domain/{model,repository,service} 
mkdir -p protobuf/{pb,pbserver} 
mkdir -p proto/{pb,pbserver}
mkdir common

go mod edit --module=github.com/wsjcko/shopcart
go mod edit --go=1.17  

gorm 有个根据创建表sql 生成model  : gormt

清除mod下载的包
go clean -modcache


### consul 微服务注册中心和配置中心
docker search --filter is-official=true --filter stars=3 consul
docker pull consul

## 生产环境要注意数据落盘  -v /data/consul:/data/consul
docker run -d -p 8500:8500 consul:latest 

### 注册中心
"github.com/asim/go-micro/plugins/registry/consul/v4"

### 配置中心
"github.com/asim/go-micro/plugins/config/source/consul/v4"

### consul数据库配置
http://127.0.0.1:8500/ui/dc1/kv/create

key: micro/config/mysql

{
  "host":"172.21.222.223",
  "user":"root",
  "pwd":"123456",
  "database":"shopdb",
  "port":"3306"
}

### 链路追踪jaeger 耶格 
[官方文档](https://www.jaegertracing.io/docs/1.32/)

docker run -d --name jaeger \
  -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
  -p 5775:5775/udp \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 14268:14268 \
  -p 9411:9411 \
  jaegertracing/all-in-one:latest

docker inspect -f '{{.HostConfig.PortBindings}}' jaeger

  http://127.0.0.1:16686/search


### 绑定链路追踪 服务端绑定handle 客户端绑定client
micro.WrapClient(opentracing4.NewClientWrapper(opentracing.GlobalTracer()))
micro.WrapClient(opentracing4.NewClientWrapper(opentracing.GlobalTracer()))

opentracing4 "github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v4"

### 创建链路追踪实例
"github.com/opentracing/opentracing-go"
"github.com/uber/jaeger-client-go"
"github.com/uber/jaeger-client-go/config"


### 集流量控制、熔断、容错，负载均衡等hystrix-go
docker search hystrix
docker pull mlabouardy/hystrix-dashboard
docker run --name hystrix-dashboard -d -p 9002:9002 mlabouardy/hystrix-dashboard:latest


### 购物车加入熔断（客户端），限流（服务端），负载均衡（客户端）

github.com/asim/go-micro/go-plugins/wrapper/ratelimiter/uber/v4 限流


docker run --rm -p 8080:8080 gharsallahmoez/micro --registry=consul --registry_address=172.21.222.223:8500 api --handler=api 
micro --registry=etcd --registry_address=172.21.222.223:8500 //最新版本支持etcd

docker run -d --network host -t -p 8088:8088 -p 9096:9096 --name shopcartapi shopcartapi:latest -ip "172.21.222.223"
docker run -d --network host -t -p 8087:8087 --name shopcart shopcart:latest -ip "172.21.222.223"

http://192.168.65.4:9096/hystrix.stream   ip:192.168.65.4  从注册中心 查看http://127.0.0.1:8500/ service api.shopCartApi 和service.shop.cart

http://127.0.0.1:8080/shopCartApi/findAll?user_id=3



### 支付paypal    refund_id唯一 幂等性
http://127.0.0.1:8080/shopPaymentApi/payPalRefund?payment_id=1&refund_id=111&money=110.10


curl -L -O https://artifacts.elastic.co/downloads/beats/filebeat/filebeat-8.1.2-linux-x86_64.tar.gz
tar xzvf filebeat-8.1.2-linux-x86_64.tar.gz

./filebeat -e -c filebeat.yml -d


## 安装docker-compose  

Docker-Desktop  Setting ->General ->Use Docker Compose V2  
docker-compose --version


或者

echo `uname -s`-`uname -m`
Linux-x86_64

sudo curl -L https://github.com/docker/compose/releases/download/2.4.1/docker-compose-`uname -s`-`uname -m` -o /usr/local/bin/docker-compose
chmod 777 /usr/local/bin/docker-compose
docker-compose --version














# Docker 使用

### 列出容器
docker ps -h 使用说明
docker ps 本地运行中
docker ps -a 所有的（运行，停止）
docker ps -s  容器文件大小，获得 2 个数值：一个是容器真实增加的大小，一个是整个容器的虚拟大小。容器虚拟大小 = 容器真实增加大小 + 容器镜像大小。
docker ps --no-trunc 即不会截断输出。该选项有点长，其中 trunc 算是 truncate 的缩写。
docker ps -l 最后被创建的容器 相当于docker ps -n 1
docker ps -n 2 最后被创建的2个容器
docker ps -q 只显示容器 ID ，清理容器时非常好用，filter 过滤显示一节有具体实例。

联合使用
docker ps --format "{{.ID}}: {{.Command}}"  --no-trunc  
docker ps --filter name=mysql57 --format "{{.ID}}: {{.Command}}"  --no-trunc
docker ps --filter name='.*mysql.*' 正则匹配

### 启动，关闭，重启
docker start consul或id
docker stop consul或id
docker restart consul或id


### 附着到容器上  run -d后台运行，创建守护式容器
docker run -d --name topdemo ubuntu /usr/bin/top -b
docker attach topdemo

attach 进入容器topdemo正在执行某个命令的终端，不能在里面操作

和 docker exec -it topdemo区别在于， exec 进入容器topdemo开启新的终端，可以操作


### 获取容器日志
docker logs consul

跟踪; 和tail -f 一样
docker logs -f consul

获取日志最后10行,会读取整个日志文件
docker logs --tail 10 consul

跟踪最新日志,不必读取整个日志文件
docker logs --tail 0 -f consul

加上时间戳
docker logs --tail 0 -ft consul


### docker日志驱动和存储驱动 ,默认json-file 为docker logs 提供基础
docker支持的日志驱动，指定 docker run --log-driver="syslog" -name consul 

none            无日志
json-file       将日志写入json-file，默认值  docker logs
syslog         将日志写入syslog，syslog必须在机器上启动
journald      将日志写入journald,journald必须在机器上启动
gelf             将日志写入GELF端点，如Graylog或Logstash
fluentd        将日志吸入fluentd，fluentd必须在机器上启动
awslogs      将日志写入亚马逊Cloudwatch
splunk         使用HTTP事件收集器将日志写入splunk
etwlogs       将日志消息作为windows时间跟踪。仅在windows平台可用
gcplogs       将日志写入Google云平台
nats             将日志发布到NATS服务器

### port
查看容器的端口映射
docker port consul
docker port mysql57
docker port shopcart

ShopProduct  : 8085
ShopCategory : 8086
ShopCart     : 8087
ShopCartApi  : 8088
ShopOrder    : 8089
ShopPayment  : 8090
ShopPaymentApi  : 8091


### 重命名
docker run --name consule 命名出错
docker rename consule consul
或
docker rename 78ec447d1388 consul

### 删除容器
docker rm consul 或 id

删除所有
docker rm `docker ps -a -q`

### 镜像
docker images
docker pull mysql:5.7

docker images mysql:5.7 查看本地镜像具体  size 448MB
docker ps -s --filter name=mysql57  :   增加大小 51.1MB (virtual 499MB = 448+51)


docker images -f xxxx
查看对应的过滤条件
这个过滤标签的格式是 “key=value”，如果有多个条件，则使用这种 --filter “key1=value” --filter “key2=value”

当前支持的过滤配置的key为

dangling：显示标记为空的镜像，值只有true和false
label：这个是根据标签进行过滤，其中lable的值，是docker在编译的时候配置的或者在Dockerfile中配置的
before：这个是根据时间来进行过滤，其中before的value表示某个镜像构建时间之前的镜像列表
since：跟before正好相反，表示的是在某个镜像构建之后构建的镜像
reference：这个是添加正则进行匹配


比如：过滤没有打标签的镜像
docker images -f "dangling=true" -q   全部没有打标签tag的ImageId
echo $(docker images -f "dangling=true" -q)
49d1b3444c46 a6dc65fb20ce 0f4e8139b3ba
 docker rmi  $(docker images -f "dangling=true" -q)

REPOSITORY   TAG       IMAGE ID       CREATED        SIZE
<none>       <none>    49d1b3444c46   17 hours ago   25.3MB
<none>       <none>    a6dc65fb20ce   18 hours ago   25.3MB
<none>       <none>    0f4e8139b3ba   19 hours ago   25.2MB

docker images --filter=reference='*:*'

REPOSITORY     TAG       IMAGE ID       CREATED        SIZE
shopcartapi    latest    7624e394f127   17 hours ago   25.3MB
shopcart       latest    201773e02aba   17 hours ago   25.2MB
shopcategory   latest    47e10a8ce53c   2 days ago     23.9MB
shopproduct    latest    f76d3de06822   2 days ago     25.1MB
micoserver     latest    97bc7bba963e   4 days ago     28.3MB
nginx          latest    605c77e624dd   3 months ago   141MB
redis          latest    7614ae9453d1   3 months ago   113MB
mysql          5.7       c20987f18b13   3 months ago   448MB
consul         latest    76802375bc5c   3 months ago   118MB

docker search jaeger
docker search --filter is-official=true --filter stars=3 jaeger

### 构建镜像
docker build -t shopcartApi:latest .  (同目录.   也可以指定-f Dockerfile)
docker commit Id  wsjcko/shopcartApi 现在不推荐
docker login  远程hub

如果一个构建失败，可以基于最后的成功构建新容器
docker build -t shopcartApi:latest .       RUN: step3 b867c0c39ccd  RUN: step 4 error
docker run -it b867c0c39ccd /bin/bash

忽略 Dockerfile的构建缓存, 每次都是最新版本，所有的都会执行一次（缓存apt-get update）
docker build  --no-cache -t shopcartApi:latest .   

### 查看镜像的Dockerfile指令
docker history id  或 wsjcko/mysql57

### 删除本地镜像
docker rmi id

## network
Docker四种网络模式，默认bridge

### 1.bridge  配置 -net=bridge (默认该模式) 桥接网络模式
在该模式中，Docker 守护进程创建了一个虚拟以太网桥docker0，新建的容器会自动桥接到这个接口，附加在其上的任何网卡之间都能自动转发数据包。

默认情况下，守护进程会创建一对对虚拟设备接口 veth pair, 将其中一个接口设置为容器的eth0接口（容器的网卡），另一个接口放置在宿主机的命名空间中，以类似vethxxxxx这样的名字命名，从而将宿主机上的所有容器都连接到这个内部网络上。


##### 新建的容器会自动桥接到这个接口
docker network inspect bridge 的Containers

或者
docker inspect -f '{{.NetworkSettings.IPAddress}}' mysql57
172.17.0.4
docker inspect -f '{{.NetworkSettings.IPAddress}}' consul
172.17.0.2
docker inspect -f '{{.NetworkSettings.IPAddress}}' jaeger
172.17.0.3
docker inspect -f '{{.NetworkSettings.IPAddress}}' hystrix
172.17.0.5
### 2.host模式 配置 -net=host 容器和宿主机共享Network namespace
如果启动容器的时候使用host模式，那么这个容器将不会获得一个独立的Network Namespace，而是和宿主机共用一个Network Namespace。容器将不会虚拟出自己的网卡，配置自己的IP等，而是使用宿主机的IP和端口。但是，容器的其他方面，如文件系统、进程列表等还是和宿主机隔离的。

使用host模式的容器可以直接使用宿主机的IP地址与外界通信，容器内部的服务端口也可以使用宿主机的端口，不需要进行NAT，host最大的优势就是网络性能比较好，但是docker host上已经使用的端口就不能再用了，网络的隔离性不好。

### 3.container模式  配置 -net=container: name or id
容器和另外一个容器共享Network namespace, kubernetes中的pod就是多个容器共享一个Network namespace

这个模式指定新创建的容器和已经存在的一个容器共享一个 Network Namespace，而不是和宿主机共享。新创建的容器不会创建自己的网卡，配置自己的 IP，而是和一个指定的容器共享 IP、端口范围等。同样，两个容器除了网络方面，其他的如文件系统、进程列表等还是隔离的。两个容器的进程可以通过 lo 网卡设备通信。

### 4.none模式 -net=none 容器有独立的Network namespace,但并没有对其进行任何网络设置，如分配veth pair和网络桥接，配置IP等
使用none模式，Docker容器拥有自己的Network Namespace，但是，并不为Docker容器进行任何网络配置。也就是说，这个Docker容器没有网卡、IP、路由等信息。需要我们自己为Docker容器添加网卡、配置IP等。

这种网络模式下容器只有lo回环网络，没有其他网卡。none模式可以在容器创建时通过–network=none来指定。这种类型的网络没有办法联网，封闭的网络能很好的保证容器的安全性。


docker network -h  查看network使用方法
docker network ls (默认，bridge，local,nono)
docker network create shopnet 创建docker网络shopnet    IP范围都可以指定  --subnet=172.20.10.0/16
docker network inspect shopnet 查看shopnet网络详情
docker network rm shopnet  删除
docker run -d --name redis --net shopnet --ip 172.20.10.5 redis:latest   创建Redis容器，指定网络和IP


docker port consul
docker start consul -v /data/www/consul:/data/consul


##### 添加已有容器到 shopnet网络

docker network connect shopnet consul

这样consul 在brige,shopnet网络中都有对应的IP

sudo iptables -t nat -L -n

##### 查看容器内的进程
docker top consul

UID       PID                 PPID    STIME    TIME                CMD
root      3116                3096    09:19    00:00:00            /usr/bin/dumb-init /bin/sh /usr/local/bin/docker-entrypoint.sh agent -dev -client 0.0.0.0
_apt      3153                3116    09:19    00:00:10            consul agent -data-dir=/consul/data -config-dir=/consul/config -dev -client 0.0.0.0

##### 统计信息
docker stats consul

CONTAINER ID   NAME      CPU %     MEM USAGE / LIMIT     MEM %     NET I/O           BLOCK I/O   PIDS
7746c468aade   consul    1.33%     24.21MiB / 12.37GiB   0.19%     15.5kB / 10.3kB   0B / 0B     18
CONTAINER ID   NAME      CPU %     MEM USAGE / LIMIT     MEM %     NET I/O           BLOCK I/O   PIDS
7746c468aade   consul    1.25%     24.21MiB / 12.37GiB   0.19%     15.5kB / 10.3kB   0B / 0B     18

##### 自动重启容器
docker run --restart=always
docker run --restart=on-failure:5 非正常退出，最多重启5次

##### 深入容器，查看更多的容器信息
docker inspect consul

docker inspect -f '{{.State.Running}}' consul




##### 监控告警Prometheus
Prometheus Agent 是 Prometheus 在 2.32.0 版本推出的实验性功能，当启用此功能后，将不会在本地文件系统上生成块，并且无法在本地查询。如果网络出现异常状态无法连接至远程端点，数据将暂存在本地磁盘，但是仅限于两个小时的缓冲。
remote-write-receiver 是 Prometheus 在 2.25 版本推出的实验性功能，当启用此功能后，prometheus 可以作为另一个 prometheus 的远程存储。



日志接入ELK，到最后k8s部署