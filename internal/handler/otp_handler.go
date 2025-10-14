package handler

import (
	"github.com/Bilal-Cplusoft/sun_ready/internal/client"
	"net/http"
	"encoding/json"
)

type OtpHandler struct {
	twilioClient *client.TwilioClient
}

func NewOtpHandler(twilioClient *client.TwilioClient) *OtpHandler {
	return &OtpHandler{
		twilioClient: twilioClient,
	}
}


// SendOTP godoc
// @Summary      Send OTP to a phone number
// @Description  Sends a one-time password (OTP) via SMS to the specified phone number using Twilio.
// @Tags         OTP
// @Accept       json
// @Produce      json
// @Param        phone   query     string  true  "Phone number with country code (e.g. +923001234567)"
// @Success      200     {object}  map[string]string  "OTP sent successfully"
// @Failure      400     {object}  map[string]string  "Missing or invalid phone parameter"
// @Failure      500     {object}  map[string]string  "Failed to send OTP"
// @Router       /api/otp/send [get]
func (h *OtpHandler) SendOTP(w http.ResponseWriter, r *http.Request) {
	phone := r.URL.Query().Get("phone")
	if phone == "" {
		http.Error(w, "missing phone parameter", http.StatusBadRequest)
		return
	}

	err := h.twilioClient.SendOTP(phone)
	if err != nil {
		http.Error(w, "failed to send OTP: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "OTP sent successfully",
	})
}

// VerifyOTP godoc
// @Summary      Verify an OTP code
// @Description  Verifies a one-time password (OTP) sent to a phone number.
// @Tags         OTP
// @Accept       json
// @Produce      json
// @Param        phone   query     string  true  "Phone number with country code (e.g. +923001234567)"
// @Param        otp     query     string  true  "OTP code received via SMS"
// @Success      200     {object}  map[string]string  "OTP verified successfully"
// @Failure      400     {object}  map[string]string  "Missing or invalid parameters / OTP verification failed"
// @Router       /api/otp/verify [get]
func (h *OtpHandler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	phone := r.URL.Query().Get("phone")
	otp := r.URL.Query().Get("otp")

	if phone == "" || otp == "" {
		http.Error(w, "missing phone or otp parameter", http.StatusBadRequest)
		return
	}

	err := h.twilioClient.VerifyOTP(phone, otp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "OTP verified successfully",
	})
}
