package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"racing-restful/internal/models"
	"strconv"
)

type RacingHandler struct {
}

type createRequest struct {
	racing models.Racing `json:"racing"`
	users  []string      `json:"users"`
}

func NewRacingHandler() *RacingHandler {
	return &RacingHandler{}
}

func (h *RacingHandler) AddToRouter(path string, router *gin.Engine) {
	//CRUD 추가하고싶음
	// C -> name,turns로 user, racing 둘 다 생성
	// R -> racing 정보 조회 -> 유저들 정보 조회
	// U -> racing 정보 수정 (최대 도는 횟수 user들 validation)
	// D -> racing 삭제 시 cascade 로 다같이 삭제 (user도)
	group := router.Group(path)

	group.GET("/:racingID", h.getRacing)
	group.POST("/", h.creatRacing)
	group.PATCH("/")
	group.DELETE("/")
}

// GetRacing: Racing 조회
func (h *RacingHandler) getRacing(c *gin.Context) {
	racingId, err := strconv.ParseInt(c.Param("racingID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "parse error : " + err.Error()})
	}

	racing, ok := models.Racings[racingId]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"message": strconv.FormatInt(racingId, 10) + "is not found"})
	}

	r, err := json.Marshal(racing)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "racing can't marshal")
	}

	c.JSON(http.StatusOK, gin.H{"racing": string(r)})
}

func (h *RacingHandler) creatRacing(c *gin.Context) {
	readCloser := c.Request.Body
	body, err := io.ReadAll(readCloser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}

	var r map[string]interface{}

	err = json.Unmarshal(body, &r)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "parse error : " + err.Error()})
	}

	models.AddRacing(r)
	c.JSON(http.StatusCreated, gin.H{"racing": r, "id": len(models.Racings)})
}
