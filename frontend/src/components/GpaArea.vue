<template>
  <div class="gpa-container">
    <div class="gpa-grid">
      <div
          v-for="(course, index) in processedData"
          :key="index"
          class="gpa-card"
          @click="showCourseTeachers(course[1])"
      >
        <div class="card-content">
          <div class="card-header">
            <span class="course-name">{{ course[1] }}</span>
            <span class="gpa-value">{{ course[2] }}</span>
          </div>
          <div class="card-meta">
            <span class="rating">
              👤{{ course[3] }}
            </span>
            <span class="std-dev">±{{ course[4] }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 教师详情模态框 -->
    <transition name="modal-fade">
      <div
          v-if="selectedCourse"
          class="modal-backdrop"
          @click.self="closeModal"
      >
        <div class="modal-box">
          <div class="modal-header">
            <h2 class="modal-title">{{ selectedCourse }} 授课教师</h2>
            <button
                class="modal-close"
                @click="closeModal"
            >
              <svg width="24" height="24" viewBox="0 0 24 24">
                <path
                    d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"/>
              </svg>
            </button>
          </div>

          <div class="modal-content">
            <div
                v-if="teachers.length === 0"
                class="no-teachers"
            >
              <div class="no-data-icon">📭</div>
              <p>暂无该课程的教师信息</p>
            </div>

            <div
                v-else
                class="teacher-grid"
            >
              <div
                  v-for="teacher in teachers"
                  :key="teacher[0]"
                  class="teacher-card"
              >
                <div class="teacher-header">
                  <span class="teacher-name">{{ teacher[0] }}</span>
                </div>

                <div class="teacher-stats">
                  <div class="stat-item">
                    <span class="stat-label">平均绩点</span>
                    <span class="stat-value gpa">{{ teacher[2] }} ± {{ teacher[4] }}</span>
                  </div>
                  <div class="stat-item">
                    <span class="stat-label">评价人数</span>
                    <span class="stat-value count">{{ teacher[3] }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import {ref, computed, inject} from 'vue'

const props = defineProps({
  data: {
    type: Array,
    required: true,
    default: () => []
  }
})

// 从上下文注入教师数据
const appData = inject('appData')
const selectedCourse = ref(null)
const teachers = ref([])

const processedData = computed(() => {
  return [...props.data].sort((a, b) => b[3] - a[3])
})

// 显示指定课程的教师
const showCourseTeachers = (courseName) => {
  selectedCourse.value = courseName
  teachers.value = appData.gpaData.filter(gpa => gpa[1] === courseName)
      .sort((a, b) => b[2] - a[2])
}

// 关闭模态框
const closeModal = () => {
  selectedCourse.value = null
  teachers.value = []
}
</script>

<style scoped>
/* 原始课程列表样式 */
.gpa-container {
  --min-column-width: min(320px, 90vw);
  --gap: 0.75rem;
}

.gpa-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(var(--min-column-width), 1fr));
  gap: var(--gap);
}

.gpa-card {
  background: #ffffff;
  border-radius: 6px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
  transition: transform 0.2s ease;
  cursor: pointer;
  break-inside: avoid;
}

.gpa-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.card-content {
  padding: 0.625rem 0.875rem;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 0.25rem;
}

.course-name {
  flex: 1;
  font-weight: 500;
  color: #2d3436;
  font-size: 0.9rem;
  line-height: 1.3;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.gpa-value {
  font-size: 0.95rem;
  font-weight: 700;
  color: #4a90e2;
  min-width: 2.8em;
  text-align: right;
}

.card-meta {
  display: flex;
  justify-content: space-between;
  font-size: 0.8rem;
  color: #636e72;
}

.rating {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
}

.std-dev {
  opacity: 0.8;
}

/* 模态框样式 */
.modal-backdrop {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-box {
  background: white;
  border-radius: 12px;
  width: min(90vw, 800px);
  max-height: 80vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
}

.modal-header {
  padding: 1.5rem;
  border-bottom: 1px solid #eee;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-title {
  margin: 0;
  font-size: 1.4rem;
  color: #2c3e50;
}

.modal-close {
  background: none;
  border: none;
  padding: 0.5rem;
  cursor: pointer;
  color: #95a5a6;
  transition: color 0.2s;
}

.modal-close:hover {
  color: #e74c3c;
}

.modal-close svg {
  display: block;
  fill: currentColor;
}

.modal-content {
  padding: 1.5rem;
  overflow-y: auto;
}

/* 教师列表 */
.teacher-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 1rem;
}

.teacher-card {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 1.25rem;
  transition: transform 0.2s;
}

.teacher-card:hover {
  transform: translateY(-2px);
  background: #f1f3f5;
}

.teacher-header {
  margin-bottom: 1rem;
}

.teacher-name {
  font-weight: 600;
  color: #2d3436;
  font-size: 1rem;
  display: block;
  margin-bottom: 0.25rem;
}

.teacher-department {
  font-size: 0.85rem;
  color: #636e72;
}

.teacher-stats {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1rem;
}

.stat-item {
  text-align: center;
  padding: 0.75rem;
  background: white;
  border-radius: 6px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
}

.stat-label {
  font-size: 0.85rem;
  color: #7f8c8d;
  margin-bottom: 0.25rem;
  display: block;
}

.stat-value {
  font-weight: 700;
  font-size: 1.1rem;
}

.stat-value.gpa {
  color: #4a90e2;
}

.stat-value.count {
  color: #2ecc71;
}

/* 空状态 */
.no-teachers {
  text-align: center;
  padding: 3rem 1rem;
}

.no-data-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
  opacity: 0.6;
}

.no-teachers p {
  color: #95a5a6;
  margin: 0;
}

/* 过渡动画 */
.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: opacity 0.3s ease;
}

.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}

@media (max-width: 640px) {
  .gpa-grid {
    grid-template-columns: 1fr;
  }

  .modal-box {
    width: 95vw;
  }

  .teacher-grid {
    grid-template-columns: 1fr;
  }
}
</style>