import React, { useState, useEffect } from 'react';
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

  const loadFunds = async (params: FundFilterParams = {}) => {
    setLoading(true);
    try {
      const response = await apiClient.getFundFilter(params);
      setFunds(response.fund_list || []);
      setFundTypes(response.fund_types || []);
      setTotalCount(response.fund_4433_count || 0);
    } catch (error) {
      console.error('加载基金数据失败:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadFunds();
  }, []);

  const handleFilter = (values: FundFilterParams) => {
    loadFunds(values);
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
        />
      </Card>
    </div>
  );
};

export default FundFilter;
