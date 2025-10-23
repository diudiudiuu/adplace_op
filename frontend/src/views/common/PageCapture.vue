<template>
    <div>
        <n-card title="🔍 页面抓取工具">
            <template #header-extra>
                <n-tooltip>
                    <template #trigger>
                        <n-tag type="info">完整下载</n-tag>
                    </template>
                    原封不动地下载整个网页，包括所有资源文件，并打包成ZIP
                </n-tooltip>
            </template>
            <n-space vertical size="large">
                <!-- 功能说明 -->
                <n-alert type="info" title="功能说明" closable>
                    <n-ul>
                        <n-li>📄 <strong>完整下载</strong>：原封不动地下载整个网页，包括HTML、CSS、JavaScript、图片等所有资源</n-li>
                        <n-li>📦 <strong>标准结构</strong>：生成标准的静态网站结构，index.html + static资源目录</n-li>
                        <n-li>🔗 <strong>链接修正</strong>：自动修正页面中的资源链接，确保离线浏览正常</n-li>
                        <n-li>⚙️ <strong>灵活配置</strong>：可选择包含或排除特定类型的资源文件</n-li>
                    </n-ul>
                </n-alert>

                <!-- 操作按钮 -->
                <n-space justify="center">
                    <n-tooltip>
                        <template #trigger>
                            <n-button type="primary" size="large" @click="captureUrl" :disabled="!form.url.trim()">
                                <template #icon>
                                    <n-icon>
                                        <CameraOutline />
                                    </n-icon>
                                </template>
                                开始抓取
                            </n-button>
                        </template>
                        抓取指定 URL 的页面内容
                    </n-tooltip>
                    <n-tooltip>
                        <template #trigger>
                            <n-button type="info" @click="clearResults">
                                <template #icon>
                                    <n-icon>
                                        <RefreshOutline />
                                    </n-icon>
                                </template>
                                清空结果
                            </n-button>
                        </template>
                        清空抓取结果
                    </n-tooltip>
                </n-space>
                <!-- 抓取配置 -->
                <n-card size="small" title="抓取配置">
                    <template #header-extra>
                        <n-tag type="success" size="small">已优化默认配置</n-tag>
                    </template>
                    <n-form :model="form" label-placement="left" label-width="120">
                        <n-form-item label="目标 URL" required>
                            <n-input v-model:value="form.url" placeholder="请输入要抓取的网页 URL，如：https://example.com"
                                @keyup.enter="captureUrl" />
                            <template #suffix>
                                <n-space>
                                    <n-dropdown :options="urlOptions" @select="selectUrl">
                                        <n-button text>
                                            <template #icon>
                                                <n-icon>
                                                    <ChevronDownOutline />
                                                </n-icon>
                                            </template>
                                        </n-button>
                                    </n-dropdown>
                                    <n-button text type="primary" @click="testConnection" :disabled="!form.url.trim()">
                                        测试连接
                                    </n-button>
                                </n-space>
                            </template>
                        </n-form-item>
                        <n-form-item label="抓取选项">
                            <n-space vertical>
                                <n-checkbox v-model:checked="options.includeImages">
                                    包含图片
                                </n-checkbox>
                                <n-checkbox v-model:checked="options.includeStyles">
                                    包含样式
                                </n-checkbox>
                                <n-checkbox v-model:checked="options.includeScripts">
                                    包含脚本
                                </n-checkbox>
                                <n-checkbox v-model:checked="options.followRedirects">
                                    跟随重定向
                                </n-checkbox>

                            </n-space>
                        </n-form-item>
                        <n-form-item label="超时时间">
                            <n-input-number v-model:value="options.timeout" :min="10" :max="180" :step="10"
                                placeholder="秒" />
                            <template #suffix>秒</template>
                        </n-form-item>
                        <n-form-item label="最大文件数">
                            <n-input-number v-model:value="options.maxFiles" :min="50" :max="1000" :step="50"
                                placeholder="个" />
                            <template #suffix>个</template>
                        </n-form-item>
                        <n-form-item label="调试模式">
                            <n-checkbox v-model:checked="debugMode">
                                显示详细错误信息
                            </n-checkbox>
                        </n-form-item>
                    </n-form>
                </n-card>

                <!-- 抓取结果 -->
                <n-card v-if="captureResult" size="small" title="抓取结果">
                    <template #header-extra>
                        <n-space>
                            <n-tag :type="captureResult.success ? 'success' : 'error'">
                                {{ captureResult.success ? '成功' : '失败' }}
                            </n-tag>
                            <n-tag type="info">
                                {{ captureResult.timestamp }}
                            </n-tag>
                        </n-space>
                    </template>

                    <n-space vertical>
                        <!-- 基本信息 -->
                        <n-descriptions :column="2" bordered size="small">
                            <n-descriptions-item label="URL">
                                <n-text>{{ captureResult.url }}</n-text>
                            </n-descriptions-item>
                            <n-descriptions-item label="状态码">
                                <n-tag :type="captureResult.statusCode === 200 ? 'success' : 'warning'">
                                    {{ captureResult.statusCode }}
                                </n-tag>
                            </n-descriptions-item>
                            <n-descriptions-item label="内容类型">
                                <n-text>{{ captureResult.contentType || '未知' }}</n-text>
                            </n-descriptions-item>
                            <n-descriptions-item label="内容大小">
                                <n-text>{{ formatBytes(captureResult.contentLength || 0) }}</n-text>
                            </n-descriptions-item>
                            <n-descriptions-item v-if="captureResult.filesCount" label="文件数量">
                                <n-text>{{ captureResult.filesCount }} 个</n-text>
                            </n-descriptions-item>
                            <n-descriptions-item v-if="captureResult.zipSize" label="ZIP大小">
                                <n-text>{{ formatBytes(captureResult.zipSize) }}</n-text>
                            </n-descriptions-item>
                        </n-descriptions>

                        <!-- ZIP下载状态 -->
                        <n-alert v-if="captureResult.zipPath" type="success" title="ZIP包已自动下载">
                            <template #icon>
                                <n-icon>
                                    <ArchiveOutline />
                                </n-icon>
                            </template>
                            <n-text>完整的网页已打包并下载，包含 {{ captureResult.filesCount }} 个文件</n-text>
                        </n-alert>

                        <!-- 错误信息 -->
                        <n-alert v-if="!captureResult.success && captureResult.error" type="error" title="抓取失败">
                            {{ captureResult.error }}
                        </n-alert>

                        <!-- 文件列表 -->
                        <div
                            v-if="captureResult.success && captureResult.downloadedFiles && captureResult.downloadedFiles.length > 0">
                            <n-card size="small" title="下载的文件">
                                <n-scrollbar style="max-height: 300px;">
                                    <n-list>
                                        <n-list-item v-for="(file, index) in captureResult.downloadedFiles"
                                            :key="index">
                                            <n-thing>
                                                <template #header>
                                                    <n-text>{{ formatFilePath(file) }}</n-text>
                                                </template>
                                                <template #description>
                                                    <n-tag size="small" :type="getFileTypeColor(file)">
                                                        {{ getFileType(file) }}
                                                    </n-tag>
                                                </template>
                                            </n-thing>
                                        </n-list-item>
                                    </n-list>
                                </n-scrollbar>
                            </n-card>
                        </div>
                    </n-space>
                </n-card>
            </n-space>
        </n-card>
    </div>
</template>

<script setup lang="ts">
import { ref, inject } from 'vue'
import { useMessage } from 'naive-ui'
import { CameraOutline, RefreshOutline, ArchiveOutline, DownloadOutline, ChevronDownOutline } from '@vicons/ionicons5'
import api from '@/api'

const message = useMessage()

// 注入全局 loading
const globalLoading = inject('globalLoading') as any

// 表单数据
const form = ref({
    url: ''
})

// 抓取选项
const options = ref({
    includeImages: true,
    includeStyles: true,
    includeScripts: true,
    followRedirects: true,
    timeout: 60,
    maxFiles: 200
})

// 抓取结果
const captureResult = ref<any>(null)

// 调试模式
const debugMode = ref(true) // 默认开启调试模式

// URL选项
const urlOptions = [
    {
        label: 'Example.com (测试)',
        key: 'https://example.com'
    },
    {
        label: 'httpbin.org (测试)',
        key: 'https://httpbin.org/html'
    },
    {
        label: 'GitHub',
        key: 'https://github.com'
    },
    {
        label: 'MDN Web Docs',
        key: 'https://developer.mozilla.org'
    },
    {
        label: 'Bootstrap',
        key: 'https://getbootstrap.com'
    }
]

// 选择URL
const selectUrl = (key: string) => {
    form.value.url = key
}

// 格式化字节大小
const formatBytes = (bytes: number): string => {
    if (bytes === 0) return '0 Bytes'
    const k = 1024
    const sizes = ['Bytes', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 抓取页面
const captureUrl = async () => {
    if (!form.value.url.trim()) {
        message.error('请输入要抓取的 URL')
        return
    }

    // 预处理URL
    let processedUrl = form.value.url.trim()

    // 如果没有协议，自动添加https://
    if (!processedUrl.startsWith('http://') && !processedUrl.startsWith('https://')) {
        processedUrl = 'https://' + processedUrl
        form.value.url = processedUrl
    }

    // 验证 URL 格式
    try {
        const url = new URL(processedUrl)
        // 检查是否为有效的域名
        if (!url.hostname || url.hostname.length < 3) {
            message.error('请输入有效的网站地址')
            return
        }
    } catch (error) {
        message.error('请输入有效的 URL 格式，例如：https://example.com')
        return
    }

    globalLoading.show(`正在抓取页面：${processedUrl}`)

    try {
        if (debugMode.value) {
            console.log('开始抓取页面:', processedUrl)
            console.log('抓取选项:', options.value)
        }

        const result = await api('capture_page', {
            url: processedUrl,
            options: JSON.stringify(options.value)
        })

        if (debugMode.value) {
            console.log('API响应:', result)
        }

        if (result.code === 200) {
            captureResult.value = {
                success: true,
                url: processedUrl,
                timestamp: new Date().toLocaleString(),
                statusCode: result.data.statusCode || 200,
                contentType: result.data.contentType,
                contentLength: result.data.contentLength,
                duration: result.data.duration,
                filesCount: result.data.filesCount,
                zipPath: result.data.zipPath,
                zipSize: result.data.zipSize,
                downloadedFiles: result.data.downloadedFiles
            }

            // 自动下载ZIP文件
            if (result.data.zipPath) {
                await autoDownloadZip(result.data.zipPath)
                message.success('页面抓取完成，ZIP文件已自动下载')
            } else {
                message.success('页面抓取成功')
            }
        } else {
            captureResult.value = {
                success: false,
                url: form.value.url,
                timestamp: new Date().toLocaleString(),
                error: result.msg || '抓取失败',
                statusCode: result.data?.statusCode || 0
            }
            message.error(result.msg || '页面抓取失败')
        }
    } catch (error) {
        console.error('Page capture error:', error)

        let errorMessage = '未知错误'
        if (error instanceof Error) {
            errorMessage = error.message
        } else if (typeof error === 'string') {
            errorMessage = error
        }

        captureResult.value = {
            success: false,
            url: form.value.url,
            timestamp: new Date().toLocaleString(),
            error: errorMessage
        }

        // 提供更友好的错误提示
        if (errorMessage.includes('网络')) {
            message.error('网络连接失败，请检查网络连接或URL是否正确')
        } else if (errorMessage.includes('超时')) {
            message.error('请求超时，请尝试增加超时时间或稍后重试')
        } else if (errorMessage.includes('格式')) {
            message.error('URL格式不正确，请检查输入的网址')
        } else {
            message.error('页面抓取失败：' + errorMessage)
        }
    } finally {
        globalLoading.hide()
    }
}

// 清空结果
const clearResults = () => {
    captureResult.value = null
    message.info('已清空抓取结果')
}

// 自动下载ZIP文件
const autoDownloadZip = async (zipPath: string) => {
    try {
        if (debugMode.value) {
            console.log('开始下载ZIP文件:', zipPath)
        }

        // 调用Go后端的文件下载方法
        const response = await api('download_file', {
            filePath: zipPath
        })

        if (debugMode.value) {
            console.log('下载API响应:', response)
            console.log('响应数据类型:', typeof response?.data)
            console.log('响应数据长度:', response?.data?.length)
        }

        if (response && response.code === 200 && response.data) {
            // 处理Base64编码的二进制数据
            let binaryData
            if (typeof response.data === 'string') {
                // 后端返回Base64编码的字符串，直接解码
                try {
                    const binaryString = atob(response.data)
                    binaryData = new Uint8Array(binaryString.length)
                    for (let i = 0; i < binaryString.length; i++) {
                        binaryData[i] = binaryString.charCodeAt(i)
                    }
                    if (debugMode.value) {
                        console.log('Base64解码成功，数据长度:', binaryData.length)
                    }
                } catch (e) {
                    console.error('Base64解码失败:', e)
                    throw new Error('Base64解码失败: ' + (e as Error).message)
                }
            } else if (Array.isArray(response.data)) {
                // 兼容旧的数组格式
                binaryData = new Uint8Array(response.data)
            } else {
                throw new Error('不支持的数据格式: ' + typeof response.data)
            }

            // 创建下载链接
            const blob = new Blob([binaryData], { type: 'application/zip' })
            const url = window.URL.createObjectURL(blob)
            const link = document.createElement('a')
            link.href = url

            // 生成文件名：网站域名_时间戳.zip
            const urlObj = new URL(captureResult.value.url)
            const domain = urlObj.hostname.replace(/[^a-zA-Z0-9]/g, '_')
            const timestamp = new Date().toISOString().slice(0, 19).replace(/[:-]/g, '')
            link.download = `${domain}_${timestamp}.zip`

            document.body.appendChild(link)
            link.click()
            document.body.removeChild(link)
            window.URL.revokeObjectURL(url)

            if (debugMode.value) {
                console.log('ZIP文件已自动下载:', link.download)
                console.log('文件大小:', blob.size, 'bytes')
            }

            message.success(`ZIP文件已下载: ${link.download}`)
        } else {
            message.error('下载失败：' + (response?.msg || '服务器返回数据为空'))
        }
    } catch (error) {
        console.error('Auto download error:', error)
        message.error('自动下载失败：' + (error as Error).message)
    }
}

// 获取文件类型
const getFileType = (fileName: string): string => {
    const ext = fileName.split('.').pop()?.toLowerCase()
    switch (ext) {
        case 'html':
        case 'htm':
            return 'HTML'
        case 'css':
            return 'CSS'
        case 'js':
            return 'JavaScript'
        case 'jpg':
        case 'jpeg':
        case 'png':
        case 'gif':
        case 'webp':
        case 'svg':
            return '图片'
        default:
            return '其他'
    }
}

// 获取文件类型颜色
const getFileTypeColor = (fileName: string): string => {
    const ext = fileName.split('.').pop()?.toLowerCase()
    switch (ext) {
        case 'html':
        case 'htm':
            return 'success'
        case 'css':
            return 'info'
        case 'js':
            return 'warning'
        case 'jpg':
        case 'jpeg':
        case 'png':
        case 'gif':
        case 'webp':
        case 'svg':
            return 'error'
        default:
            return 'default'
    }
}

// 测试连接
const testConnection = async () => {
    if (!form.value.url.trim()) {
        message.error('请先输入URL')
        return
    }

    let testUrl = form.value.url.trim()
    if (!testUrl.startsWith('http://') && !testUrl.startsWith('https://')) {
        testUrl = 'https://' + testUrl
    }

    try {
        const url = new URL(testUrl)
        message.info(`正在测试连接到: ${url.hostname}`)

        // 这里可以添加一个简单的连接测试
        // 暂时只显示URL解析结果
        message.success(`URL解析成功: ${url.protocol}//${url.hostname}${url.pathname}`)
    } catch (error) {
        message.error('URL格式错误: ' + (error as Error).message)
    }
}

// 格式化文件路径显示
const formatFilePath = (filePath: string): string => {
    // 如果是index.html，显示为根文件
    if (filePath === 'index.html') {
        return '📄 ' + filePath
    }
    // 如果在static目录下，添加文件夹图标
    if (filePath.startsWith('static/')) {
        const parts = filePath.split('/')
        if (parts.length >= 3) {
            const folder = parts[1]
            const file = parts[2]
            const folderIcon = folder === 'css' ? '🎨' : folder === 'js' ? '⚡' : folder === 'images' ? '🖼️' : '📁'
            return `${folderIcon} static/${folder}/${file}`
        }
    }
    return '📄 ' + filePath
}
</script>

<style scoped>
:deep(.n-card .n-card__header) {
    padding-bottom: 12px;
}

:deep(.n-descriptions .n-descriptions-item) {
    padding: 6px 0;
}

:deep(.n-form-item) {
    margin-bottom: 16px;
}
</style>