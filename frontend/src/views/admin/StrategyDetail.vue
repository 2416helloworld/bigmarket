<template>
  <div>
    <div class="admin-card" style="margin-bottom: 12px">
      <div class="admin-toolbar">
        <strong>策略 {{ strategyId }}</strong>
        <el-button type="success" size="small" @click="assemble">装配概率表 / 库存到 Redis</el-button>
        <el-button size="small" @click="$router.push('/admin/strategies')">返回列表</el-button>
      </div>
      <div style="color: #909399; font-size: 12px">提示：修改奖品概率、库存或规则后，需要点击装配才会生效。</div>
    </div>

    <CrudTable
      resource="strategy-awards"
      title="策略奖品与概率"
      :columns="awardColumns"
      :fields="awardFields"
      :query="query"
      :defaults="{ strategy_id: strategyId, award_count: 1, award_count_surplus: 1, award_rate: 0.1, rule_models: 'tree_luck_award', sort: 1 }"
    />
    <div class="mt12">
      <CrudTable
        resource="strategy-rules"
        title="策略规则（黑名单 / 权重）"
        :columns="ruleColumns"
        :fields="ruleFields"
        :query="query"
        :defaults="{ strategy_id: strategyId, rule_type: 1 }"
      />
    </div>
  </div>
</template>

<script setup>
import { reactive } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import CrudTable from './CrudTable.vue'
import { adminApi, toastError } from '../../api'

const route = useRoute()
const strategyId = Number(route.params.strategyId)
const query = reactive({ strategyId })

const awardColumns = [
  { prop: 'award_id', label: '奖品ID', width: 90 },
  { prop: 'award_title', label: '奖品名称', minWidth: 150 },
  { prop: 'award_subtitle', label: '副标题', minWidth: 130 },
  { prop: 'award_count', label: '总库存', width: 90 },
  { prop: 'award_count_surplus', label: '剩余', width: 90 },
  { prop: 'award_rate', label: '概率', width: 90 },
  { prop: 'rule_models', label: '规则树', width: 140 },
  { prop: 'sort', label: '排序', width: 70 }
]
const awardFields = [
  { prop: 'award_id', label: '奖品ID', type: 'number' },
  { prop: 'award_title', label: '奖品名称' },
  { prop: 'award_subtitle', label: '副标题' },
  { prop: 'award_count', label: '总库存', type: 'number' },
  { prop: 'award_count_surplus', label: '剩余库存', type: 'number' },
  { prop: 'award_rate', label: '中奖概率', type: 'number', precision: 4, step: 0.0001 },
  { prop: 'rule_models', label: '规则树ID' },
  { prop: 'sort', label: '排序', type: 'number' }
]

const ruleColumns = [
  { prop: 'rule_model', label: '规则模型', width: 140 },
  { prop: 'rule_value', label: '规则值', minWidth: 260 },
  { prop: 'rule_desc', label: '说明', minWidth: 160 },
  { prop: 'rule_type', label: '类型', width: 80 }
]
const ruleFields = [
  { prop: 'rule_model', label: '规则模型' },
  { prop: 'rule_value', label: '规则值', type: 'textarea' },
  { prop: 'rule_desc', label: '说明' },
  { prop: 'rule_type', label: '类型（1策略 2奖品）', type: 'number' }
]

async function assemble() {
  try {
    await adminApi.strategyArmory(strategyId)
    ElMessage.success('策略装配完成')
  } catch (e) {
    toastError(e)
  }
}
</script>
