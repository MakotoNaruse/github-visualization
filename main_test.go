package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_index(t *testing.T) {
	ginContext, _ := gin.CreateTestContext(httptest.NewRecorder())
	req, _ := http.NewRequest("GET", "/", nil)
	ginContext.Request = req
	store := cookie.NewStore([]byte("secret"))
	handlerFunc := sessions.Sessions("github-visualizer", store)
	handlerFunc(ginContext)

	type args struct {
		c *gin.Context
	}
	//dummyToken := "dummy"
	tests := []struct {
		name  string
		args  args
		token *string
	}{
		/*
			{
				name: "sessionがある場合",
				args: args{ c: ginContext },
				token: &dummyToken,
			}*/
		{
			name:  "sessionがない場合",
			args:  args{c: ginContext},
			token: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index(tt.args.c)
			fmt.Println(tt.args.c.Writer.Status())
		})
	}
}
