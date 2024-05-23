package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/rohhamh/go-shopping-cart-crud/database"
	"github.com/rohhamh/go-shopping-cart-crud/model"
	"github.com/rohhamh/go-shopping-cart-crud/utils/handlers"
)

type Cart struct {
	Prefix	string
}

type CartRequest struct {
    Data		string
    State       string
}
var CartRequestHandler CartRequest

func (c Cart) Handle(mux *http.ServeMux) {
	mux.HandleFunc(
		fmt.Sprintf("%s", c.Prefix),
		handlers.WithLogger(CartRequestHandler.Basket),
	)
	mux.HandleFunc(
		fmt.Sprintf("%s/{id}", c.Prefix),
		handlers.WithLogger(CartRequestHandler.Get),
	)
}

func (cr CartRequest) Basket (res http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		cr.GetAll(res, req)
	} else if req.Method == "POST" {
		cr.Create(res, req)
	}
}
func (cr CartRequest) GetAll (res http.ResponseWriter, req *http.Request) {

	db := database.Connection()
	carts := []model.Cart{}
	db.Find(&carts)

	sampleCart, err := json.Marshal(carts)
	if err != nil {
		fmt.Printf("Error marshalling carts %v\n", err)
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte("Bad carts"))
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(sampleCart)
}

func (cr CartRequest) Get (res http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := strings.Split(req.URL.Path, "/")[2]

	cart := model.Cart{}
	db := database.Connection()
	query := db.Find(&cart, id)
	if query.RowsAffected == 0 {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	sampleCart, err := json.Marshal(cart)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte("Bad cart"))
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(sampleCart)
}

func (cr CartRequest) Create (res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	cart := model.Cart{}
	err := json.NewDecoder(req.Body).Decode(&cart)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Printf("cart %+v\n", cart)

	db := database.Connection()
	query := db.Create(&cart)
	if query.Error != nil {
		fmt.Printf("err %+v", query.Error)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusCreated)
}
