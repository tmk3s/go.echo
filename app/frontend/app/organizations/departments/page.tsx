"use client";

import { useState, useEffect } from 'react';
import RootLayout from '@/components/RootLayout';

export default () => {
  useEffect(()=>{
	  console.log(1)
  },[])
  return (
    <main>
      <RootLayout>
        <>
          <h1 className='text-3xl font-bold'>部署情報</h1>
        </>
      </RootLayout>
    </main>  
  )
}
  