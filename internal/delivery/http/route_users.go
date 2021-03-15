package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
)

func (delivery *Delivery) getUsers(w http.ResponseWriter, r *http.Request) {
	userID, item, err := parseUserPath(r.URL.Path)
	switch item {
	case "information":
		delivery.getUserInformation(w, r, userID)
	case "favorites":
		delivery.getUserFavorites(w, r, userID)
	case "articles":
		delivery.getUserArticles(w, r, userID)
	default:
		delivery.logger.Noticef("user id: %v not exist but be queried, info: %v err: %v", userID, item, err)
		w.WriteHeader(http.StatusNotFound)
	}
}

func (delivery *Delivery) getUserInformation(w http.ResponseWriter, r *http.Request, userID string) {
	token := delivery.getTokenFromRequest(r)

	err := delivery.usecase.CheckPermission(token,
		[]usecase.Permission{usecase.PermissionReadUserInformation},
		map[string]string{
			"user_id": userID,
		})

	if err != nil {
		// TODO: record unauthorized access
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	dataMap, err := delivery.usecase.GetUserInformation(context.Background(), userID)
	if err != nil {
		// TODO: record error
		w.WriteHeader(http.StatusInternalServerError)
		m := map[string]string{
			"error":             "find_userrec_error",
			"error_description": err.Error(),
		}
		b, _ := json.MarshalIndent(m, "", "  ")
		_, err = w.Write(b)
		if err != nil {
			delivery.logger.Errorf("getUserInformation error response err: %w", err)
		}
		return
	}

	responseMap := map[string]interface{}{
		"data": dataMap,
	}
	responseByte, _ := json.MarshalIndent(responseMap, "", "  ")

	_, err = w.Write(responseByte)
	if err != nil {
		delivery.logger.Errorf("getUserInformation success response err: %w", err)
	}
}

func (delivery *Delivery) getUserFavorites(w http.ResponseWriter, r *http.Request, userID string) {
	token := delivery.getTokenFromRequest(r)
	err := delivery.usecase.CheckPermission(token,
		[]usecase.Permission{usecase.PermissionReadUserInformation},
		map[string]string{
			"user_id": userID,
		})

	if err != nil {
		// TODO: record unauthorized access
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	dataItems, err := delivery.usecase.GetUserFavorites(context.Background(), userID)
	if err != nil {
		delivery.logger.Errorf("failed to get user favorites: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	responseMap := map[string]interface{}{
		"data": map[string]interface{}{
			"items": dataItems,
		},
	}

	responseByte, _ := json.MarshalIndent(responseMap, "", "  ")

	w.Write(responseByte)
}

func (delivery *Delivery) getUserArticles(w http.ResponseWriter, r *http.Request, userID string) {
	token := delivery.getTokenFromRequest(r)
	err := delivery.usecase.CheckPermission(token,
		[]usecase.Permission{usecase.PermissionReadUserInformation},
		map[string]string{
			"user_id": userID,
		})

	if err != nil {
		// TODO: record unauthorized access
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// return need fix
	dataItems, err := delivery.usecase.GetUserArticles(context.Background(), userID)
	if err != nil {
		delivery.logger.Errorf("failed to get user's articles: %s\n", err)
	}

	responseMap := map[string]interface{}{
		"data": map[string]interface{}{
			"items": dataItems,
		},
	}

	responseByte, _ := json.MarshalIndent(responseMap, "", "  ")

	_, err = w.Write(responseByte)
	if err != nil {
		delivery.logger.Errorf("getUserFavorites success response err: %w", err)
	}
}
