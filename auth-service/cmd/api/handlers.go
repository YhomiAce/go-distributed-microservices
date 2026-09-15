package main

import (
	"authentication/data"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

var loggerBaseUrl = "http://logger-service:8082"

func (app *Application) authenticate(w http.ResponseWriter, r *http.Request) {
	payload := JsonResponse{
		Error: false,
		Message: "Authentication service is up and running",
	}
	err := app.writeJson(w, http.StatusOK, payload)
	if err != nil {
		http.Error(w, "Internal Server Error writing JSON response", http.StatusInternalServerError)
		// app.errorJson(w, err)
		return
	}
}

func (app *Application) registerUser(w http.ResponseWriter, r *http.Request) {
	type requestPayload struct {
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Password  string `json:"password"`
	}

	var payload requestPayload

	err := app.readJson(w, r, &payload)
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	existingUser, err := app.Models.User.GetByEmail(payload.Email)
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}

	if existingUser != nil {
		app.errorJson(w, fmt.Errorf("user already exists"), http.StatusBadRequest)
		return
	}

	user := data.User{
		Email:     payload.Email,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Active:    true,
	}

	hashedPassword, err := app.Models.User.HashPassword(payload.Password)
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}
	user.Password = hashedPassword

	userId, err := app.Models.User.Insert(user)
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}

	err = app.logRequest("register", fmt.Sprintf("%s created an account", user.Email))
	if err != nil {
		app.errorJson(w, err)
		return
	}

	payloadResponse := JsonResponse{
		Error:   false,
		Message: "User registered successfully",
		Data:    map[string]int{"user_id": userId},
	}
	err = app.writeJson(w, http.StatusCreated, payloadResponse)
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}
}

func (app *Application) loginUser(w http.ResponseWriter, r *http.Request) {
	type requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var payload requestPayload

	err := app.readJson(w, r, &payload)
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}
	
	user, err := app.Models.User.GetByEmail(payload.Email)
	if err != nil {
		app.errorJson(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}
	log.Println("User found:", user.Email, "ID:", user.ID)

	validPassword, err := user.PasswordMatches(payload.Password)
	log.Println("Password match result:", validPassword, "Error:", err)
	if err != nil || !validPassword {
		app.errorJson(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}
	log.Println("User authenticated:", user.Email, "ID:", user.ID)
	err = app.logRequest("authenticate", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		app.errorJson(w, err)
		return
	}

	payloadResponse := JsonResponse{
		Error:   false,
		Message: "User authenticated successfully",
		Data:    user,
	}
	err = app.writeJson(w, http.StatusOK, payloadResponse)
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}
}

func (app *Application) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}
	entry.Data = data
	entry.Name = name

	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	url := fmt.Sprintf("%s/%s", loggerBaseUrl, "log")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return  err
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return  err
	}
	defer res.Body.Close()
	return nil
}