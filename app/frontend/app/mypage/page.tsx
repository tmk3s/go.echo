"use client";

import { useState, useEffect } from 'react';
import RootLayout from '@/components/RootLayout';
import newApiInstance from '../api';

const Row = ({ label, value }: { label: string; value: string | null | undefined }) => (
  <tr className='border-b dark:border-gray-700'>
    <th className='px-6 py-3 w-48 bg-gray-50 dark:bg-gray-700 font-medium text-gray-700 dark:text-gray-300 text-sm'>
      {label}
    </th>
    <td className='px-6 py-3 text-sm text-gray-900 dark:text-white'>
      {value ?? '—'}
    </td>
  </tr>
);

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

      <div className='mt-6 bg-white rounded-lg shadow dark:bg-gray-800 overflow-hidden'>
        <table className='w-full'>
          <tbody>
            <Row label='姓' value={user?.UserInfo?.last_name} />
            <Row label='名' value={user?.UserInfo?.first_name} />
            <Row label='メールアドレス' value={user?.email} />
            <Row label='誕生日' value={formatDate(user?.UserInfo?.birthday)} />
            <Row label='性別' value={user?.UserInfo?.gender ? '男' : user?.UserInfo ? '女' : null} />
            <Row label='在籍状況' value={user?.UserInfo?.working ? '在籍中' : user?.UserInfo ? '離職済' : null} />
          </tbody>
        </table>
      </div>
    </RootLayout>
  );
};

export default MyPage;
