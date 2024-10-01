import React, {useState, useEffect, useMemo} from "react";
import { Outlet } from "react-router-dom";

import './App.css'

function App() {
    return (
        <>
            <Outlet /> 
        </>
    );
}

export default App;
