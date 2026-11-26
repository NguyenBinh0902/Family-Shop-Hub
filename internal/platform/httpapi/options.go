// Package httpapi contains the types and Gin server interface generated from
// api/openapi.yaml (httpapi.gen.go, never edited by hand) plus the options in
// this file that route every error through platform/problem.
package httpapi

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"familyshop-hub/internal/platform/problem"
)

// Register mounts handler on r under the spec's base URL (/api/v1), with all
// errors rendered as Problem Details.
//
// The generated defaults would send err.Error() to the client, leaking
// internal details; these options replace them.
func Register(r gin.IRouter, handler StrictServerInterface, logger *slog.Logger) {
	writeAs := func(status int) func(*gin.Context, error) {
		return func(c *gin.Context, err error) { problem.Write(c, logger, err, status) }
	}
	strict := NewStrictHandlerWithOptions(handler, nil, StrictGinServerOptions{
		RequestErrorHandlerFunc:  writeAs(http.StatusBadRequest),
		HandlerErrorFunc:         writeAs(http.StatusInternalServerError),
		ResponseErrorHandlerFunc: writeAs(http.StatusInternalServerError),
	})
	RegisterHandlersWithOptions(r, strict, GinServerOptions{
		BaseURL: BaseURL,
		// Called when path or query parameters fail to parse.
		ErrorHandler: func(c *gin.Context, err error, status int) {
			problem.Write(c, logger, err, status)
		},
	})
}

// BaseURL is the servers.url of the spec; every operation path is relative
// to it.
const BaseURL = "/api/v1"
