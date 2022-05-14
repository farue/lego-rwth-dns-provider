package main

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"time"
)

const baseUrl = "https://noc-portal.rz.rwth-aachen.de/dns-admin/api/v1"

type ApiClient struct {
	httpClient *http.Client
}

type ApiError struct {
	StatusCode    int
	RequestMethod string
	RequestUrl    string
	Body          string
	Message       string
}

func (e ApiError) Error() string {
	str := ""
	if e.Message != "" {
		str += e.Message
		str += "\n"
	}
	str += strings.Join([]string{strconv.Itoa(e.StatusCode), e.RequestMethod, e.RequestUrl}, " :: ")
	if e.Body != "" {
		str += "\n"
		str += e.Body
	}
	return str
}

// https://mholt.github.io/json-to-go/

type DeployZoneResponse struct {
	ID         int       `json:"id"`
	ZoneName   string    `json:"zone_name"`
	Status     string    `json:"status"`
	UpdatedAt  time.Time `json:"updated_at"`
	LastDeploy time.Time `json:"last_deploy"`
	Dnssec     Dnssec    `json:"dnssec"`
}

type Zone struct {
	ID         int       `json:"id"`
	ZoneName   string    `json:"zone_name"`
	Status     string    `json:"status"`
	UpdatedAt  time.Time `json:"updated_at"`
	LastDeploy time.Time `json:"last_deploy"`
	Dnssec     Dnssec    `json:"dnssec"`
}

type ZoneSigningKey struct {
	CreatedAt time.Time `json:"created_at"`
}

type KeySigningKey struct {
	CreatedAt time.Time `json:"created_at"`
}

type Dnssec struct {
	ZoneSigningKey ZoneSigningKey `json:"zone_signing_key"`
	KeySigningKey  KeySigningKey  `json:"key_signing_key"`
}

type ListZonesResponse struct {
	Zones []Zone
}

type Record struct {
	ID        int       `json:"id"`
	ZoneID    int       `json:"zone_id"`
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
	Editable  bool      `json:"editable"`
}

type ListRecordsResponse struct {
	Records []Record
}

type CreateRecordsResponse struct {
	ID        int       `json:"id"`
	ZoneID    int       `json:"zone_id"`
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DeleteRecordResponse struct {
	ID        int       `json:"id"`
	ZoneID    int       `json:"zone_id"`
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewApiClient(client *http.Client) *ApiClient {
	return &ApiClient{
		httpClient: client,
	}
}

func (c *ApiClient) listZones(apiToken string) (*ListZonesResponse, error) {
	response := ListZonesResponse{}
	err := c.doRequest("GET", baseUrl+"/list_zones", nil, apiToken, &response.Zones)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *ApiClient) listRecords(apiToken string, zoneId int, search string) (*ListRecordsResponse, error) {
	response := ListRecordsResponse{}

	bodyStr := fmt.Sprintf("zone_id=%d", zoneId)
	if search != "" {
		bodyStr += fmt.Sprintf("&search=%s", search)
	}
	body := strings.NewReader(bodyStr)

	err := c.doRequest("GET", baseUrl+"/list_records", body, apiToken, &response.Records)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *ApiClient) deployZone(apiToken string, zoneId int) (*DeployZoneResponse, error) {
	var response DeployZoneResponse
	body := strings.NewReader(fmt.Sprintf("zone_id=%d", zoneId))
	err := c.doRequest("POST", baseUrl+"/deploy_zone", body, apiToken, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *ApiClient) createRecords(apiToken string, record string) (*CreateRecordsResponse, error) {
	var response CreateRecordsResponse
	body := strings.NewReader(fmt.Sprintf("record_content=%s", record))
	err := c.doRequest("POST", baseUrl+"/create_record", body, apiToken, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *ApiClient) destroyRecords(apiToken string, recordId int) (*DeleteRecordResponse, error) {
	var response DeleteRecordResponse
	body := strings.NewReader(fmt.Sprintf("record_id=%d", recordId))
	err := c.doRequest("DELETE", baseUrl+"/destroy_record", body, apiToken, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *ApiClient) doRequest(method string, url string, body io.Reader, apiToken string, response interface{}) error {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}

	req.Header.Set("PRIVATE-TOKEN", apiToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	debugRequest(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	debugResponse(resp)

	if err = checkError(req, resp); err != nil {
		return err
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()
	err = json.Unmarshal(raw, &response)
	if err != nil {
		return fmt.Errorf("failed to unmarshal %q to type %T: %w", raw, ListZonesResponse{}, err)
	}

	return nil
}

func checkError(req *http.Request, resp *http.Response) error {
	if resp.StatusCode >= http.StatusBadRequest {
		e := ApiError{
			StatusCode:    resp.StatusCode,
			RequestMethod: req.Method,
			RequestUrl:    req.URL.String(),
			Body:          "",
			Message:       "",
		}

		body, err := io.ReadAll(resp.Body)
		if err == nil {
			e.Body = string(body)
		}
		return e
	}
	return nil
}

func debugRequest(req *http.Request) {
	reqDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Debug().Msg("failed to dump request")
	}
	log.Debug().Msgf("REQUEST:\n%s", string(reqDump))
}

func debugResponse(resp *http.Response) {
	respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Debug().Msg("failed to dump response")
	}
	log.Debug().Msgf("RESPONSE:\n%s", string(respDump))
}
