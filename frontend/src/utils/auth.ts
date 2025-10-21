// 授权管理工具函数

/**
 * 检查用户是否已授权
 */
export const isAuthorized = (): boolean => {
    const authorization = localStorage.getItem('authorization');
    return !!(authorization && authorization.trim() !== '');
};

/**
 * 获取授权token
 */
export const getAuthorization = (): string | null => {
    const authorization = localStorage.getItem('authorization');
    return authorization && authorization.trim() !== '' ? authorization.trim() : null;
};

/**
 * 设置授权token
 */
export const setAuthorization = (token: string): void => {
    if (token && token.trim() !== '') {
        localStorage.setItem('authorization', token.trim());
        console.log('Authorization token set');
    }
};

/**
 * 清除所有存储信息（logout 时使用）
 */
export const clearAllStorage = (): void => {
    // 清除 localStorage
    localStorage.clear();
    
    // 清除 sessionStorage
    sessionStorage.clear();
    
    console.log('All storage cleared (localStorage + sessionStorage)');
};

/**
 * 清除授权信息（保持向后兼容）
 */
export const clearAuthorization = (): void => {
    clearAllStorage();
};

/**
 * 检查授权状态并返回相应的默认值
 */
export const getAuthenticatedData = <T>(defaultValue: T): T => {
    return isAuthorized() ? defaultValue : (Array.isArray(defaultValue) ? [] as T : {} as T);
};