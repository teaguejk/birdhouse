import React from 'react';

import { NavLink, Outlet, useLocation } from 'react-router-dom';

import { Box, Flex, Text } from '@chakra-ui/react';

const Layout: React.FC = () => {
    return (
        <Box w='100vw'>
			<Flex p="4" bg="#333" color="white" w="100%" minW='100%'>
                <Text>Birdhouse</Text>
            </Flex>
            <Outlet />
        </Box>
    );
};

export default Layout;