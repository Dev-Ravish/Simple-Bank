package api

import (
	"net/http"
	db "simplebank/db/sqlc"
	"strconv"

	"github.com/gin-gonic/gin"
)

type createAccountParams struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required, oneof"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountParams

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err)) //this is what we are sending to the user, the first one is the status code and the second one is JSON that we want to display
		return
	}

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, account)

}

func (server *Server) getAccount(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}
	account, err := server.store.GetAccount(ctx, int64(id))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, account)

}
