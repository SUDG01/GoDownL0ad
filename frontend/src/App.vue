<script setup>
import {reactive, onMounted} from 'vue'
import {DownloadDemo} from '../wailsjs/go/main/App'
// 导入 Wails 的事件监听功能
import {EventsOn} from '../wailsjs/runtime/runtime'

const data = reactive({
  // 给个测试链接，比如这个 10MB 的测试文件
  url: "https://speed.hetzner.de/100MB.bin",
  status: "准备就绪",
  progress: 0, // 0 到 100
  isDownloading: false
})

// ⭐ 组件加载时，开始监听后端发来的 "download_progress" 信号
onMounted(() => {
  EventsOn("download_progress", (percentage) => {
    data.progress = percentage.toFixed(1); // 保留一位小数
    data.status = `正在下载... ${data.progress}%`
  })
})

function startTask() {
  if (data.isDownloading) return;
  
  data.isDownloading = true;
  data.status = "正在连接资源...";
  data.progress = 0;

  // 调用 Go 后端
  DownloadDemo(data.url).then(result => {
    data.status = result;
    data.isDownloading = false;
    if(result.includes("完成")) data.progress = 100;
  })
}
</script>

<template>
  <div class="container">
    <h1>🚀 极速下载器 Demo</h1>
    
    <div class="input-group">
      <input v-model="data.url" placeholder="输入下载链接" />
      <button @click="startTask" :disabled="data.isDownloading">
        {{ data.isDownloading ? '下载中...' : '开始下载' }}
      </button>
    </div>

    <div class="progress-container">
      <div class="progress-bar" :style="{ width: data.progress + '%' }"></div>
    </div>
    
    <p class="status">{{ data.status }}</p>
  </div>
</template>

<style>
/* CSS 美化部分 */
.container {
  max-width: 600px;
  margin: 50px auto;
  font-family: 'Segoe UI', sans-serif;
  text-align: center;
}

.input-group {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
}

input {
  flex: 1;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 5px;
}

button {
  padding: 10px 20px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  transition: background 0.3s;
}

button:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}

/* 进度条轨道 */
.progress-container {
  height: 20px;
  background-color: #e0e0e0;
  border-radius: 10px;
  overflow: hidden; /* 保证圆角 */
  margin-bottom: 10px;
}

/* 进度条本体 - 极光色渐变 */
.progress-bar {
  height: 100%;
  background: linear-gradient(90deg, #00c6ff, #0072ff);
  width: 0%; /* 初始宽度 */
  transition: width 0.1s linear; /* 让动画丝滑 */
}

.status {
  color: #666;
  font-size: 14px;
}
</style>