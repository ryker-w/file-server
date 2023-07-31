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
            :title="state.path"
            type="success"
            :closable="false"
            class="mb15"
            :description="'文件夹：'+state.folderCount+'，文件：'+state.fileCount"
        ></el-alert>
        <el-upload
            class="upload-file"
            drag
            :action="upload"
            :data="{path:[state.path]}"
            multiple
            :show-file-list="false"
            :on-success="onSuccess"
            :on-error="onError"
            :before-upload="beforeUpload"
        >
          <el-icon class="el-icon-upload">
            <upload-filled/>
          </el-icon>
          <div class="el-upload__text">
            将文件放到这里或 <em>手动选择</em>
          </div>
        </el-upload>
      </div>
      <div class="file-row">
        <el-row :gutter="10">
          <el-col :span="4" v-for="(item,index) in state.datas" :key="index">
            <el-card shadow="hover" v-if="item.type=='folder'"
                     style="margin: 5px;height: 280px; ">
              <div style="text-align: center" @click="onNext(state.path+'/'+item.name)">
<!--                <img :src="setSVG(item.suffix)" class="image">-->
                <SvgIcon :name="setSVG(item.suffix)" class="icon" />
<!--                <i class="icon" :class="setSVG(item.suffix)"></i>-->
              </div>
              <div style="text-align: center;padding: 14px">
                <el-tooltip :content="item.name" effect="dark" placement="top-end">
                <div class="content-span">{{ item.name }}</div>
                </el-tooltip>
              </div>
              <div class="bottom">
                <el-row :gutter="20">
                  <el-col :span="12">
                    <el-button size="mini"
                               icon="ele-FolderDelete"
                               round
                               @click="deleteFile(state.path,item.name)"
                               type="danger">删除
                  </el-button></el-col>
                  <el-col :span="12">
                  </el-col>
                </el-row>
              </div>
            </el-card>
            <el-card shadow="hover" v-else style="margin: 5px; height: 280px;">
              <div style="text-align: center">
<!--                <img :src="setSVG(item.suffix)" class="image">-->
                <SvgIcon :name="setSVG(item.suffix)" class="icon" />
<!--                <i class="icon" :class="setSVG(item.suffix)"></i>-->
              </div>
              <div style="text-align: center;padding: 14px">
                <el-tooltip :content="item.name" effect="dark" placement="top-end">
                  <div class="content-span">{{ item.name }}</div>
                </el-tooltip>
              </div>
              <div class="bottom">
                <el-row :gutter="20">
                  <el-col :span="12"><el-button size="mini"  icon="ele-FolderDelete" round @click="deleteFile(state.path,item.name)"
                                     type="danger">删除
                  </el-button></el-col>
                  <el-col :span="12">
                    <el-button size="mini" icon="ele-Download" round @click="download(state.path+'/'+item.name,item.name)"
                               type="primary">下载
                    </el-button>
                  </el-col>
                </el-row>


              </div>
            </el-card>
          </el-col>
        </el-row>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts" name="folder">
import {onMounted, reactive} from 'vue';
import {RouteParamValue, useRoute, useRouter} from 'vue-router';
import {downloadApi, upload, useFsApi} from "/@/api/fs";
import {ElMessage, ElMessageBox} from "element-plus";
import Ext from "/@/utils/ext";
import "../../../theme/ali/iconfont.css";

const ext = Ext

const route = useRoute()
const router = useRouter()
const fsApi = useFsApi()

// 定义变量内容
const state = reactive({
  path: '/',
  title: '/',
  datas: [],
  routes: [],
  folderCount: 0,
  fileCount: 0,
});
const beforeUpload = (file: any) => {
  const m = 1024*1024
  const limit = 50 // 50M
  let fileSizeM = file.size / m
  return fileSizeM < limit
}
const deleteFile = (path: any, name: any) => {
  ElMessageBox.confirm(`此操作将永久删除文件, 是否继续?`, '提示', {
    confirmButtonText: '删除',
    cancelButtonText: '取消',
    type: 'warning',
  }).then(() => {
    fsApi.delete({path: [path], name: name}).then(res => {
      if (res && res.code == 200) {
        ElMessage.success("删除成功!");
        updatePath()
      } else {
        ElMessage.error("删除失败！");
      }
    })
  }).catch(() => {
  });
}
const onSuccess = (response: any, file: any, fileList: any) => {
  if (response && response.code == 200) {
    ElMessage.success("上传成功!");
    updatePath()
  } else {
    ElMessage.error("上传失败！");
  }
}
const onError = (err: any, file: any, fileList: any) => {
  console.log("上传失败：", err)
  ElMessage.error("上传失败！");
}
const download = (filePath: any, fileName: any) => {
  downloadApi({
    filePath: filePath,
    fileName: fileName
  })
}

const onNext = (path) => {
  state.path = path
  // console.log("onNext", state.path)
  router.push("/fs" + state.path)
}

const setSVG =  (suffix: any) => {
  if(suffix=='folder'){
    return "iconfont icon-folder"
  }
  if (suffix) {
    for(let svg in ext){
        if(svg==suffix.toLowerCase()){
          return "iconfont "+ext[svg]
        }
    }
  }
  return "iconfont icon-file"
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
        state.folderCount = res.folders.length
      }
      if (res.files) {
        for (let file of res.files) {
          files.push({
            name: file.name,
            type: "file",
            suffix: file.ext ? file.ext : 'file'
          })
        }
        state.fileCount = res.files.length
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
.image {
  width: 200px;
  height: 200px;
}
.file-row{
  height: 1000px;
  overflow-y: auto;
}
.upload-file{
  /*height: 100px;*/
}
.bottom{
  width: 100%;
  text-align: center;
}
.icon {
  width: 1em;
  height: 1em;
  vertical-align: -0.15em;
  fill: currentColor;
  overflow: hidden;
  font-size: 150px !important;
}
.content-span{
  width: 100%;
  overflow:hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  -o-text-overflow:ellipsis;
}
</style>
