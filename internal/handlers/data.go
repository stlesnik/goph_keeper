package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/stlesnik/goph_keeper/internal/logger"
	"github.com/stlesnik/goph_keeper/internal/models"
)

// CreateData handles creation of new data item.
func (h *Handlers) CreateData(w http.ResponseWriter, r *http.Request) {
	var createDataReq models.CreateDataRequest
	err := json.NewDecoder(r.Body).Decode(&createDataReq)
	if err != nil {
		logger.Logger.Errorw("Error decoding create data request", "error", err)
		http.Error(w, "Got error while decoding create data request", http.StatusBadRequest)
		return
	}

	err = h.service.Data.Create(r.Context(), createDataReq)
	if err != nil {
		logger.Logger.Errorw("Create data error", "error", err,
			"error", err.Error(),
			"ip", r.RemoteAddr,
		)
		http.Error(w, "Unable to create data item", http.StatusBadRequest)
		return
	}

	logger.Logger.Infow("Create data success",
		"type", createDataReq.Type,
		"title", createDataReq.Title,
		"ip", r.RemoteAddr,
	)
	w.WriteHeader(http.StatusCreated)
}

// GetAllData handles retrieval of all user data items with pagination.
func (h *Handlers) GetAllData(w http.ResponseWriter, r *http.Request) {
	offset, err := strconv.Atoi(chi.URLParam(r, "offset"))
	if err != nil {
		logger.Logger.Errorw("Error decoding offset", "error", err, "offset", offset)
		offset = 0
	}
	dataItems, err := h.service.Data.GetAll(r.Context(), offset)
	if err != nil {
		logger.Logger.Errorw("Get data error",
			"error", err.Error(),
			"ip", r.RemoteAddr,
		)
		http.Error(w, "Unable to get data items", http.StatusBadRequest)
		return
	}

	logger.Logger.Infow("Get data success",
		"ip", r.RemoteAddr,
	)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(dataItems)
	if err != nil {
		logger.Logger.Errorw("Error writing response", "error", err)
		http.Error(w, "Unable to get data items", http.StatusBadRequest)
		return
	}
}

// GetDataByID handles retrieval of specific data item by ID.
func (h *Handlers) GetDataByID(w http.ResponseWriter, r *http.Request) {
	dataID := chi.URLParam(r, "id")

	dataItem, err := h.service.Data.GetByID(r.Context(), dataID)
	if err != nil {
		logger.Logger.Errorw("Get data by ID error", "error", err,
			"error", err.Error(),
			"data_id", dataID,
			"ip", r.RemoteAddr,
		)
		http.Error(w, "Unable to get data item", http.StatusBadRequest)
		return
	}

	logger.Logger.Infow("Get data by ID success",
		"data_id", dataID,
		"ip", r.RemoteAddr,
	)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(dataItem)
	if err != nil {
		logger.Logger.Errorw("Error writing response", "error", err)
		http.Error(w, "Unable to get data item", http.StatusBadRequest)
		return
	}
}

// UpdateData handles updating of existing data item.
func (h *Handlers) UpdateData(w http.ResponseWriter, r *http.Request) {
	dataID := chi.URLParam(r, "id")

	var updateReq models.UpdateDataRequest
	err := json.NewDecoder(r.Body).Decode(&updateReq)
	if err != nil {
		logger.Logger.Errorw("Error decoding update data request", "error", err)
		http.Error(w, "Got error while decoding update data request", http.StatusBadRequest)
		return
	}

	err = h.service.Data.Update(r.Context(), dataID, updateReq)
	if err != nil {
		logger.Logger.Errorw("Update data error", "error", err,
			"error", err.Error(),
			"data_id", dataID,
			"ip", r.RemoteAddr,
		)
		http.Error(w, "Unable to update data item", http.StatusBadRequest)
		return
	}

	logger.Logger.Infow("Update data success",
		"data_id", dataID,
		"type", updateReq.Type,
		"title", updateReq.Title,
		"ip", r.RemoteAddr,
	)
	w.WriteHeader(http.StatusOK)
}

// DeleteData handles deletion of data item.
func (h *Handlers) DeleteData(w http.ResponseWriter, r *http.Request) {
	dataID := chi.URLParam(r, "id")

	err := h.service.Data.Delete(r.Context(), dataID)
	if err != nil {
		logger.Logger.Errorw("Delete data error", "error", err,
			"error", err.Error(),
			"data_id", dataID,
			"ip", r.RemoteAddr,
		)
		http.Error(w, "Unable to delete data item", http.StatusBadRequest)
		return
	}

	logger.Logger.Infow("Delete data success",
		"data_id", dataID,
		"ip", r.RemoteAddr,
	)
	w.WriteHeader(http.StatusNoContent)
}
