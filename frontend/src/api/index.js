import axios from 'axios'
import { ElMessage } from 'element-plus'

const http = axios.create({ baseURL: '/api/v1', timeout: 20000 })

http.interceptors.request.use((config) => {
  const token = localStorage.getItem('bm_admin_token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

http.interceptors.response.use(
  (res) => {
    const body = res.data
    if (body && body.code && body.code !== '0000') {
      return Promise.reject(new Error(body.info || '请求失败'))
    }
    return body ? body.data : null
  },
  (err) => {
    if (err.response && err.response.status === 401) {
      localStorage.removeItem('bm_admin_token')
      if (!location.pathname.startsWith('/admin/login')) location.href = '/admin/login'
    }
    const msg = err.response?.data?.info || err.message || '网络异常'
    return Promise.reject(new Error(msg))
  }
)

function request(method, url, data, params) {
  return http({ method, url, data, params })
}

export const userApi = {
  activityList: () => request('get', '/raffle/activity/list'),
  armory: (activityId) => request('post', '/raffle/activity/armory', { activityId }),
  awards: (activityId) => request('get', '/raffle/strategy/query_raffle_award_list', null, { activityId }),
  draw: (userId, activityId) => request('post', '/raffle/activity/draw', { userId, activityId }),
  account: (userId, activityId) => request('post', '/raffle/activity/query_user_activity_account', { userId, activityId }),
  credit: (userId) => request('post', '/raffle/activity/query_user_credit_account', { userId }),
  skus: (activityId) => request('get', '/raffle/activity/query_sku_product_list_by_activity_id', null, { activityId }),
  exchange: (userId, sku) => request('post', '/raffle/activity/credit_pay_exchange_sku', { userId, sku }),
  sign: (userId) => request('post', '/raffle/activity/calendar_sign_rebate', { userId }),
  isSigned: (userId) => request('get', '/raffle/activity/is_calendar_sign_rebate', null, { userId }),
  awardRecords: (userId) => request('get', '/raffle/activity/query_user_award_records', null, { userId }),
  creditOrders: (userId) => request('get', '/raffle/activity/query_user_credit_orders', null, { userId }),
  rebateOrders: (userId) => request('get', '/raffle/activity/query_user_rebate_orders', null, { userId })
}

export const adminApi = {
  login: (username, password) => request('post', '/admin/login', { username, password }),
  dashboard: () => request('get', '/admin/dashboard'),
  list: (resource, params) => request('get', `/admin/crud/${resource}`, null, params),
  create: (resource, data) => request('post', `/admin/crud/${resource}`, data),
  update: (resource, id, data) => request('put', `/admin/crud/${resource}/${id}`, data),
  remove: (resource, id) => request('delete', `/admin/crud/${resource}/${id}`),
  orders: (type, userId) => request('get', '/admin/orders', null, { type, userId }),
  users: () => request('get', '/admin/users'),
  creditAdjust: (userId, amount, tradeType) => request('post', '/admin/users/credit', { userId, amount, tradeType }),
  armory: (activityId) => request('post', '/admin/armory', { activityId }),
  tasks: () => request('get', '/admin/tasks'),
  ruleTreeDetail: (treeId) => request('get', `/admin/rule-trees/${treeId}/detail`),
  strategyArmory: (strategyId) => request('post', '/raffle/strategy/strategy_armory', { strategyId })
}

export function toastError(err) {
  ElMessage.error(err?.message || '操作失败')
}

export default http
