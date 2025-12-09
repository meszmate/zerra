import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css';
import "./i18n";
import { QueryClientProvider, QueryClient } from '@tanstack/react-query'
import { Toaster } from "react-hot-toast"
import { createBrowserRouter, RouterProvider } from 'react-router-dom';

import * as Sentry from "@sentry/react";
import Home from './pages/index'
import RootAppLayout from './pages/app/layout'
import AuthLayout from './pages/auth/layout'
import { LoginForm } from './pages/auth/login'
import { RegisterForm } from './pages/auth/register'
import { LoginConfirm } from './pages/auth/login-confirm'
import { RegisterConfirm } from './pages/auth/register-confirm'
import { ResetPasswordForm } from './pages/auth/reset-password'

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
  },
  {
    path: "auth",
    element: <AuthLayout />,
    children: [
      {
        path: "login",
        children: [
          {
            index: true,
            element: <LoginForm />,
          },
          {
            path: "confirm",
            element: <LoginConfirm />,
          }
        ]
      },
      {
        path: "register",
        children: [
          {
            index: true,
            element: <RegisterForm />,
          },
          {
            path: "confirm",
            element: <RegisterConfirm />,
          },
        ]
      },
      {
        path: "reset-password",
        children: [
          {
            index: true,
            element: <ResetPasswordForm />,
          },
        ]
      }
    ]
  }
]);

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
      <Toaster />
    </QueryClientProvider>
  </StrictMode>,
)
