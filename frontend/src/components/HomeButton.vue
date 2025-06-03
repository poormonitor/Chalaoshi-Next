<!-- components/HomeButton.vue -->
<template>
  <div
      class="home-button"
      role="button"
      aria-label="返回主页"
      tabindex="0"
      @click="goHome"
      @keyup.enter="goHome"
  >
    <svg class="icon" viewBox="0 0 36 36">
      <!-- 优化后的居中Home图标 -->
      <path
          d="M18 5.5L7 16.5h3.5v10h5v-7h5v7h5v-10H29L18 5.5Z"
          stroke-linecap="round"
          stroke-linejoin="round"
          transform="translate(1 1)"
      />
    </svg>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()

const goHome = () => {
  router.push('/')
}
</script>

<style scoped>
.home-button {
  position: fixed;
  top: 1rem;
  left: 1rem;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: #4DA8DA;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  z-index: 1000;
}

.home-button:hover {
  background: #378dbc;
  transform: scale(1.1);
}

.icon {
  width: 60%;  /* 控制实际显示尺寸 */
  height: 60%;
  margin-top: -2px; /* 微调垂直居中 */
}

.icon path {
  stroke: #ffffff;
  stroke-width: 1.8;  /* 优化线条粗细 */
  fill: none;
  vector-effect: non-scaling-stroke; /* 保持线条粗细 */
}
</style>