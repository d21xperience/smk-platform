// Configuration for your app
// https://v2.quasar.dev/quasar-cli-vite/quasar-config-file
// const path = require('path')
import { defineConfig } from '#q-app'

export default defineConfig((/* ctx */) => {
  return {
    boot: [
      'axios', // Inisialisasi Axios dan interceptor
    ],

    css: ['app.scss'],

    extras: [
      'bootstrap-icons',
      'roboto-font', // optional, you are not bound to it
      'material-icons', // optional, you are not bound to it
    ],

    build: {
      target: {
        // browser: 'baseline-widely-available',
        // node: 'node22'
      },

      filenameBasedRouting: true,

      vueRouterMode: 'hash', // available values: 'hash', 'history'

      defineEnv: {
        QCLI_API_BASE_URL: JSON.stringify(process.env.QCLI_API_BASE_URL || '/api/v1'),
      },

      extendViteConf(viteConf) {
        viteConf.server = viteConf.server || {}
        viteConf.server.fs = viteConf.server.fs || {}
        viteConf.server.fs.allow = ['../../']
      },

      vitePlugins: [
        [
          'vite-plugin-checker',
          {
            eslint: {
              lintCommand: 'eslint -c ./eslint.config.js "./src*/**/*.{js,mjs,cjs,vue}"',
              useFlatConfig: true,
            },
          },
          { server: false },
        ],
      ],
    },

    // https://v2.quasar.dev/quasar-cli-vite/quasar-config-file#devserver
    devServer: {
      port: 9002, // <-- FIX: Mengubah port default Quasar (9000) agar tidak bentrok dengan MinIO
      open: true, // opens browser window automatically
      proxy: {
        '/api': {
          target: 'http://localhost:8081',
          changeOrigin: true,
          secure: false,
        },
      },
    },

    framework: {
      config: {},
      plugins: ['Notify', 'Dialog', 'Meta'],
    },

    animations: [],

    ssr: {
      prodPort: 3000,
      middlewares: [
        'render', // keep this as last one
      ],
      pwa: false,
    },

    pwa: {
      workboxMode: 'GenerateSW', 
    },

    cordova: {},

    capacitor: {
      hideSplashscreen: true,
    },

    electron: {
      preloadScripts: ['electron-preload'],
      inspectPort: 5858,
      bundler: 'packager',
      packager: {},
      builder: {
        appId: 'siakad-tu',
      },
    },

    bex: {
      extraScripts: [],
    },
  }
})
