package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openai "github.com/sashabaranov/go-openai"
)

type DataRequest struct {
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
	Message  string `json:"message"`
}
type DataRespnse struct {
	Response string `json:"response"`
	Username string `json:"username"`
	Id       int    `json:"id"`
}

func main() {
	r := gin.Default()
	r.POST("/data", func(ctx *gin.Context) {
		var dataResponse DataRequest
		if err := ctx.ShouldBindJSON(&dataResponse); err != nil {
			ctx.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		// send server api response
		respmsg, err := Open(dataResponse.Message)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"Err":      err.Error(),
				"username": dataResponse.Username,
				"id":       dataResponse.UserId,
			})
		}
		resp := &DataRespnse{
			Response: respmsg,
			Username: dataResponse.Username,
			Id:       dataResponse.UserId,
		}
		ctx.JSON(200, resp)
	})

	r.Run(":5005")

}
func Open(msg string) (string, error) {
	client := openai.NewClient("sk-proj-uHDfqte1xburIyeEtA1tT3BlbkFJZroDT9Eh114GmBPlDm2P")
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: msg,
				},
			},
		},
	)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}
