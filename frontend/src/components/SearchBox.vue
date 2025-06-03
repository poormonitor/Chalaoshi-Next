<!-- components/SearchBox.vue -->
<template>
  <div class="search-container" ref="searchContainer">
    <div class="search-input">
      <input
          type="text"
          v-model="query"
          @input="handleInput"
          @focus="showSuggestions = true"
          placeholder="输入教师姓名全拼或缩写"
      />
      <transition name="fade">
        <ul v-if="showSuggestions && suggestions.length" class="suggestions">
          <li
              v-for="teacher in suggestions"
              :key="teacher.id"
              @click="selectTeacher(teacher)"
          >
            <span class="name">{{ teacher.姓名 }}</span>
            <span class="meta">
              {{ teacher.学院 }} · 评分{{ teacher.评分 }}（{{ teacher.评分人数 }}人）
            </span>
          </li>
        </ul>
      </transition>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useDebounceFn } from '@vueuse/core'
import { onClickOutside } from '@vueuse/core'

const props = defineProps({
  teachers: {
    type: Array,
    required: true
  }
})

const emit = defineEmits(['search'])

const router = useRouter()
const query = ref('')
const showSuggestions = ref(false)
const searchContainer = ref(null)

// 使用VueUse的防抖函数
const debouncedSearch = useDebounceFn(() => {
  emit('search', query.value)
}, 300)

const suggestions = computed(() => {
  const q = query.value.toLowerCase()
  if (!q) return []
  return props.teachers.filter(teacher =>
      teacher.拼音.includes(q) ||
      teacher.拼音缩写.includes(q) ||
      teacher.姓名.includes(q)
  ).slice(0, 5)
})

const handleInput = () => {
  debouncedSearch()
}

const selectTeacher = (teacher) => {
  router.push(`/teacher/${teacher.id}`)
  closeSuggestions()
}

const closeSuggestions = () => {
  showSuggestions.value = false
}

onClickOutside(searchContainer, () => {
  closeSuggestions()
})
</script>

<style scoped>
.search-container {
  max-width: 600px;
  margin: 0 auto 2rem;
  position: relative;
}

.search-input {
  position: relative;
}

input {
  width: 90%;
  padding: 1rem 2rem;
  border-radius: 2rem;
  border: 2px solid var(--primary-color);
  font-size: 1.1rem;
  transition: box-shadow 0.3s ease;
}

input:focus {
  outline: none;
  box-shadow: 0 4px 12px rgba(67, 97, 238, 0.3);
}

.suggestions {
  position: absolute;
  width: 100%;
  max-height: 400px;
  overflow-y: auto;
  background: white;
  border-radius: 1rem;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  margin-top: 0.5rem;
  padding: 0.5rem 0;
  z-index: 100;
}

.suggestions li {
  padding: 1rem;
  cursor: pointer;
  transition: background 0.2s;
}

.suggestions li:hover {
  background: #f8f9fa;
}

.name {
  font-weight: 600;
  display: block;
  margin-bottom: 0.25rem;
}

.meta {
  font-size: 0.9em;
  color: #666;
}
</style>