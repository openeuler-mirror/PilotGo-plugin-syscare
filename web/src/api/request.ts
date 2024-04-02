import axios from 'axios';
import router from '@/router';

// 公共定义
export const RespCodeOK = 200
export interface ResultData {
  code?: number;
  data?: any[];
  msg?: string;
  ok?:boolean;
  page?:number;
  size?: number;
  total?: number;
}

// 1.创建axios实例
const request = axios.create({
  baseURL: '',
  timeout: 5000
});

// 2.1添加请求拦截器
request.interceptors.request.use(
  (config) => {  
    // 根据不同的请求类型设置不同的 content-type  
    if (config.method === 'post') {  
      if (config.data instanceof FormData) {  
        // 如果是 FormData，通常用于文件上传，不需要设置 content-type  
        // 因为浏览器会自动设置正确的 boundary  
      } else {  
        // 如果是普通 post 请求，设置 content-type 为 application/json  
        config.headers['Content-Type'] = 'application/json;charset=utf-8';  
      }  
    } else if (config.method === 'get') {  
      // 如果是 get 请求，通常不需要设置 content-type  
    }  
    // 其他逻辑...  
    return config;  
  },  
  (error) => {
    return Promise.reject(error);
  },
);

// 2.2添加响应拦截器
request.interceptors.response.use(
  (response: any) => {
    if (response.data && response.data.code == '401') {
      router.push('/');
    } else {
      return response;
    }
  },
  (error) => {
    if (error.response) {
      switch (error.response.status) {
        case 401:
          router.push('/');
      }
      return Promise.reject(error.response.data);
    }
  },
);

export default request;
