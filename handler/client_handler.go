package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bagasadiii/maxcloud_vps/model/req"
	"github.com/bagasadiii/maxcloud_vps/service"
	"github.com/bagasadiii/maxcloud_vps/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type ClientHandler struct {
	service service.ClientServiceImpl
	logger  *zap.Logger
}

func NewClientHandler(service service.ClientServiceImpl, logger *zap.Logger) *ClientHandler {
	return &ClientHandler{
		service: service,
		logger:  logger,
	}
}

func (ch *ClientHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	var input req.NewClient
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || r.Body == nil {
		ch.logger.Error(utils.ErrBadRequest.Error(), zap.Error(utils.ErrBadRequest))
		utils.JSONResponse(w, http.StatusBadRequest, err)
		return
	}
	if err := ch.service.CreateClientService(r.Context(), &input); err != nil {
		status := utils.ErrCheck(err)
		utils.JSONResponse(w, status, err)
		return
	}
	utils.JSONResponse(w, http.StatusCreated, fmt.Sprintf("%s created, you choose %s plan", input.Email, input.Plan))
}

func (ch *ClientHandler) GetClientInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clientIDString := vars["client_id"]
	clientID, err := uuid.Parse(clientIDString)
	if err != nil {
		info := "id not found or invalid ID"
		ch.logger.Error(utils.ErrNotFound.Error(), zap.String("error", info), zap.Error(err))
		utils.JSONResponse(w, http.StatusNotFound, err)
		return
	}
	res, err := ch.service.GetClientInfoService(r.Context(), clientID)
	if err != nil {
		status := utils.ErrCheck(err)
		utils.JSONResponse(w, status, err)
		return
	}
	utils.JSONResponse(w, http.StatusOK, res)
}
