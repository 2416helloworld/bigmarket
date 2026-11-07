<template>
  <div>
    <div class="admin-card">
      <div class="admin-toolbar">
        <strong>用户积分调整</strong>
        <el-input v-model="form.userId" placeholder="用户ID" size="small" style="width: 160px" />
        <el-input-number v-model="form.amount" :min="0.01" :precision="2" size="small" />
        <el-select v-model="form.tradeType" size="small" style="width: 120px">
          <el-option label="增加" value="forward" />
          <el-option label="扣减" value="reverse" />
        </el-select>
        <el-button type="primary" size="small" @click="adjust">提交</el-button>
        <el-button size="small" @click="load">刷新</el-button>
      </div>
      <el-table :data="data.creditAccounts || []" size="small" border max-height="300">
        <el-table-column prop="userId" label="用户" width="130" />
        <el-table-column prop="totalAmount" label="累计积分" width="120" />
        <el-table-column prop="availableAmount" label="可用积分" width="120" />
        <el-table-column prop="accountStatus" label="状态" width="100" />
        <el-table-column prop="updateTime" label="更新时间" width="180" />
      </el-table>
    </div>

    <div class="admin-card mt12">
      <div class="admin-toolbar"><strong>用户活动账户</strong></div>
      <el-table :data="data.activityAccounts || []" size="small" border max-height="360">
        <el-table-column prop="userId" label="用户" width="120" />
        <el-table-column prop="activityId" label="活动ID" width="110" />
        <el-table-column prop="totalCount" label="总次数" width="90" />
        <el-table-column prop="totalCountSurplus" label="总剩余" width="90" />
        <el-table-column prop="dayCountSurplus" label="日剩余" width="90" />
        <el-table-column prop="monthCountSurplus" label="月剩余" width="90" />
        <el-table-column prop="updateTime" label="更新时间" width="180" />
      </el-table>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { adminApi, toastError } from '../../api'

const data = ref({})
const form = reactive({ userId: '', amount: 10, tradeType: 'forward' })

async function load() {
  try {
    data.value = await adminApi.users()
  } catch (e) {
    toastError(e)
  }
}

async function adjust() {
  if (!form.userId) {
    ElMessage.warning('请输入用户ID')
    return
  }
  try {
    await adminApi.creditAdjust(form.userId, form.amount, form.tradeType)
    ElMessage.success('积分调整成功')
    await load()
  } catch (e) {
    toastError(e)
  }
}

onMounted(load)
</script>
