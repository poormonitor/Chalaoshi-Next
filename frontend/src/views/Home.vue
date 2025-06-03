<!-- views/Home.vue -->
<template>
  <div class="container">
    <transition-group name="fade" tag="div" class="grid">
      <teacher-card
          v-for="teacher in randomizedTeachers"
          :key="teacher.id"
          :teacher="teacher"
          @click="navigateToDetail(teacher.id)"
      />
    </transition-group>
  </div>
</template>

<script setup>
import {computed, inject} from 'vue'
import TeacherCard from '../components/TeacherCard.vue'

import {useRouter} from 'vue-router'

const appData = inject('appData')
const router = useRouter()

const randomizedTeachers = computed(() => {
  let teachers = appData.teachers.sort(t => -t.评分).slice(0, appData.teachers.length / 5)
  return teachers.sort(() => Math.random() - 0.5).slice(0, 10)
})

function navigateToDetail(id) {
  router.push(`/teacher/${id}`)
}
</script>

<style scoped>
.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
  margin-top: 5rem;
}

.grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1.5rem;
  padding: 1rem;
}
</style>