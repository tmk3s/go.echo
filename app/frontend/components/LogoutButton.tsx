"use client";

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import useApi from '@/app/api';

const LogoutButton = () => {
  const router = useRouter();
  const api = useApi();
  const [loading, setLoading] = useState(false);

  const logout = async () => {
    setLoading(true);
    try {
      // セッションCookieは HttpOnly のためサーバー側で失効させる
      await api.post('/sign_out');
    } catch (e) {
      // 失効に失敗してもログイン画面には戻す(期限切れなどで既に無効な場合がある)
      console.error(e);
    }
    // 前の画面のデータが残らないよう、履歴を置き換えたうえで再読み込みする
    router.replace('/sign_in');
    router.refresh();
  };

  return (
    <button
      type='button'
      onClick={logout}
      disabled={loading}
      className='flex w-full items-center gap-2 px-3 py-2 text-sm text-gray-700 hover:bg-gray-100 disabled:opacity-50 dark:text-gray-300 dark:hover:bg-gray-700'
    >
      <svg className='w-4 h-4' fill='none' stroke='currentColor' strokeWidth={2} viewBox='0 0 24 24' aria-hidden='true'>
        <path strokeLinecap='round' strokeLinejoin='round' d='M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1' />
      </svg>
      {loading ? 'ログアウト中...' : 'ログアウト'}
    </button>
  );
};

export default LogoutButton;
