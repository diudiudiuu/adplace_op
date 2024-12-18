import { type InvokeArgs, invoke } from '@tauri-apps/api/core'
import { ElLoading } from 'element-plus'

// promise api 调用 tauri api
const api = (uri: string, data: unknown) => {
    // loading
    const loading = ElLoading.service({
        target: '.content-box',
        lock: true,
        text: 'Loading',
    })
    try {
        return invoke(uri, data as InvokeArgs | undefined).then((res: unknown) => { // Update the type of 'res' to 'unknown'
            loading.close();
            return JSON.parse(res as string); // Explicitly cast 'res' as string
        }).catch((err) => {
            console.error(err);
            loading.close();
            return {};
        });

    } catch (error) {
        console.error(error);
        loading.close();
        return {};
    }
}


export default api