<template>
  <div class="croppers-container layout-pd">
    <el-card shadow="hover">
      <template #header>
        <span @click="openFolder([state.path])">{{ state.path }}</span>
      </template>
      <el-alert
          :title="state.title"
          type="success"
          :closable="false"
          class="mb15"
      ></el-alert>
      <div class="file-row">
        <el-row :gutter="10">
          <el-col :span="4" v-for="(item,index) in state.files" :key="index">
            <el-card @click="openFolder([state.title+item.name+'/'])" style="margin: 5px;">
              <img
                  :src="setFileIcon(item.suffix)"
                  class="image"
              />
              <div style="padding: 14px">
                <span>{{ item.name }}</span>
              </div>
            </el-card>
          </el-col>
        </el-row>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts" name="explorer">
import {onMounted, reactive} from 'vue';
import {RouteParamValue, useRoute,useRouter} from 'vue-router';
import {useFsApi} from "/@/api/fs";

const route = useRoute()
const router = useRouter()
const fsApi = useFsApi()

// 定义变量内容
const state = reactive({
  name: '',
  path: '/',
  paths: [],
  title: '/',
  files: [],
  suffixs: [
    "file",
    "folder",
    "gif",
    "png",
    "mp4",
    "pdf",
    "word",
    "ppt",
    "txt",
    "vue",
    "ymal",
    "excel"
  ]
});

const openFolder = async (paths: string[]) => {
  console.log(paths)
  await fsApi.explorer({
    "path": paths,
  }).then(res => {
    let files = []
    if (res && res.code == 200) {
      let len = paths[0].split("/")
      let pathList = []
      if (len.length > 0) {
        for (let item of len) {
          if (item) {
            pathList.push(item)
          }
        }
      }
      state.paths = pathList
      console.log("paths", state.paths)
      let path = '/'
      let title = '/'
      for (let i = 0; i < state.paths.length; i++) {
        if((i+1)<state.paths.length){
          path = path + state.paths[i] + "/"
        }
        console.log(i,path)
        title = title + state.paths[i] + "/"
      }
      state.path = path
      state.title = title
      if (res.folders) {
        for (let folder of res.folders) {
          files.push({
            name: folder,
            type: "folder",
            suffix: "folder"
          })
        }
      }
      // console.log(data["files"])
      if (res.files) {
        for (let file of res.files) {
          // console.log(file)
          files.push({
            name: file.name,
            type: "file",
            suffix: file.ext ? file.ext : 'file'
          })
        }
      }
      state.files = files;
      // console.log(files)
    } else {
      console.log(res.code)
    }
  }).catch(err => {
    console.log(err)
  })
}
const setFileIcon = (suffix: any) => {
  if (suffix && state.suffixs.indexOf(suffix.toLowerCase()) > -1) {
    return "./src/assets/file/" + suffix.toLowerCase() + ".svg"
  }
  return "./src/assets/file/file.svg"
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
  console.log(state.path)
}
onMounted(() => {
  // updatePath()
  let paths = ['/']
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
.image {
  width: 100%;
  height: 100%;
}
</style>
