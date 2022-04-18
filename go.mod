module github.com/wsjcko/shoppaymentApi

go 1.16

require (
	github.com/afex/hystrix-go v0.0.0-20180502004556-fa1af6a1f4f5
	github.com/asim/go-micro/plugins/config/source/consul/v4 v4.0.0-20220415015530-034ba9a0dead
	github.com/asim/go-micro/plugins/registry/consul/v4 v4.0.0-20220415015530-034ba9a0dead
	github.com/asim/go-micro/plugins/wrapper/select/roundrobin/v4 v4.0.0-20220407072057-664ee31b1a6b
	github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v4 v4.0.0-20220415015530-034ba9a0dead
	github.com/golang/protobuf v1.5.2
	github.com/jinzhu/gorm v1.9.16
	github.com/opentracing/opentracing-go v1.2.0
	github.com/plutov/paypal/v4 v4.5.1
	github.com/prometheus/client_golang v1.11.0
	github.com/prometheus/common v0.26.0
	github.com/uber/jaeger-client-go v2.30.0+incompatible
	github.com/urfave/cli/v2 v2.3.0
	github.com/wsjcko/shoppayment v1.0.0
	go-micro.dev/v4 v4.6.0
	go.uber.org/zap v1.10.0
	google.golang.org/protobuf v1.26.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
)

// This can be removed once etcd becomes go gettable, version 3.4 and 3.5 is not,
// see https://github.com/etcd-io/etcd/issues/11154 and https://github.com/etcd-io/etcd/issues/11931.
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace github.com/wsjcko/shoppaymentApi => github.com/wsjcko/shoppaymentApi v1.0.1
