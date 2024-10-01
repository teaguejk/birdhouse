import React from 'react';

import { NavLink, Outlet, useLocation } from 'react-router-dom';


const Layout: React.FC = () => {
    return (
        <div>
            <Outlet />
        </div>
    );
};

export default Layout;