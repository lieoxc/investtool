// 工具函数
export const formatNumber = (num: number, decimals: number = 2): string => {
  if (isNaN(num)) return '--';
  return num.toFixed(decimals);
};

export const formatCurrency = (num: number, unit: string = '元'): string => {
  if (isNaN(num)) return '--';
  
  if (unit === '元') {
    const yi = num / 100000000.0;
    if (Math.abs(yi) >= 1) {
      return `${yi.toFixed(2)}亿元`;
    } else if (num / 10000.0 >= 1) {
      return `${(num / 10000.0).toFixed(2)}万元`;
    }
    return `${num.toFixed(2)}元`;
  }
  
  return `${num.toFixed(2)}${unit}`;
};

export const formatPercentage = (num: number, decimals: number = 2): string => {
  if (isNaN(num)) return '--';
  return `${num.toFixed(decimals)}%`;
};

export const formatDate = (date: string | Date): string => {
  if (!date) return '--';
  const d = new Date(date);
  return d.toLocaleDateString('zh-CN');
};

export const formatDateTime = (date: string | Date): string => {
  if (!date) return '--';
  const d = new Date(date);
  return d.toLocaleString('zh-CN');
};

// 获取基金类型颜色
export const getFundTypeColor = (type: string): string => {
  const colorMap: Record<string, string> = {
    '股票型': '#1890ff',
    '混合型': '#52c41a',
    '债券型': '#faad14',
    '货币型': '#722ed1',
    '指数型': '#eb2f96',
    'QDII': '#13c2c2',
    'FOF': '#fa8c16',
  };
  return colorMap[type] || '#666666';
};

// 判断是否为4433基金
export const is4433Fund = (fund: any): boolean => {
  if (!fund || !fund.performance) return false;
  
  const p = fund.performance;
  return (
    p.year_1_rank_ratio <= 25 &&
    p.year_2_rank_ratio <= 25 &&
    p.year_3_rank_ratio <= 25 &&
    p.year_5_rank_ratio <= 25 &&
    p.this_year_rank_ratio <= 25 &&
    p.month_6_rank_ratio <= 33.33 &&
    p.month_3_rank_ratio <= 33.33
  );
};

// 计算基金评分
export const calculateFundScore = (fund: any): number => {
  if (!fund) return 0;
  
  let score = 0;
  
  // 绩效排名评分 (40%)
  const performanceScore = (
    (100 - fund.performance?.year_1_rank_ratio || 0) * 0.2 +
    (100 - fund.performance?.year_3_rank_ratio || 0) * 0.1 +
    (100 - fund.performance?.month_6_rank_ratio || 0) * 0.1
  );
  
  // 风险指标评分 (30%)
  const riskScore = (
    Math.max(0, 100 - (fund.stddev?.avg_135 || 0) * 2) * 0.15 +
    Math.max(0, (fund.sharp?.avg_135 || 0) * 10) * 0.15
  );
  
  // 基金经理评分 (20%)
  const managerScore = Math.min(100, (fund.manager?.manage_days || 0) / 365 * 20);
  
  // 规模评分 (10%)
  const scaleScore = Math.min(100, (fund.net_assets_scale || 0) / 100000000 * 10);
  
  score = performanceScore + riskScore + managerScore + scaleScore;
  
  return Math.round(Math.max(0, Math.min(100, score)));
};

// 防抖函数
export const debounce = <T extends (...args: any[]) => any>(
  func: T,
  wait: number
): ((...args: Parameters<T>) => void) => {
  let timeout: NodeJS.Timeout;
  return (...args: Parameters<T>) => {
    clearTimeout(timeout);
    timeout = setTimeout(() => func(...args), wait);
  };
};

// 节流函数
export const throttle = <T extends (...args: any[]) => any>(
  func: T,
  limit: number
): ((...args: Parameters<T>) => void) => {
  let inThrottle: boolean;
  return (...args: Parameters<T>) => {
    if (!inThrottle) {
      func(...args);
      inThrottle = true;
      setTimeout(() => (inThrottle = false), limit);
    }
  };
};

// 深拷贝
export const deepClone = <T>(obj: T): T => {
  if (obj === null || typeof obj !== 'object') return obj;
  if (obj instanceof Date) return new Date(obj.getTime()) as any;
  if (obj instanceof Array) return obj.map(item => deepClone(item)) as any;
  if (typeof obj === 'object') {
    const clonedObj = {} as any;
    for (const key in obj) {
      if (obj.hasOwnProperty(key)) {
        clonedObj[key] = deepClone(obj[key]);
      }
    }
    return clonedObj;
  }
  return obj;
};
