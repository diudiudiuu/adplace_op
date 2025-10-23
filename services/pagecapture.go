package services

import (
	"archive/zip"
	"compress/gzip"
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
	"sync"
	"time"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// PageCaptureService 页面抓取服务
type PageCaptureService struct {
	client           *http.Client
	baseURL          *url.URL
	resources        map[string]*ResourceInfo
	tempDir          string
	maxFiles         int
	fileCount        int
	mutex            sync.RWMutex // 保护共享资源的互斥锁
	debug            bool         // 调试模式
	progressCallback ProgressCallback
	progressInfo     ProgressInfo
	progressMutex    sync.RWMutex
}

// NewPageCaptureService 创建新的页面抓取服务
func NewPageCaptureService() *PageCaptureService {
	// 创建自定义的HTTP传输
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		DisableCompression:  false, // 启用压缩
	}

	return &PageCaptureService{
		client: &http.Client{
			Timeout:   30 * time.Second,
			Transport: transport,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				// 允许最多10次重定向
				if len(via) >= 10 {
					return fmt.Errorf("重定向次数过多")
				}
				// 保持请求头
				for key, val := range via[0].Header {
					req.Header[key] = val
				}
				return nil
			},
		},
		resources: make(map[string]*ResourceInfo),
		debug:     false, // 默认关闭调试
	}
}

// SetDebug 设置调试模式
func (s *PageCaptureService) SetDebug(debug bool) {
	s.debug = debug
}

// SetProgressCallback 设置进度回调
func (s *PageCaptureService) SetProgressCallback(callback ProgressCallback) {
	s.progressCallback = callback
}

// GetCurrentProgress 获取当前进度信息
func (s *PageCaptureService) GetCurrentProgress() ProgressInfo {
	s.progressMutex.RLock()
	defer s.progressMutex.RUnlock()
	
	// 返回进度信息的副本
	progress := ProgressInfo{
		Phase:          s.progressInfo.Phase,
		TotalFiles:     s.progressInfo.TotalFiles,
		CompletedFiles: s.progressInfo.CompletedFiles,
		CurrentFile:    s.progressInfo.CurrentFile,
		FileProgress:   s.progressInfo.FileProgress,
		FileList:       make([]FileInfo, len(s.progressInfo.FileList)),
	}
	
	// 复制文件列表
	copy(progress.FileList, s.progressInfo.FileList)
	
	return progress
}

// updateProgress 更新进度信息
func (s *PageCaptureService) updateProgress(phase string, currentFile string, fileProgress int) {
	s.progressMutex.Lock()
	defer s.progressMutex.Unlock()
	
	s.progressInfo.Phase = phase
	s.progressInfo.CurrentFile = currentFile
	s.progressInfo.FileProgress = fileProgress
	
	if s.progressCallback != nil {
		s.progressCallback(s.progressInfo)
	}
}

// updateFileStatus 更新文件状态
func (s *PageCaptureService) updateFileStatus(url, status string, progress int) {
	s.progressMutex.Lock()
	defer s.progressMutex.Unlock()
	
	for i := range s.progressInfo.FileList {
		if s.progressInfo.FileList[i].URL == url {
			s.progressInfo.FileList[i].Status = status
			s.progressInfo.FileList[i].Progress = progress
			break
		}
	}
	
	// 更新完成文件数
	completed := 0
	for _, file := range s.progressInfo.FileList {
		if file.Status == "completed" {
			completed++
		}
	}
	s.progressInfo.CompletedFiles = completed
	
	if s.progressCallback != nil {
		s.progressCallback(s.progressInfo)
	}
}

// updateFileSize 更新文件大小
func (s *PageCaptureService) updateFileSize(url string, size int64) {
	s.progressMutex.Lock()
	defer s.progressMutex.Unlock()
	
	for i := range s.progressInfo.FileList {
		if s.progressInfo.FileList[i].URL == url {
			s.progressInfo.FileList[i].Size = s.formatFileSize(size)
			break
		}
	}
	
	if s.progressCallback != nil {
		s.progressCallback(s.progressInfo)
	}
}

// formatFileSize 格式化文件大小
func (s *PageCaptureService) formatFileSize(bytes int64) string {
	if bytes == 0 {
		return "0 B"
	}
	
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	
	sizes := []string{"B", "KB", "MB", "GB", "TB"}
	return fmt.Sprintf("%.1f %s", float64(bytes)/float64(div), sizes[exp])
}

// updateFileSize 更新文件大小
func (s *PageCaptureService) updateFileSize(url string, size int64) {
	s.progressMutex.Lock()
	defer s.progressMutex.Unlock()
	
	for i := range s.progressInfo.FileList {
		if s.progressInfo.FileList[i].URL == url {
			s.progressInfo.FileList[i].Size = s.formatFileSize(size)
			break
		}
	}
	
	if s.progressCallback != nil {
		s.progressCallback(s.progressInfo)
	}
}

// formatFileSize 格式化文件大小
func (s *PageCaptureService) formatFileSize(bytes int64) string {
	if bytes == 0 {
		return "0 B"
	}
	
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	
	sizes := []string{"B", "KB", "MB", "GB", "TB"}
	return fmt.Sprintf("%.1f %s", float64(bytes)/float64(div), sizes[exp])
}

// debugPrintf 调试输出函数
func (s *PageCaptureService) debugPrintf(format string, args ...interface{}) {
	if s.debug {
		fmt.Printf(format, args...)
	}
}

// ProgressCallback 进度回调函数类型
type ProgressCallback func(progress ProgressInfo)

// ProgressInfo 进度信息
type ProgressInfo struct {
	Phase         string     `json:"phase"`         // analyzing, downloading, saving, complete
	TotalFiles    int        `json:"totalFiles"`    
	CompletedFiles int       `json:"completedFiles"`
	CurrentFile   string     `json:"currentFile"`   
	FileProgress  int        `json:"fileProgress"`  // 0-100
	FileList      []FileInfo `json:"fileList"`      
}

// FileInfo 文件信息
type FileInfo struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Size     string `json:"size"`
	Status   string `json:"status"`   // pending, downloading, completed, failed
	Progress int    `json:"progress"` // 0-100
	URL      string `json:"url"`
}

// CaptureOptions 抓取选项
type CaptureOptions struct {
	IncludeImages       bool   `json:"includeImages"`
	IncludeStyles       bool   `json:"includeStyles"`
	IncludeScripts      bool   `json:"includeScripts"`
	FollowRedirects     bool   `json:"followRedirects"`
	IncludeFonts        bool   `json:"includeFonts"`
	IncludeVideos       bool   `json:"includeVideos"`
	RemoveAnalytics     bool   `json:"removeAnalytics"`
	RemoveTracking      bool   `json:"removeTracking"`
	RemoveAds           bool   `json:"removeAds"`
	RemoveTagManager    bool   `json:"removeTagManager"`
	RemoveMaliciousTags bool   `json:"removeMaliciousTags"`
	Timeout             int    `json:"timeout"`
	CreateZip           bool   `json:"createZip"`
	MaxFiles            int    `json:"maxFiles"`
	MaxDepth            int    `json:"maxDepth"`
	MaxConcurrency      int    `json:"maxConcurrency"`
	ForceEncoding       string `json:"forceEncoding"`
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
	FileDetails     []FileInfo          `json:"fileDetails,omitempty"`
	SuccessCount    int                 `json:"successCount"`
	FailedCount     int                 `json:"failedCount"`
}

// ResourceInfo 资源信息
type ResourceInfo struct {
	URL       string
	LocalPath string
	Type      string
	Content   []byte
}

// DownloadTask 下载任务
type DownloadTask struct {
	URL          string
	ResourceType string
	Element      *goquery.Selection
	AttrName     string
}

// DownloadResult 下载结果
type DownloadResult struct {
	Task      DownloadTask
	LocalPath string
	Success   bool
	Error     error
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
	
	// 临时启用调试模式来诊断问题
	s.debug = true

	// 设置超时 - 对于大文件需要更长时间
	timeout := options.Timeout
	if timeout < 120 {
		timeout = 120 // 最少2分钟
	}
	s.client.Timeout = time.Duration(timeout) * time.Second
	s.debugPrintf("设置HTTP超时: %d 秒\n", timeout)

	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "page_capture_*")
	if err != nil {
		return nil, fmt.Errorf("创建临时目录失败: %v", err)
	}
	s.tempDir = tempDir
	defer os.RemoveAll(tempDir)

	// 下载主页面（带重试机制）
	htmlContent, resp, err := s.downloadPageWithRetryAndEncoding(targetURL, 3, options.ForceEncoding)
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

	// 构建详细的文件信息
	s.progressMutex.RLock()
	fileDetails := make([]FileInfo, len(s.progressInfo.FileList))
	copy(fileDetails, s.progressInfo.FileList)
	successCount := s.progressInfo.CompletedFiles
	failedCount := len(s.progressInfo.FileList) - successCount
	s.progressMutex.RUnlock()

	// 构建结果
	result := &CaptureResult{
		StatusCode:      resp.StatusCode,
		ContentType:     resp.Header.Get("Content-Type"),
		ContentLength:   int64(len(htmlContent)),
		Content:         modifiedHTML, // 添加处理后的HTML内容
		Headers:         resp.Header,
		Duration:        time.Since(startTime).Milliseconds(),
		ZipPath:         zipPath,
		ZipSize:         zipSize,
		FilesCount:      len(s.resources) + 1,
		DownloadedFiles: s.getFileList(),
		FileDetails:     fileDetails,
		SuccessCount:    successCount,
		FailedCount:     failedCount,
	}

	return result, nil
}

// downloadPage 下载页面
func (s *PageCaptureService) downloadPage(targetURL string) (string, *http.Response, error) {
	// 创建请求
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return "", nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置更完整的请求头来模拟真实浏览器
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Sec-Ch-Ua", `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Windows"`)
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	// 如果URL包含Referer信息，设置Referer头
	if parsedURL, err := url.Parse(targetURL); err == nil {
		if parsedURL.Host != "" {
			req.Header.Set("Referer", fmt.Sprintf("%s://%s/", parsedURL.Scheme, parsedURL.Host))
		}
	}

	// 执行请求
	resp, err := s.client.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return "", resp, fmt.Errorf("HTTP错误: %d %s", resp.StatusCode, resp.Status)
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", resp, fmt.Errorf("读取响应失败: %v", err)
	}

	// 检查内容是否为空
	if len(body) == 0 {
		return "", resp, fmt.Errorf("响应内容为空")
	}

	// 使用简化的编码检测
	content := s.simpleEncodingDetection(body, resp.Header.Get("Content-Type"))

	// 检查是否返回了错误页面或重定向页面
	if strings.Contains(strings.ToLower(content), "access denied") ||
		strings.Contains(strings.ToLower(content), "forbidden") ||
		strings.Contains(strings.ToLower(content), "blocked") ||
		strings.Contains(strings.ToLower(content), "captcha") {
		return "", resp, fmt.Errorf("网站拒绝访问，可能存在反爬虫机制")
	}

	return content, resp, nil
}

// decompressIfNeeded 检查并解压缩内容
func (s *PageCaptureService) decompressIfNeeded(body []byte, contentEncoding string) ([]byte, error) {
	s.debugPrintf("Content-Encoding: %s\n", contentEncoding)
	s.debugPrintf("检查是否需要解压缩，前4字节: %v\n", body[:min(len(body), 4)])
	
	// 检查 GZIP 魔数 (0x1f 0x8b)
	if len(body) >= 2 && body[0] == 0x1f && body[1] == 0x8b {
		s.debugPrintf("检测到GZIP压缩内容，开始解压缩\n")
		
		reader, err := gzip.NewReader(strings.NewReader(string(body)))
		if err != nil {
			s.debugPrintf("创建GZIP读取器失败: %v\n", err)
			return body, err
		}
		defer reader.Close()
		
		decompressed, err := io.ReadAll(reader)
		if err != nil {
			s.debugPrintf("GZIP解压缩失败: %v\n", err)
			return body, err
		}
		
		s.debugPrintf("GZIP解压缩成功，原始大小: %d，解压后大小: %d\n", len(body), len(decompressed))
		s.debugPrintf("解压后前200字符: %s\n", string(decompressed[:min(len(decompressed), 200)]))
		return decompressed, nil
	}
	
	// 如果Content-Encoding指示了压缩但没有检测到魔数
	if strings.Contains(strings.ToLower(contentEncoding), "gzip") {
		s.debugPrintf("Content-Encoding指示GZIP但未检测到魔数\n")
	}
	
	s.debugPrintf("内容未压缩或不是GZIP格式\n")
	return body, nil
}

// readWithProgress 带进度监控的读取函数
func (s *PageCaptureService) readWithProgress(reader io.Reader, totalSize int64, url string) ([]byte, error) {
	buffer := make([]byte, 32*1024) // 32KB缓冲区
	var result []byte
	var downloaded int64
	lastProgress := 0
	
	for {
		n, err := reader.Read(buffer)
		if n > 0 {
			result = append(result, buffer[:n]...)
			downloaded += int64(n)
			
			// 每下载5%更新一次进度
			if totalSize > 0 {
				progress := int(downloaded * 100 / totalSize)
				if progress >= lastProgress+5 || progress == 100 {
					s.debugPrintf("下载进度: %s - %d%% (%.2f/%.2f MB)\n", 
						url, progress, 
						float64(downloaded)/(1024*1024), 
						float64(totalSize)/(1024*1024))
					s.updateFileStatus(url, "downloading", progress)
					lastProgress = progress
				}
			}
		}
		
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
	}
	
	return result, nil
}

// downloadLargeFileInChunks 分块并发下载大文件
func (s *PageCaptureService) downloadLargeFileInChunks(url string, totalSize int64, client *http.Client) ([]byte, error) {
	const chunkSize = 2 * 1024 * 1024 // 2MB per chunk
	const maxConcurrency = 4          // 最多4个并发
	
	numChunks := (totalSize + chunkSize - 1) / chunkSize
	s.debugPrintf("分块下载: %d 个块，每块 %.2f MB\n", numChunks, float64(chunkSize)/(1024*1024))
	
	type chunkResult struct {
		index int
		data  []byte
		err   error
	}
	
	chunks := make([][]byte, numChunks)
	resultChan := make(chan chunkResult, numChunks)
	semaphore := make(chan struct{}, maxConcurrency)
	
	var wg sync.WaitGroup
	
	// 启动下载协程
	for i := int64(0); i < numChunks; i++ {
		wg.Add(1)
		go func(chunkIndex int64) {
			defer wg.Done()
			semaphore <- struct{}{} // 获取信号量
			defer func() { <-semaphore }() // 释放信号量
			
			start := chunkIndex * chunkSize
			end := start + chunkSize - 1
			if end >= totalSize {
				end = totalSize - 1
			}
			
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				resultChan <- chunkResult{int(chunkIndex), nil, err}
				return
			}
			
			req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))
			req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
			
			resp, err := client.Do(req)
			if err != nil {
				resultChan <- chunkResult{int(chunkIndex), nil, err}
				return
			}
			defer resp.Body.Close()
			
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				resultChan <- chunkResult{int(chunkIndex), nil, err}
				return
			}
			
			s.debugPrintf("块 %d/%d 下载完成 (%.2f MB)\n", chunkIndex+1, numChunks, float64(len(data))/(1024*1024))
			resultChan <- chunkResult{int(chunkIndex), data, nil}
		}(i)
	}
	
	// 等待所有下载完成
	go func() {
		wg.Wait()
		close(resultChan)
	}()
	
	// 收集结果
	for result := range resultChan {
		if result.err != nil {
			return nil, fmt.Errorf("块 %d 下载失败: %v", result.index, result.err)
		}
		chunks[result.index] = result.data
	}
	
	// 合并所有块
	var finalData []byte
	for _, chunk := range chunks {
		finalData = append(finalData, chunk...)
	}
	
	s.debugPrintf("分块下载完成，总大小: %.2f MB\n", float64(len(finalData))/(1024*1024))
	return finalData, nil
}

// getFileNameFromURL 从URL提取文件名
func (s *PageCaptureService) getFileNameFromURL(urlStr string) string {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "unknown"
	}
	
	fileName := path.Base(parsedURL.Path)
	if fileName == "" || fileName == "." || fileName == "/" {
		return "index"
	}
	
	return fileName
}

// getEncodingByName 根据编码名称获取编码
func (s *PageCaptureService) getEncodingByName(name string) encoding.Encoding {
	name = strings.ToLower(strings.TrimSpace(name))
	
	switch name {
	case "utf-8", "utf8":
		return unicode.UTF8
	case "gbk", "gb2312", "gb18030":
		return simplifiedchinese.GBK
	case "big5":
		return traditionalchinese.Big5
	case "shift_jis", "shift-jis", "sjis":
		return japanese.ShiftJIS
	case "euc-jp":
		return japanese.EUCJP
	case "euc-kr":
		return korean.EUCKR
	case "iso-8859-1", "latin1":
		return charmap.ISO8859_1
	case "windows-1252", "cp1252":
		return charmap.Windows1252
	default:
		return nil
	}
}



// min 辅助函数
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// simpleEncodingDetection 简化的编码检测
func (s *PageCaptureService) simpleEncodingDetection(body []byte, contentType string) string {
	s.debugPrintf("开始编码检测，Content-Type: %s\n", contentType)
	s.debugPrintf("原始内容前100字节: %v\n", body[:min(len(body), 100)])
	
	// 首先尝试直接转换为字符串，看看是否包含HTML标签
	directStr := string(body)
	if strings.Contains(strings.ToLower(directStr), "<html") || 
	   strings.Contains(strings.ToLower(directStr), "<head") || 
	   strings.Contains(strings.ToLower(directStr), "<body") {
		s.debugPrintf("直接字符串转换包含HTML标签，使用原始内容\n")
		return directStr
	}

	// 检查是否已经是有效的UTF-8
	if utf8.Valid(body) {
		s.debugPrintf("内容已经是有效的UTF-8\n")
		return string(body)
	}

	s.debugPrintf("内容不是有效的UTF-8，开始编码转换\n")

	// 从Content-Type头中提取编码
	if contentType != "" {
		if strings.Contains(strings.ToLower(contentType), "charset=") {
			parts := strings.Split(strings.ToLower(contentType), "charset=")
			if len(parts) > 1 {
				charset := strings.TrimSpace(strings.Split(parts[1], ";")[0])
				s.debugPrintf("从Content-Type检测到编码: %s\n", charset)
				if encoding := s.getEncodingByName(charset); encoding != nil {
					decoder := encoding.NewDecoder()
					if result, _, err := transform.Bytes(decoder, body); err == nil && utf8.Valid(result) {
						resultStr := string(result)
						if strings.Contains(strings.ToLower(resultStr), "<html") || 
						   strings.Contains(strings.ToLower(resultStr), "<head") || 
						   strings.Contains(strings.ToLower(resultStr), "<body") {
							s.debugPrintf("使用%s编码转换成功\n", charset)
							return resultStr
						}
					}
				}
			}
		}
	}

	// 尝试从HTML内容中检测编码（使用原始字节）
	bodyStr := string(body[:min(len(body), 2048)])
	if strings.Contains(strings.ToLower(bodyStr), "charset=") {
		re := regexp.MustCompile(`charset\s*=\s*["']?([^"'\s>]+)`)
		matches := re.FindStringSubmatch(strings.ToLower(bodyStr))
		if len(matches) > 1 {
			charset := matches[1]
			s.debugPrintf("从HTML meta标签检测到编码: %s\n", charset)
			if encoding := s.getEncodingByName(charset); encoding != nil {
				decoder := encoding.NewDecoder()
				if result, _, err := transform.Bytes(decoder, body); err == nil && utf8.Valid(result) {
					resultStr := string(result)
					if strings.Contains(strings.ToLower(resultStr), "<html") || 
					   strings.Contains(strings.ToLower(resultStr), "<head") || 
					   strings.Contains(strings.ToLower(resultStr), "<body") {
						s.debugPrintf("使用%s编码转换成功\n", charset)
						return resultStr
					}
				}
			}
		}
	}

	// 尝试常见的编码
	encodings := []struct {
		name string
		enc  encoding.Encoding
	}{
		{"GBK", simplifiedchinese.GBK},
		{"Big5", traditionalchinese.Big5},
		{"ISO-8859-1", charmap.ISO8859_1},
		{"Windows-1252", charmap.Windows1252},
	}

	for _, item := range encodings {
		s.debugPrintf("尝试%s编码...\n", item.name)
		decoder := item.enc.NewDecoder()
		if result, _, err := transform.Bytes(decoder, body); err == nil && utf8.Valid(result) {
			resultStr := string(result)
			s.debugPrintf("%s编码转换后前200字符: %s\n", item.name, resultStr[:min(len(resultStr), 200)])
			// 检查转换后的内容是否包含合理的HTML标签
			if strings.Contains(strings.ToLower(resultStr), "<html") || 
			   strings.Contains(strings.ToLower(resultStr), "<head") || 
			   strings.Contains(strings.ToLower(resultStr), "<body") {
				s.debugPrintf("使用%s编码转换成功\n", item.name)
				return resultStr
			}
		} else {
			s.debugPrintf("%s编码转换失败: %v\n", item.name, err)
		}
	}

	// 最后降级为直接字符串转换
	s.debugPrintf("所有编码转换都失败，使用原始字符串\n")
	s.debugPrintf("原始字符串前200字符: %s\n", directStr[:min(len(directStr), 200)])
	return directStr
}

// downloadPageWithRetry 带重试机制的页面下载
func (s *PageCaptureService) downloadPageWithRetry(targetURL string, maxRetries int) (string, *http.Response, error) {
	var lastErr error
	
	for i := 0; i <= maxRetries; i++ {
		if i > 0 {
			// 重试前等待一段时间
			waitTime := time.Duration(i) * time.Second
			time.Sleep(waitTime)
		}
		
		content, resp, err := s.downloadPage(targetURL)
		if err == nil {
			return content, resp, nil
		}
		
		lastErr = err
		
		// 如果是某些特定错误，不进行重试
		if strings.Contains(err.Error(), "access denied") ||
			strings.Contains(err.Error(), "forbidden") ||
			strings.Contains(err.Error(), "反爬虫") {
			break
		}
	}
	
	return "", nil, fmt.Errorf("重试%d次后仍然失败: %v", maxRetries, lastErr)
}

// downloadPageWithRetryAndEncoding 带重试机制和编码选项的页面下载
func (s *PageCaptureService) downloadPageWithRetryAndEncoding(targetURL string, maxRetries int, forceEncoding string) (string, *http.Response, error) {
	var lastErr error
	
	for i := 0; i <= maxRetries; i++ {
		if i > 0 {
			// 重试前等待一段时间
			waitTime := time.Duration(i) * time.Second
			time.Sleep(waitTime)
		}
		
		content, resp, err := s.downloadPageWithEncoding(targetURL, forceEncoding)
		if err == nil {
			return content, resp, nil
		}
		
		lastErr = err
		
		// 如果是某些特定错误，不进行重试
		if strings.Contains(err.Error(), "access denied") ||
			strings.Contains(err.Error(), "forbidden") ||
			strings.Contains(err.Error(), "反爬虫") {
			break
		}
	}
	
	return "", nil, fmt.Errorf("重试%d次后仍然失败: %v", maxRetries, lastErr)
}

// downloadPageWithEncoding 支持编码选项的页面下载
func (s *PageCaptureService) downloadPageWithEncoding(targetURL string, forceEncoding string) (string, *http.Response, error) {
	// 创建请求
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return "", nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头（与原函数相同）
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Sec-Ch-Ua", `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Windows"`)
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	if parsedURL, err := url.Parse(targetURL); err == nil {
		if parsedURL.Host != "" {
			req.Header.Set("Referer", fmt.Sprintf("%s://%s/", parsedURL.Scheme, parsedURL.Host))
		}
	}

	// 执行请求
	resp, err := s.client.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return "", resp, fmt.Errorf("HTTP错误: %d %s", resp.StatusCode, resp.Status)
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", resp, fmt.Errorf("读取响应失败: %v", err)
	}

	// 检查内容是否为空
	if len(body) == 0 {
		return "", resp, fmt.Errorf("响应内容为空")
	}

	// 检查并解压缩 GZIP 内容
	body, err = s.decompressIfNeeded(body, resp.Header.Get("Content-Encoding"))
	if err != nil {
		return "", resp, fmt.Errorf("解压缩失败: %v", err)
	}

	// 简化的编码处理
	var content string
	if forceEncoding != "" && forceEncoding != "auto" {
		// 使用指定编码
		if encoding := s.getEncodingByName(forceEncoding); encoding != nil {
			decoder := encoding.NewDecoder()
			if result, _, err := transform.Bytes(decoder, body); err == nil {
				content = string(result)
			} else {
				content = string(body) // 降级处理
			}
		} else {
			content = string(body)
		}
	} else {
		// 自动检测编码
		content = s.simpleEncodingDetection(body, resp.Header.Get("Content-Type"))
	}

	// 检查是否返回了错误页面或重定向页面
	if strings.Contains(strings.ToLower(content), "access denied") ||
		strings.Contains(strings.ToLower(content), "forbidden") ||
		strings.Contains(strings.ToLower(content), "blocked") ||
		strings.Contains(strings.ToLower(content), "captcha") {
		return "", resp, fmt.Errorf("网站拒绝访问，可能存在反爬虫机制")
	}

	return content, resp, nil
}

// processHTMLAndDownloadResources 处理HTML并下载资源（并发版本）
func (s *PageCaptureService) processHTMLAndDownloadResources(htmlContent string, options CaptureOptions) (string, error) {
	s.debugPrintf("HTML内容长度: %d 字符\n", len(htmlContent))
	s.debugPrintf("HTML前500字符: %s\n", htmlContent[:min(len(htmlContent), 500)])
	
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		s.debugPrintf("HTML解析失败: %v\n", err)
		return htmlContent, err
	}

	s.debugPrintf("HTML解析成功\n")

	// 收集所有下载任务
	var tasks []DownloadTask

	s.debugPrintf("开始收集资源任务...\n")
	s.debugPrintf("选项: IncludeStyles=%v, IncludeScripts=%v, IncludeImages=%v, IncludeVideos=%v\n", 
		options.IncludeStyles, options.IncludeScripts, options.IncludeImages, options.IncludeVideos)
	s.debugPrintf("HTML文档解析成功，开始查找资源...\n")

	// 收集CSS文件下载任务
	if options.IncludeStyles {
		s.debugPrintf("收集CSS文件...\n")
		cssLinks := doc.Find("link[rel=stylesheet]")
		s.debugPrintf("找到 %d 个CSS link标签\n", cssLinks.Length())
		
		cssLinks.Each(func(i int, sel *goquery.Selection) {
			if href, exists := sel.Attr("href"); exists {
				// 清理URL中的空白字符
				href = strings.TrimSpace(href)
				s.debugPrintf("发现CSS: '%s'\n", href)
				absoluteURL := s.resolveURL(href)
				if absoluteURL != "" {
					s.debugPrintf("添加CSS任务: %s\n", absoluteURL)
					tasks = append(tasks, DownloadTask{
						URL:          absoluteURL,
						ResourceType: "css",
						Element:      sel,
						AttrName:     "href",
					})
				} else {
					s.debugPrintf("CSS URL解析失败: %s\n", href)
				}
			} else {
				s.debugPrintf("CSS link标签没有href属性\n")
			}
		})
		
		// 也检查所有link标签
		allLinks := doc.Find("link")
		s.debugPrintf("总共找到 %d 个link标签\n", allLinks.Length())
	} else {
		s.debugPrintf("跳过CSS文件收集 (IncludeStyles=false)\n")
	}

	// 收集JavaScript文件下载任务
	if options.IncludeScripts {
		s.debugPrintf("收集JavaScript文件...\n")
		jsScripts := doc.Find("script[src]")
		s.debugPrintf("找到 %d 个带src的script标签\n", jsScripts.Length())
		
		jsScripts.Each(func(i int, sel *goquery.Selection) {
			if src, exists := sel.Attr("src"); exists {
				// 清理URL中的空白字符
				src = strings.TrimSpace(src)
				s.debugPrintf("发现JS: '%s'\n", src)
				absoluteURL := s.resolveURL(src)
				if absoluteURL != "" {
					s.debugPrintf("添加JS任务: %s\n", absoluteURL)
					tasks = append(tasks, DownloadTask{
						URL:          absoluteURL,
						ResourceType: "js",
						Element:      sel,
						AttrName:     "src",
					})
				} else {
					s.debugPrintf("JS URL解析失败: %s\n", src)
				}
			}
		})
		
		// 也检查所有script标签
		allScripts := doc.Find("script")
		s.debugPrintf("总共找到 %d 个script标签\n", allScripts.Length())
	} else {
		s.debugPrintf("跳过JavaScript文件收集 (IncludeScripts=false)\n")
	}

	// 收集图片下载任务
	if options.IncludeImages {
		s.debugPrintf("收集图片文件...\n")
		imgTags := doc.Find("img[src]")
		s.debugPrintf("找到 %d 个带src的img标签\n", imgTags.Length())
		
		imgTags.Each(func(i int, sel *goquery.Selection) {
			if src, exists := sel.Attr("src"); exists {
				// 清理URL中的空白字符
				src = strings.TrimSpace(src)
				s.debugPrintf("发现图片: '%s'\n", src)
				absoluteURL := s.resolveURL(src)
				if absoluteURL != "" {
					s.debugPrintf("添加图片任务: %s\n", absoluteURL)
					tasks = append(tasks, DownloadTask{
						URL:          absoluteURL,
						ResourceType: "images",
						Element:      sel,
						AttrName:     "src",
					})
				} else {
					s.debugPrintf("图片URL解析失败: %s\n", src)
				}
			}
		})
		
		// 也检查所有img标签
		allImages := doc.Find("img")
		s.debugPrintf("总共找到 %d 个img标签\n", allImages.Length())
	} else {
		s.debugPrintf("跳过图片文件收集 (IncludeImages=false)\n")
	}

	// 收集视频下载任务
	if options.IncludeVideos {
		s.debugPrintf("收集视频文件...\n")
		
		// video[src]
		videoTags := doc.Find("video[src]")
		s.debugPrintf("找到 %d 个带src的video标签\n", videoTags.Length())
		videoTags.Each(func(i int, sel *goquery.Selection) {
			if src, exists := sel.Attr("src"); exists {
				// 清理URL中的空白字符
				src = strings.TrimSpace(src)
				s.debugPrintf("发现视频: '%s'\n", src)
				absoluteURL := s.resolveURL(src)
				if absoluteURL != "" {
					s.debugPrintf("添加视频任务: %s\n", absoluteURL)
					tasks = append(tasks, DownloadTask{
						URL:          absoluteURL,
						ResourceType: "videos",
						Element:      sel,
						AttrName:     "src",
					})
				} else {
					s.debugPrintf("视频URL解析失败: %s\n", src)
				}
			}
		})

		// video source[src]
		videoSources := doc.Find("video source[src]")
		s.debugPrintf("找到 %d 个video source标签\n", videoSources.Length())
		videoSources.Each(func(i int, sel *goquery.Selection) {
			if src, exists := sel.Attr("src"); exists {
				// 清理URL中的空白字符
				src = strings.TrimSpace(src)
				s.debugPrintf("发现视频源: '%s'\n", src)
				absoluteURL := s.resolveURL(src)
				if absoluteURL != "" {
					s.debugPrintf("添加视频源任务: %s\n", absoluteURL)
					tasks = append(tasks, DownloadTask{
						URL:          absoluteURL,
						ResourceType: "videos",
						Element:      sel,
						AttrName:     "src",
					})
				} else {
					s.debugPrintf("视频源URL解析失败: %s\n", src)
				}
			}
		})

		// audio[src]
		audioTags := doc.Find("audio[src]")
		s.debugPrintf("找到 %d 个带src的audio标签\n", audioTags.Length())
		audioTags.Each(func(i int, sel *goquery.Selection) {
			if src, exists := sel.Attr("src"); exists {
				// 清理URL中的空白字符
				src = strings.TrimSpace(src)
				s.debugPrintf("发现音频: '%s'\n", src)
				absoluteURL := s.resolveURL(src)
				if absoluteURL != "" {
					s.debugPrintf("添加音频任务: %s\n", absoluteURL)
					tasks = append(tasks, DownloadTask{
						URL:          absoluteURL,
						ResourceType: "videos",
						Element:      sel,
						AttrName:     "src",
					})
				} else {
					s.debugPrintf("音频URL解析失败: %s\n", src)
				}
			}
		})

		// audio source[src]
		audioSources := doc.Find("audio source[src]")
		s.debugPrintf("找到 %d 个audio source标签\n", audioSources.Length())
		audioSources.Each(func(i int, sel *goquery.Selection) {
			if src, exists := sel.Attr("src"); exists {
				// 清理URL中的空白字符
				src = strings.TrimSpace(src)
				s.debugPrintf("发现音频源: '%s'\n", src)
				absoluteURL := s.resolveURL(src)
				if absoluteURL != "" {
					s.debugPrintf("添加音频源任务: %s\n", absoluteURL)
					tasks = append(tasks, DownloadTask{
						URL:          absoluteURL,
						ResourceType: "videos",
						Element:      sel,
						AttrName:     "src",
					})
				} else {
					s.debugPrintf("音频源URL解析失败: %s\n", src)
				}
			}
		})
	} else {
		s.debugPrintf("跳过视频文件收集 (IncludeVideos=false)\n")
	}

	// 收集字体文件下载任务
	if options.IncludeFonts {
		doc.Find("link[rel='preload'][as='font']").Each(func(i int, sel *goquery.Selection) {
			if href, exists := sel.Attr("href"); exists {
				absoluteURL := s.resolveURL(href)
				if absoluteURL != "" {
					tasks = append(tasks, DownloadTask{
						URL:          absoluteURL,
						ResourceType: "fonts",
						Element:      sel,
						AttrName:     "href",
					})
				}
			}
		})
	}

	// 并发下载所有资源
	maxConcurrency := options.MaxConcurrency
	if maxConcurrency <= 0 {
		maxConcurrency = 10 // 默认10个并发
	}
	
	// 初始化进度信息
	s.progressMutex.Lock()
	s.progressInfo.Phase = "downloading"
	s.progressInfo.TotalFiles = len(tasks)
	s.progressInfo.CompletedFiles = 0
	s.progressInfo.FileList = make([]FileInfo, len(tasks))
	
	for i, task := range tasks {
		fileName := s.getFileNameFromURL(task.URL)
		s.progressInfo.FileList[i] = FileInfo{
			Name:     fileName,
			Type:     task.ResourceType,
			Size:     "等待下载...",
			Status:   "pending",
			Progress: 0,
			URL:      task.URL,
		}
	}
	s.progressMutex.Unlock()
	
	s.debugPrintf("初始化文件列表完成，共 %d 个文件\n", len(tasks))
	
	// 如果没有找到任何资源文件，至少添加主页面到文件列表
	if len(tasks) == 0 {
		s.debugPrintf("没有找到资源文件，添加主页面到文件列表\n")
		s.progressInfo.FileList = []FileInfo{
			{
				Name:     "index.html",
				Type:     "html",
				Size:     s.formatFileSize(int64(len(htmlContent))),
				Status:   "completed",
				Progress: 100,
				URL:      s.baseURL.String(),
			},
		}
		s.progressInfo.TotalFiles = 1
		s.progressInfo.CompletedFiles = 1
	}
	
	s.updateProgress("downloading", "开始下载资源文件...", 0)
	
	s.debugPrintf("开始并发下载 %d 个资源，并发数: %d\n", len(tasks), maxConcurrency)
	results := s.downloadResourcesConcurrently(tasks, maxConcurrency)

	// 统计结果
	successCount := 0
	for _, result := range results {
		if result.Success && result.LocalPath != "" {
			result.Task.Element.SetAttr(result.Task.AttrName, result.LocalPath)
			successCount++
		}
	}
	
	s.debugPrintf("资源下载完成: 成功 %d/%d 个\n", successCount, len(tasks))

	// 处理CSS中的背景图片和字体（这些需要特殊处理，暂时保持同步）
	if options.IncludeImages || options.IncludeFonts {
		doc.Find("style").Each(func(i int, sel *goquery.Selection) {
			cssContent := sel.Text()
			if options.IncludeImages {
				cssContent = s.processCSSContent(cssContent)
			}
			if options.IncludeFonts {
				cssContent = s.processFontContent(cssContent)
			}
			sel.SetText(cssContent)
		})
	}

	html, err := doc.Html()
	if err != nil {
		return htmlContent, err
	}

	// 删除第三方跟踪代码
	if options.RemoveAnalytics || options.RemoveTracking || options.RemoveAds || options.RemoveTagManager || options.RemoveMaliciousTags {
		html = s.removeThirdPartyCode(html, options)
	}

	// 格式化HTML
	formattedHTML := s.formatHTML(html)
	return formattedHTML, nil
}

// downloadResourcesConcurrently 并发下载资源
func (s *PageCaptureService) downloadResourcesConcurrently(tasks []DownloadTask, maxConcurrency int) []DownloadResult {
	if len(tasks) == 0 {
		return nil
	}

	// 限制任务数量
	if len(tasks) > s.maxFiles {
		tasks = tasks[:s.maxFiles]
	}

	// 创建通道
	taskChan := make(chan DownloadTask, len(tasks))
	resultChan := make(chan DownloadResult, len(tasks))

	// 启动工作协程
	var wg sync.WaitGroup
	for i := 0; i < maxConcurrency && i < len(tasks); i++ {
		wg.Add(1)
		go s.downloadWorker(taskChan, resultChan, &wg)
	}

	// 发送任务
	for _, task := range tasks {
		taskChan <- task
	}
	close(taskChan)

	// 等待所有工作协程完成
	wg.Wait()
	close(resultChan)

	// 收集结果
	var results []DownloadResult
	for result := range resultChan {
		results = append(results, result)
	}

	return results
}

// downloadWorker 下载工作协程
func (s *PageCaptureService) downloadWorker(taskChan <-chan DownloadTask, resultChan chan<- DownloadResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range taskChan {
		// 检查文件数量限制
		s.mutex.RLock()
		currentCount := s.fileCount
		s.mutex.RUnlock()

		if currentCount >= s.maxFiles {
			resultChan <- DownloadResult{
				Task:    task,
				Success: false,
				Error:   fmt.Errorf("达到最大文件数限制"),
			}
			continue
		}

		// 检查是否已下载
		s.mutex.RLock()
		if resource, exists := s.resources[task.URL]; exists {
			s.mutex.RUnlock()
			resultChan <- DownloadResult{
				Task:      task,
				LocalPath: resource.LocalPath,
				Success:   true,
			}
			continue
		}
		s.mutex.RUnlock()

		// 下载资源
		localPath := s.downloadResourceSync(task.URL, task.ResourceType)
		
		resultChan <- DownloadResult{
			Task:      task,
			LocalPath: localPath,
			Success:   localPath != "",
		}
	}
}

// downloadResourceSync 同步下载单个资源（线程安全版本）
func (s *PageCaptureService) downloadResourceSync(resourceURL, resourceType string) string {
	// 检查是否已下载（只读锁）
	s.mutex.RLock()
	if resource, exists := s.resources[resourceURL]; exists {
		s.mutex.RUnlock()
		return resource.LocalPath
	}
	
	// 检查文件数量限制
	if s.fileCount >= s.maxFiles {
		s.mutex.RUnlock()
		return ""
	}
	s.mutex.RUnlock()

	// 在锁外进行网络下载（这是耗时操作）
	s.debugPrintf("开始下载: %s (%s)\n", resourceURL, resourceType)
	s.updateFileStatus(resourceURL, "downloading", 0)
	
	// 为大文件创建更长超时的客户端
	client := s.client
	if resourceType == "videos" {
		// 视频文件使用更长的超时时间（10分钟）
		client = &http.Client{
			Timeout:   10 * time.Minute,
			Transport: s.client.Transport,
		}
		s.debugPrintf("使用扩展超时(10分钟)下载视频文件\n")
	}
	
	resp, err := client.Get(resourceURL)
	if err != nil {
		s.debugPrintf("下载失败: %s - %v\n", resourceURL, err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.debugPrintf("HTTP错误: %s - %d\n", resourceURL, resp.StatusCode)
		return ""
	}

	// 获取文件大小信息
	contentLength := resp.ContentLength
	if contentLength > 0 {
		s.debugPrintf("下载中: %s - 大小: %.2f MB\n", resourceURL, float64(contentLength)/(1024*1024))
		// 更新文件大小信息
		s.updateFileSize(resourceURL, contentLength)
	}

	// 为大文件添加分块并发下载
	var content []byte
	if contentLength > 10*1024*1024 && resourceType == "videos" {
		s.debugPrintf("大视频文件，启用分块并发下载...\n")
		resp.Body.Close() // 关闭当前连接
		content, err = s.downloadLargeFileInChunks(resourceURL, contentLength, client)
	} else if contentLength > 5*1024*1024 {
		s.debugPrintf("大文件下载，启用进度监控...\n")
		content, err = s.readWithProgress(resp.Body, contentLength, resourceURL)
	} else {
		content, err = io.ReadAll(resp.Body)
	}
	
	if err != nil {
		s.debugPrintf("读取失败: %s - %v\n", resourceURL, err)
		s.updateFileStatus(resourceURL, "failed", 0)
		return ""
	}

	s.debugPrintf("下载完成: %s - 实际大小: %.2f MB\n", resourceURL, float64(len(content))/(1024*1024))
	s.updateFileStatus(resourceURL, "completed", 100)

	// 对于CSS和JS文件，尝试处理编码问题
	if resourceType == "css" || resourceType == "js" {
		if !utf8.Valid(content) {
			// 使用简化的编码检测
			convertedContent := s.simpleEncodingDetection(content, resp.Header.Get("Content-Type"))
			content = []byte(convertedContent)
		}
	}

	// 生成本地路径
	localPath := s.generateLocalPath(resourceURL, resourceType)

	// 获取写锁来存储结果（只锁定存储操作）
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 再次检查是否已被其他协程下载
	if resource, exists := s.resources[resourceURL]; exists {
		return resource.LocalPath
	}

	// 再次检查文件数量限制
	if s.fileCount >= s.maxFiles {
		return ""
	}

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

// downloadResource 下载单个资源（保持向后兼容）
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
	s.debugPrintf("解析URL: %s (基础URL: %s)\n", resourceURL, s.baseURL.String())
	
	// 清理资源URL - 移除所有前后空白字符
	resourceURL = strings.TrimSpace(resourceURL)
	s.debugPrintf("清理后的URL: '%s'\n", resourceURL)
	
	if strings.HasPrefix(resourceURL, "http://") || strings.HasPrefix(resourceURL, "https://") {
		s.debugPrintf("绝对URL: %s\n", resourceURL)
		return resourceURL
	}

	if strings.HasPrefix(resourceURL, "//") {
		result := s.baseURL.Scheme + ":" + resourceURL
		s.debugPrintf("协议相对URL: %s\n", result)
		return result
	}

	if strings.HasPrefix(resourceURL, "/") {
		result := s.baseURL.Scheme + "://" + s.baseURL.Host + resourceURL
		s.debugPrintf("根相对URL: %s\n", result)
		return result
	}

	// 相对路径处理
	baseDir := path.Dir(s.baseURL.Path)
	if baseDir == "." || baseDir == "/" {
		baseDir = ""
	}
	
	// 处理 ./ 开头的相对路径
	if strings.HasPrefix(resourceURL, "./") {
		resourceURL = resourceURL[2:] // 移除 "./"
	}
	
	// 确保路径正确拼接
	var result string
	if baseDir == "" {
		result = s.baseURL.Scheme + "://" + s.baseURL.Host + "/" + resourceURL
	} else {
		result = s.baseURL.Scheme + "://" + s.baseURL.Host + baseDir + "/" + resourceURL
	}
	
	s.debugPrintf("相对路径解析: %s -> %s (baseDir: %s)\n", resourceURL, result, baseDir)
	return result
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
	case "videos":
		return ".mp4"
	case "fonts":
		return ".woff2"
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

// processFontContent 处理CSS内容中的字体文件URL
func (s *PageCaptureService) processFontContent(cssContent string) string {
	if s.fileCount >= s.maxFiles {
		return cssContent
	}

	// 处理 @font-face 规则中的 src 属性
	fontFaceRegex := regexp.MustCompile(`@font-face\s*\{[^}]*src\s*:\s*([^;}]+)[;}]`)
	urlRegex := regexp.MustCompile(`url\s*\(\s*['"]?([^'")]+)['"]?\s*\)`)

	return fontFaceRegex.ReplaceAllStringFunc(cssContent, func(match string) string {
		return urlRegex.ReplaceAllStringFunc(match, func(urlMatch string) string {
			submatches := urlRegex.FindStringSubmatch(urlMatch)
			if len(submatches) > 1 {
				resourceURL := submatches[1]
				// 检查是否是字体文件
				if s.isFontFile(resourceURL) {
					absoluteURL := s.resolveURL(resourceURL)
					if absoluteURL != "" {
						localPath := s.downloadResource(absoluteURL, "fonts")
						if localPath != "" {
							return strings.Replace(urlMatch, resourceURL, localPath, 1)
						}
					}
				}
			}
			return urlMatch
		})
	})
}

// isFontFile 检查是否是字体文件
func (s *PageCaptureService) isFontFile(url string) bool {
	fontExtensions := []string{".woff", ".woff2", ".ttf", ".otf", ".eot"}
	lowerURL := strings.ToLower(url)
	for _, ext := range fontExtensions {
		if strings.Contains(lowerURL, ext) {
			return true
		}
	}
	return false
}

// saveAllFiles 保存所有文件
func (s *PageCaptureService) saveAllFiles(htmlContent string) error {
	s.debugPrintf("开始保存文件到: %s\n", s.tempDir)
	
	// 保存index.html
	indexPath := filepath.Join(s.tempDir, "index.html")
	err := os.WriteFile(indexPath, []byte(htmlContent), 0644)
	if err != nil {
		return fmt.Errorf("保存index.html失败: %v", err)
	}
	s.debugPrintf("保存index.html成功: %s\n", indexPath)

	// 保存资源文件
	s.debugPrintf("开始保存 %d 个资源文件\n", len(s.resources))
	savedCount := 0
	for url, resource := range s.resources {
		fullPath := filepath.Join(s.tempDir, resource.LocalPath)
		s.debugPrintf("保存资源: %s -> %s\n", url, fullPath)

		// 创建目录
		dir := filepath.Dir(fullPath)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			s.debugPrintf("创建目录失败: %s - %v\n", dir, err)
			continue // 跳过无法创建目录的文件
		}

		// 保存文件
		err = os.WriteFile(fullPath, resource.Content, 0644)
		if err != nil {
			s.debugPrintf("保存文件失败: %s - %v\n", fullPath, err)
			continue // 跳过无法保存的文件
		}
		savedCount++
	}
	
	s.debugPrintf("资源文件保存完成: %d/%d 个成功\n", savedCount, len(s.resources))
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

// removeThirdPartyCode 删除第三方跟踪代码
func (s *PageCaptureService) removeThirdPartyCode(html string, options CaptureOptions) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return html
	}

	// 删除统计分析代码
	if options.RemoveAnalytics {
		s.removeAnalyticsCode(doc)
	}

	// 删除跟踪代码
	if options.RemoveTracking {
		s.removeTrackingCode(doc)
	}

	// 删除广告代码
	if options.RemoveAds {
		s.removeAdsCode(doc)
	}

	// 删除标签管理器
	if options.RemoveTagManager {
		s.removeTagManagerCode(doc)
	}

	// 删除恶意标签
	if options.RemoveMaliciousTags {
		s.removeMaliciousTags(doc)
	}

	result, err := doc.Html()
	if err != nil {
		return html
	}

	return result
}

// removeAnalyticsCode 删除统计分析代码
func (s *PageCaptureService) removeAnalyticsCode(doc *goquery.Document) {
	// Google Analytics 相关
	analyticsSelectors := []string{
		"script[src*='google-analytics.com']",
		"script[src*='googletagmanager.com/gtag']",
		"script[src*='gtag/js']",
		"script[src*='analytics.js']",
		"script[src*='gtag.js']",
		"script[src*='ga.js']",
		// 百度统计
		"script[src*='hm.baidu.com']",
		"script[src*='cnzz.com']",
		// Mixpanel
		"script[src*='mixpanel.com']",
		// Segment
		"script[src*='segment.com']",
		"script[src*='segment.io']",
	}

	for _, selector := range analyticsSelectors {
		doc.Find(selector).Remove()
	}

	// 删除包含分析代码的内联脚本
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		content := s.Text()
		if containsAnalyticsCode(content) {
			s.Remove()
		}
	})
}

// removeTrackingCode 删除跟踪代码
func (s *PageCaptureService) removeTrackingCode(doc *goquery.Document) {
	// Facebook Pixel, TikTok Pixel, Hotjar 等
	trackingSelectors := []string{
		"script[src*='connect.facebook.net']",
		"script[src*='analytics.tiktok.com']",
		"script[src*='snapchat.com/web-sdk']",
		"script[src*='hotjar.com']",
		"script[src*='crazyegg.com']",
		"script[src*='clarity.ms']",
		"script[src*='mouseflow.com']",
		"script[src*='fullstory.com']",
	}

	for _, selector := range trackingSelectors {
		doc.Find(selector).Remove()
	}

	// 删除包含跟踪代码的内联脚本
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		content := s.Text()
		if containsTrackingCode(content) {
			s.Remove()
		}
	})
}

// removeAdsCode 删除广告代码
func (s *PageCaptureService) removeAdsCode(doc *goquery.Document) {
	// Google Ads, Taboola, PopAds 等
	adsSelectors := []string{
		"script[src*='googlesyndication.com']",
		"script[src*='doubleclick.net']",
		"script[src*='taboola.com']",
		"script[src*='outbrain.com']",
		"script[src*='popads.net']",
		"script[src*='propellerads.com']",
		"script[src*='adcash.com']",
		"script[src*='affiliate.js']",
		"script[src*='redirect.js']",
		"ins.adsbygoogle",
	}

	for _, selector := range adsSelectors {
		doc.Find(selector).Remove()
	}

	// 删除包含广告代码的内联脚本
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		content := s.Text()
		if containsAdsCode(content) {
			s.Remove()
		}
	})
}

// removeTagManagerCode 删除标签管理器代码
func (s *PageCaptureService) removeTagManagerCode(doc *goquery.Document) {
	// Google Tag Manager
	tagManagerSelectors := []string{
		"script[src*='googletagmanager.com/gtm.js']",
		"noscript iframe[src*='googletagmanager.com']",
	}

	for _, selector := range tagManagerSelectors {
		doc.Find(selector).Remove()
	}

	// 删除包含GTM代码的内联脚本
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		content := s.Text()
		if containsTagManagerCode(content) {
			s.Remove()
		}
	})
}

// containsAnalyticsCode 检查是否包含分析代码
func containsAnalyticsCode(content string) bool {
	analyticsKeywords := []string{
		"google-analytics.com",
		"gtag(",
		"ga(",
		"_gaq",
		"GoogleAnalyticsObject",
		"hm.baidu.com",
		"_hmt",
		"mixpanel",
		"analytics.track",
		"segment.com",
	}

	contentLower := strings.ToLower(content)
	for _, keyword := range analyticsKeywords {
		if strings.Contains(contentLower, strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}

// containsTrackingCode 检查是否包含跟踪代码
func containsTrackingCode(content string) bool {
	trackingKeywords := []string{
		"fbq(",
		"facebook.net",
		"ttq.track",
		"tiktok",
		"snaptr(",
		"hotjar",
		"hj(",
		"crazyegg",
		"clarity",
		"mouseflow",
		"fullstory",
	}

	contentLower := strings.ToLower(content)
	for _, keyword := range trackingKeywords {
		if strings.Contains(contentLower, strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}

// containsAdsCode 检查是否包含广告代码
func containsAdsCode(content string) bool {
	adsKeywords := []string{
		"googlesyndication",
		"adsbygoogle",
		"doubleclick",
		"taboola",
		"outbrain",
		"popads",
		"propellerads",
		"adcash",
		"affiliate",
		"redirect",
	}

	contentLower := strings.ToLower(content)
	for _, keyword := range adsKeywords {
		if strings.Contains(contentLower, strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}

// containsTagManagerCode 检查是否包含标签管理器代码
func containsTagManagerCode(content string) bool {
	tagManagerKeywords := []string{
		"googletagmanager.com",
		"dataLayer",
		"GTM-",
	}

	contentLower := strings.ToLower(content)
	for _, keyword := range tagManagerKeywords {
		if strings.Contains(contentLower, strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}

// removeMaliciousTags 删除可能被恶意利用的HTML标签
func (s *PageCaptureService) removeMaliciousTags(doc *goquery.Document) {
	// 删除 <base> 标签 - 可能被用于劫持所有相对链接
	doc.Find("base").Remove()

	// 删除 <meta http-equiv="refresh"> 标签 - 可能被用于自动跳转到恶意网站
	doc.Find("meta[http-equiv='refresh']").Remove()
	doc.Find("meta[http-equiv='Refresh']").Remove()
	doc.Find("meta[http-equiv='REFRESH']").Remove()

	// 删除 <meta name="referrer"> 标签 - 可能被用于伪造来源
	doc.Find("meta[name='referrer']").Remove()
	doc.Find("meta[name='Referrer']").Remove()
	doc.Find("meta[name='REFERRER']").Remove()

	// 删除其他可能的恶意 meta 标签
	// 删除可能用于重定向的 meta 标签
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		httpEquiv, exists := s.Attr("http-equiv")
		if exists {
			httpEquivLower := strings.ToLower(httpEquiv)
			// 删除各种可能的重定向相关的 meta 标签
			if httpEquivLower == "refresh" || httpEquivLower == "location" || httpEquivLower == "redirect" {
				s.Remove()
			}
		}

		name, exists := s.Attr("name")
		if exists {
			nameLower := strings.ToLower(name)
			// 删除可能影响安全的 meta 标签
			if nameLower == "referrer" || nameLower == "redirect" || nameLower == "location" {
				s.Remove()
			}
		}

		content, exists := s.Attr("content")
		if exists {
			contentLower := strings.ToLower(content)
			// 检查 content 属性中是否包含可疑的跳转指令
			if strings.Contains(contentLower, "url=") && (strings.Contains(contentLower, "http://") || strings.Contains(contentLower, "https://")) {
				// 如果 content 包含 URL 跳转，可能是恶意的
				httpEquiv, hasHttpEquiv := s.Attr("http-equiv")
				if hasHttpEquiv && strings.ToLower(httpEquiv) == "refresh" {
					s.Remove()
				}
			}
		}
	})

	// 删除可能包含恶意跳转的内联脚本
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		content := s.Text()
		if containsMaliciousCode(content) {
			s.Remove()
		}
	})
}

// containsMaliciousCode 检查是否包含恶意代码
func containsMaliciousCode(content string) bool {
	maliciousKeywords := []string{
		"window.location.href",
		"window.location.replace",
		"window.location.assign",
		"document.location.href",
		"document.location.replace",
		"location.href",
		"location.replace",
		"top.location",
		"parent.location",
		// 检查是否有可疑的重定向模式
		"setTimeout.*location",
		"setInterval.*location",
		// 检查是否有 base 标签操作
		"document.createElement.*base",
		"<base",
		// 检查是否有动态创建 meta refresh
		"http-equiv.*refresh",
		"meta.*refresh",
	}

	contentLower := strings.ToLower(content)
	for _, keyword := range maliciousKeywords {
		if strings.Contains(contentLower, strings.ToLower(keyword)) {
			return true
		}
	}

	// 使用正则表达式检查更复杂的模式
	maliciousPatterns := []string{
		`location\s*=\s*['"][^'"]*['"]`, // location = "url"
		`location\s*\.\s*href\s*=`,      // location.href =
		`window\s*\.\s*open\s*\(`,       // window.open(
		`document\s*\.\s*write.*<base`,  // document.write 包含 base 标签
		`setTimeout\s*\(.*location`,     // setTimeout 中包含 location
		`setInterval\s*\(.*location`,    // setInterval 中包含 location
	}

	for _, pattern := range maliciousPatterns {
		matched, err := regexp.MatchString(pattern, contentLower)
		if err == nil && matched {
			return true
		}
	}

	return false
}
