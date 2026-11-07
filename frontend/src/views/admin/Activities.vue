<template>
  <div>
    <CrudTable
      resource="activities"
      title="抽奖活动"
      :columns="activityColumns"
      :fields="activityFields"
      :defaults="{ state: 'open', strategy_id: 100006 }"
    >
      <template #row-actions="{ row }">
        <el-button link type="success" size="small" @click="assemble(row)">装配</el-button>
      </template>
    </CrudTable>

    <div class="mt12">
      <CrudTable
        resource="counts"
        title="活动次数配置"
        :columns="countColumns"
        :fields="countFields"
        :defaults="{ total_count: 10, day_count: 10, month_count: 10 }"
      />
    </div>

    <div class="mt12">
      <CrudTable
        resource="skus"
        title="活动商品 SKU"
        :columns="skuColumns"
        :fields="skuFields"
        :defaults="{ stock_count: 1000, stock_count_surplus: 1000, product_amount: 10 }"
      />
    </div>
  </div>
</template>

<script setup>
import { ElMessage } from 'element-plus'
import CrudTable from './CrudTable.vue'
import { adminApi, toastError } from '../../api'

const activityColumns = [
  { prop: 'activity_id', label: '活动ID', width: 100 },
  { prop: 'activity_name', label: '活动名称', minWidth: 150 },
  { prop: 'activity_desc', label: '描述', minWidth: 150 },
  { prop: 'begin_date_time', label: '开始时间', width: 170 },
  { prop: 'end_date_time', label: '结束时间', width: 170 },
  { prop: 'strategy_id', label: '策略ID', width: 100 },
  { prop: 'state', label: '状态', width: 90 }
]
const activityFields = [
  { prop: 'activity_id', label: '活动ID', type: 'number' },
  { prop: 'activity_name', label: '活动名称' },
  { prop: 'activity_desc', label: '活动描述' },
  { prop: 'begin_date_time', label: '开始时间', type: 'datetime' },
  { prop: 'end_date_time', label: '结束时间', type: 'datetime' },
  { prop: 'strategy_id', label: '策略ID', type: 'number' },
  {
    prop: 'state',
    label: '状态',
    type: 'select',
    options: [
      { label: '创建', value: 'create' },
      { label: '开放', value: 'open' },
      { label: '关闭', value: 'close' }
    ]
  }
]

const countColumns = [
  { prop: 'activity_count_id', label: '次数配置ID', width: 130 },
  { prop: 'total_count', label: '总次数', width: 100 },
  { prop: 'day_count', label: '日次数', width: 100 },
  { prop: 'month_count', label: '月次数', width: 100 }
]
const countFields = [
  { prop: 'activity_count_id', label: '次数配置ID', type: 'number' },
  { prop: 'total_count', label: '总次数', type: 'number' },
  { prop: 'day_count', label: '日次数', type: 'number' },
  { prop: 'month_count', label: '月次数', type: 'number' }
]

const skuColumns = [
  { prop: 'sku', label: 'SKU', width: 100 },
  { prop: 'activity_id', label: '活动ID', width: 100 },
  { prop: 'activity_count_id', label: '次数配置ID', width: 130 },
  { prop: 'stock_count', label: '总库存', width: 100 },
  { prop: 'stock_count_surplus', label: '剩余库存', width: 100 },
  { prop: 'product_amount', label: '积分价格', width: 100 }
]
const skuFields = [
  { prop: 'sku', label: 'SKU', type: 'number' },
  { prop: 'activity_id', label: '活动ID', type: 'number' },
  { prop: 'activity_count_id', label: '次数配置ID', type: 'number' },
  { prop: 'stock_count', label: '总库存', type: 'number' },
  { prop: 'stock_count_surplus', label: '剩余库存', type: 'number' },
  { prop: 'product_amount', label: '积分价格', type: 'number', precision: 2, step: 0.01 }
]

async function assemble(row) {
  try {
    await adminApi.armory(row.activity_id)
    ElMessage.success(`活动 ${row.activity_id} 装配完成，概率表与库存已写入 Redis`)
  } catch (e) {
    toastError(e)
  }
}
</script>
