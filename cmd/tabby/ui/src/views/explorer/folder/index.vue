<template>
  <div class="croppers-container layout-pd">
    <el-card shadow="hover">
      <template #header>
        <div style="height: 20px;">
          <el-breadcrumb separator-class="el-icon-arrow-right">
            <el-breadcrumb-item v-for="(item,index) in state.routes" :key="index" :to="{ path: item.path }">
              {{ item.name }}
            </el-breadcrumb-item>
          </el-breadcrumb>
        </div>
      </template>
      <div>
        <el-alert
            :title="state.title"
            type="success"
            :closable="false"
            class="mb15"
        ></el-alert>
        <el-upload
            class="upload-demo"
            drag
            :action="upload"
            multiple
        >
          <el-icon class="el-icon-upload">
            <upload-filled/>
          </el-icon>
          <div class="el-upload__text">
            将文件放到这里或 <em>手动选择</em>
          </div>
        </el-upload>
      </div>
    </el-card>
    <div class="file-row">
      <el-row :gutter="10">
        <el-col :span="4" v-for="(item,index) in state.datas" :key="index">
          <el-card v-if="item.type=='folder'" @click="onNext(state.path+'/'+item.name)" style="margin: 5px;">
            <img :src="setSVG(item.suffix)" class="image">
            <div style="padding: 14px">
              <span>{{ item.name }}</span>
            </div>
          </el-card>
          <el-card v-else style="margin: 5px;">
            <img :src="setSVG(item.suffix)" class="image">
            <div style="padding: 14px">
              <span>{{ item.name }}</span>
              <el-button @click="download(state.path+'/'+item.name,item.name)" text style="float: right;">下载</el-button>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup lang="ts" name="folder">
import {onMounted, reactive} from 'vue';
import {RouteParamValue, useRoute, useRouter} from 'vue-router';
import {downloadApi, upload, useFsApi} from "/@/api/fs";

const route = useRoute()
const router = useRouter()
const fsApi = useFsApi()
// 定义变量内容
const state = reactive({
  path: '/',
  title: '/',
  datas: [],
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
  ],
  routes: []
});

const download = (filePath: any, fileName: any) => {
  downloadApi({
    filePath: filePath,
    fileName: fileName
  })
}

const onNext = (path) => {
  state.path = state.path + path
  // console.log("n", state.path)
  router.push("/fs" + state.path)
}

const setSVG = (suffix: any) => {
  if (suffix && state.suffixs.indexOf(suffix.toLowerCase()) > -1) {
    return "./src/assets/file/" + suffix.toLowerCase() + ".svg"
  }
  return "./src/assets/file/file.svg"
}
const updatePath = async () => {
  let paths = await <RouteParamValue[]>route.params["path"]
  // console.log("paths", paths)
  let routes = []
  routes.push({
    path: "/fs",
    name: "/"
  })
  if (paths && paths.length > 0) {
    var newArr = await paths.filter(function (item, index) {
      return paths.indexOf(item) === index;  // 因为indexOf 只能查找到第一个
    });
    let path = ""
    for (let i in newArr) {
      if (newArr[i]) {
        path = path + '/' + newArr[i]
        routes.push({
          path: "/fs" + path,
          name: newArr[i]
        })
      }
    }
    state.path = path
  }
  state.routes = routes
  onOpenFolder()
}
const onOpenFolder = async () => {
  await openFolder([state.path])
}
const openFolder = async (paths: string[]) => {
  let files = []
  await fsApi.explorer({
    "path": paths,
  }).then(res => {
    if (res && res.code == 200) {
      if (res.folders) {
        for (let folder of res.folders) {
          files.push({
            name: folder,
            type: "folder",
            suffix: "folder"
          })
        }
      }
      if (res.files) {
        for (let file of res.files) {
          files.push({
            name: file.name,
            type: "file",
            suffix: file.ext ? file.ext : 'file'
          })
        }
      }
      state.datas = files;
    } else {
      console.log(res.code)
    }
  }).catch(err => {
    console.log(err)
  })
}
onMounted(() => {
  updatePath()
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
