import React, { useState, useEffect, useCallback } from 'react';
import { Card, Typography, Alert, Table, Tag, Button, Space, Row, Col, Statistic, Select, InputNumber, Input } from 'antd';
import { LinkOutlined, TeamOutlined } from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import apiClient from '../services/api';
import { FundManagerParams, FundManagerInfo } from '../types/fund';
import { formatCurrency, formatPercentage } from '../utils';

const { Title } = Typography;
const { Option } = Select;

const FundManagers: React.FC = () => {
  const [managers, setManagers] = useState<FundManagerInfo[]>([]);
  const [loading, setLoading] = useState(false);
  const [params, setParams] = useState<FundManagerParams>({
    min_working_years: 8,
    min_yieldse: 15.0,
    max_current_fund_count: 10,
    min_scale: 60.0,
    page_num: 1,
    page_size: 20,
    sort: 'yieldse',
    fund_type: '',
    name: ''
  });

  const loadManagers = useCallback(async (newParams: FundManagerParams = params) => {
    setLoading(true);
    try {
      const response = await apiClient.getFundManagers(newParams);
      setManagers(response.managers || []);
    } catch (error) {
      console.error('加载基金经理数据失败:', error);
    } finally {
      setLoading(false);
    }
  }, [params]);

  useEffect(() => {
    loadManagers();
  }, [loadManagers]);

  const handleFilterChange = (key: keyof FundManagerParams, value: any) => {
    const newParams = { ...params, [key]: value, page_num: 1 };
    setParams(newParams);
    loadManagers(newParams);
  };

  const columns: ColumnsType<FundManagerInfo> = [
    {
      title: '基金经理',
      dataIndex: 'name',
      key: 'name',
      width: 150,
      render: (name: string, record: FundManagerInfo) => (
        <Space>
          <TeamOutlined />
          <Button
            type="link"
            onClick={() => window.open(`https://appunit.1234567.com.cn/fundmanager/manager.html?managerid=${record.id}`, '_blank')}
            icon={<LinkOutlined />}
          >
            {name}
          </Button>
        </Space>
      ),
    },
    {
      title: '从业年限',
      dataIndex: 'working_years',
      key: 'working_years',
      width: 100,
      render: (years: number) => `${years}年`,
      sorter: (a, b) => a.working_years - b.working_years,
    },
    {
      title: '年化回报',
      dataIndex: 'yieldse',
      key: 'yieldse',
      width: 100,
      render: (yieldse: number) => (
        <span style={{ color: yieldse >= 15 ? '#52c41a' : '#ff4d4f' }}>
          {formatPercentage(yieldse)}
        </span>
      ),
      sorter: (a, b) => a.yieldse - b.yieldse,
    },
    {
      title: '管理基金数',
      dataIndex: 'current_fund_count',
      key: 'current_fund_count',
      width: 120,
      render: (count: number) => (
        <Tag color={count <= 5 ? 'green' : count <= 10 ? 'orange' : 'red'}>
          {count}只
        </Tag>
      ),
      sorter: (a, b) => a.current_fund_count - b.current_fund_count,
    },
    {
      title: '管理规模',
      dataIndex: 'scale',
      key: 'scale',
      width: 120,
      render: (scale: number) => formatCurrency(scale),
      sorter: (a, b) => a.scale - b.scale,
    },
    {
      title: '代表基金',
      dataIndex: 'current_best_fund_code',
      key: 'current_best_fund_code',
      width: 120,
      render: (code: string, record: FundManagerInfo) => (
        <Space>
          <Button
            type="link"
            size="small"
            onClick={() => window.open(`http://fund.eastmoney.com/${code}.html`, '_blank')}
          >
            {code}
          </Button>
          {record.best_fund_is_4433 && (
            <Tag color="green">4433</Tag>
          )}
        </Space>
      ),
    },
  ];

  return (
    <div className="investool-container">
      <Card>
        <Title level={2}>基金经理筛选</Title>
        
        <Alert
          message="筛选说明"
          description="根据基金经理的从业年限、年化回报、管理规模等指标筛选优秀的基金经理。"
          type="info"
          style={{ marginBottom: 16 }}
        />

        <Card title="筛选条件" style={{ marginBottom: 16 }}>
          <Row gutter={16}>
            <Col xs={24} sm={6}>
              <div style={{ marginBottom: 8 }}>
                <label>最低从业年限</label>
                <InputNumber
                  min={0}
                  max={50}
                  value={params.min_working_years}
                  onChange={(value) => handleFilterChange('min_working_years', value)}
                  style={{ width: '100%' }}
                />
              </div>
            </Col>
            <Col xs={24} sm={6}>
              <div style={{ marginBottom: 8 }}>
                <label>最低年化回报(%)</label>
                <InputNumber
                  min={0}
                  max={100}
                  step={0.1}
                  value={params.min_yieldse}
                  onChange={(value) => handleFilterChange('min_yieldse', value)}
                  style={{ width: '100%' }}
                />
              </div>
            </Col>
            <Col xs={24} sm={6}>
              <div style={{ marginBottom: 8 }}>
                <label>最大管理基金数</label>
                <InputNumber
                  min={1}
                  max={50}
                  value={params.max_current_fund_count}
                  onChange={(value) => handleFilterChange('max_current_fund_count', value)}
                  style={{ width: '100%' }}
                />
              </div>
            </Col>
            <Col xs={24} sm={6}>
              <div style={{ marginBottom: 8 }}>
                <label>最小管理规模(亿)</label>
                <InputNumber
                  min={0}
                  step={10}
                  value={params.min_scale}
                  onChange={(value) => handleFilterChange('min_scale', value)}
                  style={{ width: '100%' }}
                />
              </div>
            </Col>
          </Row>
          <Row gutter={16}>
            <Col xs={24} sm={8}>
              <div style={{ marginBottom: 8 }}>
                <label>排序方式</label>
                <Select
                  value={params.sort}
                  onChange={(value) => handleFilterChange('sort', value)}
                  style={{ width: '100%' }}
                >
                  <Option value="yieldse">年化回报</Option>
                  <Option value="scale">管理规模</Option>
                  <Option value="score">综合评分</Option>
                  <Option value="an">获奖数量</Option>
                  <Option value="fc">基金数量</Option>
                  <Option value="cbr">当前最佳回报</Option>
                  <Option value="wbr">任职最佳回报</Option>
                </Select>
              </div>
            </Col>
            <Col xs={24} sm={8}>
              <div style={{ marginBottom: 8 }}>
                <label>基金类型</label>
                <Select
                  value={params.fund_type}
                  onChange={(value) => handleFilterChange('fund_type', value)}
                  style={{ width: '100%' }}
                  allowClear
                >
                  <Option value="">全部类型</Option>
                  <Option value="股票型">股票型</Option>
                  <Option value="混合型">混合型</Option>
                  <Option value="债券型">债券型</Option>
                  <Option value="指数型">指数型</Option>
                </Select>
              </div>
            </Col>
            <Col xs={24} sm={8}>
              <div style={{ marginBottom: 8 }}>
                <label>基金经理姓名</label>
                <Input
                  value={params.name}
                  onChange={(e) => handleFilterChange('name', e.target.value)}
                  style={{ width: '100%' }}
                  placeholder="输入姓名搜索"
                />
              </div>
            </Col>
          </Row>
        </Card>

        <Row gutter={16} style={{ marginBottom: 16 }}>
          <Col xs={24} sm={8}>
            <Statistic
              title="筛选结果"
              value={managers.length}
            />
          </Col>
          <Col xs={24} sm={8}>
            <Statistic
              title="平均从业年限"
              value={managers.length > 0 ? (managers.reduce((sum, m) => sum + m.working_years, 0) / managers.length).toFixed(1) : 0}
              suffix="年"
            />
          </Col>
          <Col xs={24} sm={8}>
            <Statistic
              title="平均年化回报"
              value={managers.length > 0 ? (managers.reduce((sum, m) => sum + m.yieldse, 0) / managers.length).toFixed(1) : 0}
              suffix="%"
            />
          </Col>
        </Row>

        <Table
          columns={columns}
          dataSource={managers}
          loading={loading}
          rowKey="id"
          pagination={{
            pageSize: 20,
            showSizeChanger: true,
            showQuickJumper: true,
            showTotal: (total, range) => 
              `第 ${range[0]}-${range[1]} 条，共 ${total} 条`,
          }}
          scroll={{ x: 800 }}
        />
      </Card>
    </div>
  );
};

export default FundManagers;
