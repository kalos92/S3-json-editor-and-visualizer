<script setup>

import { NModal } from 'naive-ui'

const props = defineProps({
  error: String,
  open: Boolean,
  title: String,
})

function getTitle(){
  return props.title
}

</script>
<template>


  <n-modal
    v-model:show="props.open"
    :mask-closable="false"
    preset="dialog"
    :title="getTitle()" 
    positive-text="Got it!"
    @positive-click="$emit('close-modal')"
    @negative-click="$emit('close-modal')"
    @close="$emit('close-modal')"
  >S3rver says: {{props.error}}</n-modal>


</template>
