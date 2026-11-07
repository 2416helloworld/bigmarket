<template>
  <div class="admin-card">
    <div class="admin-toolbar">
      <strong>异步任务（MQ 补偿）</strong>
      <el-button size="small" @click="load">刷新</el-button>
    </div>
    <el-table :data="rows" size="small" border v-loading="loading" max-height="600">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="userId" label="用户" width="110" />
      <el-table-column prop="topic" label="Topic" min-width="200" />
      <el-table-column prop="messageId" label="消息ID" width="170" />
      <el-table-column prop="state" label="状态" width="100" />
      <el-table-column prop="createTime" label="创建时间" width="180" />
    </el-table>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { adminApi, toastError } from '../../api'

const rows = ref([])
const loading = ref(false)
async function load() {
  loading.value = true
  try {
    rows.value = (await adminApi.tasks()) || []
  } catch (e) {
    toastError(e)
  } finally {
    loading.value = false
  }
}
onMounted(load)
</script>
