import api from '@/api'
import { isAuthorized } from '@/utils/auth'

// 数据管理器 - 统一管理服务器数据的缓存和获取
// 注意：缓存只在前端进行，后端每次都从KV服务获取最新数据
class DataManager {
    private static instance: DataManager
    private serverData: any[] = []
    private isLoading = false
    private hasLoaded = false
    private readonly SERVER_DATA_CACHE_KEY = 'serverData'
    private readonly SERVER_DATA_TIMESTAMP_KEY = 'serverDataTimestamp'
    private readonly CACHE_DURATION = 2 * 60 * 60 * 1000 // 2小时缓存

    static getInstance(): DataManager {
        if (!DataManager.instance) {
            console.log('DataManager: Creating new instance')
            DataManager.instance = new DataManager()
            // 初始化时尝试从缓存加载
            DataManager.instance.loadFromCache()
        } else {
            console.log('DataManager: Returning existing instance')
        }
        return DataManager.instance
    }

    // 从缓存加载数据
    private loadFromCache(): void {
        console.log('DataManager: loadFromCache called')
        try {
            const cachedData = localStorage.getItem(this.SERVER_DATA_CACHE_KEY)
            const cachedTimestamp = localStorage.getItem(this.SERVER_DATA_TIMESTAMP_KEY)
            
            console.log('DataManager: Cache check', {
                hasCachedData: !!cachedData,
                hasCachedTimestamp: !!cachedTimestamp,
                cachedDataLength: cachedData ? cachedData.length : 0
            })

            if (cachedData && cachedTimestamp) {
                const timestamp = parseInt(cachedTimestamp)
                const now = Date.now()
                const age = now - timestamp
                
                console.log('DataManager: Cache age check', {
                    timestamp,
                    now,
                    age,
                    cacheDuration: this.CACHE_DURATION,
                    isExpired: age >= this.CACHE_DURATION
                })

                // 检查缓存是否过期
                if (age < this.CACHE_DURATION) {
                    this.serverData = JSON.parse(cachedData)
                    this.hasLoaded = true
                    console.log('DataManager: Loaded server data from cache', {
                        count: this.serverData.length,
                        age: Math.round(age / 1000 / 60) + ' minutes'
                    })
                } else {
                    console.log('DataManager: Cache expired, will reload')
                    this.clearCache()
                }
            } else {
                console.log('DataManager: No cache data found')
            }
        } catch (error) {
            console.error('DataManager: Failed to load from cache:', error)
            this.clearCache()
        }
    }

    // 保存数据到缓存
    private saveToCache(): void {
        try {
            localStorage.setItem(this.SERVER_DATA_CACHE_KEY, JSON.stringify(this.serverData))
            localStorage.setItem(this.SERVER_DATA_TIMESTAMP_KEY, Date.now().toString())
            console.log('DataManager: Saved server data to cache', {
                count: this.serverData.length
            })
        } catch (error) {
            console.error('DataManager: Failed to save to cache:', error)
        }
    }

    // 清除缓存
    private clearCache(): void {
        localStorage.removeItem(this.SERVER_DATA_CACHE_KEY)
        localStorage.removeItem(this.SERVER_DATA_TIMESTAMP_KEY)
        this.hasLoaded = false
    }

    // 获取服务器数据（优先使用缓存）
    async getServerData(forceRefresh = false): Promise<any[]> {
        console.log('DataManager: getServerData called', {
            forceRefresh,
            hasLoaded: this.hasLoaded,
            isLoading: this.isLoading,
            cacheCount: this.serverData.length,
            isAuthorized: isAuthorized()
        })

        // 如果强制刷新，清除缓存状态
        if (forceRefresh) {
            console.log('DataManager: Force refresh requested, clearing cache')
            this.hasLoaded = false
            this.clearCache()
        }

        // 如果正在加载，等待完成
        if (this.isLoading) {
            console.log('DataManager: Already loading, waiting...')
            return new Promise((resolve) => {
                const checkLoading = () => {
                    if (!this.isLoading) {
                        resolve(this.serverData)
                    } else {
                        setTimeout(checkLoading, 100)
                    }
                }
                checkLoading()
            })
        }

        // 如果已经加载过且不是强制刷新，直接返回缓存数据
        if (this.hasLoaded && !forceRefresh) {
            console.log('DataManager: Returning cached server data', {
                count: this.serverData.length
            })
            return Promise.resolve(this.serverData)
        }

        // 检查授权
        if (!isAuthorized()) {
            console.log('DataManager: No authorization, returning empty data')
            return Promise.resolve([])
        }

        this.isLoading = true
        console.log('DataManager: Loading server data from API...')

        try {
            const serverList = await api('list', {})
            console.log('DataManager: Received server data:', serverList)

            if (Array.isArray(serverList)) {
                this.serverData = serverList
                this.hasLoaded = true
                this.saveToCache()
                console.log(`DataManager: Successfully loaded ${this.serverData.length} servers`)
            } else {
                console.log('DataManager: Invalid server data format')
                this.serverData = []
            }
        } catch (error) {
            console.error('DataManager: Failed to load server data:', error)
            // 如果有缓存数据，继续使用
            if (this.serverData.length > 0) {
                console.log('DataManager: Using cached data due to API error')
            }
        } finally {
            this.isLoading = false
        }

        return this.serverData
    }

    // 获取特定服务器信息
    async getServerById(serverId: string, forceRefresh = false): Promise<any | null> {
        const servers = await this.getServerData(forceRefresh)
        return servers.find(server => server.server_id === serverId) || null
    }

    // 获取特定项目信息
    async getProjectById(projectId: string, forceRefresh = false): Promise<any | null> {
        const servers = await this.getServerData(forceRefresh)
        for (const server of servers) {
            if (server.project_list) {
                const project = server.project_list.find((p: any) => p.project_id === projectId)
                if (project) {
                    return project
                }
            }
        }
        return null
    }

    // 强制刷新数据
    async refreshData(): Promise<any[]> {
        console.log('DataManager: Force refreshing server data...')
        return this.getServerData(true)
    }

    // 清除所有数据（用于退出登录）
    clearAllData(): void {
        this.serverData = []
        this.hasLoaded = false
        this.isLoading = false
        this.clearCache()
        console.log('DataManager: All data cleared')
    }

    // 获取缓存的服务器数据（同步方法，不触发API调用）
    getCachedServerData(): any[] {
        return this.serverData
    }

    // 获取缓存状态信息
    getCacheInfo(): { hasCache: boolean; age?: number; count: number } {
        const timestamp = localStorage.getItem(this.SERVER_DATA_TIMESTAMP_KEY)
        if (timestamp && this.hasLoaded) {
            const age = Math.round((Date.now() - parseInt(timestamp)) / 1000 / 60)
            return {
                hasCache: true,
                age,
                count: this.serverData.length
            }
        }
        return {
            hasCache: false,
            count: this.serverData.length
        }
    }

    // 数据变更后的回调（用于其他操作后刷新缓存）
    async onDataChanged(): Promise<void> {
        console.log('DataManager: Data changed, refreshing cache...')
        await this.refreshData()
    }
}

// 导出单例实例
export default DataManager.getInstance()
export { DataManager }