<template>
  <CrudTable resource="strategies" title="抽奖策略" :columns="columns" :fields="fields">
    <template #row-actions="{ row }">
      <el-button link type="primary" size="small" @click="$router.push(`/admin/strategies/${row.strategy_id}`)">配置</el-button>
    </template>
  </CrudTable>
</template>

<script setup>
import CrudTable from './CrudTable.vue'

const columns = [
  { prop: 'strategy_id', label: '策略ID', width: 120 },
  { prop: 'strategy_desc', label: '策略描述', minWidth: 180 },
  { prop: 'rule_models', label: '规则模型', minWidth: 200 }
]
const fields = [
  { prop: 'strategy_id', label: '策略ID', type: 'number' },
  { prop: 'strategy_desc', label: '策略描述' },
  { prop: 'rule_models', label: '规则模型（逗号分隔）', type: 'textarea' }
]
</script>
