import { invoke } from '@tauri-apps/api'
import { InvokeArgs } from '@tauri-apps/api/tauri'

// promise api 调用 tauri api
const api = (uri: string, data: InvokeArgs | undefined) => {
    try {
        return invoke(uri, data).then((res: unknown) => { // Update the type of 'res' to 'unknown'
            return JSON.parse(res as string); // Explicitly cast 'res' as string
        }).catch((err) => {
            console.error(err);
            return {};
        });
        
    } catch (error) {
        console.error(error);
        return {};
    }
}


export default api