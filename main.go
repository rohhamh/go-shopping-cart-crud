package main

import (
	"fmt"
	"net/http"

	"github.com/rohhamh/go-shopping-cart-crud/database"
	"github.com/rohhamh/go-shopping-cart-crud/handlers"
	"github.com/rohhamh/go-shopping-cart-crud/middlewares"
)

func main()  {
	database.Connect()
	mux := http.NewServeMux()

	handlers.Cart {
        Prefix: "/basket",
        Middlewares: &[]middlewares.Middleware{
            middlewares.Logger,
            middlewares.Authorize,
        }}.Handle(mux)

	handlers.User {
        Prefix: "/user",
        Middlewares: &[]middlewares.Middleware{
            middlewares.Logger,
        },
    }.Handle(mux)

	fmt.Println("Listening on port 3000...")
	http.ListenAndServe("localhost:3000", mux)
}
