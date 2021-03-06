package hyper

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type FipCli struct {
	hyperCli *HyperConn
}

func NewFipCli(client *HyperConn) *FipCli {
	return &FipCli{
		hyperCli: client,
	}
}

func (f *FipCli) AllocateFip(count int) (int, []FipResponse, error) {
	var (
		result     string
		httpStatus int
		err        error
	)
	method := "POST"
	endpoint := fmt.Sprintf("/api/v1/hyper/fips?count=%v", count)
	result, httpStatus, err = f.hyperCli.SockRequest(method, endpoint, nil, "")
	if err != nil {
		log.Fatalf("send request error: %v", err)
	} else if httpStatus != http.StatusCreated {
		log.Fatalf("response error: %v - %v", httpStatus, result)
	}
	var fipListAllocated []FipResponse
	if err = json.Unmarshal([]byte(result), &fipListAllocated); err != nil {
		log.Fatalf("failed to parse allocated fip list")
	}
	return httpStatus, fipListAllocated, nil
}

func (f *FipCli) ListFips() (int, []FipResponse, error) {
	method := "GET"
	endpoint := "/api/v1/hyper/fips"

	result, httpStatus, err := f.hyperCli.SockRequest(method, endpoint, nil, "")
	if err != nil {
		log.Fatalf("send request error: %v", err)
	} else if httpStatus != http.StatusOK {
		log.Fatalf("response error: %v - %v", httpStatus, result)
	}
	var fipList []FipResponse
	json.Unmarshal([]byte(result), &fipList)
	return httpStatus, fipList, nil
}

func (f *FipCli) GetFip(ip string) (int, *FipResponse, error) {
	if ip == "" {
		log.Fatal("Please specify ip")
	}

	method := "GET"
	endpoint := fmt.Sprintf("/api/v1/hyper/fips/%v", url.QueryEscape(ip))

	result, httpStatus, err := f.hyperCli.SockRequest(method, endpoint, nil, "")
	if err != nil {
		log.Fatalf("send request error: %v", err)
	} else if httpStatus != http.StatusOK {
		log.Fatalf("response error: %v - %v", httpStatus, result)
	}
	var fip FipResponse
	err = json.Unmarshal([]byte(result), &fip)
	if err != nil {
		log.Fatalf("failed to convert result to fip:%v", err)
	}
	return httpStatus, &fip, nil
}

func (f *FipCli) NameFip(ip, name string) (int, string, error) {
	if ip == "" {
		log.Fatal("Please specify ip")
	}
	if name == "" {
		log.Fatal("Please specify --name")
	}
	method := "POST"
	endpoint := fmt.Sprintf("/api/v1/hyper/fips/%v", ip)
	data := fmt.Sprintf(`{"name":"%v"}`, name)
	result, httpStatus, err := f.hyperCli.SockRequest(method, endpoint, strings.NewReader(data), "application/json")
	if err != nil {
		log.Fatalf("send request error: %v", err)
	} else if httpStatus != http.StatusNoContent {
		log.Fatalf("response error: %v - %v", httpStatus, result)
	}
	return httpStatus, result, nil
}

func (f *FipCli) ReleaseFip(ip string) (int, string) {
	if ip == "" {
		log.Fatal("Please specify ip")
	}

	method := "DELETE"
	endpoint := fmt.Sprintf("/api/v1/hyper/fips/%v", ip)

	result, httpStatus, err := f.hyperCli.SockRequest(method, endpoint, nil, "")
	if err != nil {
		log.Fatalf("send request error: %v", err)
	} else if httpStatus != http.StatusNoContent {
		log.Fatalf("response error: %v - %v", httpStatus, result)
	}
	return httpStatus, result
}

func (f *FipCli) ReleaseAllFips() {
	method := "GET"
	endpoint := "/api/v1/hyper/fips"

	result, httpStatus, err := f.hyperCli.SockRequest(method, endpoint, nil, "")
	if err != nil {
		log.Fatalf("send request error: %v", err)
	} else if httpStatus != http.StatusOK {
		log.Fatalf("response error: %v - %v", httpStatus, result)
	}

	var fipList []FipResponse
	err = json.Unmarshal([]byte(result), &fipList)
	if err != nil {
		log.Fatalf("failed to parse fip list:%v", err)
	}
	for _, i := range fipList {
		f.ReleaseFip(i.Fip)
	}
}
