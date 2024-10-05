import React from 'react';

const Home: React.FC = () => {
    return (
        <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', height: '100vh' }}>
            <h1>Welcome to the Home page!</h1>
            <p>This is a basic React page.</p>
        </div>
    );
};

export default Home;