import { extendTheme } from "@chakra-ui/react";
import { mode } from '@chakra-ui/theme-tools'

const theme = extendTheme({
    styles: {
        global: (props: any) => ({
            body: {
                color: mode('gray.800', '#dadada')(props),
                bg: mode('white', 'dark.main')(props),
            },
            '*::placeholder': {
                color: mode('gray.400', 'gray.200')(props),
            },
            '*, *::before, &::after': {
                // borderColor: 'gray.400',
                wordWrap: 'break-word',
            },
        }),
    },
    colors: {
        dark: {
            //main: "#595b5e"
            main: "#525252"
        },
    },
    components: {
        Drawer: {
            baseStyle: (props: any) => ({
                dialog: {
                    // ...props.theme.components.Drawer.baseStyle.dialog,
                    bg: mode('white', 'dark.main')(props),
                    border: "1px solid",
                    borderRadius: 'md',
                    boxShadow: "lg",
                },
            }),
        },
        Modal: {
            baseStyle: (props: any) => ({
                dialog: {
                    bg: mode('white', 'dark.main')(props),
                    border: "1px solid",
                    borderRadius: 'md',
                    boxShadow: "lg",
                },
            }),
        },
        Card: {
            baseStyle: (props: any) => ({
                container: {
                    bg: mode('white', '#626262')(props),
                    borderRadius: 'lg',
                    _hover: {
                        bg: mode('#D3D3D3', '#444')(props),
                    },
                },
            }),
        },
        Tabs: {
            baseStyle: (props: any) => ({
                tab: {
                    boxShadow: "none",
                    outline: 'none',
                    border: 'none',
                    mb: '-2px',
                    _focus: {
                        bg: 'none',
                        boxShadow: "none",
                        outline: 'none',
                        border: 'none',
                        mb: '-2px',
                    },
                    _selected: {
                        color: 'orange.200',
                        bg: 'none',
                        boxShadow: "none",
                        outline: 'none',
                        border: 'none',
                        mb: '-2px',
                    },
                },
            }),
        },
    },
});

export default theme;
