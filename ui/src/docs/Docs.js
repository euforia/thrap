import React, { Component } from 'react';
import { RedocStandalone } from 'redoc';
import thrap from '../api/thrap.js';

class APIDocs extends Component {
    render() {
        return (
            <RedocStandalone 
                specUrl={thrap.addr() + '/swagger.json'}
                options={{
                    hideDownloadButton: true,
                    disableSearch: true,
                    pathInMiddlePanel: true,
                    requiredPropsFirst: true,
                    sortPropsAlphabetically: false,
                }}
            />
        );
    }
}

export default APIDocs;