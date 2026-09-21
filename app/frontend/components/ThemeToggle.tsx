"use client";

import { useEffect, useState } from 'react';

export type Theme = 'system' | 'light' | 'dark';

export const THEME_STORAGE_KEY = 'theme';

const prefersDark = () => window.matchMedia('(prefers-color-scheme: dark)').matches;

// Tailwind は darkMode: 'class' 設定なので、html の dark クラスで表示が切り替わる
const applyTheme = (theme: Theme) => {
  const isDark = theme === 'dark' || (theme === 'system' && prefersDark());
  document.documentElement.classList.toggle('dark', isDark);
};

const readTheme = (): Theme => {
  try {
    const saved = window.localStorage.getItem(THEME_STORAGE_KEY);
    if (saved === 'light' || saved === 'dark' || saved === 'system') {
      return saved;
    }
  } catch {
    // プライベートモードなどで localStorage が使えない場合はシステム設定に従う
  }
  return 'system';
};

const NEXT_THEME: Record<Theme, Theme> = {
  system: 'light',
  light: 'dark',
  dark: 'system',
};

const LABEL: Record<Theme, string> = {
  system: 'システム',
  light: 'ライト',
  dark: 'ダーク',
};

const Icon = ({ theme }: { theme: Theme }) => {
  const common = { className: 'w-4 h-4', fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', strokeWidth: 2 };
  if (theme === 'light') {
    return (
      <svg {...common} aria-hidden='true'>
        <path strokeLinecap='round' strokeLinejoin='round' d='M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z' />
      </svg>
    );
  }
  if (theme === 'dark') {
    return (
      <svg {...common} aria-hidden='true'>
        <path strokeLinecap='round' strokeLinejoin='round' d='M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z' />
      </svg>
    );
  }
  return (
    <svg {...common} aria-hidden='true'>
      <path strokeLinecap='round' strokeLinejoin='round' d='M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z' />
    </svg>
  );
};

const ThemeToggle = () => {
  const [theme, setTheme] = useState<Theme>('system');
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    const current = readTheme();
    setTheme(current);
    applyTheme(current);
    setMounted(true);
  }, []);

  // システム設定に追従しているときは、OS 側の切り替えもその場で反映する
  useEffect(() => {
    if (!mounted || theme !== 'system') return;
    const mql = window.matchMedia('(prefers-color-scheme: dark)');
    const onChange = () => applyTheme('system');
    mql.addEventListener('change', onChange);
    return () => mql.removeEventListener('change', onChange);
  }, [mounted, theme]);

  const cycle = () => {
    const next = NEXT_THEME[theme];
    setTheme(next);
    applyTheme(next);
    try {
      window.localStorage.setItem(THEME_STORAGE_KEY, next);
    } catch {
      // 保存できなくても表示の切り替えは行う(次回はシステム設定に戻る)
    }
  };

  // サーバー側では実際のテーマが分からないため、マウント前は場所だけ確保して
  // hydration の不一致とレイアウトのズレを避ける
  if (!mounted) {
    return <div className='h-9 w-[104px]' aria-hidden='true' />;
  }

  return (
    <button
      type='button'
      onClick={cycle}
      aria-label={`表示テーマ: ${LABEL[theme]}。クリックで切り替え`}
      title={`表示テーマ: ${LABEL[theme]}`}
      className='inline-flex h-9 w-[104px] items-center justify-center gap-1.5 rounded-lg border border-black/10 bg-white/80 px-3 text-sm font-medium text-gray-700 hover:bg-white dark:border-white/20 dark:bg-gray-800/80 dark:text-gray-200 dark:hover:bg-gray-800'
    >
      <Icon theme={theme} />
      {LABEL[theme]}
    </button>
  );
};

export default ThemeToggle;
