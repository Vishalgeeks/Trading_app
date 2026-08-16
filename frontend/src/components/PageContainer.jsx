export default function PageContainer({ children, title = 'Henna Booking' }) {
  return (
    <div className="min-h-screen flex flex-col">
      <main className="flex-grow">
        {title && (
          <div className="bg-rose-500 py-12">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
              <h1 className="text-3xl font-bold text-white">{title}</h1>
            </div>
          </div>
        )}
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          {children}
        </div>
      </main>
    </div>
  );
}
