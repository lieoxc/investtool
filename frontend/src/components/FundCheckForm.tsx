import React, { useState } from 'react';
import { Form, InputNumber, Button, Card, Row, Col, Switch, message, Input } from 'antd';
import { SearchOutlined } from '@ant-design/icons';
import { FundCheckParams } from '../types/fund';

interface FundCheckFormProps {
  onSubmit: (values: FundCheckParams) => void;
  loading?: boolean;
}

const FundCheckForm: React.FC<FundCheckFormProps> = ({
  onSubmit,
  loading = false
}) => {
  const [form] = Form.useForm();
  const [checkStocks, setCheckStocks] = useState(false);

  const handleSubmit = (values: FundCheckParams) => {
    if (!values.fundcode?.trim()) {
      message.error('请填写基金代码');
      return;
    }
    onSubmit({ ...values, check_stocks: checkStocks });
  };

  return (
    <Card title="基金检测" style={{ marginBottom: 16 }}>
      <Form
        form={form}
        layout="vertical"
        onFinish={handleSubmit}
        initialValues={{
          min_scale: 2,
          max_scale: 50,
          min_manager_years: 5,
          year_1_rank_ratio: 25.0,
          this_year_235_rank_ratio: 25.0,
          month_6_rank_ratio: 33.33,
          month_3_rank_ratio: 33.33,
          max_135_avg_stddev: 25.0,
          min_135_avg_sharp: 1.0,
          max_135_avg_retr: 25.0,
        }}
      >
        <Row gutter={16}>
          <Col xs={24}>
            <Form.Item
              label="基金代码"
              name="fundcode"
              rules={[{ required: true, message: '请输入基金代码' }]}
            >
              <Input
                placeholder="请输入基金代码，多个代码用空格分隔"
                style={{ width: '100%' }}
              />
            </Form.Item>
          </Col>
        </Row>

        <Row gutter={16}>
          <Col xs={24} sm={12} md={8}>
            <Form.Item
              label="近1年绩效排名前百分之"
              name="year_1_rank_ratio"
            >
              <InputNumber
                min={0}
                max={100}
                step={0.01}
                style={{ width: '100%' }}
                placeholder="25.00"
              />
            </Form.Item>
          </Col>
          <Col xs={24} sm={12} md={8}>
            <Form.Item
              label="近2,3,5年及今年来绩效排名前百分之"
              name="this_year_235_rank_ratio"
            >
              <InputNumber
                min={0}
                max={100}
                step={0.01}
                style={{ width: '100%' }}
                placeholder="25.00"
              />
            </Form.Item>
          </Col>
          <Col xs={24} sm={12} md={8}>
            <Form.Item
              label="近6个月绩效排名前百分之"
              name="month_6_rank_ratio"
            >
              <InputNumber
                min={0}
                max={100}
                step={0.01}
                style={{ width: '100%' }}
                placeholder="33.33"
              />
            </Form.Item>
          </Col>
        </Row>

        <Row gutter={16}>
          <Col xs={24} sm={12} md={8}>
            <Form.Item
              label="近3个月绩效排名前百分之"
              name="month_3_rank_ratio"
            >
              <InputNumber
                min={0}
                max={100}
                step={0.01}
                style={{ width: '100%' }}
                placeholder="33.33"
              />
            </Form.Item>
          </Col>
          <Col xs={24} sm={12} md={8}>
            <Form.Item
              label="基金规模最小值（亿）"
              name="min_scale"
            >
              <InputNumber
                min={0}
                step={1}
                style={{ width: '100%' }}
                placeholder="2"
              />
            </Form.Item>
          </Col>
          <Col xs={24} sm={12} md={8}>
            <Form.Item
              label="基金规模最大值（亿）"
              name="max_scale"
            >
              <InputNumber
                min={0}
                step={1}
                style={{ width: '100%' }}
                placeholder="50"
              />
            </Form.Item>
          </Col>
        </Row>

        <Row gutter={16}>
          <Col xs={24} sm={12} md={8}>
            <Form.Item
              label="基金经理管理该基金最低年限"
              name="min_manager_years"
            >
              <InputNumber
                min={0}
                step={1}
                style={{ width: '100%' }}
                placeholder="5"
              />
            </Form.Item>
          </Col>
          <Col xs={24} sm={12} md={8}>
            <Form.Item
              label="近1,3,5年波动率平均值的最大值(%)"
              name="max_135_avg_stddev"
            >
              <InputNumber
                min={0}
                step={1.0}
                style={{ width: '100%' }}
                placeholder="25.0"
              />
            </Form.Item>
          </Col>
          <Col xs={24} sm={12} md={8}>
            <Form.Item
              label="近1,3,5年夏普比率平均值的最小值(%)"
              name="min_135_avg_sharp"
            >
              <InputNumber
                min={0}
                step={1.0}
                style={{ width: '100%' }}
                placeholder="1.0"
              />
            </Form.Item>
          </Col>
        </Row>

        <Row gutter={16}>
          <Col xs={24} sm={12} md={8}>
            <Form.Item
              label="近1,3,5年最大回撤率平均值的最大值(%)"
              name="max_135_avg_retr"
            >
              <InputNumber
                min={0}
                step={1.0}
                style={{ width: '100%' }}
                placeholder="25.0"
              />
            </Form.Item>
          </Col>
          <Col xs={24} sm={12} md={8}>
            <Form.Item
              label="检测持仓个股"
              name="check_stocks"
            >
              <Switch
                checked={checkStocks}
                onChange={setCheckStocks}
                checkedChildren="开启"
                unCheckedChildren="关闭"
              />
            </Form.Item>
          </Col>
        </Row>

        <Form.Item>
          <Button
            type="primary"
            htmlType="submit"
            icon={<SearchOutlined />}
            loading={loading}
            size="large"
            style={{ width: '100%' }}
          >
            检测基金
          </Button>
        </Form.Item>
      </Form>
    </Card>
  );
};

export default FundCheckForm;
