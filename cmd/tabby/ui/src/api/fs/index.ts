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
		upload: (data: object) => {
			return request({
				url: '/api/fs/upload',
				method: 'post',
				data,
			});
		},
		delete: (data: object) => {
			return request({
				url: '/api/fs/delete',
				method: 'post',
				data,
			});
		},
	};
}

export const upload = "/api/fs/upload";

export const downloadApi = (p: object) => exportXLS("/api/fs/download",p)