package main

import (
	"context"
	"net"
	"net/http"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/asim/go-micro/plugins/registry/consul/v4"
	"github.com/asim/go-micro/plugins/wrapper/select/roundrobin/v4"
	opentracing4 "github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v4"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/opentracing/opentracing-go"
	cli2 "github.com/urfave/cli/v2"
	paymentPb "github.com/wsjcko/shoppayment/protobuf/pb"
	"github.com/wsjcko/shoppaymentApi/common"
	"github.com/wsjcko/shoppaymentApi/handler"
	"github.com/wsjcko/shoppaymentApi/logger"
	pb "github.com/wsjcko/shoppaymentApi/protobuf/pb"
	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
	log "go-micro.dev/v4/logger"
	"go-micro.dev/v4/registry"
)

var (
	MICRO_API_NAME        = "go.micro.api.shopPaymentApi" //决定路由： shopPaymentApi/
	MICRO_SERVICE_NAME    = "go.micro.service.shop.payment"
	MICRO_VERSION         = "latest"
	MICRO_ADDRESS         = "0.0.0.0:8088"
	MICRO_HYSTRIX_HOST    = "0.0.0.0"
	MICRO_HYSTRIX_PORT    = "9096"
	MICRO_CONSUL_ADDRESS  = "127.0.0.1:8500"
	MICRO_JAEGER_ADDRESS  = "127.0.0.1:6831"
	DOCKER_HOST           = "127.0.0.1"
	MICRO_PROMETHEUS_PORT = "9092"
)

func SetDockerHost(host string) {
	DOCKER_HOST = host
	MICRO_CONSUL_ADDRESS = host + ":8500"
	MICRO_JAEGER_ADDRESS = host + ":6831"
}

func main() {

	function := micro.NewFunction(
		micro.Flags(
			&cli2.StringFlag{ //micro 多个选项 --ip
				Name:  "ip",
				Usage: "docker Host IP(ubuntu)",
				Value: "0.0.0.0",
			},
		),
	)

	function.Init(
		micro.Action(func(c *cli2.Context) error {
			ipstr := c.Value("ip").(string)
			log.Info("docker Host IP(ubuntu)1111", ipstr)
			if net.ParseIP(ipstr) == nil {
				return nil
			}
			SetDockerHost(ipstr)
			return nil
		}),
	)

	log.Info("DOCKER_HOST ", DOCKER_HOST)

	logger.Init("ShopPaymentApi", "micro.log")

	//注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			MICRO_CONSUL_ADDRESS,
		}
	})

	//链路追踪
	t, io, err := common.NewTracer(MICRO_API_NAME, MICRO_JAEGER_ADDRESS)
	if err != nil {
		logger.Error(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	// 熔断器
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	// 启动端口 启动监听 上报熔断状态
	go func() {
		err = http.ListenAndServe(net.JoinHostPort(MICRO_HYSTRIX_HOST, MICRO_HYSTRIX_PORT), hystrixStreamHandler)
		if err == http.ErrServerClosed {
			logger.Info("httpserver shutdown cased: ", err)
		} else {
			logger.Error(err)
		}
	}()

	// 暴露监控地址
	common.PrometheusBoot(MICRO_PROMETHEUS_PORT)

	// New Service
	srv := micro.NewService(
		micro.Name(MICRO_API_NAME),
		micro.Version(MICRO_VERSION),
		micro.Address(MICRO_ADDRESS),
		//添加 consul 注册中心
		micro.Registry(consulRegistry),
		//添加链路追踪 服务端绑定handle 客户端绑定client
		micro.WrapHandler(opentracing4.NewHandlerWrapper(opentracing.GlobalTracer())),
		micro.WrapClient(opentracing4.NewClientWrapper(opentracing.GlobalTracer())),
		//添加熔断
		micro.WrapClient(NewClientHystrixWrapper()),
		//添加负载均衡
		micro.WrapClient(roundrobin.NewClientWrapper()),
	)

	// Initialise service
	srv.Init()

	// 调用后端服务
	shopPaymentService := paymentPb.NewShopPaymentService(MICRO_SERVICE_NAME, srv.Client())

	// Register Handler
	if err := pb.RegisterShopPaymentApiHandler(srv.Server(), &handler.ShopPaymentApi{PaymentService: shopPaymentService}); err != nil {
		logger.Error(err)
	}

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}

type clientWrapper struct {
	client.Client
}

func (c *clientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	return hystrix.Do(req.Service()+"."+req.Endpoint(), func() error {
		//run 正常执行
		logger.Info(req.Service() + "." + req.Endpoint())
		return c.Client.Call(ctx, req, rsp, opts...)
	}, func(err error) error {
		logger.Error(err)
		return err
	})
}

func NewClientHystrixWrapper() client.Wrapper {
	return func(i client.Client) client.Client {
		return &clientWrapper{i}
	}
}
