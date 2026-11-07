<template>
  <div>
    <div class="admin-stat-grid">
      <div class="admin-stat"><div class="num">{{ data.activityCount || 0 }}</div><div class="label">活动数量</div></div>
      <div class="admin-stat"><div class="num">{{ data.strategyCount || 0 }}</div><div class="label">策略数量</div></div>
      <div class="admin-stat"><div class="num">{{ data.drawToday || 0 }}</div><div class="label">今日抽奖次数</div></div>
      <div class="admin-stat"><div class="num">{{ data.awardToday || 0 }}</div><div class="label">今日中奖笔数</div></div>
      <div class="admin-stat"><div class="num">{{ data.activityOrderCount || 0 }}</div><div class="label">额度订单</div></div>
      <div class="admin-stat"><div class="num">{{ data.awardCount || 0 }}</div><div class="label">累计中奖</div></div>
      <div class="admin-stat"><div class="num">{{ data.creditOrderCount || 0 }}</div><div class="label">积分流水</div></div>
      <div class="admin-stat"><div class="num">{{ data.taskCount || 0 }}</div><div class="label">异步任务</div></div>
    </div>

    <div class="admin-card mt12">
      <div class="admin-toolbar"><strong>最近中奖记录</strong><el-button size="small" @click="load">刷新</el-button></div>
      <el-table :data="data.recentAwards || []" size="small" border>
        <el-table-column prop="userId" label="用户" width="110" />
        <el-table-column prop="awardId" label="奖品ID" width="90" />
        <el-table-column prop="awardTitle" label="奖品" min-width="140" />
        <el-table-column prop="awardState" label="状态" width="100" />
        <el-table-column prop="createTime" label="时间" width="180" />
      </el-table>
    </div>

    <div class="admin-card mt12">
      <div class="admin-toolbar"><strong>最近额度订单</strong></div>
      <el-table :data="data.recentOrders || []" size="small" border>
        <el-table-column prop="userId" label="用户" width="110" />
        <el-table-column prop="sku" label="SKU" width="90" />
        <el-table-column prop="activityName" label="活动" min-width="140" />
        <el-table-column prop="totalCount" label="次数" width="80" />
        <el-table-column prop="payAmount" label="积分" width="90" />
        <el-table-column prop="state" label="状态" width="100" />
        <el-table-column prop="createTime" label="时间" width="180" />
      </el-table>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { adminApi, toastError } from '../../api'

const data = ref({})
async function load() {
  try {
    data.value = await adminApi.dashboard()
  } catch (e) {
    toastError(e)
  }
}
onMounted(load)
</script>
