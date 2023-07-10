import request from '/@/utils/request';

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
	};
}
