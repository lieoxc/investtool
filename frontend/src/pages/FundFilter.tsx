import React, { useState, useEffect, useCallback } from 'react';
import { Card, Row, Col, Statistic, Typography, Alert } from 'antd';
import { InfoCircleOutlined } from '@ant-design/icons';
import FundTable from '../components/FundTable';
import FundFilterForm from '../components/FundFilterForm';
import apiClient from '../services/api';
import { Fund, FundFilterParams } from '../types/fund';

const { Title } = Typography;

const FundFilter: React.FC = () => {
  const [funds, setFunds] = useState<Fund[]>([]);
  const [loading, setLoading] = useState(false);
  const [fundTypes, setFundTypes] = useState<string[]>([]);
  const [totalCount, setTotalCount] = useState(0);
  const [total, setTotal] = useState(0);
  const [params, setParams] = useState<FundFilterParams>({
    page_num: 1,
    page_size: 20,
    sort: 0,
    type: ''
  });

  const loadFunds = useCallback(async (newParams?: FundFilterParams) => {
    const paramsToUse = newParams || params;
    setLoading(true);
    try {
      console.log('正在加载基金筛选数据，参数:', paramsToUse);
      const response = await apiClient.getFundFilter(paramsToUse);
      console.log('API响应:', response);
      
      setFunds(response.fund_list || []);
      setFundTypes(response.fund_types || []);
      setTotalCount(response.fund_4433_count || 0);
      // 设置分页总数
      setTotal(response.pagination?.total || response.fund_list?.length || 0);
    } catch (error) {
      console.error('加载基金数据失败:', error);
    } finally {
      setLoading(false);
    }
  }, [params]);

  useEffect(() => {
    loadFunds();
  }, [loadFunds]);

  const handleFilter = (values: FundFilterParams) => {
    // 重置到第一页，保留排序和分页大小设置
    const newParams: FundFilterParams = { 
      ...values, 
      page_num: 1, 
      page_size: params.page_size || 20,
      sort: values.sort !== undefined ? values.sort : params.sort || 0
    };
    setParams(newParams);
    loadFunds(newParams);
  };

  const handlePageChange = (page: number, pageSize?: number) => {
    const newParams = { 
      ...params, 
      page_num: page,
      page_size: pageSize || params.page_size || 20
    };
    setParams(newParams);
    loadFunds(newParams);
  };

  return (
    <div className="investool-container">
      <Card>
        <Title level={2}>4433基金严选</Title>
        
        <Alert
          message="筛选说明"
          description="根据4433法则和风险指标对基金进行严格筛选，帮助您找到更优质的基金产品。"
          type="info"
          icon={<InfoCircleOutlined />}
          style={{ marginBottom: 16 }}
        />

        <FundFilterForm
          onSubmit={handleFilter}
          loading={loading}
          fundTypes={fundTypes}
        />

        <Row gutter={16} style={{ marginBottom: 16 }}>
          <Col xs={24} sm={8}>
            <Statistic
              title="筛选结果"
              value={funds.length}
              suffix={`/ ${totalCount}`}
            />
          </Col>
          <Col xs={24} sm={8}>
            <Statistic
              title="基金类型"
              value={fundTypes.length}
              suffix="种"
            />
          </Col>
          <Col xs={24} sm={8}>
            <Statistic
              title="筛选条件"
              value="自定义"
            />
          </Col>
        </Row>

        <FundTable
          data={funds}
          loading={loading}
          showScore={true}
          pagination={{
            current: params.page_num || 1,
            pageSize: params.page_size || 20,
            total: total,
            showSizeChanger: true,
            showQuickJumper: true,
            showTotal: (total, range) => 
              `第 ${range[0]}-${range[1]} 条，共 ${total} 条`,
            onChange: handlePageChange,
            onShowSizeChange: handlePageChange,
          }}
        />
      </Card>
    </div>
  );
};

export default FundFilter;
