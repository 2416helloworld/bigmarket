<template>
  <div class="admin-card">
    <div class="admin-toolbar">
      <el-radio-group v-model="type" size="small" @change="load">
        <el-radio-button label="activity">额度订单</el-radio-button>
        <el-radio-button label="raffle">抽奖单</el-radio-button>
        <el-radio-button label="award">中奖记录</el-radio-button>
        <el-radio-button label="credit">积分流水</el-radio-button>
        <el-radio-button label="rebate">返利订单</el-radio-button>
      </el-radio-group>
      <el-input v-model="userId" placeholder="按用户ID筛选" clearable size="small" style="width: 200px" @keyup.enter="load" />
      <el-button size="small" @click="load">查询</el-button>
    </div>
    <el-table :data="rows" size="small" border v-loading="loading" max-height="600">
      <el-table-column v-for="col in columns" :key="col.prop" :prop="col.prop" :label="col.label" :width="col.width" :min-width="col.minWidth || 100" show-overflow-tooltip />
    </el-table>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { adminApi, toastError } from '../../api'

const type = ref('activity')
const userId = ref('')
const rows = ref([])
const loading = ref(false)

const configs = {
  activity: [
    { prop: 'userId', label: '用户', width: 110 },
    { prop: 'sku', label: 'SKU', width: 90 },
    { prop: 'activityName', label: '活动', minWidth: 150 },
    { prop: 'totalCount', label: '总次数', width: 90 },
    { prop: 'dayCount', label: '日次数', width: 90 },
    { prop: 'monthCount', label: '月次数', width: 90 },
    { prop: 'payAmount', label: '积分', width: 90 },
    { prop: 'state', label: '状态', width: 100 },
    { prop: 'outBusinessNo', label: '外部单号', minWidth: 170 },
    { prop: 'createTime', label: '创建时间', width: 180 }
  ],
  raffle: [
    { prop: 'userId', label: '用户', width: 110 },
    { prop: 'activityId', label: '活动ID', width: 100 },
    { prop: 'strategyId', label: '策略ID', width: 100 },
    { prop: 'orderId', label: '抽奖单号', minWidth: 180 },
    { prop: 'orderState', label: '状态', width: 100 },
    { prop: 'createTime', label: '创建时间', width: 180 }
  ],
  award: [
    { prop: 'userId', label: '用户', width: 110 },
    { prop: 'awardId', label: '奖品ID', width: 90 },
    { prop: 'awardTitle', label: '奖品', minWidth: 150 },
    { prop: 'awardState', label: '发放状态', width: 110 },
    { prop: 'orderId', label: '抽奖单号', minWidth: 180 },
    { prop: 'createTime', label: '时间', width: 180 }
  ],
  credit: [
    { prop: 'userId', label: '用户', width: 110 },
    { prop: 'tradeName', label: '交易名称', width: 120 },
    { prop: 'tradeType', label: '方向', width: 100 },
    { prop: 'tradeAmount', label: '积分', width: 100 },
    { prop: 'outBusinessNo', label: '外部单号', minWidth: 170 },
    { prop: 'createTime', label: '时间', width: 180 }
  ],
  rebate: [
    { prop: 'userId', label: '用户', width: 110 },
    { prop: 'behaviorType', label: '行为', width: 100 },
    { prop: 'rebateDesc', label: '返利说明', minWidth: 160 },
    { prop: 'rebateType', label: '类型', width: 100 },
    { prop: 'rebateAmount', label: '金额', width: 90 },
    { prop: 'outBusinessNo', label: '外部单号', minWidth: 170 }
  ]
}
const columns = computed(() => configs[type.value] || [])

async function load() {
  loading.value = true
  try {
    rows.value = (await adminApi.orders(type.value, userId.value)) || []
  } catch (e) {
    toastError(e)
  } finally {
    loading.value = false
  }
}
onMounted(load)
</script>
