// Package bse provides a client for the BSE StAR MF API.
// BSE StAR MF is India's national exchange-based mutual fund platform.
//
// API credentials (MemberCode, UserID, Password) are loaded from environment
// variables via the config package and must never be hard-coded.
package bse

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// OrderFlag distinguishes new orders from amendments.
type OrderFlag string

const (
	OrderFlagNew    OrderFlag = "N"
	OrderFlagAmend  OrderFlag = "A"
	OrderFlagCancel OrderFlag = "C"
)

// TransactionType maps to the BSE transaction type codes.
type TransactionType string

const (
	TransactionTypePurchase  TransactionType = "P"
	TransactionTypeRedemption TransactionType = "R"
	TransactionTypeSIPNew    TransactionType = "SN"
)

// OrderRequest is sent to the BSE PlaceOrder endpoint.
type OrderRequest struct {
	MemberCode      string          `json:"TransCode"`
	ClientCode      string          `json:"ClientCode"`
	SchemeCode      string          `json:"SchemeCode"`
	ISIN            string          `json:"ISIN"`
	BuySell         TransactionType `json:"BuySell"`
	BuySellType     string          `json:"BuySellType"` // FRESH / ADDITIONAL
	Amount          float64         `json:"OrderVal"`
	Qty             float64         `json:"Qty"`
	OrderFlag       OrderFlag       `json:"OrderFlag"`
	FolioNo         string          `json:"FolioNo"`
	Remarks         string          `json:"Remarks"`
	KYCStatus       string          `json:"KYCStatus"`
	RefNo           string          `json:"RefNo"`
	SubBrokerCode   string          `json:"SubBrCode"`
	DPC             string          `json:"DPC"`
	IPAdd           string          `json:"IPAdd"`
	Password        string          `json:"Password"`
	PassKey         string          `json:"PassKey"`
}

// OrderResponse is returned by the BSE API on order submission.
type OrderResponse struct {
	StatusCode    string `json:"BseCode"`
	StatusMessage string `json:"BseMsg"`
	OrderID       string `json:"OrderID"`
	RejectCode    string `json:"RejCode"`
}

// SchemeListItem holds minimal scheme metadata returned by BSE.
type SchemeListItem struct {
	SchemeCode  string `json:"SchemeCode"`
	ISINDiv     string `json:"ISINDiv"`
	ISINGrowth  string `json:"ISINGro"`
	SchemeName  string `json:"SchemeName"`
	AMCCode     string `json:"AMCCode"`
	MinPurchase string `json:"MinPurchaseAmount"`
	MinSIP      string `json:"SIPMinAmount"`
}

// NAVRecord holds daily NAV data for a scheme.
type NAVRecord struct {
	SchemeCode string  `json:"scheme_code"`
	Date       string  `json:"date"`
	NAV        float64 `json:"nav"`
}

// Client is a BSE StAR MF API client.
type Client struct {
	memberCode string
	userID     string
	password   string
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new BSE StAR MF API client.
func NewClient(memberCode, userID, password, baseURL string) *Client {
	return &Client{
		memberCode: memberCode,
		userID:     userID,
		password:   password,
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// PlaceOrder submits a mutual fund order to BSE StAR MF.
func (c *Client) PlaceOrder(ctx context.Context, req *OrderRequest) (*OrderResponse, error) {
	if req == nil {
		return nil, errors.New("bse: order request must not be nil")
	}

	req.MemberCode = c.memberCode
	req.Password = c.password

	endpoint := fmt.Sprintf("%s/GetOrderEntryParam", c.baseURL)
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("bse: marshal order request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("bse: create HTTP request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-BSE-Member", c.memberCode)
	httpReq.Header.Set("X-BSE-User", c.userID)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("bse: HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bse: unexpected HTTP status %d", resp.StatusCode)
	}

	var orderResp OrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&orderResp); err != nil {
		return nil, fmt.Errorf("bse: decode response: %w", err)
	}

	return &orderResp, nil
}

// GetSchemeList retrieves the full list of active mutual fund schemes from BSE.
func (c *Client) GetSchemeList(ctx context.Context) ([]SchemeListItem, error) {
	params := url.Values{}
	params.Set("Flag", "0")
	params.Set("Euin", "")

	endpoint := fmt.Sprintf("%s/GetSchemeList?%s", c.baseURL, params.Encode())
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("bse: create HTTP request: %w", err)
	}
	httpReq.Header.Set("X-BSE-Member", c.memberCode)
	httpReq.Header.Set("X-BSE-User", c.userID)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("bse: HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bse: unexpected HTTP status %d", resp.StatusCode)
	}

	var schemes []SchemeListItem
	if err := json.NewDecoder(resp.Body).Decode(&schemes); err != nil {
		return nil, fmt.Errorf("bse: decode response: %w", err)
	}

	return schemes, nil
}

// GetOrderStatus polls BSE for the current status of an order.
// BSE does not provide webhooks, so callers must poll periodically.
func (c *Client) GetOrderStatus(ctx context.Context, orderID string) (*OrderResponse, error) {
	if orderID == "" {
		return nil, errors.New("bse: orderID must not be empty")
	}

	params := url.Values{}
	params.Set("OrderId", orderID)
	params.Set("MemberCode", c.memberCode)

	endpoint := fmt.Sprintf("%s/GetOrderStatus?%s", c.baseURL, params.Encode())
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("bse: create HTTP request: %w", err)
	}
	httpReq.Header.Set("X-BSE-Member", c.memberCode)
	httpReq.Header.Set("X-BSE-User", c.userID)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("bse: HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bse: unexpected HTTP status %d", resp.StatusCode)
	}

	var statusResp OrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&statusResp); err != nil {
		return nil, fmt.Errorf("bse: decode response: %w", err)
	}

	return &statusResp, nil
}
