import React from 'react';
import { Card, Typography, Alert, Row, Col, Statistic } from 'antd';

const { Title, Paragraph, Text } = Typography;

const About: React.FC = () => {
  return (
    <div className="investool-container">
      <Card>
        <Row gutter={24}>
          <Col xs={24} lg={16}>
            <Title level={2}>关于 InvesTool</Title>
            
            <Alert
              message="免责声明"
              description="以下所有数据与信息仅供参考，不构成投资建议。投资有风险，入市需谨慎。"
              type="warning"
              style={{ marginBottom: 24 }}
            />

            <Title level={3}>项目介绍</Title>
            <Paragraph>
              InvesTool 是一个用于股票基金投资分析的辅助工具网站，主要提供公司财报分析、股票基本面分析、
              基本面选股、基金经理排行榜、基金4433法则筛选、基金排行榜等功能。
            </Paragraph>

            <Title level={3}>主要功能</Title>
            <Row gutter={16}>
              <Col xs={24} sm={12}>
                <Card size="small" style={{ marginBottom: 16 }}>
                  <Statistic title="4433基金筛选" value="智能筛选" />
                  <Text type="secondary">根据4433法则自动筛选优质基金</Text>
                </Card>
              </Col>
              <Col xs={24} sm={12}>
                <Card size="small" style={{ marginBottom: 16 }}>
                  <Statistic title="基金检测" value="全面检测" />
                  <Text type="secondary">多维度检测基金投资价值</Text>
                </Card>
              </Col>
              <Col xs={24} sm={12}>
                <Card size="small" style={{ marginBottom: 16 }}>
                  <Statistic title="基金经理分析" value="专业分析" />
                  <Text type="secondary">基金经理业绩和能力分析</Text>
                </Card>
              </Col>
              <Col xs={24} sm={12}>
                <Card size="small" style={{ marginBottom: 16 }}>
                  <Statistic title="持仓相似度" value="关联分析" />
                  <Text type="secondary">分析基金持仓相似度</Text>
                </Card>
              </Col>
            </Row>

            <Title level={3}>4433法则</Title>
            <Paragraph>
              由台大财务金融学系邱显比教授提出的选基法则：
            </Paragraph>
            <ul>
              <li><strong>4</strong>: 最近1年收益率排名在同类型基金前1/4</li>
              <li><strong>4</strong>: 最近2年、3年、5年及今年来收益率排名均在同类型基金前1/4</li>
              <li><strong>3</strong>: 最近6个月收益率排名在同类型基金前1/3</li>
              <li><strong>3</strong>: 最近3个月收益率排名在同类型基金前1/3</li>
            </ul>

            <Title level={3}>技术栈</Title>
            <Row gutter={16}>
              <Col xs={24} sm={8}>
                <Card size="small">
                  <Title level={5}>后端</Title>
                  <Text>Go + Gin + GORM</Text>
                </Card>
              </Col>
              <Col xs={24} sm={8}>
                <Card size="small">
                  <Title level={5}>前端</Title>
                  <Text>React + TypeScript + Ant Design</Text>
                </Card>
              </Col>
              <Col xs={24} sm={8}>
                <Card size="small">
                  <Title level={5}>数据源</Title>
                  <Text>东方财富、新浪财经等</Text>
                </Card>
              </Col>
            </Row>
          </Col>
{/*           
          <Col xs={24} lg={8}>
            <Card>
              <Title level={4}>项目信息</Title>
              <Divider />
              
              <div style={{ marginBottom: 16 }}>
                <Text strong>版本：</Text>
                <Text>v1.0.0</Text>
              </div>
              
              <div style={{ marginBottom: 16 }}>
                <Text strong>作者：</Text>
                <Text>axiaoxin</Text>
              </div>
              
              <div style={{ marginBottom: 16 }}>
                <Text strong>开源协议：</Text>
                <Text>MIT</Text>
              </div>
              
              <div style={{ marginBottom: 16 }}>
                <Text strong>GitHub：</Text>
                <br />
                <Button 
                  type="link" 
                  icon={<GithubOutlined />}
                  onClick={() => window.open('https://github.com/axiaoxin-com/investool', '_blank')}
                >
                  axiaoxin-com/investool
                </Button>
              </div>
              
              <Divider />
              
              <div style={{ textAlign: 'center' }}>
                <Text type="secondary">
                  Made with <HeartOutlined style={{ color: '#ff4d4f' }} /> by axiaoxin
                </Text>
                <br />
                <Text type="secondary">
                  © 2021-{new Date().getFullYear()} All rights reserved.
                </Text>
              </div>
            </Card>
          </Col> */}
        </Row>
      </Card>
    </div>
  );
};

export default About;
