// controllers/user_controller.go
package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go-auth/config"
	"go-auth/models"
	"go-auth/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// REGISTER
func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body: " + err.Error(),
		})
		return
	}

	// Validation
	if user.Username == "" || user.Email == "" || user.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Please enter a valid username, email, and password",
		})
		return
	}

	collection := config.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if user already exists
	var existingUser models.User
	err := collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "User with this email already exists",
		})
		return
	} else if err != mongo.ErrNoDocuments {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Database error: " + err.Error(),
		})
		return
	}

	// Hash password
	hashed, err := utils.HashPassword(user.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Error hashing password: " + err.Error(),
		})
		return
	}
	user.Password = hashed

	// Create MongoDB ObjectID
	user.ID = primitive.NewObjectID()

	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Error creating user: " + err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User registered successfully",
		"user": map[string]interface{}{
			"id":       user.ID.Hex(),
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

func Login(W http.ResponseWriter, r *http.Request) {
	W.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(W, "Method is not allowed ", http.StatusMethodNotAllowed)
		return
	}

	// login request structure

	var loginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		W.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(W).Encode(map[string]string{
			"error": "Invalid request body: " + err.Error(),
		})
		return

	}
	// validation

	if loginReq.Email == "" || loginReq.Password == "" {
		W.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(W).Encode(map[string]string{
			"error": "Please enter a valid email and password",
		})
		return

	}

	collection := config.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// find the user by email

	var user models.User

	err := collection.FindOne(ctx, bson.M{"email": loginReq.Email}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			W.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(W).Encode(map[string]string{
				"error": "Invalid email or password",
			})
			return
		}
		W.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(W).Encode(map[string]string{
			"error": "Database error: " + err.Error(),
		})
		return
	}

	// checked password
	err = utils.CheckPassword(user.Password, loginReq.Password)

	if err != nil {
		W.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(W).Encode(map[string]string{
			"error": "Invalid email or password",
		})
		return
	}
	token, err := utils.GenerateToken(user.ID.Hex(), user.Email, user.Username)
	if err != nil {
		W.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(W).Encode(map[string]string{
			"error": "Error generating token: " + err.Error(),
		})
		return
	}
	W.WriteHeader(http.StatusOK)
	json.NewEncoder(W).Encode(map[string]interface{}{
		"message": "Login successful",
		"token":   token,
	})
}
