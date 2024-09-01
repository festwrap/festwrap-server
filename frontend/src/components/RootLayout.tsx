export default function RootLayout({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  return (
    <div className="flex flex-col h-full items-center justify-center bg-gray-50">
      {children}
    </div>
  )
}
