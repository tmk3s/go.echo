"use client";

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useForm } from 'react-hook-form';
import RootLayout from '@/components/RootLayout';
import useApi from '@/app/api';

type FormValues = {
  name: string;
};

const inputClass = 'w-full bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white';

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
            <button
              type='submit'
              className='text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 dark:bg-blue-600 dark:hover:bg-blue-700'
            >
              保存
            </button>
            <button
              type='button'
              onClick={() => router.push('/setting')}
              className='text-gray-700 bg-white border border-gray-300 hover:bg-gray-100 font-medium rounded-lg text-sm px-5 py-2.5 dark:bg-gray-800 dark:text-gray-400 dark:border-gray-600 dark:hover:bg-gray-700'
            >
              キャンセル
            </button>
          </div>
        </div>

        <div className='mt-6 bg-white rounded-lg shadow dark:bg-gray-800 p-6 max-w-lg'>
          <div>
            <label className='block mb-1 text-sm font-medium text-gray-700 dark:text-gray-300'>
              会社名
            </label>
            <input
              className={inputClass}
              {...register('name', { required: '会社名は必須です' })}
            />
            {errors.name && (
              <p className='mt-1 text-xs text-red-500'>{errors.name.message}</p>
            )}
          </div>
        </div>
      </form>
    </RootLayout>
  );
};

export default SettingEditPage;
