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

type ProductController struct {
	service *service.ProductService
	router  *gin.RouterGroup
}

func NewProductController(router *gin.RouterGroup) *ProductController {
	c := &ProductController{
		service: service.NewProductService(),
		router:  router,
	}

	product := c.router.Group("/products")
	product.POST("/", c.Create)
	product.GET("/", c.Search)
	product.DELETE("/:id", c.Delete)
	product.PUT("/:id", c.Update)
	product.POST("/:id/stock", c.AdjustStock)
	product.GET("/:id/movement", c.History)

	return c
}

// @Tags product
// @Accept json
// @Param parameter body dto.ProductRequest true "param"
// @Produce json
// @Success 201 {object} object{metadata=model.MetadataResponse} "OK"
// @Router /products [post]
func (c *ProductController) Create(ctx *gin.Context) {
	startTime := time.Now()
	var resp model.Response

	var param dto.ProductRequest
	if err := ctx.BindJSON(&param); err != nil {
		log.Println(err)
		resp.Metadata.Message = "invalid json format"

		util.SetMetadataResponse(ctx, http.StatusBadRequest, startTime, &resp)
		return
	}

	resp = c.service.Create(param)
	statusCode := http.StatusCreated

	switch resp.Metadata.Message {
	case "SKU already exist":
		statusCode = http.StatusConflict
	case "failed to create product":
		statusCode = http.StatusInternalServerError
	}

	util.SetMetadataResponse(ctx, statusCode, startTime, &resp)
}

// @Tags product
// @Accept json
// @Produce json
// @Param search query string false "incase-sensitive"
// @Param page query int false "page"
// @Param limit query int false "limit"
// @Success 200 {object} object{metadata=model.MetadataResponse, data=model.PaginatedResponse} "OK"
// @Router /products [get]
func (c *ProductController) Search(ctx *gin.Context) {
	startTime := time.Now()
	var resp model.Response

	var param dto.ProductQueryRequest
	if err := ctx.ShouldBindQuery(&param); err != nil {
		log.Println(err)
		resp.Metadata.Message = "invalid queary parameter"
		util.SetMetadataResponse(ctx, http.StatusBadRequest, startTime, &resp)
		return
	}

	resp = c.service.Search(param)
	statusCode := http.StatusCreated

	if resp.Metadata.Message == "failed to count product" || resp.Metadata.Message == "failed to fetch product" {
		statusCode = http.StatusInternalServerError
	}

	util.SetMetadataResponse(ctx, statusCode, startTime, &resp)
}

// @Tags product
// @Produce json
// @Param id path int true "product id"
// @Success 204 "No Content"
// @Router /products/{id} [delete]
func (c *ProductController) Delete(ctx *gin.Context) {
	starTime := time.Now()
	var resp model.Response

	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		log.Println(err)
		resp.Metadata.Message = "invalid product id"
		util.SetMetadataResponse(ctx, http.StatusBadRequest, starTime, &resp)
		return
	}

	resp = c.service.Delete(id)
	statusCode := http.StatusInternalServerError

	switch resp.Metadata.Message {
	case "OK":
		ctx.Status(http.StatusNoContent)
		return
	case "product has movement hoistory":
		statusCode = http.StatusConflict
	}

	util.SetMetadataResponse(ctx, statusCode, starTime, &resp)
}

// @Tags product
// @Accept json
// @Produce json
// @Param id path int true "product id"
// @Param parameter body dto.ProductUpdateRequest true "parameter"
// @Success 201 {object} object{metadata=model.MetadataResponse} "OK"
// @Router /products/{id} [put]
func (c *ProductController) Update(ctx *gin.Context) {
	startTime := time.Now()
	var resp model.Response

	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		resp.Metadata.Message = "invalid prooduct id"
		util.SetMetadataResponse(ctx, http.StatusBadRequest, startTime, &resp)
		return
	}

	var param dto.ProductUpdateRequest
	if err := ctx.BindJSON(&param); err != nil {
		resp.Metadata.Message = "invalid json format"
		util.SetMetadataResponse(ctx, http.StatusBadRequest, startTime, &resp)
		return
	}

	resp = c.service.Update(id, param)
	statusCode := http.StatusOK

	switch {
	case resp.Metadata.Message == "product not found":
		statusCode = http.StatusNotFound
	case resp.Metadata.Message == "SKU already exists":
		statusCode = http.StatusConflict
	case resp.Metadata.Message != "OK":
		statusCode = http.StatusInternalServerError
	}

	util.SetMetadataResponse(ctx, statusCode, startTime, &resp)
}

// @Tags product
// @Accept json
// @Produce json
// @Param id path int true "product id"
// @Param parameter body dto.StockMovementRequest true "parameter"
// @Success 201 {object} object{metadata=model.MetadataResponse} "OK"
// @Router /products/{id}/stock [post]
func (c *ProductController) AdjustStock(ctx *gin.Context) {
	startTime := time.Now()
	var resp model.Response

	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		log.Println(err)
		resp.Metadata.Message = "invalid product id"
		util.SetMetadataResponse(ctx, http.StatusBadRequest, startTime, &resp)
		return
	}

	var param dto.StockMovementRequest
	if err := ctx.BindJSON(&param); err != nil {
		log.Println(err)
		resp.Metadata.Message = "invalid json format or validation"
		util.SetMetadataResponse(ctx, http.StatusBadRequest, startTime, &resp)
		return
	}

	resp = c.service.AdjustStock(id, param)
	statusCode := http.StatusOK

	switch {
	case resp.Metadata.Message == "product not found":
		statusCode = http.StatusNotFound
	case resp.Metadata.Message == "insufficient stock":
		statusCode = http.StatusConflict
	case resp.Metadata.Message != "OK":
		statusCode = http.StatusInternalServerError
	}

	util.SetMetadataResponse(ctx, statusCode, startTime, &resp)
}

// @Tags product
// @Produce json
// @Param id path int true "product id"
// @Success 201 ""
// @Router /products/{id}/movement [get]
func (c *ProductController) History(ctx *gin.Context) {
	startTime := time.Now()
	var resp model.Response

	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		log.Println(err)
		resp.Metadata.Message = "invalid productd  id"
		util.SetMetadataResponse(ctx, http.StatusBadRequest, startTime, &resp)
		return
	}

	resp = c.service.Movement(id)
	util.SetMetadataResponse(ctx, http.StatusOK, startTime, &resp)
}
