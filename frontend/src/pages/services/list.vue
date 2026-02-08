<template>
    <el-card shadow="never" class="border-0 mt-0">
        <div class="flex justify-between items-center">
            <!-- 左侧：搜索和刷新 -->
            <div class="flex gap-2 items-center">
                <el-input v-model="searchKeyword" placeholder="搜索" prefix-icon="Search" size="mini" class="w-86"
                    @keyup.enter="handleSearch" />
                <el-button type="primary" size="mini" @click="handleSearch">搜索</el-button>
                <el-button type="default" size="mini" @click="refreshData">刷新</el-button>
            </div>
            <!-- 右侧：新增按钮组 -->
            <div class="flex gap-2">
                <el-button type="primary" size="mini" @click="openAddHttpDrawer">新增Http</el-button>
                <el-button type="primary" size="mini">新增Tcp</el-button>
                <el-button type="primary" size="mini">新增Grpc</el-button>
            </div>
        </div>

        <el-table :data="tableData" height="250" style="width: 100%" class="mt-2" v-loading="loading">
            <el-table-column prop="service_name" label="服务名称" width="160" />
            <el-table-column prop="service_desc" label="服务描述" width="160" show-overflow-tooltip />
            <el-table-column prop="load_type" label="负载类型" width="80" align="center">
                <template #default="scope">
                    <span v-if="scope.row.load_type === 0">HTTP</span>
                    <span v-else-if="scope.row.load_type === 1">TCP</span>
                    <span v-else-if="scope.row.load_type === 2">GRPC</span>
                    <span v-else>{{ scope.row.load_type }}</span>
                </template>
            </el-table-column>
            <el-table-column prop="service_addr" label="服务地址" width="200" align="center" />
            <el-table-column prop="qps" label="QPS" width="80" align="center" />
            <el-table-column prop="qpd" label="QPD" width="80" align="center" />
            <el-table-column prop="total_node" label="节点数量" width="80" align="center" />
            <el-table-column label="操作" width="150" align="center">
                <template #default="scope">
                    <el-button type="primary" size="small">编辑</el-button>
                    <el-button type="danger" size="small">删除</el-button>
                </template>
            </el-table-column>
        </el-table>

        <!-- 分页组件 -->
        <div class="mt-5 flex justify-center">
            <el-pagination v-model:current-page="currentPage" v-model:page-size="pageSize" :page-sizes="[10, 20]"
                layout="total, prev,pager,next,sizes" :total="total" @size-change="handleSizeChange"
                @current-change="handleCurrentChange" />
        </div>

        <!-- 新增服务抽屉 -->
        <FormDrawer ref="formDrawerRef" title="新增Http服务" :size="'60%'" @dad_submit="handleSubmit" >
            <el-form ref="formRef" :model="formData" :rules="rules" label-width="120px">
                <el-collapse v-model="activeNames">
                    <!-- 基本信息 -->
                    <el-collapse-item title="基本信息" name="basic">
                        <el-form-item label="服务名称" prop="service_name">
                            <el-input v-model="formData.service_name" placeholder="请输入服务名称" />
                        </el-form-item>
                        <el-form-item label="服务描述" prop="service_desc">
                            <el-input v-model="formData.service_desc" type="textarea" :rows="2" placeholder="请输入服务描述" />
                        </el-form-item>
                    </el-collapse-item>

                    <!-- 接入配置 -->
                    <el-collapse-item title="接入配置" name="access">
                        <el-form-item label="接入类型" prop="rule_type">
                            <el-radio-group v-model="formData.rule_type">
                                <el-radio :label="0">域名</el-radio>
                                <el-radio :label="1">前缀</el-radio>
                            </el-radio-group>
                        </el-form-item>
                        <el-form-item label="接入路径" prop="rule">
                            <el-input v-model="formData.rule" placeholder="请输入域名或前缀" />
                        </el-form-item>

                        <el-form-item label="支持HTTPS" prop="need_https" class="w-80">
                            <el-switch v-model="formData.need_https" :active-value="1" :inactive-value="0" />
                        </el-form-item>

                        <el-form-item label="启用strip_uri" prop="need_strip_uri">
                            <el-switch v-model="formData.need_strip_uri" :active-value="1" :inactive-value="0" />
                        </el-form-item>

                        <el-form-item label="支持WebSocket" prop="need_websocket">
                            <el-switch v-model="formData.need_websocket" :active-value="1" :inactive-value="0" />
                        </el-form-item>
                        <el-form-item label="URL重写" prop="url_rewrite">
                            <el-input v-model="formData.url_rewrite" placeholder="请输入URL重写规则" />
                        </el-form-item>
                        <el-form-item label="Header转换" prop="header_transfor">
                            <el-input v-model="formData.header_transfor" type="textarea" :rows="2"
                                placeholder="请输入Header转换规则" />
                        </el-form-item>
                    </el-collapse-item>

                    <!-- 权限控制 -->
                    <el-collapse-item title="权限控制" name="auth">
                        <el-form-item label="开启权限" prop="open_auth">
                            <el-switch v-model="formData.open_auth" :active-value="1" :inactive-value="0" />
                        </el-form-item>
                        <el-form-item label="黑名单IP" prop="black_list">
                            <el-input v-model="formData.black_list" type="textarea" :rows="2"
                                placeholder="请输入黑名单IP，多个IP用逗号分隔" />
                        </el-form-item>
                        <el-form-item label="白名单IP" prop="white_list">
                            <el-input v-model="formData.white_list" type="textarea" :rows="2"
                                placeholder="请输入白名单IP，多个IP用逗号分隔" />
                        </el-form-item>
                        <el-form-item label="客户端IP限流" prop="clientip_flow_limit">
                            <el-input-number v-model="formData.clientip_flow_limit" :min="0"
                                controls-position="right" />
                        </el-form-item>
                        <el-form-item label="服务端限流" prop="service_flow_limit">
                            <el-input-number v-model="formData.service_flow_limit" :min="0" controls-position="right" />
                        </el-form-item>
                    </el-collapse-item>

                    <!-- 负载均衡 -->
                    <el-collapse-item title="负载均衡" name="loadbalance">
                        <el-form-item label="轮询方式" prop="round_type">
                            <el-select v-model="formData.round_type" placeholder="请选择轮询方式">
                                <el-option label="随机" :value="0" />
                                <el-option label="轮询" :value="1" />
                                <el-option label="加权轮询" :value="2" />
                                <el-option label="IP哈希" :value="3" />
                            </el-select>
                        </el-form-item>
                        <el-form-item label="IP列表" prop="ip_list">
                            <el-input v-model="formData.ip_list" type="textarea" :rows="3"
                                placeholder="请输入IP列表，格式：127.0.0.1:8001,127.0.0.1:8002" />
                        </el-form-item>
                        <el-form-item label="权重列表" prop="weight_list">
                            <el-input v-model="formData.weight_list" placeholder="请输入权重列表，格式：1,2,3" />
                        </el-form-item>
                        <el-form-item label="连接超时(秒)" prop="upstream_connect_timeout">
                            <el-input-number v-model="formData.upstream_connect_timeout" :min="0"
                                controls-position="right" />
                        </el-form-item>
                        <el-form-item label="Header超时(秒)" prop="upstream_header_timeout">
                            <el-input-number v-model="formData.upstream_header_timeout" :min="0"
                                controls-position="right" />
                        </el-form-item>
                        <el-form-item label="空闲超时(秒)" prop="upstream_idle_timeout">
                            <el-input-number v-model="formData.upstream_idle_timeout" :min="0"
                                controls-position="right" />
                        </el-form-item>
                        <el-form-item label="最大空闲连接数" prop="upstream_max_idle">
                            <el-input-number v-model="formData.upstream_max_idle" :min="0" controls-position="right" />
                        </el-form-item>
                    </el-collapse-item>
                </el-collapse>
            </el-form>
        </FormDrawer>
    </el-card>
</template>

<script setup>
import { ref } from 'vue'
import { getServiceList, addService } from '~/api/services.js'
import { ElMessage } from 'element-plus'
import FormDrawer from '~/components/FormDrawer.vue'

const searchKeyword = ref('')
const tableData = ref([])
const loading = ref(false)
//分页
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

//新增表单部分
const formDrawerRef = ref(null)
const formRef = ref(null)
const activeNames = ref(['basic', 'access', 'auth', 'loadbalance'])

const formData = ref({
    service_name: '',
    service_desc: '',
    rule_type: 0,
    rule: '',
    need_https: 0,
    need_strip_uri: 0,
    need_websocket: 0,
    url_rewrite: '',
    header_transfor: '',
    open_auth: 0,
    black_list: '',
    white_list: '',
    clientip_flow_limit: 0,
    service_flow_limit: 0,
    round_type: 0,
    ip_list: '',
    weight_list: '',
    upstream_connect_timeout: 0,
    upstream_header_timeout: 0,
    upstream_idle_timeout: 0,
    upstream_max_idle: 0
})

const rules = {
    service_name: [
        { required: true, message: '请输入服务名称', trigger: 'blur' },
        { min: 1, max: 50, message: '长度在 1 到 50 个字符', trigger: 'blur' }
    ],
    service_desc: [
        { required: true, message: '请输入服务描述', trigger: 'blur' },
        { min: 1, max: 255, message: '长度在 1 到 255 个字符', trigger: 'blur' }
    ],
    rule_type: [
        { required: true, message: '请选择接入类型', trigger: 'change' }
    ],
    rule: [
        { required: true, message: '请输入接入路径', trigger: 'blur' }
    ],
    round_type: [
        { required: true, message: '请选择轮询方式', trigger: 'change' }
    ],
    ip_list: [
        { required: true, message: '请输入IP列表', trigger: 'blur' }
    ],
    weight_list: [
        { required: true, message: '请输入权重列表', trigger: 'blur' }
    ]
}

const handleSubmit = () => {
    formRef.value.validate((valid) => {
        if (valid) {
            console.log('表单数据:', formData.value)
            ElMessage.success('提交成功')
            formDrawerRef.value.close()
            getTableData(1)
        } else {
            ElMessage.error('请填写完整表单')
        }
    })
}

const openAddHttpDrawer = () => {
    formDrawerRef.value.open()
    formData.value = {
        service_name: '',
        service_desc: '',
        rule_type: 0,
        rule: '',
        need_https: 0,
        need_strip_uri: 0,
        need_websocket: 0,
        url_rewrite: '',
        header_transfor: '',
        open_auth: 0,
        black_list: '',
        white_list: '',
        clientip_flow_limit: 0,
        service_flow_limit: 0,
        round_type: 0,
        ip_list: '',
        weight_list: '',
        upstream_connect_timeout: 0,
        upstream_header_timeout: 0,
        upstream_idle_timeout: 0,
        upstream_max_idle: 0
    }

}


// 获取服务列表分页数据
function getTableData(p = null) {
    if (typeof p == "number") {
        currentPage.value = p
    }

    // 防止重复请求
    if (loading.value) return

    loading.value = true

    // 记录开始时间
    const startTime = Date.now()
    const minLoadingTime = 200 // 最小加载时间（毫秒）

    getServiceList(currentPage.value, pageSize.value)
        .then(res => {
            if (res.data.success === true) {
                tableData.value = res.data.data.list
                total.value = res.data.data.total
                currentPage.value = res.data.data.page
                pageSize.value = res.data.data.size
            } else {
                tableData.value = []
                ElMessage.error("获取服务列表失败")
            }
        }).finally(() => {
            // 计算执行时间，确保至少显示指定的加载时间
            const elapsed = Date.now() - startTime
            if (elapsed < minLoadingTime) {
                // 如果执行时间不足，延迟关闭加载状态
                setTimeout(() => {
                    loading.value = false
                }, minLoadingTime - elapsed)
            } else {
                // 执行时间足够，直接关闭加载状态
                loading.value = false
            }
        })
}

getTableData()

const handleSearch = () => {
    console.log('搜索关键词:', searchKeyword.value)
}

const refreshData = () => {
    console.log('刷新数据')
    getTableData(1)
}


// function handleCurrentChange(val) {
//     getTableData(val)
// }

// 分页方法
function handleSizeChange(val) {
    pageSize.value = val
    getTableData(1) // 切换每页条数时重置到第一页
}

const handleCurrentChange = (val) => {
    getTableData(val)
}
</script>

<style scoped>
h1 {
    @apply text-2xl font-bold text-gray-800;
}

:deep(.el-card__body) {
    padding-top: 0px;
}
</style>
