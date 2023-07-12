import request, {exportXLS} from '/@/utils/request';

export function useFsApi() {
	return {
		explorer: (data: object) => {
			return request({
				url: '/api/fs/explorer',
				method: 'post',
				data,
			});
		},
		download: (data: object) => {
			return request({
				url: '/api/fs/download',
				method: 'post',
				data,
			});
		},
		upload: (data: object) => {
			return request({
				url: '/api/fs/upload',
				method: 'post',
				data,
			});
		},
	};
}

export const upload = "/api/fs/upload";
export const downloadApi = (p: object) => exportXLS("/api/fs/download",p)