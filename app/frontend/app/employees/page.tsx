"use client";

import { useEffect, useRef, useState, useMemo } from 'react';
import { useRouter } from 'next/navigation';
import RootLayout from '@/components/RootLayout';
import useApi from '@/app/api';

const CSV_COLUMNS = [
  { field: 'staff_code',      label: 'スタッフコード',   required: true },
  { field: 'last_name',       label: '姓',               required: true },
  { field: 'first_name',      label: '名',               required: true },
  { field: 'last_name_kana',  label: '姓（カナ）',       required: false },
  { field: 'first_name_kana', label: '名（カナ）',       required: false },
  { field: 'email',           label: 'メールアドレス',   required: true },
] as const;

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
  const [showExportModal, setShowExportModal] = useState(false);
  const importRef = useRef<HTMLInputElement>(null);

  const fetchEmployees = () => {
    api.get('/api/employees').then((res) => setEmployees(res.data ?? []));
  };

  useEffect(() => {
    fetchEmployees();
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

  const handleExport = async () => {
    const res = await api.get('/api/employees/export', { responseType: 'blob' });
    const url = URL.createObjectURL(new Blob([res.data], { type: 'text/csv' }));
    const a = document.createElement('a');
    a.href = url;
    a.download = 'employees.csv';
    a.click();
    URL.revokeObjectURL(url);
    setShowExportModal(false);
  };

  const handleImport = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    const form = new FormData();
    form.append('file', file);
    await api.post('/api/employees/import', form, {
      headers: { 'Content-Type': 'multipart/form-data' },
    });
    e.target.value = '';
    fetchEmployees();
  };

  return (
    <RootLayout>
      <div className='flex items-center justify-between'>
        <h1 className='text-3xl font-bold'>社員一覧</h1>
        <div className='flex gap-2'>
          <button
            onClick={() => setShowExportModal(true)}
            className='text-gray-700 bg-white border border-gray-300 hover:bg-gray-100 font-medium rounded-lg text-sm px-5 py-2.5 dark:bg-gray-800 dark:text-gray-400 dark:border-gray-600 dark:hover:bg-gray-700'
          >
            エクスポート
          </button>
          <button
            onClick={() => importRef.current?.click()}
            className='text-gray-700 bg-white border border-gray-300 hover:bg-gray-100 font-medium rounded-lg text-sm px-5 py-2.5 dark:bg-gray-800 dark:text-gray-400 dark:border-gray-600 dark:hover:bg-gray-700'
          >
            インポート
          </button>
          <input
            ref={importRef}
            type='file'
            accept='.csv'
            className='hidden'
            onChange={handleImport}
          />
          <button
            onClick={() => router.push('/employees/new')}
            className='text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800'
          >
            新規作成
          </button>
        </div>
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
      {showExportModal && (
        <div className='fixed inset-0 z-50 flex items-center justify-center bg-black/50'>
          <div className='bg-white dark:bg-gray-800 rounded-lg shadow-lg w-full max-w-lg mx-4'>
            <div className='flex items-center justify-between px-6 py-4 border-b dark:border-gray-700'>
              <h2 className='text-lg font-semibold text-gray-900 dark:text-white'>CSVエクスポート</h2>
              <button
                onClick={() => setShowExportModal(false)}
                className='text-gray-400 hover:text-gray-600 dark:hover:text-gray-300'
              >
                <svg className='w-5 h-5' fill='none' stroke='currentColor' viewBox='0 0 24 24'>
                  <path strokeLinecap='round' strokeLinejoin='round' strokeWidth={2} d='M6 18L18 6M6 6l12 12' />
                </svg>
              </button>
            </div>

            <div className='px-6 py-4'>
              <p className='text-sm text-gray-500 dark:text-gray-400 mb-4'>
                出力されるCSVの列は以下の通りです。
              </p>
              <table className='w-full text-sm text-left'>
                <thead>
                  <tr className='border-b dark:border-gray-600'>
                    <th className='pb-2 font-medium text-gray-700 dark:text-gray-300'>列名（ヘッダー）</th>
                    <th className='pb-2 font-medium text-gray-700 dark:text-gray-300'>項目名</th>
                    <th className='pb-2 font-medium text-gray-700 dark:text-gray-300 text-center'>必須</th>
                  </tr>
                </thead>
                <tbody>
                  {CSV_COLUMNS.map((col) => (
                    <tr key={col.field} className='border-b dark:border-gray-700 last:border-0'>
                      <td className='py-2 font-mono text-xs text-gray-600 dark:text-gray-400'>{col.field}</td>
                      <td className='py-2 text-gray-800 dark:text-gray-200'>{col.label}</td>
                      <td className='py-2 text-center'>
                        {col.required ? (
                          <span className='px-2 py-0.5 text-xs font-medium bg-red-100 text-red-700 rounded dark:bg-red-900 dark:text-red-300'>必須</span>
                        ) : (
                          <span className='px-2 py-0.5 text-xs font-medium bg-gray-100 text-gray-500 rounded dark:bg-gray-700 dark:text-gray-400'>任意</span>
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            <div className='flex justify-end gap-2 px-6 py-4 border-t dark:border-gray-700'>
              <button
                onClick={() => setShowExportModal(false)}
                className='text-gray-700 bg-white border border-gray-300 hover:bg-gray-100 font-medium rounded-lg text-sm px-5 py-2.5 dark:bg-gray-800 dark:text-gray-400 dark:border-gray-600 dark:hover:bg-gray-700'
              >
                キャンセル
              </button>
              <button
                onClick={handleExport}
                className='text-white bg-blue-700 hover:bg-blue-800 font-medium rounded-lg text-sm px-5 py-2.5 dark:bg-blue-600 dark:hover:bg-blue-700'
              >
                エクスポート
              </button>
            </div>
          </div>
        </div>
      )}
    </RootLayout>
  );
};

export default Employees;
