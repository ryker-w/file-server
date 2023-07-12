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
.croppers-container {

.cropper-img-warp {
  text-align: center;

.cropper-img {
  margin: auto;
  width: 150px;
  height: 150px;
  border-radius: 100%;
}

}
}
.image {
  width: 100%;
  height: 100%;
}
</style>
