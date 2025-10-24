import React, { useState, useEffect, useCallback } from 'react';
import { Card, Row, Col, Statistic, Select, Button, Alert, message } from 'antd';
import { ReloadOutlined, InfoCircleOutlined } from '@ant-design/icons';
import FundTable from '../components/FundTable';
import apiClient from '../services/api';
import { Fund, FundIndexParams } from '../types/fund';
import { formatDateTime } from '../utils';

const { Option } = Select;

const FundIndex: React.FC = () => {
  const [funds, setFunds] = useState<Fund[]>([]);
  const [loading, setLoading] = useState(false);
  const [fundTypes, setFundTypes] = useState<string[]>([]);
  const [updatedAt, setUpdatedAt] = useState<string>('');
  const [allFundCount, setAllFundCount] = useState(0);
  const [fund4433Count, setFund4433Count] = useState(0);
  const [params, setParams] = useState<FundIndexParams>({
    page_num: 1,
    page_size: 20,
    sort: 0,
    type: ''
  });

  const loadFunds = useCallback(async (newParams: FundIndexParams = params) => {
    setLoading(true);
    try {
      console.log('正在加载基金数据，参数:', newParams);
      const response = await apiClient.getFundIndex(newParams);
      console.log('API响应:', response);
      
      setFunds(response.fund_list || []);
      setFundTypes(response.fund_types || []);
      setUpdatedAt(response.updated_at || '');
      setAllFundCount(response.all_fund_count || 0);
      setFund4433Count(response.fund_4433_count || 0);
      message.success('数据加载成功');
    } catch (error: any) {
      console.error('加载基金数据失败:', error);
      const errorMessage = error.response?.data?.message || error.message || '加载基金数据失败';
      message.error(`加载失败: ${errorMessage}`);
    } finally {
      setLoading(false);
    }
  }, [params]);

  useEffect(() => {
    loadFunds();
  }, [loadFunds]);

  const handleTypeChange = (type: string) => {
    const newParams = { ...params, type, page_num: 1 };
    setParams(newParams);
    loadFunds(newParams);
  };

  const handleSortChange = (sort: number) => {
    const newParams = { ...params, sort, page_num: 1 };
    setParams(newParams);
    loadFunds(newParams);
  };

  const handleRefresh = () => {
    loadFunds();
  };

  return (
    <div className="investool-container">
      <Card>
        <Row gutter={16} style={{ marginBottom: 16 }}>
          <Col xs={24} sm={8}>
            <Statistic
              title="4433基金总数"
              value={fund4433Count}
              suffix={`/ ${allFundCount}`}
            />
          </Col>
          <Col xs={24} sm={8}>
            <Statistic
              title="筛选总数"
              value={allFundCount}
            />
          </Col>
          <Col xs={24} sm={8}>
            <Statistic
              title="更新时间"
              value={updatedAt ? formatDateTime(updatedAt) : '--'}
            />
          </Col>
        </Row>

        <Alert
          message="4433法则说明"
          description={
            <div>
              <p>由台大财务金融学系邱显比教授提出的选基法则：</p>
              <ul>
                <li><strong>4</strong>: 最近1年收益率排名在同类型基金前1/4</li>
                <li><strong>4</strong>: 最近2年、3年、5年及今年来收益率排名均在同类型基金前1/4</li>
                <li><strong>3</strong>: 最近6个月收益率排名在同类型基金前1/3</li>
                <li><strong>3</strong>: 最近3个月收益率排名在同类型基金前1/3</li>
              </ul>
            </div>
          }
          type="info"
          icon={<InfoCircleOutlined />}
          style={{ marginBottom: 16 }}
        />

        <Row gutter={16} style={{ marginBottom: 16 }}>
          <Col xs={24} sm={8}>
            <Select
              placeholder="选择基金类型"
              style={{ width: '100%' }}
              value={params.type}
              onChange={handleTypeChange}
              allowClear
            >
              <Option value="">全部类型</Option>
              {fundTypes.map(type => (
                <Option key={type} value={type}>{type}</Option>
              ))}
            </Select>
          </Col>
          <Col xs={24} sm={8}>
            <Select
              placeholder="选择排序方式"
              style={{ width: '100%' }}
              value={params.sort}
              onChange={handleSortChange}
            >
              <Option value={0}>按周收益率排序</Option>
              <Option value={1}>按月收益率排序</Option>
              <Option value={2}>按年收益率排序</Option>
              <Option value={3}>按基金规模排序</Option>
            </Select>
          </Col>
          <Col xs={24} sm={8}>
            <Button
              icon={<ReloadOutlined />}
              onClick={handleRefresh}
              loading={loading}
              style={{ width: '100%' }}
            >
              刷新数据
            </Button>
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

export default FundIndex;
