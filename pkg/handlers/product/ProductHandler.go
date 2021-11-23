package product

import (
	"github.com/labstack/echo-contrib/jaegertracing"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pangpanglabs/echoswagger/v2"
	"github.com/sirupsen/logrus"
	"kanberskyecho/pkg/infrastructure/product/entities"
	"kanberskyecho/pkg/logging"
	. "kanberskyecho/pkg/services/product"
	"net/http"
	"strconv"
)

func NewProductHandler(productService ProductService, logger logging.LoggerService){
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowHeaders:[]string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// Enable tracing middleware
	c := jaegertracing.New(e, nil)
	defer c.Close()

	se := echoswagger.New(e, "doc/", &echoswagger.Info{
		Title:          "Product Api",
		Description:    "This is a sample product crud operation",
		Version:        "1.0.0",
		TermsOfService: "http://swagger.io/terms/",
		Contact: &echoswagger.Contact{
			Name: "Şefik Can Kanber",
		},
		License: &echoswagger.License{
			Name: "Apache 2.0",
			URL:  "http://www.apache.org/licenses/LICENSE-2.0.html",
		},
	})

	se.SetExternalDocs("Find out more about Swagger", "http://swagger.io").
		SetResponseContentType("application/xml", "application/json").
		SetUI(echoswagger.UISetting{DetachSpec: true, HideTop: true})
		//SetScheme("https", "http")

	handler := &productHandler{productService, logger}

	p := se.Group("product","/api/v1/products")

	p.POST("", handler.AddProduct).
		AddParamBody(entities.Product{}, "body", "Product object that needs to be added to the database", true).
		AddResponse(http.StatusCreated, "success", nil, nil).
		AddResponse(http.StatusBadRequest, "invalid input", nil, nil).
		AddResponse(http.StatusInternalServerError, "Db operation failed", nil, nil).
		SetRequestContentType("application/json", "application/xml").
		SetSummary("Add a new product to the database")

	p.PUT("/:id", handler.UpdateProduct).
		AddParamBody(entities.Product{},"body", "Product object that needs to be update to the store", true).
		AddParamPath(int64(0), "id", "Product id to update").
		AddResponse(http.StatusBadRequest, "Invalid ID supplied", nil, nil).
		AddResponse(http.StatusNotFound, "Product not found", nil, nil).
		AddResponse(http.StatusOK, "Success", nil, nil).
		AddResponse(http.StatusInternalServerError, "Db operation failed", nil, nil).
		SetSummary("Update a product")

	p.DELETE("/:id", handler.DeleteProduct).
		AddParamPath(int64(0), "id", "Product id to delete").
		AddResponse(http.StatusBadRequest, "Invalid ID supplied", nil, nil).
		AddResponse(http.StatusNotFound, "Product not found", nil, nil).
		AddResponse(http.StatusNoContent, "Success", nil, nil).
		AddResponse(http.StatusInternalServerError, "Db operation failed", nil, nil).
		SetSummary("Delete a product")

	logger.CreateLog().Error(se.Echo().Start(":7500"))
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
		return c.JSON(http.StatusInternalServerError,echo.Map{
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
		return c.JSON(http.StatusInternalServerError,echo.Map{
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
		return c.JSON(http.StatusInternalServerError,echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, "")
}
