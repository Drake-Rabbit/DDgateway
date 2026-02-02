<script setup>
import {useUserStore} from '~/store/user.js'
const userstore = useUserStore()


</script>

<template>
 dashborad
 <div>
    <h1>{{userstore.userinfo}},欢迎来到仪表盘</h1>
    
 </div>
</template>

<style scoped>


</style>