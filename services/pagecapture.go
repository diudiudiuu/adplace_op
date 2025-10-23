package services

import (
	"archive/zip"
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// PageCaptureService 页面抓取服务
type PageCaptureService struct {
	client          *http.Client
	baseURL         *url.URL
	downloadedFiles map[string]string // URL -> 本地路径映射
	tempDir         string
}

// NewPageCaptureService 创建新的页面抓取服务
func NewPageCaptureService() *PageCaptureService {
	return &PageCaptureService{
		downloadedFiles: make(map[string]string),
	}
}

// CaptureOptions 抓取选项
type CaptureOptions struct {
	IncludeImages   bool `json:"includeImages"`
	IncludeStyles   bool `json:"includeStyles"`
	IncludeScripts  bool `json:"includeScripts"`
	FollowRedirects bool `json:"followRedirects"`
	Timeout         int  `json:"timeout"`
	CreateZip       bool `json:"createZip"` // 是否创建ZIP包
	MaxDepth        int  `json:"maxDepth"`  // 最大抓取深度
	MaxFiles        int  `json:"maxFiles"`  // 最大文件数量
}

// CaptureResult 抓取结果
type CaptureResult struct {
	StatusCode      int                 `json:"statusCode"`
	ContentType     string              `json:"contentType"`
	ContentLength   int64               `json:"contentLength"`
	Content         string              `json:"content"`
	Headers         map[string][]string `json:"headers"`
	Duration        int64               `json:"duration"`
	ZipPath         string              `json:"zipPath,omitempty"`         // ZIP文件路径
	ZipSize         int64               `json:"zipSize,omitempty"`         // ZIP文件大小
	FilesCount      int                 `json:"filesCount,omitempty"`      // 下载的文件数量
	DownloadedFiles []string            `json:"downloadedFiles,omitempty"` // 下载的文件列表
}

// ResourceInfo 资源信息
type ResourceInfo struct {
	URL       string
	LocalPath string
	Type      string // css, js, img, etc.
}

// CapturePage 抓取页面内容
func (s *PageCaptureService) CapturePage(targetURL string, options CaptureOptions) (*CaptureResult, error) {
	startTime := time.Now()

	// 验证URL
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return nil, fmt.Errorf("无效的URL格式: %v", err)
	}

	// 确保URL有协议
	if parsedURL.Scheme == "" {
		parsedURL.Scheme = "https"
		targetURL = parsedURL.String()
		parsedURL, _ = url.Parse(targetURL)
	}

	s.baseURL = parsedURL

	// 创建HTTP客户端
	s.client = &http.Client{
		Timeout: time.Duration(options.Timeout) * time.Second,
	}

	// 如果不跟随重定向，设置重定向策略
	if !options.FollowRedirects {
		s.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	// 如果需要创建ZIP，创建临时目录
	if options.CreateZip {
		tempDir, err := os.MkdirTemp("", "page_capture_*")
		if err != nil {
			return nil, fmt.Errorf("创建临时目录失败: %v", err)
		}
		s.tempDir = tempDir
		defer os.RemoveAll(tempDir) // 清理临时目录
	}

	// 设置默认值
	if options.MaxFiles <= 0 {
		options.MaxFiles = 200
	}
	if options.MaxDepth <= 0 {
		options.MaxDepth = 2
	}

	// 下载主页面
	content, resp, err := s.downloadPage(targetURL)
	if err != nil {
		return nil, fmt.Errorf("下载主页面失败: %v", err)
	}

	result := &CaptureResult{
		StatusCode:      resp.StatusCode,
		ContentType:     resp.Header.Get("Content-Type"),
		ContentLength:   resp.ContentLength,
		Headers:         resp.Header,
		Duration:        time.Since(startTime).Milliseconds(),
		FilesCount:      1,
		DownloadedFiles: []string{"index.html"},
	}

	if options.CreateZip {
		// 保存主页面到临时目录
		indexPath := filepath.Join(s.tempDir, "index.html")
		if err := os.WriteFile(indexPath, []byte(content), 0644); err != nil {
			return nil, fmt.Errorf("保存主页面失败: %v", err)
		}

		// 解析并下载资源
		modifiedContent, downloadedCount := s.downloadResources(content, options)
		result.Content = modifiedContent
		result.FilesCount += downloadedCount

		// 更新主页面内容
		if err := os.WriteFile(indexPath, []byte(modifiedContent), 0644); err != nil {
			return nil, fmt.Errorf("更新主页面失败: %v", err)
		}

		// 创建ZIP文件
		zipPath, zipSize, err := s.createZipFile()
		if err != nil {
			return nil, fmt.Errorf("创建ZIP文件失败: %v", err)
		}

		result.ZipPath = zipPath
		result.ZipSize = zipSize

		// 获取下载的文件列表
		files, _ := s.getDownloadedFilesList()
		result.DownloadedFiles = files
	} else {
		// 不创建ZIP，只处理内容
		if !options.IncludeScripts {
			content = s.removeScripts(content)
		}
		if !options.IncludeStyles {
			content = s.removeStyles(content)
		}
		if !options.IncludeImages {
			content = s.removeImages(content)
		}
		result.Content = content
	}

	return result, nil
}

// downloadPage 下载页面内容
func (s *PageCaptureService) downloadPage(targetURL string) (string, *http.Response, error) {
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return "", nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置User-Agent和请求头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Cache-Control", "max-age=0")

	// 发送请求
	resp, err := s.client.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查HTTP状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return "", resp, fmt.Errorf("HTTP错误: %d %s", resp.StatusCode, resp.Status)
	}

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", resp, fmt.Errorf("读取响应内容失败: %v", err)
	}

	// 检查内容类型
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(strings.ToLower(contentType), "text/html") &&
		!strings.Contains(strings.ToLower(contentType), "application/xhtml") {
		return "", resp, fmt.Errorf("不支持的内容类型: %s", contentType)
	}

	// 转换为UTF-8字符串
	content := string(body)

	// 简单的编码检测和转换
	if !isValidUTF8(content) {
		// 尝试从GBK转换
		if gbkContent := convertFromGBK(body); gbkContent != "" {
			content = gbkContent
		}
	}

	return content, resp, nil
}

// isValidUTF8 检查字符串是否为有效的UTF-8
func isValidUTF8(s string) bool {
	for _, r := range s {
		if r == '\uFFFD' {
			return false
		}
	}
	return true
}

// convertFromGBK 尝试从GBK编码转换为UTF-8
func convertFromGBK(data []byte) string {
	// 这里可以添加更复杂的编码检测和转换逻辑
	// 暂时返回原始字符串
	return string(data)
}

// downloadResources 下载页面中的资源文件
func (s *PageCaptureService) downloadResources(content string, options CaptureOptions) (string, int) {
	downloadedCount := 0
	modifiedContent := content

	// 下载CSS文件
	if options.IncludeStyles {
		cssRegex := regexp.MustCompile(`<link[^>]*href\s*=\s*["']([^"']+)["'][^>]*rel\s*=\s*["']stylesheet["'][^>]*>|<link[^>]*rel\s*=\s*["']stylesheet["'][^>]*href\s*=\s*["']([^"']+)["'][^>]*>`)
		matches := cssRegex.FindAllStringSubmatch(content, -1)
		for _, match := range matches {
			cssURL := match[1]
			if cssURL == "" {
				cssURL = match[2]
			}
			if cssURL != "" && downloadedCount < options.MaxFiles {
				localPath := s.downloadResource(cssURL, "css")
				if localPath != "" {
					modifiedContent = strings.ReplaceAll(modifiedContent, cssURL, localPath)
					downloadedCount++
				}
			}
		}
	}

	// 下载JavaScript文件
	if options.IncludeScripts {
		jsRegex := regexp.MustCompile(`<script[^>]*src\s*=\s*["']([^"']+)["'][^>]*>`)
		matches := jsRegex.FindAllStringSubmatch(content, -1)
		for _, match := range matches {
			jsURL := match[1]
			if jsURL != "" && downloadedCount < options.MaxFiles {
				localPath := s.downloadResource(jsURL, "js")
				if localPath != "" {
					modifiedContent = strings.ReplaceAll(modifiedContent, jsURL, localPath)
					downloadedCount++
				}
			}
		}
	}

	// 下载图片文件
	if options.IncludeImages {
		imgRegex := regexp.MustCompile(`<img[^>]*src\s*=\s*["']([^"']+)["'][^>]*>`)
		matches := imgRegex.FindAllStringSubmatch(content, -1)
		for _, match := range matches {
			imgURL := match[1]
			if imgURL != "" && downloadedCount < options.MaxFiles {
				localPath := s.downloadResource(imgURL, "img")
				if localPath != "" {
					modifiedContent = strings.ReplaceAll(modifiedContent, imgURL, localPath)
					downloadedCount++
				}
			}
		}

		// 下载CSS中的背景图片
		bgImgRegex := regexp.MustCompile(`background-image\s*:\s*url\s*\(\s*["']?([^"')]+)["']?\s*\)`)
		matches = bgImgRegex.FindAllStringSubmatch(content, -1)
		for _, match := range matches {
			imgURL := match[1]
			if imgURL != "" && downloadedCount < options.MaxFiles {
				localPath := s.downloadResource(imgURL, "img")
				if localPath != "" {
					modifiedContent = strings.ReplaceAll(modifiedContent, imgURL, localPath)
					downloadedCount++
				}
			}
		}
	}

	return modifiedContent, downloadedCount
}

// downloadResource 下载单个资源文件
func (s *PageCaptureService) downloadResource(resourceURL, resourceType string) string {
	// 解析资源URL
	fullURL, err := s.resolveURL(resourceURL)
	if err != nil {
		return ""
	}

	// 检查是否已经下载过
	if localPath, exists := s.downloadedFiles[fullURL]; exists {
		return localPath
	}

	// 下载资源
	resp, err := s.client.Get(fullURL)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return ""
	}

	// 读取内容
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	// 生成本地文件名
	fileName := s.generateFileName(fullURL, resourceType)
	// 使用static目录结构，根据资源类型分类
	subDir := s.getStaticSubDir(resourceType)
	localPath := filepath.Join("static", subDir, fileName)
	fullLocalPath := filepath.Join(s.tempDir, localPath)

	// 创建目录
	dir := filepath.Dir(fullLocalPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return ""
	}

	// 保存文件
	if err := os.WriteFile(fullLocalPath, content, 0644); err != nil {
		return ""
	}

	// 记录下载的文件
	s.downloadedFiles[fullURL] = localPath

	return localPath
}

// resolveURL 解析相对URL为绝对URL
func (s *PageCaptureService) resolveURL(resourceURL string) (string, error) {
	if strings.HasPrefix(resourceURL, "http://") || strings.HasPrefix(resourceURL, "https://") {
		return resourceURL, nil
	}

	if strings.HasPrefix(resourceURL, "//") {
		return s.baseURL.Scheme + ":" + resourceURL, nil
	}

	if strings.HasPrefix(resourceURL, "/") {
		return s.baseURL.Scheme + "://" + s.baseURL.Host + resourceURL, nil
	}

	// 相对路径
	baseDir := path.Dir(s.baseURL.Path)
	if baseDir == "." {
		baseDir = ""
	}
	return s.baseURL.Scheme + "://" + s.baseURL.Host + baseDir + "/" + resourceURL, nil
}

// generateFileName 生成本地文件名
func (s *PageCaptureService) generateFileName(fullURL, resourceType string) string {
	parsedURL, err := url.Parse(fullURL)
	if err != nil {
		// 使用MD5作为文件名
		hash := fmt.Sprintf("%x", md5.Sum([]byte(fullURL)))
		return hash + s.getExtensionByType(resourceType)
	}

	fileName := path.Base(parsedURL.Path)
	if fileName == "" || fileName == "." || fileName == "/" {
		// 使用MD5作为文件名
		hash := fmt.Sprintf("%x", md5.Sum([]byte(fullURL)))
		fileName = hash + s.getExtensionByType(resourceType)
	}

	// 确保文件有正确的扩展名
	if !strings.Contains(fileName, ".") {
		fileName += s.getExtensionByType(resourceType)
	}

	return fileName
}

// getExtensionByType 根据资源类型获取扩展名
func (s *PageCaptureService) getExtensionByType(resourceType string) string {
	switch resourceType {
	case "css":
		return ".css"
	case "js":
		return ".js"
	case "img":
		return ".jpg"
	default:
		return ".txt"
	}
}

// getStaticSubDir 根据资源类型获取静态资源子目录
func (s *PageCaptureService) getStaticSubDir(resourceType string) string {
	switch resourceType {
	case "css":
		return "css"
	case "js":
		return "js"
	case "img":
		return "images"
	default:
		return "assets"
	}
}

// createZipFile 创建ZIP文件
func (s *PageCaptureService) createZipFile() (string, int64, error) {
	// 创建ZIP文件
	zipFileName := fmt.Sprintf("webpage_%d.zip", time.Now().Unix())
	zipPath := filepath.Join(os.TempDir(), zipFileName)

	zipFile, err := os.Create(zipPath)
	if err != nil {
		return "", 0, err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// 遍历临时目录，添加所有文件到ZIP
	err = filepath.Walk(s.tempDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// 计算相对路径
		relPath, err := filepath.Rel(s.tempDir, filePath)
		if err != nil {
			return err
		}

		// 在ZIP中创建文件
		zipFileWriter, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		// 读取文件内容
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		// 写入ZIP
		_, err = zipFileWriter.Write(fileContent)
		return err
	})

	if err != nil {
		return "", 0, err
	}

	// 获取ZIP文件大小
	zipInfo, err := zipFile.Stat()
	if err != nil {
		return zipPath, 0, nil
	}

	return zipPath, zipInfo.Size(), nil
}

// getDownloadedFilesList 获取下载的文件列表
func (s *PageCaptureService) getDownloadedFilesList() ([]string, error) {
	var files []string

	err := filepath.Walk(s.tempDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			relPath, err := filepath.Rel(s.tempDir, filePath)
			if err != nil {
				return err
			}
			files = append(files, relPath)
		}
		return nil
	})

	return files, err
}

// removeScripts 移除脚本标签
func (s *PageCaptureService) removeScripts(content string) string {
	scriptRegex := regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`)
	return scriptRegex.ReplaceAllString(content, "")
}

// removeStyles 移除样式标签和内联样式
func (s *PageCaptureService) removeStyles(content string) string {
	// 移除style标签
	styleRegex := regexp.MustCompile(`(?i)<style[^>]*>.*?</style>`)
	content = styleRegex.ReplaceAllString(content, "")

	// 移除link标签中的样式表
	linkRegex := regexp.MustCompile(`(?i)<link[^>]*rel\s*=\s*["']stylesheet["'][^>]*>`)
	content = linkRegex.ReplaceAllString(content, "")

	return content
}

// removeImages 移除图片标签
func (s *PageCaptureService) removeImages(content string) string {
	imgRegex := regexp.MustCompile(`(?i)<img[^>]*>`)
	return imgRegex.ReplaceAllString(content, "")
}
