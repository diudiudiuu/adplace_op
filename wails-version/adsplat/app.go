package main

import (
	"adsplat/services"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/crypto/ssh"
)

// App struct
type App struct {
	ctx         context.Context
	jsonService *services.JsonService
	aesService  *services.AesService
	kvService   *services.KvService
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		jsonService: services.NewJsonService(),
		aesService:  services.NewAesService(),
		kvService:   services.NewKvService(),
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
func (a *App) List(authorization string) string {
	log.Printf("List called with authorization: %s", authorization)

	// 检查授权
	if authorization == "" || strings.TrimSpace(authorization) == "" {
		log.Printf("No authorization provided")
		response := ApiResponse{Code: 401, Msg: "Authorization required"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	servers, kvResponse, err := a.jsonService.LoadJsonFileWithResponse(authorization)
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
func (a *App) ServerInfo(serverID, authorization string) string {
	log.Printf("ServerInfo called with serverID: %s, authorization: %s", serverID, authorization)

	// 检查授权
	if authorization == "" || strings.TrimSpace(authorization) == "" {
		response := ApiResponse{Code: 401, Msg: "Authorization required"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	server, err := a.jsonService.GetServerByID(serverID, authorization)
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
func (a *App) ProjectInfo(projectID, authorization string) string {
	log.Printf("ProjectInfo called with projectID: %s, authorization: %s", projectID, authorization)

	// 检查授权
	if authorization == "" || strings.TrimSpace(authorization) == "" {
		response := ApiResponse{Code: 401, Msg: "Authorization required"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	project, err := a.jsonService.GetProjectByID(projectID, authorization)
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
func (a *App) ProjectForm(serverID, projectInfo, authorization string) string {
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

	err := a.jsonService.AddOrUpdateProject(serverID, project, authorization)
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
func (a *App) ProjectDelete(serverID, projectID, authorization string) string {
	log.Printf("ProjectDelete called with serverID: %s, projectID: %s, authorization: %s",
		serverID, projectID, authorization)

	// 检查授权
	if authorization == "" || strings.TrimSpace(authorization) == "" {
		response := ApiResponse{Code: 401, Msg: "Authorization required"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	err := a.jsonService.DeleteProject(serverID, projectID, authorization)
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
func (a *App) Exec(projectID, sql, sqlType, authorization string) string {
	log.Printf("Exec called with projectID: %s, sql: %s, sqlType: %s, authorization: %s",
		projectID, sql, sqlType, authorization)

	// 获取项目信息
	project, err := a.jsonService.GetProjectByID(projectID, authorization)
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

	return string(body)
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
func (a *App) ServerAdd(serverID, serverName, serverIP, serverPort, serverUser, serverPassword, authorization string) string {
	log.Printf("ServerAdd called with serverID: %s, serverName: %s, authorization: %s", serverID, serverName, authorization)

	// 检查授权
	if authorization == "" || strings.TrimSpace(authorization) == "" {
		response := ApiResponse{Code: 401, Msg: "Authorization required"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	// 创建新服务器数据
	newServer := services.ServerData{
		ServerID:       serverID,
		ServerName:     serverName,
		ServerIP:       serverIP,
		ServerPort:     serverPort,
		ServerUser:     serverUser,
		ServerPassword: serverPassword,
		ProjectList:    []services.ProjectData{},
	}

	err := a.jsonService.AddServer(newServer, authorization)
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
func (a *App) ServerUpdate(oldServerID, newServerID, serverName, serverIP, serverPort, serverUser, serverPassword, authorization string) string {
	log.Printf("ServerUpdate called with oldServerID: %s, newServerID: %s, serverName: %s, authorization: %s", oldServerID, newServerID, serverName, authorization)

	// 检查授权
	if authorization == "" || strings.TrimSpace(authorization) == "" {
		response := ApiResponse{Code: 401, Msg: "Authorization required"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	// 创建更新的服务器数据
	updatedServer := services.ServerData{
		ServerID:       newServerID,
		ServerName:     serverName,
		ServerIP:       serverIP,
		ServerPort:     serverPort,
		ServerUser:     serverUser,
		ServerPassword: serverPassword,
	}

	err := a.jsonService.UpdateServerWithNewID(oldServerID, updatedServer, authorization)
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
func (a *App) ServerDelete(serverID, authorization string) string {
	log.Printf("ServerDelete called with serverID: %s, authorization: %s", serverID, authorization)

	// 检查授权
	if authorization == "" || strings.TrimSpace(authorization) == "" {
		response := ApiResponse{Code: 401, Msg: "Authorization required"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	err := a.jsonService.DeleteServer(serverID, authorization)
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
func (a *App) TestStoredServerSSH(serverID, authorization string) string {
	log.Printf("TestStoredServerSSH called with serverID: %s", serverID)

	// 检查授权
	if authorization == "" || strings.TrimSpace(authorization) == "" {
		response := ApiResponse{Code: 401, Msg: "Authorization required"}
		result, _ := json.Marshal(response)
		return string(result)
	}

	// 获取服务器信息
	server, err := a.jsonService.GetServerByID(serverID, authorization)
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
	err = a.jsonService.UpdateServerConnectionStatus(serverID, testResult, authorization)
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
		User: serverUser,
		Auth: []ssh.AuthMethod{},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 注意：生产环境应该验证主机密钥
		Timeout: 10 * time.Second,
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
			Msg: fmt.Sprintf("SSH连接失败: %v", err),
			Data: map[string]interface{}{
				"connected": false,
				"error": err.Error(),
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
			Msg: fmt.Sprintf("创建SSH会话失败: %v", err),
			Data: map[string]interface{}{
				"connected": true,
				"session_error": err.Error(),
				"test_time": time.Now().Format("2006-01-02 15:04:05"),
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
		Msg: "SSH连接成功",
		Data: map[string]interface{}{
			"connected": true,
			"user": strings.TrimSpace(string(output)),
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
