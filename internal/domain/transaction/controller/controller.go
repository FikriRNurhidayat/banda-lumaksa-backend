package transaction_controller

import (
	"net/http"

	common_errors "github.com/fikrirnurhidayat/banda-lumaksa/internal/common/errors"
	common_schema "github.com/fikrirnurhidayat/banda-lumaksa/internal/common/schema"
	common_service "github.com/fikrirnurhidayat/banda-lumaksa/internal/common/service"

	transaction_service "github.com/fikrirnurhidayat/banda-lumaksa/internal/domain/transaction/service"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TransactionController interface {
	Register(*echo.Echo)
	ListTransactions(c echo.Context) error
	GetTransaction(c echo.Context) error
}

type TransactionControllerImpl struct {
	transactionService transaction_service.TransactionService
}

func (ctl *TransactionControllerImpl) Register(e *echo.Echo) {
	e.GET("/v1/transactions/:id", ctl.GetTransaction)
	e.GET("/v1/transactions", ctl.ListTransactions)
}

func (ctl *TransactionControllerImpl) GetTransaction(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return common_errors.ErrInvalidUUID
	}

	params := &transaction_service.GetTransactionParams{
		ID: id,
	}

	result, err := ctl.transactionService.GetTransaction(c.Request().Context(), params)
	if err != nil {
		return err
	}

	response := &GetTransactionResponse{
		Transaction: NewTransactionResponse(result.Transaction),
	}

	return c.JSON(http.StatusOK, response)
}

func (ctl *TransactionControllerImpl) ListTransactions(c echo.Context) error {
	params := &transaction_service.ListTransactionsParams{
		Pagination: common_service.PaginationParams{},
	}

	if err := echo.QueryParamsBinder(c).
		String("description_like", &params.DescriptionLike).
		Uint32("page", &params.Pagination.Page).
		Uint32("page_size", &params.Pagination.PageSize).
		FailFast(true).
		BindError(); err != nil {
		c.Logger().Error(err.Error())
		return err
	}

	result, err := ctl.transactionService.ListTransactions(c.Request().Context(), params)
	if err != nil {
		return err
	}

	response := &ListTransactionsResponse{
		PaginationResponse: common_schema.NewPaginationResponse(result.Pagination),
		Transactions:       NewTransactionsResponse(result.Transactions),
	}

	return c.JSON(http.StatusOK, response)
}

func New(transactionService transaction_service.TransactionService) TransactionController {
	return &TransactionControllerImpl{
		transactionService: transactionService,
	}
}
