import { ReactNode } from 'react';
import Card from './Card';

export type RowProps = {
  label: string;
  value: ReactNode;
};

// 「項目名 / 値」を並べる表の1行。
// マイページ・会社情報・社員詳細で同じ定義が3つ重複していたため共通化した。
export const Row = ({ label, value }: RowProps) => (
  <tr className='border-b border-line last:border-b-0'>
    <th className='px-6 py-3 w-48 bg-surface-muted font-medium text-body text-sm text-left'>
      {label}
    </th>
    <td className='px-6 py-3 text-sm text-body'>
      {value ?? '—'}
    </td>
  </tr>
);

export type DescriptionTableProps = {
  children: ReactNode;
  className?: string;
};

const DescriptionTable = ({ children, className = '' }: DescriptionTableProps) => (
  <Card className={className}>
    <table className='w-full'>
      <tbody>{children}</tbody>
    </table>
  </Card>
);

export default DescriptionTable;
