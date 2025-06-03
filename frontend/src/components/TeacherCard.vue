<!-- components/TeacherCard.vue -->
<template>
  <article class="teacher-card" @click="navigateToDetail">
    <div class="content">
      <h3 class="name">{{ teacher.姓名 }}</h3>
      <div class="stats">
        <div class="rating">
          <div class="value">{{ teacher.评分 }}</div>
          <div class="count">({{ teacher.评分人数 }}人评价)</div>
        </div>
        <div class="college">{{ teacher.学院 }}</div>
      </div>
    </div>
    <div class="heat" :style="heatStyle">
      {{ heatText }}
    </div>
  </article>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'

const props = defineProps({
  teacher: {
    type: Object,
    required: true
  }
})

const router = useRouter()

const heatStyle = computed(() => ({
  backgroundColor: `hsl(${props.teacher.评分 * 36}, 70%, 50%)`
}))

const heatText = computed(() => {
  if (props.teacher.评分 > 9.0) return '超好评'
  if (props.teacher.评分 > 7.0) return '好评'
  return '一般'
})

const navigateToDetail = () => {
  router.push(`/teacher/${props.teacher.id}`)
}
</script>

<style scoped>
.teacher-card {
  background: white;
  border-radius: 1rem;
  padding: 1.5rem;
  cursor: pointer;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  position: relative;
  overflow: hidden;
}

.teacher-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.1);
}

.content {
  position: relative;
  z-index: 1;
}

.name {
  color: var(--secondary-color);
  margin-bottom: 1rem;
}

.stats {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.rating .value {
  font-size: 1.5em;
  font-weight: 700;
  color: var(--primary-color);
}

.rating .count {
  color: #666;
  font-size: 0.9em;
}

.college {
  background: var(--primary-color);
  color: white;
  padding: 0.25rem 0.75rem;
  border-radius: 1rem;
  font-size: 0.8em;
}

.heat {
  position: absolute;
  right: -2rem;
  top: -2rem;
  width: 6rem;
  height: 6rem;
  border-radius: 50%;
  transform: rotate(45deg);
  color: white;
  display: flex;
  align-items: flex-end;
  justify-content: center;
  padding-bottom: 1rem;
  font-size: 0.9em;
  opacity: 0.3;
  transition: opacity 0.3s ease;
}

.teacher-card:hover .heat {
  opacity: 0.35;
}
</style>