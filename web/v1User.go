package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) v1UserGETHandle(ctx *gin.Context) {
	userID := ctx.Param("user")
	user, err := s.query.GetUser(userID)
	if err != nil {
		// TODO handle later.
		return
	}

	ctx.JSON(http.StatusOK, user)
}
