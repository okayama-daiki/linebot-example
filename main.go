package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

var (
	channelSecret      = os.Getenv("LINE_CHANNEL_SECRET")
	channelAccessToken = os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")
	bot                *messaging_api.MessagingApiAPI
)

func init() {
	var err error
	bot, err = messaging_api.NewMessagingApiAPI(
		channelAccessToken,
	)
	if err != nil {
		log.Fatal(err)
	}
}

func lambdaHandler(req events.LambdaFunctionURLRequest) (events.APIGatewayProxyResponse, error) {
	httpReq := http.Request{
		Header: http.Header{
			"X-Line-Signature": []string{req.Headers["x-line-signature"]},
		},
		Body: io.NopCloser(strings.NewReader(req.Body)),
	}
	cb, err := webhook.ParseRequest(channelSecret, &httpReq)
	if err != nil {
		log.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
		}, err
	}

	for _, event := range cb.Events {
		switch e := event.(type) {
		case webhook.MessageEvent:
			switch message := e.Message.(type) {
			case webhook.TextMessageContent:
				if _, err := bot.ReplyMessage(
					&messaging_api.ReplyMessageRequest{
						ReplyToken: e.ReplyToken,
						Messages: []messaging_api.MessageInterface{
							messaging_api.TextMessage{
								Text: message.Text,
							},
						},
					},
				); err != nil {
					log.Println(err)
				} else {
					log.Println("Replied to message", message.Text)
				}
			}
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(lambdaHandler)
}
