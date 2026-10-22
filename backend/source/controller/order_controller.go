package controller

import (
	"log"
	"net/http"
	"stockflow/dto"
	"stockflow/model"
	"stockflow/service"
	"stockflow/util"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	service *service.OrderService
	router  *gin.RouterGroup
}

func NewOrderController(router *gin.RouterGroup) *OrderController {
	c := &OrderController{
		service: service.NewOrderService(),
		router:  router,
	}

	order := c.router.Group("/orders")
	order.POST("/", c.Create)
	order.POST("/:id/cancel", c.Cancel)
	order.GET("/:id", c.GetDetail)

	return c
}

// @Tags order
// @Accept json
// @Param parameter body dto.OrderCreateRequest true "param"
// @Produce json
// @Success 201 {object} object{metadata=model.MetadataResponse} "OK"
// @Router /orders [post]
func (o *OrderController) Create(ctx *gin.Context) {
	startTime := time.Now()
	var resp model.Response

	var param dto.OrderCreateRequest
	if err := ctx.BindJSON(&param); err != nil {
		log.Println(err)
		resp.Metadata.Message = "invalid json or validation"
		util.SetMetadataResponse(ctx, http.StatusBadRequest, startTime, &resp)
		return
	}

	resp = o.service.Create(param)
	statusCode := http.StatusCreated

	switch {
	case resp.Metadata.Message == "duplicate sku in request":
		statusCode = http.StatusBadRequest
	case resp.Metadata.Message == "insufficient stock":
		statusCode = http.StatusConflict
	case resp.Metadata.Message != "OK":
		statusCode = http.StatusInternalServerError
	}

	util.SetMetadataResponse(ctx, statusCode, startTime, &resp)
}

// @Tags order
// @Accept json
// @Param id path int true "param"
// @Produce json
// @Success 201 {object} object{metadata=model.MetadataResponse} "OK"
// @Router /orders/{id}/cancel [post]
func (o *OrderController) Cancel(ctx *gin.Context) {
	startTime := time.Now()
	var resp model.Response

	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		log.Println(err)
		resp.Metadata.Message = "invalid order id"
		util.SetMetadataResponse(ctx, http.StatusBadRequest, startTime, &resp)
		return
	}

	resp = o.service.Cancel(id)
	statusCode := http.StatusOK

	switch {
	case resp.Metadata.Message == "order already cancelled":
		statusCode = http.StatusConflict
	case resp.Metadata.Message != "OK":
		statusCode = http.StatusInternalServerError
	}

	util.SetMetadataResponse(ctx, statusCode, startTime, &resp)
}

// @Tags order
// @Produce json
// @Param id path int true "page"
// @Success 200 {object} object{metadata=model.MetadataResponse, data=model.PaginatedResponse} "OK"
// @Router /orders/{id} [get]
func (o *OrderController) GetDetail(ctx *gin.Context) {
	startTime := time.Now()
	var resp model.Response

	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		log.Println(err)
		resp.Metadata.Message = "invalid order id"
		util.SetMetadataResponse(ctx, http.StatusBadRequest, startTime, &resp)
		return
	}

	resp = o.service.GetDetail(id)
	statusCode := http.StatusOK

	if resp.Metadata.Message == "order not found" {
		statusCode = http.StatusNotFound
	}

	util.SetMetadataResponse(ctx, statusCode, startTime, &resp)
}
