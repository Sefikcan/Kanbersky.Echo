package main

import (
	"github.com/sirupsen/logrus"
	"kanberskyecho/pkg/handlers/product"
	. "kanberskyecho/pkg/infrastructure/product"
	"kanberskyecho/pkg/logging"
	. "kanberskyecho/pkg/services/product"
)

func main() {
	logger := logging.NewLoggerService(logrus.New())
	productRepository := NewProductRepository(Connect())
	productService := NewProductService(productRepository, logger)
	product.NewProductHandler(productService, logger)

	// TODO: Register Jaeger
	// TODO: Register Prometheous
	// TODO: Register Casbin
}
