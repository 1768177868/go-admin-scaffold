import Cookies from 'js-cookie'

const TokenKey = 'go-admin-token'

export function getToken() {
  return Cookies.get(TokenKey)
}

export function setToken(token) {
  return Cookies.set(TokenKey, token, { expires: 7 }) // 7天过期
}

export function removeToken() {
  return Cookies.remove(TokenKey)
} 