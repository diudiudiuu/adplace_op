package services

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// ServerData 服务器数据结构
type ServerData struct {
	ServerID         string        `json:"server_id"`
	ServerName       string        `json:"server_name"`
	ServerIP         string        `json:"server_ip"`
	ServerPort       string        `json:"server_port"`
	ServerUser       string        `json:"server_user"`
	ServerPassword   string        `json:"server_password"`
	DefaultPath      string        `json:"default_path"`
	ProjectList      []ProjectData `json:"project_list"`
	ConnectionStatus string        `json:"connection_status,omitempty"` // "connected", "disconnected", "unknown"
	LastTestTime     string        `json:"last_test_time,omitempty"`
	LastTestResult   string        `json:"last_test_result,omitempty"`
}

// ProjectData 项目数据结构
type ProjectData struct {
	ProjectID        string `json:"project_id"`
	ProjectName      string `json:"project_name"`
	ProjectManageURL string `json:"project_manage_url"`
	ProjectAPIURL    string `json:"project_api_url"`
	APIPort          string `json:"api_port"`
	FrontPort        string `json:"front_port"`
}

// 移除固定的KV_KEY，改为使用传入的参数

// JsonService JSON数据管理服务
type JsonService struct {
	kvService     *KvService
	cache         []ServerData
	cacheTime     int64
	cacheDuration int64 // 缓存持续时间（毫秒）
}

// NewJsonService 创建JSON服务实例
func NewJsonService() *JsonService {
	return &JsonService{
		kvService:     NewKvService(),
		cache:         []ServerData{},
		cacheTime:     0,
		cacheDuration: 2 * 60 * 60 * 1000, // 2小时缓存（毫秒）
	}
}

// LoadJsonFile 加载JSON数据
func (s *JsonService) LoadJsonFile(authorization, clientJson string) ([]ServerData, error) {
	servers, _, err := s.LoadJsonFileWithResponse(authorization, clientJson)
	return servers, err
}

// LoadJsonFileWithResponse 加载JSON数据并返回KV响应
func (s *JsonService) LoadJsonFileWithResponse(authorization, clientJson string) ([]ServerData, *KvResponse, error) {
	log.Printf("LoadJsonFileWithResponse called with authorization: %s, clientJson: %s", authorization, clientJson)

	// 检查缓存是否有效
	currentTime := time.Now().UnixMilli()
	if s.cacheTime > 0 && (currentTime-s.cacheTime) < s.cacheDuration && len(s.cache) > 0 {
		log.Printf("Using cached data, age: %d ms", currentTime-s.cacheTime)
		return s.cache, nil, nil
	}

	log.Printf("Cache expired or empty, fetching from KV store")
	resp, err := s.kvService.GetKey(clientJson, authorization)
	if err != nil {
		log.Printf("Failed to get KV data: %v", err)
		// 如果有缓存数据，返回缓存数据
		if len(s.cache) > 0 {
			log.Printf("Returning cached data due to KV error")
			return s.cache, nil, err
		}
		return []ServerData{}, nil, err
	}

	log.Printf("KV response: Code=%d, Data=%v", resp.Code, resp.Data != nil)

	// 如果是 401 错误，直接返回
	if resp.Code == 401 {
		log.Printf("KV service returned 401 Unauthorized")
		return []ServerData{}, resp, nil
	}

	if resp.Code == 200 && resp.Data != nil {
		log.Printf("KV data value: %s", resp.Data.Value)
		var servers []ServerData
		if err := json.Unmarshal([]byte(resp.Data.Value), &servers); err != nil {
			log.Printf("Failed to unmarshal JSON: %v", err)
			return []ServerData{}, resp, nil
		}

		// 数据迁移：为缺少 default_path 的服务器添加默认值
		needsSave := false
		for i := range servers {
			if servers[i].DefaultPath == "" {
				servers[i].DefaultPath = "/adplace"
				needsSave = true
				log.Printf("Added default path '/adplace' to server: %s", servers[i].ServerID)
			}

			// 数据迁移：为项目添加默认端口值
			for j := range servers[i].ProjectList {
				if servers[i].ProjectList[j].APIPort == "" {
					servers[i].ProjectList[j].APIPort = "9000"
					needsSave = true
					log.Printf("Added default API port '9000' to project: %s", servers[i].ProjectList[j].ProjectID)
				}
				if servers[i].ProjectList[j].FrontPort == "" {
					servers[i].ProjectList[j].FrontPort = "3000"
					needsSave = true
					log.Printf("Added default front port '3000' to project: %s", servers[i].ProjectList[j].ProjectID)
				}
			}
		}

		// 如果有数据被修改，保存回去
		if needsSave {
			log.Printf("Saving updated server data with default paths")
			s.SaveJsonFile(servers, authorization, clientJson)
		}

		// 更新缓存
		s.cache = servers
		s.cacheTime = time.Now().UnixMilli()
		log.Printf("Successfully loaded %d servers and updated cache", len(servers))
		return servers, resp, nil
	}

	// 如果没有数据，创建空数据
	log.Printf("No data found, creating empty data")
	emptyData := []ServerData{}
	s.SaveJsonFile(emptyData, authorization, clientJson)
	
	// 更新缓存
	s.cache = emptyData
	s.cacheTime = time.Now().UnixMilli()
	return emptyData, resp, nil
}

// SaveJsonFile 保存JSON数据
func (s *JsonService) SaveJsonFile(data []ServerData, authorization, clientJson string) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// 调试：打印要保存的 JSON 数据
	log.Printf("Saving JSON data: %s", string(jsonData))

	_, err = s.kvService.UpdateKey(clientJson, string(jsonData), authorization)
	
	// 清除缓存，强制下次重新获取
	if err == nil {
		s.clearCache()
		log.Printf("Cache cleared after successful save")
	}
	
	return err
}

// clearCache 清除缓存
func (s *JsonService) clearCache() {
	s.cache = []ServerData{}
	s.cacheTime = 0
}

// GetServerByID 根据ID获取服务器信息
func (s *JsonService) GetServerByID(serverID, authorization, clientJson string) (*ServerData, error) {
	servers, err := s.LoadJsonFile(authorization, clientJson)
	if err != nil {
		return nil, err
	}

	for _, server := range servers {
		if server.ServerID == serverID {
			return &server, nil
		}
	}

	return nil, nil
}

// GetProjectByID 根据ID获取项目信息
func (s *JsonService) GetProjectByID(projectID, authorization, clientJson string) (*ProjectData, error) {
	servers, err := s.LoadJsonFile(authorization, clientJson)
	if err != nil {
		return nil, err
	}

	for _, server := range servers {
		for _, project := range server.ProjectList {
			if project.ProjectID == projectID {
				return &project, nil
			}
		}
	}

	return nil, nil
}

// AddOrUpdateProject 添加或更新项目
func (s *JsonService) AddOrUpdateProject(serverID string, projectInfo ProjectData, authorization, clientJson string) error {
	servers, err := s.LoadJsonFile(authorization, clientJson)
	if err != nil {
		return err
	}

	// 查找服务器
	for i, server := range servers {
		if server.ServerID == serverID {
			// 查找项目是否存在
			found := false
			for j, project := range server.ProjectList {
				if project.ProjectID == projectInfo.ProjectID {
					// 更新现有项目
					servers[i].ProjectList[j] = projectInfo
					found = true
					break
				}
			}

			// 如果项目不存在，添加新项目
			if !found {
				servers[i].ProjectList = append(servers[i].ProjectList, projectInfo)
			}

			return s.SaveJsonFile(servers, authorization, clientJson)
		}
	}

	return fmt.Errorf("服务器ID %s 不存在", serverID)
}

// UpdateServerConnectionStatus 更新服务器连接状态
func (s *JsonService) UpdateServerConnectionStatus(serverID, testResult, authorization, clientJson string) error {
	servers, err := s.LoadJsonFile(authorization, clientJson)
	if err != nil {
		return err
	}

	// 解析测试结果
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(testResult), &result); err != nil {
		return fmt.Errorf("解析测试结果失败: %v", err)
	}

	// 查找并更新服务器
	for i, server := range servers {
		if server.ServerID == serverID {
			// 更新连接状态
			if data, ok := result["data"].(map[string]interface{}); ok {
				if connected, ok := data["connected"].(bool); ok {
					if connected {
						servers[i].ConnectionStatus = "connected"
					} else {
						servers[i].ConnectionStatus = "disconnected"
					}
				}

				if testTime, ok := data["test_time"].(string); ok {
					servers[i].LastTestTime = testTime
				}
			}

			if msg, ok := result["msg"].(string); ok {
				servers[i].LastTestResult = msg
			}

			return s.SaveJsonFile(servers, authorization, clientJson)
		}
	}

	return fmt.Errorf("服务器ID %s 不存在", serverID)
}

// DeleteProject 删除项目
func (s *JsonService) DeleteProject(serverID, projectID, authorization, clientJson string) error {
	servers, err := s.LoadJsonFile(authorization, clientJson)
	if err != nil {
		return err
	}

	// 查找服务器
	for i, server := range servers {
		if server.ServerID == serverID {
			// 查找并删除项目
			for j, project := range server.ProjectList {
				if project.ProjectID == projectID {
					// 删除项目
					servers[i].ProjectList = append(
						servers[i].ProjectList[:j],
						servers[i].ProjectList[j+1:]...,
					)
					return s.SaveJsonFile(servers, authorization, clientJson)
				}
			}
		}
	}

	return nil
}

// AddServer 添加新服务器
func (s *JsonService) AddServer(serverData ServerData, authorization, clientJson string) error {
	servers, err := s.LoadJsonFile(authorization, clientJson)
	if err != nil {
		return err
	}

	// 检查服务器ID是否已存在
	for _, server := range servers {
		if server.ServerID == serverData.ServerID {
			return fmt.Errorf("服务器ID %s 已存在", serverData.ServerID)
		}
	}

	// 添加新服务器
	servers = append(servers, serverData)

	return s.SaveJsonFile(servers, authorization, clientJson)
}

// UpdateServer 更新服务器信息
func (s *JsonService) UpdateServer(serverID string, updatedServer ServerData, authorization, clientJson string) error {
	servers, err := s.LoadJsonFile(authorization, clientJson)
	if err != nil {
		return err
	}

	// 查找并更新服务器
	for i, server := range servers {
		if server.ServerID == serverID {
			// 保留原有的项目列表
			updatedServer.ProjectList = server.ProjectList
			servers[i] = updatedServer
			return s.SaveJsonFile(servers, authorization, clientJson)
		}
	}

	return nil
}

// UpdateServerWithNewID 更新服务器信息（支持更改服务器ID）
func (s *JsonService) UpdateServerWithNewID(oldServerID string, updatedServer ServerData, authorization, clientJson string) error {
	servers, err := s.LoadJsonFile(authorization, clientJson)
	if err != nil {
		return err
	}

	// 如果新ID与旧ID不同，需要检查新ID是否已存在
	if oldServerID != updatedServer.ServerID {
		for _, server := range servers {
			if server.ServerID == updatedServer.ServerID {
				return fmt.Errorf("服务器ID %s 已存在", updatedServer.ServerID)
			}
		}
	}

	// 查找并更新服务器
	for i, server := range servers {
		if server.ServerID == oldServerID {
			// 保留原有的项目列表和连接状态信息
			updatedServer.ProjectList = server.ProjectList
			updatedServer.ConnectionStatus = server.ConnectionStatus
			updatedServer.LastTestTime = server.LastTestTime
			updatedServer.LastTestResult = server.LastTestResult
			servers[i] = updatedServer
			return s.SaveJsonFile(servers, authorization, clientJson)
		}
	}

	return fmt.Errorf("服务器ID %s 不存在", oldServerID)
}

// DeleteServer 删除服务器
func (s *JsonService) DeleteServer(serverID, authorization, clientJson string) error {
	servers, err := s.LoadJsonFile(authorization, clientJson)
	if err != nil {
		return err
	}

	// 查找并删除服务器
	for i, server := range servers {
		if server.ServerID == serverID {
			servers = append(servers[:i], servers[i+1:]...)
			return s.SaveJsonFile(servers, authorization, clientJson)
		}
	}

	return nil
}
