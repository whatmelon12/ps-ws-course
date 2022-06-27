package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/whatmelon12/ps-ws-course/middleware/cors"
	"github.com/whatmelon12/ps-ws-course/model"
	"github.com/whatmelon12/ps-ws-course/services"
)

const path = "product"

func handleProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getProducts(w)
	case http.MethodPost:
		createProduct(w, r)
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleProduct(w http.ResponseWriter, r *http.Request) {
	pathSegments := strings.Split(r.URL.Path, fmt.Sprintf("%s/", path))
	if len(pathSegments[1:]) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	productId := pathSegments[len(pathSegments)-1]
	switch r.Method {
	case http.MethodGet:
		getProduct(productId, w)
	case http.MethodPut:
		updateProduct(productId, w, r)
	case http.MethodDelete:
		deleteProduct(productId, w)
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getProducts(w http.ResponseWriter) {
	products := services.GetProducts()
	data, err := json.Marshal(products)
	if err != nil {
		log.Fatal(err)
	}
	_, err = w.Write([]byte(data))
	if err != nil {
		log.Fatal(err)
	}
}

func getProduct(productId string, w http.ResponseWriter) {
	product := services.GetProduct(productId)
	if product == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	data, err := json.Marshal(product)
	if err != nil {
		log.Fatal(err)
	}
	_, err = w.Write([]byte(data))
	if err != nil {
		log.Fatal(err)
	}
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	var product model.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	insertId, err := services.CreateProduct(product)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"productId": %d}`, insertId)))
}

func updateProduct(productId string, w http.ResponseWriter, r *http.Request) {
	var product model.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = services.UpdateProduct(productId, product)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func deleteProduct(productId string, w http.ResponseWriter) {
	err := services.DeleteProduct(productId)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func SetupProductRoutes(apiBasePath string) {
	productsHandler := http.HandlerFunc(handleProducts)
	productHandler := http.HandlerFunc(handleProduct)

	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, path), cors.Middleware(productsHandler))
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, path), cors.Middleware(productHandler))
}
