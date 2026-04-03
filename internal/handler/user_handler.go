package handler

import (
	"LiveChat/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	model.Service
}

func NewHandler(s model.Service) *Handler {
	return &Handler{
		Service: s,
	}
}

func (h *Handler) CreateUser(c *gin.Context) {
	var u model.CreateUserReq
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.Service.CreateUser(c.Request.Context(), &u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h *Handler) Login(c *gin.Context) {
	var user model.LoginUserReq
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := h.Service.Login(c.Request.Context(), &user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := &model.LoginUserRes{
		AccessToken: u.AccessToken,
		Username:    u.Username,
		ID:          u.ID,
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h *Handler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}
