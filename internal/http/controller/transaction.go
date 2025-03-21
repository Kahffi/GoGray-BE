package controller

import (
	"net/http"

	"github.com/Kahffi/GoGray-BE/internal/service"
	"github.com/Kahffi/GoGray-BE/models"

	"github.com/labstack/echo"
)

type TransactionController interface {
	Create(ctx echo.Context) error
}

type transactionControllerImpl struct {
	TransactionService service.TransactionService
}

func NewTransactionController(transactionService service.TransactionService) TransactionController {
	return &transactionControllerImpl{
		TransactionService: transactionService,
	}
}

func (controller *transactionControllerImpl) Create(ctx echo.Context) error {
	request := new(models.TransactionCreateRequest)

	if err := ctx.Bind(request); err != nil {
		return ctx.JSON(http.StatusBadRequest, models.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "Bad Request",
			Data:   nil,
		})
	}

	err := controller.TransactionService.CreateTransaction(*request)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, models.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "Internal Server Error",
			Data:   nil,
		})
	}

	return ctx.JSON(http.StatusAccepted, models.WebResponse{
		Code:   http.StatusAccepted,
		Status: "Image Uploaded",
		Data:   nil,
	})
}
