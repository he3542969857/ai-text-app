<script setup lang="ts">
import { computed, ref, watch, nextTick } from 'vue'

// StreamOutput 展示流式输出文本:
//  - 打字机效果:文本随 props.text 增量更新;
//  - 虚拟滚动:行数超过阈值时只渲染可视区行(加分项②大文本性能优化)。
const props = withDefaults(
  defineProps<{
    text: string
    autoScroll?: boolean
    itemHeight?: number
    viewportHeight?: number
  }>(),
  { autoScroll: true, itemHeight: 24, viewportHeight: 360 },
)

const lines = computed(() => props.text.split('\n'))
const VIRTUAL_THRESHOLD = 100

const scrollTop = ref(0)
const container = ref<HTMLElement | null>(null)

const useVirtual = computed(() => lines.value.length > VIRTUAL_THRESHOLD)
const totalHeight = computed(() => lines.value.length * props.itemHeight)

const visibleCount = computed(() => Math.ceil(props.viewportHeight / props.itemHeight) + 4)
const startIndex = computed(() =>
  Math.max(0, Math.floor(scrollTop.value / props.itemHeight) - 2),
)
const visibleLines = computed(() =>
  lines.value.slice(startIndex.value, startIndex.value + visibleCount.value),
)
const offsetY = computed(() => startIndex.value * props.itemHeight)

function onScroll(e: Event) {
  scrollTop.value = (e.target as HTMLElement).scrollTop
}

// 新内容到来时自动滚到底(打字机跟随)。
watch(
  () => props.text,
  async () => {
    if (!props.autoScroll || !container.value) return
    await nextTick()
    container.value.scrollTop = container.value.scrollHeight
  },
)
</script>

<template>
  <div
    ref="container"
    class="stream-output"
    :style="{ height: viewportHeight + 'px' }"
    @scroll="onScroll"
  >
    <!-- 虚拟滚动:撑高总高度,只渲染可视区切片 -->
    <template v-if="useVirtual">
      <div :style="{ height: totalHeight + 'px', position: 'relative' }">
        <div :style="{ transform: `translateY(${offsetY}px)` }">
          <div
            v-for="(line, i) in visibleLines"
            :key="startIndex + i"
            class="line"
            :style="{ height: itemHeight + 'px' }"
          >
            {{ line || ' ' }}
          </div>
        </div>
      </div>
    </template>
    <!-- 普通渲染 -->
    <template v-else>
      <pre class="plain">{{ text || '等待输出…' }}</pre>
    </template>
  </div>
</template>

<style scoped>
.stream-output {
  overflow-y: auto;
  border: 1px solid var(--n-border-color, #e0e0e6);
  border-radius: 6px;
  padding: 12px;
  background: rgba(128, 128, 128, 0.04);
  font-family: 'Menlo', 'Consolas', monospace;
  font-size: 14px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-word;
}
.line {
  white-space: pre-wrap;
  word-break: break-word;
}
.plain {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
  font-family: inherit;
}
</style>
