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

	"github.com/PuerkitoBio/goquery"
)

// PageCaptureService 页面抓取服务
type PageCaptureService struct {
	client    *http.Client
	baseURL   *url.URL
	resources map[string]*ResourceInfo
	tempDir   string
	maxFiles  int
	fileCount int
}

// NewPageCaptureService 创建新的页面抓取服务
func NewPageCaptureService() *PageCaptureService {
	return &PageCaptureService{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		resources: make(map[string]*ResourceInfo),
	}
}

// CaptureOptions 抓取选项
type CaptureOptions struct {
	IncludeImages   bool `json:"includeImages"`
	IncludeStyles   bool `json:"includeStyles"`
	IncludeScripts  bool `json:"includeScripts"`
	FollowRedirects bool `json:"followRedirects"`
	Timeout         int  `json:"timeout"`
	CreateZip       bool `json:"createZip"`
	MaxFiles        int  `json:"maxFiles"`
	MaxDepth        int  `json:"maxDepth"`
}

// CaptureResult 抓取结果
type CaptureResult struct {
	StatusCode      int                 `json:"statusCode"`
	ContentType     string              `json:"contentType"`
	ContentLength   int64               `json:"contentLength"`
	Content         string              `json:"content"`
	Headers         map[string][]string `json:"headers"`
	Duration        int64               `json:"duration"`
	ZipPath         string              `json:"zipPath,omitempty"`
	ZipSize         int64               `json:"zipSize,omitempty"`
	FilesCount      int                 `json:"filesCount,omitempty"`
	DownloadedFiles []string            `json:"downloadedFiles,omitempty"`
}

// ResourceInfo 资源信息
type ResourceInfo struct {
	URL       string
	LocalPath string
	Type      string
	Content   []byte
}

// CapturePage 抓取页面内容
func (s *PageCaptureService) CapturePage(targetURL string, options CaptureOptions) (*CaptureResult, error) {
	startTime := time.Now()

	// 验证URL
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return nil, fmt.Errorf("无效的URL格式: %v", err)
	}
	if parsedURL.Scheme == "" {
		parsedURL.Scheme = "https"
		targetURL = parsedURL.String()
		parsedURL, _ = url.Parse(targetURL)
	}

	s.baseURL = parsedURL
	s.maxFiles = options.MaxFiles
	if s.maxFiles <= 0 {
		s.maxFiles = 200
	}

	// 设置超时
	s.client.Timeout = time.Duration(options.Timeout) * time.Second

	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "page_capture_*")
	if err != nil {
		return nil, fmt.Errorf("创建临时目录失败: %v", err)
	}
	s.tempDir = tempDir
	defer os.RemoveAll(tempDir)

	// 下载主页面
	htmlContent, resp, err := s.downloadPage(targetURL)
	if err != nil {
		return nil, fmt.Errorf("下载主页面失败: %v", err)
	}

	// 解析HTML并下载资源
	modifiedHTML, err := s.processHTMLAndDownloadResources(htmlContent, options)
	if err != nil {
		return nil, fmt.Errorf("处理HTML失败: %v", err)
	}

	// 保存文件
	err = s.saveAllFiles(modifiedHTML)
	if err != nil {
		return nil, fmt.Errorf("保存文件失败: %v", err)
	}

	// 创建ZIP
	zipPath, zipSize, err := s.createZipFile()
	if err != nil {
		return nil, fmt.Errorf("创建ZIP失败: %v", err)
	}

	// 构建结果
	result := &CaptureResult{
		StatusCode:      resp.StatusCode,
		ContentType:     resp.Header.Get("Content-Type"),
		ContentLength:   int64(len(htmlContent)),
		Headers:         resp.Header,
		Duration:        time.Since(startTime).Milliseconds(),
		ZipPath:         zipPath,
		ZipSize:         zipSize,
		FilesCount:      len(s.resources) + 1,
		DownloadedFiles: s.getFileList(),
	}

	return result, nil
}

// downloadPage 下载页面
func (s *PageCaptureService) downloadPage(targetURL string) (string, *http.Response, error) {
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return "", nil, err
	}

	// 设置请求头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return "", resp, fmt.Errorf("HTTP错误: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", resp, err
	}

	return string(body), resp, nil
}

// processHTMLAndDownloadResources 处理HTML并下载资源
func (s *PageCaptureService) processHTMLAndDownloadResources(htmlContent string, options CaptureOptions) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return htmlContent, err
	}

	// 下载CSS文件
	if options.IncludeStyles {
		doc.Find("link[rel=stylesheet]").Each(func(i int, sel *goquery.Selection) {
			if s.fileCount >= s.maxFiles {
				return
			}
			href, exists := sel.Attr("href")
			if exists {
				absoluteURL := s.resolveURL(href)
				if absoluteURL != "" {
					localPath := s.downloadResource(absoluteURL, "css")
					if localPath != "" {
						sel.SetAttr("href", localPath)
					}
				}
			}
		})
	}

	// 下载JavaScript文件
	if options.IncludeScripts {
		doc.Find("script[src]").Each(func(i int, sel *goquery.Selection) {
			if s.fileCount >= s.maxFiles {
				return
			}
			src, exists := sel.Attr("src")
			if exists {
				absoluteURL := s.resolveURL(src)
				if absoluteURL != "" {
					localPath := s.downloadResource(absoluteURL, "js")
					if localPath != "" {
						sel.SetAttr("src", localPath)
					}
				}
			}
		})
	}

	// 下载图片
	if options.IncludeImages {
		doc.Find("img[src]").Each(func(i int, sel *goquery.Selection) {
			if s.fileCount >= s.maxFiles {
				return
			}
			src, exists := sel.Attr("src")
			if exists {
				absoluteURL := s.resolveURL(src)
				if absoluteURL != "" {
					localPath := s.downloadResource(absoluteURL, "images")
					if localPath != "" {
						sel.SetAttr("src", localPath)
					}
				}
			}
		})

		// 处理CSS中的背景图片
		doc.Find("style").Each(func(i int, sel *goquery.Selection) {
			cssContent := sel.Text()
			modifiedCSS := s.processCSSContent(cssContent)
			sel.SetText(modifiedCSS)
		})
	}

	html, err := doc.Html()
	if err != nil {
		return htmlContent, err
	}

	// 格式化HTML
	formattedHTML := s.formatHTML(html)
	return formattedHTML, nil
}

// downloadResource 下载单个资源
func (s *PageCaptureService) downloadResource(resourceURL, resourceType string) string {
	if s.fileCount >= s.maxFiles {
		return ""
	}

	// 检查是否已下载
	if resource, exists := s.resources[resourceURL]; exists {
		return resource.LocalPath
	}

	// 下载资源
	resp, err := s.client.Get(resourceURL)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return ""
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	// 生成本地路径
	localPath := s.generateLocalPath(resourceURL, resourceType)

	// 存储资源信息
	s.resources[resourceURL] = &ResourceInfo{
		URL:       resourceURL,
		LocalPath: localPath,
		Type:      resourceType,
		Content:   content,
	}

	s.fileCount++
	return localPath
}

// resolveURL 解析URL
func (s *PageCaptureService) resolveURL(resourceURL string) string {
	if strings.HasPrefix(resourceURL, "http://") || strings.HasPrefix(resourceURL, "https://") {
		return resourceURL
	}

	if strings.HasPrefix(resourceURL, "//") {
		return s.baseURL.Scheme + ":" + resourceURL
	}

	if strings.HasPrefix(resourceURL, "/") {
		return s.baseURL.Scheme + "://" + s.baseURL.Host + resourceURL
	}

	// 相对路径
	baseDir := path.Dir(s.baseURL.Path)
	if baseDir == "." {
		baseDir = ""
	}
	return s.baseURL.Scheme + "://" + s.baseURL.Host + baseDir + "/" + resourceURL
}

// generateLocalPath 生成本地路径
func (s *PageCaptureService) generateLocalPath(resourceURL, resourceType string) string {
	parsedURL, err := url.Parse(resourceURL)
	if err != nil {
		// 使用MD5生成文件名
		hash := fmt.Sprintf("%x", md5.Sum([]byte(resourceURL)))
		return fmt.Sprintf("static/%s/%s%s", resourceType, hash, s.getExtensionByType(resourceType))
	}

	filename := path.Base(parsedURL.Path)
	if filename == "" || filename == "." {
		hash := fmt.Sprintf("%x", md5.Sum([]byte(resourceURL)))
		filename = hash + s.getExtensionByType(resourceType)
	}

	// 确保有扩展名
	if !strings.Contains(filename, ".") {
		filename += s.getExtensionByType(resourceType)
	}

	return fmt.Sprintf("static/%s/%s", resourceType, filename)
}

// getExtensionByType 获取扩展名
func (s *PageCaptureService) getExtensionByType(resourceType string) string {
	switch resourceType {
	case "css":
		return ".css"
	case "js":
		return ".js"
	case "images":
		return ".jpg"
	default:
		return ".txt"
	}
}

// processCSSContent 处理CSS内容中的URL
func (s *PageCaptureService) processCSSContent(cssContent string) string {
	if s.fileCount >= s.maxFiles {
		return cssContent
	}

	urlRegex := regexp.MustCompile(`url\s*\(\s*['"]?([^'")]+)['"]?\s*\)`)

	return urlRegex.ReplaceAllStringFunc(cssContent, func(match string) string {
		submatches := urlRegex.FindStringSubmatch(match)
		if len(submatches) > 1 {
			resourceURL := submatches[1]
			absoluteURL := s.resolveURL(resourceURL)
			if absoluteURL != "" {
				localPath := s.downloadResource(absoluteURL, "images")
				if localPath != "" {
					return strings.Replace(match, resourceURL, localPath, 1)
				}
			}
		}
		return match
	})
}

// saveAllFiles 保存所有文件
func (s *PageCaptureService) saveAllFiles(htmlContent string) error {
	// 保存index.html
	indexPath := filepath.Join(s.tempDir, "index.html")
	err := os.WriteFile(indexPath, []byte(htmlContent), 0644)
	if err != nil {
		return err
	}

	// 保存资源文件
	for _, resource := range s.resources {
		fullPath := filepath.Join(s.tempDir, resource.LocalPath)

		// 创建目录
		dir := filepath.Dir(fullPath)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			continue
		}

		// 保存文件
		err = os.WriteFile(fullPath, resource.Content, 0644)
		if err != nil {
			continue
		}
	}

	return nil
}

// createZipFile 创建ZIP文件
func (s *PageCaptureService) createZipFile() (string, int64, error) {
	zipFileName := fmt.Sprintf("webpage_%d.zip", time.Now().Unix())
	zipPath := filepath.Join(os.TempDir(), zipFileName)

	zipFile, err := os.Create(zipPath)
	if err != nil {
		return "", 0, err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(s.tempDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(s.tempDir, filePath)
		if err != nil {
			return err
		}

		zipFileWriter, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		_, err = zipFileWriter.Write(fileContent)
		return err
	})

	if err != nil {
		return "", 0, err
	}

	zipInfo, err := zipFile.Stat()
	if err != nil {
		return zipPath, 0, nil
	}

	return zipPath, zipInfo.Size(), nil
}

// formatHTML 格式化HTML代码
func (s *PageCaptureService) formatHTML(html string) string {
	// 简单的HTML格式化
	formatted := html

	// 添加适当的换行和缩进
	formatted = strings.ReplaceAll(formatted, "><", ">\n<")
	formatted = strings.ReplaceAll(formatted, "</head>", "</head>\n")
	formatted = strings.ReplaceAll(formatted, "</body>", "\n</body>")
	formatted = strings.ReplaceAll(formatted, "</html>", "\n</html>")

	// 处理常见标签的换行
	tags := []string{"div", "p", "h1", "h2", "h3", "h4", "h5", "h6", "ul", "ol", "li", "nav", "section", "article", "header", "footer", "main"}
	for _, tag := range tags {
		formatted = strings.ReplaceAll(formatted, fmt.Sprintf("</%s>", tag), fmt.Sprintf("</%s>\n", tag))
		formatted = strings.ReplaceAll(formatted, fmt.Sprintf("<%s", tag), fmt.Sprintf("\n<%s", tag))
	}

	// 清理多余的空行
	lines := strings.Split(formatted, "\n")
	var cleanLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			cleanLines = append(cleanLines, trimmed)
		}
	}

	// 添加基本缩进
	var indentedLines []string
	indent := 0
	for _, line := range cleanLines {
		trimmed := strings.TrimSpace(line)

		// 减少缩进（闭合标签）
		if strings.HasPrefix(trimmed, "</") && !strings.Contains(trimmed, "<meta") && !strings.Contains(trimmed, "<link") && !strings.Contains(trimmed, "<img") {
			indent--
			if indent < 0 {
				indent = 0
			}
		}

		// 添加缩进
		indentStr := strings.Repeat("  ", indent)
		indentedLines = append(indentedLines, indentStr+trimmed)

		// 增加缩进（开放标签）
		if strings.HasPrefix(trimmed, "<") && !strings.HasPrefix(trimmed, "</") &&
			!strings.Contains(trimmed, "<meta") && !strings.Contains(trimmed, "<link") &&
			!strings.Contains(trimmed, "<img") && !strings.HasSuffix(trimmed, "/>") {
			indent++
		}
	}

	return strings.Join(indentedLines, "\n")
}

// getFileList 获取文件列表
func (s *PageCaptureService) getFileList() []string {
	files := []string{"index.html"}

	for _, resource := range s.resources {
		files = append(files, resource.LocalPath)
	}

	return files
}
