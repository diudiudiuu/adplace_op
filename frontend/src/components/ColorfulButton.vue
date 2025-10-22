<template>
    <n-button 
        :class="['colorful-btn', `colorful-btn--${type}`]"
        :size="size"
        quaternary
        @click="$emit('click')"
    >
        <template #icon>
            <n-icon :size="iconSize">
                <component :is="iconComponent" />
            </n-icon>
        </template>
        <span v-if="!iconOnly" class="btn-text">{{ text }}</span>
    </n-button>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { CreateOutline, TrashOutline, AddCircleOutline, RefreshOutline } from '@vicons/ionicons5'

interface Props {
    type: 'edit' | 'delete' | 'add' | 'refresh'
    size?: 'tiny' | 'small' | 'medium' | 'large'
    text?: string
    iconOnly?: boolean
}

const props = withDefaults(defineProps<Props>(), {
    size: 'small',
    text: '',
    iconOnly: false
})

const emit = defineEmits<{
    click: []
}>()

const iconSize = computed(() => {
    switch (props.size) {
        case 'tiny': return 12
        case 'small': return 14
        case 'medium': return 16
        case 'large': return 18
        default: return 14
    }
})

const iconComponent = computed(() => {
    switch (props.type) {
        case 'edit': return CreateOutline
        case 'delete': return TrashOutline
        case 'add': return AddCircleOutline
        case 'refresh': return RefreshOutline
        default: return CreateOutline
    }
})
</script>

<style scoped>
.colorful-btn {
    border-radius: 6px !important;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1) !important;
    border: 1px solid transparent !important;
    position: relative !important;
    overflow: hidden !important;
}

.colorful-btn::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    border-radius: 6px;
    padding: 1px;
    background: linear-gradient(135deg, var(--btn-gradient-from), var(--btn-gradient-to));
    -webkit-mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
    -webkit-mask-composite: subtract;
    mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
    mask-composite: subtract;
    opacity: 0;
    transition: opacity 0.3s ease;
}

.colorful-btn:hover::before {
    opacity: 1;
}

.colorful-btn:hover {
    transform: translateY(-1px);
    box-shadow: 0 4px 12px var(--btn-shadow);
}

.colorful-btn:active {
    transform: translateY(0);
}

/* 编辑按钮 - 蓝色渐变 */
.colorful-btn--edit {
    --btn-gradient-from: #3B82F6;
    --btn-gradient-to: #1D4ED8;
    --btn-shadow: rgba(59, 130, 246, 0.3);
    background: linear-gradient(135deg, rgba(59, 130, 246, 0.1), rgba(29, 78, 216, 0.1)) !important;
    color: #3B82F6 !important;
}

.colorful-btn--edit:hover {
    background: linear-gradient(135deg, #3B82F6, #1D4ED8) !important;
    color: white !important;
}

/* 删除按钮 - 红色渐变 */
.colorful-btn--delete {
    --btn-gradient-from: #EF4444;
    --btn-gradient-to: #DC2626;
    --btn-shadow: rgba(239, 68, 68, 0.3);
    background: linear-gradient(135deg, rgba(239, 68, 68, 0.1), rgba(220, 38, 38, 0.1)) !important;
    color: #EF4444 !important;
}

.colorful-btn--delete:hover {
    background: linear-gradient(135deg, #EF4444, #DC2626) !important;
    color: white !important;
}

/* 添加按钮 - 绿色渐变 */
.colorful-btn--add {
    --btn-gradient-from: #10B981;
    --btn-gradient-to: #059669;
    --btn-shadow: rgba(16, 185, 129, 0.3);
    background: linear-gradient(135deg, rgba(16, 185, 129, 0.1), rgba(5, 150, 105, 0.1)) !important;
    color: #10B981 !important;
}

.colorful-btn--add:hover {
    background: linear-gradient(135deg, #10B981, #059669) !important;
    color: white !important;
}

/* 刷新按钮 - 青色渐变 */
.colorful-btn--refresh {
    --btn-gradient-from: #06B6D4;
    --btn-gradient-to: #0891B2;
    --btn-shadow: rgba(6, 182, 212, 0.3);
    background: linear-gradient(135deg, rgba(6, 182, 212, 0.1), rgba(8, 145, 178, 0.1)) !important;
    color: #06B6D4 !important;
}

.colorful-btn--refresh:hover {
    background: linear-gradient(135deg, #06B6D4, #0891B2) !important;
    color: white !important;
}

.btn-text {
    margin-left: 4px;
    font-weight: 500;
}

/* 图标样式 */
:deep(.n-icon) {
    transition: all 0.3s ease;
}

.colorful-btn:hover :deep(.n-icon) {
    transform: scale(1.1);
}

/* 小尺寸按钮特殊处理 */
.colorful-btn[data-size="small"] {
    min-width: 28px !important;
    height: 28px !important;
    padding: 0 8px !important;
}

.colorful-btn[data-size="tiny"] {
    min-width: 24px !important;
    height: 24px !important;
    padding: 0 6px !important;
}
</style>