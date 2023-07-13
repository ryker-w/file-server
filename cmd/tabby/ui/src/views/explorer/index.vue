<template>
  <div class="croppers-container layout-pd">
  </div>
</template>

<script setup lang="ts" name="explorer">
import {onMounted} from 'vue';
import {useRouter} from "vue-router";

const router = useRouter()

onMounted(() => {
  next()
})
//中转
const next = () => {
  router.push("/fs/")
}
</script>

<style scoped lang="scss">

</style>
