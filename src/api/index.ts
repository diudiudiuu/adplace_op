import { invoke } from '@tauri-apps/api'
import { InvokeArgs } from '@tauri-apps/api/tauri'

// promise api 调用 tauri api
const api = async (uri: string, data: InvokeArgs | undefined) => {
    try {
        const response = await invoke(uri, data)
        return response
    } catch (error) {
        console.error(error)
    }
}


export default api