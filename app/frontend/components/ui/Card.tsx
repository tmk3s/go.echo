import { HTMLAttributes } from 'react';

export type CardProps = HTMLAttributes<HTMLDivElement>;

// カード(白い面)の共通ラッパー。ページ背景(canvas)より一段明るい面を作る。
const Card = ({ className = '', ...props }: CardProps) => (
  <div className={`bg-surface rounded-lg shadow overflow-hidden ${className}`.trim()} {...props} />
);

export default Card;
