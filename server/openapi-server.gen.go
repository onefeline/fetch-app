package server

import (
	"fmt"
	"github.com/labstack/echo"
	"net/http"

	"github.com/oapi-codegen/runtime"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Autogenerated using oapi-codegen

// Item defines model for Item.
type Item struct {
	// Price The total price payed for this item.
	Price string `json:"price"`

	// ShortDescription The Short Product Description for the item.
	ShortDescription string `json:"shortDescription"`
}

// Receipt defines model for Receipt.
type Receipt struct {
	Items []Item `json:"items"`

	// PurchaseDate The date of the purchase printed on the receipt.
	PurchaseDate openapi_types.Date `json:"purchaseDate"`

	// PurchaseTime The time of the purchase printed on the receipt. 24-hour time expected.
	PurchaseTime string `json:"purchaseTime"`

	// Retailer The name of the retailer or store the receipt is from.
	Retailer string `json:"retailer"`

	// Total The total amount paid on the receipt.
	Total string `json:"total"`
}

// PostReceiptsProcessJSONRequestBody defines body for PostReceiptsProcess for application/json ContentType.
type PostReceiptsProcessJSONRequestBody = Receipt

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Submits a receipt for processing
	// (POST /receipts/process)
	PostReceiptsProcess(ctx echo.Context) error
	// Returns the points awarded for the receipt
	// (GET /receipts/{id}/points)
	GetReceiptsIdPoints(ctx echo.Context, id string) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// PostReceiptsProcess converts echo context to params.
func (w *ServerInterfaceWrapper) PostReceiptsProcess(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostReceiptsProcess(ctx)
	return err
}

// GetReceiptsIdPoints converts echo context to params.
func (w *ServerInterfaceWrapper) GetReceiptsIdPoints(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetReceiptsIdPoints(ctx, id)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST(baseURL+"/receipts/process", wrapper.PostReceiptsProcess)
	router.GET(baseURL+"/receipts/:id/points", wrapper.GetReceiptsIdPoints)

}
