"use client";

import Link from 'next/link';
import { usePathname } from 'next/navigation';

type MenuItem = {
  label: string;
  href?: string; // href が無い項目は未実装
};

const MENU_ITEMS: MenuItem[] = [
  { label: 'マイページ', href: '/mypage' },
  { label: '会社情報設定', href: '/setting' },
  { label: '給与情報設定' },
  { label: '社員一覧', href: '/employees' },
  { label: '組織情報設定', href: '/organizations' },
];

const itemBase = 'flex items-center justify-between px-3 py-2 text-sm';

const Layout = ({ children }: { children: React.ReactNode }) => {
  const pathname = usePathname();
  const isActive = (href: string) => pathname === href || pathname.startsWith(`${href}/`);

  return (
    <main>
      <div className="lg:flex lg:justify-center">
        <nav className='Menu lg:w-60 lg:min-w-[240px] m-16 h-fit rounded-md shadow-md bg-white border border-gray-200 dark:bg-gray-800 dark:border-gray-700 lg:sticky lg:top-16 lg:self-start'>
          <div className='h-16 font-bold pt-5 pl-3 text-gray-900 dark:text-white'>メニュー</div>
          <ul className='pb-2'>
            {MENU_ITEMS.map((item) => (
              <li key={item.label}>
                {item.href ? (
                  <Link
                    href={item.href}
                    aria-current={isActive(item.href) ? 'page' : undefined}
                    className={
                      isActive(item.href)
                        ? `${itemBase} bg-blue-50 text-blue-700 font-medium dark:bg-gray-700 dark:text-white`
                        : `${itemBase} text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700`
                    }
                  >
                    {item.label}
                  </Link>
                ) : (
                  // 未実装の項目は押せないことが分かるようにする
                  <span
                    aria-disabled='true'
                    className={`${itemBase} text-gray-400 cursor-not-allowed dark:text-gray-500`}
                  >
                    {item.label}
                    <span className='px-1.5 py-0.5 text-xs rounded bg-gray-100 text-gray-500 dark:bg-gray-700 dark:text-gray-400'>
                      準備中
                    </span>
                  </span>
                )}
              </li>
            ))}
          </ul>
        </nav>
        <div className='m-16 lg:w-[1920px]'>
          {children}
        </div>
      </div>
    </main>
  )
}

export default Layout
