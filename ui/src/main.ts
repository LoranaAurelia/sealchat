import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { i18n, setLocale, setLocaleByNavigatorWithStorage } from './lang'

import App from './App.vue'
import router from './router'

const app = createApp(App)

app.use(i18n)
app.use(createPinia())
app.use(router)

import '@imengyu/vue3-context-menu/lib/vue3-context-menu.css'
import ContextMenu from '@imengyu/vue3-context-menu'

app.use(ContextMenu)

setLocaleByNavigatorWithStorage()

import './assets/main.css'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'
import 'dayjs/locale/ja'

dayjs.locale(document.documentElement.lang);
dayjs.extend(relativeTime)

import { useUserStore } from './stores/user'

router.beforeEach(async (to, from, next) => {
  if (to.name === 'user-signin' || to.name === 'user-signup') {
    return next();
  }

  const user = useUserStore();
  const r = await user.checkUserSession();
  if (r) {
    return next();
  }

  next({ name: 'user-signin' })
  // window.location.href = '//' + window.location.hostname + ":4455/login";
  return;
})

// import AutoImport from 'unplugin-auto-import/vite'
// import { VueHooksPlusResolver } from '@vue-hooks-plus/resolvers'

// export const AutoImportDeps = () =>
//   AutoImport({
//     imports: ['vue', 'vue-router'],
//     include: [/\.[tj]sx?$/, /\.vue$/, /\.vue\?vue/, /\.md$/],
//     dts: 'src/auto-imports.d.ts',
//     resolvers: [VueHooksPlusResolver()],
//   })

// 这几句详见 https://www.naiveui.com/zh-CN/os-theme/docs/style-conflict
const meta = document.createElement('meta')
meta.name = 'naive-ui-style'
document.head.appendChild(meta)

app.mount('#app')
