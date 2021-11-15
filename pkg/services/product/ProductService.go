package ProductService

import (
	"github.com/sirupsen/logrus"
	. "kanberskyecho/pkg/infrastructure/product"
	"kanberskyecho/pkg/infrastructure/product/entities"
	"kanberskyecho/pkg/logging"
)

type ProductService interface {
	GetProductById(entities.Product)(entities.Product, error)
	AddProduct(entities.Product) (entities.Product, error)
	RemoveProduct(entities.Product) error
	UpdateProduct(entities.Product) (entities.Product, error)
}

type productService struct {
	productRepository ProductRepository
	loggerService     logging.LoggerService
}

func (p productService) AddProduct(product entities.Product) (entities.Product, error) {
	resp, err := p.productRepository.InsertProduct(product)
	if err != nil {
		p.loggerService.CreateLog().WithFields(logrus.Fields {
			"methodName": "AddProduct",
			"type":"AddProduct_Service_Action",
		}).Error(err.Error())
		return entities.Product{}, err
	}

	return resp, nil
}

func (p productService) GetProductById(product entities.Product) (entities.Product, error) {
	resp, err := p.productRepository.GetProductById(product)
	if err != nil {
		p.loggerService.CreateLog().WithFields(logrus.Fields {
			"methodName": "GetProductById",
			"type":"GetProductById_Service_Action",
		}).Error(err.Error())
		return entities.Product{}, err
	}

	return resp, nil
}

func (p productService) RemoveProduct(product entities.Product) error {
	err := p.productRepository.DeleteProduct(product)
	if err != nil {
		p.loggerService.CreateLog().WithFields(logrus.Fields {
			"methodName": "RemoveProduct",
			"type":"RemoveProduct_Service_Action",
		}).Error(err.Error())
		return err
	}

	return nil
}

func (p productService) UpdateProduct(product entities.Product) (entities.Product, error) {
	resp, err := p.productRepository.UpdateProduct(product)
	if err != nil {
		p.loggerService.CreateLog().WithFields(logrus.Fields {
			"methodName": "UpdateProduct",
			"type":"UpdateProduct_Service_Action",
		}).Error(err.Error())
		return entities.Product{}, err
	}

	return resp, nil
}

func NewProductService(p ProductRepository, logger logging.LoggerService) ProductService{
	return &productService{p, logger}
}
