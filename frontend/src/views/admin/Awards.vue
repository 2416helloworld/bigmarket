<template>
  <CrudTable resource="awards" title="奖品配置" :columns="columns" :fields="fields" />
</template>

<script setup>
import CrudTable from './CrudTable.vue'

const columns = [
  { prop: 'award_id', label: '奖品ID', width: 100 },
  { prop: 'award_key', label: '发放类型', width: 180 },
  { prop: 'award_config', label: '发放配置', minWidth: 180 },
  { prop: 'award_desc', label: '奖品说明', minWidth: 180 }
]
const fields = [
  { prop: 'award_id', label: '奖品ID', type: 'number' },
  {
    prop: 'award_key',
    label: '发放类型',
    type: 'select',
    options: [
      { label: '随机积分', value: 'user_credit_random' },
      { label: '黑名单积分', value: 'user_credit_blacklist' },
      { label: 'OpenAI 使用次数', value: 'openai_use_count' },
      { label: 'OpenAI 模型', value: 'openai_model' }
    ]
  },
  { prop: 'award_config', label: '发放配置' },
  { prop: 'award_desc', label: '奖品说明' }
]
</script>
