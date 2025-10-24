import React from 'react';
import { Routes, Route } from 'react-router-dom';
import { Layout } from 'antd';
import Header from './components/Header';
import FundIndex from './pages/FundIndex';
import FundFilter from './pages/FundFilter';
import FundCheck from './pages/FundCheck';
import FundManagers from './pages/FundManagers';
import FundSimilarity from './pages/FundSimilarity';
import QueryByStock from './pages/QueryByStock';
import About from './pages/About';

const { Content } = Layout;

const App: React.FC = () => {
  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header />
      <Content>
        <Routes>
          <Route path="/" element={<FundIndex />} />
          <Route path="/fund" element={<FundIndex />} />
          <Route path="/fund/filter" element={<FundFilter />} />
          <Route path="/fund/check" element={<FundCheck />} />
          <Route path="/fund/managers" element={<FundManagers />} />
          <Route path="/fund/similarity" element={<FundSimilarity />} />
          <Route path="/fund/query_by_stock" element={<QueryByStock />} />
          <Route path="/about" element={<About />} />
        </Routes>
      </Content>
    </Layout>
  );
};

export default App;
