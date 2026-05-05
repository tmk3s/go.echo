"use client";

import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { useForm } from 'react-hook-form';
import RootLayout from '@/components/RootLayout';
import useApi from '@/app/api';

type FormValues = {
  staff_code: string;
  last_name: string;
  first_name: string;
  last_name_kana: string;
  first_name_kana: string;
  email: string;
  post_code: string;
  prefecture_id: string;
  city: string;
  address_line1: string;
  address_line2: string;
  tel: string;
};

type TenureForm = {
  ID: number;
  joined_on: string;
  resignation_on: string;
  resignation_type: string;
  status: string;
};

type Department = { ID: number; name: string; depth: number };
type DepartmentOption = { id: number; label: string };
type Prefecture = { ID: number; name: string };

// order_no順に並んだ部署リストから depth を使って "親 / 子" 形式のラベルを生成
function buildDepartmentOptions(departments: Department[]): DepartmentOption[] {
  const stack: string[] = [];
  return departments.map((dept) => {
    stack.length = dept.depth;
    stack.push(dept.name);
    return { id: dept.ID, label: stack.join(' / ') };
  });
}

const Section = ({ title, children }: { title: string; children: React.ReactNode }) => (
  <div className='mt-6'>
    <h2 className='text-lg font-semibold mb-2 text-gray-700 dark:text-gray-300'>{title}</h2>
    <div className='bg-white rounded-lg shadow dark:bg-gray-800 p-6'>{children}</div>
  </div>
);

const Field = ({ label, children }: { label: string; children: React.ReactNode }) => (
  <div>
    <label className='block mb-1 text-sm font-medium text-gray-700 dark:text-gray-300'>{label}</label>
    {children}
  </div>
);

const inputClass = 'w-full bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white';

const EmployeeEditPage = () => {
  const { id } = useParams();
  const router = useRouter();
  const api = useApi();

  const [tenures, setTenures] = useState<TenureForm[]>([]);
  const [departmentOptions, setDepartmentOptions] = useState<DepartmentOption[]>([]);
  const [selectedDeptIds, setSelectedDeptIds] = useState<(number | null)[]>([null]);
  const [prefectures, setPrefectures] = useState<Prefecture[]>([]);

  const { register, handleSubmit, reset } = useForm<FormValues>();

  useEffect(() => {
    Promise.all([
      api.get(`/api/employees/${id}`),
      api.get('/api/departments'),
      api.get('/api/prefectures'),
    ]).then(([detailRes, deptRes, prefRes]) => {
      const { Employee: emp, Address, Tenures, Departments } = detailRes.data;

      reset({
        staff_code: emp.staff_code,
        last_name: emp.last_name,
        first_name: emp.first_name,
        last_name_kana: emp.last_name_kana,
        first_name_kana: emp.first_name_kana,
        email: emp.email,
        post_code: Address?.post_code ?? '',
        prefecture_id: Address?.prefecture_id?.toString() ?? '',
        city: Address?.city ?? '',
        address_line1: Address?.address_line1 ?? '',
        address_line2: Address?.address_line2 ?? '',
        tel: Address?.tel ?? '',
      });

      setTenures(
        (Tenures ?? []).map((t: any) => ({
          ID: t.ID,
          joined_on: t.joined_on?.substring(0, 10) ?? '',
          resignation_on: t.resignation_on?.substring(0, 10) ?? '',
          resignation_type: t.resignation_type ?? '',
          status: t.status ?? '',
        }))
      );

      setDepartmentOptions(buildDepartmentOptions(deptRes.data ?? []));
      const assignedIds = (Departments ?? []).map((d: Department) => d.ID as number | null);
      setSelectedDeptIds(assignedIds.length > 0 ? assignedIds : [null]);
      setPrefectures(prefRes.data ?? []);
    });
  }, [id]);

  const onSubmit = handleSubmit(async (data) => {
    await api.put(`/api/employees/${id}`, {
      last_name: data.last_name,
      first_name: data.first_name,
      last_name_kana: data.last_name_kana,
      first_name_kana: data.first_name_kana,
      email: data.email,
      staff_code: data.staff_code,
      post_code: data.post_code || null,
      prefecture_id: data.prefecture_id ? Number(data.prefecture_id) : null,
      city: data.city || null,
      address_line1: data.address_line1 || null,
      address_line2: data.address_line2 || null,
      tel: data.tel || null,
      tenures: tenures.map((t) => ({
        id: t.ID,
        joined_on: t.joined_on,
        resignation_on: t.resignation_on || null,
        resignation_type: t.resignation_type || null,
        status: t.status || null,
      })),
      department_ids: selectedDeptIds.filter((id): id is number => id !== null),
    });
    router.push(`/employees/${id}`);
  });

  const updateTenureField = (idx: number, field: keyof Omit<TenureForm, 'ID'>, value: string) => {
    setTenures((prev) => prev.map((t, i) => i === idx ? { ...t, [field]: value } : t));
  };

  return (
    <RootLayout>
      <form onSubmit={onSubmit}>
        <div className='flex items-center justify-between'>
          <h1 className='text-3xl font-bold'>社員編集</h1>
          <div className='flex gap-2'>
            <button
              type='submit'
              className='text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 dark:bg-blue-600 dark:hover:bg-blue-700'
            >
              保存
            </button>
            <button
              type='button'
              onClick={() => router.push(`/employees/${id}`)}
              className='text-gray-700 bg-white border border-gray-300 hover:bg-gray-100 font-medium rounded-lg text-sm px-5 py-2.5 dark:bg-gray-800 dark:text-gray-400 dark:border-gray-600 dark:hover:bg-gray-700'
            >
              キャンセル
            </button>
          </div>
        </div>

        {/* 基本情報 */}
        <Section title='基本情報'>
          <div className='grid grid-cols-2 gap-4'>
            <Field label='スタッフコード'>
              <input className={inputClass} {...register('staff_code')} />
            </Field>
            <Field label='メールアドレス'>
              <input type='email' className={inputClass} {...register('email')} />
            </Field>
            <Field label='姓'>
              <input className={inputClass} {...register('last_name')} />
            </Field>
            <Field label='名'>
              <input className={inputClass} {...register('first_name')} />
            </Field>
            <Field label='姓（カナ）'>
              <input className={inputClass} {...register('last_name_kana')} />
            </Field>
            <Field label='名（カナ）'>
              <input className={inputClass} {...register('first_name_kana')} />
            </Field>
          </div>
        </Section>

        {/* 住所情報 */}
        <Section title='住所情報'>
          <div className='grid grid-cols-2 gap-4'>
            <Field label='郵便番号'>
              <input className={inputClass} {...register('post_code')} />
            </Field>
            <Field label='電話番号'>
              <input className={inputClass} {...register('tel')} />
            </Field>
            <Field label='都道府県'>
              <select className={inputClass} {...register('prefecture_id')}>
                <option value=''>選択してください</option>
                {prefectures.map((p) => (
                  <option key={p.ID} value={p.ID}>{p.name}</option>
                ))}
              </select>
            </Field>
            <Field label='市区町村'>
              <input className={inputClass} {...register('city')} />
            </Field>
            <Field label='住所1'>
              <input className={inputClass} {...register('address_line1')} />
            </Field>
            <Field label='住所2'>
              <input className={inputClass} {...register('address_line2')} />
            </Field>
          </div>
        </Section>

        {/* 在籍情報 */}
        <Section title='在籍情報'>
          {tenures.length === 0 && <p className='text-sm text-gray-500'>登録なし</p>}
          {tenures.map((tenure, idx) => (
            <div key={tenure.ID} className='mb-6 pb-6 border-b last:border-b-0 dark:border-gray-700'>
              <div className='grid grid-cols-2 gap-4'>
                <Field label='入社日'>
                  <input type='date' className={inputClass} value={tenure.joined_on}
                    onChange={(e) => updateTenureField(idx, 'joined_on', e.target.value)} />
                </Field>
                <Field label='退職日'>
                  <input type='date' className={inputClass} value={tenure.resignation_on}
                    onChange={(e) => updateTenureField(idx, 'resignation_on', e.target.value)} />
                </Field>
                <Field label='退職区分'>
                  <input className={inputClass} value={tenure.resignation_type}
                    onChange={(e) => updateTenureField(idx, 'resignation_type', e.target.value)} />
                </Field>
                <Field label='ステータス'>
                  <input className={inputClass} value={tenure.status}
                    onChange={(e) => updateTenureField(idx, 'status', e.target.value)} />
                </Field>
              </div>
            </div>
          ))}
        </Section>

        {/* 所属部署 */}
        <Section title='所属部署'>
          <div className='space-y-2'>
            {selectedDeptIds.map((deptId, idx) => (
              <div key={idx} className='flex items-center gap-2'>
                <select
                  className={inputClass + ' max-w-sm'}
                  value={deptId ?? ''}
                  onChange={(e) => {
                    const val = e.target.value ? Number(e.target.value) : null;
                    setSelectedDeptIds((prev) => prev.map((v, i) => i === idx ? val : v));
                  }}
                >
                  <option value=''>選択してください</option>
                  {departmentOptions.map((opt) => (
                    <option key={opt.id} value={opt.id}>{opt.label}</option>
                  ))}
                </select>
                {selectedDeptIds.length > 1 && (
                  <button
                    type='button'
                    onClick={() => setSelectedDeptIds((prev) => prev.filter((_, i) => i !== idx))}
                    className='text-red-500 hover:text-red-700 text-lg font-bold px-2'
                  >
                    ×
                  </button>
                )}
              </div>
            ))}
            <button
              type='button'
              onClick={() => setSelectedDeptIds((prev) => [...prev, null])}
              className='mt-2 text-blue-600 hover:text-blue-800 text-sm font-medium flex items-center gap-1'
            >
              ＋ 部署を追加
            </button>
          </div>
        </Section>

        {/* 下部にも保存ボタン */}
        <div className='mt-6 flex justify-end gap-2'>
          <button
            type='submit'
            className='text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 dark:bg-blue-600 dark:hover:bg-blue-700'
          >
            保存
          </button>
          <button
            type='button'
            onClick={() => router.push(`/employees/${id}`)}
            className='text-gray-700 bg-white border border-gray-300 hover:bg-gray-100 font-medium rounded-lg text-sm px-5 py-2.5 dark:bg-gray-800 dark:text-gray-400 dark:border-gray-600 dark:hover:bg-gray-700'
          >
            キャンセル
          </button>
        </div>
      </form>
    </RootLayout>
  );
};

export default EmployeeEditPage;
