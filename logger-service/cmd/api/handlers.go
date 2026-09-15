package main

import (
	"logger/data"
	"net/http"
)

type RequestPayload struct {
	Name	string 	`json:"name"`
	Data	string 	`json:"data"`
}

func (app *Application) writeLog(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	_ = app.readJson(w, r, &requestPayload)

	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	response := JsonResponse {
		Error: false,
		Message: "Logged",
	}
	app.writeJson(w, http.StatusAccepted, response)
}
