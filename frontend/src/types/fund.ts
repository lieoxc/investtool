// 基金相关类型定义
export interface Fund {
  code: string;
  name: string;
  type: string;
  net_assets_scale: number;
  manager: FundManager;
  performance: FundPerformance;
  stddev: FundStddev;
  sharp: FundSharp;
  max_retracement: FundMaxRetracement;
  stocks?: FundStock[];
}

export interface FundManager {
  id: string;
  name: string;
  manage_days: number;
  manage_repay: number;
}

export interface FundPerformance {
  year_1_rank_ratio: number;
  year_2_rank_ratio: number;
  year_3_rank_ratio: number;
  year_5_rank_ratio: number;
  this_year_rank_ratio: number;
  month_6_rank_ratio: number;
  month_3_rank_ratio: number;
}

export interface FundStddev {
  avg_135: number;
  year_1: number;
  year_3: number;
  year_5: number;
}

export interface FundSharp {
  avg_135: number;
  year_1: number;
  year_3: number;
  year_5: number;
}

export interface FundMaxRetracement {
  avg_135: number;
  year_1: number;
  year_3: number;
  year_5: number;
}

export interface FundStock {
  code: string;
  name: string;
  hold_ratio: number;
  industry: string;
  adjust_ratio: number;
}

// API 请求参数
export interface FundIndexParams {
  page_num?: number;
  page_size?: number;
  sort?: number;
  type?: string;
}

export interface FundFilterParams {
  page_num?: number;
  page_size?: number;
  sort?: number;
  type?: string;
  year_1_rank_ratio?: number;
  this_year_235_rank_ratio?: number;
  month_6_rank_ratio?: number;
  month_3_rank_ratio?: number;
  min_scale?: number;
  max_scale?: number;
  min_estab_years?: number;
  min_manager_years?: number;
  types?: string[];
  max_135_avg_stddev?: number;
  min_135_avg_sharp?: number;
  max_135_avg_retr?: number;
}

export interface FundCheckParams {
  fundcode: string;
  min_scale?: number;
  max_scale?: number;
  min_manager_years?: number;
  year_1_rank_ratio?: number;
  this_year_235_rank_ratio?: number;
  month_6_rank_ratio?: number;
  month_3_rank_ratio?: number;
  max_135_avg_stddev?: number;
  min_135_avg_sharp?: number;
  max_135_avg_retr?: number;
  check_stocks?: boolean;
}

export interface FundManagerParams {
  name?: string;
  min_working_years?: number;
  min_yieldse?: number;
  max_current_fund_count?: number;
  min_scale?: number;
  page_num?: number;
  page_size?: number;
  sort?: string;
  fund_type?: string;
}

export interface FundSimilarityParams {
  codes: string;
}

export interface QueryByStockParams {
  keywords: string;
}

// API 响应类型
export interface ApiResponse<T = any> {
  data?: T;
  error?: string;
  message?: string;
}

export interface FundIndexResponse {
  fund_list: Fund[];
  pagination: Pagination;
  updated_at: string;
  all_fund_count: number;
  fund_4433_count: number;
  fund_types: string[];
}

export interface FundCheckResponse {
  funds: Fund[];
  param: FundCheckParams;
  stock_check_results?: Record<string, any>;
}

export interface FundManagerResponse {
  managers: FundManagerInfo[];
  pagination: Pagination;
}

export interface FundManagerInfo {
  id: string;
  name: string;
  working_years: number;
  yieldse: number;
  current_fund_count: number;
  scale: number;
  best_fund_is_4433: boolean;
  current_best_fund_code: string;
}

export interface Pagination {
  page_num: number;
  page_size: number;
  total: number;
  total_pages: number;
  start_index: number;
  end_index: number;
}
