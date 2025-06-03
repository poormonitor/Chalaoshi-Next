<!-- views/TeacherDetail.vue -->
<template>
  <div class="main-content">
    <div class="container" v-if="teacher">
      <div class="header">
        <div class="header-content">
          <div class="teacher-info">
            <div class="name">{{ teacher.姓名 }}</div>
            <div class="meta-group">
              <span class="college">{{ teacher.学院 }}</span>
              <span class="divider">|</span>
              <span class="course-count">{{ gpaCourses.length }} 门课程</span>
            </div>
          </div>

          <div class="rating-card">
            <div class="rating-label">综合评价</div>
            <div class="rating-value">
              {{ teacher.评分 }}
              <span class="rating-scale">/10</span>
            </div>
          </div>
        </div>
      </div>

      <div class="content">
        <section class="gpa-section">
          <h2>课程绩点</h2>
          <gpa-area :data="gpaCourses"/>
        </section>

        <section class="comments-section">
          <h2>学生评价（{{ comments.length }}条）</h2>
          <div v-for="comment in comments" :key="comment.评论id" class="comment-item">
            <comment-item :comment="comment"/>
          </div>
        </section>
      </div>
    </div>
  </div>
</template>

<script setup>
import {computed, inject, onMounted} from 'vue'
import GpaArea from '../components/GpaArea.vue'
import CommentItem from '../components/CommentItem.vue'

const props = defineProps({
  id: {
    type: [String, Number],
    required: true
  }
})

const appData = inject('appData')

const teacher = computed(() =>
    appData.teachers.find(t => t.id === props.id)
)

const gpaCourses = computed(() =>
    appData.gpaData.filter(t => t[0] === teacher.value.姓名)?.sort((a, b) => b[1] - a[1]) || []
)

const comments = computed(() => {
  return appData.comments[teacher.value.学院].filter(c => c.老师id === teacher.value.id)
          .sort((a, b) => new Date(b.发表时间) - new Date(a.发表时间))
      || []
})

// get to the top of the page when activated
onMounted(() => {
  window.scrollTo(0, 0)
})
</script>

<style scoped>
.header {
  background: linear-gradient(135deg, #4a90e2, #63b3ed);
  color: white;
  border-radius: 1rem;
  margin-bottom: 1.5rem;
  padding: 1.25rem;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
}

.teacher-info {
  flex: 1;
}

.name {
  font-size: 1.75rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
  letter-spacing: -0.5px;
}

.meta-group {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  font-size: 0.95rem;
  opacity: 0.9;
  margin-bottom: 0.6rem;
}

.divider {
  opacity: 0.5;
}

.rating-card {
  background: rgba(255, 255, 255, 0.15);
  padding: 0.75rem 1.25rem;
  border-radius: 0.75rem;
  min-width: 120px;
  text-align: center;
  backdrop-filter: blur(4px);
}

.rating-label {
  font-size: 0.9rem;
  margin-bottom: 0.25rem;
  opacity: 0.9;
}

.rating-value {
  font-size: 2rem;
  font-weight: 700;
  line-height: 1;
  display: flex;
  align-items: baseline;
  justify-content: center;
}

.rating-scale {
  font-size: 1rem;
  margin-left: 0.25rem;
  opacity: 0.8;
}

@media (max-width: 768px) {
  .header-content {
    flex-direction: column;
    align-items: flex-start;
  }

  .name {
    font-size: 1.5rem;
  }

  .rating-card {
    width: 100%;
    padding: 0.75rem;
  }

  .rating-value {
    font-size: 1.75rem;
  }
}

.main-content {
  margin: 6rem 2rem 2rem;
}

.content {
  display: grid;
  gap: 1rem;
}

.comments-section, .gpa-section {
  padding: 0 1rem;
}
</style>