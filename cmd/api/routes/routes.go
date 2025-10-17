package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (app *application) routes() http.Handler {
	g := gin.Default()

	return g
}
