import React from 'react';
import { Table, Tag, Button, Space } from 'antd';
import { LinkOutlined, CopyOutlined } from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import { Fund } from '../types/fund';
import { formatCurrency, formatPercentage, getFundTypeColor, calculateFundScore } from '../utils';

interface FundTableProps {
  data: Fund[];
  loading?: boolean;
  showScore?: boolean;
  onRowClick?: (record: Fund) => void;
}

const FundTable: React.FC<FundTableProps> = ({ 
  data, 
  loading = false, 
  showScore = true,
  onRowClick 
}) => {
  const handleCopyCode = (code: string) => {
    navigator.clipboard.writeText(code);
    // 这里可以添加成功提示
  };

  const columns: ColumnsType<Fund> = [
    {
      title: '基金代码',
      dataIndex: 'code',
      key: 'code',
      width: 100,
      render: (code: string) => (
        <Space>
          <Button 
            type="link" 
            size="small"
            onClick={() => handleCopyCode(code)}
            icon={<CopyOutlined />}
          />
          {code}
        </Space>
      ),
    },
    {
      title: '基金名称',
      dataIndex: 'name',
      key: 'name',
      width: 200,
      render: (name: string, record: Fund) => (
        <Button 
          type="link" 
          onClick={() => window.open(`http://fund.eastmoney.com/${record.code}.html`, '_blank')}
          icon={<LinkOutlined />}
        >
          {name}
        </Button>
      ),
    },
    {
      title: '基金类型',
      dataIndex: 'type',
      key: 'type',
      width: 100,
      render: (type: string) => (
        <Tag color={getFundTypeColor(type)}>
          {type}
        </Tag>
      ),
    },
    {
      title: '基金规模',
      dataIndex: 'net_assets_scale',
      key: 'net_assets_scale',
      width: 120,
      render: (scale: number) => formatCurrency(scale),
      sorter: (a, b) => a.net_assets_scale - b.net_assets_scale,
    },
    {
      title: '基金经理',
      key: 'manager',
      width: 150,
      render: (_, record: Fund) => (
        <div>
          <div>{record.manager?.name || '--'}</div>
          <div style={{ fontSize: '12px', color: '#666' }}>
            管理{(record.manager?.manage_days || 0) / 365}年
          </div>
        </div>
      ),
    },
    {
      title: '近1年排名',
      dataIndex: ['performance', 'year_1_rank_ratio'],
      key: 'year_1_rank',
      width: 100,
      render: (ratio: number) => (
        <span style={{ color: ratio <= 25 ? '#52c41a' : '#ff4d4f' }}>
          {formatPercentage(ratio)}
        </span>
      ),
      sorter: (a, b) => (a.performance?.year_1_rank_ratio || 0) - (b.performance?.year_1_rank_ratio || 0),
    },
    {
      title: '近3年排名',
      dataIndex: ['performance', 'year_3_rank_ratio'],
      key: 'year_3_rank',
      width: 100,
      render: (ratio: number) => (
        <span style={{ color: ratio <= 25 ? '#52c41a' : '#ff4d4f' }}>
          {formatPercentage(ratio)}
        </span>
      ),
      sorter: (a, b) => (a.performance?.year_3_rank_ratio || 0) - (b.performance?.year_3_rank_ratio || 0),
    },
    {
      title: '波动率',
      dataIndex: ['stddev', 'avg_135'],
      key: 'stddev',
      width: 100,
      render: (stddev: number) => formatPercentage(stddev),
      sorter: (a, b) => (a.stddev?.avg_135 || 0) - (b.stddev?.avg_135 || 0),
    },
    {
      title: '夏普比率',
      dataIndex: ['sharp', 'avg_135'],
      key: 'sharp',
      width: 100,
      render: (sharp: number) => formatPercentage(sharp),
      sorter: (a, b) => (a.sharp?.avg_135 || 0) - (b.sharp?.avg_135 || 0),
    },
  ];

  // 如果显示评分，添加评分列
  if (showScore) {
    columns.push({
      title: '综合评分',
      key: 'score',
      width: 100,
      render: (_, record: Fund) => {
        const score = calculateFundScore(record);
        return (
          <Tag color={score >= 80 ? 'green' : score >= 60 ? 'orange' : 'red'}>
            {score}分
          </Tag>
        );
      },
      sorter: (a, b) => calculateFundScore(a) - calculateFundScore(b),
    });
  }

  return (
    <Table
      columns={columns}
      dataSource={data}
      loading={loading}
      rowKey="code"
      pagination={{
        pageSize: 20,
        showSizeChanger: true,
        showQuickJumper: true,
        showTotal: (total, range) => 
          `第 ${range[0]}-${range[1]} 条，共 ${total} 条`,
      }}
      scroll={{ x: 1000 }}
      onRow={(record) => ({
        onClick: () => onRowClick?.(record),
        style: { cursor: 'pointer' },
      })}
    />
  );
};

export default FundTable;
