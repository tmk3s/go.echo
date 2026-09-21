const ErrorToast = ({message, clostToast}: {message: string, clostToast: any}) => {
  setTimeout(() => {
    clostToast()
  }, 5000);

  return (
    <div id="toast-warning" className="flex absolute top-0 right-0 items-center w-full max-w-xs p-4 text-muted bg-surface rounded-lg shadow" role="alert">
      <div className="inline-flex items-center justify-center flex-shrink-0 w-8 h-8 text-orange-500 bg-orange-100 rounded-lg dark:bg-orange-700 dark:text-orange-200">
        <svg className="w-5 h-5" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="currentColor" viewBox="0 0 20 20">
          <path d="M10 .5a9.5 9.5 0 1 0 9.5 9.5A9.51 9.51 0 0 0 10 .5ZM10 15a1 1 0 1 1 0-2 1 1 0 0 1 0 2Zm1-4a1 1 0 0 1-2 0V6a1 1 0 0 1 2 0v5Z" />
        </svg>
        <span className="sr-only">Warning icon</span>
      </div>
      <div className="ms-3 text-sm font-normal">{message}</div>
      <button
        type="button"
        className="ms-auto -mx-1.5 -my-1.5 bg-surface text-muted hover:text-body rounded-lg focus:ring-2 focus:ring-line p-1.5 hover:bg-surface-muted inline-flex items-center justify-center h-8 w-8"
        data-dismiss-target="#toast-warning"
        aria-label="Close"
        onClick={() => {clostToast()}}
      >
        <span className="sr-only">Close</span>
        <svg className="w-3 h-3" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 14 14">
          <path stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="m1 1 6 6m0 0 6 6M7 7l6-6M7 7l-6 6" />
        </svg>
      </button>
    </div>
  )
}

export default ErrorToast;