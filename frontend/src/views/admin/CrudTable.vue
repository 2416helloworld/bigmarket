<template>
  <div class="admin-card">
    <div class="admin-toolbar">
      <strong>{{ title }}</strong>
      <el-button type="primary" size="small" @click="openCreate">新增</el-button>
      <el-button size="small" @click="load">刷新</el-button>
      <slot name="toolbar" />
      <span style="margin-left: auto; color: #909399; font-size: 12px">共 {{ rows.length }} 条</span>
    </div>
    <el-table :data="rows" size="small" border v-loading="loading" max-height="560">
      <el-table-column v-for="col in columns" :key="col.prop" :prop="col.prop" :label="col.label" :width="col.width" :min-width="col.minWidth || 90" show-overflow-tooltip />
      <el-table-column label="操作" width="170" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" size="small" @click="openEdit(row)">编辑</el-button>
          <slot name="row-actions" :row="row" />
          <el-button link type="danger" size="small" @click="removeRow(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialog" :title="editing ? '编辑' : '新增'" width="520px" append-to-body>
      <el-form label-width="120px">
        <el-form-item v-for="field in fields" :key="field.prop" :label="field.label">
          <el-input-number v-if="field.type === 'number'" v-model="form[field.prop]" :precision="field.precision" :step="field.step || 1" style="width: 100%" />
          <el-select v-else-if="field.type === 'select'" v-model="form[field.prop]" style="width: 100%">
            <el-option v-for="opt in field.options" :key="opt.value" :label="opt.label" :value="opt.value" />
          </el-select>
          <el-date-picker v-else-if="field.type === 'datetime'" v-model="form[field.prop]" type="datetime" value-format="YYYY-MM-DD HH:mm:ss" style="width: 100%" />
          <el-input v-else-if="field.type === 'textarea'" v-model="form[field.prop]" type="textarea" :rows="3" />
          <el-input v-else v-model="form[field.prop]" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialog = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="save">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { adminApi, toastError } from '../../api'

const props = defineProps({
  resource: { type: String, required: true },
  title: { type: String, default: '' },
  columns: { type: Array, default: () => [] },
  fields: { type: Array, default: () => [] },
  query: { type: Object, default: () => ({}) },
  idKey: { type: String, default: 'id' },
  defaults: { type: Object, default: () => ({}) }
})

const rows = ref([])
const loading = ref(false)
const saving = ref(false)
const dialog = ref(false)
const editing = ref(null)
const form = reactive({})

async function load() {
  loading.value = true
  try {
    rows.value = (await adminApi.list(props.resource, props.query)) || []
  } catch (e) {
    toastError(e)
  } finally {
    loading.value = false
  }
}

function resetForm(row) {
  Object.keys(form).forEach((k) => delete form[k])
  const base = { ...props.defaults }
  if (row) {
    props.fields.forEach((f) => {
      base[f.prop] = row[f.prop]
    })
  }
  Object.assign(form, base)
}

function openCreate() {
  editing.value = null
  resetForm(null)
  dialog.value = true
}

function openEdit(row) {
  editing.value = row
  resetForm(row)
  dialog.value = true
}

async function save() {
  saving.value = true
  try {
    if (editing.value) {
      await adminApi.update(props.resource, editing.value[props.idKey], form)
    } else {
      await adminApi.create(props.resource, form)
    }
    ElMessage.success('保存成功')
    dialog.value = false
    await load()
  } catch (e) {
    toastError(e)
  } finally {
    saving.value = false
  }
}

async function removeRow(row) {
  try {
    await ElMessageBox.confirm('确认删除该配置？删除后需重新装配缓存。', '提示', { type: 'warning' })
    await adminApi.remove(props.resource, row[props.idKey])
    ElMessage.success('已删除')
    await load()
  } catch (e) {
    if (e !== 'cancel') toastError(e)
  }
}

watch(() => props.query, () => load(), { deep: true })
onMounted(load)
defineExpose({ load })
</script>
