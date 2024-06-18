"use client";

import Axios from 'axios';
import { useForm } from "react-hook-form";
import { useState, useEffect } from 'react';
import newApiInstance from "../api"

interface Todo {
  id: number;
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
          <div className="relative bg-white rounded-lg shadow dark:bg-gray-700">
            {/* <!-- Modal body --> */}
            <div className="p-4 md:p-5 space-y-4">
              <input type="hidden" id="id" {...register("id")} />
              <textarea
                id="title"
                className="block p-2.5 w-full text-sm text-gray-900 bg-gray-50 rounded-lg border border-gray-300 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                required
                {...register("title")}
              />
            </div>
            {/* <!-- Modal footer --> */}
            <div className="flex items-center p-4 md:p-5 border-t border-gray-200 rounded-b dark:border-gray-600">
              <button type="submit" className="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800">
                {btnName}
              </button>
              <button type="button" onClick={() => setOpenModal(false)} className="py-2.5 px-5 ms-3 text-sm font-medium text-gray-900 focus:outline-none bg-white rounded-lg border border-gray-200 hover:bg-gray-100 hover:text-blue-700 focus:z-10 focus:ring-4 focus:ring-gray-100 dark:focus:ring-gray-700 dark:bg-gray-800 dark:text-gray-400 dark:border-gray-600 dark:hover:text-white dark:hover:bg-gray-700">
                Cancel
              </button>
              {
                obj && (
                  <button type="button" onClick={() => deleteTodo(obj)} className="py-2.5 px-5 ms-3 text-sm font-medium text-gray-900 focus:outline-none bg-white rounded-lg border border-gray-200 hover:bg-gray-100 hover:text-blue-700 focus:z-10 focus:ring-4 focus:ring-red-100 dark:focus:ring-red-700 dark:bg-red-800 dark:text-red-400 dark:border-red-600 dark:hover:text-white dark:hover:bg-red-700">
                    Delete
                  </button>
                )
              }
            </div>
          </div>
        </div>
      </form>
    </div>
  )
}

export default () => {
  const [openAddModal, setOpenAddModal] = useState(false);
  const [openCompleteModal, setOpenCompleteModal] = useState(false);
  const [todo, setTodo] = useState({});
  const [todos, setTodos] = useState<any>() // useState<Array<Todo>>();
  const [login, setLogin] = useState(false);

  const unauthorized = (e: any) => {
    console.error(e);
    if (Axios.isAxiosError(e) && e.response && e.response.status === 401) {
      console.log('401 Error!!');
      console.log(e.message); //Axiosの例外オブジェクトとして扱える
      localStorage.setItem("unauthorizedError", "ログインしてください");
      window.location.href = '/sign_in'
    }
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
      setTodos(response.data);
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
      setTodos(response.data);
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
      setTodos(response.data);
    } catch (e) {
      unauthorized(e);
    }
  }

  useEffect(() => {
    fetchTodos();
  }, [])

  return (
    <main className="min-h-screen flex-col items-center justify-between p-24">
      {/* Modal toggle */}
      { login && (
        <button onClick={() => setOpenAddModal(true)} className="mb-16 block text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800" type="button">
          Add
        </button>
      )}
      
      { openAddModal && (
        <EditModal obj={null} onSubmit={addTodo} setOpenModal={setOpenAddModal} deleteTodo={deleteTodo}/>
      )}
      <div className="grid lg:grid-cols-4 md:grid-cols-3 xs:grid-cols-2 gap-4 text-center">
        {
          todos?.map((todo: Todo) => {
            return (
              <div 
                key={todo.id}
                className="relative block group min-h-72 p-6 bg-white border border-gray-200 rounded-lg shadow hover:bg-gray-100 dark:bg-gray-800 dark:border-gray-700 dark:hover:bg-gray-700"
                onClick={() => {
                  setTodo({id: todo.id, title: todo.title})
                  setOpenCompleteModal(true)
                }}
              >
                {/* https://qiita.com/yuji38kwmt/items/ba8d59eb0abef1956bae relativeではbreak-words効かない*/}
                {/* calsの部分は親要素からpadding分引いたものを指定している */}
                <p className="absolute w-max max-w-[calc(100%_-_48px)] break-words line-clamp-[10] font-normal text-gray-700 dark:text-gray-400">
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
  );
}
