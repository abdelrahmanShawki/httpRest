package main

import (
	"errors"
	"fmt"
	"httpRest/internal/data"
	"httpRest/internal/validator"
	"net/http"
	"time"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJson(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err) //review
		return
	}

	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if data.ValidateUser(v, user); !v.Valid() {
		app.validationErrorResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "is used before")
			app.validationErrorResponse(w, r, v.Errors)
			return

		default:
			app.logger.PrintError(err, nil)
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	// make a goroutine to make the email in , this make it decoupled from the function return enhancing performance
	// defer function of recover to handle any panic in this goroutine
	// add wait group to handle shutdown
	// move to helper later with arbitary excution
	go func() {
		app.waitGroup.Add(1)
		defer app.waitGroup.Done()

		defer func() {
			if err := recover(); err != nil {
				app.logger.PrintError(fmt.Errorf("%s ", err), nil)
			}
		}()

		token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)

		err = app.mailer.Send(user.Email, user.Name, token.PlainText)
		if err != nil {
			app.logger.PrintError(err, nil)
			return
		}
	}()

	err = app.writeJson(w, http.StatusAccepted, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		TokenPlainText string `json:"token"`
	}

	err := app.readJson(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()
	data.ValidateTokenPlainText(v, input.TokenPlainText)
	if !v.Valid() {
		app.validationErrorResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetForToken(data.ScopeActivation, input.TokenPlainText)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			app.validationErrorResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	user.Activated = true

	err = app.models.Users.Update(user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.models.Tokens.DeleteToken(data.ScopeActivation, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJson(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
