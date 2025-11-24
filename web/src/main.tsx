import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import { QueryClientProvider, QueryClient } from '@tanstack/react-query'
import { ReactQueryDevtools } from '@tanstack/react-query-devtools'
import { Toaster } from "react-hot-toast"
import { createBrowserRouter, RouterProvider } from 'react-router-dom';

import * as Sentry from "@sentry/react";
import Home from './pages/Index'
import RootAppLayout from './pages/app/layout'

Sentry.init({
  dsn: "https://412466daced4b1d85ee040eef66efc95@o4510248538472448.ingest.us.sentry.io/4510248563113984",
  sendDefaultPii: true,
  environment: import.meta.env.MODE
})

const queryClient = new QueryClient();

const router = createBrowserRouter([
  {
    index: true,
    element: <Home />
  },
  {
    path: "app",
    element: <RootAppLayout />,
    children: [
    ]
  }
]);

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
      <Toaster />
      {import.meta.env.DEV && <ReactQueryDevtools initialIsOpen={false} />}
    </QueryClientProvider>
  </StrictMode>,
)
