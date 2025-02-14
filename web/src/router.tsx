import * as React from "react";
import { createBrowserRouter, redirect } from "react-router-dom";

import App from "./App"

import Home from "./routes/home/Home";

// import { getCookie, setCookie } from "./utils/Cookies";

import Layout from "./layout/Layout";

const routes = [
    {
        path: "/",
        element: <App />,

        // errorElement: <ErrorPage />,
        children: [
            {
                element: <Layout />,
                children: [
                    {
                        element: <Home />,
                        index: true,
                    },
                ]
            },
        ]
    },
]

export const router = createBrowserRouter(routes);
