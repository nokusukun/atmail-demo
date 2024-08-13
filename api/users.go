package api

import (
	"atmail-demo/config"
	"atmail-demo/database"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nokusukun/faust"
	"github.com/nokusukun/faust/param"
	"log"
	"net/http"
	"strings"
	"time"
)

type UserPayload struct {
	Username    string `json:"username,omitempty"`
	Email       string `json:"email,omitempty"`
	Age         int    `json:"age,omitempty"`
	Permissions string `json:"permissions,omitempty"`
}

func UserPayloadValidator(u UserPayload) error {
	if u.Username == "" {
		return fmt.Errorf("Username is required")
	}
	if u.Email == "" {
		return fmt.Errorf("Email is required")
	}
	if u.Age < 1 {
		return fmt.Errorf("Age must be 1 or older")
	}
	if !strings.Contains(u.Email, "@") {
		return fmt.Errorf("Email is invalid")
	}
	return nil
}

func ReturnJSON(w http.ResponseWriter, v any, code ...int) {
	w.Header().Set("Content-Type", "application/json")
	if len(code) > 0 {
		w.WriteHeader(code[0])
	} else {
		w.WriteHeader(200)
	}
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		log.Println("Error returning JSON response:", err)
	}
}

func ReturnError(w http.ResponseWriter, err error, code int) {
	ReturnJSON(w, map[string]any{
		"error": err.Error(),
	}, code)
}

var LoggingMiddleware mux.MiddlewareFunc = func(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		handler.ServeHTTP(w, r)
		log.Printf("[%v] %v %v\n", r.Method, r.URL.Path, time.Since(now))
	})
}

var AuthMiddleware mux.MiddlewareFunc = func(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("AuthMiddleware")
		auth := r.Header.Get("Authorization")
		if auth == "" {
			ReturnJSON(w, map[string]any{
				"error": "Authorization header is required",
			}, 401)
			return
		}
		if !strings.HasPrefix(auth, "Basic ") {
			ReturnJSON(w, map[string]any{
				"error": "Authorization header must be a Basic auth",
			}, 401)
			return
		}
		if auth != fmt.Sprintf("Basic %v:%v", config.USERNAME, config.PASSWORD) {
			ReturnJSON(w, map[string]any{
				"error": "Invalid credentials",
			}, 401)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func UsersEndpoint(server *faust.API, db *database.Database) {
	router := server.Subrouter("/users")

	router.Get("/{id}", func(e *faust.Endpoint) http.HandlerFunc {
		e.Description("Fetch a user by ID")
		e.Middlewares(LoggingMiddleware, AuthMiddleware)
		id := param.Path[uint](e, "id", param.Info{
			Description: "The ID of the user to fetch",
		})

		return func(w http.ResponseWriter, r *http.Request) {
			user, err := db.GetUser(id.Value(r))
			if err != nil {
				code := 500
				if strings.Contains(err.Error(), "record not found") {
					code = 404
				}
				ReturnError(w, err, code)
				return
			}
			ReturnJSON(w, user)
		}
	})

	router.Post("/", func(e *faust.Endpoint) http.HandlerFunc {
		e.Description("Create a new user")
		e.Middlewares(LoggingMiddleware, AuthMiddleware)
		payload := param.Json[UserPayload](e, "payload", param.Info{
			Description: "The payload to create a new user",
		}).Validate(UserPayloadValidator)

		return func(w http.ResponseWriter, r *http.Request) {
			newUser, err := payload.ValueWithError(r)
			if err != nil {
				fmt.Println("Error validating payload", err)
				ReturnError(w, err, 400)
				return
			}

			if newUser.Permissions == "" {
				newUser.Permissions = "PUT,DELETE"
			}

			u, err := db.NewUser(&database.User{
				Username:    newUser.Username,
				Email:       newUser.Email,
				Age:         newUser.Age,
				Permissions: newUser.Permissions,
			})
			if err != nil {
				ReturnError(w, err, 500)
				return
			}

			ReturnJSON(w, map[string]any{
				"ok":   true,
				"user": u,
			})
		}
	})

	router.Put("/{id}", func(e *faust.Endpoint) http.HandlerFunc {
		e.Description("Update a user by ID")
		e.Middlewares(LoggingMiddleware, AuthMiddleware)
		id := param.Path[uint](e, "id", param.Info{
			Description: "The ID of the user to update",
		})

		payload := param.Json[UserPayload](e, "payload", param.Info{
			Description: "The payload to update a user",
		})

		return func(w http.ResponseWriter, r *http.Request) {
			updateUser, err := payload.ValueWithError(r)
			if err != nil {
				fmt.Println("Error validating payload", err)
				ReturnError(w, err, 400)
				return
			}
			oldUser, err := db.GetUser(id.Value(r))
			if err != nil {
				if strings.Contains(err.Error(), "record not found") {
					ReturnError(w, fmt.Errorf("User not found"), 404)
					return
				}
				ReturnError(w, err, 500)
				return
			}

			if !strings.Contains(oldUser.Permissions, "PUT") {
				ReturnError(w, fmt.Errorf("You cannot update this user"), 403)
				return
			}

			err = db.UpdateUser(id.Value(r), &database.User{
				Username:    updateUser.Username,
				Email:       updateUser.Email,
				Age:         updateUser.Age,
				Permissions: updateUser.Permissions,
			})
			if err != nil {
				ReturnError(w, err, 500)
				return
			}

			ReturnJSON(w, map[string]any{
				"ok": true,
			})
		}
	})

	router.Delete("/{id}", func(e *faust.Endpoint) http.HandlerFunc {
		e.Description("Delete a user by ID")
		e.Middlewares(LoggingMiddleware, AuthMiddleware)
		id := param.Path[uint](e, "id", param.Info{
			Description: "The ID of the user to delete",
		})

		return func(w http.ResponseWriter, r *http.Request) {
			oldUser, err := db.GetUser(id.Value(r))
			if err != nil {
				if strings.Contains(err.Error(), "record not found") {
					ReturnError(w, fmt.Errorf("User not found"), 404)
					return
				}
				ReturnError(w, err, 500)
				return
			}

			if !strings.Contains(oldUser.Permissions, "DELETE") {
				ReturnError(w, fmt.Errorf("You cannot delete this user"), 403)
				return
			}

			err = db.DeleteUser(id.Value(r))
			if err != nil {
				ReturnError(w, err, 500)
				return
			}

			ReturnJSON(w, map[string]any{
				"ok": true,
			})
		}
	})
}
