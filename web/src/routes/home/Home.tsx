import React from 'react';
import { useEffect, useState } from 'react';

import { 
    Text,
    Flex, 
    Spacer,
    Grid 
} from '@chakra-ui/react';

const Home: React.FC = () => {

    return (
        <Flex> {/* style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', height: '100vh' }}> */}
            <Spacer />
            <Flex p='6' direction='column' textAlign='center'>
                <Text>Latest image from the Birdhouse:</Text>
                <img 
                    src={`http://localhost:8080/v1/images/latest`} 
                    alt="Latest image from the Birdhouse" 
                    style={{ width: '100%', maxWidth: '600px', height: 'auto' }} 
                />
            </Flex>
        </Flex>
    );
};

export default Home;