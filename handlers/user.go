package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
    "crypto/subtle"

	"github.com/rohhamh/go-shopping-cart-crud/config"
	"github.com/rohhamh/go-shopping-cart-crud/database"
	"github.com/rohhamh/go-shopping-cart-crud/model"
	"github.com/rohhamh/go-shopping-cart-crud/utils/handlers"
	"golang.org/x/crypto/argon2"
	// "gopkg.in/validator.v2"
)

type User struct {
	Prefix	string
}

type UserRequest struct {
	Name		string
	Email		string
	Password	string
}
var UserRequestHandler UserRequest

func (u User) Handle(mux *http.ServeMux) {
	mux.HandleFunc(
		fmt.Sprintf("%s", u.Prefix),
		handlers.WithLogger(UserRequestHandler.Create),
	)
	mux.HandleFunc(
		fmt.Sprintf("%s/login", u.Prefix),
		handlers.WithLogger(UserRequestHandler.Login),
	)
}

func (urh UserRequest) Create (res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	user := UserRequest{}

	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	// if errs := validator.Validate(user); errs != nil {
	// 	fmt.Printf("bad user %+v\n", errs)
	// 	res.WriteHeader(http.StatusBadRequest)
	// 	res.Write([]byte(errs.Error()))
	// }
	db := database.Connection()

	salt := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		fmt.Print("Err generating random salt")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	passwordHash := argon2.IDKey(
		[]byte(user.Password), salt,
		config.PasswordTime, config.PasswordMemory, config.PasswordThreads, config.PasswordKeyLen,
	)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
    b64Hash := base64.RawStdEncoding.EncodeToString(passwordHash)

	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, config.PasswordMemory, config.PasswordTime, config.PasswordThreads,
		b64Salt, b64Hash,
	)
	dbUser := model.User{
		Name: user.Name,
		Email: user.Email,
		Password: encodedHash,
	}
    query := db.Where("email = ?", user.Email).Find(&dbUser)
    if query.RowsAffected > 0 {
		res.WriteHeader(http.StatusBadRequest)
        res.Write([]byte("User already exists!"))
        return
    }
	insertion := db.Create(&dbUser)
	if insertion.Error != nil {
		res.WriteHeader(http.StatusInternalServerError)
        return
	}
	res.WriteHeader(http.StatusCreated)
}

func (urh UserRequest) Login(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	user := UserRequest{}

	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Invalid body"))
		return
	}

    user.Email = strings.TrimSpace(user.Email)
    if user.Email == "" || user.Password == "" {
        res.WriteHeader(http.StatusBadRequest)
        return
    }

	dbUser := model.User{}
	db := database.Connection()
    query := db.Where("email = ?", user.Email).Find(&dbUser)
    if query.RowsAffected <= 0 {
        fmt.Printf("email %s not found\n", user.Email)
		res.WriteHeader(http.StatusForbidden)
        return
    }

	encodedHash := dbUser.Password

	vals := strings.Split(encodedHash, "$")
    if len(vals) != 6 {
		res.WriteHeader(http.StatusForbidden)
		return
    }

    var version int
    _, err = fmt.Sscanf(vals[2], "v=%d", &version)
    if err != nil {
		res.WriteHeader(http.StatusForbidden)
		return
    }
    if version != argon2.Version {
		res.WriteHeader(http.StatusForbidden)
		return
    }

	var time uint32
	var memory uint32
	var threads uint8
    _, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
    if err != nil {
        fmt.Printf("reading argon parameters failed for user %s, error: %v\n", user.Email, err)
		res.WriteHeader(http.StatusForbidden)
		return
    }

	salt, err := base64.RawStdEncoding.Strict().DecodeString(vals[4])
    if err != nil {
        fmt.Printf("salt b64 decode failed for user %s, error: %v\n", user.Email, err)
		res.WriteHeader(http.StatusForbidden)
		return
    }

	dbPassword, err := base64.RawStdEncoding.Strict().DecodeString(vals[5])
    if err != nil {
        fmt.Printf("password b64 decode failed for user %s, error: %v\n", user.Email, err)
		res.WriteHeader(http.StatusForbidden)
		return
    }
	keyLength := uint32(len(dbPassword))

    reqPassword := argon2.IDKey([]byte(user.Password), salt, time, memory, threads, keyLength)

    if subtle.ConstantTimeCompare(reqPassword, dbPassword) == 0 {
        fmt.Printf("bad password for user %s\n", user.Email)
        res.WriteHeader(http.StatusForbidden)
        return
    }

	res.WriteHeader(http.StatusOK)
}
