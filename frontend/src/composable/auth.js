import { useCookies } from "@vueuse/integrations/useCookies";
const TokenKey = 'token'
const cookie = useCookies(TokenKey)

//获取token
export function getToken() {
    return cookie.get(TokenKey)
}

//设置token
export function setToken(token) {
    cookie.set(TokenKey, token)
}

//删除token
export function removeToken() {
    cookie.remove(TokenKey)
}