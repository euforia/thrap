import React from 'react';
import ReactDOM from 'react-dom';

import { BrowserRouter } from 'react-router-dom';
import { MuiThemeProvider, createMuiTheme } from '@material-ui/core';
// import blue from '@material-ui/core/colors/blue';

import './index.css';
import App from './App';

const basePath = "/ui/";

const theme = createMuiTheme({
    // palette: {
        // primary: { 
        //     main: '#6958a0',
        // },
        // secondary: {
        //     main: '#ec6565',
        // }
    // },
    // anchor: {
		// main: blue[500],
		// selected: blue[700]
	// },
});

ReactDOM.render(
    <BrowserRouter basename={basePath}>
        <MuiThemeProvider theme={theme}>
            <App />
        </MuiThemeProvider>
    </BrowserRouter>, 
    document.getElementById('root')
);

