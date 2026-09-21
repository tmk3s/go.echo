"use client";

import { useState } from 'react';
import RootLayout from '@/components/RootLayout';
import Button from '@/components/ui/Button';
import Card from '@/components/ui/Card';
import useApi from '@/app/api';

const DepartmentsCsv = () => {
  const api = useApi();
  const [file, setFile] = useState<File | null>(null);
  const [uploading, setUploading] = useState(false);
  const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const selected = e.target.files?.[0] ?? null;
    setFile(selected);
    setMessage(null);
  };

  const handleUpload = async () => {
    if (!file) return;
    const formData = new FormData();
    formData.append('file', file);

    setUploading(true);
    setMessage(null);
    try {
      await api.post('/api/departments/upload', formData, {
        headers: { 'Content-Type': 'multipart/form-data' },
      });
      setMessage({ type: 'success', text: 'インポートが完了しました' });
      setFile(null);
    } catch (e: any) {
      setMessage({ type: 'error', text: e.response?.data?.message ?? 'エラーが発生しました' });
    } finally {
      setUploading(false);
    }
  };

  return (
    <main>
      <RootLayout>
        <>
          <h1 className='text-3xl font-bold'>部署CSV</h1>
          <Card className='mt-20'>
            <div className={`p-5`}>
              <div className='mt-10'>
                <h2 className='font-bold'>CSVアップロード</h2>
                <p className='mb-2 text-sm text-muted'>
                  1行目にヘッダー <code>name</code> を持つCSVをアップロードしてください。<br />
                  CSVに存在しない部署名は新規作成されます。
                </p>

                <div className="flex items-center justify-center w-full mt-4">
                  <label htmlFor="dropzone-file" className="flex flex-col items-center justify-center w-full h-48 border-2 border-line border-dashed rounded-lg cursor-pointer bg-surface-muted hover:bg-canvas">
                    <div className="flex flex-col items-center justify-center pt-5 pb-6">
                      <svg className="w-8 h-8 mb-4 text-muted" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 20 16">
                        <path stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M13 13h3a3 3 0 0 0 0-6h-.025A5.56 5.56 0 0 0 16 6.5 5.5 5.5 0 0 0 5.207 5.021C5.137 5.017 5.071 5 5 5a4 4 0 0 0 0 8h2.167M10 15V6m0 0L8 8m2-2 2 2"/>
                      </svg>
                      {file
                        ? <p className="text-sm text-body font-semibold">{file.name}</p>
                        : <p className="text-sm text-muted"><span className="font-semibold">クリックして選択</span>またはドラッグ＆ドロップ</p>
                      }
                      <p className="text-xs text-muted mt-1">.csv のみ</p>
                    </div>
                    <input id="dropzone-file" type="file" accept=".csv" className="hidden" onChange={handleFileChange} />
                  </label>
                </div>

                {message && (
                  <p className={`mt-3 text-sm ${message.type === 'success' ? 'text-success' : 'text-danger'}`}>
                    {message.text}
                  </p>
                )}

                <Button className='mt-4' onClick={handleUpload} disabled={!file || uploading}>
                  {uploading ? 'アップロード中...' : 'インポート実行'}
                </Button>
              </div>
            </div>
          </Card>
        </>
      </RootLayout>
    </main>
  );
};

export default DepartmentsCsv;
