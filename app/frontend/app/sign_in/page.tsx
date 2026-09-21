"use client";

import axios from 'axios';
import ErrorToast from '@/components/ErrorToast'
import Button from '@/components/ui/Button';
import TextInput, { Label } from '@/components/ui/TextInput';

import { useRouter } from 'next/navigation'
import { useForm } from "react-hook-form";
import { useState, useEffect } from 'react';

const SignIn = () => {
  const router = useRouter();
  const [openErrorToast, setOpenErrorToast] = useState(true);
  const [unauthorizedError, setUnauthorizedError] = useState('');

  const signIn = async (email: string, password: string) => {
    try {
      console.log('exec sign in');
      axios.defaults.baseURL = 'http://localhost:1323';
      // https://sheltie-garage.xyz/tech/2023/07/cookie%E3%81%8C%E3%81%AA%E3%81%8B%E3%81%AA%E3%81%8B%E3%81%A7%E3%81%8D%E3%81%9A%E3%81%AB%E3%83%8F%E3%83%9E%E3%81%A3%E3%81%9F%E8%A9%B1/
      const response = await axios.post('/sign_in', { email: email, password: password }, { withCredentials: true });
      console.log(response);
      router.push('/mypage');
    } catch (e: any) {
      setUnauthorizedError(e.response?.data?.message ?? 'ネットワークエラーが発生しました');
      setOpenErrorToast(true)
      console.error(e);
    }
  }

  const clostToast = () => {
    window.localStorage.removeItem("unauthorizedError");
    setOpenErrorToast(false)
  }

  const {
    register,
    handleSubmit,
    watch,
    formState: { errors }
  } = useForm({
    defaultValues: {
      email: '',
      password: '',
    }
  });

  useEffect(() => {
    // https://sentry.io/answers/referenceerror-localstorage-is-not-defined-in-next-js/
    const message: string = window.localStorage.getItem("unauthorizedError") || ''
    setUnauthorizedError(message);
  }, [])

  return (
    <main className="flex min-h-screen flex-col items-center justify-between p-24">
      { (openErrorToast && unauthorizedError) && (<ErrorToast message={unauthorizedError} clostToast={clostToast}/>) }
      <form
        className="max-w-sm mx-auto"
        onSubmit={handleSubmit((data) => {
          signIn(data.email, data.password)
        })}
      >
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
          <TextInput type="password" id="password" required {...register("password")} />
        </div>
        <Button type="submit" className="w-full sm:w-auto">ログイン</Button>
      </form>
      <div>
        <p>サンプル</p>
        <ul>
          <li>* admin@example.com / password1</li>
          <li>* user1@example.com / password2</li>
          <li>* user2@example.com / password3</li>
          <li>* other@example.com / password4 (別会社)</li>
        </ul>
      </div>
    </main>
  );
}

export default SignIn
