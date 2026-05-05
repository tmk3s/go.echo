"use client";

import { useEffect, useState, useMemo } from 'react';
import { useParams, useRouter } from 'next/navigation';
import RootLayout from '@/components/RootLayout';
import useApi from '@/app/api';

type Employee = {
  ID: number;
  last_name: string;
  first_name: string;
  last_name_kana: string;
  first_name_kana: string;
  email: string;
  staff_code: string;
};

type EmployeeAddress = {
  post_code: string | null;
  prefecture_name: string | null;
  city: string | null;
  address_line1: string | null;
  address_line2: string | null;
  tel: string | null;
};

type EmployeeTenure = {
  ID: number;
  joined_on: string;
  resignation_on: string | null;
  resignation_type: string | null;
  status: string | null;
};

type Department = {
  ID: number;
  name: string;
  depth: number;
};

type EmployeeDetail = {
  Employee: Employee;
  Address: EmployeeAddress | null;
  Tenures: EmployeeTenure[];
  Departments: Department[];
};

// order_no順の全部署リストから id → "親 / 子" ラベルのマップを生成
function buildDeptLabelMap(allDepts: Department[]): Map<number, string> {
  const stack: string[] = [];
  const map = new Map<number, string>();
  for (const dept of allDepts) {
    stack.length = dept.depth;
    stack.push(dept.name);
    map.set(dept.ID, stack.join(' / '));
  }
  return map;
}

const Row = ({ label, value }: { label: string; value: string | null | undefined }) => (
  <tr className='border-b dark:border-gray-700'>
    <th className='px-6 py-3 w-48 bg-gray-50 dark:bg-gray-700 font-medium text-gray-700 dark:text-gray-300 text-sm'>
      {label}
    </th>
    <td className='px-6 py-3 text-sm text-gray-900 dark:text-white'>
      {value ?? '—'}
    </td>
  </tr>
);

const Section = ({ title, children }: { title: string; children: React.ReactNode }) => (
  <div className='mt-6'>
    <h2 className='text-lg font-semibold mb-2 text-gray-700 dark:text-gray-300'>{title}</h2>
    <div className='bg-white rounded-lg shadow dark:bg-gray-800 overflow-hidden'>
      {children}
    </div>
  </div>
);

const formatDate = (iso: string | null | undefined) => {
  if (!iso) return null;
  return new Date(iso).toLocaleDateString('ja-JP');
};

const EmployeeDetailPage = () => {
  const { id } = useParams();
  const router = useRouter();
  const api = useApi();
  const [detail, setDetail] = useState<EmployeeDetail | null>(null);
  const [allDepts, setAllDepts] = useState<Department[]>([]);

  useEffect(() => {
    Promise.all([
      api.get(`/api/employees/${id}`),
      api.get('/api/departments'),
    ]).then(([detailRes, deptRes]) => {
      setDetail(detailRes.data);
      setAllDepts(deptRes.data ?? []);
    });
  }, [id]);

  const deptLabelMap = useMemo(() => buildDeptLabelMap(allDepts), [allDepts]);

  if (!detail) {
    return (
      <RootLayout>
        <p className='text-gray-500'>読み込み中...</p>
      </RootLayout>
    );
  }

  const { Employee: emp, Address, Tenures, Departments } = detail;

  return (
    <RootLayout>
      <div className='flex items-center justify-between'>
        <h1 className='text-3xl font-bold'>社員詳細</h1>
        <div className='flex gap-2'>
          <button
            onClick={() => router.push(`/employees/${id}/edit`)}
            className='text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 dark:bg-blue-600 dark:hover:bg-blue-700'
          >
            編集
          </button>
          <button
            onClick={() => router.push('/employees')}
            className='text-gray-700 bg-white border border-gray-300 hover:bg-gray-100 font-medium rounded-lg text-sm px-5 py-2.5 dark:bg-gray-800 dark:text-gray-400 dark:border-gray-600 dark:hover:bg-gray-700'
          >
            一覧に戻る
          </button>
        </div>
      </div>

      <Section title='基本情報'>
        <table className='w-full'>
          <tbody>
            <Row label='スタッフコード' value={emp.staff_code} />
            <Row label='氏名' value={`${emp.last_name} ${emp.first_name}`} />
            <Row label='氏名（カナ）' value={`${emp.last_name_kana} ${emp.first_name_kana}`} />
            <Row label='メールアドレス' value={emp.email} />
          </tbody>
        </table>
      </Section>

      <Section title='住所情報'>
        {Address ? (
          <table className='w-full'>
            <tbody>
              <Row label='郵便番号' value={Address.post_code} />
              <Row label='都道府県' value={Address.prefecture_name} />
              <Row label='市区町村' value={Address.city} />
              <Row label='住所1' value={Address.address_line1} />
              <Row label='住所2' value={Address.address_line2} />
              <Row label='電話番号' value={Address.tel} />
            </tbody>
          </table>
        ) : (
          <p className='px-6 py-4 text-sm text-gray-500'>登録なし</p>
        )}
      </Section>

      <Section title='在籍情報'>
        {Tenures.length > 0 ? (
          <table className='w-full text-sm text-left text-gray-500 dark:text-gray-400'>
            <thead className='text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400'>
              <tr>
                <th className='px-6 py-3'>入社日</th>
                <th className='px-6 py-3'>退職日</th>
                <th className='px-6 py-3'>退職区分</th>
                <th className='px-6 py-3'>ステータス</th>
              </tr>
            </thead>
            <tbody>
              {Tenures.map((tenure) => (
                <tr key={tenure.ID} className='border-b dark:border-gray-700'>
                  <td className='px-6 py-3'>{formatDate(tenure.joined_on)}</td>
                  <td className='px-6 py-3'>{formatDate(tenure.resignation_on) ?? '—'}</td>
                  <td className='px-6 py-3'>{tenure.resignation_type ?? '—'}</td>
                  <td className='px-6 py-3'>{tenure.status ?? '—'}</td>
                </tr>
              ))}
            </tbody>
          </table>
        ) : (
          <p className='px-6 py-4 text-sm text-gray-500'>登録なし</p>
        )}
      </Section>

      <Section title='所属部署'>
        {Departments.length > 0 ? (
          <table className='w-full text-sm text-left text-gray-500 dark:text-gray-400'>
            <thead className='text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400'>
              <tr>
                <th className='px-6 py-3'>部署名</th>
              </tr>
            </thead>
            <tbody>
              {Departments.map((dept) => (
                <tr key={dept.ID} className='border-b dark:border-gray-700'>
                  <td className='px-6 py-3'>{deptLabelMap.get(dept.ID) ?? dept.name}</td>
                </tr>
              ))}
            </tbody>
          </table>
        ) : (
          <p className='px-6 py-4 text-sm text-gray-500'>登録なし</p>
        )}
      </Section>
    </RootLayout>
  );
};

export default EmployeeDetailPage;
