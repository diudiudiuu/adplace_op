<template>
    <component :is="iconComponent" />
</template>

<script setup lang="ts">
import { computed, h } from 'vue'

interface Props {
    name: string
    size?: number
}

const props = withDefaults(defineProps<Props>(), {
    size: 18
})

// 彩色图标组件映射
const colorfulIcons = {
    Grid: () => h('svg', {
        width: props.size,
        height: props.size,
        viewBox: '0 0 24 24',
        fill: 'none'
    }, [
        h('rect', {
            x: '3',
            y: '3',
            width: '7',
            height: '7',
            rx: '2',
            fill: 'url(#grid-gradient-1)'
        }),
        h('rect', {
            x: '14',
            y: '3',
            width: '7',
            height: '7',
            rx: '2',
            fill: 'url(#grid-gradient-2)'
        }),
        h('rect', {
            x: '14',
            y: '14',
            width: '7',
            height: '7',
            rx: '2',
            fill: 'url(#grid-gradient-3)'
        }),
        h('rect', {
            x: '3',
            y: '14',
            width: '7',
            height: '7',
            rx: '2',
            fill: 'url(#grid-gradient-4)'
        }),
        h('defs', {}, [
            h('linearGradient', { id: 'grid-gradient-1', x1: '0%', y1: '0%', x2: '100%', y2: '100%' }, [
                h('stop', { offset: '0%', 'stop-color': '#FF6B6B' }),
                h('stop', { offset: '100%', 'stop-color': '#FF8E8E' })
            ]),
            h('linearGradient', { id: 'grid-gradient-2', x1: '0%', y1: '0%', x2: '100%', y2: '100%' }, [
                h('stop', { offset: '0%', 'stop-color': '#4ECDC4' }),
                h('stop', { offset: '100%', 'stop-color': '#6EDDD6' })
            ]),
            h('linearGradient', { id: 'grid-gradient-3', x1: '0%', y1: '0%', x2: '100%', y2: '100%' }, [
                h('stop', { offset: '0%', 'stop-color': '#45B7D1' }),
                h('stop', { offset: '100%', 'stop-color': '#67C3DD' })
            ]),
            h('linearGradient', { id: 'grid-gradient-4', x1: '0%', y1: '0%', x2: '100%', y2: '100%' }, [
                h('stop', { offset: '0%', 'stop-color': '#96CEB4' }),
                h('stop', { offset: '100%', 'stop-color': '#B2D8C4' })
            ])
        ])
    ]),

    Platform: () => h('svg', {
        width: props.size,
        height: props.size,
        viewBox: '0 0 24 24',
        fill: 'none'
    }, [
        h('path', {
            d: 'M4 6C4 4.89543 4.89543 4 6 4H18C19.1046 4 20 4.89543 20 6V8H4V6Z',
            fill: 'url(#server-gradient-1)'
        }),
        h('path', {
            d: 'M4 10H20V12C20 13.1046 19.1046 14 18 14H6C4.89543 14 4 13.1046 4 12V10Z',
            fill: 'url(#server-gradient-2)'
        }),
        h('path', {
            d: 'M4 16H20V18C20 19.1046 19.1046 20 18 20H6C4.89543 20 4 19.1046 4 18V16Z',
            fill: 'url(#server-gradient-3)'
        }),
        h('circle', { cx: '7', cy: '6.5', r: '1', fill: '#FFFFFF' }),
        h('circle', { cx: '7', cy: '11', r: '1', fill: '#FFFFFF' }),
        h('circle', { cx: '7', cy: '17.5', r: '1', fill: '#FFFFFF' }),
        h('defs', {}, [
            h('linearGradient', { id: 'server-gradient-1', x1: '0%', y1: '0%', x2: '100%', y2: '0%' }, [
                h('stop', { offset: '0%', 'stop-color': '#667eea' }),
                h('stop', { offset: '100%', 'stop-color': '#764ba2' })
            ]),
            h('linearGradient', { id: 'server-gradient-2', x1: '0%', y1: '0%', x2: '100%', y2: '0%' }, [
                h('stop', { offset: '0%', 'stop-color': '#f093fb' }),
                h('stop', { offset: '100%', 'stop-color': '#f5576c' })
            ]),
            h('linearGradient', { id: 'server-gradient-3', x1: '0%', y1: '0%', x2: '100%', y2: '0%' }, [
                h('stop', { offset: '0%', 'stop-color': '#4facfe' }),
                h('stop', { offset: '100%', 'stop-color': '#00f2fe' })
            ])
        ])
    ]),

    Add: () => h('svg', {
        width: props.size,
        height: props.size,
        viewBox: '0 0 24 24',
        fill: 'none'
    }, [
        h('circle', {
            cx: '12',
            cy: '12',
            r: '10',
            fill: 'url(#add-gradient)'
        }),
        h('path', {
            d: 'M12 8V16M8 12H16',
            stroke: '#FFFFFF',
            'stroke-width': '2.5',
            'stroke-linecap': 'round'
        }),
        h('defs', {}, [
            h('linearGradient', { id: 'add-gradient', x1: '0%', y1: '0%', x2: '100%', y2: '100%' }, [
                h('stop', { offset: '0%', 'stop-color': '#56ab2f' }),
                h('stop', { offset: '100%', 'stop-color': '#a8e6cf' })
            ])
        ])
    ])
}

const iconComponent = computed(() => {
    return colorfulIcons[props.name as keyof typeof colorfulIcons] || colorfulIcons.Grid
})
</script>