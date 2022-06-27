package services

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/whatmelon12/ps-ws-course/database"
	"github.com/whatmelon12/ps-ws-course/model"
)

func GetProducts() []model.Product {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, "SELECT productId, productName, sku, pricePerUnit FROM products")
	if err != nil {
		log.Fatal(err)
		return nil
	}

	defer results.Close()
	products := make([]model.Product, 0)
	for results.Next() {
		product := &model.Product{}
		results.Scan(&product.ProductID, &product.ProductName, &product.Sku, &product.PricePerUnit)
		products = append(products, *product)
	}
	return products
}

func GetProduct(productId string) *model.Product {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	results := database.DbConn.QueryRowContext(ctx, `SELECT productId, productName, sku, pricePerUnit 
	FROM products
	WHERE productId = ?`, productId)
	product := &model.Product{}
	err := results.Scan(&product.ProductID, &product.ProductName, &product.Sku, &product.PricePerUnit)
	if err == sql.ErrNoRows {
		return nil
	}
	return product
}

func CreateProduct(product model.Product) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	result, err := database.DbConn.ExecContext(ctx, `INSERT INTO products
	(manufacturer, 
		sku,
		upc,
		pricePerUnit,
		quantityOnHand,
		productName) VALUES (?, ?, ?, ?, ?, ?)`,
		product.Manufacturer,
		product.Sku,
		product.Upc,
		product.PricePerUnit,
		product.QuantityOnHand,
		product.ProductName)
	if err != nil {
		return 0, nil
	}

	insertId, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}
	return int(insertId), nil
}

func UpdateProduct(productId string, product model.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_, err := database.DbConn.ExecContext(ctx, `UPDATE products SET 
		manufacturer=?, 
		sku=?,
		upc=?,
		pricePerUnit=?,
		quantityOnHand=?,
		productName=?
		WHERE productId=?`,
		product.Manufacturer,
		product.Sku,
		product.Upc,
		product.PricePerUnit,
		product.QuantityOnHand,
		product.ProductName,
		productId)
	if err != nil {
		return err
	}
	return nil
}

func DeleteProduct(productId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_, err := database.DbConn.ExecContext(ctx, "DELETE FROM products WHERE productId=?", productId)
	if err != nil {
		return err
	}
	return nil
}
