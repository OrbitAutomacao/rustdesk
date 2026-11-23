<template>
  <div class="login-container">
    <div class="login-hero">
      <img src="@/assets/orbit-logo.webp" alt="Orbit Automação" class="hero-logo"/>
      <div class="hero-content">
        <p class="hero-kicker">ORBIT AUTOMAÇÃO</p>
        <h1>Suporte remoto,<br>em um só lugar.</h1>
        <p>Gerencie dispositivos, organize sua equipe e encontre rapidamente os computadores dos clientes.</p>
      </div>
      <span class="hero-footer">Orbit Suporte Remoto</span>
    </div>
    <div class="login-card">
      <div class="login-card-head">
        <img src="@/assets/orbit-logo.webp" alt="Orbit Automação" class="login-logo"/>
        <h2>Bem-vindo de volta</h2>
        <p>Entre para acessar o painel de suporte.</p>
      </div>

      <el-form v-if="!disablePwd" label-position="top" class="login-form">
        <el-form-item :label="T('Username')">
          <el-input v-model="form.username" type="username" class="login-input"></el-input>
        </el-form-item>

        <el-form-item :label="T('Password')">
          <el-input v-model="form.password" type="password" @keyup.enter.native="login" show-password
                    class="login-input"></el-input>
        </el-form-item>
        <el-form-item :label="T('Captcha')" v-if="captchaCode">
          <el-input v-model="form.captcha" @keyup.enter.native="login"  class="login-input captcha-input">
            <template #append>
              <img :src="captchaCode.b64" @click="loadCaptcha" class="captcha" alt="captcha"/>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-button @click="login" type="primary" class="login-button">{{ T('Login') }}</el-button>
          <el-button v-if="allowRegister" @click="register" class="login-button">{{ T('Register') }}</el-button>
        </el-form-item>
      </el-form>

      <div class="divider" v-if="options.length > 0 && !disablePwd">
        <span>{{ T('or login in with') }}</span>
      </div>

      <div class="oidc-options">
        <div v-for="(option, index) in options" :key="index" class="oidc-option">
          <el-button @click="handleOIDCLogin(option.name)" class="oidc-btn">
            <img :src="getProviderImage(option.name)" alt="provider" class="oidc-icon"/>
            <span>{{ T(option.name) }}</span>
          </el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
  import { reactive, onMounted, ref } from 'vue'
  import { useUserStore } from '@/store/user'
  import { ElMessage } from 'element-plus'
  import { T } from '@/utils/i18n'
  import { useRoute, useRouter } from 'vue-router'
  import { loginOptions, captcha } from '@/api/login'
  import { getCode, removeCode } from '@/utils/auth'

  const oauthInfo = ref({})
  const userStore = useUserStore()
  const route = useRoute()
  const router = useRouter()
  const options = reactive([]) // 存储 OIDC 登录选项

  let platform = window.navigator.platform
  if (navigator.platform.indexOf('Mac') === 0) {
    platform = 'mac'
  } else if (navigator.platform.indexOf('Win') === 0) {
    platform = 'windows'
  } else if (navigator.platform.indexOf('Linux armv') === 0) {
    platform = 'android'
  } else if (navigator.platform.indexOf('Linux') === 0) {
    platform = 'linux'
  }
  const userAgent = navigator.userAgent
  let browser = 'Unknown Browser'
  if (/chrome|crios/i.test(userAgent)) browser = 'Chrome'
  else if (/firefox|fxios/i.test(userAgent)) browser = 'Firefox'
  else if (/safari/i.test(userAgent) && !/chrome/i.test(userAgent)) browser = 'Safari'
  else if (/edg/i.test(userAgent)) browser = 'Edge'

  const form = reactive({
    username: '',
    password: '',
    platform: platform,
    captcha: '',
    captcha_id: ''
  })

  const captchaCode = ref('')
  const redirect = route.query?.redirect
  const login = async () => {
    const res = await userStore.login(form).catch(e => e)
    if (!res.code) {
      ElMessage.success(T('LoginSuccess'))
      router.push({ path: redirect || '/', replace: true })
      return
    }
    if (res.code === 110) {
      // need captcha
      loadCaptcha()
    }
  }

  const loadCaptcha = async () => {
    const captchaRes = await captcha().catch(_ => false)
    console.log(captchaRes)
    captchaCode.value = captchaRes.data.captcha
    form.captcha_id = captchaRes.data.captcha.id
  }

  const handleOIDCLogin = (provider) => {
    userStore.oidc(provider, platform, browser)
  }

  import googleImage from '@/assets/google.png'
  import githubImage from '@/assets/github.png'
  import oidcImage from '@/assets/oidc.png'
  import webauthImage from '@/assets/webauth.png'
  import defaultImage from '@/assets/oidc.png'

  const providerImageMap = {
    google: googleImage,
    github: githubImage,
    oidc: oidcImage,
    // WebAuth: webauthImage,
    default: defaultImage,
  }

  const getProviderImage = (provider) => {
    return providerImageMap[provider.toLowerCase()] || providerImageMap.default
  }

  const allowRegister = ref(false)
  const disablePwd = ref(false)
  const loadLoginOptions = async () => {
    try {
      const res = await loginOptions().catch(_ => false)
      if (!res || !res.data) return console.error('No valid response received')
      res.data.ops.map(option => (options.push({ name: option }))) // 创建新的对象数组
      if (res.data.auto_oidc) {
        // 如果有自动OIDC登录选项，直接调用第一个
        handleOIDCLogin(res.data.ops[0])
      }
      disablePwd.value = res.data.disable_pwd
      allowRegister.value = res.data.register
      if (res.data.need_captcha) {
        loadCaptcha()
      }
    } catch (error) {
      console.error('Error loading login options:', error.message)
    }
  }

  onMounted(async () => {
    const code = getCode()
    if (code) {
      // 如果code存在，进行query获取user info
      const res = await userStore.query(code)
      if (res) {
        // 删除code，确保跳转之前对code进行清楚
        removeCode()
        ElMessage.success(T('LoginSuccess'))
        router.push({ path: redirect || '/', replace: true })
      }
    } else {
      // 如果code不存在, 现实登陆页面
      loadLoginOptions() // 组件挂载后调用登录选项加载函数
    }
  })

  const register = () => {
    router.push('/register')
  }
</script>

<style scoped lang="scss">
.login-container {
  display: grid;
  grid-template-columns: minmax(320px, 1fr) minmax(420px, 1fr);
  min-height: 100vh;
  background: #f4f8fb;
}
.login-hero {
  position: relative;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  overflow: hidden;
  padding: clamp(36px, 6vw, 84px);
  background: radial-gradient(circle at 85% 18%, #0b5378 0, transparent 36%), linear-gradient(150deg, #071b33, #0b2d4b 57%, #0a415c);
  color: #fff;
}
.login-hero::after {
  content: "";
  position: absolute;
  width: 420px;
  height: 420px;
  right: -170px;
  bottom: -190px;
  border: 1px solid rgba(80, 208, 241, .25);
  border-radius: 50%;
  box-shadow: 0 0 0 65px rgba(49, 180, 222, .04), 0 0 0 130px rgba(49, 180, 222, .03);
}
.hero-logo { filter: brightness(0) invert(1); width: 172px; height: auto; position: relative; z-index: 1; }
.hero-content { max-width: 540px; position: relative; z-index: 1; }
.hero-kicker { margin: 0 0 22px; color: #58d1f0; font-size: 12px; font-weight: 800; letter-spacing: .22em; }
.hero-content h1 { margin: 0 0 22px; font-size: clamp(38px, 4.5vw, 66px); line-height: 1.07; letter-spacing: -.045em; }
.hero-content > p:last-child { max-width: 490px; margin: 0; color: #c3d9e6; font-size: 17px; line-height: 1.65; }
.hero-footer { color: #9db9c9; font-size: 13px; position: relative; z-index: 1; }
.login-card {
  width: min(440px, calc(100% - 48px));
  margin: auto;
  padding: 38px;
  border: 1px solid #e0eaf0;
  border-radius: 22px;
  background: #fff;
  box-shadow: 0 20px 60px rgba(24, 56, 80, .08);
}
.login-card-head { margin-bottom: 32px; }
.login-logo { display: block; width: 150px; height: auto; margin-bottom: 28px; }
.login-card-head h2 { margin: 0 0 8px; color: #15334b; font-size: 28px; letter-spacing: -.025em; }
.login-card-head p { margin: 0; color: #668095; font-size: 14px; }
.login-form { margin: 0; }
.login-input { width: 100%; }
.login-button { width: 100%; min-height: 44px; margin: 12px 0 0; }
.captcha { cursor: pointer; width: 150px; }
.captcha-input :deep(.el-input-group__append) { padding: 0; overflow: hidden; border-radius: 0 9px 9px 0; }
:deep(.el-form-item__label) { color: #344e64; font-weight: 650; }
:deep(.el-input__wrapper) { min-height: 44px; background: #fbfdfe; }
:deep(.el-input__inner) { color: #16334b; }
.divider { display: flex; align-items: center; gap: 12px; margin: 26px 0 18px; color: #7c91a1; font-size: 12px; }
.divider::before, .divider::after { content: ""; height: 1px; background: #dde7ee; flex: 1; }
.oidc-options { display: flex; flex-direction: column; gap: 10px; }
.oidc-btn { display: flex; align-items: center; justify-content: center; gap: 10px; width: 100%; min-height: 44px; }
.oidc-icon { width: 22px; height: 22px; }
@media (max-width: 900px) {
  .login-container { grid-template-columns: 1fr; }
  .login-hero { min-height: 200px; padding: 28px; }
  .hero-logo { width: 115px; }
  .hero-content h1 { font-size: 26px; margin: 0; }
  .hero-kicker, .hero-content > p:last-child, .hero-footer { display: none; }
  .login-card { margin: -20px auto 32px; position: relative; z-index: 1; }
}
@media (max-width: 520px) {
  .login-card { padding: 26px; width: calc(100% - 28px); }
  .login-card-head h2 { font-size: 24px; }
}
</style>