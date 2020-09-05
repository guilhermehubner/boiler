package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"strconv"
	"strings"

	"boiler/pkg/iface"
	"boiler/pkg/store/log"

	"github.com/go-chi/chi"
)

// AddUserHandle handle an AddUser request
func AddUserHandle(service iface.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := struct {
			Name     string `json:"name"`
			Password string `json:"password"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			Fail(w, r, http.StatusBadRequest, "could not parse payload")
			return
		}

		payload.Name = strings.TrimSpace(payload.Name)
		if len(payload.Name) == 0 {
			Fail(w, r, http.StatusBadRequest, "empty name")
			return
		}

		userID, err := service.AddUser(r.Context(), payload.Name, payload.Password)
		if err != nil {
			log.Log(err)
			Fail(w, r, http.StatusInternalServerError, "service failed")
			return
		}

		JSON(w, r, map[string]interface{}{
			"user_id": userID,
		})
	}
}

// ListUsersHandle handle an ListUsers request
func ListUsersHandle(service iface.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit := 100
		rawLimit := r.URL.Query()["limit"]
		if len(rawLimit) > 0 {
			var err error
			limit, err = strconv.Atoi(rawLimit[0])
			if err != nil || limit <= 0 {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "invalid limit \"%s\"", rawLimit[0])
				return
			}
		}

		users, err := service.FilterUsers(r.Context(), iface.FilterUsers{Limit: uint(limit)})
		if err != nil {
			log.Log(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "service failed")
			return
		}

		JSON(w, r, map[string]interface{}{
			"users": users,
		})
	}
}

// DeleteUserHandle handle an DeleteUser request
func DeleteUserHandle(service iface.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
		if err != nil || userID == 0 {
			Fail(w, r, http.StatusBadRequest, "invalid user ID")
			return
		}

		err = service.DeleteUser(r.Context(), userID)
		if err != nil {
			log.Log(err)
			Fail(w, r, http.StatusInternalServerError, "service failed")
			return
		}

		JSON(w, r, nil)
	}
}

// GetUserHandle handle an GetUser request
func GetUserHandle(service iface.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
		if err != nil || userID == 0 {
			Fail(w, r, http.StatusBadRequest, "invalid user ID")
			return
		}

		user, err := service.GetUserByID(r.Context(), userID)
		if err != nil {
			log.Log(err)
			Fail(w, r, http.StatusInternalServerError, "service failed")
			return
		}

		JSON(w, r, map[string]interface{}{
			"user": user,
		})
	}
}

// AddEmailHandle handle an AddEmail request
func AddEmailHandle(service iface.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := struct {
			UserID  int64  `json:"user_id"`
			Address string `json:"address"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			log.Log(err)
			Fail(w, r, http.StatusBadRequest, "invalid payload")
			return
		}

		email, err := mail.ParseAddress(payload.Address)
		if err != nil {
			Fail(w, r, http.StatusBadRequest, "invalid email address")
			return
		}

		if payload.UserID < 1 {
			Fail(w, r, http.StatusBadRequest, "invalid user ID")
			return
		}

		emailID, err := service.AddEmail(r.Context(), payload.UserID, email.Address)
		if err != nil {
			log.Log(err)
			Fail(w, r, http.StatusInternalServerError, "service failed")
			return
		}

		JSON(w, r, struct {
			EmailID int64 `json:"email_id"`
		}{emailID})
	}
}

// DeleteEmailHandle handle an DeleteEmail request
func DeleteEmailHandle(service iface.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		emailID, err := strconv.ParseInt(chi.URLParam(r, "emailID"), 10, 64)
		if err != nil || emailID <= 0 {
			Fail(w, r, http.StatusBadRequest, "invalid email ID")
			return
		}

		err = service.DeleteEmail(r.Context(), emailID)
		if err != nil {
			log.Log(err)
			Fail(w, r, http.StatusInternalServerError, "service failed")
			return
		}

		JSON(w, r, nil)
	}
}

// ListEmailsHandle handle an ListEmails request
func ListEmailsHandle(service iface.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()["user_id"]
		if len(params) == 0 {
			Fail(w, r, http.StatusBadRequest, "missing URL query user_id")
			return
		}

		userID, err := strconv.ParseInt(params[0], 10, 64)
		if err != nil || userID == 0 {
			Fail(w, r, http.StatusBadRequest, "invalid URL query user_id")
			return
		}

		emails, err := service.FilterEmails(r.Context(), iface.FilterEmails{UserID: userID})
		if err != nil {
			log.Log(err)
			Fail(w, r, http.StatusInternalServerError, "service failed")
			return
		}

		JSON(w, r, map[string]interface{}{
			"emails": emails,
		})
	}
}

// LoginHandle handle an authentication request
func LoginHandle(service iface.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			Fail(w, r, http.StatusBadRequest, "could not parse payload")
			return
		}

		user, token, err := service.AuthUser(r.Context(), payload.Email, payload.Password)
		if err != nil {
			log.Log(err)
			Fail(w, r, http.StatusInternalServerError, "service failed")
			return
		}

		JSON(w, r, map[string]interface{}{
			"user":  user,
			"token": token,
		})
	}
}
