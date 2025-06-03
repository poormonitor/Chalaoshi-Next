<!-- App.vue -->
<script setup>
import * as zip from "@zip.js/zip.js";
import {GetFileBase64, GetStatus} from "../wailsjs/go/main/App.js";
import Papa from 'papaparse';
import {provide, reactive, onMounted, ref, computed} from 'vue'
import SearchBox from "./components/SearchBox.vue";
import HomeButton from "./components/HomeButton.vue";

const password = "5f4bf44780d48556ee2cc0af48b5fba0fd68ebd3006fd4d268efea252397efd4"

const status = ref(0)
const appData = reactive({
  teachers: [],
  gpaData: {},
  comments: new Map()
})

provide('appData', appData)

const statusStr = computed(() => {
  switch (status.value) {
    case 0:
      return "加载成功";
    case 1:
      return "初始化下载";
    case 2:
      return "正在从服务器下载";
    case 3:
      return "正在寻找 Peers";
    case 4:
      return "正在连接 Peers";
    case 5:
      return "正在从 Peers 下载";
    default:
      return "";
  }
})

async function loadDataWrapper() {
  if (!appData.teachers || appData.teachers.length === 0) {
    status.value = await GetStatus()
    if (status.value === 0) await loadData();
    setTimeout(loadDataWrapper, 1000 * 3)
  } else {
    setTimeout(loadDataWrapper, 1000 * 60 * 10); // 每10分钟重新加载数据
  }
}

const base64ToBlob = (base64, contentType) => {
  const byteCharacters = atob(base64);
  const byteArrays = [];
  for (let i = 0; i < byteCharacters.length; i++) {
    byteArrays.push(byteCharacters.charCodeAt(i));
  }
  const byteArray = new Uint8Array(byteArrays);
  return new Blob([byteArray], {type: contentType});
}

async function loadData() {
  const fileBase64 = await GetFileBase64()
  if (!fileBase64) {
    console.error("Failed to load data files.");
    return;
  }

  const fileMap = {};
  const fileBlob = base64ToBlob(fileBase64, 'application/zip');
  const zipFile = new zip.ZipReader(new zip.BlobReader(fileBlob));
  for (const file of await zipFile.getEntries()) {
    fileMap[file.filename] = await file.getData(new zip.TextWriter(), {password: password});
  }

  const teachers = Papa.parse(fileMap['teachers.csv'], {
    skipEmptyLines: true,
    header: true
  });

  const gpa = [];
  const gpaFile = JSON.parse(fileMap['gpa.json'])
  for (let teacher of Object.keys(gpaFile)) {
    const courses = gpaFile[teacher];
    for (let course of courses) {
      gpa.push([teacher, ...course]);
    }
  }

  const comments = {};
  for (let entry of Object.keys(fileMap)) {
    if (entry.startsWith("comment_")) {
      // 使用正则表达式提取单位名称
      const unitName = entry.match(/^comment_(.*?)\.csv$/)?.[1];
      if (unitName) {
        comments[unitName] = Papa.parse(fileMap[entry], {
          skipEmptyLines: true,
          header: true
        }).data;
      }
    }
  }

  appData.teachers = teachers.data.map(t => ({
    ...t,
    热度: parseInt(t.热度),
    评分人数: parseInt(t.评分人数),
    评分: parseFloat(t.评分).toFixed(2)
  }))
  appData.gpaData = gpa
  appData.comments = comments
}

onMounted(() => {
  loadDataWrapper();
})
</script>

<template>
  <div v-if="status === 0 && appData.teachers.length !== 0">
    <home-button/>
    <div class="search-box-wrapper">
      <search-box :teachers="appData.teachers"/>
    </div>
    <router-view/>
  </div>
  <div v-else class="loading-container">
    <div class="loading-indicator">
      <div class="spinner"></div>
      <div class="loading-text">{{ statusStr }}</div>
    </div>
  </div>
</template>

<style>

.loading-container {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(3px);
}

.loading-indicator {
  text-align: center;
  transform: translateY(-20%);
}

.spinner {
  width: 48px;
  height: 48px;
  margin: 0 auto 16px;
  border: 4px solid var(--primary-color);
  border-radius: 50%;
  border-top-color: transparent;
  animation: spin 1s linear infinite;
}

.loading-text {
  color: var(--secondary-color);
  font-size: 1.2em;
  font-weight: 500;
  letter-spacing: 0.05em;
  background: linear-gradient(45deg, var(--primary-color), var(--secondary-color));
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  padding: 0 8px;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

:root {
  --primary-color: #4361ee;
  --secondary-color: #3a0ca3;
  --background: #f8f9fa;
}

body {
  font-family: 'Helvetica Neue', sans-serif;
  margin: 0;
  background: var(--background);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.search-box-wrapper {
  position: fixed;
  top: 20px;
  width: 100%;
  z-index: 999;
}
</style>