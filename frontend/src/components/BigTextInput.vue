<script setup lang="ts">
import { computed, ref } from 'vue'

// BigTextInput:面向大文本输入的编辑器。
//  - 编辑用 native <textarea>(单节点,浏览器原生高效);
//  - 左侧行号槽采用「虚拟滚动」:无论多少行,只渲染可视区的行号 DOM,
//    随 textarea 滚动同步,避免大文本下数万 DOM 节点导致卡顿(加分项②)。
const props = withDefaults(
  defineProps<{
    modelValue: string
    placeholder?: string
    minRows?: number
    lineHeight?: number
  }>(),
  { placeholder: '', minRows: 6, lineHeight: 22 },
)
const emit = defineEmits<{ 'update:modelValue': [string] }>()

const lh = props.lineHeight
const viewportH = computed(() => Math.max(props.minRows, 6) * lh + 24)

const scrollTop = ref(0)
const ta = ref<HTMLTextAreaElement | null>(null)

const lines = computed(() => props.modelValue.split('\n'))
const lineCount = computed(() => lines.value.length)
const charCount = computed(() => props.modelValue.length)
const VIRTUAL_THRESHOLD = 200
const virtual = computed(() => lineCount.value > VIRTUAL_THRESHOLD)

const totalHeight = computed(() => lineCount.value * lh)
const visibleCount = computed(() => Math.ceil(viewportH.value / lh) + 4)
const startLine = computed(() => Math.max(0, Math.floor(scrollTop.value / lh) - 2))
const visibleNumbers = computed(() =>
  Array.from({ length: Math.min(visibleCount.value, lineCount.value - startLine.value) }, (_, i) => startLine.value + i + 1),
)
const offsetY = computed(() => startLine.value * lh)

function onScroll(e: Event) {
  scrollTop.value = (e.target as HTMLTextAreaElement).scrollTop
}
function onInput(e: Event) {
  emit('update:modelValue', (e.target as HTMLTextAreaElement).value)
}
</script>

<template>
  <div class="bti">
    <div class="bti-body" :style="{ height: viewportH + 'px' }">
      <!-- 虚拟滚动行号槽 -->
      <div class="gutter" :style="{ width: virtual ? '52px' : '40px' }">
        <div :style="{ height: totalHeight + 'px', position: 'relative' }">
          <div :style="{ transform: `translateY(${offsetY}px)` }">
            <div v-for="n in visibleNumbers" :key="n" class="ln" :style="{ height: lh + 'px' }">
              {{ n }}
            </div>
          </div>
        </div>
      </div>
      <textarea
        ref="ta"
        class="ta"
        :value="modelValue"
        :placeholder="placeholder"
        :style="{ lineHeight: lh + 'px' }"
        spellcheck="false"
        @input="onInput"
        @scroll="onScroll"
      ></textarea>
    </div>
    <div class="bti-foot">
      <span>字符 {{ charCount.toLocaleString() }} · 行 {{ lineCount.toLocaleString() }}</span>
      <span v-if="virtual" class="vbadge">⚡️ 虚拟滚动已启用</span>
    </div>
  </div>
</template>

<style scoped>
.bti {
  border: 0.5px solid var(--border);
  border-radius: 14px;
  overflow: hidden;
  background: var(--bg-elevated);
}
.bti-body {
  display: flex;
  overflow: hidden;
}
.gutter {
  flex: none;
  overflow: hidden;
  background: color-mix(in srgb, var(--text) 4%, transparent);
  border-right: 0.5px solid var(--border);
  color: var(--text-tertiary);
  font-family: 'SF Mono', 'Menlo', 'Consolas', monospace;
  font-size: 12px;
  text-align: right;
  user-select: none;
}
.ln {
  padding-right: 8px;
  line-height: v-bind('lh + "px"');
}
.ta {
  flex: 1;
  border: none;
  outline: none;
  resize: none;
  padding: 0 14px;
  background: transparent;
  color: var(--text);
  font-family: 'SF Mono', 'Menlo', 'Consolas', monospace;
  font-size: 14px;
  white-space: pre;
  overflow: auto;
}
.bti-foot {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 7px 14px;
  border-top: 0.5px solid var(--border);
  font-size: 12px;
  color: var(--text-tertiary);
  background: color-mix(in srgb, var(--text) 3%, transparent);
}
.vbadge {
  color: var(--accent);
  font-weight: 500;
}
</style>
