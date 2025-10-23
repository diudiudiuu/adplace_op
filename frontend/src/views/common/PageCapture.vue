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
                            <n-input 
                                v-model:value="form.url" 
                                placeholder="请输入要抓取的网页 URL，如：https://example.com"
                                @keyup.enter="captureUrl"
                            />
                            <template #suffix>
                                <n-dropdown :options="urlOptions" @select="selectUrl">
                                    <n-button text>
                                        <template #icon>
                                            <n-icon>
                                                <ChevronDownOutline />
                                            </n-icon>
                                        </template>
                                    </n-button>
                                </n-dropdown>
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
                                <n-checkbox v-model:checked="options.createZip">
                                    创建ZIP包（完整下载）
                                </n-checkbox>
                            </n-space>
                        </n-form-item>
                        <n-form-item label="超时时间">
                            <n-input-number 
                                v-model:value="options.timeout" 
                                :min="10" 
                                :max="180" 
                                :step="10"
                                placeholder="秒"
                            />
                            <template #suffix>秒</template>
                        </n-form-item>
                        <n-form-item v-if="options.createZip" label="最大文件数">
                            <n-input-number 
                                v-model:value="options.maxFiles" 
                                :min="50" 
                                :max="1000" 
                                :step="50"
                                placeholder="个"
                            />
                            <template #suffix>个</template>
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

                        <!-- ZIP下载 -->
                        <n-alert v-if="captureResult.zipPath" type="success" title="ZIP包已生成">
                            <template #icon>
                                <n-icon>
                                    <ArchiveOutline />
                                </n-icon>
                            </template>
                            <n-space vertical>
                                <n-text>完整的网页已打包为ZIP文件，包含 {{ captureResult.filesCount }} 个文件</n-text>
                                <n-button type="primary" @click="downloadZip">
                                    <template #icon>
                                        <n-icon>
                                            <DownloadOutline />
                                        </n-icon>
                                    </template>
                                    下载ZIP包
                                </n-button>
                            </n-space>
                        </n-alert>

                        <!-- 错误信息 -->
                        <n-alert v-if="!captureResult.success && captureResult.error" type="error" title="抓取失败">
                            {{ captureResult.error }}
                        </n-alert>

                        <!-- 页面内容预览 -->
                        <div v-if="captureResult.success && captureResult.content">
                            <n-tabs type="line" animated>
                                <n-tab-pane name="preview" tab="页面预览">
                                    <n-scrollbar style="max-height: 400px;">
                                        <div class="page-preview" v-html="captureResult.content"></div>
                                    </n-scrollbar>
                                </n-tab-pane>
                                <n-tab-pane name="source" tab="源代码">
                                    <n-scrollbar style="max-height: 400px;">
                                        <pre class="source-code">{{ captureResult.content }}</pre>
                                    </n-scrollbar>
                                </n-tab-pane>
                                <n-tab-pane name="files" tab="文件列表" v-if="captureResult.downloadedFiles && captureResult.downloadedFiles.length > 0">
                                    <n-scrollbar style="max-height: 400px;">
                                        <n-list>
                                            <n-list-item v-for="(file, index) in captureResult.downloadedFiles" :key="index">
                                                <n-thing>
                                                    <template #header>
                                                        <n-text>{{ formatFilePath(file) }}</n-text>
                                                    </template>
                                                    <template #description>
                                                        <n-space>
                                                            <n-tag size="small" :type="getFileTypeColor(file)">
                                                                {{ getFileType(file) }}
                                                            </n-tag>
                                                            <n-text depth="3" style="font-size: 12px;">{{ file }}</n-text>
                                                        </n-space>
                                                    </template>
                                                </n-thing>
                                            </n-list-item>
                                        </n-list>
                                    </n-scrollbar>
                                </n-tab-pane>
                                <n-tab-pane name="info" tab="详细信息">
                                    <n-descriptions :column="1" bordered size="small">
                                        <n-descriptions-item label="响应头">
                                            <pre class="headers-code">{{ JSON.stringify(captureResult.headers || {}, null, 2) }}</pre>
                                        </n-descriptions-item>
                                        <n-descriptions-item label="抓取时间">
                                            <n-text>{{ captureResult.duration }}ms</n-text>
                                        </n-descriptions-item>
                                    </n-descriptions>
                                </n-tab-pane>
                            </n-tabs>
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
    createZip: true,
    maxFiles: 200
})

// 抓取结果
const captureResult = ref<any>(null)

// URL选项
const urlOptions = [
    {
        label: 'Example.com',
        key: 'https://example.com'
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

    // 验证 URL 格式
    try {
        new URL(form.value.url)
    } catch (error) {
        message.error('请输入有效的 URL 格式')
        return
    }

    globalLoading.show(`正在抓取页面：${form.value.url}`)

    try {
        const result = await api('capture_page', {
            url: form.value.url,
            options: JSON.stringify(options.value)
        })

        if (result.code === 200) {
            captureResult.value = {
                success: true,
                url: form.value.url,
                timestamp: new Date().toLocaleString(),
                statusCode: result.data.statusCode || 200,
                contentType: result.data.contentType,
                contentLength: result.data.contentLength,
                content: result.data.content,
                headers: result.data.headers,
                duration: result.data.duration,
                filesCount: result.data.filesCount,
                zipPath: result.data.zipPath,
                zipSize: result.data.zipSize,
                downloadedFiles: result.data.downloadedFiles
            }
            message.success('页面抓取成功')
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
        captureResult.value = {
            success: false,
            url: form.value.url,
            timestamp: new Date().toLocaleString(),
            error: (error as Error).message
        }
        message.error('页面抓取异常：' + (error as Error).message)
    } finally {
        globalLoading.hide()
    }
}

// 清空结果
const clearResults = () => {
    captureResult.value = null
    message.info('已清空抓取结果')
}

// 下载ZIP文件
const downloadZip = async () => {
    if (!captureResult.value?.zipPath) {
        message.error('ZIP文件路径不存在')
        return
    }

    globalLoading.show('正在下载ZIP文件...')

    try {
        // 调用Go后端的文件下载方法
        const response = await api('download_file', {
            filePath: captureResult.value.zipPath
        })
        
        if (response && response.code === 200) {
            message.success('ZIP文件下载功能已实现，文件已保存到临时目录')
            message.info(`文件路径: ${captureResult.value.zipPath}`)
        } else {
            message.error('下载失败：' + (response?.msg || '未知错误'))
        }
    } catch (error) {
        console.error('Download error:', error)
        message.error('下载失败：' + (error as Error).message)
    } finally {
        globalLoading.hide()
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
.page-preview {
    border: 1px solid #e0e0e0;
    border-radius: 6px;
    padding: 12px;
    background: #fafafa;
    font-size: 12px;
    line-height: 1.4;
}

.source-code {
    background: #f5f5f5;
    border: 1px solid #e0e0e0;
    border-radius: 6px;
    padding: 12px;
    font-size: 11px;
    line-height: 1.3;
    white-space: pre-wrap;
    word-break: break-all;
    margin: 0;
}

.headers-code {
    background: #f8f9fa;
    border: 1px solid #e9ecef;
    border-radius: 4px;
    padding: 8px;
    font-size: 11px;
    line-height: 1.3;
    margin: 0;
}

:deep(.n-card .n-card__header) {
    padding-bottom: 12px;
}

:deep(.n-descriptions .n-descriptions-item) {
    padding: 6px 0;
}

:deep(.n-form-item) {
    margin-bottom: 16px;
}

:deep(.n-tabs .n-tabs-pane) {
    padding: 12px 0;
}
</style>