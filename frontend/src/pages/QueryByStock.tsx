import React, { useState } from 'react';
import { Card, Typography, Alert, Form, Input, Button } from 'antd';
import { SearchOutlined } from '@ant-design/icons';
import apiClient from '../services/api';
import { QueryByStockParams } from '../types/fund';

const { Title } = Typography;
const { TextArea } = Input;

const QueryByStock: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<any>(null);

  const handleSubmit = async (values: QueryByStockParams) => {
    if (!values.keywords?.trim()) {
      return;
    }
    
    setLoading(true);
    try {
      const response = await apiClient.postQueryByStock(values);
      setResult(response);
    } catch (error) {
      console.error('股票选基失败:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="investool-container">
      <Card>
        <Title level={2}>股票选基</Title>
        
        <Alert
          message="功能说明"
          description="输入股票名称或代码，系统将查找持有这些股票的基金，帮助您通过股票反向选择基金。"
          type="info"
          style={{ marginBottom: 16 }}
        />

        <Card title="输入股票信息" style={{ marginBottom: 16 }}>
          <Form
            layout="vertical"
            onFinish={handleSubmit}
            initialValues={{ keywords: '' }}
          >
            <Form.Item
              label="股票名称或代码"
              name="keywords"
              rules={[{ required: true, message: '请输入股票名称或代码' }]}
            >
              <TextArea
                rows={4}
                placeholder="请输入持仓股票名称或代码，多个股票用空格或换行分隔&#10;例如：&#10;贵州茅台&#10;000001&#10;中国平安"
              />
            </Form.Item>
            <Form.Item>
              <Button
                type="primary"
                htmlType="submit"
                icon={<SearchOutlined />}
                loading={loading}
                size="large"
                style={{ width: '100%' }}
              >
                查询基金
              </Button>
            </Form.Item>
          </Form>
        </Card>

        {result && (
          <Card title="查询结果">
            <div style={{ textAlign: 'center', padding: '40px 0' }}>
              <Title level={4}>股票选基结果</Title>
              <p>查询完成，共找到 {result.fund_count || 0} 只相关基金</p>
              <p>匹配股票数量：{result.stock_count || 0} 只</p>
            </div>
          </Card>
        )}
      </Card>
    </div>
  );
};

export default QueryByStock;
