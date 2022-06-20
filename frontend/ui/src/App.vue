<script setup>
import {  ref } from 'vue'
import JsonVisualizer from './components/JsonVisualizer.vue';
import {
  darkTheme, NConfigProvider, NButton, NCard, NSpace, NGlobalStyle, NH2, useOsTheme, enUS,
  NGrid, NGridItem
} from 'naive-ui'
import { } from 'naive-ui'
import dateEnUs from 'naive-ui/es/locales/date/enUS';


const locale = ref(enUS)
const localeDate = ref(dateEnUs)
const theme = ref(useOsTheme().value === 'dark' ? darkTheme : null)

if (global === undefined) {
  var global = window;
}
</script>

<template>
  <n-config-provider :theme="theme" :locale="locale" :date-locale="localeDate">
    <n-card :bordered="true" class="header">
      <n-grid cols="1 s:2" item-responsive responsive="screen">
        
          <n-grid-item span="1 400:1 600:1 800:1">
            <n-h2>S3 JSON Editor and Visualizer</n-h2>
          </n-grid-item>
          <n-grid-item span="1 400:1 600:1 800:1">
            <n-space justify="end">
            <n-button @click="theme = darkTheme">Dark</n-button>
            <n-button @click="theme = null">Light</n-button>
            </n-space >
          </n-grid-item>

      </n-grid>

    </n-card>
    <JsonVisualizer />


    <n-global-style />
  </n-config-provider>


</template>
<style>
.boxtest {

  border: 10px solid red;
}

.header {
  height: 10vh;
}
</style>