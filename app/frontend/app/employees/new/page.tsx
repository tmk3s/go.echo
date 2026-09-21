"use client";

import { useRouter } from 'next/navigation';
import { useForm } from 'react-hook-form';
import RootLayout from '@/components/RootLayout';
import Button from '@/components/ui/Button';
import TextInput, { Label } from '@/components/ui/TextInput';
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
        <Button variant='secondary' onClick={() => router.push('/employees')}>
          一覧に戻る
        </Button>
      </div>

      <form onSubmit={handleSubmit(onSubmit)} className='mt-6 max-w-lg space-y-4'>
        {fields.map(({ label, name, type }) => (
          <div key={name}>
            <Label>{label}</Label>
            <TextInput type={type ?? 'text'} {...register(name, { required: true })} />
            {errors[name] && (
              <p className='mt-1 text-xs text-danger'>必須項目です</p>
            )}
          </div>
        ))}

        <Button type='submit'>作成する</Button>
      </form>
    </RootLayout>
  );
};

export default EmployeeNew;
