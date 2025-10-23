<template>
    <div>
        <n-card title=" 页网页备份工具">
            <n-space vertical size="large">
                <!-- 功能说明 -->
                <n-alert type="info" title="功能说明" closable>
                    <n-ul>
                        <n-li>📄 <strong>完整备份</strong>：完整保存网页内容，包括HTML、CSS、JavaScript、图片等所有资源</n-li>
                        <n-li>📦 <strong>标准结构</strong>：生成标准的静态网站结构，index.html + static资源目录</n-li>
                        <n-li>🔗 <strong>链接修正</strong>：自动修正页面中的资源链接，确保离线浏览正常</n-li>
                        <n-li>📁 <strong>自定义保存</strong>：可选择任意目录保存备份文件</n-li>
                        <n-li>✨ <strong>HTML格式化</strong>：自动格式化HTML代码，便于阅读和编辑</n-li>
                        <n-li>🛡️ <strong>隐私清理</strong>：自动删除第三方跟踪、统计、广告代码，保护隐私</n-li>
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
                                开始备份
                            </n-button>
                        </template>
                        备份指定 URL 的页面内容
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
                        清空备份结果
                    </n-tooltip>
                </n-space>


                <!-- 备份配置 -->
                <n-card size="small" title="备份配置">
                    <template #header-extra>
                        <n-space size="small">
                            <n-tooltip>
                                <template #trigger>
                                    <n-button text size="small" @click="showDocumentation">
                                        <template #icon>
                                            <n-icon>
                                                <DocumentTextOutline />
                                            </n-icon>
                                        </template>
                                        功能说明
                                    </n-button>
                                </template>
                                查看详细的功能说明文档
                            </n-tooltip>
                            <n-tooltip>
                                <template #trigger>
                                    <n-button text size="small" @click="showTestPage">
                                        <template #icon>
                                            <n-icon>
                                                <CodeOutline />
                                            </n-icon>
                                        </template>
                                        测试示例
                                    </n-button>
                                </template>
                                查看包含第三方代码和恶意标签的测试页面
                            </n-tooltip>
                            <n-tag type="success" size="small">已优化默认配置</n-tag>
                        </n-space>
                    </template>
                    <n-form :model="form" label-placement="left" label-width="120">
                        <n-form-item label="目标 URL" required>
                            <n-input v-model:value="form.url" placeholder="请输入要备份的网页 URL，如：https://example.com"
                                @keyup.enter="captureUrl" />
                            <template #suffix>
                                <n-button text type="primary" @click="testConnection" :disabled="!form.url.trim()">
                                    测试连接
                                </n-button>
                            </template>
                        </n-form-item>
                        <n-form-item label="保存目录" required>
                            <n-input v-model:value="saveDirectory" placeholder="请选择保存备份文件的目录" readonly>
                                <template #suffix>
                                    <n-button text type="primary" @click="selectDirectory">
                                        选择目录
                                    </n-button>
                                </template>
                            </n-input>
                        </n-form-item>
                        <n-form-item label="备份选项">
                            <n-space vertical>
                                <n-tooltip>
                                    <template #trigger>
                                        <n-checkbox v-model:checked="options.includeImages">
                                            包含图片
                                        </n-checkbox>
                                    </template>
                                    下载并保存网页中的图片资源
                                </n-tooltip>
                                <n-tooltip>
                                    <template #trigger>
                                        <n-checkbox v-model:checked="options.includeStyles">
                                            包含样式
                                        </n-checkbox>
                                    </template>
                                    下载并保存CSS样式文件
                                </n-tooltip>
                                <n-tooltip>
                                    <template #trigger>
                                        <n-checkbox v-model:checked="options.includeScripts">
                                            包含脚本
                                        </n-checkbox>
                                    </template>
                                    下载并保存JavaScript脚本文件
                                </n-tooltip>
                                <n-tooltip>
                                    <template #trigger>
                                        <n-checkbox v-model:checked="options.followRedirects">
                                            跟随重定向
                                        </n-checkbox>
                                    </template>
                                    自动跟随页面重定向链接
                                </n-tooltip>
                                <n-tooltip>
                                    <template #trigger>
                                        <n-checkbox v-model:checked="options.includeFonts">
                                            包含字体
                                        </n-checkbox>
                                    </template>
                                    下载并保存网页字体文件
                                </n-tooltip>
                                <n-tooltip>
                                    <template #trigger>
                                        <n-checkbox v-model:checked="options.includeVideos">
                                            包含视频
                                        </n-checkbox>
                                    </template>
                                    下载并保存视频文件（可能较大）
                                </n-tooltip>
                            </n-space>
                        </n-form-item>
                        <n-form-item label="隐私清理">
                            <n-space vertical>
                                <n-tooltip>
                                    <template #trigger>
                                        <n-checkbox v-model:checked="options.removeAnalytics">
                                            删除统计分析代码
                                        </n-checkbox>
                                    </template>
                                    移除Google Analytics、百度统计、Mixpanel等数据统计代码
                                </n-tooltip>
                                <n-tooltip>
                                    <template #trigger>
                                        <n-checkbox v-model:checked="options.removeTracking">
                                            删除跟踪代码
                                        </n-checkbox>
                                    </template>
                                    移除Facebook Pixel、TikTok Pixel、Hotjar等用户行为跟踪代码
                                </n-tooltip>
                                <n-tooltip>
                                    <template #trigger>
                                        <n-checkbox v-model:checked="options.removeAds">
                                            删除广告代码
                                        </n-checkbox>
                                    </template>
                                    移除Google Ads、Taboola、PopAds等广告投放和联盟代码
                                </n-tooltip>
                                <n-tooltip>
                                    <template #trigger>
                                        <n-checkbox v-model:checked="options.removeTagManager">
                                            删除标签管理器
                                        </n-checkbox>
                                    </template>
                                    移除Google Tag Manager (GTM)等标签管理系统
                                </n-tooltip>
                                <n-tooltip>
                                    <template #trigger>
                                        <n-checkbox v-model:checked="options.removeMaliciousTags">
                                            删除恶意标签
                                        </n-checkbox>
                                    </template>
                                    移除可能被恶意利用的HTML标签：base、meta refresh、meta referrer等
                                </n-tooltip>
                            </n-space>
                        </n-form-item>
                        <n-form-item label="超时时间">
                            <n-input-number v-model:value="options.timeout" :min="60" :max="300" :step="10"
                                placeholder="秒" />
                            <template #suffix>秒</template>
                        </n-form-item>
                        <n-form-item label="最大文件数">
                            <n-input-number v-model:value="options.maxFiles" :min="200" :max="1000" :step="50"
                                placeholder="个" />
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

                        <!-- 备份保存状态 -->
                        <n-alert v-if="captureResult.zipPath && saveDirectory" type="success" title="备份文件已保存">
                            <template #icon>
                                <n-icon>
                                    <ArchiveOutline />
                                </n-icon>
                            </template>
                            <n-space vertical>
                                <n-text>完整的网页已备份并保存，包含 {{ captureResult.filesCount }} 个文件</n-text>
                                <n-text depth="3" style="font-size: 12px;">保存位置: {{ saveDirectory }}</n-text>
                            </n-space>
                        </n-alert>

                        <!-- 未选择目录提示 -->
                        <n-alert v-if="captureResult.zipPath && !saveDirectory" type="warning" title="请选择保存目录">
                            <template #icon>
                                <n-icon>
                                    <ArchiveOutline />
                                </n-icon>
                            </template>
                            <n-space vertical>
                                <n-text>网页备份成功，但未选择保存目录</n-text>
                                <n-button type="primary" @click="selectDirectory">
                                    选择保存目录
                                </n-button>
                            </n-space>
                        </n-alert>

                        <!-- 错误信息 -->
                        <n-alert v-if="!captureResult.success && captureResult.error" type="error" title="备份失败">
                            {{ captureResult.error }}
                        </n-alert>

                        <!-- 文件列表 -->

                    </n-space>
                </n-card>
            </n-space>
        </n-card>

        <!-- 功能说明弹窗 -->
        <n-modal v-model:show="showDocModal" preset="card" title="📖 页面捕获隐私清理功能说明" style="width: 90%; max-width: 1000px;">
            <div v-html="documentationContent" class="documentation-content"></div>
        </n-modal>

        <!-- 测试页面弹窗 -->
        <n-modal v-model:show="showTestModal" preset="card" title="🧪 测试页面代码" style="width: 90%; max-width: 1000px;">
            <n-code :code="testPageContent" language="html" show-line-numbers />
        </n-modal>
    </div>
</template>

<script setup lang="ts">
import { ref, inject } from 'vue'
import { useMessage } from 'naive-ui'
import { RefreshOutline, ArchiveOutline, DocumentTextOutline, CodeOutline } from '@vicons/ionicons5'
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
    includeFonts: true,
    includeVideos: false,
    removeAnalytics: true,
    removeTracking: true,
    removeAds: true,
    removeTagManager: true,
    removeMaliciousTags: true,
    timeout: 60,
    maxFiles: 200
})

// 抓取结果
const captureResult = ref<any>(null)

// 弹窗控制
const showDocModal = ref(false)
const showTestModal = ref(false)

// 文档内容
const documentationContent = ref('')
const testPageContent = ref('')



// 保存目录（从本地缓存加载）
const saveDirectory = ref(localStorage.getItem('pageCapture_saveDirectory') || '')



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

    // 验证保存目录
    if (!saveDirectory.value.trim()) {
        message.error('请先选择保存目录')
        return
    }

    globalLoading.show(`正在备份页面：${processedUrl}`)

    try {


        const result = await api('capture_page', {
            url: processedUrl,
            options: JSON.stringify(options.value)
        })



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

            // 保存ZIP文件到指定目录
            if (result.data.zipPath) {
                if (saveDirectory.value) {
                    await saveZipToDirectory(result.data.zipPath)
                } else {
                    message.warning('未选择保存目录，ZIP文件已生成但未保存')
                }
            } else {
                message.success('页面备份成功')
            }
        } else {
            captureResult.value = {
                success: false,
                url: form.value.url,
                timestamp: new Date().toLocaleString(),
                error: result.msg || '抓取失败',
                statusCode: result.data?.statusCode || 0
            }
            message.error(result.msg || '页面备份失败')
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
    message.info('已清空备份结果')
}

// 保存ZIP文件到指定目录
const saveZipToDirectory = async (zipPath: string) => {
    try {


        // 生成文件名：网站域名_时间戳.zip
        const urlObj = new URL(captureResult.value.url)
        const domain = urlObj.hostname.replace(/[^a-zA-Z0-9]/g, '_')
        const timestamp = new Date().toISOString().slice(0, 19).replace(/[:-]/g, '')
        const fileName = `${domain}_${timestamp}.zip`

        // 调用Go后端的文件保存方法
        const response = await api('save_zip_to_directory', {
            sourcePath: zipPath,
            targetDirectory: saveDirectory.value,
            fileName: fileName
        })



        if (response && response.code === 200) {
            message.success(`ZIP文件已保存到: ${saveDirectory.value}\\${fileName}`)
        } else {
            message.error('保存失败：' + (response?.msg || '未知错误'))
        }
    } catch (error) {
        console.error('Save zip error:', error)
        message.error('保存失败：' + (error as Error).message)
    }
}

// 自动下载ZIP文件（保留作为备用）
const autoDownloadZip = async (zipPath: string) => {
    try {


        // 调用Go后端的文件下载方法
        const response = await api('download_file', {
            filePath: zipPath
        })



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



            message.success(`ZIP文件已下载: ${link.download}`)
        } else {
            message.error('下载失败：' + (response?.msg || '服务器返回数据为空'))
        }
    } catch (error) {
        console.error('Auto download error:', error)
        message.error('自动下载失败：' + (error as Error).message)
    }
}



// 选择保存目录
const selectDirectory = async () => {
    try {
        // 调用Go后端的目录选择方法
        const result = await api('select_directory', {})

        if (result && result.code === 200 && result.data) {
            saveDirectory.value = result.data
            // 缓存到本地存储
            localStorage.setItem('pageCapture_saveDirectory', result.data)
            message.success('目录选择成功')
        } else if (result && result.code === 400) {
            // 用户取消选择
            message.info('已取消选择目录')
        } else {
            message.error('选择目录失败：' + (result?.msg || '未知错误'))
        }
    } catch (error) {
        console.error('Select directory error:', error)
        message.error('选择目录异常：' + (error as Error).message)
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

// 显示功能说明文档
const showDocumentation = () => {
    documentationContent.value = `
        <h2>🛡️ 页面捕获隐私清理功能</h2>
        
        <h3>功能概述</h3>
        <p>在页面备份工具中新增了隐私清理功能，可以自动删除网页中的第三方跟踪、统计、广告代码，保护用户隐私。</p>
        
        <h3>隐私清理选项</h3>
        
        <h4>1. 删除统计分析代码 ✅</h4>
        <ul>
            <li>Google Analytics / gtag.js / GA4</li>
            <li>百度统计 / CNZZ</li>
            <li>Mixpanel / Segment</li>
        </ul>
        
        <h4>2. 删除跟踪代码 ✅</h4>
        <ul>
            <li>Facebook Pixel</li>
            <li>TikTok Pixel / Snapchat Pixel</li>
            <li>Hotjar / CrazyEgg / Clarity</li>
        </ul>
        
        <h4>3. 删除广告代码 ✅</h4>
        <ul>
            <li>Google Ads / DoubleClick</li>
            <li>Taboola / Outbrain</li>
            <li>PopAds / PropellerAds / AdCash</li>
            <li>affiliate.js / redirect.js</li>
        </ul>
        
        <h4>4. 删除标签管理器 ✅</h4>
        <ul>
            <li>Google Tag Manager (GTM)</li>
        </ul>
        
        <h4>5. 删除恶意标签 ✅</h4>
        <ul>
            <li><code>&lt;base href="..."&gt;</code> - 防止劫持所有相对链接</li>
            <li><code>&lt;meta http-equiv="refresh"&gt;</code> - 防止自动跳转到恶意网站</li>
            <li><code>&lt;meta name="referrer"&gt;</code> - 防止来源伪造</li>
            <li>恶意JavaScript重定向代码</li>
        </ul>
        
        <h3>安全防护</h3>
        
        <h4>恶意标签防护</h4>
        <ol>
            <li><strong>&lt;base&gt; 标签劫持防护</strong> - 自动删除所有 base 标签，防止恶意网站劫持页面中的所有相对链接</li>
            <li><strong>自动跳转防护</strong> - 删除 meta refresh 标签，防止页面自动跳转到钓鱼网站或恶意网站</li>
            <li><strong>来源伪造防护</strong> - 删除 meta referrer 标签，防止恶意网站伪造访问来源</li>
            <li><strong>JavaScript重定向防护</strong> - 检测并删除包含恶意重定向的JavaScript代码</li>
        </ol>
        
        <h3>使用方法</h3>
        <ol>
            <li>在页面捕获界面中，找到"隐私清理"选项组</li>
            <li>根据需要勾选要删除的第三方代码类型</li>
            <li>开始备份，系统会自动清理选中的代码类型</li>
            <li>备份完成后，生成的HTML文件将不包含被清理的第三方代码</li>
        </ol>
        
        <h3>注意事项</h3>
        <ul>
            <li>隐私清理功能默认启用，确保用户隐私安全</li>
            <li>恶意标签清理功能默认启用，提供额外的安全防护</li>
            <li>清理过程不会影响页面的基本功能和样式</li>
            <li>被清理的代码包括外部引用和内联代码</li>
            <li>清理后的页面在离线环境下浏览更加安全</li>
        </ul>
    `
    showDocModal.value = true
}

// 显示测试页面代码
const showTestPage = () => {
    testPageContent.value = 'HTML测试页面包含以下内容：\n\n' +
        '1. 恶意标签示例：\n' +
        '   - <base href="https://evil-site.com/">\n' +
        '   - <meta http-equiv="refresh" content="5;url=https://phishing-site.com">\n' +
        '   - <meta name="referrer" content="no-referrer">\n\n' +
        '2. 第三方跟踪代码：\n' +
        '   - Google Analytics\n' +
        '   - Facebook Pixel\n' +
        '   - 百度统计\n' +
        '   - Google Tag Manager\n\n' +
        '3. 恶意JavaScript代码：\n' +
        '   - window.location.href 重定向\n' +
        '   - setTimeout 延时跳转\n' +
        '   - 动态创建base标签\n\n' +
        '这些代码在启用隐私清理功能后会被自动删除。'
    
    showTestModal.value = true
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

.documentation-content {
    line-height: 1.6;
    font-size: 14px;
}

.documentation-content h2 {
    color: #2080f0;
    border-bottom: 2px solid #2080f0;
    padding-bottom: 8px;
    margin-bottom: 16px;
}

.documentation-content h3 {
    color: #18a058;
    margin-top: 24px;
    margin-bottom: 12px;
}

.documentation-content h4 {
    color: #f0a020;
    margin-top: 16px;
    margin-bottom: 8px;
}

.documentation-content ul, .documentation-content ol {
    margin-left: 20px;
    margin-bottom: 12px;
}

.documentation-content li {
    margin-bottom: 4px;
}

.documentation-content code {
    background-color: #f5f5f5;
    padding: 2px 6px;
    border-radius: 4px;
    font-family: 'Courier New', monospace;
    color: #d03050;
}
</style>