package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	models "racing-restful/internal/models"
	"strconv"
)

type RacingHandler struct {
}

type CreateRequest struct {
	Racing models.Racing `json:"racing"`
	Users  []string      `json:"users"`
}

type UpdateRequest struct {
	MaxTurns int `json:"max_turns"`
}

func NewRacingHandler() *RacingHandler {
	return &RacingHandler{}
}

func (h *RacingHandler) AddToRouter(path string, router *gin.Engine) {
	group := router.Group(path)

	group.GET("/:racingID", h.get)
	group.POST("/", h.create)
	group.PUT("/:racingID", h.update)
	group.DELETE("/:racingID", h.delete)
}

func (h *RacingHandler) get(c *gin.Context) {
	racingId, err := strconv.ParseInt(c.Param("racingID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "parse error : " + err.Error()})
	}

	racing, ok := models.Racings[racingId]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"message": strconv.FormatInt(racingId, 10) + " is not found"})
	}

	var users []models.User
	for _, user := range models.Users {
		if user.RacingId == racingId {
			users = append(users, *user)
		}
	}

	r, err := json.Marshal(racing)
	u, err := json.Marshal(users)

	if err != nil {
		c.JSON(http.StatusInternalServerError, "racing can't marshal")
	}

	c.JSON(http.StatusOK, gin.H{"racing": string(r), "users": string(u)})
}

func (h *RacingHandler) create(c *gin.Context) {
	readCloser := c.Request.Body
	body, err := io.ReadAll(readCloser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}

	var r CreateRequest
	err = json.Unmarshal(body, &r)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "parse error : " + err.Error()})
	}

	racing := r.Racing
	racingId := models.AddRacing(&racing)

	var userIds []int64

	for _, name := range r.Users {
		userId := models.AddUser(models.NewUser(name, racing.MaxTurns, racingId))
		userIds = append(userIds, userId)
	}

	c.JSON(http.StatusCreated, gin.H{"racingId": racingId, "userIds": userIds})
}

func (h *RacingHandler) update(c *gin.Context) {
	racingId, err := strconv.ParseInt(c.Param("racingId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "parse error : " + err.Error()})
	}

	readCloser := c.Request.Body
	body, err := io.ReadAll(readCloser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}

	var r UpdateRequest
	err = json.Unmarshal(body, &r)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "parse error : " + err.Error()})
	}

	var userIds []int64
	for id, user := range models.Users {
		if user.RacingId == racingId {
			models.UpdateTurns(id, r.MaxTurns)
			userIds = append(userIds, id)
		}
	}

	c.JSON(http.StatusOK, gin.H{"UserIds": userIds})
}

func (h *RacingHandler) delete(c *gin.Context) {
	racingId, err := strconv.ParseInt(c.Param("racingID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "parse error : " + err.Error()})
	}

	delete(models.Racings, racingId)

	for idx, user := range models.Users {
		if user.RacingId == racingId {
			delete(models.Users, idx)
		}
	}

	fmt.Print(models.Racings)
	c.JSON(http.StatusOK, gin.H{"message": "delete complete"})
}
