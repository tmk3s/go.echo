import Link from "next/link";

const cardClass =
  "group rounded-lg border border-transparent px-5 py-4 transition-colors hover:border-line hover:bg-surface";

export default function Home() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-between p-24">

      <div className="mb-32 grid text-center lg:mb-0 lg:w-full lg:max-w-5xl lg:grid-cols-1 lg:text-left">
        <Link href="/sign_in" className={cardClass}>
          <h2 className="mb-3 text-2xl font-semibold">
            SignIn
            <span className="inline-block transition-transform group-hover:translate-x-1 motion-reduce:transform-none">
              -&gt;
            </span>
          </h2>
          <p className="m-0 max-w-[30ch] text-sm opacity-50">
            ログイン画面へ
          </p>
        </Link>

        <Link href="/company_registration" className={cardClass}>
          <h2 className="mb-3 text-2xl font-semibold">
            会社登録{" "}
            <span className="inline-block transition-transform group-hover:translate-x-1 motion-reduce:transform-none">
              -&gt;
            </span>
          </h2>
          <p className="m-0 max-w-[30ch] text-sm opacity-50">
            会社とユーザーを新規登録
          </p>
        </Link>
      </div>
    </main>
  );
}
