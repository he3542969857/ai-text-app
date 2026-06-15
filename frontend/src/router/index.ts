import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  { path: '/', name: 'list', component: () => import('../views/ListView.vue') },
  { path: '/translate', name: 'translate', component: () => import('../views/TranslateView.vue') },
  { path: '/summarize', name: 'summarize', component: () => import('../views/SummarizeView.vue') },
  { path: '/history', name: 'history', component: () => import('../views/HistoryView.vue') },
]

export default createRouter({
  history: createWebHistory(),
  routes,
})
