import { ButtonHTMLAttributes, forwardRef } from 'react';

type Variant = 'primary' | 'secondary' | 'outline' | 'success' | 'danger';
type Size = 'sm' | 'md';

const base =
  'inline-flex items-center justify-center gap-1.5 font-medium rounded-lg text-sm ' +
  'focus:outline-none focus:ring-4 disabled:opacity-50 disabled:cursor-not-allowed';

const sizes: Record<Size, string> = {
  sm: 'px-3 py-1.5',
  md: 'px-5 py-2.5',
};

const variants: Record<Variant, string> = {
  primary: 'text-brand-contrast bg-brand hover:bg-brand-hover focus:ring-brand/30',
  secondary: 'text-body bg-surface border border-line hover:bg-surface-muted focus:ring-line/50',
  outline: 'text-brand border border-brand hover:bg-brand hover:text-brand-contrast focus:ring-brand/30',
  success: 'text-success border border-success hover:bg-success hover:text-white focus:ring-success/30',
  danger: 'text-danger border border-danger hover:bg-danger hover:text-white focus:ring-danger/30',
};

export type ButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: Variant;
  size?: Size;
};

// アプリ内のボタンはすべてこれを使う。
// 以前は同じクラス文字列が各ページにコピーされており、focus リングの有無などが
// 少しずつ食い違っていた。
const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  ({ variant = 'primary', size = 'md', type = 'button', className = '', ...props }, ref) => (
    <button
      ref={ref}
      type={type}
      className={`${base} ${sizes[size]} ${variants[variant]} ${className}`.trim()}
      {...props}
    />
  )
);

Button.displayName = 'Button';

export default Button;
