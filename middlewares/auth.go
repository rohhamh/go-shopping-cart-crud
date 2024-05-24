package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rohhamh/go-shopping-cart-crud/utils/jwt"
	"github.com/rohhamh/go-shopping-cart-crud/model"
	"github.com/rohhamh/go-shopping-cart-crud/database"
)

func Authorize(next *RequestHandler) RequestHandler {
	return func (res http.ResponseWriter, req *http.Request)  {
        jwtCookie, err := req.Cookie("jwt")
        if err != nil {
            res.WriteHeader(http.StatusForbidden)
            return
        }

        token, err := jwt.ValidateToken(jwtCookie.Value)
        if err != nil || !token.Valid {
            fmt.Printf("invalid token %v\n", err)
            res.WriteHeader(http.StatusForbidden)
            return
        }
        userEmail, err := token.Claims.GetSubject()
        if err != nil {
            fmt.Printf("Get token subject err %v\n", err)
            res.WriteHeader(http.StatusForbidden)
            return
        }

        db := database.Connection()
        user := model.User {}
        query := db.Where("email = ?", userEmail).Find(&user)
        if query.RowsAffected <= 0 {
            res.WriteHeader(http.StatusForbidden)
            return
        }
        ctx := context.WithValue(req.Context(), "user", user)
        if next != nil {
            (*next)(res, req.WithContext(ctx))
        }
	}
}
