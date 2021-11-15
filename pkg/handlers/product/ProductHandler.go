package product

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"kanberskyecho/pkg/infrastructure/product/entities"
	"kanberskyecho/pkg/logging"
	. "kanberskyecho/pkg/services/product"
	"net/http"
	"strconv"
)

func NewProductHandler(productService ProductService, logger logging.LoggerService){
	e := echo.New()

	handler := &productHandler{productService, logger}

	e.POST("/api/v1/products", handler.AddProduct)
	e.PUT("/api/v1/products/:id", handler.UpdateProduct)
	e.DELETE("/api/v1/products/:id", handler.DeleteProduct)

	logger.CreateLog().Error(e.Start(":7500"))
}

type productHandler struct {
	p       ProductService
	logger  logging.LoggerService
}

func (p *productHandler) AddProduct(c echo.Context) error {
	var product entities.Product
	if err:= c.Bind(&product); err != nil {
		p.logger.CreateLog().WithFields(logrus.Fields{
			"methodName": "AddProduct",
			"type":"AddProduct_Handler_Bind_Operation",
		}).Error(err.Error())
		return c.JSON(http.StatusBadRequest,echo.Map{
			"error": err.Error(),
		})
	}

	resp, err := p.p.AddProduct(product)
	if err != nil {
		p.logger.CreateLog().WithFields(logrus.Fields{
			"methodName": "AddProduct",
			"type":"AddProduct_Handler_Insert_Operation",
		}).Error(err.Error())
		return c.JSON(http.StatusBadRequest,echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"data": resp,
	})
}

func (p *productHandler) UpdateProduct(c echo.Context) error {
	var product entities.Product

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		p.logger.CreateLog().WithFields(logrus.Fields{
			"methodName": "UpdateProduct",
			"type":"UpdateProduct_Handler_Convert_Id_Operation",
		}).Error(err.Error())
		return c.JSON(http.StatusBadRequest,echo.Map{
			"error": err.Error(),
		})
	}

	if err:= c.Bind(&product); err != nil {
		p.logger.CreateLog().WithFields(logrus.Fields{
			"methodName": "UpdateProduct",
			"type":"UpdateProduct_Handler_Bind_Operation",
		}).Error(err.Error())
		return c.JSON(http.StatusBadRequest,echo.Map{
			"error": err.Error(),
		})
	}

	product.Id = id

	_, err = p.p.GetProductById(product)
	if err != nil {
		p.logger.CreateLog().WithFields(logrus.Fields{
			"methodName": "UpdateProduct",
			"type":"UpdateProduct_Handler_IsExists_Operation",
		}).Error(err.Error())
		return c.JSON(http.StatusNotFound,echo.Map{
			"error": err.Error(),
		})
	}

	resp, err := p.p.UpdateProduct(product)
	if err != nil {
		p.logger.CreateLog().WithFields(logrus.Fields{
			"methodName": "UpdateProduct",
			"type":"UpdateProduct_Handler_Update_Operation",
		}).Error(err.Error())
		return c.JSON(http.StatusBadRequest,echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": resp,
	})
}

func (p *productHandler) DeleteProduct(c echo.Context) error {
	var product entities.Product

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		p.logger.CreateLog().WithFields(logrus.Fields{
			"methodName": "DeleteProduct",
			"type":"DeleteProduct_Handler_Bind_Id_Operation",
		}).Error(err.Error())
		return c.JSON(http.StatusBadRequest,echo.Map{
			"error": err.Error(),
		})
	}

	product.Id = id

	_, err = p.p.GetProductById(product)
	if err != nil {
		p.logger.CreateLog().WithFields(logrus.Fields{
			"methodName": "DeleteProduct",
			"type":"DeleteProduct_Handler_IsExists_Operation",
		}).Error(err.Error())
		return c.JSON(http.StatusNotFound,echo.Map{
			"error": err.Error(),
		})
	}

	err = p.p.RemoveProduct(product)
	if err != nil {
		p.logger.CreateLog().WithFields(logrus.Fields{
			"methodName": "DeleteProduct",
			"type":"DeleteProduct_Handler_Delete_Operation",
		}).Error(err.Error())
		return c.JSON(http.StatusBadRequest,echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, "")
}
