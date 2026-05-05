 const Layout = ({children}: { children: React.ReactNode }) => {
  return (
    <main>
      <div className="lg:flex lg:justify-center">
        <div className='Menu lg:w-60 lg:min-w-[240px] m-16 h-fit rounded-md shadow-md dark:bg-gray-800 dark:border-gray-700 lg:sticky lg:top-16 lg:self-start'>
          <div className='h-16 font-bold pt-5 pl-3'>メニュー</div>
          <ul className=''>
            <li className='h-10 pt-2 pl-3 dark:hover:bg-gray-700'>
              <a className='w-full h-full top-0 bottom-0 block' href="/mypage">マイページ</a>
            </li>
            <li className='h-10 pt-2 pl-3 dark:hover:bg-gray-700'>
              <a className='w-full h-full top-0 bottom-0 block' href="/setting">会社情報設定</a>
            </li>
            <li className='h-10 pt-2 pl-3 dark:hover:bg-gray-700'>
              <a className='w-full h-full top-0 bottom-0 block' href="#">給与情報設定</a>
            </li>
            <li className='h-10 pt-2 pl-3 dark:hover:bg-gray-700'>
              <a className='w-full h-full top-0 bottom-0 block' href="/employees">社員一覧</a>
            </li>
            <li className='h-10 pt-2 pl-3 dark:hover:bg-gray-700'>
              <a className='w-full h-full top-0 bottom-0 block' href="/organizations">組織情報設定</a>
            </li>
          </ul>
        </div>
        <div className='m-16 lg:w-[1920px]'>
          {children}
        </div> 
      </div> 
    </main>  
  )
}
export default Layout