"use client";

import { useRouter } from 'next/navigation'
import { useForm, SubmitHandler } from "react-hook-form";
import newApiInstance from "../api"

// async function signIn(email: string, password: string) {
//   try {
//     const api = newApiInstance();
//     const response = await api.post('/sign_in', { email: email, password: password });
//     console.log(response);
    
//   } catch (error) {
//     console.error(error);
//   }
// }

export default () => {
  const router = useRouter();

  const signIn = async (email: string, password: string) => {
    try {
      const api = newApiInstance();
      const response = await api.post('/sign_in', { email: email, password: password });
      console.log(response);
      router.push('/todos');
    } catch (error) {
      console.error(error);
    }
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
  return (
    <main className="flex min-h-screen flex-col items-center justify-between p-24">
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
        <div className="flex items-start mb-5">
          <div className="flex items-center h-5">
            <input id="remember" type="checkbox" value="" className="w-4 h-4 border border-gray-300 rounded bg-gray-50 focus:ring-3 focus:ring-blue-300 dark:bg-gray-700 dark:border-gray-600 dark:focus:ring-blue-600 dark:ring-offset-gray-800 dark:focus:ring-offset-gray-800" required />
          </div>
          <label className="ms-2 text-sm font-medium text-gray-900 dark:text-gray-300">Remember me</label>
        </div>
        <button type="submit" className="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm w-full sm:w-auto px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800">
          Submit
        </button>
      </form>
    </main>
  );
}
