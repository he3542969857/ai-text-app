<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { fetchFunctions, type FunctionDef } from '../api/client'

const router = useRouter()
const functions = ref<FunctionDef[]>([])
const error = ref('')

const emojiFor = (fn: FunctionDef) => {
  if (fn.type === 'summarize') return '📝'
  return fn.params?.from === 'zh' ? '🇨🇳' : '🇬🇧'
}
const routeFor = (type: string) => (type === 'summarize' ? '/summarize' : '/translate')

onMounted(async () => {
  try {
    functions.value = await fetchFunctions()
  } catch (e) {
    error.value = (e as Error).message
  }
})
</script>

<template>
  <section class="hero ap-fade-up">
    <div class="hero-glow"></div>
    <h1 class="ap-title hero-title">文字，<br />更聪明的处理方式。</h1>
    <p class="ap-subtitle hero-sub">中英互译与智能总结 · 由大模型实时流式驱动 ⚡️</p>
  </section>

  <p v-if="error" class="err">⚠️ {{ error }}</p>

  <section class="grid">
    <article
      v-for="(fn, i) in functions"
      :key="i"
      class="ap-card ap-card-hover feature ap-fade-up"
      :style="{ animationDelay: 0.1 + i * 0.08 + 's' }"
      @click="router.push(routeFor(fn.type))"
    >
      <div class="feature-emoji">{{ emojiFor(fn) }}</div>
      <h3 class="feature-name">{{ fn.name }}</h3>
      <p class="feature-desc">{{ fn.description }}</p>
      <span class="feature-go">开始使用 →</span>
    </article>
  </section>
</template>

<style scoped>
.hero {
  position: relative;
  text-align: center;
  padding: 40px 0 56px;
}
.hero-glow {
  position: absolute;
  inset: -80px 0 auto 0;
  height: 360px;
  background: var(--hero-glow);
  pointer-events: none;
  z-index: -1;
}
.hero-title {
  font-size: clamp(34px, 6vw, 56px);
}
.hero-sub {
  font-size: clamp(16px, 2.5vw, 21px);
  margin-top: 18px;
}
.err {
  text-align: center;
  color: #ff453a;
}
.grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 18px;
}
.feature {
  padding: 28px 24px 24px;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  min-height: 200px;
}
.feature-emoji {
  font-size: 40px;
  line-height: 1;
  margin-bottom: 16px;
  filter: drop-shadow(0 4px 10px rgba(0, 0, 0, 0.12));
}
.feature-name {
  font-size: 21px;
  font-weight: 600;
  letter-spacing: -0.02em;
  margin: 0 0 6px;
  color: var(--text);
}
.feature-desc {
  font-size: 15px;
  color: var(--text-secondary);
  margin: 0;
  flex: 1;
}
.feature-go {
  margin-top: 16px;
  font-size: 14px;
  font-weight: 500;
  color: var(--accent);
}

@media (max-width: 720px) {
  .grid {
    grid-template-columns: 1fr;
  }
}
</style>
