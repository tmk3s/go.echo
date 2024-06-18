"use client";

import axios, {AxiosInstance} from 'axios';

function getCookieValue(key: string): string {
  const cookies = document.cookie.split(';')
  const foundCookie = cookies.find(
    (cookie) => cookie.split('=')[0].trim() === key.trim()
  )
  if (foundCookie) {
    const cookieValue = decodeURIComponent(foundCookie.split('=')[1])
    return cookieValue
  }
  return ''
}

export default (): AxiosInstance  => {
  const session: string = getCookieValue('session')
  const instance = axios.create({
    baseURL: 'http://localhost:1323',
    headers: {
      'Content-Type': 'application/json',
      'X-Requested-With': 'XMLHttpRequest',
      'Authorization': `Bearer ${session}`,
    },
    responseType: 'json'  
  })
  return instance;
}