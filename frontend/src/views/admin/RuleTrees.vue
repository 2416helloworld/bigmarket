<template>
  <div>
    <CrudTable resource="rule-trees" title="规则树" :columns="columns" :fields="fields">
      <template #row-actions="{ row }">
        <el-button link type="primary" size="small" @click="openDetail(row)">节点</el-button>
      </template>
    </CrudTable>

    <el-drawer v-model="drawer" :title="`规则树 ${currentTree} 节点与连线`" size="70%" append-to-body>
      <div v-if="currentTree">
        <CrudTable
          resource="rule-tree-nodes"
          title="规则树节点"
          :columns="nodeColumns"
          :fields="nodeFields"
          :query="nodeQuery"
          :defaults="{ tree_id: currentTree }"
        />
        <div class="mt12">
          <CrudTable
            resource="rule-tree-lines"
            title="规则树连线"
            :columns="lineColumns"
            :fields="lineFields"
            :query="lineQuery"
            :defaults="{ tree_id: currentTree, rule_limit_type: 'EQUAL', rule_limit_value: 'ALLOW' }"
          />
        </div>
      </div>
    </el-drawer>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import CrudTable from './CrudTable.vue'

const drawer = ref(false)
const currentTree = ref('')
const nodeQuery = reactive({ treeId: '' })
const lineQuery = reactive({ treeId: '' })

const columns = [
  { prop: 'tree_id', label: '规则树ID', width: 140 },
  { prop: 'tree_name', label: '名称', width: 160 },
  { prop: 'tree_desc', label: '说明', minWidth: 180 },
  { prop: 'tree_node_rule_key', label: '根节点', width: 140 }
]
const fields = [
  { prop: 'tree_id', label: '规则树ID' },
  { prop: 'tree_name', label: '名称' },
  { prop: 'tree_desc', label: '说明' },
  { prop: 'tree_node_rule_key', label: '根节点' }
]

const nodeColumns = [
  { prop: 'rule_key', label: '节点KEY', width: 150 },
  { prop: 'rule_desc', label: '说明', minWidth: 200 },
  { prop: 'rule_value', label: '规则值', minWidth: 140 }
]
const nodeFields = [
  { prop: 'rule_key', label: '节点KEY' },
  { prop: 'rule_desc', label: '说明' },
  { prop: 'rule_value', label: '规则值' }
]

const lineColumns = [
  { prop: 'rule_node_from', label: '起始节点', width: 150 },
  { prop: 'rule_node_to', label: '目标节点', width: 150 },
  { prop: 'rule_limit_type', label: '限定类型', width: 110 },
  { prop: 'rule_limit_value', label: '限定值', width: 110 }
]
const lineFields = [
  { prop: 'rule_node_from', label: '起始节点' },
  { prop: 'rule_node_to', label: '目标节点' },
  { prop: 'rule_limit_type', label: '限定类型' },
  { prop: 'rule_limit_value', label: '限定值' }
]

function openDetail(row) {
  currentTree.value = row.tree_id
  nodeQuery.treeId = row.tree_id
  lineQuery.treeId = row.tree_id
  drawer.value = true
}
</script>
