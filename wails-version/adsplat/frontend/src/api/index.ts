import { decryptAes } from '@/utils';
import { isAuthorized, getAuthorization, clearAllStorage } from '@/utils/auth';
import { clearMenus } from '@/components/menu';

// 全局消息和加载实例
let globalMessage: any = null;
let globalLoading: any = null;

// 设置全局实例的函数
export const setGlobalInstances = (message: any, loading: any) => {
    globalMessage = message;
    globalLoading = loading;
};

// 扩展 Window 接口以支持 Wails
declare global {
    interface Window {
        go?: {
            main?: {
                App?: any;
            };
        };
    }
}

// 处理未授权错误（401）
const handleUnauthorized = () => {
    console.log('Handling 401 Unauthorized - logging out user');

    // 显示 401 错误提示
    if (globalMessage) {
        globalMessage.error('授权已过期，请重新登录', {
            duration: 3000
        });
    }

    // 清除所有存储
    clearAllStorage();

    // 清除菜单缓存
    clearMenus();

    // 重定向到登录页面
    if (typeof window !== 'undefined') {
        // 延迟一点时间确保清理完成
        setTimeout(() => {
            if (window.location.pathname !== '/lock') {
                window.location.href = '/lock';
            }
        }, 1500); // 延长时间让用户看到提示
    }
};

// 简化的 API 调用映射
const apiMap: Record<string, (data: any) => Promise<string>> = {
    'list': (data: any) => window.go!.main!.App!.List(data.authorization),
    'server_info': (data: any) => window.go!.main!.App!.ServerInfo(data.serverId, data.authorization),
    'project_info': (data: any) => window.go!.main!.App!.ProjectInfo(data.projectId, data.authorization),
    'project_form': (data: any) => window.go!.main!.App!.ProjectForm(data.serverId, data.projectInfo, data.authorization),
    'project_delete': (data: any) => window.go!.main!.App!.ProjectDelete(data.serverId, data.projectId, data.authorization),
    'exec': (data: any) => window.go!.main!.App!.Exec(data.projectId, data.sql, data.sqlType, data.authorization),
    'test_401': () => window.go!.main!.App!.TestUnauthorized(),
};

// 简化的 API 调用函数
const api = async (uri: string, data: any, showLoading: boolean = true) => {
    let loadingInstance: any = null;

    // 开始全屏加载
    if (showLoading && globalLoading) {
        loadingInstance = globalLoading.create({
            show: true,
            description: '请稍候...'
        });
    }

    try {
        // 检查 Wails 环境
        if (!window.go?.main?.App) {
            console.warn(`API ${uri}: Wails not available, returning empty data`);
            if (globalMessage) {
                globalMessage.warning('应用环境未准备就绪，请稍后重试');
            }
            return uri === 'list' ? [] : {};
        }

        // 检查授权
        if (!isAuthorized()) {
            return uri === 'list' ? [] : {};
        }

        // 准备请求数据
        const requestData = { ...data, authorization: getAuthorization() };

        const apiFunction = apiMap[uri];
        if (!apiFunction) return uri === 'list' ? [] : {};

        const res = await apiFunction(requestData);
        const parsedData = JSON.parse(res);

        // 检查是否返回 401 未授权错误
        if (parsedData?.code === 401) {
            console.warn(`API ${uri}: Received 401 Unauthorized, logging out...`);
            handleUnauthorized();
            return uri === 'list' ? [] : {};
        }

        // 检查其他错误状态码
        if (parsedData?.code && parsedData.code !== 200) {
            const errorMsg = parsedData.msg || `请求失败 (${parsedData.code})`;
            if (globalMessage) {
                globalMessage.error(errorMsg);
            }
            console.error(`API ${uri} error:`, parsedData);
        }

        // 处理 list API
        if (uri === 'list') {
            return Array.isArray(parsedData) ? parsedData : [];
        }

        // 解密数据
        if (parsedData?.data) {
            parsedData.data = decryptAes(parsedData.data);
        }

        return parsedData || {};
    } catch (error) {
        console.error(`API ${uri} error:`, error);
        
        if (globalMessage) {
            globalMessage.error(`请求失败: ${error}`);
        }

        return uri === 'list' ? [] : {};
    } finally {
        // 关闭加载
        if (loadingInstance) {
            loadingInstance.destroy();
        }
    }
}

// 测试 401 处理的函数
const test401 = async () => {
    console.log('Testing 401 handling...');
    const result = await api('test_401', {});
    console.log('Test result:', result);
};

// 测试 loading 效果的函数
const testLoading = async () => {
    console.log('Testing loading effect...');
    // 模拟一个延迟的 API 调用
    if (globalLoading) {
        const loading = globalLoading.create({
            show: true,
            description: '测试全屏加载中...'
        });
        setTimeout(() => {
            loading.destroy();
            if (globalMessage) {
                globalMessage.success('Loading 测试完成');
            }
        }, 2000);
    }
};

// 在开发环境暴露测试函数
if (typeof window !== 'undefined') {
    (window as any).test401 = test401;
    (window as any).testLoading = testLoading;
    (window as any).testMessage = () => {
        if (globalMessage) {
            globalMessage.info('这是一个测试消息');
        }
    };
}

// 导出未授权处理函数，供其他地方使用
export { handleUnauthorized };

export default api