package main

import "net/http"

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