<template>
    <n-config-provider :theme="null" :theme-overrides="lightThemeOverrides">
        <n-loading-bar-provider>
            <n-message-provider>
                <n-dialog-provider>
                    <n-notification-provider>
                        <router-view />
                        
                        <!-- 简单的全屏Loading -->
                        <n-modal 
                            v-model:show="isLoading" 
                            :mask-closable="false" 
                            :closable="false"
                            :show-icon="false"
                            transform-origin="center"
                        >
                            <n-card 
                                style="width: 300px; text-align: center"
                                :bordered="false"
                                size="large"
                            >
                                <n-space vertical align="center">
                                    <n-spin size="large" />
                                    <div>
                                        <h3 style="margin: 0 0 8px 0">{{ loadingText }}</h3>
                                        <p style="margin: 0; color: #999">请稍候...</p>
                                    </div>
                                </n-space>
                            </n-card>
                        </n-modal>
                    </n-notification-provider>
                </n-dialog-provider>
            </n-message-provider>
        </n-loading-bar-provider>
    </n-config-provider>
</template>

<script setup lang="ts">
import { lightThemeOverrides } from '@/utils/theme'
import { initFontConfig } from '@/utils/font-config'
import { isLoading, loadingText } from '@/composables/useLoading'

// 初始化字体配置
initFontConfig()
</script>
<style>
</style>
