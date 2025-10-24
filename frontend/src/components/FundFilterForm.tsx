import React from 'react';
import { Form, InputNumber, Select, Button, Card, Row, Col, Space } from 'antd';
import { SearchOutlined, ReloadOutlined } from '@ant-design/icons';
import { FundFilterParams } from '../types/fund';

interface FundFilterFormProps {
  onSubmit: (values: FundFilterParams) => void;
  loading?: boolean;
  fundTypes?: string[];
  initialValues?: FundFilterParams;
}

const FundFilterForm: React.FC<FundFilterFormProps> = ({
  onSubmit,
  loading = false,
  fundTypes = [],
  initialValues = {}
}) => {
  const [form] = Form.useForm();

  const handleSubmit = (values: FundFilterParams) => {
    onSubmit(values);
  };

  const handleReset = () => {
    form.resetFields();
    onSubmit({});
  };

  return (
    <Card title="基金筛选条件" style={{ marginBottom: 16 }}>
      <Form
        form={form}
        layout="vertical"
        onFinish={handleSubmit}
        initialValues={{
          year_1_rank_ratio: 25.0,
          this_year_235_rank_ratio: 25.0,
          month_6_rank_ratio: 33.33,
          month_3_rank_ratio: 33.33,
          min_scale: 2,
          max_scale: 50,
          min_estab_years: 5,
          min_manager_years: 5,
          max_135_avg_stddev: 25.0,
          min_135_avg_sharp: 1.0,
          max_135_avg_retr: 25.0,
          ...initialValues
        }}
      >
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
              label="该基金成立最低年限"
              name="min_estab_years"
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
              label="指定基金类型"
              name="types"
            >
              <Select
                mode="multiple"
                placeholder="选择基金类型"
                style={{ width: '100%' }}
                options={fundTypes.map(type => ({ label: type, value: type }))}
              />
            </Form.Item>
          </Col>
        </Row>

        <Row gutter={16}>
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
        </Row>

        <Form.Item>
          <Space>
            <Button
              type="primary"
              htmlType="submit"
              icon={<SearchOutlined />}
              loading={loading}
            >
              筛选
            </Button>
            <Button
              icon={<ReloadOutlined />}
              onClick={handleReset}
            >
              重置
            </Button>
          </Space>
        </Form.Item>
      </Form>
    </Card>
  );
};

export default FundFilterForm;
