"use client";

import Axios from 'axios';
import { useForm } from "react-hook-form";
import { useRouter } from 'next/navigation'
import { useState, useEffect } from 'react';
import RootLayout from '@/components/RootLayout';
import Button from '@/components/ui/Button';
import Card from '@/components/ui/Card';
import ModalShell from '@/components/ui/Modal';
import TextInput, { Label } from '@/components/ui/TextInput';
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
    <ModalShell>
      <form
        onSubmit={handleSubmit((data) => {
          onSubmit({ID: data.ID, parentId: data.parentId, name: data.name});
        })}
      >
        <div>
          <div className="p-4 md:p-5 space-y-4">
            <input type="hidden" id="parentId" {...register("parentId")} />
            <Label>部署名</Label>
            <TextInput id="name" required {...register("name")} />
          </div>
          <div className="flex justify-center p-4 md:p-5 border-t border-line rounded-b">
            <Button type="submit">保存する</Button>
            <Button variant="secondary" className="ms-3" onClick={() => setOpenModal(false)}>
              キャンセル
            </Button>
          </div>
        </div>
      </form>
    </ModalShell>
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
            <Button
              onClick={() => {
                setTargetDepartment(null)
                setCreateOpenModal(true)}
              }>
              新規に部署を追加
            </Button>
            <Button
              variant='secondary'
              className='ml-8'
              onClick={() => {
                router.push('/organizations/departments/csv')
              }}>
              一括操作
            </Button>
          </div>
          <Card className='rounded-md shadow-md'>
            <div className='grid grid-cols-12 pt-5 pb-5 border-b border-line bg-surface-muted'>
              <div className='col-span-10 ml-5'>部署名</div>
              <div className='mr-5'>操作</div>
            </div>
            <div>
              {
                departments?.map((department: any, index: number) => {
                  return (
                    <div key={index} className='grid grid-cols-12 gap-2 mt-3 mb-3'>
                      <div className={`${calcDepth(department.depth)} grid-cols-9 col-span-9 break-words`}>
                        <span className='text-wrap'>{department.name}</span>
                      </div>
                      <Button
                        variant='outline'
                        size='sm'
                        className='col-span-1'
                        onClick={() => {
                          setTargetDepartment(department)
                          setCreateOpenModal(true)}
                        }>
                        追加
                      </Button>
                      <Button
                        variant='success'
                        size='sm'
                        className='col-span-1'
                        onClick={() => {
                          setTargetDepartment(department)
                          setUpdateOpenModal(true)}
                        }>
                        編集
                      </Button>
                      <Button
                        variant='danger'
                        size='sm'
                        className='col-span-1'
                        onClick={() => {
                          setTargetDepartment(department)
                          setDeleteOpenModal(true)}
                        }>
                        削除
                      </Button>
                    </div>
                  )
                })
              }
            </div>
          </Card>
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
        <ModalShell>
          <div className="p-4 md:p-5 space-y-4">
            <p>本当に削除しますか？</p>
          </div>
          <div className="flex justify-center p-4 md:p-5 border-t border-line rounded-b">
            <Button variant="danger" onClick={() => deleteDepartment(targetDepartment.ID)}>
              削除
            </Button>
            <Button variant="secondary" className="ms-3" onClick={() => setDeleteOpenModal(false)}>
              キャンセル
            </Button>
          </div>
        </ModalShell>
      )}
    </main>  
  )
}
  
export default Departments