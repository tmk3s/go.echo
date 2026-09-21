"use client";

import { useForm } from "react-hook-form";
import { useState, useEffect } from 'react';
import RootLayout from '@/components/RootLayout';
import Button from '@/components/ui/Button';
import LogoutButton from '@/components/LogoutButton';
import newApiInstance from "../api"

interface Todo {
  ID: number;
  userId: number;
  title: string;
  completed: boolean;
}

const EditModal = ({obj, onSubmit, setOpenModal, deleteTodo}: {obj: any, onSubmit: any, setOpenModal: any, deleteTodo: any} ): React.ReactNode => {
  const {
    register,
    handleSubmit,
    watch,
    formState: { errors }
  } = useForm({
    defaultValues: {
      id: obj?.id,
      title: obj?.title,
    }
  });

  const btnName: string = obj ? 'Done' : 'Add'

  return (
    <div className="overflow-y-auto overflow-x-hidden fixed top-0 right-0 left-0 z-50 justify-center items-center w-full md:inset-0 h-[calc(100%-1rem)] max-h-full">
      <form
        className="max-w-sm mx-auto"
        onSubmit={handleSubmit((data) => {
          onSubmit({id: data.id, title: data.title});
        })}
      >
        <div className="relative p-4 w-full max-w-2xl max-h-full">
          {/* <!-- Modal content --> */}
          <div className="relative bg-surface rounded-lg shadow">
            {/* <!-- Modal body --> */}
            <div className="p-4 md:p-5 space-y-4">
              <input type="hidden" id="id" {...register("id")} />
              <textarea
                id="title"
                className="block p-2.5 w-full text-sm text-body bg-surface-muted rounded-lg border border-line focus:ring-brand focus:border-brand"
                required
                {...register("title")}
              />
            </div>
            {/* <!-- Modal footer --> */}
            <div className="flex items-center p-4 md:p-5 border-t border-line rounded-b">
              <Button type="submit">
                {btnName}
              </Button>
              <Button variant="secondary" className="ms-3" onClick={() => setOpenModal(false)}>
                Cancel
              </Button>
              {
                obj && (
                  <Button variant="danger" className="ms-3" onClick={() => deleteTodo(obj)}>
                    Delete
                  </Button>
                )
              }
            </div>
          </div>
        </div>
      </form>
    </div>
  )
}

const Todos =  () => {
  const [openAddModal, setOpenAddModal] = useState(false);
  const [openCompleteModal, setOpenCompleteModal] = useState(false);
  const [todo, setTodo] = useState({});
  const [todos, setTodos] = useState<any>() // useState<Array<Todo>>();
  const [login, setLogin] = useState(false);

  // 401（セッション切れ）はapi.tsxの共通インターセプターがトップページへリダイレクトする
  const unauthorized = (e: any) => {
    console.error(e);
  }

  const fetchTodos = async () => {
    try {
      const api = newApiInstance();
      const response = await api.get('/api/todos')
  
      console.log(response);
      setTodos(response.data);
      setLogin(true);
    } catch (e) {
      unauthorized(e);
    }
  }
  
  const addTodo = async (obj: any) => {
    try {
      const api = newApiInstance();
      const formData = new FormData();
      formData.append("title", obj.title);
      const response = await api.post('/api/todo', formData)
  
      console.log(response);
      setOpenAddModal(false);
      await fetchTodos();
    } catch (e) {
      unauthorized(e);
    }
  }
  
  const completeTodo = async (obj: any) => {
    try {
      const api = newApiInstance();
      const formData = new FormData();
      const response = await api.put(`/api/todo/${obj.id}/completed`, formData)
  
      console.log(response);
      setOpenCompleteModal(false);
      await fetchTodos();
    } catch (e) {
      unauthorized(e);
    }
  }
  
  const deleteTodo = async (obj: any) => {
    try {
      const api = newApiInstance();
      const response = await api.delete(`/api/todo/${obj.id}`)
  
      console.log(response);
      setOpenCompleteModal(false);
      await fetchTodos();
    } catch (e) {
      unauthorized(e);
    }
  }

  useEffect(() => {
    fetchTodos();
  }, [fetchTodos])

  return (
    <RootLayout>
    <main className="min-h-screen flex-col items-center justify-between p-24">
      {/* Modal toggle */}
      { login && (
        <>
          <Button className="mb-16" onClick={() => setOpenAddModal(true)}>
            Add
          </Button>

          <div className="mb-16 w-40">
            <LogoutButton />
          </div>
        </>
      )}
      
      { openAddModal && (
        <EditModal obj={null} onSubmit={addTodo} setOpenModal={setOpenAddModal} deleteTodo={deleteTodo}/>
      )}
      <div className="grid lg:grid-cols-4 md:grid-cols-3 xs:grid-cols-2 gap-4 text-center">
        {
          todos?.map((todo: Todo) => {
            return (
              <div 
                key={todo.ID}
                className="relative block group min-h-72 p-6 bg-surface border border-line rounded-lg shadow hover:bg-surface-muted"
                onClick={() => {
                  setTodo({id: todo.ID, title: todo.title})
                  setOpenCompleteModal(true)
                }}
              >
                {/* https://qiita.com/yuji38kwmt/items/ba8d59eb0abef1956bae relativeではbreak-words効かない*/}
                {/* calsの部分は親要素からpadding分引いたものを指定している */}
                <p className="absolute w-max max-w-[calc(100%_-_48px)] break-words line-clamp-[10] font-normal text-body">
                  {todo.title}
                </p>
                {todo.completed && (
                    <img
                      className="absolute inset-0"
                      src="/done-256.svg"
                      alt="Next.js Logo"
                      width={75}
                      height={37}
                    />
                )}
              </div>
            )
          })
        }
       </div>
      { openCompleteModal && (
        <EditModal obj={todo} onSubmit={completeTodo} setOpenModal={setOpenCompleteModal} deleteTodo={deleteTodo}/>
      )}
    </main>
    </RootLayout>
  );
}

export default Todos