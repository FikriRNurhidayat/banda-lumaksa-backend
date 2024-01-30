package transaction

import (
	"net/http"

	"github.com/fikrirnurhidayat/banda-lumaksa/internal/common/schema"
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/common/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TransactionController interface {
	Register(*echo.Echo)
	ListTransactions(c echo.Context) error
	GetTransaction(c echo.Context) error
}

type TransactionControllerImpl struct {
	transactionService TransactionService
}

func (ctl *TransactionControllerImpl) Register(e *echo.Echo) {
	e.GET("/v1/transactions/:id", ctl.GetTransaction)
	e.GET("/v1/transactions", ctl.ListTransactions)
}

func (ctl *TransactionControllerImpl) GetTransaction(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Logger().Error(err.Error())
		return err
	}

	params := &GetTransactionParams{
		ID: id,
	}

	result, err := ctl.transactionService.GetTransaction(c.Request().Context(), params)
	if err != nil {
		c.Logger().Error(err.Error())
		return err
	}

	response := &GetTransactionResponse{
		Transaction: NewTransactionResponse(result.Transaction),
	}

	return c.JSON(http.StatusOK, response)
}

func (ctl *TransactionControllerImpl) ListTransactions(c echo.Context) error {
	params := &ListTransactionsParams{
		Pagination: service.PaginationParams{},
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
		c.Logger().Error(err.Error())
		return err
	}

	response := &ListTransactionsResponse{
		PaginationResponse: schema.NewPaginationResponse(result.Pagination),
		Transactions:       NewTransactionsResponse(result.Transactions),
	}

	return c.JSON(http.StatusOK, response)
}

func NewTransactionController(transactionService TransactionService) TransactionController {
	return &TransactionControllerImpl{
		transactionService: transactionService,
	}
}
