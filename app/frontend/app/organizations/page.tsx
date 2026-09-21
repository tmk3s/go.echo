"use client";

import Link from 'next/link';
import RootLayout from '@/components/RootLayout';

const cardBase = 'rounded-lg p-4 shadow-lg border-b-2 border-line block';

const Organizations = () => {
  return (
    <RootLayout>
      <h1 className='text-3xl font-bold'>組織情報</h1>
      <div className="grid grid-cols-2 gap-4 mt-6">
        <Link
          href='/organizations/departments'
          className={`${cardBase} bg-surface hover:bg-surface-muted`}
        >
          <h2 className="mt-2 font-bold text-body">部署</h2>
          <p className='text-sm text-muted'>部署情報の確認、編集ができます</p>
        </Link>

        {/* 職種は未実装のため押せないことが分かる見た目にする */}
        <div
          aria-disabled='true'
          className={`${cardBase} bg-surface cursor-not-allowed opacity-60`}
        >
          <h2 className="mt-2 font-bold text-muted">
            職種
            <span className='ml-2 px-1.5 py-0.5 text-xs font-normal rounded bg-surface-muted text-muted'>
              準備中
            </span>
          </h2>
          <p className='text-sm text-muted'>職種情報の確認、編集ができます</p>
        </div>
      </div>
    </RootLayout>
  )
}

export default Organizations
