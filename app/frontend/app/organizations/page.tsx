"use client";

import RootLayout from '@/components/RootLayout';

const Organizations = () => {
  return (
    <main>
      <RootLayout>
        <>
          <h1 className='text-3xl font-bold'>組織情報</h1>
          <div className="grid grid-cols-2 gap-4 mt-6">
            <a className='rounded-lg p-4 shadow-lg border-b-2 pt-2 pl-3 hover:bg-gray-700' href='/organizations/departments'>
              <h2 className="mt-2 font-bold">部署</h2>
              <p>部署情報の確認、編集ができます</p>
            </a>
            <a className='rounded-lg p-4 shadow-lg border-b-2 pt-2 pl-3 hover:bg-gray-700'>
              <h2 className="mt-2 font-bold">職種</h2>
              <p>職種情報の確認、編集ができます</p>
            </a>
          </div>
        </>
        
      </RootLayout>
    </main>  
  )
}

export default Organizations
  