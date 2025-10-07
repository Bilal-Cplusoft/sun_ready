package handler

import (
	"github.com/Bilal-Cplusoft/sun_ready/internal/service"
	"net/http"
	"encoding/json"
)

type QuoteHandler struct {
	quoteService *service.QuoteService
}

func NewQuoteHandler(quoteService *service.QuoteService) *QuoteHandler {
	return &QuoteHandler{quoteService: quoteService}
}


// GetQuote godoc
// @Summary      Calculate solar quote
// @Description  Takes input parameters for a solar system and returns a detailed quote with costs, savings, and payback period.
// @Tags         quote
// @Accept       json
// @Produce      json
// @Param        quote  body      service.QuoteInput  true  "Quote input payload"
// @Success      200    {object}  service.QuoteResult
// @Failure      400    {object}  map[string]string  "Invalid request payload"
// @Failure      500    {object}  map[string]string  "Failed to calculate quote"
// @Router       /quote [post]
func (h *QuoteHandler) GetQuote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input service.QuoteInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request payload: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	result, err := h.quoteService.CalculateQuote(input)
	if err != nil {
		http.Error(w, "failed to calculate quote: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
