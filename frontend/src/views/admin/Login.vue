<template>
  <div style="min-height: 100vh; display: flex; align-items: center; justify-content: center; background: linear-gradient(135deg, #1f2733, #3c1f1f)">
    <div class="admin-card" style="width: 380px">
      <h2 style="margin: 0 0 6px">大营销运营后台</h2>
      <p style="color: #909399; font-size: 13px; margin-top: 0">Go + Vue 全栈重构版</p>
      <el-form :model="form" label-position="top" @submit.prevent="login">
        <el-form-item label="账号">
          <el-input v-model="form.username" placeholder="admin" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="form.password" type="password" show-password placeholder="admin123" @keyup.enter="login" />
        </el-form-item>
        <el-button type="primary" style="width: 100%" :loading="loading" @click="login">登录</el-button>
      </el-form>
      <div style="margin-top: 14px; font-size: 12px; color: #909399">默认账号 admin / admin123</div>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { adminApi, toastError } from '../../api'

const router = useRouter()
const form = reactive({ username: 'admin', password: 'admin123' })
const loading = ref(false)

async function login() {
  loading.value = true
  try {
    const data = await adminApi.login(form.username, form.password)
    localStorage.setItem('bm_admin_token', data.token)
    localStorage.setItem('bm_admin_name', data.nickname || data.username)
    ElMessage.success('登录成功')
    router.push('/admin/dashboard')
  } catch (e) {
    toastError(e)
  } finally {
    loading.value = false
  }
}
</script>
