package client

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/stlesnik/goph_keeper/internal/config"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/stlesnik/goph_keeper/internal/models"
)

// Client represents the GophKeeper client
type Client struct {
	config     config.ClientConfig
	httpClient *http.Client
	token      string
	encryption *ClientEncryption
}

// NewClient creates a new client instance
func NewClient(config config.ClientConfig) *Client {
	caCert, err := os.ReadFile(config.CertFile)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}

	return &Client{
		config:     config,
		httpClient: client,
		token:      "",
	}
}

// makeRequest makes HTTP request to server
func (c *Client) makeRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, c.config.ServerURL+endpoint, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", c.token)
	}
	return c.httpClient.Do(req)
}

// Register registers a new user
func (c *Client) Register(username, email, password string) error {
	req := models.RegisterUserRequest{
		Username: username,
		Email:    email,
		Password: password,
	}

	resp, err := c.makeRequest("POST", "/user/register", req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("unexpected error while closing body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("registration failed: %s", string(body))
	}

	token := resp.Header.Get("Authorization")
	salt := resp.Header.Get("X-User-Salt")
	if token != "" {
		c.token = token
		c.encryption = NewClientEncryption(password, salt)
	}

	return nil
}

// Login authenticates user
func (c *Client) Login(email, password string) error {
	req := models.LoginUserRequest{
		Email:    email,
		Password: password,
	}

	resp, err := c.makeRequest("POST", "/user/login", req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("unexpected error while closing body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("login failed: %s", string(body))
	}

	token := resp.Header.Get("Authorization")
	salt := resp.Header.Get("X-User-Salt")
	if token != "" {
		c.token = token
		c.encryption = NewClientEncryption(password, salt)
	}

	return nil
}

// GetProfile gets user profile
func (c *Client) GetProfile() (*models.UserProfile, error) {
	resp, err := c.makeRequest("GET", "/user/profile", nil)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("unexpected error while closing body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get profile: %s", string(body))
	}

	var profile models.UserProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, err
	}

	return &profile, nil
}

// GetAllData gets all user data items
func (c *Client) GetAllData(offset int) ([]models.DataItemResponse, error) {
	strOffset := strconv.Itoa(offset)
	resp, err := c.makeRequest("GET", "/data/"+strOffset, nil)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("unexpected error while closing body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get data: %s", string(body))
	}

	var items []models.DataItemResponse
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, err
	}

	return items, nil
}

// GetDataByID gets specific data item
func (c *Client) GetDataByID(id string) (*models.DataItemResponse, error) {
	resp, err := c.makeRequest("GET", "/data/item/"+id, nil)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("unexpected error while closing body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get data item: %s", string(body))
	}

	var item models.DataItemResponse
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		return nil, err
	}

	if c.encryption != nil && item.Data != "" {
		decryptedData, err := c.encryption.DecryptData(item.Data, item.Metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt data: %v", err)
		}
		item.Data = decryptedData
	}

	return &item, nil
}

// CreateData creates new data item
func (c *Client) CreateData(dataType, title, data, metadata string) error {
	if c.encryption == nil {
		return fmt.Errorf("not authenticated - please login first")
	}

	encryptedData, iv, err := c.encryption.EncryptData(data)
	if err != nil {
		return fmt.Errorf("failed to encrypt data: %v", err)
	}

	req := models.CreateDataRequest{
		Type:          dataType,
		Title:         title,
		EncryptedData: encryptedData,
		IV:            iv,
		Metadata:      metadata,
	}

	resp, err := c.makeRequest("POST", "/data/item", req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("unexpected error while closing body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create data: %s", string(body))
	}

	return nil
}

// UpdateData updates existing data item
func (c *Client) UpdateData(id, dataType, title, data, metadata string) error {
	if c.encryption == nil {
		return fmt.Errorf("not authenticated, please login first")
	}

	encryptedData, iv, err := c.encryption.EncryptData(data)
	if err != nil {
		return fmt.Errorf("failed to encrypt data: %v", err)
	}

	req := models.UpdateDataRequest{
		Type:          dataType,
		Title:         title,
		EncryptedData: encryptedData,
		IV:            iv,
		Metadata:      metadata,
	}

	resp, err := c.makeRequest("PUT", "/data/item/"+id, req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("unexpected error while closing body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to update data: %s", string(body))
	}

	return nil
}

// DeleteData deletes data item
func (c *Client) DeleteData(id string) error {
	resp, err := c.makeRequest("DELETE", "/data/item/"+id, nil)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("unexpected error while closing body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete data: %s", string(body))
	}

	return nil
}

// ChangePassword changes user password
func (c *Client) ChangePassword(currentPassword, newPassword string) error {
	req := models.ChangePasswordRequest{
		CurrentPassword: currentPassword,
		NewPassword:     newPassword,
	}

	resp, err := c.makeRequest("PUT", "/user/password", req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("unexpected error while closing body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to change password: %s", string(body))
	}

	return nil
}
