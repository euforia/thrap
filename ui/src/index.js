import React from 'react';

import ReactDOM from 'react-dom';
import { BrowserRouter } from 'react-router-dom';

import './index.css';
import App from './App';
// import registerServiceWorker from './registerServiceWorker';

import { MuiThemeProvider, createMuiTheme } from '@material-ui/core';

const theme = createMuiTheme({
    palette: {
        type: 'dark', 
        primary: { 
            main: '#6958a0',
        },
        //secondary: { main: '#11cb5f' }
    },
    overrides: {
        MuiInput: {
            underline: {
                color: '#abbbc6',
                '&:before': {
                    borderBottom: '1px solid #668295',// when input is not touched
                },
            },
        },
    },
});

ReactDOM.render((
    <BrowserRouter basename={'/ui'}>
        <MuiThemeProvider theme={theme}>
            <App />
        </MuiThemeProvider>
    </BrowserRouter>
    ), document.getElementById('root'));

// registerServiceWorker();
