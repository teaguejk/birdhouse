import React from 'react'
import { RouterProvider } from "react-router-dom";
import ReactDOM from 'react-dom/client'
import { ChakraProvider } from '@chakra-ui/react'
import { ToastContainer } from 'react-toastify';

import { router } from './router.tsx'
import theme from './theme.ts'

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <meta name="viewport" content="width=device-width, initial-scale=1" /> 
    {localStorage.getItem('chakra-ui-color-mode')? <></>:<>{localStorage.setItem('chakra-ui-color-mode', 'dark')}</>}
    <ChakraProvider theme={theme}>
      <RouterProvider router={router} />
      <ToastContainer
                position="bottom-right"
                autoClose={5000}
                hideProgressBar={false}
                newestOnTop={false}
                closeOnClick
                rtl={false}
                pauseOnFocusLoss
                draggable
                pauseOnHover
                theme="light"
            />    
    </ChakraProvider>
  </React.StrictMode>,
)
