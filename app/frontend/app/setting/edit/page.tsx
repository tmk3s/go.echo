"use client";

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useForm } from 'react-hook-form';
import RootLayout from '@/components/RootLayout';
import Button from '@/components/ui/Button';
import Card from '@/components/ui/Card';
import TextInput, { Label } from '@/components/ui/TextInput';
import useApi from '@/app/api';

type FormValues = {
  name: string;
};

const SettingEditPage = () => {
  const router = useRouter();
  const api = useApi();
  const { register, handleSubmit, reset, formState: { errors } } = useForm<FormValues>();

  useEffect(() => {
    api.get('/api/company').then((res) => {
      reset({ name: res.data.name });
    });
  }, []);

  const onSubmit = handleSubmit(async (data) => {
    await api.put('/api/company', data);
    router.push('/setting');
  });

  return (
    <RootLayout>
      <form onSubmit={onSubmit}>
        <div className='flex items-center justify-between'>
          <h1 className='text-3xl font-bold'>会社情報編集</h1>
          <div className='flex gap-2'>
            <Button type='submit'>保存</Button>
            <Button variant='secondary' onClick={() => router.push('/setting')}>
              キャンセル
            </Button>
          </div>
        </div>

        <Card className='mt-6 p-6 max-w-lg'>
          <div>
            <Label>会社名</Label>
            <TextInput {...register('name', { required: '会社名は必須です' })} />
            {errors.name && (
              <p className='mt-1 text-xs text-red-500'>{errors.name.message}</p>
            )}
          </div>
        </Card>
      </form>
    </RootLayout>
  );
};

export default SettingEditPage;
