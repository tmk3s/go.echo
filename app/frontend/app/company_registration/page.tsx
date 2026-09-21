"use client";

import axios from 'axios';
import ErrorToast from '@/components/ErrorToast'
import Button from '@/components/ui/Button';
import TextInput, { Label } from '@/components/ui/TextInput';

import { useRouter } from 'next/navigation'
import { useForm } from "react-hook-form";
import { useState } from 'react';

type CompanyRegistrationForm = {
  company_name: string;
  last_name: string;
  first_name: string;
  email: string;
  password: string;
};

const CompanyRegistration = () => {
  const router = useRouter();
  const [openErrorToast, setOpenErrorToast] = useState(false);
  const [errorMessage, setErrorMessage] = useState('');

  const {
    register,
    handleSubmit,
    formState: { isSubmitting }
  } = useForm<CompanyRegistrationForm>({
    defaultValues: {
      company_name: '',
      last_name: '',
      first_name: '',
      email: '',
      password: '',
    }
  });

  // 会社と最初のユーザーを同時に作成する。作成後はそのままログインできる。
  const registerCompany = async (data: CompanyRegistrationForm) => {
    try {
      axios.defaults.baseURL = 'http://localhost:1323';
      await axios.post('/companies', data, { withCredentials: true });
      router.push('/sign_in');
    } catch (e: any) {
      setErrorMessage(e.response?.data?.message ?? 'ネットワークエラーが発生しました');
      setOpenErrorToast(true);
      console.error(e);
    }
  }

  const clostToast = () => setOpenErrorToast(false);

  return (
    <main className="flex min-h-screen flex-col items-center justify-between p-24">
      {(openErrorToast && errorMessage) && (<ErrorToast message={errorMessage} clostToast={clostToast} />)}
      <form
        className="max-w-sm mx-auto"
        onSubmit={handleSubmit(registerCompany)}
      >
        <h1 className="mb-6 text-xl font-semibold text-body">会社登録</h1>

        <div className="mb-5">
          <Label>会社名</Label>
          <TextInput
            type="text"
            id="company_name"
            placeholder="株式会社サンプル"
            required
            {...register("company_name")}
          />
        </div>

        <div className="mb-5 flex gap-3">
          <div className="flex-1">
            <Label>姓</Label>
            <TextInput type="text" id="last_name" placeholder="山田" required {...register("last_name")} />
          </div>
          <div className="flex-1">
            <Label>名</Label>
            <TextInput type="text" id="first_name" placeholder="太郎" required {...register("first_name")} />
          </div>
        </div>

        <div className="mb-5">
          <Label>メールアドレス</Label>
          <TextInput
            type="email"
            id="email"
            placeholder="name@example.com"
            required
            {...register("email")}
          />
        </div>

        <div className="mb-5">
          <Label>パスワード</Label>
          <TextInput
            type="password"
            id="password"
            required
            minLength={8}
            {...register("password")}
          />
          <p className="mt-1 text-xs text-muted">8文字以上で入力してください</p>
        </div>

        <Button type="submit" disabled={isSubmitting} className="w-full sm:w-auto">
          登録
        </Button>
      </form>
    </main>
  );
}

export default CompanyRegistration
