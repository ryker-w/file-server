import axios, { AxiosInstance } from 'axios';
import { ElMessage, ElMessageBox } from 'element-plus';
import { Session } from '/@/utils/storage';
import qs from 'qs';

// 配置新建一个 axios 实例
const service: AxiosInstance = axios.create({
	baseURL: import.meta.env.VITE_API_URL,
	timeout: 50000,
	headers: { 'Content-Type': 'application/json' },
	paramsSerializer: {
		serialize(params) {
			return qs.stringify(params, { allowDots: true });
		},
	},
});

// 添加请求拦截器
service.interceptors.request.use(
	(config) => {
		// 在发送请求之前做些什么 token
		if (Session.get('token')) {
			config.headers!['Authorization'] = `${Session.get('token')}`;
		}
		return config;
	},
	(error) => {
		// 对请求错误做些什么
		return Promise.reject(error);
	}
);

// 添加响应拦截器
service.interceptors.response.use(
	(response) => {
		// 对响应数据做点什么
		const res = response.data;
		/*if (res.code && res.code !== 0) {
			// `token` 过期或者账号已在别处登录
			if (res.code === 401 || res.code === 4001) {
				Session.clear(); // 清除浏览器全部临时缓存
				window.location.href = '/'; // 去登录页
				ElMessageBox.alert('你已被登出，请重新登录', '提示', {})
					.then(() => {})
					.catch(() => {});
			}
			return Promise.reject(service.interceptors.response);
		} else {*/
		// 导出报表的接口需要headers
		if (response.headers["content-disposition"] && response.headers["content-disposition"] !== null) {
			return response
		}
		return res;
		// }
	},
	(error) => {
		// 对响应错误做点什么
		if (error.message.indexOf('timeout') != -1) {
			ElMessage.error('网络超时');
		} else if (error.message == 'Network Error') {
			ElMessage.error('网络连接错误');
		} else {
			if (error.response.data) ElMessage.error(error.response.statusText);
			else ElMessage.error('接口路径找不到');
		}
		return Promise.reject(error);
	}
);
export function exportXLS(url: string, params: object) {
	return service.request({
		method: 'POST',
		url: url,
		data: params,
		headers: {
			'Content-Type': 'application/json'
		},
		responseType: 'blob'
	}).then(res => {
		console.log(res)
		const link = document.createElement('a')
		const blob = new Blob([res.data], { type: 'application/octet-stream' })
		link.style.display = 'none'
		link.href = URL.createObjectURL(blob)
		let fileName = res.headers["content-disposition"]
		console.log(fileName)
		fileName = decodeURI(escape(fileName.substring(fileName.indexOf("=")+1)))
		link.setAttribute('download', fileName)
		document.body.appendChild(link)
		link.click()
		document.body.removeChild(link)
	})
}
// 导出 axios 实例
export default service;
