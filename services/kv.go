package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const KV_BASE_URL = "https://kv.adswds.com/1ep2d8wb"

// KvData KV存储数据结构
type KvData struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// KvResponse KV API响应结构
type KvResponse struct {
	Code int     `json:"code"`
	Msg  string  `json:"msg"`
	Data *KvData `json:"data,omitempty"`
}

// KvService KV存储服务
type KvService struct {
	client *http.Client
}

// NewKvService 创建KV服务实例
func NewKvService() *KvService {
	return &KvService{
		client: &http.Client{},
	}
}

// request 通用请求方法
func (s *KvService) request(method, authorization, key, value string) (*KvResponse, error) {
	reqURL := fmt.Sprintf("%s?key=%s", KV_BASE_URL, url.QueryEscape(key))
	if value != "" {
		reqURL = fmt.Sprintf("%s&value=%s", reqURL, url.QueryEscape(value))
	}

	req, err := http.NewRequest(method, reqURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authorization))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var kvResp KvResponse
	if err := json.Unmarshal(body, &kvResp); err != nil {
		return nil, err
	}

	return &kvResp, nil
}

// CreateKey 创建键值对
func (s *KvService) CreateKey(key, value, authorization string) (*KvResponse, error) {
	return s.request("POST", authorization, key, value)
}

// GetKey 获取键值对
func (s *KvService) GetKey(key, authorization string) (*KvResponse, error) {
	return s.request("GET", authorization, key, "")
}

// UpdateKey 更新键值对
func (s *KvService) UpdateKey(key, value, authorization string) (*KvResponse, error) {
	return s.request("PUT", authorization, key, value)
}

// DeleteKey 删除键值对
func (s *KvService) DeleteKey(key, authorization string) (*KvResponse, error) {
	return s.request("DELETE", authorization, key, "")
}
