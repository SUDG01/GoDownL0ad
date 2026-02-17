<script setup>
import {reactive, onMounted} from 'vue'
// 注意：这里请确保路径和你实际项目一致，如果报错请改回你原来的路径
import {DownloadDemo} from '../wailsjs/go/main/App'
import {EventsOn} from '../wailsjs/runtime/runtime'
// 引入背景图片 (请确保你放了一张图片在 src/assets/bg.jpg)
// 如果没有图片，请注释掉这一行，并将 CSS 中的 background-image 改为纯色
import bgImage from './assets/bg.jpg'

const data = reactive({
  url: "",
  status: "准备就绪 (´• ω •`)ﾉ", // 改了个可爱点的初始状态
  progress: 0,
  speed: "0 B/s",
  isDownloading: false
})

function formatSpeed(bytes) {
  if (bytes === 0) return '0 B/s';
  const k = 1024;
  const sizes = ['B/s', 'KB/s', 'MB/s', 'GB/s'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  // 修复了原代码中的一个小括号位置错误导致计算不准的问题
  return (bytes / Math.pow(k, i)).toFixed(1) + ' ' + sizes[i];
}

onMounted(() => {
  EventsOn("download_progress", (percentage) => {
    data.progress = percentage.toFixed(1);
    // 加点可爱的颜文字
    data.status = `少女祈祷中... ${data.progress}% (o゜▽゜)o☆`
  })

  EventsOn("download_speed", (bytesPerSecond) => {
    data.speed = formatSpeed(bytesPerSecond);
  })
})

function startTask() {
  if (data.isDownloading) return;

  if (!data.url) {
    data.status = "呐，链接还没填呢！(＃°Д°)"
    return;
  }
  
  data.isDownloading = true;
  data.status = "正在建立契约...";
  data.progress = 0;

  DownloadDemo(data.url).then(result => {
    data.status = result + " (๑•̀ㅂ•́)و✧";
    data.isDownloading = false;
    if(result.includes("完成")) data.progress = 100;
  }).catch(err => {
     data.status = "发生错误了呜呜呜: " + err;
     data.isDownloading = false;
  })
}
</script>

<template>
  <div class="app-background" :style="{ backgroundImage: `url(${bgImage})` }">
    <div class="overlay"></div>

    <div class="acg-container">
      <h1 class="acg-title">
        <span class="title-icon"></span>
        GoDownL0ad <span class="highlight">Demo</span>
        <span class="title-decoration"></span>
      </h1>
      
      <div class="input-group acg-input-group">
        <input v-model="data.url" placeholder="请把链接投喂到这里... (*/ω＼*)" class="acg-input" />
        <button @click="startTask" :disabled="data.isDownloading" class="acg-button check-button">
          <span v-if="data.isDownloading">下载中...</span>
          <span v-else>开始契约</span>
        </button>
      </div>

      <div class="progress-wrapper">
        <div class="progress-container acg-progress-track">
          <div class="progress-bar acg-progress-fill progress-striped" :style="{ width: data.progress + '%' }">
              <span class="progress-text" v-if="data.progress > 5">{{ data.progress }}%</span>
          </div>
        </div>
      </div>
      
      <div class="info-row acg-info-row">
        <div class="info-bubble status-bubble">
          <span class="icon"></span>
          <span class="text">{{ data.status }}</span>
        </div>
        <div class="info-bubble speed-bubble">
          <span class="icon"></span>
          <span class="text speed-text">{{ data.speed }}</span>
        </div>
      </div>

    </div>
  </div>
</template>

<style scoped>
/* 引入一种看起来比较圆润可爱的字体 (可选，如果系统没有会回退) */
@import url('https://fonts.googleapis.com/css2?family=M+PLUS+Rounded+1c:wght@400;700&display=swap');

/* --- 全局背景与布局 --- */
.app-background {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  display: flex;
  justify-content: center;
  align-items: center;
  font-family: 'M PLUS Rounded 1c', 'Segoe UI', sans-serif; /* 使用圆体字 */
  /* 如果没有背景图，用这个渐变代替 */
  background: linear-gradient(135deg, #fce3ec 0%, #dbeaff 100%);
}

.overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  /* 叠加一层半透明白色，让背景变淡，突出前景 */
  background: rgba(255, 255, 255, 0.4);
  backdrop-filter: blur(4px); /* 稍微模糊一下背景 */
  z-index: 1;
}

/* --- 主体卡片 --- */
.acg-container {
  position: relative;
  z-index: 2;
  width: 90%;
  max-width: 650px;
  /* 毛玻璃效果背景 */
  background: rgba(255, 255, 255, 0.75);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-radius: 30px; /* 大圆角 */
  padding: 40px;
  /* 柔和的粉色光晕阴影 */
  box-shadow: 0 10px 40px rgba(255, 182, 193, 0.4), inset 0 0 0 2px rgba(255, 255, 255, 0.6);
  border: 2px solid rgba(255, 255, 255, 0.8);
  text-align: center;
  transition: transform 0.3s ease;
}
.acg-container:hover {
    transform: translateY(-5px); /* 鼠标悬停微微上浮 */
}

/* --- 标题 --- */
.acg-title {
  font-size: 2.5em;
  color: #555;
  margin-bottom: 35px;
  /* 文字渐变 */
  background: linear-gradient(to right, #ff7eb3, #8e44ad);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  font-weight: 800;
  position: relative;
}
.highlight {
    color: #ff7eb3; /* 虽然上面渐变覆盖了，但留着做后备 */
}
.title-decoration {
    position: absolute;
    top: -10px;
    right: 20px;
    font-size: 0.6em;
    animation: twinkle 2s infinite ease-in-out;
}
@keyframes twinkle {
    0%, 100% { opacity: 1; transform: scale(1); }
    50% { opacity: 0.5; transform: scale(0.8); }
}

/* --- 输入框组合 --- */
.acg-input-group {
  display: flex;
  gap: 15px;
  margin-bottom: 30px;
  background: rgba(255, 255, 255, 0.6);
  padding: 5px;
  border-radius: 50px; /* 胶囊形状 */
  box-shadow: inset 0 2px 5px rgba(0,0,0,0.05);
}

.acg-input {
  flex: 1;
  padding: 15px 25px;
  border: 2px solid transparent;
  background: transparent;
  border-radius: 50px;
  font-size: 16px;
  color: #666;
  outline: none;
  transition: all 0.3s;
}

.acg-input:focus {
  border-color: #ffb6c1; /* 聚焦时显示粉色边框 */
  background: rgba(255, 255, 255, 0.9);
  box-shadow: 0 0 15px rgba(255, 182, 193, 0.5);
}
/* 修改 placeholder 颜色 */
.acg-input::placeholder {
    color: #aaa;
    font-size: 14px;
}

/* --- 按钮 --- */
.acg-button {
  padding: 15px 35px;
  /* 粉色到紫色的渐变果冻按钮 */
  background: linear-gradient(135deg, #ff9a9e 0%, #fad0c4 99%, #fad0c4 100%);
  color: white;
  border: none;
  border-radius: 50px;
  font-size: 16px;
  font-weight: bold;
  cursor: pointer;
  box-shadow: 0 5px 15px rgba(255, 154, 158, 0.4);
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
  white-space: nowrap;
}

.acg-button:hover:not(:disabled) {
  transform: scale(1.05) translateY(-2px);
  box-shadow: 0 8px 20px rgba(255, 154, 158, 0.6);
  filter: brightness(1.05);
}

.acg-button:active:not(:disabled) {
    transform: scale(0.95);
}

.acg-button:disabled {
  /* 禁用时的灰色渐变 */
  background: linear-gradient(135deg, #d3cce3, #e9e4f0);
  color: #999;
  box-shadow: none;
  cursor: not-allowed;
}

/* --- 进度条 --- */
.progress-wrapper {
    padding: 5px;
    background: rgba(255,255,255,0.5);
    border-radius: 25px;
    margin-bottom: 25px;
}
.acg-progress-track {
  height: 28px; /* 加高一点 */
  background-color: #ffe4e1; /* 浅粉色轨道 */
  border-radius: 20px;
  overflow: hidden;
  box-shadow: inset 0 2px 5px rgba(0,0,0,0.05);
  border: 2px solid #fff;
}

.acg-progress-fill {
  height: 100%;
  /* 更梦幻的渐变：粉色 -> 蓝色 */
  background: linear-gradient(90deg, #ff9a44, #fc6767, #ec008c, #6a82fb);
  background-size: 300% 100%;
  width: 0%;
  transition: width 0.2s ease-out;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  box-shadow: 0 0 10px rgba(236, 0, 140, 0.5);
  position: relative;
}

/* 进度条上的条纹动画 */
.progress-striped {
    background-image: linear-gradient(
        45deg,
        rgba(255, 255, 255, 0.15) 25%,
        transparent 25%,
        transparent 50%,
        rgba(255, 255, 255, 0.15) 50%,
        rgba(255, 255, 255, 0.15) 75%,
        transparent 75%,
        transparent
    );
    background-size: 40px 40px;
    animation: progress-stripes 1s linear infinite;
}
@keyframes progress-stripes {
    from { background-position: 40px 0; }
    to { background-position: 0 0; }
}


.progress-text {
    margin-right: 10px;
    color: white;
    font-weight: bold;
    font-size: 14px;
    text-shadow: 1px 1px 2px rgba(0,0,0,0.2);
}

/* --- 信息气泡行 --- */
.acg-info-row {
  display: flex;
  justify-content: space-between;
  margin-top: 20px;
  gap: 15px;
}

.info-bubble {
    flex: 1;
    background: rgba(255,255,255,0.8);
    padding: 12px 20px;
    border-radius: 20px;
    display: flex;
    align-items: center;
    gap: 10px;
    box-shadow: 0 4px 10px rgba(0,0,0,0.05);
    border: 2px solid #fff;
    transition: transform 0.3s;
}
.info-bubble:hover {
    transform: translateY(-3px);
}

.status-bubble {
    border-color: #e0c3fc; /* 浅紫色边框 */
}
.speed-bubble {
    border-color: #c2e9fb; /* 浅蓝色边框 */
}

.info-bubble .icon {
    font-size: 1.2em;
}
.info-bubble .text {
    font-size: 14px;
    color: #555;
    text-align: left;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.speed-text {
  color: #4facfe; /* 鲜艳的蓝色 */
  font-weight: 800;
  font-family: 'M PLUS Rounded 1c', monospace;
}
</style>