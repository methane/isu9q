package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
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

var (
	muCounter  sync.Mutex
	counter    = 0
	choconURLs = []string{
		"http://app01:5000",
		"http://app02:5000",
		"http://app03:5000",
	}
)

const retryLimit = 10

func getChoconURL() string {
	muCounter.Lock()
	defer muCounter.Unlock()
	counter++
	return choconURLs[counter%3]
}

func APIPaymentToken(paymentURL string, param *APIPaymentServiceTokenReq) (*APIPaymentServiceTokenRes, error) {
	b, _ := json.Marshal(param)
	return APIPaymentTokenTry(paymentURL, b, 0)
}

func APIPaymentTokenTry(paymentURL string, b []byte, count int) (*APIPaymentServiceTokenRes, error) {
	req, err := http.NewRequest(http.MethodPost, getChoconURL()+"/token", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/json")
	u, _ := url.Parse(paymentURL)
	if u.Scheme == "https" {
		req.Host = u.Hostname() + ".ccnproxy-https"
		log.Print(req.Host)
	} else {
		req.Host = u.Hostname()
	}

	log.Printf("api: %v", req.URL)
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
			time.Sleep(100 * time.Millisecond)
			return APIPaymentTokenTry(paymentURL, b, count+1)
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

func APIShipmentCreate(shipmentURL string, param *APIShipmentCreateReq) (*APIShipmentCreateRes, error) {
	b, _ := json.Marshal(param)
	return APIShipmentCreateTry(shipmentURL, b, 0)
}
func APIShipmentCreateTry(shipmentURL string, b []byte, count int) (*APIShipmentCreateRes, error) {
	req, err := http.NewRequest(http.MethodPost, getChoconURL()+"/create", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", IsucariAPIToken)
	u, _ := url.Parse(shipmentURL)
	if u.Scheme == "https" {
		req.Host = u.Hostname() + ".ccnproxy-https"
		log.Print(req.Host)
	} else {
		req.Host = u.Hostname()
	}

	log.Printf("api: %v", req.URL)
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
			time.Sleep(100 * time.Millisecond)
			return APIShipmentCreateTry(shipmentURL, b, count+1)
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

func APIShipmentRequest(shipmentURL string, param *APIShipmentRequestReq) ([]byte, error) {
	b, _ := json.Marshal(param)
	return APIShipmentRequestTry(shipmentURL, b, 0)
}

func APIShipmentRequestTry(shipmentURL string, b []byte, count int) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, getChoconURL()+"/request", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", IsucariAPIToken)
	u, _ := url.Parse(shipmentURL)
	if u.Scheme == "https" {
		req.Host = u.Hostname() + ".ccnproxy-https"
		log.Print(req.Host)
	} else {
		req.Host = u.Hostname()
	}

	log.Printf("api: %v", req.URL)
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
			time.Sleep(100 * time.Millisecond)
			return APIShipmentRequestTry(shipmentURL, b, count+1)
		}
		return nil, fmt.Errorf("status code: %d; body: %s", res.StatusCode, rb)
	}

	return ioutil.ReadAll(res.Body)
}

func APIShipmentStatus(shipmentURL string, param *APIShipmentStatusReq) (*APIShipmentStatusRes, error) {
	b, _ := json.Marshal(param)
	return APIShipmentStatusTry(shipmentURL, b, 0)
}

func APIShipmentStatusTry(shipmentURL string, b []byte, count int) (*APIShipmentStatusRes, error) {
	req, err := http.NewRequest(http.MethodGet, getChoconURL()+"/status", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", IsucariAPIToken)
	u, _ := url.Parse(shipmentURL)
	if u.Scheme == "https" {
		req.Host = u.Hostname() + ".ccnproxy-https"
		log.Print(req.Host)
	} else {
		req.Host = u.Hostname()
	}

	log.Printf("api: %v", req.URL)
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
			time.Sleep(100 * time.Millisecond)
			return APIShipmentStatusTry(shipmentURL, b, count+1)
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
