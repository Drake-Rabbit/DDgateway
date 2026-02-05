<template>
    <el-drawer v-model="showDrawer" :title="title" :size="size" 
    :close-on-click-modal="false"
    :destroy-on-close="destroyOnClose"
    >
        <div class="formDrawer">
                <slot></slot>
        </div>
        <div class="formDrawerFooter">
            <el-button type="primary" @click="son_submit">{{confirmText}}</el-button>
            <el-button type="default" @click="close">取消</el-button>
        </div>


    </el-drawer>
</template>

<script setup>
import { ref } from 'vue'

const showDrawer = ref(false)

//打开
const open = () => { showDrawer.value = true }
//关闭
const close = () => { showDrawer.value = false }
//定义emit
const emit = defineEmits(['dad_submit'])
//提交
const son_submit = () => {emit('dad_submit')}




//定义props
const props = defineProps({
    title: String,
    size: {
        type: String,
        default: '45%'
    },
    destroyOnClose: {
        type: Boolean,
        default: true
    },
    confirmText: {
        type: String,
        default: '提交'
    }

})



//向父组件暴露方法
defineExpose({
    open,
    close
})

</script>


<style scoped></style>
