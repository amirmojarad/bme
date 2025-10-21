package controller

import (
	"bme/pkg/errorext"
	"bme/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type HeaderEntity struct {
	UserID uint `header:"X-Auth-User-ID"`
}

type HeaderEntityBindingRequired struct {
	UserID uint `header:"X-Auth-User-ID" binding:"required"`
}

func writeValidationErrors(ctx *gin.Context, validationError *utils.ValidationError) {
	ctx.JSON(http.StatusOK, *validationError)
}

func writeBindingErrorResponse(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
}

func writeBadRequestErrorResponse(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
}

func writeErrorResponse(ctx *gin.Context, err error, logger *logrus.Entry) {
	var extErr *errorext.Error

	if errors.As(err, &extErr) {
		logger.
			WithField("path", ctx.Request.URL.Path).
			WithField("status", extErr.StatusCode).
			WithField("stackTrace", extErr.StackTrace).
			Error(extErr.Error())

		ctx.JSON(extErr.StatusCode, gin.H{"message": extErr.Error(), "key": extErr.Key})

		return
	}

	ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
}

func convertQ[T uint](q *string) (*T, *string) {
	if q == nil {
		return nil, nil
	}

	numericVal, err := strconv.Atoi(*q)
	if err != nil {
		return nil, q
	}

	numericValAsUint := T(uint(numericVal))
	return &numericValAsUint, nil
}
