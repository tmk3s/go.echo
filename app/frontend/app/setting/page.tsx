"use client";

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import RootLayout from '@/components/RootLayout';
import Button from '@/components/ui/Button';
import DescriptionTable, { Row } from '@/components/ui/DescriptionTable';
import useApi from '@/app/api';

type Company = {
  ID: number;
  name: string;
};

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
        <p className='text-muted'>読み込み中...</p>
      </RootLayout>
    );
  }

  return (
    <RootLayout>
      <div className='flex items-center justify-between'>
        <h1 className='text-3xl font-bold'>会社情報</h1>
        <Button onClick={() => router.push('/setting/edit')}>編集</Button>
      </div>

      <DescriptionTable className='mt-6'>
        <Row label='会社名' value={company.name} />
      </DescriptionTable>
    </RootLayout>
  );
};

export default SettingPage;
