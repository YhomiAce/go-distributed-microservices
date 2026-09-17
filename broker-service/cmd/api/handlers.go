package main

import (
	"broker/lib/events"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/rpc"
	"time"
)

var authserviceBaseUrl = "http://auth-service:8081"
var loggerBaseUrl = "http://logger-service:8082"
var mailerBaseUrl = "http://mailer-service:8083"

type RequestPayload struct {
	Action string 				`json:"action"`
	Auth AuthPayload			`json:"auth,omitempty"`
	Register RegisterPayload	`json:"register,omitempty"`
	Log LogPayload				`json:"log,omitempty"`
	Mail MailPayload			`json:"mail,omitempty"`
}

type MailPayload struct {
	From	string 	`json:"from"`
	To		string 		`json:"to"`
	Subject	string 	`json:"subject"`
	Message	string 	`json:"message"`
}

type AuthPayload struct {
	Email	string 	`json:"email"`
	Password	string `json:"password"`
}

type RegisterPayload struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

type LogPayload struct {
	Name	string 		`json:"name"`
	Data	string 		`json:"data"`
}

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:  100,
		MaxIdleConnsPerHost:  10,
		IdleConnTimeout: 90 * time.Second,
	},
}

func (app *Application) broker(w http.ResponseWriter, r *http.Request) {
	payload := JsonResponse{
		Error: false,
		Message: "Broker service is up and running",
	}
	err := app.writeJson(w, http.StatusOK, payload)
	if err != nil {
		http.Error(w, "Internal Server Error writing JSON response", http.StatusInternalServerError)
		// app.errorJson(w, err)
		return
	}
}

func (app *Application) handleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	err := app.readJson(w, r, &requestPayload)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	switch requestPayload.Action {
		case "register":
			app.register(w, requestPayload.Register)

		case "auth":
			app.authenticate(w, requestPayload.Auth)

		case "log":
			app.log(w, requestPayload.Log)

		case "log-mq":
			app.logViaRabbitMQ(w, requestPayload.Log)

		case "log-rpc":
			app.logViaRPC(w, requestPayload.Log)

		case "mail":
			app.sendMail(w, requestPayload.Mail)

		default: app.errorJson(w, errors.New("Invalid Action"))
	}
}

func (app *Application) authenticate(w http.ResponseWriter, payload AuthPayload) {
	jsonData, _ := json.MarshalIndent(payload, "", "\t")
	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", authserviceBaseUrl, "login"), bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJson(w, err)
		return
	}

	res, err := httpClient.Do(request)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusUnauthorized {
		app.errorJson(w, errors.New("Invalid Credentials"))
		return
	} else if res.StatusCode != http.StatusOK {
		app.errorJson(w, errors.New("Error calling auth service"))
		return
	}

	var response JsonResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	if response.Error {
		app.errorJson(w, errors.New(response.Message), http.StatusBadRequest)
		return
	}

	var responsePayload JsonResponse
	responsePayload.Error = false
	responsePayload.Message = "Success"
	responsePayload.Data = response.Data
	err = app.writeJson(w, http.StatusOK, responsePayload)
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}

}

func (app *Application) register(w http.ResponseWriter, payload RegisterPayload) {
	jsonData, _ := json.Marshal(payload)
	 res, err := http.Post(
		fmt.Sprintf("%s/%s", authserviceBaseUrl, "register"), 
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	defer res.Body.Close()

	var response JsonResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	if response.Error {
		app.errorJson(w, errors.New(response.Message), http.StatusBadRequest)
		return
	}
	var responsePayload JsonResponse
	responsePayload.Error = false
	responsePayload.Message = "Success"
	responsePayload.Data = response.Data
	err = app.writeJson(w, http.StatusOK, responsePayload)
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}
}

func (app *Application) log(w http.ResponseWriter, payload LogPayload) {
	jsonData, _ := json.Marshal(payload)
	var url = fmt.Sprintf("%s/log",loggerBaseUrl)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJson(w, err)
		return
	}
	res, err := httpClient.Do(req)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusAccepted {
		log.Println(res.StatusCode)
		app.errorJson(w, errors.New("error logging message"), http.StatusInternalServerError)
		return
	}
	var response JsonResponse
	response.Error = false
	response.Message = "Logged"
	err = app.writeJson(w, http.StatusAccepted, response)
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}
}

func (app *Application) sendMail(w http.ResponseWriter, payload MailPayload) {
	jsonData, _ := json.Marshal(payload)

	var url = fmt.Sprintf("%s/send", mailerBaseUrl)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJson(w, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := httpClient.Do(req)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusAccepted {
		log.Println(res.StatusCode)
		app.errorJson(w, errors.New("error sending email"), http.StatusInternalServerError)
		return
	}
	var response JsonResponse
	response.Error = false
	response.Message = "Mail Sent"
	err = app.writeJson(w, http.StatusAccepted, response)
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}
}

func (app *Application) logViaRabbitMQ(w http.ResponseWriter, payload LogPayload) {
	err := app.publishRabbitMQEvent(payload.Name, payload.Data)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	var response JsonResponse
	response.Error = false
	response.Message = "Logged Via RabbitMQ"

	app.writeJson(w, http.StatusAccepted, response)
}

func (app *Application) publishRabbitMQEvent(name, msg string) error {
	emitter, err := events.NewEventEmitter(app.RabbitConn)
	if err != nil {
		return  err
	}
	payload := LogPayload {
		Name: name,
		Data: msg,
	}
	jsonData, _ := json.Marshal(payload)
	err = emitter.Push(string(jsonData), "log.INFO")
	if err != nil {
		return  err
	}
	return nil
}

type RPCPayload struct{
	Name string
	Data string
}

func (app *Application) logViaRPC(w http.ResponseWriter, payload LogPayload) {
	client, err := rpc.Dial("tcp", "logger-service:5002")
	if err != nil {
		app.errorJson(w,err)
		return
	}
	rpcPayload := RPCPayload{
		Name: payload.Name,
		Data: payload.Data,
	}

	var result string

	err = client.Call("RPCServer.LogInfo", rpcPayload, &result)
	if err != nil {
		app.errorJson(w,err)
		return
	}
	response := JsonResponse {
		Error: false,
		Message: result,
	}
	app.writeJson(w, http.StatusAccepted, response)
}
