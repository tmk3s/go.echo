"use client";

import Axios from 'axios';
import { useForm } from "react-hook-form";
import { useRouter } from 'next/navigation'
import { useState, useEffect } from 'react';
import RootLayout from '@/components/RootLayout';
import {Form, PrimaryBtn, DefaultBtn, GreenBtn, RedBtn, DarkBtn, BaseModal} from '@/consts/styles';
import newApiInstance from "../../api"


type department = {
  ID: number,
  name: string,
  parentId: number | null
}

const Modal = ({obj, onSubmit, mode, setOpenModal}: {obj: department | null, onSubmit: any, mode: string, setOpenModal: any } ): React.ReactNode => {
  const {
    register,
    handleSubmit,
    watch,
    formState: { errors }
  } = useForm({
    defaultValues: {
      ID: obj?.ID,
      parentId: obj?.ID,
      name: mode === 'create' ? null : obj?.name,
    }
  });

  return (
    <div className={`${BaseModal}`}>
      <form
        onSubmit={handleSubmit((data) => {
          onSubmit({ID: data.ID, parentId: data.parentId, name: data.name});
        })}
      >
        {/* <!-- Modal content --> */}
        <div className="w-96 rounded-lg shadow dark:bg-gray-700">
          {/* <!-- Modal body --> */}
          <div className="p-4 md:p-5 space-y-4">
            <input type="hidden" id="parentId" {...register("parentId")} />
            <label>部署名</label>
            <input
              id="name"
              className="block p-2.5 w-full text-sm text-gray-900 bg-gray-50 rounded-lg border border-gray-300 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
              required
              {...register("name")}
            />
          </div>
          {/* <!-- Modal footer --> */}
          <div className="flex justify-center p-4 md:p-5 border-t border-gray-200 rounded-b dark:border-gray-600">
            <button type="submit" className="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800">
              保存する
            </button>
            <button type="button" onClick={() => setOpenModal(false)} className="py-2.5 px-5 ms-3 text-sm font-medium text-gray-900 focus:outline-none bg-white rounded-lg border border-gray-200 hover:bg-gray-100 hover:text-blue-700 focus:z-10 focus:ring-4 focus:ring-gray-100 dark:focus:ring-gray-700 dark:bg-gray-800 dark:text-gray-400 dark:border-gray-600 dark:hover:text-white dark:hover:bg-gray-700">
              キャンセル
            </button>
          </div>
        </div>
      </form>
    </div>
  )
}

const Departments = () => {
  const router = useRouter();
  const [openCreateModal, setCreateOpenModal] = useState(false);
  const [openUpdateModal, setUpdateOpenModal] = useState(false);
  const [openDeleteModal, setDeleteOpenModal] = useState(false);
  const [departments, setDepartments] = useState([]);
  const [targetDepartment, setTargetDepartment] = useState<department | null>(null);

  const calcDepth = (depth: number) => {
    // Tailwind がクラス名を抽出する方法の最も重要な点は、 ソースファイル中に完全な文字列として存在するクラスのみを検出することです。
    // もし、文字列の補間をしたり、クラス名の一部を連結したりすると、Tailwind はそれを見つけられず、対応する CSS を生成することができません。
    switch (depth) {
      case 0:
        return 'ml-[20px]'
      case 1:
        return 'ml-[40px]'
      case 2:
        return 'ml-[60px]'
      case 3:
        return 'ml-[80px]'
      case 4:
        return 'ml-[100px]'
      case 5:
        return 'ml-[120px]'
      case 6:
        return 'ml-[140px]'
      default:
        return 'ml-[20px]'
    }
  }

  const fetchDepartments = async () => {
    try {
      const api = newApiInstance();
      const response = await api.get('/api/departments')
  
      console.log(response);
      setDepartments(response.data);
    } catch (e) {
      console.log(e);
    }
  }

  const createDepartment = async (obj: department) => {
    try {
      const api = newApiInstance();
      // const formData = new FormData(); 数値が扱えないので使用しない
      const response = await api.post('/api/department', {
        parent_id: obj.parentId ? obj.parentId : null,
        name: obj.name
      })
  
      console.log(response);
      fetchDepartments();
      setCreateOpenModal(false);
    } catch (e) {
      console.log(e);
    }
  }

  const updateDepartment = async (obj: department) => {
    try {
      const api = newApiInstance();
      // const formData = new FormData(); 数値が扱えないので使用しない
      const response = await api.put(`/api/departments/${obj.ID}`, {
        name: obj.name
      })
  
      console.log(response);
      fetchDepartments();
      setUpdateOpenModal(false);
    } catch (e) {
      console.log(e);
    }
  }

  const deleteDepartment = async (id: number) => {
    try {
      const api = newApiInstance();
      // const formData = new FormData(); 数値が扱えないので使用しない
      const response = await api.delete(`/api/departments/${id}`)
  
      console.log(response);
      fetchDepartments();
      setDeleteOpenModal(false);
    } catch (e) {
      console.log(e);
    }
  }

  useEffect(()=>{
	  fetchDepartments()
  },[])

  return (
    <main>
      <RootLayout>
        <>
          <h1 className='text-3xl font-bold'>部署情報</h1>
          <div className='mt-8 mb-8'>
            <button
              className={PrimaryBtn}
              onClick={() => {
                setTargetDepartment(null)
                setCreateOpenModal(true)}
              }>
              新規に部署を追加
            </button>
            <button
              className={`${DarkBtn} px-5 py-2.5 ml-8` }
              onClick={() => {
                router.push('/organizations/departments/csv')
              }}>
              一括操作
            </button>
          </div>
          <div className={`rounded-md shadow-md dark:bg-gray-900 dark:border-gray-700`}>
            <div className='grid grid-cols-12 pt-5 pb-5 border-b border-gray-200 dark:bg-gray-950 dark:border-gray-600'>
              <div className='col-span-10 ml-5'>部署名</div>
              <div className='mr-5'>操作</div>
            </div>
            <div>
              {
                departments?.map((department: any, index: number) => {
                  return (
                    <div key={index} className='grid grid-cols-12 gap-2 mt-3 mb-3'>
                      <div className={`${calcDepth(department.depth)} grid-cols-9 col-span-9 dark:border-gray-600 break-words`}>
                        <span className='text-wrap'>{department.name}</span>
                      </div>
                      <button
                        className={`${DefaultBtn} col-span-1`}
                        onClick={() => {
                          setTargetDepartment(department)
                          setCreateOpenModal(true)}
                        }>
                        追加
                      </button>
                      <button
                        className={`${GreenBtn} col-span-1`}
                        onClick={() => {
                          setTargetDepartment(department)
                          setUpdateOpenModal(true)}
                        }>
                        編集
                      </button>
                      <button
                        className={`${RedBtn} col-span-1`}
                        onClick={() => {
                          setTargetDepartment(department)
                          setDeleteOpenModal(true)}
                        }>
                        削除
                      </button>
                    </div>
                  )
                })
              }
            </div>
          </div>
        </>
      </RootLayout>
      { openCreateModal && (
        <Modal
          mode='create'
          obj={targetDepartment}
          onSubmit={createDepartment}
          setOpenModal={setCreateOpenModal} 
        />
      )}
      { openUpdateModal && (
        <Modal
          mode='update'
          obj={targetDepartment}
          onSubmit={updateDepartment}
          setOpenModal={setUpdateOpenModal} 
        />
      )}
      { (openDeleteModal && targetDepartment)&& (
        <div className={`${BaseModal}`}>
            {/* <!-- Modal content --> */}
          <div className="w-96 rounded-lg shadow dark:bg-gray-700">
            {/* <!-- Modal body --> */}
            <div className="p-4 md:p-5 space-y-4">
              <p>本当に削除しますか？</p>
            </div>
            {/* <!-- Modal footer --> */}
            <div className="flex justify-center p-4 md:p-5 border-t border-gray-200 rounded-b dark:border-gray-600">
              <button type="submit" className={`${RedBtn} font-medium rounded-lg text-sm px-5 py-2.5 text-center`} onClick={() => deleteDepartment(targetDepartment.ID)}>
                削除
              </button>
              <button type="button" onClick={() => setDeleteOpenModal(false)} className="py-2.5 px-5 ms-3 text-sm font-medium text-gray-900 focus:outline-none bg-white rounded-lg border border-gray-200 hover:bg-gray-100 hover:text-blue-700 focus:z-10 focus:ring-4 focus:ring-gray-100 dark:focus:ring-gray-700 dark:bg-gray-800 dark:text-gray-400 dark:border-gray-600 dark:hover:text-white dark:hover:bg-gray-700">
                キャンセル
              </button>
            </div>
          </div>
        </div>
      )}
    </main>  
  )
}
  
export default Departments