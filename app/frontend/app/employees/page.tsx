"use client";

import { useEffect, useState, useMemo } from 'react';
import { useRouter } from 'next/navigation';
import RootLayout from '@/components/RootLayout';
import useApi from '@/app/api';

type Employee = {
  ID: number;
  last_name: string;
  first_name: string;
  last_name_kana: string;
  first_name_kana: string;
  email: string;
  staff_code: string;
};

const Employees = () => {
  const router = useRouter();
  const api = useApi();
  const [employees, setEmployees] = useState<Employee[]>([]);
  const [query, setQuery] = useState('');

  useEffect(() => {
    api.get('/api/employees').then((res) => {
      setEmployees(res.data ?? []);
    });
  }, []);

  const filtered = useMemo(() => {
    if (!query) return employees;
    const q = query.toLowerCase();
    return employees.filter((e) =>
      e.staff_code.toLowerCase().startsWith(q) ||
      (e.last_name + e.first_name).toLowerCase().startsWith(q) ||
      (e.last_name_kana + e.first_name_kana).toLowerCase().startsWith(q) ||
      e.email.toLowerCase().startsWith(q)
    );
  }, [employees, query]);

  return (
    <RootLayout>
      <div className='flex items-center justify-between'>
        <h1 className='text-3xl font-bold'>社員一覧</h1>
        <button
          onClick={() => router.push('/employees/new')}
          className='text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800'
        >
          新規作成
        </button>
      </div>

      <div className='mt-4'>
        <input
          type='text'
          placeholder='スタッフコード・氏名・メールアドレスで検索'
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          className='w-full max-w-md bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white'
        />
      </div>

      <div className='mt-4 overflow-x-auto'>
        <table className='w-full text-sm text-left text-gray-500 dark:text-gray-400'>
          <thead className='text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400'>
            <tr>
              <th className='px-6 py-3'>スタッフコード</th>
              <th className='px-6 py-3'>氏名</th>
              <th className='px-6 py-3'>氏名（カナ）</th>
              <th className='px-6 py-3'>メールアドレス</th>
            </tr>
          </thead>
          <tbody>
            {filtered.map((employee) => (
              <tr
                key={employee.ID}
                onClick={() => router.push(`/employees/${employee.ID}`)}
                className='bg-white border-b dark:bg-gray-800 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600 cursor-pointer'
              >
                <td className='px-6 py-4'>{employee.staff_code}</td>
                <td className='px-6 py-4'>{employee.last_name} {employee.first_name}</td>
                <td className='px-6 py-4'>{employee.last_name_kana} {employee.first_name_kana}</td>
                <td className='px-6 py-4'>{employee.email}</td>
              </tr>
            ))}
          </tbody>
        </table>
        {filtered.length === 0 && (
          <p className='mt-4 text-center text-gray-500'>該当する社員が見つかりません</p>
        )}
      </div>
    </RootLayout>
  );
};

export default Employees;
