package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// CloudflareService Cloudflare DNS 管理服务
type CloudflareService struct {
	client *http.Client
}

// CloudflareConfig Cloudflare 配置
type CloudflareConfig struct {
	APIToken string `json:"api_token"`
	ZoneID   string `json:"zone_id"`
}

// DNSRecord DNS 记录结构
type DNSRecord struct {
	ID      string `json:"id,omitempty"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Proxied bool   `json:"proxied"`
	TTL     int    `json:"ttl"`
}

// CloudflareResponse Cloudflare API 响应
type CloudflareResponse struct {
	Success bool        `json:"success"`
	Errors  []APIError  `json:"errors"`
	Result  interface{} `json:"result"`
}

// APIError API 错误结构
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// PagesCustomDomain Pages 自定义域名结构
type PagesCustomDomain struct {
	Name   string `json:"name"`
	Status string `json:"status,omitempty"`
	ID     string `json:"id,omitempty"`
}

// PagesProject Pages 项目信息
type PagesProject struct {
	Name    string              `json:"name"`
	Domains []PagesCustomDomain `json:"domains,omitempty"`
}

// NewCloudflareService 创建 Cloudflare 服务实例
func NewCloudflareService() *CloudflareService {
	return &CloudflareService{
		client: &http.Client{},
	}
}

// request 通用请求方法
func (s *CloudflareService) request(method, endpoint string, config CloudflareConfig, data interface{}) (*CloudflareResponse, error) {
	apiURL := fmt.Sprintf("https://api.cloudflare.com/client/v4%s", endpoint)

	var body io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("序列化请求数据失败: %v", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, apiURL, body)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.APIToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var cfResp CloudflareResponse
	if err := json.Unmarshal(respBody, &cfResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if !cfResp.Success {
		if len(cfResp.Errors) > 0 {
			return nil, fmt.Errorf("Cloudflare API 错误: %s", cfResp.Errors[0].Message)
		}
		return nil, fmt.Errorf("Cloudflare API 请求失败")
	}

	return &cfResp, nil
}

// GetDNSRecords 获取 DNS 记录
func (s *CloudflareService) GetDNSRecords(config CloudflareConfig, name, recordType string) ([]DNSRecord, error) {
	endpoint := fmt.Sprintf("/zones/%s/dns_records", config.ZoneID)

	// 添加查询参数
	params := url.Values{}
	if name != "" {
		params.Add("name", name)
	}
	if recordType != "" {
		params.Add("type", recordType)
	}

	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	resp, err := s.request("GET", endpoint, config, nil)
	if err != nil {
		return nil, err
	}

	// 解析结果
	resultBytes, err := json.Marshal(resp.Result)
	if err != nil {
		return nil, fmt.Errorf("解析结果失败: %v", err)
	}

	var records []DNSRecord
	if err := json.Unmarshal(resultBytes, &records); err != nil {
		return nil, fmt.Errorf("解析 DNS 记录失败: %v", err)
	}

	return records, nil
}

// CreateDNSRecord 创建 DNS 记录
func (s *CloudflareService) CreateDNSRecord(config CloudflareConfig, record DNSRecord) (*DNSRecord, error) {
	endpoint := fmt.Sprintf("/zones/%s/dns_records", config.ZoneID)

	// 设置默认 TTL
	if record.TTL == 0 {
		if record.Proxied {
			record.TTL = 1 // 代理模式下 TTL 必须为 1
		} else {
			record.TTL = 3600
		}
	}

	resp, err := s.request("POST", endpoint, config, record)
	if err != nil {
		return nil, err
	}

	// 解析结果
	resultBytes, err := json.Marshal(resp.Result)
	if err != nil {
		return nil, fmt.Errorf("解析结果失败: %v", err)
	}

	var createdRecord DNSRecord
	if err := json.Unmarshal(resultBytes, &createdRecord); err != nil {
		return nil, fmt.Errorf("解析创建的记录失败: %v", err)
	}

	return &createdRecord, nil
}

// UpdateDNSRecord 更新 DNS 记录
func (s *CloudflareService) UpdateDNSRecord(config CloudflareConfig, recordID string, record DNSRecord) (*DNSRecord, error) {
	endpoint := fmt.Sprintf("/zones/%s/dns_records/%s", config.ZoneID, recordID)

	// 设置默认 TTL
	if record.TTL == 0 {
		if record.Proxied {
			record.TTL = 1
		} else {
			record.TTL = 3600
		}
	}

	resp, err := s.request("PUT", endpoint, config, record)
	if err != nil {
		return nil, err
	}

	// 解析结果
	resultBytes, err := json.Marshal(resp.Result)
	if err != nil {
		return nil, fmt.Errorf("解析结果失败: %v", err)
	}

	var updatedRecord DNSRecord
	if err := json.Unmarshal(resultBytes, &updatedRecord); err != nil {
		return nil, fmt.Errorf("解析更新的记录失败: %v", err)
	}

	return &updatedRecord, nil
}

// DeleteDNSRecord 删除 DNS 记录
func (s *CloudflareService) DeleteDNSRecord(config CloudflareConfig, recordID string) error {
	endpoint := fmt.Sprintf("/zones/%s/dns_records/%s", config.ZoneID, recordID)

	_, err := s.request("DELETE", endpoint, config, nil)
	return err
}

// ConfigureDNSRecord 配置 DNS 记录（创建或更新）
func (s *CloudflareService) ConfigureDNSRecord(config CloudflareConfig, record DNSRecord) (*DNSRecord, string, error) {
	// 查找现有记录
	existingRecords, err := s.GetDNSRecords(config, record.Name, record.Type)
	if err != nil {
		return nil, "", fmt.Errorf("查询现有记录失败: %v", err)
	}

	if len(existingRecords) > 0 {
		// 更新现有记录
		existingRecord := existingRecords[0]
		updatedRecord, err := s.UpdateDNSRecord(config, existingRecord.ID, record)
		if err != nil {
			return nil, "", fmt.Errorf("更新记录失败: %v", err)
		}
		return updatedRecord, "updated", nil
	} else {
		// 创建新记录
		createdRecord, err := s.CreateDNSRecord(config, record)
		if err != nil {
			return nil, "", fmt.Errorf("创建记录失败: %v", err)
		}
		return createdRecord, "created", nil
	}
}

// GetAccountID 获取账户ID（Pages API 需要）
func (s *CloudflareService) GetAccountID(config CloudflareConfig) (string, error) {
	resp, err := s.request("GET", "/accounts", config, nil)
	if err != nil {
		return "", err
	}

	// 解析结果
	resultBytes, err := json.Marshal(resp.Result)
	if err != nil {
		return "", fmt.Errorf("解析结果失败: %v", err)
	}

	var accounts []map[string]interface{}
	if err := json.Unmarshal(resultBytes, &accounts); err != nil {
		return "", fmt.Errorf("解析账户信息失败: %v", err)
	}

	if len(accounts) == 0 {
		return "", fmt.Errorf("未找到账户信息")
	}

	// 返回第一个账户的ID
	if accountID, ok := accounts[0]["id"].(string); ok {
		return accountID, nil
	}

	return "", fmt.Errorf("无法获取账户ID")
}

// GetPagesProjects 获取 Pages 项目列表
func (s *CloudflareService) GetPagesProjects(config CloudflareConfig, accountID string) ([]PagesProject, error) {
	endpoint := fmt.Sprintf("/accounts/%s/pages/projects", accountID)

	resp, err := s.request("GET", endpoint, config, nil)
	if err != nil {
		return nil, err
	}

	// 解析结果
	resultBytes, err := json.Marshal(resp.Result)
	if err != nil {
		return nil, fmt.Errorf("解析结果失败: %v", err)
	}

	var projects []PagesProject
	if err := json.Unmarshal(resultBytes, &projects); err != nil {
		return nil, fmt.Errorf("解析 Pages 项目失败: %v", err)
	}

	return projects, nil
}

// AddPagesCustomDomain 为 Pages 项目添加自定义域名
func (s *CloudflareService) AddPagesCustomDomain(config CloudflareConfig, accountID, projectName, domain string) (*PagesCustomDomain, error) {
	endpoint := fmt.Sprintf("/accounts/%s/pages/projects/%s/domains", accountID, projectName)

	domainData := map[string]string{
		"name": domain,
	}

	resp, err := s.request("POST", endpoint, config, domainData)
	if err != nil {
		return nil, err
	}

	// 解析结果
	resultBytes, err := json.Marshal(resp.Result)
	if err != nil {
		return nil, fmt.Errorf("解析结果失败: %v", err)
	}

	var customDomain PagesCustomDomain
	if err := json.Unmarshal(resultBytes, &customDomain); err != nil {
		return nil, fmt.Errorf("解析自定义域名失败: %v", err)
	}

	return &customDomain, nil
}

// GetPagesCustomDomains 获取 Pages 项目的自定义域名列表
func (s *CloudflareService) GetPagesCustomDomains(config CloudflareConfig, accountID, projectName string) ([]PagesCustomDomain, error) {
	endpoint := fmt.Sprintf("/accounts/%s/pages/projects/%s/domains", accountID, projectName)

	resp, err := s.request("GET", endpoint, config, nil)
	if err != nil {
		return nil, err
	}

	// 解析结果
	resultBytes, err := json.Marshal(resp.Result)
	if err != nil {
		return nil, fmt.Errorf("解析结果失败: %v", err)
	}

	var domains []PagesCustomDomain
	if err := json.Unmarshal(resultBytes, &domains); err != nil {
		return nil, fmt.Errorf("解析自定义域名失败: %v", err)
	}

	return domains, nil
}

// DeletePagesCustomDomain 删除 Pages 项目的自定义域名
func (s *CloudflareService) DeletePagesCustomDomain(config CloudflareConfig, accountID, projectName, domain string) error {
	endpoint := fmt.Sprintf("/accounts/%s/pages/projects/%s/domains/%s", accountID, projectName, domain)

	_, err := s.request("DELETE", endpoint, config, nil)
	return err
}
