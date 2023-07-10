<template>
  <div class="croppers-container layout-pd">
    <el-card shadow="hover" :header="state.path">
      <el-alert
          :title="state.title"
          type="success"
          :closable="false"
          class="mb15"
      ></el-alert>
    </el-card>
  </div>
</template>

<script setup lang="ts" name="explorer">
import { onMounted, reactive} from 'vue';
import {RouteParamValue, useRoute} from 'vue-router';
import { useFsApi } from "/@/api/fs";

const route = useRoute()
const fsApi = useFsApi()

// 定义变量内容
const state = reactive({
  path: '/1/2./4',
  title: '文件夹',
});

const openFolder = async (paths: string[]) => {
  await fsApi.explorer({
    "path": paths,
  }).then(res => {
    if (res && res.code == 200) {
      let data = res.data
      console.log(data["folders"])
      console.log(data["files"])
    } else {
      console.log(res.code)
    }
  }).catch(err => {
    console.log(err)
  })
}

const onOpenFolder = async () => {
  let paths = <RouteParamValue[]>route.params["path"]
  await openFolder(paths)
}

const updatePath = () => {
  console.log(route.params)
  let paths = <RouteParamValue[]>route.params["path"]
  let path = ""
  for (let i in paths) {
    path += '/' + paths[i]
  }
  state.path = path
}
onMounted(() => {
  updatePath()
  let paths = ["/"]
  openFolder(paths)
})
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
</style>
