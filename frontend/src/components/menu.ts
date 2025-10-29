import type { Menus } from '@/types/menu';
import dataManager from '@/utils/dataManager';
import { isAuthorized } from '@/utils/auth';
import { addDynamicRoutes } from '@/router';

// 菜单状态管理
class MenuManager {
    private static instance: MenuManager;
    private menus: Menus[] = [];
    private isLoading = false;
    private hasLoaded = false;
    private readonly MENU_CACHE_KEY = 'menuData';
    private readonly MENU_CACHE_TIMESTAMP_KEY = 'menuDataTimestamp';
    private readonly CACHE_DURATION = 5 * 60 * 1000; // 5分钟缓存

    static getInstance(): MenuManager {
        if (!MenuManager.instance) {
            MenuManager.instance = new MenuManager();
            // 初始化时尝试从缓存加载
            MenuManager.instance.loadFromCache();
        }
        return MenuManager.instance;
    }

    // 从缓存加载菜单数据
    private loadFromCache(): void {
        try {
            const cachedData = localStorage.getItem(this.MENU_CACHE_KEY);
            const cachedTimestamp = localStorage.getItem(this.MENU_CACHE_TIMESTAMP_KEY);

            if (cachedData && cachedTimestamp) {
                const timestamp = parseInt(cachedTimestamp);
                const now = Date.now();

                // 检查缓存是否过期
                if (now - timestamp < this.CACHE_DURATION) {
                    this.menus = JSON.parse(cachedData);
                    this.hasLoaded = true;
                    console.log('MenuManager: Loaded menus from cache');

                    // 添加动态路由
                    const serverMenus = this.menus.filter(menu => menu.id !== 'welcome');
                    if (serverMenus.length > 0) {
                        addDynamicRoutes(serverMenus);
                    }
                } else {
                    console.log('MenuManager: Cache expired, will reload');
                    this.clearCache();
                }
            }
        } catch (error) {
            console.error('MenuManager: Failed to load from cache:', error);
            this.clearCache();
        }
    }

    // 保存菜单数据到缓存
    private saveToCache(): void {
        try {
            localStorage.setItem(this.MENU_CACHE_KEY, JSON.stringify(this.menus));
            localStorage.setItem(this.MENU_CACHE_TIMESTAMP_KEY, Date.now().toString());
            console.log('MenuManager: Saved menus to cache');
        } catch (error) {
            console.error('MenuManager: Failed to save to cache:', error);
        }
    }

    // 清除缓存
    private clearCache(): void {
        localStorage.removeItem(this.MENU_CACHE_KEY);
        localStorage.removeItem(this.MENU_CACHE_TIMESTAMP_KEY);
    }

    // 获取基础菜单
    private getBaseMenu(): Menus[] {
        return [
            {
                id: 'welcome',
                title: '主页',
                index: '/welcome',
                icon: 'Grid',
            }
        ];
    }



    // 构建服务器菜单项
    private buildServerMenu(serverData: any): Menus {
        const children: any[] = [];

        // 添加项目列表
        if (Array.isArray(serverData.project_list)) {
            serverData.project_list.forEach((project: any) => {
                if (project?.project_id && project?.project_name) {
                    children.push({
                        id: project.project_id,
                        pid: serverData.server_id,
                        index: `/project/${serverData.server_id}/${project.project_id}`,
                        title: project.project_name,
                    });
                }
            });
        }

        // 添加"添加客户"选项
        children.push({
            id: `form_${serverData.server_id}`,
            pid: serverData.server_id,
            index: `/project_form/${serverData.server_id}`,
            title: '➕ 添加客户',
        });

        return {
            id: serverData.server_id,
            title: serverData.server_name,
            index: `/project/${serverData.server_id}`,
            icon: 'Platform',
            children: children
        };
    }

    // 加载菜单数据
    async loadMenus(useExistingData = false): Promise<Menus[]> {
        // 如果正在加载，等待完成
        if (this.isLoading) {
            return new Promise((resolve) => {
                const checkLoading = () => {
                    if (!this.isLoading) {
                        resolve(this.menus);
                    } else {
                        setTimeout(checkLoading, 100);
                    }
                };
                checkLoading();
            });
        }

        // 如果已经加载过，直接返回
        if (this.hasLoaded) {
            return Promise.resolve(this.menus);
        }

        // 重置为基础菜单
        this.menus = this.getBaseMenu();

        // 检查授权
        if (!isAuthorized()) {
            console.log('MenuManager: No authorization, returning base menu');
            return Promise.resolve(this.menus);
        }

        this.isLoading = true;
        console.log('MenuManager: Loading server data...');

        try {
            // 使用数据管理器获取服务器数据
            const serverList = useExistingData
                ? dataManager.getCachedServerData()
                : await dataManager.getServerData();
            console.log('MenuManager: Received server list:', serverList, 'useExistingData:', useExistingData);

            if (!Array.isArray(serverList) || serverList.length === 0) {
                console.log('MenuManager: No server data available');
                return this.menus;
            }

            // 构建服务器菜单
            const serverMenus = serverList
                .filter(server => server?.server_id && server?.server_name)
                .map(server => this.buildServerMenu(server));

            // 合并菜单 - 将抓页面菜单放在最后
            const baseMenus = this.getBaseMenu();
            const pageCaptureMenu = {
                id: 'page-capture',
                title: '网页备份',
                index: '/page-capture',
                icon: 'Platform',
            };
            this.menus = [...baseMenus, ...serverMenus, pageCaptureMenu];

            // 添加动态路由
            if (serverMenus.length > 0) {
                console.log(`MenuManager: Adding ${serverMenus.length} server menus`);
                addDynamicRoutes(serverMenus);
            }

            this.hasLoaded = true;

            // 保存到缓存
            this.saveToCache();

            console.log(`MenuManager: Successfully loaded ${this.menus.length} menu items`);

        } catch (error) {
            console.error('MenuManager: Failed to load menus:', error);
        } finally {
            this.isLoading = false;
        }

        return this.menus;
    }

    // 重新加载菜单（用于授权后刷新）
    async reloadMenus(skipDataRefresh = false): Promise<Menus[]> {
        this.hasLoaded = false;
        this.menus = [];
        this.clearCache(); // 清除缓存
        // 只有在需要时才刷新数据管理器的数据
        if (!skipDataRefresh) {
            await dataManager.refreshData();
        }
        return this.loadMenus(skipDataRefresh);
    }

    // 清除菜单（用于退出登录）
    clearMenus(): void {
        const baseMenus = this.getBaseMenu();
        const pageCaptureMenu = {
            id: 'page-capture',
            title: '📄 网页备份',
            index: '/page-capture',
            icon: 'Platform',
        };
        this.menus = [...baseMenus, pageCaptureMenu];
        this.hasLoaded = false;
        this.isLoading = false;
        this.clearCache();
        // 同时清除数据管理器的数据
        dataManager.clearAllData();
        console.log('MenuManager: Menus cleared');
    }

    // 获取当前菜单
    getCurrentMenus(): Menus[] {
        return this.menus;
    }
}

// 导出单例实例
const menuManager = MenuManager.getInstance();

// 兼容原有接口
export const getMenus = () => menuManager.loadMenus();
export const reloadMenus = (skipDataRefresh = false) => menuManager.reloadMenus(skipDataRefresh);
export const clearMenus = () => menuManager.clearMenus();
export const getCurrentMenus = () => menuManager.getCurrentMenus();

// 导出菜单数据（兼容性）
export const menus = menuManager.getCurrentMenus();

