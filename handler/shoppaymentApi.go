package handler

import (
	"context"
	"errors"
	"github.com/plutov/paypal/v4"
	paymentPb "github.com/wsjcko/shoppayment/protobuf/pb"
	"github.com/wsjcko/shoppaymentApi/logger"
	pb "github.com/wsjcko/shoppaymentApi/protobuf/pb"

	"strconv"
)

type ShopPaymentApi struct {
	PaymentService paymentPb.ShopPaymentService
}

var (
	//写你自己paypal上面的应用的ClientID
	ClientID string = "AYS3EN1rZs8vKUAXdCSagPe6qJwicatmDwWMIEJ0ny5cxYYqB7qIn7CTiX7Br2Rs9j7s988aWjn2xaGL"
)

// paymentPb.PayPalRefund 通过API向外暴露为/shopPaymentApi/payPalRefund，接收http请求
// 即：/shopPaymentApi/payPalRefund 请求会调用 go.micro.api.shop.paymentApi 服务的paymentPb.PayPalRefund
func (e *ShopPaymentApi) PayPalRefund(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	//验证payment 支付通道是否赋值
	if err := CheckParam("payment_id", req); err != nil {
		rsp.StatusCode = 500
		return err
	}
	//验证 退款号
	if err := CheckParam("refund_id", req); err != nil {
		rsp.StatusCode = 500
		return err
	}
	//验证 退款金额
	if err := CheckParam("money", req); err != nil {
		rsp.StatusCode = 500
		return err
	}

	//获取paymentID
	payID, err := strconv.ParseInt(req.Get["payment_id"].Values[0], 10, 64)
	if err != nil {
		logger.Error(err)
		return err
	}
	//获取支付通道信息
	paymentInfo, err := e.PaymentService.FindPaymentByID(ctx, &paymentPb.PaymentID{PaymentId: payID})
	if err != nil {
		logger.Error(err)
		return err
	}
	//SID 获取 paymentInfo.PaymentSid
	//支付模式
	status := paypal.APIBaseSandBox //沙盒环境
	if paymentInfo.PaymentStatus {
		status = paypal.APIBaseLive //生产环境
	}
	//退款例子
	refundId := req.Get["refund_id"].Values[0]
	payout := paypal.Payout{
		SenderBatchHeader: &paypal.SenderBatchHeader{
			EmailSubject: refundId + " wsjcko 提醒你收款！",
			EmailMessage: refundId + " 您有一个收款信息！",
			//每笔转账都要唯一,paypal服务端幂等性
			SenderBatchID: refundId,
		},
		Items: []paypal.PayoutItem{
			{
				RecipientType: "EMAIL",
				Receiver:      "sb-sgroy15724941@personal.example.com",
				Amount: &paypal.AmountPayout{
					//币种
					Currency: "USD",
					Value:    req.Get["money"].Values[0],
				},
				Note:         refundId,
				SenderItemID: refundId,
			},
		},
	}
	//创建支付客户端
	payPalClient, err := paypal.NewClient(ClientID, paymentInfo.PaymentSid, status)
	if err != nil {
		logger.Error(err)
	}
	// 获取 token
	_, err = payPalClient.GetAccessToken()
	if err != nil {
		logger.Error(err)
	}
	paymentResult, err := payPalClient.CreateSinglePayout(payout)
	if err != nil {
		logger.Error(err)
	}
	logger.Info(paymentResult)
	rsp.Body = refundId + " 支付成功！"
	return err
}

func CheckParam(key string, req *pb.Request) error {
	if _, ok := req.Get[key]; !ok {
		err := errors.New(key + " 参数异常")
		logger.Error(err)
		return err
	}
	return nil
}
