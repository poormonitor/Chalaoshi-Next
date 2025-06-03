<!-- components/CommentItem.vue -->
<template>
  <article class="comment-item">
    <div class="header">
      <span class="time">{{ formattedTime }}</span>
      <div class="votes">
        <div class="vote-item">
          <span class="icon">👍</span>
          <span class="upvote">{{ comment.点赞量 }}</span>
        </div>
        <div class="vote-item">
          <span class="icon">👎</span>
          <span class="downvote">{{ comment.点踩量 }}</span>
        </div>
        <span
            class="net-vote"
            :class="{ positive: netVotes > 0, negative: netVotes < 0 }"
        >
          {{ netVotes >= 0 ? '+' : '' }}{{ netVotes }}
        </span>
      </div>
    </div>
    <div class="content" v-html="processedContent"></div>
  </article>
</template>

<script setup>
/* 逻辑部分保持不变 */
import { computed } from 'vue'
import dayjs from 'dayjs'

const props = defineProps({
  comment: {
    type: Object,
    required: true
  }
})

const netVotes = computed(() =>
    props.comment.点赞减去点踩数量 ||
    (props.comment.点赞量 - props.comment.点踩量)
)

const processedContent = computed(() => {
  return props.comment.内容.replace(/\\n\\n/g, '<br />')
})

const formattedTime = computed(() => {
  const date = dayjs(props.comment.发表时间)
  return date.isValid() ? date.format("YYYY-MM-DD") : '未知时间'
})
</script>

<style scoped>
.comment-item {
  background: white;
  padding: 1.5rem;
  border-radius: 12px;
  margin-bottom: 1.25rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  transition: all 0.2s ease;
  border: 1px solid #e9ecef;
}

.comment-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.25rem;
  padding-bottom: 0.75rem;
  border-bottom: 1px solid #e9ecef;
}

.time {
  color: #6c757d;
  font-size: 0.85rem;
  font-weight: 500;
}

.votes {
  display: flex;
  gap: 1.25rem;
  align-items: center;
}

.vote-item {
  display: flex;
  align-items: center;
  gap: 0.4rem;
}

.icon {
  font-size: 0.9em;
  opacity: 0.8;
}

.upvote {
  color: #2a9d8f;
  font-weight: 500;
}

.downvote {
  color: #e76f51;
  font-weight: 500;
}

.net-vote {
  font-weight: 600;
  padding: 0.25rem 0.75rem;
  border-radius: 8px;
  background: #f1f3f5;
  font-size: 0.85rem;
  margin-left: 0.5rem;
}

.net-vote.positive {
  background: #d4f3e6;
  color: #2a9d8f;
}

.net-vote.negative {
  background: #ffe5d9;
  color: #e76f51;
}

.content {
  color: #495057;
  line-height: 1.7;
  font-size: 0.95rem;
  white-space: pre-wrap;
}

@media (max-width: 480px) {
  .comment-item {
    padding: 1.2rem;
    border-radius: 8px;
  }

  .votes {
    gap: 1rem;
  }

  .net-vote {
    padding: 0.2rem 0.6rem;
  }
}
</style>