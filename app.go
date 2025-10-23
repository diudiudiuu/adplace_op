package main

import (
	"adsplat/services"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/crypto/ssh"
)

// App struct
type App struct {
	ctx                context.Context
	jsonService        *services.JsonService
	aesService         *services.AesService
	kvService          *services.KvService
	cloudflareService  *services.CloudflareService
	pageCaptureService *services.PageCaptureService
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		jsonService:        services.NewJsonService(),
		aesService:         services.NewAesService(),
		kvService:          services.NewKvService(),
		cloudflareService:  services.NewCloudflareService(),
		pageCaptureService: services.NewPageCaptureService(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// ApiResponse 通用API响应结构
type ApiResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// List 获取服务器列表 (对应 Rust 的 list 函数)
func (a *App) List(authorization, clientJson string) string {
	log.Printf("List called with authorization: %s", authorization)

	// 检查授权
	if authorization == "" || strings.TrimSpace(authorization) == "" {
		log.Printf("No authorization provided")
		response := ApiResponse{Code: 401, Msg: "Authorization required"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	servers, kvResponse, err := a.jsonService.LoadJsonFileWithResponse(authorization, clientJson)
	if err != nil {
		log.Printf("Failed to load JSON file: %v", err)
		response := ApiResponse{Code: 500, Msg: "Internal server error"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	// 检查 KV 响应的状态码
	if kvResponse != nil && kvResponse.Code == 401 {
		log.Printf("KV service returned 401 Unauthorized")
		response := ApiResponse{Code: 401, Msg: "Unauthorized"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	log.Printf("Loaded %d servers from JSON", len(servers))

	result, err := json.Marshal(servers)
	if err != nil {
		log.Printf("Failed to marshal servers: %v", err)
		response := ApiResponse{Code: 500, Msg: "Failed to marshal data"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	return string(result)
}

// ServerInfo 获取服务器信息 (对应 Rust 的 server_info 函数)
func (a *App) ServerInfo(serverID, authorization, clientJson string) string {
	log.Printf("ServerInfo called with serverID: %s, authorization: %s", serverID, authorization)

	// 检查授权
	if authorization == "" || strings.TrimSpace(authorization) == "" {
		response := ApiResponse{Code: 401, Msg: "Authorization required"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	server, err := a.jsonService.GetServerByID(serverID, authorization, clientJson)
	if err != nil {
		log.Printf("Failed to get server info: %v", err)
		response := ApiResponse{Code: 500, Msg: "Internal server error"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	if server == nil {
		response := ApiResponse{Code: 404, Msg: "Server not found"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	result, err := json.Marshal(server)
	if err != nil {
		log.Printf("Failed to marshal server: %v", err)
		response := ApiResponse{Code: 500, Msg: "Failed to marshal data"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	return string(result)
}

// ProjectInfo 获取项目信息 (对应 Rust 的 project_info 函数)
func (a *App) ProjectInfo(projectID, authorization, clientJson string) string {
	log.Printf("ProjectInfo called with projectID: %s, authorization: %s", projectID, authorization)

	// 检查授权
	if authorization == "" || strings.TrimSpace(authorization) == "" {
		response := ApiResponse{Code: 401, Msg: "Authorization required"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	project, err := a.jsonService.GetProjectByID(projectID, authorization, clientJson)
	if err != nil {
		log.Printf("Failed to get project info: %v", err)
		response := ApiResponse{Code: 500, Msg: "Internal server error"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	if project == nil {
		response := ApiResponse{Code: 404, Msg: "Project not found"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	result, err := json.Marshal(project)
	if err != nil {
		log.Printf("Failed to marshal project: %v", err)
		response := ApiResponse{Code: 500, Msg: "Failed to marshal data"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	return string(result)
}

// ProjectForm 创建或更新项目 (对应 Rust 的 project_form 函数)
func (a *App) ProjectForm(serverID, projectInfo, authorization, clientJson string) string {
	log.Printf("ProjectForm called with serverID: %s, projectInfo: %s, authorization: %s",
		serverID, projectInfo, authorization)

	// 检查授权
	if authorization == "" || strings.TrimSpace(authorization) == "" {
		response := ApiResponse{Code: 401, Msg: "Authorization required"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	var project services.ProjectData
	if err := json.Unmarshal([]byte(projectInfo), &project); err != nil {
		log.Printf("Failed to unmarshal project info: %v", err)
		response := ApiResponse{Code: 400, Msg: "Invalid project data"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	err := a.jsonService.AddOrUpdateProject(serverID, project, authorization, clientJson)
	if err != nil {
		log.Printf("Failed to add/update project: %v", err)
		response := ApiResponse{Code: 500, Msg: err.Error()}
		result, _ := json.Marshal(response)
		return string(result)
	}

	response := ApiResponse{Code: 200, Msg: "Success"}
	result, _ := json.Marshal(response)
	return string(result)
}

// ProjectDelete 删除项目 (对应 Rust 的 project_delete 函数)
func (a *App) ProjectDelete(serverID, projectID, authorization, clientJson string) string {
	log.Printf("ProjectDelete called with serverID: %s, projectID: %s, authorization: %s",
		serverID, projectID, authorization)

	// 检查授权
	if authorization == "" || strings.TrimSpace(authorization) == "" {
		response := ApiResponse{Code: 401, Msg: "Authorization required"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	err := a.jsonService.DeleteProject(serverID, projectID, authorization, clientJson)
	if err != nil {
		log.Printf("Failed to delete project: %v", err)
		response := ApiResponse{Code: 500, Msg: "Failed to delete project"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	response := ApiResponse{Code: 200, Msg: "Success"}
	result, _ := json.Marshal(response)
	return string(result)
}

// Exec 执行SQL (对应 Rust 的 exec 函数)
func (a *App) Exec(projectID, sql, sqlType, authorization, clientJson string) string {
	log.Printf("Exec called with projectID: %s, sql: %s, sqlType: %s, authorization: %s",
		projectID, sql, sqlType, authorization)

	// 获取项目信息
	project, err := a.jsonService.GetProjectByID(projectID, authorization, clientJson)
	if err != nil {
		log.Printf("Failed to get project info: %v", err)
		return `{"code": 500, "msg": "Failed to get project info"}`
	}

	if project == nil {
		return `{"code": 404, "msg": "Project not found"}`
	}

	// 加密 SQL
	signature, err := a.aesService.Encrypt(sql)
	if err != nil {
		log.Printf("Failed to encrypt SQL: %v", err)
		return `{"code": 500, "msg": "Failed to encrypt SQL"}`
	}

	// 构建请求URL
	apiURL := fmt.Sprintf("%s/dbexec", project.ProjectAPIURL)

	// 准备表单数据
	formData := url.Values{}
	formData.Set("sql_type", sqlType)
	formData.Set("signature", signature)

	// 发送POST请求
	resp, err := http.PostForm(apiURL, formData)
	if err != nil {
		log.Printf("Failed to send request: %v", err)
		return fmt.Sprintf(`{"code": 500, "msg": "%s"}`, err.Error())
	}
	defer resp.Body.Close()

	// 读取响应
	body := make([]byte, 0)
	buffer := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			body = append(body, buffer[:n]...)
		}
		if err != nil {
			break
		}
	}

	responseBody := string(body)
	log.Printf("API response status: %d, body: %s", resp.StatusCode, responseBody)

	// 检查HTTP状态码
	if resp.StatusCode != 200 {
		log.Printf("API request failed with status: %d", resp.StatusCode)
		response := ApiResponse{
			Code: resp.StatusCode,
			Msg:  fmt.Sprintf("API请求失败，状态码: %d", resp.StatusCode),
			Data: responseBody,
		}
		result, _ := json.Marshal(response)
		return string(result)
	}

	// 尝试验证响应是否为有效JSON
	var testJson interface{}
	if err := json.Unmarshal(body, &testJson); err != nil {
		log.Printf("API response is not valid JSON: %v", err)
		// 如果不是有效JSON，包装成标准响应格式
		response := ApiResponse{
			Code: 200,
			Msg:  "Success",
			Data: responseBody,
		}
		result, _ := json.Marshal(response)
		return string(result)
	}

	// 如果是有效JSON，直接返回
	return responseBody
}

// ShowMessage 显示消息对话框
func (a *App) ShowMessage(title, message string) {
	runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:    runtime.InfoDialog,
		Title:   title,
		Message: message,
	})
}

// Greet returns a greeting for the given name (保留原有方法)
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// ServerAdd 添加新服务器
func (a *App) ServerAdd(serverID, serverName, serverIP, serverPort, serverUser, serverPassword, defaultPath, authorization, clientJson string) string {
	log.Printf("ServerAdd called with serverID: %s, serverName: %s, defaultPath: %s, authorization: %s", serverID, serverName, defaultPath, authorization)

	// 检查授权
	if authorization == "" || strings.TrimSpace(authorization) == "" {
		response := ApiResponse{Code: 401, Msg: "Authorization required"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	// 如果 defaultPath 为空，设置默认值
	if defaultPath == "" {
		defaultPath = "/adplace"
	}

	// 创建新服务器数据
	newServer := services.ServerData{
		ServerID:       serverID,
		ServerName:     serverName,
		ServerIP:       serverIP,
		ServerPort:     serverPort,
		ServerUser:     serverUser,
		ServerPassword: serverPassword,
		DefaultPath:    defaultPath,
		ProjectList:    []services.ProjectData{},
	}

	log.Printf("Created newServer with DefaultPath: %s", newServer.DefaultPath)

	err := a.jsonService.AddServer(newServer, authorization, clientJson)
	if err != nil {
		log.Printf("Failed to add server: %v", err)
		response := ApiResponse{Code: 500, Msg: err.Error()}
		result, _ := json.Marshal(response)
		return string(result)
	}

	response := ApiResponse{Code: 200, Msg: "Server added successfully"}
	result, _ := json.Marshal(response)
	return string(result)
}

// ServerUpdate 更新服务器信息
func (a *App) ServerUpdate(oldServerID, newServerID, serverName, serverIP, serverPort, serverUser, serverPassword, defaultPath, authorization, clientJson string) string {
	log.Printf("ServerUpdate called with oldServerID: %s, newServerID: %s, serverName: %s, defaultPath: %s, authorization: %s", oldServerID, newServerID, serverName, defaultPath, authorization)

	// 检查授权
	if authorization == "" || strings.TrimSpace(authorization) == "" {
		response := ApiResponse{Code: 401, Msg: "Authorization required"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	// 如果 defaultPath 为空，设置默认值
	if defaultPath == "" {
		defaultPath = "/adplace"
	}

	// 创建更新的服务器数据
	updatedServer := services.ServerData{
		ServerID:       newServerID,
		ServerName:     serverName,
		ServerIP:       serverIP,
		ServerPort:     serverPort,
		ServerUser:     serverUser,
		ServerPassword: serverPassword,
		DefaultPath:    defaultPath,
	}

	log.Printf("Created updatedServer with DefaultPath: %s", updatedServer.DefaultPath)

	err := a.jsonService.UpdateServerWithNewID(oldServerID, updatedServer, authorization, clientJson)
	if err != nil {
		log.Printf("Failed to update server: %v", err)
		response := ApiResponse{Code: 500, Msg: err.Error()}
		result, _ := json.Marshal(response)
		return string(result)
	}

	response := ApiResponse{Code: 200, Msg: "Server updated successfully"}
	result, _ := json.Marshal(response)
	return string(result)
}

// ServerDelete 删除服务器
func (a *App) ServerDelete(serverID, authorization, clientJson string) string {
	log.Printf("ServerDelete called with serverID: %s, authorization: %s", serverID, authorization)

	// 检查授权
	if authorization == "" || strings.TrimSpace(authorization) == "" {
		response := ApiResponse{Code: 401, Msg: "Authorization required"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	err := a.jsonService.DeleteServer(serverID, authorization, clientJson)
	if err != nil {
		log.Printf("Failed to delete server: %v", err)
		response := ApiResponse{Code: 500, Msg: "Failed to delete server"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	response := ApiResponse{Code: 200, Msg: "Server deleted successfully"}
	result, _ := json.Marshal(response)
	return string(result)
}

// TestStoredServerSSH 测试已存储服务器的SSH连接
func (a *App) TestStoredServerSSH(serverID, authorization, clientJson string) string {
	log.Printf("TestStoredServerSSH called with serverID: %s", serverID)

	// 检查授权
	if authorization == "" || strings.TrimSpace(authorization) == "" {
		response := ApiResponse{Code: 401, Msg: "Authorization required"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	// 获取服务器信息
	server, err := a.jsonService.GetServerByID(serverID, authorization, clientJson)
	if err != nil {
		log.Printf("Failed to get server info: %v", err)
		response := ApiResponse{Code: 500, Msg: "获取服务器信息失败"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	if server == nil {
		response := ApiResponse{Code: 404, Msg: "服务器不存在"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	// 执行SSH测试
	testResult := a.performSSHTest(server.ServerIP, server.ServerPort, server.ServerUser, server.ServerPassword)

	// 更新服务器的连接状态
	err = a.jsonService.UpdateServerConnectionStatus(serverID, testResult, authorization, clientJson)
	if err != nil {
		log.Printf("Failed to update server connection status: %v", err)
	}

	return testResult
}

// performSSHTest 执行SSH测试的核心逻辑
func (a *App) performSSHTest(serverIP, serverPort, serverUser, serverPassword string) string {
	log.Printf("performSSHTest called with serverIP: %s, serverPort: %s, serverUser: %s", serverIP, serverPort, serverUser)

	// 检查必要参数
	if serverIP == "" {
		response := ApiResponse{Code: 400, Msg: "服务器IP不能为空"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	// 设置默认端口
	if serverPort == "" {
		serverPort = "22"
	}

	// 设置默认用户名
	if serverUser == "" {
		serverUser = "root"
	}

	// 创建SSH配置
	config := &ssh.ClientConfig{
		User:            serverUser,
		Auth:            []ssh.AuthMethod{},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 注意：生产环境应该验证主机密钥
		Timeout:         10 * time.Second,
	}

	// 添加密码认证（如果提供了密码）
	if serverPassword != "" {
		config.Auth = append(config.Auth, ssh.Password(serverPassword))
	}

	// 尝试连接
	address := fmt.Sprintf("%s:%s", serverIP, serverPort)
	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		log.Printf("SSH connection failed: %v", err)
		response := ApiResponse{
			Code: 500,
			Msg:  fmt.Sprintf("SSH连接失败: %v", err),
			Data: map[string]interface{}{
				"connected": false,
				"error":     err.Error(),
				"test_time": time.Now().Format("2006-01-02 15:04:05"),
			},
		}
		result, _ := json.Marshal(response)
		return string(result)
	}
	defer client.Close()

	// 测试执行简单命令
	session, err := client.NewSession()
	if err != nil {
		log.Printf("Failed to create SSH session: %v", err)
		response := ApiResponse{
			Code: 500,
			Msg:  fmt.Sprintf("创建SSH会话失败: %v", err),
			Data: map[string]interface{}{
				"connected":     true,
				"session_error": err.Error(),
				"test_time":     time.Now().Format("2006-01-02 15:04:05"),
			},
		}
		result, _ := json.Marshal(response)
		return string(result)
	}
	defer session.Close()

	// 执行简单的whoami命令测试
	output, err := session.Output("whoami")
	if err != nil {
		log.Printf("Failed to execute test command: %v", err)
	}

	response := ApiResponse{
		Code: 200,
		Msg:  "SSH连接成功",
		Data: map[string]interface{}{
			"connected": true,
			"user":      strings.TrimSpace(string(output)),
			"test_time": time.Now().Format("2006-01-02 15:04:05"),
		},
	}
	result, _ := json.Marshal(response)
	return string(result)
}

// TestSSHConnection 测试SSH连接
func (a *App) TestSSHConnection(serverIP, serverPort, serverUser, serverPassword string) string {
	return a.performSSHTest(serverIP, serverPort, serverUser, serverPassword)
}

// TestUnauthorized 测试 401 响应 (用于测试前端的 401 处理)
func (a *App) TestUnauthorized() string {
	response := ApiResponse{Code: 401, Msg: "Test unauthorized response"}
	result, _ := json.Marshal(response)
	return string(result)
}

// CloudflareGetDNSRecords 获取 Cloudflare DNS 记录
func (a *App) CloudflareGetDNSRecords(apiToken, zoneID, name, recordType string) string {
	log.Printf("CloudflareGetDNSRecords called with name: %s, type: %s", name, recordType)

	config := services.CloudflareConfig{
		APIToken: apiToken,
		ZoneID:   zoneID,
	}

	records, err := a.cloudflareService.GetDNSRecords(config, name, recordType)
	if err != nil {
		log.Printf("Failed to get DNS records: %v", err)
		response := ApiResponse{Code: 500, Msg: fmt.Sprintf("获取 DNS 记录失败: %v", err)}
		result, _ := json.Marshal(response)
		return string(result)
	}

	response := ApiResponse{Code: 200, Msg: "Success", Data: records}
	result, _ := json.Marshal(response)
	return string(result)
}

// CloudflareConfigureDNSRecord 配置 Cloudflare DNS 记录
func (a *App) CloudflareConfigureDNSRecord(apiToken, zoneID, name, recordType, content string, proxied bool) string {
	log.Printf("CloudflareConfigureDNSRecord called with name: %s, type: %s, content: %s, proxied: %t",
		name, recordType, content, proxied)

	config := services.CloudflareConfig{
		APIToken: apiToken,
		ZoneID:   zoneID,
	}

	record := services.DNSRecord{
		Type:    recordType,
		Name:    name,
		Content: content,
		Proxied: proxied,
	}

	result, action, err := a.cloudflareService.ConfigureDNSRecord(config, record)
	if err != nil {
		log.Printf("Failed to configure DNS record: %v", err)
		response := ApiResponse{Code: 500, Msg: fmt.Sprintf("配置 DNS 记录失败: %v", err)}
		responseResult, _ := json.Marshal(response)
		return string(responseResult)
	}

	responseData := map[string]interface{}{
		"record": result,
		"action": action,
	}

	response := ApiResponse{Code: 200, Msg: fmt.Sprintf("DNS 记录%s成功",
		map[string]string{"created": "创建", "updated": "更新"}[action]), Data: responseData}
	responseResult, _ := json.Marshal(response)
	return string(responseResult)
}

// CloudflareDeleteDNSRecord 删除 Cloudflare DNS 记录
func (a *App) CloudflareDeleteDNSRecord(apiToken, zoneID, recordID string) string {
	log.Printf("CloudflareDeleteDNSRecord called with recordID: %s", recordID)

	config := services.CloudflareConfig{
		APIToken: apiToken,
		ZoneID:   zoneID,
	}

	err := a.cloudflareService.DeleteDNSRecord(config, recordID)
	if err != nil {
		log.Printf("Failed to delete DNS record: %v", err)
		response := ApiResponse{Code: 500, Msg: fmt.Sprintf("删除 DNS 记录失败: %v", err)}
		result, _ := json.Marshal(response)
		return string(result)
	}

	response := ApiResponse{Code: 200, Msg: "DNS 记录删除成功"}
	result, _ := json.Marshal(response)
	return string(result)
}

// CloudflareBatchConfigureDNS 批量配置 DNS 记录
func (a *App) CloudflareBatchConfigureDNS(apiToken, zoneID, recordsJson string) string {
	log.Printf("CloudflareBatchConfigureDNS called with records: %s", recordsJson)

	config := services.CloudflareConfig{
		APIToken: apiToken,
		ZoneID:   zoneID,
	}

	var records []services.DNSRecord
	if err := json.Unmarshal([]byte(recordsJson), &records); err != nil {
		log.Printf("Failed to unmarshal records: %v", err)
		response := ApiResponse{Code: 400, Msg: "记录数据格式错误"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	results := make([]map[string]interface{}, 0, len(records))

	for _, record := range records {
		result, action, err := a.cloudflareService.ConfigureDNSRecord(config, record)
		if err != nil {
			log.Printf("Failed to configure DNS record %s: %v", record.Name, err)
			results = append(results, map[string]interface{}{
				"name":   record.Name,
				"type":   record.Type,
				"error":  err.Error(),
				"action": "failed",
			})
		} else {
			results = append(results, map[string]interface{}{
				"record": result,
				"action": action,
			})
		}
	}

	response := ApiResponse{Code: 200, Msg: "批量配置完成", Data: results}
	responseResult, _ := json.Marshal(response)
	return string(responseResult)
}

// ProjectPortUpdate 更新项目端口信息
func (a *App) ProjectPortUpdate(projectID, apiPort, frontPort, authorization, clientJson string) string {
	log.Printf("ProjectPortUpdate called with projectID: %s, apiPort: %s, frontPort: %s, authorization: %s",
		projectID, apiPort, frontPort, authorization)

	// 检查授权
	if authorization == "" || strings.TrimSpace(authorization) == "" {
		response := ApiResponse{Code: 401, Msg: "Authorization required"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	// 获取项目信息
	project, err := a.jsonService.GetProjectByID(projectID, authorization, clientJson)
	if err != nil {
		log.Printf("Failed to get project info: %v", err)
		response := ApiResponse{Code: 500, Msg: "Internal server error"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	if project == nil {
		response := ApiResponse{Code: 404, Msg: "Project not found"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	// 更新端口信息
	if apiPort != "" {
		project.APIPort = apiPort
	}
	if frontPort != "" {
		project.FrontPort = frontPort
	}

	// 找到项目所属的服务器并更新
	servers, err := a.jsonService.LoadJsonFile(authorization, clientJson)
	if err != nil {
		log.Printf("Failed to load servers: %v", err)
		response := ApiResponse{Code: 500, Msg: "Internal server error"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	for _, server := range servers {
		for _, proj := range server.ProjectList {
			if proj.ProjectID == projectID {
				err := a.jsonService.AddOrUpdateProject(server.ServerID, *project, authorization, clientJson)
				if err != nil {
					log.Printf("Failed to update project: %v", err)
					response := ApiResponse{Code: 500, Msg: err.Error()}
					result, _ := json.Marshal(response)
					return string(result)
				}

				response := ApiResponse{Code: 200, Msg: "Project ports updated successfully"}
				result, _ := json.Marshal(response)
				return string(result)
			}
		}
	}

	response := ApiResponse{Code: 404, Msg: "Project not found in any server"}
	result, _ := json.Marshal(response)
	return string(result)
}

// CloudflarePagesAddDomain 为 Cloudflare Pages 项目添加自定义域名
func (a *App) CloudflarePagesAddDomain(apiToken, zoneID, projectName, domain string) string {
	log.Printf("CloudflarePagesAddDomain called with projectName: %s, domain: %s", projectName, domain)

	config := services.CloudflareConfig{
		APIToken: apiToken,
		ZoneID:   zoneID,
	}

	// 获取账户ID
	log.Printf("Getting account ID...")
	accountID, err := a.cloudflareService.GetAccountID(config)
	if err != nil {
		log.Printf("Failed to get account ID via /accounts API: %v", err)
		log.Printf("Trying backup method via zone info...")

		// 尝试备用方案：从Zone信息获取账户ID
		accountID, err = a.cloudflareService.GetAccountIDFromZone(config)
		if err != nil {
			log.Printf("Failed to get account ID via zone info: %v", err)
			response := ApiResponse{Code: 500, Msg: fmt.Sprintf("获取账户ID失败: %v。请确保API Token有Account:Read权限或Zone:Read权限", err)}
			result, _ := json.Marshal(response)
			return string(result)
		}
	}
	log.Printf("Got account ID: %s", accountID)

	// 检查域名是否已存在
	log.Printf("Checking existing domains for project: %s", projectName)
	existingDomains, err := a.cloudflareService.GetPagesCustomDomains(config, accountID, projectName)
	if err != nil {
		log.Printf("Failed to get existing domains: %v", err)
		// 不阻断流程，继续尝试添加
	} else {
		log.Printf("Found %d existing domains", len(existingDomains))
		// 检查域名是否已存在
		for _, existingDomain := range existingDomains {
			log.Printf("Existing domain: %s", existingDomain.Name)
			if existingDomain.Name == domain {
				log.Printf("Domain %s already exists for project %s", domain, projectName)
				response := ApiResponse{
					Code: 200,
					Msg:  "域名已存在",
					Data: map[string]interface{}{
						"domain": existingDomain,
						"action": "exists",
					},
				}
				result, _ := json.Marshal(response)
				return string(result)
			}
		}
	}

	// 添加自定义域名
	log.Printf("Adding custom domain %s to project %s", domain, projectName)
	customDomain, err := a.cloudflareService.AddPagesCustomDomain(config, accountID, projectName, domain)
	if err != nil {
		log.Printf("Failed to add custom domain: %v", err)
		response := ApiResponse{Code: 500, Msg: fmt.Sprintf("添加自定义域名失败: %v", err)}
		result, _ := json.Marshal(response)
		return string(result)
	}
	log.Printf("Successfully added custom domain: %+v", customDomain)

	responseData := map[string]interface{}{
		"domain": customDomain,
		"action": "created",
	}

	response := ApiResponse{Code: 200, Msg: "自定义域名添加成功", Data: responseData}
	result, _ := json.Marshal(response)
	return string(result)
}

// CloudflarePagesGetDomains 获取 Cloudflare Pages 项目的自定义域名列表
func (a *App) CloudflarePagesGetDomains(apiToken, zoneID, projectName string) string {
	log.Printf("CloudflarePagesGetDomains called with projectName: %s", projectName)

	config := services.CloudflareConfig{
		APIToken: apiToken,
		ZoneID:   zoneID,
	}

	// 获取账户ID
	accountID, err := a.cloudflareService.GetAccountID(config)
	if err != nil {
		log.Printf("Failed to get account ID via /accounts API: %v", err)
		// 尝试备用方案
		accountID, err = a.cloudflareService.GetAccountIDFromZone(config)
		if err != nil {
			log.Printf("Failed to get account ID via zone info: %v", err)
			response := ApiResponse{Code: 500, Msg: fmt.Sprintf("获取账户ID失败: %v", err)}
			result, _ := json.Marshal(response)
			return string(result)
		}
	}

	// 获取自定义域名列表
	domains, err := a.cloudflareService.GetPagesCustomDomains(config, accountID, projectName)
	if err != nil {
		log.Printf("Failed to get custom domains: %v", err)
		response := ApiResponse{Code: 500, Msg: fmt.Sprintf("获取自定义域名失败: %v", err)}
		result, _ := json.Marshal(response)
		return string(result)
	}

	response := ApiResponse{Code: 200, Msg: "Success", Data: domains}
	result, _ := json.Marshal(response)
	return string(result)
}

// CloudflarePagesDeleteDomain 删除 Cloudflare Pages 项目的自定义域名
func (a *App) CloudflarePagesDeleteDomain(apiToken, zoneID, projectName, domain string) string {
	log.Printf("CloudflarePagesDeleteDomain called with projectName: %s, domain: %s", projectName, domain)

	config := services.CloudflareConfig{
		APIToken: apiToken,
		ZoneID:   zoneID,
	}

	// 获取账户ID
	accountID, err := a.cloudflareService.GetAccountID(config)
	if err != nil {
		log.Printf("Failed to get account ID via /accounts API: %v", err)
		// 尝试备用方案
		accountID, err = a.cloudflareService.GetAccountIDFromZone(config)
		if err != nil {
			log.Printf("Failed to get account ID via zone info: %v", err)
			response := ApiResponse{Code: 500, Msg: fmt.Sprintf("获取账户ID失败: %v", err)}
			result, _ := json.Marshal(response)
			return string(result)
		}
	}

	// 删除自定义域名
	err = a.cloudflareService.DeletePagesCustomDomain(config, accountID, projectName, domain)
	if err != nil {
		log.Printf("Failed to delete custom domain: %v", err)
		response := ApiResponse{Code: 500, Msg: fmt.Sprintf("删除自定义域名失败: %v", err)}
		result, _ := json.Marshal(response)
		return string(result)
	}

	response := ApiResponse{Code: 200, Msg: "自定义域名删除成功"}
	result, _ := json.Marshal(response)
	return string(result)
}

// GenerateProjectConfig 生成项目配置文件并上传到服务器（从数据库获取最新数据）
func (a *App) GenerateProjectConfig(serverID, authorization, clientJson string) string {
	log.Printf("GenerateProjectConfig called with serverID: %s", serverID)

	// 检查授权
	if authorization == "" || strings.TrimSpace(authorization) == "" {
		response := ApiResponse{Code: 401, Msg: "Authorization required"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	// 获取服务器信息
	server, err := a.jsonService.GetServerByID(serverID, authorization, clientJson)
	if err != nil {
		log.Printf("Failed to get server info: %v", err)
		response := ApiResponse{Code: 500, Msg: "获取服务器信息失败"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	if server == nil {
		response := ApiResponse{Code: 404, Msg: "服务器不存在"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	// 生成项目配置 - 使用当前项目配置逻辑
	projectConfig := make(map[string]map[string]string)
	for _, project := range server.ProjectList {
		apiDomain := extractDomainFromURL(project.ProjectAPIURL)
		projectConfig[project.ProjectID] = map[string]string{
			"api_port":   getPortOrDefault(project.APIPort, "9000"),
			"web_port":   getPortOrDefault(project.FrontPort, "3000"),
			"api_domain": apiDomain,
		}
	}

	configJSON, err := json.MarshalIndent(projectConfig, "", "  ")
	if err != nil {
		log.Printf("Failed to generate project config: %v", err)
		response := ApiResponse{Code: 500, Msg: "生成配置文件失败"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	// 处理 release.zip 并上传配置文件
	err = a.processReleaseAndUploadConfig(server, "project_config.json", string(configJSON))
	if err != nil {
		log.Printf("Failed to process release and upload config: %v", err)
		response := ApiResponse{Code: 500, Msg: fmt.Sprintf("处理发布包和上传配置文件失败: %v", err)}
		result, _ := json.Marshal(response)
		return string(result)
	}

	response := ApiResponse{
		Code: 200,
		Msg:  "项目配置文件生成并上传成功",
		Data: map[string]interface{}{
			"config": projectConfig,
			"path":   fmt.Sprintf("%s/project_config.json", server.DefaultPath),
		},
	}
	result, _ := json.Marshal(response)
	return string(result)
}

// UploadProjectConfig 直接上传前端生成的项目配置JSON到服务器
func (a *App) UploadProjectConfig(serverDataJson, projectConfigJson, authorization string) string {
	log.Printf("UploadProjectConfig called")

	// 检查授权
	if authorization == "" || strings.TrimSpace(authorization) == "" {
		response := ApiResponse{Code: 401, Msg: "Authorization required"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	// 解析前端传入的服务器数据（只需要服务器连接信息）
	var server services.ServerData
	if err := json.Unmarshal([]byte(serverDataJson), &server); err != nil {
		log.Printf("Failed to unmarshal server data: %v", err)
		response := ApiResponse{Code: 400, Msg: "服务器数据格式错误"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	log.Printf("Uploading frontend-generated config JSON: %s", projectConfigJson)

	// 直接使用前端传入的JSON配置，处理 release.zip 并上传配置文件
	err := a.processReleaseAndUploadConfig(&server, "project_config.json", projectConfigJson)
	if err != nil {
		log.Printf("Failed to process release and upload config: %v", err)
		response := ApiResponse{Code: 500, Msg: fmt.Sprintf("处理发布包和上传配置文件失败: %v", err)}
		result, _ := json.Marshal(response)
		return string(result)
	}

	// 解析配置JSON以返回给前端预览
	var configMap map[string]interface{}
	json.Unmarshal([]byte(projectConfigJson), &configMap)

	response := ApiResponse{
		Code: 200,
		Msg:  "项目配置文件上传成功",
		Data: map[string]interface{}{
			"config":      configMap,
			"path":        fmt.Sprintf("%s/project_config.json", server.DefaultPath),
			"data_source": "frontend",
		},
	}
	result, _ := json.Marshal(response)
	return string(result)
}

// GenerateProjectConfigForSingleProject 上传前端生成的项目配置（兼容绑定文件）
// 注意：由于绑定文件的限制，参数名可能显示为serverID, projectID等，但实际用途如下：
// 第1个参数：服务器数据JSON
// 第2个参数：项目配置JSON
// 第3个参数：未使用
// 第4个参数：授权token
func (a *App) GenerateProjectConfigForSingleProject(serverDataJson, projectConfigJson, unused, authorization string) string {
	log.Printf("GenerateProjectConfigForSingleProject called - uploading frontend config")
	log.Printf("Received parameters:")
	log.Printf("  serverDataJson length: %d", len(serverDataJson))
	log.Printf("  serverDataJson first 100 chars: %s", func() string {
		if len(serverDataJson) > 100 {
			return serverDataJson[:100] + "..."
		}
		return serverDataJson
	}())
	log.Printf("  projectConfigJson length: %d", len(projectConfigJson))
	log.Printf("  projectConfigJson: %s", projectConfigJson)
	log.Printf("  unused: %s", unused)
	log.Printf("  authorization length: %d", len(authorization))

	// 检查授权
	if authorization == "" || strings.TrimSpace(authorization) == "" {
		log.Printf("Authorization is empty or whitespace")
		response := ApiResponse{Code: 401, Msg: "Authorization required"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	log.Printf("Attempting to parse server data JSON...")
	log.Printf("Server data JSON (first 200 chars): %s", func() string {
		if len(serverDataJson) > 200 {
			return serverDataJson[:200] + "..."
		}
		return serverDataJson
	}())

	// 先尝试解析为通用map，然后提取需要的字段
	var rawData map[string]interface{}
	if err := json.Unmarshal([]byte(serverDataJson), &rawData); err != nil {
		log.Printf("Failed to unmarshal as map: %v", err)
		log.Printf("Raw server data length: %d", len(serverDataJson))
		log.Printf("Raw server data (first 500 chars): %s", func() string {
			if len(serverDataJson) > 500 {
				return serverDataJson[:500] + "..."
			}
			return serverDataJson
		}())
		response := ApiResponse{Code: 400, Msg: fmt.Sprintf("服务器数据格式错误: %v", err)}
		result, _ := json.Marshal(response)
		return string(result)
	}

	// 手动构建ServerData结构
	server := services.ServerData{
		ServerID:       getStringFromMap(rawData, "server_id"),
		ServerName:     getStringFromMap(rawData, "server_name"),
		ServerIP:       getStringFromMap(rawData, "server_ip"),
		ServerPort:     getStringFromMap(rawData, "server_port"),
		ServerUser:     getStringFromMap(rawData, "server_user"),
		ServerPassword: getStringFromMap(rawData, "server_password"),
		DefaultPath:    getStringFromMap(rawData, "default_path"),
	}

	// 如果关键字段为空，设置默认值
	if server.DefaultPath == "" {
		server.DefaultPath = "/adplace"
	}
	if server.ServerPort == "" {
		server.ServerPort = "22"
	}

	log.Printf("Successfully parsed server data: ServerID=%s, IP=%s, DefaultPath=%s",
		server.ServerID, server.ServerIP, server.DefaultPath)

	log.Printf("Uploading frontend-generated config to server: %s", server.ServerID)

	// 直接使用前端传入的JSON配置，处理 release.zip 并上传配置文件
	err := a.processReleaseAndUploadConfig(&server, "project_config.json", projectConfigJson)
	if err != nil {
		log.Printf("Failed to process release and upload config: %v", err)
		response := ApiResponse{Code: 500, Msg: fmt.Sprintf("处理发布包和上传配置文件失败: %v", err)}
		result, _ := json.Marshal(response)
		return string(result)
	}

	// 解析配置JSON以返回给前端预览
	var configMap map[string]interface{}
	json.Unmarshal([]byte(projectConfigJson), &configMap)

	response := ApiResponse{
		Code: 200,
		Msg:  "项目配置文件上传成功（前端生成）",
		Data: map[string]interface{}{
			"config":      configMap,
			"path":        fmt.Sprintf("%s/project_config.json", server.DefaultPath),
			"data_source": "frontend",
		},
	}
	result, _ := json.Marshal(response)
	return string(result)
}

// generateCurrentProjectConfigJSON 生成当前项目的配置JSON
func (a *App) generateCurrentProjectConfigJSON(server *services.ServerData, projectID string) (map[string]map[string]string, string, error) {
	projectConfig := make(map[string]map[string]string)

	log.Printf("Generating config for current project: %s", projectID)

	// 查找当前项目
	var currentProject *services.ProjectData
	for _, project := range server.ProjectList {
		if project.ProjectID == projectID {
			currentProject = &project
			break
		}
	}

	if currentProject == nil {
		return nil, "", fmt.Errorf("project %s not found in server %s", projectID, server.ServerID)
	}

	// 提取API域名
	apiDomain := extractDomainFromURL(currentProject.ProjectAPIURL)

	log.Printf("Processing current project %s: API=%s, APIPort=%s, FrontPort=%s",
		currentProject.ProjectID, apiDomain, currentProject.APIPort, currentProject.FrontPort)

	// 只包含您指定的三个字段
	projectConfig[currentProject.ProjectID] = map[string]string{
		"api_port":   getPortOrDefault(currentProject.APIPort, "9000"),
		"web_port":   getPortOrDefault(currentProject.FrontPort, "3000"),
		"api_domain": apiDomain,
	}

	// 转换为JSON
	configJSON, err := json.MarshalIndent(projectConfig, "", "  ")
	if err != nil {
		return nil, "", fmt.Errorf("failed to marshal project config: %v", err)
	}

	log.Printf("Generated current project config JSON: %s", string(configJSON))
	log.Printf("Config contains %d projects", len(projectConfig))
	for id := range projectConfig {
		log.Printf("Config project ID: %s", id)
	}
	return projectConfig, string(configJSON), nil
}

// ProjectInit SSH执行项目初始化
func (a *App) ProjectInit(serverID, projectID, authorization, clientJson string) string {
	log.Printf("ProjectInit called with serverID: %s, projectID: %s", serverID, projectID)

	// 检查授权
	if authorization == "" || strings.TrimSpace(authorization) == "" {
		response := ApiResponse{Code: 401, Msg: "Authorization required"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	// 获取服务器信息
	server, err := a.jsonService.GetServerByID(serverID, authorization, clientJson)
	if err != nil {
		log.Printf("Failed to get server info: %v", err)
		response := ApiResponse{Code: 500, Msg: "获取服务器信息失败"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	if server == nil {
		response := ApiResponse{Code: 404, Msg: "服务器不存在"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	// 执行SSH命令 - 先设置可执行权限，然后执行脚本
	command := fmt.Sprintf("cd %s && chmod +x codedeploy.sh && ./codedeploy.sh init %s", server.DefaultPath, projectID)
	output, err := a.executeSSHCommand(server, command)
	if err != nil {
		log.Printf("Failed to execute init command: %v", err)
		response := ApiResponse{Code: 500, Msg: fmt.Sprintf("执行初始化命令失败: %v", err)}
		result, _ := json.Marshal(response)
		return string(result)
	}

	response := ApiResponse{
		Code: 200,
		Msg:  "项目初始化成功",
		Data: map[string]interface{}{
			"command": command,
			"output":  output,
		},
	}
	result, _ := json.Marshal(response)
	return string(result)
}

// ProjectUpdate SSH执行项目更新（从数据库获取最新数据）
func (a *App) ProjectUpdate(serverID, projectID, authorization, clientJson string) string {
	log.Printf("ProjectUpdate called with serverID: %s, projectID: %s", serverID, projectID)

	// 检查授权
	if authorization == "" || strings.TrimSpace(authorization) == "" {
		response := ApiResponse{Code: 401, Msg: "Authorization required"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	// 获取服务器信息
	server, err := a.jsonService.GetServerByID(serverID, authorization, clientJson)
	if err != nil {
		log.Printf("Failed to get server info: %v", err)
		response := ApiResponse{Code: 500, Msg: "获取服务器信息失败"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	if server == nil {
		response := ApiResponse{Code: 404, Msg: "服务器不存在"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	// 执行SSH命令 - 先设置可执行权限，然后执行脚本
	command := fmt.Sprintf("cd %s && chmod +x codedeploy.sh && ./codedeploy.sh update %s", server.DefaultPath, projectID)
	output, err := a.executeSSHCommand(server, command)
	if err != nil {
		log.Printf("Failed to execute update command: %v", err)
		response := ApiResponse{Code: 500, Msg: fmt.Sprintf("执行更新命令失败: %v", err)}
		result, _ := json.Marshal(response)
		return string(result)
	}

	response := ApiResponse{
		Code: 200,
		Msg:  "项目更新成功",
		Data: map[string]interface{}{
			"command": command,
			"output":  output,
		},
	}
	result, _ := json.Marshal(response)
	return string(result)
}

// ProjectInitWithData 使用前端传入的服务器数据执行项目初始化
func (a *App) ProjectInitWithData(serverID, projectID, serverDataJson, authorization string) string {
	log.Printf("ProjectInitWithData called with serverID: %s, projectID: %s", serverID, projectID)

	// 检查授权
	if authorization == "" || strings.TrimSpace(authorization) == "" {
		response := ApiResponse{Code: 401, Msg: "Authorization required"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	// 解析前端传入的服务器数据
	var server services.ServerData
	if err := json.Unmarshal([]byte(serverDataJson), &server); err != nil {
		log.Printf("Failed to unmarshal server data: %v", err)
		response := ApiResponse{Code: 400, Msg: "服务器数据格式错误"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	// 验证服务器ID是否匹配
	if server.ServerID != serverID {
		log.Printf("Server ID mismatch: expected %s, got %s", serverID, server.ServerID)
		response := ApiResponse{Code: 400, Msg: "服务器ID不匹配"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	// 验证项目是否存在
	var targetProject *services.ProjectData
	for _, project := range server.ProjectList {
		if project.ProjectID == projectID {
			targetProject = &project
			break
		}
	}

	if targetProject == nil {
		log.Printf("Project %s not found in server %s", projectID, serverID)
		response := ApiResponse{Code: 404, Msg: "项目不存在"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	log.Printf("Using frontend data to init project: %s (%s)", targetProject.ProjectName, projectID)

	// 执行SSH命令 - 先设置可执行权限，然后执行脚本
	command := fmt.Sprintf("cd %s && chmod +x codedeploy.sh && ./codedeploy.sh init %s", server.DefaultPath, projectID)
	output, err := a.executeSSHCommand(&server, command)
	if err != nil {
		log.Printf("Failed to execute init command: %v", err)
		response := ApiResponse{Code: 500, Msg: fmt.Sprintf("执行初始化命令失败: %v", err)}
		result, _ := json.Marshal(response)
		return string(result)
	}

	response := ApiResponse{
		Code: 200,
		Msg:  "项目初始化成功（使用前端数据）",
		Data: map[string]interface{}{
			"command":     command,
			"output":      output,
			"project":     targetProject,
			"data_source": "frontend",
		},
	}
	result, _ := json.Marshal(response)
	return string(result)
}

// ProjectUpdateWithData 使用前端传入的服务器数据执行项目更新
func (a *App) ProjectUpdateWithData(serverID, projectID, serverDataJson, authorization string) string {
	log.Printf("ProjectUpdateWithData called with serverID: %s, projectID: %s", serverID, projectID)

	// 检查授权
	if authorization == "" || strings.TrimSpace(authorization) == "" {
		response := ApiResponse{Code: 401, Msg: "Authorization required"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	// 解析前端传入的服务器数据
	var server services.ServerData
	if err := json.Unmarshal([]byte(serverDataJson), &server); err != nil {
		log.Printf("Failed to unmarshal server data: %v", err)
		response := ApiResponse{Code: 400, Msg: "服务器数据格式错误"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	// 验证服务器ID是否匹配
	if server.ServerID != serverID {
		log.Printf("Server ID mismatch: expected %s, got %s", serverID, server.ServerID)
		response := ApiResponse{Code: 400, Msg: "服务器ID不匹配"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	// 验证项目是否存在
	var targetProject *services.ProjectData
	for _, project := range server.ProjectList {
		if project.ProjectID == projectID {
			targetProject = &project
			break
		}
	}

	if targetProject == nil {
		log.Printf("Project %s not found in server %s", projectID, serverID)
		response := ApiResponse{Code: 404, Msg: "项目不存在"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	log.Printf("Using frontend data to update project: %s (%s)", targetProject.ProjectName, projectID)

	// 执行SSH命令 - 先设置可执行权限，然后执行脚本
	command := fmt.Sprintf("cd %s && chmod +x codedeploy.sh && ./codedeploy.sh update %s", server.DefaultPath, projectID)
	output, err := a.executeSSHCommand(&server, command)
	if err != nil {
		log.Printf("Failed to execute update command: %v", err)
		response := ApiResponse{Code: 500, Msg: fmt.Sprintf("执行更新命令失败: %v", err)}
		result, _ := json.Marshal(response)
		return string(result)
	}

	response := ApiResponse{
		Code: 200,
		Msg:  "项目更新成功（使用前端数据）",
		Data: map[string]interface{}{
			"command":     command,
			"output":      output,
			"project":     targetProject,
			"data_source": "frontend",
		},
	}
	result, _ := json.Marshal(response)
	return string(result)
}

// 辅助函数：从URL中提取域名
func extractDomainFromURL(urlStr string) string {
	if urlStr == "" {
		return ""
	}

	// 如果不包含协议，添加默认协议
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		urlStr = "https://" + urlStr
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		// 如果解析失败，尝试简单的字符串处理
		parts := strings.Split(urlStr, "/")
		if len(parts) > 0 {
			return strings.Replace(parts[0], "https://", "", 1)
		}
		return urlStr
	}

	return parsedURL.Hostname()
}

// 辅助函数：获取端口或默认值
func getPortOrDefault(port, defaultPort string) string {
	if port == "" {
		return defaultPort
	}
	return port
}

// 辅助函数：从map中安全获取字符串值
func getStringFromMap(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// 辅助函数：通过SSH上传文件
func (a *App) uploadFileViaSSH(server *services.ServerData, filename, content string) error {
	// 创建SSH配置
	config := &ssh.ClientConfig{
		User:            server.ServerUser,
		Auth:            []ssh.AuthMethod{ssh.Password(server.ServerPassword)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}

	// 连接SSH
	address := fmt.Sprintf("%s:%s", server.ServerIP, server.ServerPort)
	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return fmt.Errorf("SSH连接失败: %v", err)
	}
	defer client.Close()

	// 创建会话
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("创建SSH会话失败: %v", err)
	}
	defer session.Close()

	// 创建目标文件路径
	targetPath := fmt.Sprintf("%s/%s", server.DefaultPath, filename)

	// 使用cat命令写入文件内容
	command := fmt.Sprintf("cat > %s", targetPath)
	session.Stdin = strings.NewReader(content)

	err = session.Run(command)
	if err != nil {
		return fmt.Errorf("上传文件失败: %v", err)
	}

	return nil
}

// 辅助函数：执行SSH命令
func (a *App) executeSSHCommand(server *services.ServerData, command string) (string, error) {
	// 创建SSH配置
	config := &ssh.ClientConfig{
		User:            server.ServerUser,
		Auth:            []ssh.AuthMethod{ssh.Password(server.ServerPassword)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}

	// 连接SSH
	address := fmt.Sprintf("%s:%s", server.ServerIP, server.ServerPort)
	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return "", fmt.Errorf("SSH连接失败: %v", err)
	}
	defer client.Close()

	// 创建会话
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("创建SSH会话失败: %v", err)
	}
	defer session.Close()

	// 执行命令
	output, err := session.CombinedOutput(command)
	if err != nil {
		return string(output), fmt.Errorf("命令执行失败: %v", err)
	}

	return string(output), nil
}

// processReleaseAndUploadConfig 处理 release.zip 并上传配置文件
func (a *App) processReleaseAndUploadConfig(server *services.ServerData, filename, content string) error {
	// 创建SSH配置
	config := &ssh.ClientConfig{
		User:            server.ServerUser,
		Auth:            []ssh.AuthMethod{ssh.Password(server.ServerPassword)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}

	// 连接SSH
	address := fmt.Sprintf("%s:%s", server.ServerIP, server.ServerPort)
	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return fmt.Errorf("SSH连接失败: %v", err)
	}
	defer client.Close()

	// 1. 检查 release.zip 是否存在
	log.Printf("Checking for release.zip in %s", server.DefaultPath)
	checkSession, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("创建检查会话失败: %v", err)
	}

	releaseZipPath := fmt.Sprintf("%s/release.zip", server.DefaultPath)
	checkCommand := fmt.Sprintf("test -f %s && echo 'exists' || echo 'not_exists'", releaseZipPath)
	checkOutput, err := checkSession.CombinedOutput(checkCommand)
	checkSession.Close()

	if err != nil {
		log.Printf("检查 release.zip 失败: %v", err)
	}

	// 2. 如果 release.zip 存在，进行解压和删除操作
	if strings.TrimSpace(string(checkOutput)) == "exists" {
		log.Printf("Found release.zip, processing...")

		// 解压 release.zip（强制覆盖）
		unzipSession, err := client.NewSession()
		if err != nil {
			return fmt.Errorf("创建解压会话失败: %v", err)
		}

		unzipCommand := fmt.Sprintf("cd %s && unzip -o release.zip", server.DefaultPath)
		log.Printf("Executing unzip command: %s", unzipCommand)
		unzipOutput, err := unzipSession.CombinedOutput(unzipCommand)
		unzipSession.Close()

		if err != nil {
			log.Printf("解压 release.zip 失败: %v, output: %s", err, string(unzipOutput))
			return fmt.Errorf("解压 release.zip 失败: %v", err)
		}

		log.Printf("Unzip successful, output: %s", string(unzipOutput))

		// 删除 release.zip
		deleteSession, err := client.NewSession()
		if err != nil {
			return fmt.Errorf("创建删除会话失败: %v", err)
		}

		deleteCommand := fmt.Sprintf("rm -f %s", releaseZipPath)
		log.Printf("Executing delete command: %s", deleteCommand)
		deleteOutput, err := deleteSession.CombinedOutput(deleteCommand)
		deleteSession.Close()

		if err != nil {
			log.Printf("删除 release.zip 失败: %v, output: %s", err, string(deleteOutput))
			return fmt.Errorf("删除 release.zip 失败: %v", err)
		}

		log.Printf("Release.zip deleted successfully")
	} else {
		log.Printf("No release.zip found, skipping extraction")
	}

	// 3. 上传配置文件
	log.Printf("Uploading config file: %s", filename)
	uploadSession, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("创建上传会话失败: %v", err)
	}
	defer uploadSession.Close()

	// 创建目标文件路径
	targetPath := fmt.Sprintf("%s/%s", server.DefaultPath, filename)

	// 使用cat命令写入文件内容
	command := fmt.Sprintf("cat > %s", targetPath)
	uploadSession.Stdin = strings.NewReader(content)

	err = uploadSession.Run(command)
	if err != nil {
		return fmt.Errorf("上传配置文件失败: %v", err)
	}

	log.Printf("Config file uploaded successfully to: %s", targetPath)
	return nil
}

// CapturePage 抓取页面内容
func (a *App) CapturePage(targetURL, optionsJson string) string {
	log.Printf("CapturePage called with URL: %s, options: %s", targetURL, optionsJson)

	// 验证URL
	if targetURL == "" || strings.TrimSpace(targetURL) == "" {
		response := ApiResponse{Code: 400, Msg: "URL不能为空"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	// 解析选项
	var options services.CaptureOptions
	if optionsJson != "" {
		if err := json.Unmarshal([]byte(optionsJson), &options); err != nil {
			log.Printf("Failed to unmarshal options: %v", err)
			// 使用默认选项
			options = services.CaptureOptions{
				IncludeImages:   true,
				IncludeStyles:   true,
				IncludeScripts:  true,
				FollowRedirects: true,
				Timeout:         60,
				CreateZip:       true,
				MaxFiles:        200,
			}
		}
	} else {
		// 默认选项
		options = services.CaptureOptions{
			IncludeImages:   true,
			IncludeStyles:   true,
			IncludeScripts:  true,
			FollowRedirects: true,
			Timeout:         60,
			CreateZip:       true,
			MaxFiles:        200,
		}
	}

	// 设置超时范围
	if options.Timeout < 60 {
		options.Timeout = 60
	}
	if options.Timeout > 300 {
		options.Timeout = 300
	}

	// 设置文件数量范围
	if options.MaxFiles < 200 {
		options.MaxFiles = 200
	}
	if options.MaxFiles > 1000 {
		options.MaxFiles = 1000
	}

	// 强制启用ZIP创建
	options.CreateZip = true

	log.Printf("Using options: %+v", options)

	// 执行页面抓取
	result, err := a.pageCaptureService.CapturePage(targetURL, options)
	if err != nil {
		log.Printf("Failed to capture page: %v", err)

		// 根据错误类型返回更具体的错误信息
		var errorMsg string
		if strings.Contains(err.Error(), "HTTP错误") {
			errorMsg = fmt.Sprintf("无法访问页面: %v", err)
		} else if strings.Contains(err.Error(), "不支持的内容类型") {
			errorMsg = fmt.Sprintf("页面类型不支持: %v", err)
		} else if strings.Contains(err.Error(), "请求失败") {
			errorMsg = fmt.Sprintf("网络请求失败: %v", err)
		} else {
			errorMsg = fmt.Sprintf("页面抓取失败: %v", err)
		}

		response := ApiResponse{Code: 500, Msg: errorMsg}
		responseResult, _ := json.Marshal(response)
		return string(responseResult)
	}

	// 检查结果是否有效
	if result == nil {
		log.Printf("Page capture returned nil result")
		response := ApiResponse{Code: 500, Msg: "页面抓取返回空结果"}
		responseResult, _ := json.Marshal(response)
		return string(responseResult)
	}

	log.Printf("Page captured successfully: status=%d, contentLength=%d, duration=%dms",
		result.StatusCode, result.ContentLength, result.Duration)

	// 检查内容是否为空或乱码
	if result.Content == "" {
		log.Printf("Warning: Captured content is empty")
		// 不设置默认内容，让前端知道内容为空但抓取成功
	}

	response := ApiResponse{Code: 200, Msg: "页面抓取成功", Data: result}
	responseResult, _ := json.Marshal(response)
	return string(responseResult)
}

// DownloadFile 下载文件并返回API响应格式
func (a *App) DownloadFile(filePath string) string {
	log.Printf("DownloadFile called with path: %s", filePath)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		response := ApiResponse{Code: 404, Msg: fmt.Sprintf("文件不存在: %s", filePath)}
		result, _ := json.Marshal(response)
		return string(result)
	}

	// 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		response := ApiResponse{Code: 500, Msg: fmt.Sprintf("读取文件失败: %v", err)}
		result, _ := json.Marshal(response)
		return string(result)
	}

	log.Printf("File downloaded successfully, size: %d bytes", len(content))

	// 使用Base64编码传输二进制数据，确保数据完整性
	base64Data := base64.StdEncoding.EncodeToString(content)

	// 返回成功响应，包含Base64编码的文件内容
	response := ApiResponse{Code: 200, Msg: "文件下载成功", Data: base64Data}
	result, _ := json.Marshal(response)
	return string(result)
}

// SelectDirectory 选择目录
func (a *App) SelectDirectory() string {
	log.Printf("SelectDirectory called")

	// 使用Wails的目录选择对话框
	selectedDir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择保存目录",
	})

	if err != nil {
		log.Printf("Failed to open directory dialog: %v", err)
		response := ApiResponse{Code: 500, Msg: fmt.Sprintf("打开目录选择对话框失败: %v", err)}
		result, _ := json.Marshal(response)
		return string(result)
	}

	if selectedDir == "" {
		// 用户取消选择
		response := ApiResponse{Code: 400, Msg: "用户取消选择目录"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	log.Printf("Directory selected: %s", selectedDir)
	response := ApiResponse{Code: 200, Msg: "目录选择成功", Data: selectedDir}
	result, _ := json.Marshal(response)
	return string(result)
}

// SaveZipToDirectory 保存ZIP文件到指定目录
func (a *App) SaveZipToDirectory(sourcePath, targetDirectory, fileName string) string {
	log.Printf("SaveZipToDirectory called: source=%s, target=%s, fileName=%s", sourcePath, targetDirectory, fileName)

	// 检查源文件是否存在
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		response := ApiResponse{Code: 404, Msg: fmt.Sprintf("源文件不存在: %s", sourcePath)}
		result, _ := json.Marshal(response)
		return string(result)
	}

	// 检查目标目录是否存在
	if _, err := os.Stat(targetDirectory); os.IsNotExist(err) {
		response := ApiResponse{Code: 404, Msg: fmt.Sprintf("目标目录不存在: %s", targetDirectory)}
		result, _ := json.Marshal(response)
		return string(result)
	}

	// 构建目标文件路径
	targetPath := filepath.Join(targetDirectory, fileName)

	// 复制文件
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		response := ApiResponse{Code: 500, Msg: fmt.Sprintf("打开源文件失败: %v", err)}
		result, _ := json.Marshal(response)
		return string(result)
	}
	defer sourceFile.Close()

	targetFile, err := os.Create(targetPath)
	if err != nil {
		response := ApiResponse{Code: 500, Msg: fmt.Sprintf("创建目标文件失败: %v", err)}
		result, _ := json.Marshal(response)
		return string(result)
	}
	defer targetFile.Close()

	_, err = io.Copy(targetFile, sourceFile)
	if err != nil {
		response := ApiResponse{Code: 500, Msg: fmt.Sprintf("复制文件失败: %v", err)}
		result, _ := json.Marshal(response)
		return string(result)
	}

	log.Printf("File saved successfully to: %s", targetPath)
	response := ApiResponse{Code: 200, Msg: "文件保存成功", Data: targetPath}
	result, _ := json.Marshal(response)
	return string(result)
}
