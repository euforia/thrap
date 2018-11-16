import React, { Component } from 'react';
import { RedocStandalone } from 'redoc';
import thrap from './thrap.js';

class APIDocs extends Component {
    render() {
        
        return (
            <RedocStandalone 
                specUrl={thrap.addr() + '/swagger.json'}
                options={{
                    hideDownloadButton: true,
                    disableSearch: true,
                    pathInMiddlePanel: true,
                }}
            />
        );
    }
}

export default APIDocs;