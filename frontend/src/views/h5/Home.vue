<template>
  <div class="h5-page">
    <div class="h5-header">
      <div class="h5-brand">
        <div class="h5-logo">抽</div>
        <div>
          <div class="h5-title">大营销抽奖</div>
          <div class="h5-sub">Spring To Go · 全栈重构版</div>
        </div>
      </div>
      <div class="h5-user" @click="openUserDialog">👤 {{ userId || '设置用户' }}</div>
    </div>

    <div class="h5-card">
      <el-select v-model="activityId" size="large" style="width: 100%" @change="loadActivity">
        <el-option v-for="item in activities" :key="item.activityId" :label="item.activityName" :value="item.activityId" />
      </el-select>
      <div class="h5-quota" style="margin-top: 14px">
        <div class="h5-quota-item">
          <div class="h5-quota-num">{{ account.totalCountSurplus || 0 }}</div>
          <div class="h5-quota-label">总剩余次数</div>
        </div>
        <div class="h5-quota-item">
          <div class="h5-quota-num">{{ account.dayCountSurplus || 0 }}</div>
          <div class="h5-quota-label">今日剩余</div>
        </div>
        <div class="h5-quota-item">
          <div class="h5-quota-num">{{ credit.availableAmount ?? 0 }}</div>
          <div class="h5-quota-label">可用积分</div>
        </div>
      </div>
      <div class="h5-row">
        <button class="h5-btn-ghost" :disabled="signed" @click="signIn">{{ signed ? '今日已签到' : '每日签到领积分' }}</button>
        <button class="h5-btn-ghost" @click="openUserDialog">切换用户</button>
      </div>
    </div>

    <div class="h5-section-title">
      <span>奖品池</span>
      <small>共 {{ awards.length }} 个奖品</small>
    </div>
    <div class="h5-prizes">
      <div v-for="(item, index) in awards" :key="item.awardId" class="h5-prize" :class="{ active: activeIndex === index }">
        <div class="h5-prize-icon">{{ prizeIcon(item.awardId) }}</div>
        <div class="h5-prize-name">{{ item.awardTitle }}</div>
        <div class="h5-prize-rate">{{ (item.awardRate * 100).toFixed(2) }}%</div>
      </div>
    </div>
    <button class="h5-draw" :disabled="drawing || !awards.length" @click="handleDraw">
      {{ drawing ? '抽奖中...' : '立即抽奖' }}
    </button>

    <div class="h5-section-title">
      <span>积分兑换次数</span>
      <small>兑换后获得抽奖额度</small>
    </div>
    <div v-for="sku in skus" :key="sku.sku" class="h5-sku">
      <div class="h5-sku-info">
        <div class="h5-sku-name">SKU {{ sku.sku }} · {{ sku.activityName }}</div>
        <div class="h5-sku-desc">可得 {{ sku.totalCount }} 次（日 {{ sku.dayCount }} 次）｜库存 {{ sku.stockCountSurplus }}</div>
      </div>
      <div>
        <span class="h5-sku-price">{{ sku.productAmount }} 积分</span>
        <button class="h5-sku-btn" :disabled="exchanging" @click="exchange(sku)">兑换</button>
      </div>
    </div>

    <div class="h5-section-title"><span>我的记录</span></div>
    <div class="h5-card">
      <el-tabs v-model="activeTab" stretch>
        <el-tab-pane label="中奖" name="award">
          <div class="h5-record-list">
            <div v-for="r in records.awards" :key="r.orderId" class="h5-record-item">
              <span>{{ r.awardTitle }}</span>
              <span><small>{{ formatTime(r.createTime) }}</small> {{ r.awardState === 'completed' ? '已发放' : '发放中' }}</span>
            </div>
            <div v-if="!records.awards.length" class="h5-empty">暂无中奖记录，先去抽一次吧</div>
          </div>
        </el-tab-pane>
        <el-tab-pane label="积分" name="credit">
          <div class="h5-record-list">
            <div v-for="r in records.credits" :key="r.orderId" class="h5-record-item">
              <span>{{ r.tradeName }} · {{ r.tradeType === 'forward' ? '增加' : '扣减' }}</span>
              <span>{{ r.tradeAmount }} <small>{{ formatTime(r.createTime) }}</small></span>
            </div>
            <div v-if="!records.credits.length" class="h5-empty">暂无积分记录</div>
          </div>
        </el-tab-pane>
        <el-tab-pane label="返利" name="rebate">
          <div class="h5-record-list">
            <div v-for="r in records.rebates" :key="r.orderId" class="h5-record-item">
              <span>{{ r.rebateDesc }}</span>
              <span>{{ r.rebateType === 'sku' ? '抽奖额度' : r.rebateAmount + ' 积分' }}</span>
            </div>
            <div v-if="!records.rebates.length" class="h5-empty">暂无返利记录</div>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>

    <a class="h5-admin-link" href="/admin">进入运营后台 →</a>

    <el-dialog v-model="userDialog" title="设置当前用户" width="88%" append-to-body>
      <el-input v-model="tempUserId" placeholder="请输入用户 ID，例如 demo / user001" clearable />
      <div style="margin-top: 8px; font-size: 12px; color: #909399">user001 在黑名单中，只会中兜底积分；demo 为普通用户。</div>
      <template #footer>
        <el-button @click="userDialog = false">取消</el-button>
        <el-button type="primary" @click="saveUser">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElNotification } from 'element-plus'
import { userApi, toastError } from '../../api'

const userId = ref(localStorage.getItem('bm_user_id') || 'demo')
const tempUserId = ref(userId.value)
const userDialog = ref(!userId.value)
const activities = ref([])
const activityId = ref(null)
const awards = ref([])
const skus = ref([])
const account = reactive({})
const credit = reactive({ availableAmount: 0 })
const records = reactive({ awards: [], credits: [], rebates: [] })
const activeTab = ref('award')
const drawing = ref(false)
const activeIndex = ref(-1)
const signed = ref(false)
const exchanging = ref(false)

const icons = ['🎁', '🧧', '🎟️', '💎', '📱', '🚗', '💡', '🪑', '🕹️', '👾']
function prizeIcon(id) {
  return icons[(id || 0) % icons.length]
}
function formatTime(t) {
  if (!t) return ''
  return String(t).replace('T', ' ').slice(5, 16)
}
function saveUser() {
  if (!tempUserId.value.trim()) {
    ElMessage.warning('请输入用户 ID')
    return
  }
  userId.value = tempUserId.value.trim()
  localStorage.setItem('bm_user_id', userId.value)
  userDialog.value = false
  loadActivity()
}
function openUserDialog() {
  tempUserId.value = userId.value
  userDialog.value = true
}

async function init() {
  if (!userId.value) return
  try {
    activities.value = await userApi.activityList()
    if (!activities.value.length) return
    activityId.value = activities.value[0].activityId
    await loadActivity()
  } catch (e) {
    toastError(e)
  }
}

async function loadActivity() {
  if (!userId.value || !activityId.value) return
  activeIndex.value = -1
  try {
    await userApi.armory(activityId.value)
    awards.value = await userApi.awards(activityId.value)
    Object.assign(account, await userApi.account(userId.value, activityId.value))
    Object.assign(credit, await userApi.credit(userId.value))
    skus.value = await userApi.skus(activityId.value)
    signed.value = await userApi.isSigned(userId.value)
    await loadRecords()
  } catch (e) {
    toastError(e)
  }
}

async function loadRecords() {
  try {
    const [awardsList, creditsList, rebatesList] = await Promise.all([
      userApi.awardRecords(userId.value),
      userApi.creditOrders(userId.value),
      userApi.rebateOrders(userId.value)
    ])
    records.awards = awardsList || []
    records.credits = creditsList || []
    records.rebates = rebatesList || []
  } catch (e) {
    // 记录加载失败不阻断主流程
  }
}

async function signIn() {
  try {
    await userApi.sign(userId.value)
    signed.value = true
    ElMessage.success('签到成功，积分与抽奖额度已发放')
    await loadActivity()
  } catch (e) {
    toastError(e)
  }
}

async function handleDraw() {
  if (!userId.value) {
    userDialog.value = true
    return
  }
  drawing.value = true
  let index = 0
  const timer = setInterval(() => {
    activeIndex.value = index % Math.max(awards.value.length, 1)
    index += 1
  }, 90)
  try {
    const result = await userApi.draw(userId.value, activityId.value)
    await new Promise((r) => setTimeout(r, 500))
    clearInterval(timer)
    const idx = awards.value.findIndex((a) => a.awardId === result.awardId)
    activeIndex.value = idx
    ElNotification({
      title: '恭喜中奖',
      message: `${result.awardTitle}（奖品ID ${result.awardId}）`,
      type: 'success',
      duration: 4000
    })
    await loadActivity()
  } catch (e) {
    clearInterval(timer)
    activeIndex.value = -1
    toastError(e)
  } finally {
    drawing.value = false
  }
}

async function exchange(sku) {
  exchanging.value = true
  try {
    await userApi.exchange(userId.value, sku.sku)
    ElMessage.success('兑换成功，抽奖次数已到账')
    await loadActivity()
  } catch (e) {
    toastError(e)
  } finally {
    exchanging.value = false
  }
}

onMounted(init)
</script>
