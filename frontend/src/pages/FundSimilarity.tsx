import React, { useState } from 'react';
import { Card, Typography, Alert, Form, Input, Button } from 'antd';
import { SearchOutlined } from '@ant-design/icons';
import apiClient from '../services/api';
import { FundSimilarityParams } from '../types/fund';

const { Title } = Typography;
const { TextArea } = Input;

const FundSimilarity: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<any>(null);

  const handleSubmit = async (values: FundSimilarityParams) => {
    if (!values.codes?.trim()) {
      return;
    }
    
    setLoading(true);
    try {
      const response = await apiClient.getFundSimilarity(values);
      setResult(response);
    } catch (error) {
      console.error('基金相似度检测失败:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="investool-container">
      <Card>
        <Title level={2}>基金持仓相似度检测</Title>
        
        <Alert
          message="检测说明"
          description="输入多个基金代码，系统将分析这些基金的持仓相似度，帮助您了解基金之间的关联性。"
          type="info"
          style={{ marginBottom: 16 }}
        />

        <Card title="输入基金代码" style={{ marginBottom: 16 }}>
          <Form
            layout="vertical"
            onFinish={handleSubmit}
            initialValues={{ codes: '' }}
          >
            <Form.Item
              label="基金代码"
              name="codes"
              rules={[{ required: true, message: '请输入基金代码' }]}
            >
              <TextArea
                rows={4}
                placeholder="请输入需要比较的基金代码，多个代码用空格或换行分隔&#10;例如：&#10;000001&#10;000002&#10;000003"
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
                检测相似度
              </Button>
            </Form.Item>
          </Form>
        </Card>

        {result && (
          <Card title="检测结果">
            <div style={{ textAlign: 'center', padding: '40px 0' }}>
              <Title level={4}>相似度分析结果</Title>
              <p>检测完成，共分析了 {result.fund_count || 0} 只基金</p>
              <p>平均相似度：{result.avg_similarity || 0}%</p>
            </div>
          </Card>
        )}
      </Card>
    </div>
  );
};

export default FundSimilarity;
