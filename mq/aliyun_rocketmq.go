package mq

import (
	"strings"
	"time"

	mq_http_sdk "github.com/aliyunmq/mq-http-go-sdk"
	"github.com/glory-go/glory/config"
	"github.com/gogap/errors"
)

type aliyunRocketMqSrv struct {
	config *aliyunRocketMqConf
	client *mq_http_sdk.AliyunMQClient
}

type aliyunRocketMqConf struct {
	endpoint  string
	accessKey string
	secretKey string

	instanceId      string
	consumerGroupID string
}

func newAliyunRocketMqSrv() *aliyunRocketMqSrv {
	return &aliyunRocketMqSrv{
		config: &aliyunRocketMqConf{},
	}
}

func (r *aliyunRocketMqSrv) loadConfig(conf *config.MQConfig) error {
	r.config.endpoint = conf.Endpoint
	r.config.accessKey = conf.AccessKey
	r.config.secretKey = conf.SecretKey
	r.config.instanceId = conf.InstanceID
	r.config.consumerGroupID = conf.ConsumerGroupID
	return nil
}

func (r *aliyunRocketMqSrv) Send(topic string, msg []byte) (string, error) {
	return r.send(topic, msg, nil)
}

func (r *aliyunRocketMqSrv) DelaySend(topic string, msg []byte, handleTime time.Time) (string, error) {
	return r.send(topic, msg, &handleTime)
}

func (r *aliyunRocketMqSrv) Handler(topic string, handler MQMsgHandler) error {
	consumer := r.client.GetConsumer(r.config.instanceId, topic, r.config.consumerGroupID, "")
	var exitErr error
	for {
		if exitErr != nil {
			break
		}
		endChan := make(chan int)
		respChan := make(chan mq_http_sdk.ConsumeMessageResponse)
		errChan := make(chan error)

		go func() {
			select {
			case resp := <-respChan:
				// 处理业务逻辑
				successHandles := []string{}
				for _, v := range resp.Messages {
					msg := v.MessageBody
					err := handler([]byte(msg))
					if err == nil {
						// 当处理成功后，确认消息消费成功
						successHandles = append(successHandles, v.ReceiptHandle)
					}
					// 处理不成功时，消息后续会被重新投递
				}

				// TODO: 某些消息的句柄可能超时了会导致确认不成功
				consumer.AckMessage(successHandles)
				endChan <- 1
			case err := <-errChan:
				// 没有消息，则继续循环
				if !strings.Contains(err.(errors.ErrCode).Error(), "MessageNotExist") {
					exitErr = err
				}
				endChan <- 1
			case <-time.After(35 * time.Second):
				// 长连接超时，此时重新连接
				endChan <- 1
			}
		}()

		// 长轮询消费消息
		// 长轮询表示如果topic没有消息则请求会在服务端挂住3s，3s内如果有消息可以消费则立即返回
		consumer.ConsumeMessage(respChan, errChan,
			3, // 一次最多消费3条(最多可设置为16条)
			3, // 长轮询时间3秒（最多可设置为30秒）
		)
		<-endChan
	}
	return exitErr
}

func (r *aliyunRocketMqSrv) send(topic string, msg []byte, deliverTime *time.Time) (string, error) {
	produser := r.client.GetProducer(r.config.instanceId, topic)
	mqMsg := mq_http_sdk.PublishMessageRequest{
		MessageBody: string(msg),         // 消息内容
		MessageTag:  "",                  // 消息标签
		Properties:  map[string]string{}, // 消息属性
	}
	if deliverTime != nil {
		mqMsg.StartDeliverTime = deliverTime.Unix()
	}
	ret, err := produser.PublishMessage(mqMsg)

	if err != nil {
		return "", err
	}
	return ret.MessageId, nil
}
