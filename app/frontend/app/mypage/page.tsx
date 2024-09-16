"use client";

import RootLayout from '@/components/RootLayout';
import { useState, useEffect } from 'react';
import newApiInstance from "../api"

const MyPage = () => {
  const [user, setUser] = useState<any>() // useState<Array<Todo>>();

  const fetchUser = async () => {
    try {
      const api = newApiInstance();
      const response = await api.get('/api/user')
  
      console.log(response);
      setUser(response.data)
    } catch (e) {
      console.log(e)
    }
  }

  useEffect(() => {
    fetchUser();
  }, [])

  return (
    <main>
      <RootLayout>
        <>
          <h1 className='text-3xl font-bold'>マイページ</h1>
          <div className="grid grid-cols-2 mt-6 rounded-md shadow-md dark:bg-gray-800 dark:border-gray-700">
            <div className='p-5'>
              <p className='text-s'>
                姓
              </p>
            </div>
            <div className='p-5'>
              <p className='text-s'>
                {user?.UserInfo?.last_name}
              </p>
            </div>
            <div className='p-5'>
              <p className='text-s'>
                名
              </p>
            </div>
            <div className='p-5'>
              <p className='text-s'>
                {user?.UserInfo?.first_name}
              </p>
            </div>
            <div className='p-5'>
              <p className='text-s'>
                メールアドレス
              </p>
            </div>
            <div className='p-5'>
              <p className='text-s'>
                {user?.email}
              </p>
            </div>
            <div className='p-5'>
              <p className='text-s'>
                誕生日
              </p>
            </div>
            <div className='p-5'>
              <p className='text-s'>
                {user?.UserInfo?.birthday}
              </p>
            </div>
            <div className='p-5'>
              <p className='text-s'>
                性別
              </p>
            </div>
            <div className='p-5'>
              <p className='text-s'>
                {user?.UserInfo?.gender ? '男' : '女'}
              </p>
            </div>
            <div className='p-5'>
              <p className='text-s'>
                在籍状況
              </p>
            </div>
            <div className='p-5'>
              <p className='text-s'>
                {user?.UserInfo?.working ? '在籍中' : '離職済'}
              </p>
            </div>
          </div>
        </>
        
      </RootLayout>
    </main>  
  )
}

export default MyPage
  