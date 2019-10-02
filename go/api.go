package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	IsucariAPIToken = "Bearer 75ugk2m37a750fwir5xr-22l6h4wmue1bwrubzwd0"

	userAgent = "isucon9-qualify-webapp"
)

type APIPaymentServiceTokenReq struct {
	ShopID string `json:"shop_id"`
	Token  string `json:"token"`
	APIKey string `json:"api_key"`
	Price  int    `json:"price"`
}

type APIPaymentServiceTokenRes struct {
	Status string `json:"status"`
}

type APIShipmentCreateReq struct {
	ToAddress   string `json:"to_address"`
	ToName      string `json:"to_name"`
	FromAddress string `json:"from_address"`
	FromName    string `json:"from_name"`
}

type APIShipmentCreateRes struct {
	ReserveID   string `json:"reserve_id"`
	ReserveTime int64  `json:"reserve_time"`
}

type APIShipmentRequestReq struct {
	ReserveID string `json:"reserve_id"`
}

type APIShipmentStatusRes struct {
	Status      string `json:"status"`
	ReserveTime int64  `json:"reserve_time"`
}

type APIShipmentStatusReq struct {
	ReserveID string `json:"reserve_id"`
}

const retryLimit = 20

func APIPaymentToken(ctx context.Context, paymentURL string, param *APIPaymentServiceTokenReq) (*APIPaymentServiceTokenRes, error) {
	b, _ := json.Marshal(param)
	return APIPaymentTokenTry(ctx, paymentURL, b, 0)
}

func APIPaymentTokenTry(ctx context.Context, paymentURL string, b []byte, count int) (*APIPaymentServiceTokenRes, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, paymentURL+"/token", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/json")
	//log.Printf("api: %v", req.URL)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		rb, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read res.Body and the status code of the response from shipment service was not 200: %v", err)
		}
		if count < retryLimit {
			log.Printf("retry: %v", string(rb))
			return APIPaymentTokenTry(ctx, paymentURL, b, count+1)
		}
		return nil, fmt.Errorf("status code: %d; body: %s", res.StatusCode, rb)
	}

	pstr := &APIPaymentServiceTokenRes{}
	err = json.NewDecoder(res.Body).Decode(pstr)
	if err != nil {
		return nil, err
	}

	return pstr, nil
}

func APIShipmentCreate(ctx context.Context, shipmentURL string, param *APIShipmentCreateReq) (*APIShipmentCreateRes, error) {
	b, _ := json.Marshal(param)
	return APIShipmentCreateTry(ctx, shipmentURL, b, 0)
}
func APIShipmentCreateTry(ctx context.Context, shipmentURL string, b []byte, count int) (*APIShipmentCreateRes, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, shipmentURL+"/create", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", IsucariAPIToken)

	//log.Printf("api: %v", req.URL)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		rb, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read res.Body and the status code of the response from shipment service was not 200: %v", err)
		}
		if count < retryLimit {
			log.Printf("retry: %v", string(rb))
			return APIShipmentCreateTry(ctx, shipmentURL, b, count+1)
		}
		return nil, fmt.Errorf("status code: %d; body: %s", res.StatusCode, rb)
	}

	scr := &APIShipmentCreateRes{}
	err = json.NewDecoder(res.Body).Decode(&scr)
	if err != nil {
		return nil, err
	}

	return scr, nil
}

func APIShipmentRequest(ctx context.Context, shipmentURL string, param *APIShipmentRequestReq) ([]byte, error) {
	b, _ := json.Marshal(param)
	return APIShipmentRequestTry(ctx, shipmentURL, b, 0)
}

func APIShipmentRequestTry(ctx context.Context, shipmentURL string, b []byte, count int) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, shipmentURL+"/request", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", IsucariAPIToken)

	//log.Printf("api: %v", req.URL)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		rb, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read res.Body and the status code of the response from shipment service was not 200: %v", err)
		}
		if count < retryLimit {
			log.Printf("retry: %v", string(rb))
			return APIShipmentRequestTry(ctx, shipmentURL, b, count+1)
		}
		return nil, fmt.Errorf("status code: %d; body: %s", res.StatusCode, rb)
	}

	return ioutil.ReadAll(res.Body)
}

func APIShipmentStatus(ctx context.Context, shipmentURL string, param *APIShipmentStatusReq) (*APIShipmentStatusRes, error) {
	b, _ := json.Marshal(param)
	return APIShipmentStatusTry(ctx, shipmentURL, b, 0)
}

func APIShipmentStatusTry(ctx context.Context, shipmentURL string, b []byte, count int) (*APIShipmentStatusRes, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, shipmentURL+"/status", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", IsucariAPIToken)

	//log.Printf("api: %v", req.URL)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		rb, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read res.Body and the status code of the response from shipment service was not 200: %v", err)
		}
		if count < retryLimit {
			log.Printf("retry: %v", string(rb))
			return APIShipmentStatusTry(ctx, shipmentURL, b, count+1)
		}
		return nil, fmt.Errorf("status code: %d; body: %s", res.StatusCode, rb)
	}

	ssr := &APIShipmentStatusRes{}
	err = json.NewDecoder(res.Body).Decode(&ssr)
	if err != nil {
		return nil, err
	}

	return ssr, nil
}
