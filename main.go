package main

import (
	"fmt"
	"net/http"

	"github.com/rohhamh/go-shopping-cart-crud/database"
	"github.com/rohhamh/go-shopping-cart-crud/handlers"
)

func main()  {
	database.Connect()
	mux := http.NewServeMux()

	handlers.Cart { Prefix: "/basket" }.Handle(mux)
	handlers.User { Prefix: "/user" }.Handle(mux)

	fmt.Println("Listening on port 3000...")
	http.ListenAndServe("localhost:3000", mux)
}
