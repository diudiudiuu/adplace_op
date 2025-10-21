package services

import (
	"encoding/json"
	"fmt"
	"log"
)

// ServerData 服务器数据结构
type ServerData struct {
	ServerID       string        `json:"server_id"`
	ServerName     string        `json:"server_name"`
	ServerIP       string        `json:"server_ip"`
	ServerPort     string        `json:"server_port"`
	ServerUser     string        `json:"server_user"`
	ServerPassword string        `json:"server_password"`
	ProjectList    []ProjectData `json:"project_list"`
}

// ProjectData 项目数据结构
type ProjectData struct {
	ProjectID        string `json:"project_id"`
	ProjectName      string `json:"project_name"`
	ProjectManageURL string `json:"project_manage_url"`
	ProjectAPIURL    string `json:"project_api_url"`
}

const KV_KEY = "client_json"

// JsonService JSON数据管理服务
type JsonService struct {
	kvService *KvService
}

// NewJsonService 创建JSON服务实例
func NewJsonService() *JsonService {
	return &JsonService{
		kvService: NewKvService(),
	}
}

// LoadJsonFile 加载JSON数据
func (s *JsonService) LoadJsonFile(authorization string) ([]ServerData, error) {
	servers, _, err := s.LoadJsonFileWithResponse(authorization)
	return servers, err
}

// LoadJsonFileWithResponse 加载JSON数据并返回KV响应
func (s *JsonService) LoadJsonFileWithResponse(authorization string) ([]ServerData, *KvResponse, error) {
	log.Printf("LoadJsonFileWithResponse called with authorization: %s", authorization)

	resp, err := s.kvService.GetKey(KV_KEY, authorization)
	if err != nil {
		log.Printf("Failed to get KV data: %v", err)
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
		log.Printf("Successfully loaded %d servers", len(servers))
		return servers, resp, nil
	}

	// 如果没有数据，创建空数据
	log.Printf("No data found, creating empty data")
	emptyData := []ServerData{}
	s.SaveJsonFile(emptyData, authorization)
	return emptyData, resp, nil
}

// SaveJsonFile 保存JSON数据
func (s *JsonService) SaveJsonFile(data []ServerData, authorization string) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = s.kvService.UpdateKey(KV_KEY, string(jsonData), authorization)
	return err
}

// GetServerByID 根据ID获取服务器信息
func (s *JsonService) GetServerByID(serverID, authorization string) (*ServerData, error) {
	servers, err := s.LoadJsonFile(authorization)
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
func (s *JsonService) GetProjectByID(projectID, authorization string) (*ProjectData, error) {
	servers, err := s.LoadJsonFile(authorization)
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
func (s *JsonService) AddOrUpdateProject(serverID string, projectInfo ProjectData, authorization string) error {
	servers, err := s.LoadJsonFile(authorization)
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

			return s.SaveJsonFile(servers, authorization)
		}
	}

	return fmt.Errorf("服务器ID %s 不存在", serverID)
}

// DeleteProject 删除项目
func (s *JsonService) DeleteProject(serverID, projectID, authorization string) error {
	servers, err := s.LoadJsonFile(authorization)
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
					return s.SaveJsonFile(servers, authorization)
				}
			}
		}
	}

	return nil
}

// AddServer 添加新服务器
func (s *JsonService) AddServer(serverData ServerData, authorization string) error {
	servers, err := s.LoadJsonFile(authorization)
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

	return s.SaveJsonFile(servers, authorization)
}

// UpdateServer 更新服务器信息
func (s *JsonService) UpdateServer(serverID string, updatedServer ServerData, authorization string) error {
	servers, err := s.LoadJsonFile(authorization)
	if err != nil {
		return err
	}

	// 查找并更新服务器
	for i, server := range servers {
		if server.ServerID == serverID {
			// 保留原有的项目列表
			updatedServer.ProjectList = server.ProjectList
			servers[i] = updatedServer
			return s.SaveJsonFile(servers, authorization)
		}
	}

	return nil
}

// UpdateServerWithNewID 更新服务器信息（支持更改服务器ID）
func (s *JsonService) UpdateServerWithNewID(oldServerID string, updatedServer ServerData, authorization string) error {
	servers, err := s.LoadJsonFile(authorization)
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
			// 保留原有的项目列表
			updatedServer.ProjectList = server.ProjectList
			servers[i] = updatedServer
			return s.SaveJsonFile(servers, authorization)
		}
	}

	return fmt.Errorf("服务器ID %s 不存在", oldServerID)
}

// DeleteServer 删除服务器
func (s *JsonService) DeleteServer(serverID string, authorization string) error {
	servers, err := s.LoadJsonFile(authorization)
	if err != nil {
		return err
	}

	// 查找并删除服务器
	for i, server := range servers {
		if server.ServerID == serverID {
			servers = append(servers[:i], servers[i+1:]...)
			return s.SaveJsonFile(servers, authorization)
		}
	}

	return nil
}
