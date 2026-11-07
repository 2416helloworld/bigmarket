<template>
  <div class="admin-shell">
    <aside class="admin-aside">
      <div class="admin-logo">🎯 大营销后台</div>
      <el-menu class="admin-menu" :default-active="route.path" router background-color="#1f2733" text-color="#cfd6df" active-text-color="#f6c453">
        <el-menu-item index="/admin/dashboard">数据看板</el-menu-item>
        <el-menu-item index="/admin/activities">活动管理</el-menu-item>
        <el-menu-item index="/admin/strategies">策略管理</el-menu-item>
        <el-menu-item index="/admin/awards">奖品配置</el-menu-item>
        <el-menu-item index="/admin/rule-trees">规则树</el-menu-item>
        <el-menu-item index="/admin/orders">订单记录</el-menu-item>
        <el-menu-item index="/admin/users">用户资产</el-menu-item>
        <el-menu-item index="/admin/tasks">异步任务</el-menu-item>
      </el-menu>
    </aside>
    <main class="admin-main">
      <div class="admin-header">
        <h2>{{ pageTitle }}</h2>
        <div style="display: flex; gap: 12px; align-items: center">
          <el-button size="small" tag="a" href="/" target="_blank">打开 H5</el-button>
          <span style="font-size: 13px; color: #606266">{{ adminName }}</span>
          <el-button size="small" @click="logout">退出</el-button>
        </div>
      </div>
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const adminName = localStorage.getItem('bm_admin_name') || 'admin'
const pageTitle = computed(() => {
  const map = {
    '/admin/dashboard': '数据看板',
    '/admin/activities': '活动管理',
    '/admin/strategies': '策略管理',
    '/admin/awards': '奖品配置',
    '/admin/rule-trees': '规则树',
    '/admin/orders': '订单记录',
    '/admin/users': '用户资产',
    '/admin/tasks': '异步任务'
  }
  if (route.path.startsWith('/admin/strategies/')) return '策略详情'
  return map[route.path] || '运营后台'
})
function logout() {
  localStorage.removeItem('bm_admin_token')
  router.push('/admin/login')
}
</script>
