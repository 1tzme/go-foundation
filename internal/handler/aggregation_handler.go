package handler

import (
	"hot-coffee/internal/service"
	"hot-coffee/pkg/logger"
	"net/http"
)

type AggregationHandler struct {
	aggregationService service.AggregationServiceInterface
	logger             *logger.Logger
}

func NewAggregationHandler(s service.AggregationServiceInterface, log *logger.Logger) *AggregationHandler {
	return &AggregationHandler{
		aggregationService: s,
		logger:             log.WithComponent("aggregation_handler"),
	}
}

func (h *AggregationHandler) GetTotalSales(w http.ResponseWriter, r *http.Request) {
	report, err := h.aggregationService.GetTotalSales()
	if err != nil {
		// logger err
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to generate sales report")
		return
	}

	writeJSONResponse(w, http.StatusOK, report)
}

func (h *AggregationHandler) GetPopularItems(w http.ResponseWriter, r *http.Request) {
	report, err := h.aggregationService.GetPopularItems()
	if err != nil {
		// logger err
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to get popular items")
		return
	}

	writeJSONResponse(w, http.StatusOK, report)
}
