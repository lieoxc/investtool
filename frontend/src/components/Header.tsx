import React from 'react';
import { Layout, Menu, Typography } from 'antd';
import { 
  FundOutlined, 
  FilterOutlined, 
  SearchOutlined, 
  TeamOutlined,
  BarChartOutlined,
  InfoCircleOutlined
} from '@ant-design/icons';
import { useNavigate, useLocation } from 'react-router-dom';

const { Header: AntHeader } = Layout;
const { Title } = Typography;

const Header: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();

  const menuItems = [
    {
      key: '/fund',
      icon: <FundOutlined />,
      label: '4433基金',
      onClick: () => navigate('/fund')
    },
    {
      key: '/fund/filter',
      icon: <FilterOutlined />,
      label: '基金严选',
      onClick: () => navigate('/fund/filter')
    },
    {
      key: '/fund/check',
      icon: <SearchOutlined />,
      label: '基金检测',
      onClick: () => navigate('/fund/check')
    },
    {
      key: '/fund/managers',
      icon: <TeamOutlined />,
      label: '基金经理',
      onClick: () => navigate('/fund/managers')
    },
    {
      key: '/fund/similarity',
      icon: <BarChartOutlined />,
      label: '持仓相似度',
      onClick: () => navigate('/fund/similarity')
    },
    {
      key: '/about',
      icon: <InfoCircleOutlined />,
      label: '关于',
      onClick: () => navigate('/about')
    }
  ];

  return (
    <AntHeader className="investool-header">
      <div className="investool-container">
        <div className="investool-logo">
          <FundOutlined />
          <Title level={3} style={{ color: 'white', margin: 0 }}>
            InvesTool
          </Title>
        </div>
        <Menu
          theme="dark"
          mode="horizontal"
          selectedKeys={[location.pathname]}
          items={menuItems}
          style={{ 
            backgroundColor: 'transparent',
            borderBottom: 'none',
            marginTop: 16
          }}
        />
      </div>
    </AntHeader>
  );
};

export default Header;
