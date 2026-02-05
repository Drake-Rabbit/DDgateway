<template>
    <div>
        <h1>服务列表</h1>
    </div>
</template>

<script setup>
import {useUserStore} from '~/store/user.js'
const userstore = useUserStore()
</script>

<style scoped>
    h1 {
        @apply text-2xl font-bold text-gray-800;
    }
</style>
