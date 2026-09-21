"use client";

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import LogoutButton from '@/components/LogoutButton';

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
        <nav className='Menu lg:w-60 lg:min-w-[240px] m-16 h-fit rounded-md shadow-md bg-surface border border-line lg:sticky lg:top-16 lg:self-start'>
          <div className='h-16 font-bold pt-5 pl-3 text-body'>メニュー</div>
          <ul className='pb-2'>
            {MENU_ITEMS.map((item) => (
              <li key={item.label}>
                {item.href ? (
                  <Link
                    href={item.href}
                    aria-current={isActive(item.href) ? 'page' : undefined}
                    className={
                      isActive(item.href)
                        ? `${itemBase} bg-brand-subtle text-brand-text font-medium`
                        : `${itemBase} text-body hover:bg-surface-muted`
                    }
                  >
                    {item.label}
                  </Link>
                ) : (
                  // 未実装の項目は押せないことが分かるようにする
                  <span
                    aria-disabled='true'
                    className={`${itemBase} text-muted cursor-not-allowed`}
                  >
                    {item.label}
                    <span className='px-1.5 py-0.5 text-xs rounded bg-surface-muted text-muted'>
                      準備中
                    </span>
                  </span>
                )}
              </li>
            ))}
          </ul>
          <div className='border-t border-line py-2'>
            <LogoutButton />
          </div>
        </nav>
        <div className='m-16 lg:w-[1920px]'>
          {children}
        </div>
      </div>
    </main>
  )
}

export default Layout
