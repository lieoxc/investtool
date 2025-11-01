import axios, { AxiosInstance, AxiosResponse } from 'axios';
import {
  FundIndexParams,
  FundIndexResponse,
  FundFilterParams,
  FundCheckParams,
  FundCheckResponse,
  FundManagerParams,
  FundManagerResponse,
  FundSimilarityParams,
  QueryByStockParams,
  ApiResponse
} from '../types/fund';

class ApiClient {
  private client: AxiosInstance;

  constructor() {
    // 在生产环境使用相对路径，让 nginx 代理到后端
    // 在开发环境使用环境变量或通过 react-scripts 的 proxy 配置
    const getBaseURL = () => {
      // 如果设置了环境变量，优先使用
      if (process.env.REACT_APP_API_BASE_URL) {
        return process.env.REACT_APP_API_BASE_URL;
      }
      // 生产环境（已构建的应用）使用相对路径，通过 nginx 代理
      if (process.env.NODE_ENV === 'production') {
        return '/api';
      }
      // 开发环境可以使用 proxy 配置或直接使用 localhost
      // react-scripts 的 proxy 配置会自动处理相对路径的请求
      return '';
    };

    this.client = axios.create({
      baseURL: getBaseURL(),
      timeout: 30000,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    // 请求拦截器
    this.client.interceptors.request.use(
      (config) => {
        console.log('API Request:', config.method?.toUpperCase(), config.url);
        return config;
      },
      (error) => {
        console.error('Request Error:', error);
        return Promise.reject(error);
      }
    );

    // 响应拦截器
    this.client.interceptors.response.use(
      (response: AxiosResponse) => {
        console.log('API Response:', response.status, response.config.url);
        return response;
      },
      (error) => {
        console.error('Response Error:', error.response?.status, error.message);
        return Promise.reject(error);
      }
    );
  }

  // 基金首页 - 4433基金列表
  async getFundIndex(params: FundIndexParams = {}): Promise<FundIndexResponse> {
    const response = await this.client.get('/api/fund', { params });
    // 后端返回的数据结构是 { code, message, data }
    if (response.data && response.data.data) {
      return response.data.data;
    }
    return response.data;
  }

  // 基金筛选
  async getFundFilter(params: FundFilterParams = {}): Promise<FundIndexResponse> {
    const response = await this.client.get('/api/fund/filter', { params });
    // 后端返回的数据结构是 { code, message, data }
    if (response.data && response.data.data) {
      return response.data.data;
    }
    return response.data;
  }

  // 基金检测
  async postFundCheck(params: FundCheckParams): Promise<FundCheckResponse> {
    const response = await this.client.post('/api/fund/check', params);
    // 后端返回的数据结构是 { code, message, data }
    if (response.data && response.data.data) {
      return response.data.data;
    }
    return response.data;
  }

  // 基金经理筛选
  async getFundManagers(params: FundManagerParams = {}): Promise<FundManagerResponse> {
    const response = await this.client.get('/api/fund/managers', { params });
    // 后端返回的数据结构是 { code, message, data }
    if (response.data && response.data.data) {
      return response.data.data;
    }
    return response.data;
  }

  // 基金持仓相似度
  async getFundSimilarity(params: FundSimilarityParams): Promise<ApiResponse> {
    const response = await this.client.get('/api/fund/similarity', { params });
    // 后端返回的数据结构是 { code, message, data }
    if (response.data && response.data.data) {
      return response.data.data;
    }
    return response.data;
  }

  // 股票选基
  async postQueryByStock(params: QueryByStockParams): Promise<ApiResponse> {
    const response = await this.client.post('/api/fund/query_by_stock', params);
    // 后端返回的数据结构是 { code, message, data }
    if (response.data && response.data.data) {
      return response.data.data;
    }
    return response.data;
  }

  // 通用 GET 请求
  async get<T = any>(url: string, params?: any): Promise<T> {
    const response = await this.client.get(url, { params });
    return response.data;
  }

  // 通用 POST 请求
  async post<T = any>(url: string, data?: any): Promise<T> {
    const response = await this.client.post(url, data);
    return response.data;
  }
}

// 创建单例实例
const apiClient = new ApiClient();

export default apiClient;
