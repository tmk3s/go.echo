"use client";

import Axios from 'axios';
import { useForm } from "react-hook-form";
import { useState, useEffect } from 'react';
import {Form, PrimaryBtn, DefaultBtn, GreenBtn, RedBtn, DarkBtn, BaseModal} from '@/consts/styles';
import RootLayout from '@/components/RootLayout';

const DepartmentsCsv = () => {
  return (
    <main>
      <RootLayout>
        <>
          <h1 className='text-3xl font-bold'>部署CSV</h1>
          <div className={`${Form} mt-20`}>
            <div className={`p-5`}>
              <div>
                <h2 className='font-bold'>CSVダウンロード</h2>
                <p className='mb-5'>CSVのダウンロードができます</p>
                <button
                  className={PrimaryBtn}
                  onClick={() => {
                  }}>
                  <div className='flex'>
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-6">
                      <path strokeLinecap="round" strokeLinejoin="round" d="M12 9.75v6.75m0 0-3-3m3 3 3-3m-8.25 6a4.5 4.5 0 0 1-1.41-8.775 5.25 5.25 0 0 1 10.233-2.33 3 3 0 0 1 3.758 3.848A3.752 3.752 0 0 1 18 19.5H6.75Z" />
                    </svg>
                    <p className='mt-[3px]'>CSVダウンロード</p>
                  </div>
                </button>
              </div>
              <div className='mt-10'>
                <h2 className='font-bold'>CSVアップロード</h2>
                <p className='mb-5'>CSVのアップロードができます</p>

                {/* https://flowbite.com/docs/forms/file-input/ */}
                <div className="flex items-center justify-center w-full">
                  <label htmlFor="dropzone-file" className="flex flex-col items-center justify-center w-full h-64 border-2 border-gray-300 border-dashed rounded-lg cursor-pointer bg-gray-50 dark:hover:bg-gray-800 dark:bg-gray-700 hover:bg-gray-100 dark:border-gray-600 dark:hover:border-gray-500 dark:hover:bg-gray-600">
                    <div className="flex flex-col items-center justify-center pt-5 pb-6">
                      <svg className="w-8 h-8 mb-4 text-gray-500 dark:text-gray-400" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 20 16">
                        <path stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M13 13h3a3 3 0 0 0 0-6h-.025A5.56 5.56 0 0 0 16 6.5 5.5 5.5 0 0 0 5.207 5.021C5.137 5.017 5.071 5 5 5a4 4 0 0 0 0 8h2.167M10 15V6m0 0L8 8m2-2 2 2"/>
                      </svg>
                      <p className="mb-2 text-sm text-gray-500 dark:text-gray-400"><span className="font-semibold">Click to upload</span> or drag and drop</p>
                      <p className="text-xs text-gray-500 dark:text-gray-400">SVG, PNG, JPG or GIF (MAX. 800x400px)</p>
                    </div>
                    <input id="dropzone-file" type="file" className="hidden" />
                  </label>
                </div>
              </div>
            </div>
          </div>
        </>
      </RootLayout>
    </main>  
  )
}
  
export default DepartmentsCsv