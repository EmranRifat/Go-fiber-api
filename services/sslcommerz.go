// services/sslcommerz.go
package services

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type SSLCommerzInitResponse struct {
	Status         string `json:"status"`
	SessionKey     string `json:"sessionkey"`
	GatewayPageURL string `json:"GatewayPageURL"`
	FailedReason   string `json:"failedreason"`
}

type SSLCommerzValidationResponse struct {
	Status    string `json:"status"`
	TranID    string `json:"tran_id"`
	ValID     string `json:"val_id"`
	Amount    string `json:"amount"`
	Currency  string `json:"currency"`
	RiskLevel string `json:"risk_level"`
}

func sslBaseURL() string {
	if os.Getenv("SSLC_IS_SANDBOX") == "true" {
		return "https://sandbox.sslcommerz.com"
	}
	return "https://securepay.sslcommerz.com"
}

func sslCredentials() (storeID, storePass string, err error) {
	storeID = os.Getenv("SSLC_STORE_ID")
	storePass = os.Getenv("SSLC_STORE_PASSWORD")
	if storeID == "" || storePass == "" {
		return "", "", errors.New("SSLCommerz store credentials are not configured")
	}
	return storeID, storePass, nil
}

func CreateSSLSession(tranID string, amount float64, currency, customerName, customerEmail, customerPhone string) (*SSLCommerzInitResponse, error) {
	storeID, storePass, err := sslCredentials()
	if err != nil {
		return nil, err
	}

	backendURL := os.Getenv("BACKEND_URL")
	if backendURL == "" {
		return nil, errors.New("BACKEND_URL is not configured")
	}
	if currency == "" {
		currency = "BDT"
	}

	data := url.Values{}
	data.Set("store_id", storeID)
	data.Set("store_passwd", storePass)
	data.Set("total_amount", strconv.FormatFloat(amount, 'f', 2, 64))
	data.Set("currency", currency)
	data.Set("tran_id", tranID)

	data.Set("success_url", backendURL+"/api/payment/ssl/success")
	data.Set("fail_url", backendURL+"/api/payment/ssl/fail")
	data.Set("cancel_url", backendURL+"/api/payment/ssl/cancel")
	data.Set("ipn_url", backendURL+"/api/payment/ssl/ipn")

	data.Set("cus_name", customerName)
	data.Set("cus_email", customerEmail)
	data.Set("cus_add1", "Dhaka")
	data.Set("cus_city", "Dhaka")
	data.Set("cus_country", "Bangladesh")
	data.Set("cus_phone", customerPhone)

	data.Set("shipping_method", "NO")
	data.Set("product_name", "Booking Payment")
	data.Set("product_category", "Service")
	data.Set("product_profile", "general")

	req, err := http.NewRequest(
		http.MethodPost,
		sslBaseURL()+"/gwprocess/v4/api.php",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result SSLCommerzInitResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Status != "SUCCESS" || result.GatewayPageURL == "" {
		return nil, errors.New(result.FailedReason)
	}

	return &result, nil
}

func ValidateSSLTransaction(valID string) (*SSLCommerzValidationResponse, error) {
	if valID == "" {
		return nil, errors.New("val_id is required")
	}

	storeID, storePass, err := sslCredentials()
	if err != nil {
		return nil, err
	}

	q := url.Values{}
	q.Set("val_id", valID)
	q.Set("store_id", storeID)
	q.Set("store_passwd", storePass)
	q.Set("format", "json")
	q.Set("v", "1")

	req, err := http.NewRequest(
		http.MethodGet,
		sslBaseURL()+"/validator/api/validationserverAPI.php?"+q.Encode(),
		nil,
	)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result SSLCommerzValidationResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	switch result.Status {
	case "VALID", "VALIDATED":
		return &result, nil
	default:
		return nil, errors.New("payment validation failed: " + result.Status)
	}
}