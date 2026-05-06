"use client";

import { useEffect, useRef, useState, useMemo } from 'react';
import { useRouter } from 'next/navigation';
import RootLayout from '@/components/RootLayout';
import useApi from '@/app/api';

type CsvColumn = { label: string; required: boolean };
type CsvSection = { title: string; note?: string; columns: CsvColumn[] };

const SAMPLE_CSV_ROWS = [
  ['スタッフコード', '姓', '名', '姓（カナ）', '名（カナ）', 'メールアドレス', '郵便番号', '都道府県', '市区町村', '住所1', '住所2', '電話番号', '入社日', '退職日', '退職区分', 'ステータス', '部署1'],
  ['S001', '山田', '太郎', 'ヤマダ', 'タロウ', 'yamada.taro@example.com', '100-0001', '東京都', '千代田区', '千代田1-1-1', '', '03-1234-5678', '2020-04-01', '', '', '在籍', '営業部'],
  ['S002', '鈴木', '花子', 'スズキ', 'ハナコ', 'suzuki.hanako@example.com', '530-0001', '大阪府', '大阪市北区', '梅田2-2-2', 'ビル3F', '06-9876-5432', '2019-10-01', '2023-03-31', '自己都合', '退職', '開発部'],
];

const downloadSampleCSV = () => {
  const bom = '﻿';
  const csv = SAMPLE_CSV_ROWS.map((row) =>
    row.map((cell) => (cell.includes(',') || cell.includes('"') ? `"${cell.replace(/"/g, '""')}"` : cell)).join(',')
  ).join('\n');
  const blob = new Blob([bom + csv], { type: 'text/csv' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = 'employees_sample.csv';
  a.click();
  URL.revokeObjectURL(url);
};

const CSV_SECTIONS: CsvSection[] = [
  {
    title: '基本情報',
    columns: [
      { label: 'スタッフコード', required: true },
      { label: '姓',             required: true },
      { label: '名',             required: true },
      { label: '姓（カナ）',     required: false },
      { label: '名（カナ）',     required: false },
      { label: 'メールアドレス', required: true },
    ],
  },
  {
    title: '住所情報',
    columns: [
      { label: '郵便番号', required: false },
      { label: '都道府県', required: false },
      { label: '市区町村', required: false },
      { label: '住所1',    required: false },
      { label: '住所2',    required: false },
      { label: '電話番号', required: false },
    ],
  },
  {
    title: '在籍情報',
    columns: [
      { label: '入社日',     required: false },
      { label: '退職日',     required: false },
      { label: '退職区分',   required: false },
      { label: 'ステータス', required: false },
    ],
  },
  {
    title: '所属部署',
    note: '複数ある場合は連番（部署1、部署2…）で出力されます',
    columns: [
      { label: '部署名', required: false },
    ],
  },
];

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
  const [showDropdown, setShowDropdown] = useState(false);
  const [showExportModal, setShowExportModal] = useState(false);
  const [showBulkCreateModal, setShowBulkCreateModal] = useState(false);
  const [showBulkUpdateModal, setShowBulkUpdateModal] = useState(false);
  const [importError, setImportError] = useState<string | null>(null);
  const dropdownRef = useRef<HTMLDivElement>(null);
  const bulkCreateRef = useRef<HTMLInputElement>(null);
  const bulkUpdateRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    const handler = (e: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(e.target as Node)) {
        setShowDropdown(false);
      }
    };
    document.addEventListener('mousedown', handler);
    return () => document.removeEventListener('mousedown', handler);
  }, []);

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

  const handleImport = async (e: React.ChangeEvent<HTMLInputElement>, endpoint: string) => {
    const file = e.target.files?.[0];
    if (!file) return;
    const form = new FormData();
    form.append('file', file);
    setImportError(null);
    try {
      await api.post(endpoint, form, {
        headers: { 'Content-Type': 'multipart/form-data' },
      });
      fetchEmployees();
    } catch (err: any) {
      const msg = err?.response?.data?.message;
      if (msg) setImportError(msg);
    } finally {
      e.target.value = '';
    }
  };

  return (
    <RootLayout>
      <div className='flex items-center justify-between'>
        <h1 className='text-3xl font-bold'>社員一覧</h1>
        <div className='flex gap-2'>
          <div className='relative' ref={dropdownRef}>
            <button
              onClick={() => setShowDropdown((v) => !v)}
              className='flex items-center gap-1 text-gray-700 bg-white border border-gray-300 hover:bg-gray-100 font-medium rounded-lg text-sm px-5 py-2.5 dark:bg-gray-800 dark:text-gray-400 dark:border-gray-600 dark:hover:bg-gray-700'
            >
              一括操作
              <svg className='w-4 h-4' fill='none' stroke='currentColor' viewBox='0 0 24 24'>
                <path strokeLinecap='round' strokeLinejoin='round' strokeWidth={2} d='M19 9l-7 7-7-7' />
              </svg>
            </button>
            {showDropdown && (
              <div className='absolute right-0 mt-1 w-36 bg-white dark:bg-gray-700 rounded-lg shadow-lg border border-gray-200 dark:border-gray-600 z-10 overflow-hidden'>
                {[
                  { label: 'CSV出力',  onClick: () => { setShowDropdown(false); setShowExportModal(true); } },
                  { label: '一括登録', onClick: () => { setShowDropdown(false); setShowBulkCreateModal(true); } },
                  { label: '一括更新', onClick: () => { setShowDropdown(false); setShowBulkUpdateModal(true); } },
                ].map((item) => (
                  <button
                    key={item.label}
                    onClick={item.onClick}
                    className='w-full text-left px-4 py-2.5 text-sm text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-600'
                  >
                    {item.label}
                  </button>
                ))}
              </div>
            )}
          </div>
          <input
            ref={bulkCreateRef}
            type='file'
            accept='.csv'
            className='hidden'
            onChange={(e) => handleImport(e, '/api/employees/import/create')}
          />
          <input
            ref={bulkUpdateRef}
            type='file'
            accept='.csv'
            className='hidden'
            onChange={(e) => handleImport(e, '/api/employees/import/update')}
          />
          <button
            onClick={() => router.push('/employees/new')}
            className='text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800'
          >
            新規作成
          </button>
        </div>
      </div>

      {importError && (
        <div className='mt-4 flex items-start gap-3 p-4 rounded-lg bg-red-50 border border-red-200 dark:bg-red-900/20 dark:border-red-800'>
          <svg className='w-5 h-5 text-red-500 shrink-0 mt-0.5' fill='currentColor' viewBox='0 0 20 20'>
            <path fillRule='evenodd' d='M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z' clipRule='evenodd' />
          </svg>
          <p className='text-sm text-red-700 dark:text-red-400 flex-1'>{importError}</p>
          <button onClick={() => setImportError(null)} className='text-red-400 hover:text-red-600 dark:hover:text-red-300'>
            <svg className='w-4 h-4' fill='none' stroke='currentColor' viewBox='0 0 24 24'>
              <path strokeLinecap='round' strokeLinejoin='round' strokeWidth={2} d='M6 18L18 6M6 6l12 12' />
            </svg>
          </button>
        </div>
      )}

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
      {showBulkCreateModal && (
        <ImportModal
          title='一括登録'
          description='スタッフコードが存在しない社員を新規作成します。すでに登録済みのスタッフコードがある場合はエラーになります。'
          onClose={() => setShowBulkCreateModal(false)}
          onSelect={() => { bulkCreateRef.current?.click(); setShowBulkCreateModal(false); }}
        />
      )}

      {showBulkUpdateModal && (
        <ImportModal
          title='一括更新'
          description='スタッフコードが存在する社員の情報を更新します。存在しないスタッフコードがある場合はエラーになります。'
          onClose={() => setShowBulkUpdateModal(false)}
          onSelect={() => { bulkUpdateRef.current?.click(); setShowBulkUpdateModal(false); }}
        />
      )}

      {showExportModal && (
        <div className='fixed inset-0 z-50 flex items-center justify-center bg-black/50'>
          <div className='bg-white dark:bg-gray-800 rounded-lg shadow-lg w-full max-w-lg mx-4 flex flex-col max-h-[90vh]'>
            <div className='flex items-center justify-between px-6 py-4 border-b dark:border-gray-700 shrink-0'>
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

            <div className='px-6 py-4 overflow-y-auto'>
              <p className='text-sm text-gray-500 dark:text-gray-400 mb-4'>
                出力されるCSVの項目は以下の通りです。
              </p>
              <div className='space-y-5'>
                {CSV_SECTIONS.map((section) => (
                  <div key={section.title}>
                    <p className='text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-1'>
                      {section.title}
                    </p>
                    {section.note && (
                      <p className='text-xs text-blue-600 dark:text-blue-400 mb-2'>{section.note}</p>
                    )}
                    <div className='bg-gray-50 dark:bg-gray-700/50 rounded-lg divide-y divide-gray-200 dark:divide-gray-600'>
                      {section.columns.map((col) => (
                        <div key={col.label} className='flex items-center justify-between px-4 py-2'>
                          <span className='text-sm text-gray-800 dark:text-gray-200'>{col.label}</span>
                          {col.required ? (
                            <span className='px-2 py-0.5 text-xs font-medium bg-red-100 text-red-700 rounded dark:bg-red-900 dark:text-red-300'>必須</span>
                          ) : (
                            <span className='px-2 py-0.5 text-xs font-medium bg-gray-200 text-gray-500 rounded dark:bg-gray-600 dark:text-gray-400'>任意</span>
                          )}
                        </div>
                      ))}
                    </div>
                  </div>
                ))}
              </div>
            </div>

            <div className='flex items-center justify-between px-6 py-4 border-t dark:border-gray-700 shrink-0'>
              <button
                onClick={downloadSampleCSV}
                className='flex items-center gap-1.5 text-sm text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300 font-medium'
              >
                <svg className='w-4 h-4' fill='none' stroke='currentColor' viewBox='0 0 24 24'>
                  <path strokeLinecap='round' strokeLinejoin='round' strokeWidth={2} d='M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4' />
                </svg>
                サンプルをダウンロード
              </button>
              <div className='flex gap-2'>
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
        </div>
      )}
    </RootLayout>
  );
};

type ImportModalProps = {
  title: string;
  description: string;
  onClose: () => void;
  onSelect: () => void;
};

const ImportModal = ({ title, description, onClose, onSelect }: ImportModalProps) => (
  <div className='fixed inset-0 z-50 flex items-center justify-center bg-black/50'>
    <div className='bg-white dark:bg-gray-800 rounded-lg shadow-lg w-full max-w-lg mx-4 flex flex-col max-h-[90vh]'>
      <div className='flex items-center justify-between px-6 py-4 border-b dark:border-gray-700 shrink-0'>
        <h2 className='text-lg font-semibold text-gray-900 dark:text-white'>{title}</h2>
        <button onClick={onClose} className='text-gray-400 hover:text-gray-600 dark:hover:text-gray-300'>
          <svg className='w-5 h-5' fill='none' stroke='currentColor' viewBox='0 0 24 24'>
            <path strokeLinecap='round' strokeLinejoin='round' strokeWidth={2} d='M6 18L18 6M6 6l12 12' />
          </svg>
        </button>
      </div>

      <div className='px-6 py-4 overflow-y-auto'>
        <p className='text-sm text-gray-500 dark:text-gray-400 mb-4'>{description}</p>
        <p className='text-sm text-gray-500 dark:text-gray-400 mb-4'>取り込むCSVの項目は以下の通りです。</p>
        <div className='space-y-5'>
          {CSV_SECTIONS.map((section) => (
            <div key={section.title}>
              <p className='text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-1'>
                {section.title}
              </p>
              {section.note && (
                <p className='text-xs text-blue-600 dark:text-blue-400 mb-2'>{section.note}</p>
              )}
              <div className='bg-gray-50 dark:bg-gray-700/50 rounded-lg divide-y divide-gray-200 dark:divide-gray-600'>
                {section.columns.map((col) => (
                  <div key={col.label} className='flex items-center justify-between px-4 py-2'>
                    <span className='text-sm text-gray-800 dark:text-gray-200'>{col.label}</span>
                    {col.required ? (
                      <span className='px-2 py-0.5 text-xs font-medium bg-red-100 text-red-700 rounded dark:bg-red-900 dark:text-red-300'>必須</span>
                    ) : (
                      <span className='px-2 py-0.5 text-xs font-medium bg-gray-200 text-gray-500 rounded dark:bg-gray-600 dark:text-gray-400'>任意</span>
                    )}
                  </div>
                ))}
              </div>
            </div>
          ))}
        </div>
      </div>

      <div className='flex justify-end gap-2 px-6 py-4 border-t dark:border-gray-700 shrink-0'>
        <button
          onClick={onClose}
          className='text-gray-700 bg-white border border-gray-300 hover:bg-gray-100 font-medium rounded-lg text-sm px-5 py-2.5 dark:bg-gray-800 dark:text-gray-400 dark:border-gray-600 dark:hover:bg-gray-700'
        >
          キャンセル
        </button>
        <button
          onClick={onSelect}
          className='text-white bg-blue-700 hover:bg-blue-800 font-medium rounded-lg text-sm px-5 py-2.5 dark:bg-blue-600 dark:hover:bg-blue-700'
        >
          CSVを選択して実行
        </button>
      </div>
    </div>
  </div>
);

export default Employees;
