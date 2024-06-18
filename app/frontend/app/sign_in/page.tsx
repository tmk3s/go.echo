"use client";

import axios from 'axios';
import ErrorToast from '@/components/ErrorToast'

import { useRouter } from 'next/navigation'
import { useForm } from "react-hook-form";
import { useState, useEffect } from 'react';

export default () => {
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
      router.push('/todos');
    } catch (e) {
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
          <label className="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Your email</label>
          <input
            type="email"
            id="email"
            className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
            placeholder="name@flowbite.com"
            required
            {...register("email")}
          />
        </div>
        <div className="mb-5">
          <label className="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Your password</label>
          <input
            type="password"
            id="password"
            className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
            required
            {...register("password")}
          />
        </div>
        <button type="submit" className="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm w-full sm:w-auto px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800">
          Submit
        </button>
      </form>
    </main>
  );
}
