<template>
    <div class="captcha-container">
        <canvas 
            ref="captchaCanvas" 
            :width="captchaWidth" 
            :height="captchaHeight"
            @click="generateCaptcha"
            class="captcha-canvas"
            title="点击刷新验证码"
        ></canvas>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'

interface Props {
    width?: number
    height?: number
}

const props = withDefaults(defineProps<Props>(), {
    width: 120,
    height: 40
})

const emit = defineEmits<{
    change: [code: string]
}>()

const captchaCanvas = ref<HTMLCanvasElement>()
const captchaWidth = ref(props.width)
const captchaHeight = ref(props.height)
const currentCode = ref('')

// 生成随机字符串
const generateRandomCode = (): string => {
    const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789'
    let result = ''
    for (let i = 0; i < 4; i++) {
        result += chars.charAt(Math.floor(Math.random() * chars.length))
    }
    return result
}

// 生成随机颜色
const getRandomColor = (): string => {
    const colors = ['#FF6B6B', '#4ECDC4', '#45B7D1', '#96CEB4', '#FFEAA7', '#DDA0DD', '#98D8C8']
    return colors[Math.floor(Math.random() * colors.length)]
}

// 绘制验证码
const generateCaptcha = () => {
    if (!captchaCanvas.value) return
    
    const canvas = captchaCanvas.value
    const ctx = canvas.getContext('2d')
    if (!ctx) return

    // 清空画布
    ctx.clearRect(0, 0, captchaWidth.value, captchaHeight.value)
    
    // 设置背景
    ctx.fillStyle = '#f8f9fa'
    ctx.fillRect(0, 0, captchaWidth.value, captchaHeight.value)
    
    // 生成验证码
    currentCode.value = generateRandomCode()
    
    // 绘制字符
    const fontSize = 24
    ctx.font = `bold ${fontSize}px Arial`
    ctx.textBaseline = 'middle'
    
    for (let i = 0; i < currentCode.value.length; i++) {
        const char = currentCode.value[i]
        const x = 15 + i * 22
        const y = captchaHeight.value / 2
        
        // 随机旋转角度
        const angle = (Math.random() - 0.5) * 0.4
        
        ctx.save()
        ctx.translate(x, y)
        ctx.rotate(angle)
        ctx.fillStyle = getRandomColor()
        ctx.fillText(char, 0, 0)
        ctx.restore()
    }
    
    // 添加干扰线
    for (let i = 0; i < 3; i++) {
        ctx.strokeStyle = getRandomColor()
        ctx.lineWidth = 1
        ctx.beginPath()
        ctx.moveTo(Math.random() * captchaWidth.value, Math.random() * captchaHeight.value)
        ctx.lineTo(Math.random() * captchaWidth.value, Math.random() * captchaHeight.value)
        ctx.stroke()
    }
    
    // 添加干扰点
    for (let i = 0; i < 30; i++) {
        ctx.fillStyle = getRandomColor()
        ctx.beginPath()
        ctx.arc(
            Math.random() * captchaWidth.value,
            Math.random() * captchaHeight.value,
            1,
            0,
            2 * Math.PI
        )
        ctx.fill()
    }
    
    // 发送验证码给父组件
    emit('change', currentCode.value)
}

onMounted(() => {
    nextTick(() => {
        generateCaptcha()
    })
})

defineExpose({
    refresh: generateCaptcha
})
</script>

<style scoped>
.captcha-container {
    display: inline-block;
}

.captcha-canvas {
    border: 1px solid #e0e0e6;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.3s ease;
}

.captcha-canvas:hover {
    border-color: #36ad6a;
    box-shadow: 0 0 0 2px rgba(54, 173, 106, 0.2);
}
</style>