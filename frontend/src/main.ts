import { createApp } from 'vue';
import { createPinia } from 'pinia';
import App from './App.vue';
import router from './router';

// 样式文件
import './assets/css/main.css';
import './assets/css/theme.css';

// Naive UI
import naive from 'naive-ui';
import 'vfonts/Lato.css';
import 'vfonts/FiraCode.css';

const app = createApp(App);
app.use(createPinia());
app.use(router);
app.use(naive);

app.mount('#app');