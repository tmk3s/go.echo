"use client";

import { useState, useEffect } from 'react';
import RootLayout from '@/components/RootLayout';
import DescriptionTable, { Row } from '@/components/ui/DescriptionTable';
import newApiInstance from '../api';

const formatDate = (iso: string | null | undefined) => {
  if (!iso) return null;
  return new Date(iso).toLocaleDateString('ja-JP');
};

const MyPage = () => {
  const [user, setUser] = useState<any>(null);

  useEffect(() => {
    const api = newApiInstance();
    api.get('/api/user').then((res) => setUser(res.data)).catch(console.error);
  }, []);

  return (
    <RootLayout>
      <h1 className='text-3xl font-bold'>マイページ</h1>

      <DescriptionTable className='mt-6'>
        <Row label='姓' value={user?.UserInfo?.last_name} />
        <Row label='名' value={user?.UserInfo?.first_name} />
        <Row label='メールアドレス' value={user?.email} />
        <Row label='誕生日' value={formatDate(user?.UserInfo?.birthday)} />
        <Row label='性別' value={user?.UserInfo?.gender ? '男' : user?.UserInfo ? '女' : null} />
        <Row label='在籍状況' value={user?.UserInfo?.working ? '在籍中' : user?.UserInfo ? '離職済' : null} />
      </DescriptionTable>
    </RootLayout>
  );
};

export default MyPage;
