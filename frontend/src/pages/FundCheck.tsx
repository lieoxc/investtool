import React, { useState } from 'react';
import { Card, Typography, Alert, Table, Tag, Button, Row, Col, Statistic } from 'antd';
import { LinkOutlined, CheckCircleOutlined, CloseCircleOutlined } from '@ant-design/icons';
import FundCheckForm from '../components/FundCheckForm';
import apiClient from '../services/api';
import { FundCheckParams, FundCheckResponse, Fund } from '../types/fund';
import { formatCurrency, formatPercentage } from '../utils';

const { Title } = Typography;

const FundCheck: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [checkResult, setCheckResult] = useState<FundCheckResponse | null>(null);

  const handleCheck = async (values: FundCheckParams) => {
    setLoading(true);
    try {
      const response = await apiClient.postFundCheck(values);
      setCheckResult(response);
    } catch (error) {
      console.error('基金检测失败:', error);
    } finally {
      setLoading(false);
    }
  };

  const renderCheckResult = (fund: Fund, param: FundCheckParams) => {
    const p = fund.performance;
    const checkResults = [
      {
        name: `近1年绩效排名前${param.year_1_rank_ratio}%`,
        value: `${p?.year_1_rank_ratio?.toFixed(2) || '--'}%`,
        passed: (p?.year_1_rank_ratio || 0) <= (param.year_1_rank_ratio || 25)
      },
      {
        name: `近2年绩效排名前${param.this_year_235_rank_ratio}%`,
        value: `${p?.year_2_rank_ratio?.toFixed(2) || '--'}%`,
        passed: (p?.year_2_rank_ratio || 0) <= (param.this_year_235_rank_ratio || 25)
      },
      {
        name: `近3年绩效排名前${param.this_year_235_rank_ratio}%`,
        value: `${p?.year_3_rank_ratio?.toFixed(2) || '--'}%`,
        passed: (p?.year_3_rank_ratio || 0) <= (param.this_year_235_rank_ratio || 25)
      },
      {
        name: `近5年绩效排名前${param.this_year_235_rank_ratio}%`,
        value: `${p?.year_5_rank_ratio?.toFixed(2) || '--'}%`,
        passed: (p?.year_5_rank_ratio || 0) <= (param.this_year_235_rank_ratio || 25)
      },
      {
        name: `今年来绩效排名前${param.this_year_235_rank_ratio}%`,
        value: `${p?.this_year_rank_ratio?.toFixed(2) || '--'}%`,
        passed: (p?.this_year_rank_ratio || 0) <= (param.this_year_235_rank_ratio || 25)
      },
      {
        name: `近6个月绩效排名前${param.month_6_rank_ratio}%`,
        value: `${p?.month_6_rank_ratio?.toFixed(2) || '--'}%`,
        passed: (p?.month_6_rank_ratio || 0) <= (param.month_6_rank_ratio || 33.33)
      },
      {
        name: `近3个月绩效排名前${param.month_3_rank_ratio}%`,
        value: `${p?.month_3_rank_ratio?.toFixed(2) || '--'}%`,
        passed: (p?.month_3_rank_ratio || 0) <= (param.month_3_rank_ratio || 33.33)
      },
      {
        name: `基金规模最低${param.min_scale}亿`,
        value: formatCurrency(fund.net_assets_scale),
        passed: (fund.net_assets_scale || 0) / 100000000 >= (param.min_scale || 2)
      },
      {
        name: `基金规模最高${param.max_scale}亿`,
        value: formatCurrency(fund.net_assets_scale),
        passed: (fund.net_assets_scale || 0) / 100000000 <= (param.max_scale || 50)
      },
      {
        name: `基金经理管理该基金不低于${param.min_manager_years}年`,
        value: `${((fund.manager?.manage_days || 0) / 365).toFixed(2)}年`,
        passed: (fund.manager?.manage_days || 0) / 365 >= (param.min_manager_years || 5)
      },
      {
        name: `近1,3,5年波动率平均值不高于${param.max_135_avg_stddev}%`,
        value: formatPercentage(fund.stddev?.avg_135 || 0),
        passed: (fund.stddev?.avg_135 || 0) <= (param.max_135_avg_stddev || 25)
      },
      {
        name: `近1,3,5年夏普比率平均值不低于${param.min_135_avg_sharp}%`,
        value: formatPercentage(fund.sharp?.avg_135 || 0),
        passed: (fund.sharp?.avg_135 || 0) >= (param.min_135_avg_sharp || 1)
      },
      {
        name: `近1,3,5年最大回撤率平均值不高于${param.max_135_avg_retr}%`,
        value: formatPercentage(fund.max_retracement?.avg_135 || 0),
        passed: (fund.max_retracement?.avg_135 || 0) <= (param.max_135_avg_retr || 25)
      }
    ];

    const passedCount = checkResults.filter(r => r.passed).length;
    const totalCount = checkResults.length;

    return (
      <Card key={fund.code} style={{ marginBottom: 16 }}>
        <Row gutter={16} style={{ marginBottom: 16 }}>
          <Col xs={24} sm={12}>
            <Title level={4}>
              <Button
                type="link"
                onClick={() => window.open(`http://fund.eastmoney.com/${fund.code}.html`, '_blank')}
                icon={<LinkOutlined />}
              >
                {fund.name} ({fund.code})
              </Button>
            </Title>
          </Col>
          <Col xs={24} sm={12}>
            <Statistic
              title="检测通过率"
              value={passedCount}
              suffix={`/ ${totalCount}`}
              valueStyle={{ color: passedCount === totalCount ? '#52c41a' : '#faad14' }}
            />
          </Col>
        </Row>

        <Table
          dataSource={checkResults}
          columns={[
            {
              title: '检测指标',
              dataIndex: 'name',
              key: 'name',
              width: '40%'
            },
            {
              title: '实际值',
              dataIndex: 'value',
              key: 'value',
              width: '30%'
            },
            {
              title: '结果',
              key: 'result',
              width: '30%',
              render: (_, record) => (
                <Tag color={record.passed ? 'green' : 'red'} icon={record.passed ? <CheckCircleOutlined /> : <CloseCircleOutlined />}>
                  {record.passed ? '通过' : '不通过'}
                </Tag>
              )
            }
          ]}
          pagination={false}
          size="small"
        />

        {fund.manager && (
          <div style={{ marginTop: 16 }}>
            <Title level={5}>基金经理信息</Title>
            <Row gutter={16}>
              <Col xs={24} sm={8}>
                <Statistic title="基金经理" value={fund.manager.name} />
              </Col>
              <Col xs={24} sm={8}>
                <Statistic title="管理年限" value={((fund.manager.manage_days || 0) / 365).toFixed(2)} suffix="年" />
              </Col>
              <Col xs={24} sm={8}>
                <Statistic title="任职回报" value={fund.manager.manage_repay?.toFixed(2)} suffix="%" />
              </Col>
            </Row>
          </div>
        )}
      </Card>
    );
  };

  return (
    <div className="investool-container">
      <Card>
        <Title level={2}>基金检测</Title>
        
        <Alert
          message="检测说明"
          description="输入基金代码，系统将根据4433法则和风险指标对基金进行全面检测，帮助您评估基金的投资价值。"
          type="info"
          style={{ marginBottom: 16 }}
        />

        <FundCheckForm
          onSubmit={handleCheck}
          loading={loading}
        />

        {checkResult && (
          <div style={{ marginTop: 24 }}>
            <Title level={3}>检测结果</Title>
            {checkResult.funds?.map(fund => renderCheckResult(fund, checkResult.param))}
          </div>
        )}
      </Card>
    </div>
  );
};

export default FundCheck;
