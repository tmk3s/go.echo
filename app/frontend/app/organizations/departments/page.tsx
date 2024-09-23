"use client";

import Axios from 'axios';
import { useForm } from "react-hook-form";
import { useState, useEffect } from 'react';
import RootLayout from '@/components/RootLayout';
import newApiInstance from "../../api"


type department = {
  ID: number,
  name: string,
  parentId: number | null
}

const Modal = ({obj, onSubmit, setOpenModal}: {obj: department | null, onSubmit: any, setOpenModal: any } ): React.ReactNode => {
  const {
    register,
    handleSubmit,
    watch,
    formState: { errors }
  } = useForm({
    defaultValues: {
      parentId: obj?.ID,
      name: '',
    }
  });


  return (
    <div className="overflow-y-auto overflow-x-hidden fixed top-0 right-0 left-0 z-50 justify-center items-center w-full md:inset-0 h-[calc(100%-1rem)] max-h-full">
      <form
        className="max-w-sm mx-auto"
        onSubmit={handleSubmit((data) => {
          onSubmit({parentId: data.parentId, name: data.name});
        })}
      >
        <div className="relative p-4 w-full max-w-2xl max-h-full">
          {/* <!-- Modal content --> */}
          <div className="relative bg-white rounded-lg shadow dark:bg-gray-700">
            {/* <!-- Modal body --> */}
            <div className="p-4 md:p-5 space-y-4">
              <input type="hidden" id="parentId" {...register("parentId")} />
              <textarea
                id="name"
                className="block p-2.5 w-full text-sm text-gray-900 bg-gray-50 rounded-lg border border-gray-300 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                required
                {...register("name")}
              />
            </div>
            {/* <!-- Modal footer --> */}
            <div className="flex items-center p-4 md:p-5 border-t border-gray-200 rounded-b dark:border-gray-600">
              <button type="submit" className="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800">
                Save
              </button>
              <button type="button" onClick={() => setOpenModal(false)} className="py-2.5 px-5 ms-3 text-sm font-medium text-gray-900 focus:outline-none bg-white rounded-lg border border-gray-200 hover:bg-gray-100 hover:text-blue-700 focus:z-10 focus:ring-4 focus:ring-gray-100 dark:focus:ring-gray-700 dark:bg-gray-800 dark:text-gray-400 dark:border-gray-600 dark:hover:text-white dark:hover:bg-gray-700">
                Cancel
              </button>
            </div>
          </div>
        </div>
      </form>
    </div>
  )
}

const Departments = () => {
  const [openModal, setOpenModal] = useState(false);
  const [departments, setDepartments] = useState([]);
  const [targetDepartment, setTargetDepartment] = useState(null);

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
      setDepartments(response.data);
      setOpenModal(false);
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
          <button
            className="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800"
            onClick={() => {
              setTargetDepartment(null)
              setOpenModal(true)}
            }>
            追加
          </button>
          <div className="grid">
          {
            departments?.map((department: any) => {
              return (
                <div key={department.ID} className='grid-rows-12  dark:hover:bg-gray-700'>
                  <p
                    onClick={() => {
                      setTargetDepartment(department)
                      setOpenModal(true)}
                    }
                  >
                    {'-'.repeat(department.depth)}{department.name}
                  </p>
                </div>
              )
            })
          }
          </div> 
        </>
      </RootLayout>
      { openModal && (
        <Modal obj={targetDepartment} onSubmit={createDepartment} setOpenModal={setOpenModal} />
      )}
    </main>  
  )
}
  
export default Departments