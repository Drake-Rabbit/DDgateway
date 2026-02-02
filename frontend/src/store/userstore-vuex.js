import { createStore } from 'vuex'

const userstore = createStore({
    state() {
        return {
            userinfo: {}
        }
    },
    mutations: {
        // 设置用户信息
        setUserInfo(state, userinfo) {
            state.userinfo = userinfo
        }
    }
})

export default userstore
