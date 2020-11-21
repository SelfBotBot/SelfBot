package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) v1TreeGETHandle(ctx *gin.Context) {
	userID := ctx.Param("user")

	type hierarchy struct {
		Name     string      `json:"name"`
		Children []hierarchy `json:"children"`
	}

	from, infected, err := s.query.GetUserInfections(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	var ret = hierarchy{Name: from}
	for _, v := range infected {
		ret.Children = append(ret.Children, hierarchy{Name: v})
	}

	ctx.JSON(http.StatusOK, ret)
}
