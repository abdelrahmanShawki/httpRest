package main

import (
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (app *application) readIDparam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64) // base , 64 bit size
	if err != nil || id < 1 {

		return 0, errors.New("invalid id parameter")
	}
	return id, nil
}

type envelop map[string]interface{}

func (app *application) writeJson(w http.ResponseWriter, status int, data envelop, headers http.Header) error {

	//  i want to hilight that json encoder dont use heap memory allocation as same as Marshal
	js, err := json.Marshal(data)
	if err != nil {
		//return some error ,
	}
	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}
