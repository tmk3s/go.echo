"use client";

import Link from 'next/link';
import RootLayout from '@/components/RootLayout';

const cardBase = 'rounded-lg p-4 shadow-lg border-b-2 border-gray-200 dark:border-gray-700 block';

const Organizations = () => {
  return (
    <RootLayout>
      <h1 className='text-3xl font-bold'>組織情報</h1>
      <div className="grid grid-cols-2 gap-4 mt-6">
        <Link
          href='/organizations/departments'
          className={`${cardBase} bg-white hover:bg-gray-50 dark:bg-gray-800 dark:hover:bg-gray-700`}
        >
          <h2 className="mt-2 font-bold text-gray-900 dark:text-white">部署</h2>
          <p className='text-sm text-gray-500 dark:text-gray-400'>部署情報の確認、編集ができます</p>
        </Link>

        {/* 職種は未実装のため押せないことが分かる見た目にする */}
        <div
          aria-disabled='true'
          className={`${cardBase} bg-white cursor-not-allowed opacity-60 dark:bg-gray-800`}
        >
          <h2 className="mt-2 font-bold text-gray-500 dark:text-gray-400">
            職種
            <span className='ml-2 px-1.5 py-0.5 text-xs font-normal rounded bg-gray-100 text-gray-500 dark:bg-gray-700 dark:text-gray-400'>
              準備中
            </span>
          </h2>
          <p className='text-sm text-gray-400 dark:text-gray-500'>職種情報の確認、編集ができます</p>
        </div>
      </div>
    </RootLayout>
  )
}

export default Organizations
