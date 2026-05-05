"use client";

import { useRouter } from 'next/navigation';
import { useForm } from 'react-hook-form';
import RootLayout from '@/components/RootLayout';
import useApi from '@/app/api';

type FormValues = {
  staff_code: string;
  last_name: string;
  first_name: string;
  last_name_kana: string;
  first_name_kana: string;
  email: string;
};

const EmployeeNew = () => {
  const router = useRouter();
  const api = useApi();
  const { register, handleSubmit, formState: { errors } } = useForm<FormValues>();

  const onSubmit = async (data: FormValues) => {
    await api.post('/api/employee', data);
    router.push('/employees');
  };

  const fields: { label: string; name: keyof FormValues; type?: string }[] = [
    { label: 'スタッフコード', name: 'staff_code' },
    { label: '姓', name: 'last_name' },
    { label: '名', name: 'first_name' },
    { label: '姓（カナ）', name: 'last_name_kana' },
    { label: '名（カナ）', name: 'first_name_kana' },
    { label: 'メールアドレス', name: 'email', type: 'email' },
  ];

  return (
    <RootLayout>
      <div className='flex items-center justify-between'>
        <h1 className='text-3xl font-bold'>社員新規作成</h1>
        <button
          onClick={() => router.push('/employees')}
          className='text-gray-700 bg-white border border-gray-300 hover:bg-gray-100 font-medium rounded-lg text-sm px-5 py-2.5 dark:bg-gray-800 dark:text-gray-400 dark:border-gray-600 dark:hover:bg-gray-700'
        >
          一覧に戻る
        </button>
      </div>

      <form onSubmit={handleSubmit(onSubmit)} className='mt-6 max-w-lg space-y-4'>
        {fields.map(({ label, name, type }) => (
          <div key={name}>
            <label className='block mb-1 text-sm font-medium text-gray-700 dark:text-gray-300'>
              {label}
            </label>
            <input
              type={type ?? 'text'}
              {...register(name, { required: true })}
              className='w-full bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white'
            />
            {errors[name] && (
              <p className='mt-1 text-xs text-red-500'>必須項目です</p>
            )}
          </div>
        ))}

        <button
          type='submit'
          className='text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800'
        >
          作成する
        </button>
      </form>
    </RootLayout>
  );
};

export default EmployeeNew;
