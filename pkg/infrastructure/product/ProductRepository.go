package ProductRepository

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"kanberskyecho/pkg/infrastructure/product/entities"
)

type ProductRepository interface {
	GetProductById(entities.Product) (entities.Product, error)
	InsertProduct(entities.Product) (entities.Product, error)
	DeleteProduct(entities.Product) error
	UpdateProduct(entities.Product) (entities.Product, error)
}

type productRepository struct {
	db *gorm.DB
}

func (p productRepository) InsertProduct(product entities.Product) (entities.Product, error) {
	if result := p.db.Create(&product); result.Error != nil {
		return product, result.Error
	}

	return product, nil
}

func (p productRepository) GetProductById(product entities.Product) (entities.Product, error) {
	result := p.db.First(&product, &product.Id)
	if result.Error != nil {
		return entities.Product{}, result.Error
	}

	return product, nil
}

func (p productRepository) DeleteProduct(product entities.Product) error {
	if result := p.db.Delete(product); result.Error != nil {
		return result.Error
	}

	return nil
}

func (p productRepository) UpdateProduct(product entities.Product) (entities.Product, error) {
	if result := p.db.Save(&product); result.Error != nil {
		return product, result.Error
	}

	return product, nil
}

func NewProductRepository(db *gorm.DB) ProductRepository{
	return &productRepository{db: db}
}

func Connect() *gorm.DB{
	db, err := gorm.Open(mysql.Open("admin:password@tcp(127.0.0.1:3306)/db?parseTime=true"))
	if err != nil {
		panic("Could not connect to the database")
	}

	db.AutoMigrate(&entities.Product{})
	return db
}
