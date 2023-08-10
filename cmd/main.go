package main

import (
	"github.com/gin-gonic/gin"
	"racing-restful/internal/handlers"
)

func setupRouter() *gin.Engine {
	// Logger And Recovery Middleware 붙여진 상태로 나옴
	// Recover -> 패닉 발생 되면 앱종료가 아닌 오류 보내기
	// Logger -> 각 요청에 대한 로그 출력
	r := gin.Default()

	handlers.NewRacingHandler().AddToRouter("racing", r)

	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080

	err := r.Run(":8080")
	if err != nil {
		panic("gin start error")
	}
}
