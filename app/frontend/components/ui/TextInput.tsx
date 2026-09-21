import { InputHTMLAttributes, forwardRef } from 'react';

const base =
  'block w-full bg-surface-muted border border-line text-body text-sm rounded-lg p-2.5 ' +
  'placeholder:text-muted focus:ring-brand focus:border-brand';

// select など input 以外の入力要素でも同じ見た目を使うために公開する
export const controlClass = base;

export type TextInputProps = InputHTMLAttributes<HTMLInputElement>;

const TextInput = forwardRef<HTMLInputElement, TextInputProps>(
  ({ className = '', ...props }, ref) => (
    <input ref={ref} className={`${base} ${className}`.trim()} {...props} />
  )
);

TextInput.displayName = 'TextInput';

export default TextInput;

export const Label = ({ children, className = '' }: { children: React.ReactNode; className?: string }) => (
  <label className={`block mb-2 text-sm font-medium text-body ${className}`.trim()}>
    {children}
  </label>
);
