package client

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"
    "os"
    "log"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type TwilioClient struct {
	client        *twilio.RestClient
	fromNumber    string
	otpStore      map[string]OTPData
	mu            sync.RWMutex
	otpExpiration time.Duration
}

type OTPData struct {
	Code      string
	ExpiresAt time.Time
	Attempts  int
}

func InitializeTwilio() *TwilioClient {
	accountSID := os.Getenv("TWILIO_SID")
	authToken := os.Getenv("TWILIO_AUTH")
	fromNumber := os.Getenv("TWILIO_FROM")
	if accountSID == "" || authToken == "" || fromNumber == "" {
     log.Fatal("Twilio SID, Auth Token, or From Number not set in environment")
	}
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSID,
		Password: authToken,
	})

	return &TwilioClient{
		client:        client,
		fromNumber:    fromNumber,
		otpStore:      make(map[string]OTPData),
		otpExpiration: 10 * time.Minute,
	}
}

func (tc *TwilioClient) SendOTP(phoneNumber string) error {
	otp, err := generateOTP(6)
	if err != nil {
		return fmt.Errorf("failed to generate OTP: %w", err)
	}

	tc.mu.Lock()
	tc.otpStore[phoneNumber] = OTPData{
		Code:      otp,
		ExpiresAt: time.Now().Add(tc.otpExpiration),
		Attempts:  0,
	}
	tc.mu.Unlock()

	body := fmt.Sprintf("Your verification code is %s", otp)
	params := &openapi.CreateMessageParams{}
	params.SetTo(phoneNumber)
	params.SetFrom(tc.fromNumber)
	params.SetBody(body)

	_, err = tc.client.Api.CreateMessage(params)
	if err != nil {
		tc.mu.Lock()
		delete(tc.otpStore, phoneNumber)
		tc.mu.Unlock()
		return fmt.Errorf("failed to send OTP SMS: %w", err)
	}

	fmt.Printf("OTP sent to %s: %s\n", phoneNumber, otp)
	return nil
}

func (tc *TwilioClient) VerifyOTP(phoneNumber, otp string) error {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	otpData, exists := tc.otpStore[phoneNumber]
	if !exists {
		return errors.New("no OTP found for this phone number")
	}
	if time.Now().After(otpData.ExpiresAt) {
		delete(tc.otpStore, phoneNumber)
		return errors.New("OTP has expired")
	}
	if otpData.Attempts >= 3 {
		delete(tc.otpStore, phoneNumber)
		return errors.New("maximum verification attempts exceeded")
	}
	if otpData.Code != otp {
		otpData.Attempts++
		tc.otpStore[phoneNumber] = otpData
		return fmt.Errorf("invalid OTP (attempt %d/3)", otpData.Attempts)
	}

	delete(tc.otpStore, phoneNumber)
	fmt.Printf("OTP verified successfully for %s\n", phoneNumber)
	return nil
}

func generateOTP(length int) (string, error) {
	const digits = "0123456789"
	otp := make([]byte, length)
	for i := range otp {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		otp[i] = digits[num.Int64()]
	}
	return string(otp), nil
}

func (tc *TwilioClient) CleanupExpiredOTPs() {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	now := time.Now()
	for phone, data := range tc.otpStore {
		if now.After(data.ExpiresAt) {
			delete(tc.otpStore, phone)
		}
	}
}
