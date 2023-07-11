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

<script setup lang="ts" name="folder">
import { onMounted, reactive} from 'vue';
import {RouteParamValue, useRoute} from 'vue-router';

const route = useRoute()

// 定义变量内容
const state = reactive({
	path: '/1/2./4',
  title: '文件夹',
});

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
