"use client";

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import RootLayout from '@/components/RootLayout';
import useApi from '@/app/api';

type Company = {
  ID: number;
  name: string;
};

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

const SettingPage = () => {
  const router = useRouter();
  const api = useApi();
  const [company, setCompany] = useState<Company | null>(null);

  useEffect(() => {
    api.get('/api/company').then((res) => setCompany(res.data));
  }, []);

  if (!company) {
    return (
      <RootLayout>
        <p className='text-gray-500'>読み込み中...</p>
      </RootLayout>
    );
  }

  return (
    <RootLayout>
      <div className='flex items-center justify-between'>
        <h1 className='text-3xl font-bold'>会社情報</h1>
        <button
          onClick={() => router.push('/setting/edit')}
          className='text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 dark:bg-blue-600 dark:hover:bg-blue-700'
        >
          編集
        </button>
      </div>

      <div className='mt-6 bg-white rounded-lg shadow dark:bg-gray-800 overflow-hidden'>
        <table className='w-full'>
          <tbody>
            <Row label='会社名' value={company.name} />
          </tbody>
        </table>
      </div>
    </RootLayout>
  );
};

export default SettingPage;
