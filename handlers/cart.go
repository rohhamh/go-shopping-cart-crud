package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/rohhamh/go-shopping-cart-crud/database"
	"github.com/rohhamh/go-shopping-cart-crud/middlewares"
	"github.com/rohhamh/go-shopping-cart-crud/model"
	middlewareUtils "github.com/rohhamh/go-shopping-cart-crud/utils/middlewares"
	"gorm.io/datatypes"
)

type Cart struct {
	Prefix	        string
    Middlewares     *[]middlewares.Middleware
}

type CartRequest struct {
    Data		datatypes.JSON
    State       string
}
var CartRequestHandler CartRequest

func (c Cart) Handle(mux *http.ServeMux) {
    var basketRoutes    middlewares.RequestHandler = CartRequestHandler.Basket
    var basketIdRoutes  middlewares.RequestHandler = CartRequestHandler.BasketID

    mux.HandleFunc(
        fmt.Sprintf("%s", c.Prefix),
        middlewareUtils.Chain(c.Middlewares, &basketRoutes),
    )
    mux.HandleFunc(
        fmt.Sprintf("%s/{id}", c.Prefix),
        middlewareUtils.Chain(c.Middlewares, &basketIdRoutes),
    )
}

func (cr CartRequest) Basket (res http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		cr.GetAll(res, req)
	} else if req.Method == "POST" {
		cr.Create(res, req)
	} else {
        res.WriteHeader(http.StatusMethodNotAllowed)
        return
    }
}

func (cr CartRequest) BasketID (res http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		cr.Get(res, req)
	} else if req.Method == "PATCH" {
		cr.Update(res, req)
    } else if req.Method == "DELETE" {
        cr.Delete(res, req)
    } else {
        res.WriteHeader(http.StatusMethodNotAllowed)
        return
    }
}

func (cr CartRequest) GetAll (res http.ResponseWriter, req *http.Request) {
	db := database.Connection()
	carts := []model.Cart{}
    user := req.Context().Value("user").(model.User)

	db.Where("user_id = ?", user.ID).Find(&carts)

    validCarts := []model.Cart {}
    for _, cart := range(carts) {
        if cart.State != "DELETED" {
            validCarts = append(validCarts, cart)
        }
    }

	sampleCart, err := json.Marshal(validCarts)
	if err != nil {
		fmt.Printf("Error marshalling carts %v\n", err)
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte("Bad carts"))
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(sampleCart)
}

func (cr CartRequest) Get (res http.ResponseWriter, req *http.Request) {
	cart := model.Cart{}
    cart.User = req.Context().Value("user").(model.User)
    idStr := strings.Split(req.URL.Path, "/")[2]
	id, err := strconv.Atoi(idStr)
    if err != nil {
        fmt.Printf("got invalid id %s from user %d\n", idStr, cart.User.ID)
        res.WriteHeader(http.StatusBadRequest)
        return
    }
    cart.ID = int64(id)


	db := database.Connection()
	query := db.Where("user_id = ?", cart.User.ID).Find(&cart)
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
	cartRequest := CartRequest{}
	err := json.NewDecoder(req.Body).Decode(&cartRequest)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

    if cartRequest.State != "COMPLETED" && cartRequest.State != "PENDING" {
        res.WriteHeader(http.StatusBadRequest)
        res.Write([]byte("invalid basket state"))
        return
    }

    cart := model.Cart {
        Data:   cartRequest.Data,
        State:  cartRequest.State,
    }
    cart.User = req.Context().Value("user").(model.User)

	db := database.Connection()
	query := db.Create(&cart)
	if query.Error != nil {
		fmt.Printf("err %+v\n", query.Error)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

    cartJson, err := json.Marshal(cart)
    if err != nil {
        fmt.Printf("error marshalling created cart %v\n", err)
        return
    }
    res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
    res.Write(cartJson)
}

func (cr CartRequest) Update (res http.ResponseWriter, req *http.Request) {
	cartRequest := CartRequest{}
	err := json.NewDecoder(req.Body).Decode(&cartRequest)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("invalid body"))
		return
	}

    cart := model.Cart {
        Data:   cartRequest.Data,
        State:  cartRequest.State,
    }
    cart.User = req.Context().Value("user").(model.User)

    idStr := strings.Split(req.URL.Path, "/")[2]
	id, err := strconv.Atoi(idStr)
    if err != nil {
        fmt.Printf("got invalid id %s from user %d\n", idStr, cart.User.ID)
        res.WriteHeader(http.StatusBadRequest)
        return
    }
    cart.ID = int64(id)

	db := database.Connection()
	query := db.Where("user_id = ?", cart.User.ID).Find(&cart)
	if query.RowsAffected <= 0 || query.Error != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("invalid basket id"))
		return
	}

    if strings.ToUpper(cart.State) == "DELETED" {
        res.WriteHeader(http.StatusNotFound)
        return
    }

    if strings.ToUpper(cart.State) != "PENDING" {
        res.WriteHeader(http.StatusBadRequest)
        res.Write([]byte("invalid basket state"))
        return
    }

    if cartRequest.State != "COMPLETED" && cartRequest.State != "PENDING" {
        res.WriteHeader(http.StatusBadRequest)
        res.Write([]byte("invalid basket state"))
        return
    }

    cart.Data = cartRequest.Data
    cart.State = cartRequest.State
    db.Save(&cart)

	res.WriteHeader(http.StatusOK)
}

func (cr CartRequest) Delete (res http.ResponseWriter, req *http.Request) {
    cart := model.Cart {}
    cart.User = req.Context().Value("user").(model.User)

    idStr := strings.Split(req.URL.Path, "/")[2]
	id, err := strconv.Atoi(idStr)
    if err != nil {
        fmt.Printf("got invalid id %s from user %d\n", idStr, cart.User.ID)
        res.WriteHeader(http.StatusBadRequest)
        return
    }
    cart.ID = int64(id)

	db := database.Connection()
	query := db.Where("user_id = ?", cart.User.ID).Find(&cart)
	if query.RowsAffected <= 0 || query.Error != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("invalid basket id"))
		return
	}

    if strings.ToUpper(cart.State) == "DELETED" {
        res.WriteHeader(http.StatusNotFound)
        return
    }

    if strings.ToUpper(cart.State) != "PENDING" {
        res.WriteHeader(http.StatusBadRequest)
        res.Write([]byte("invalid basket state"))
        return
    }

    cart.State = "DELETED"
    db.Save(&cart)

	res.WriteHeader(http.StatusOK)
}
