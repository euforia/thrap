import React, { Component } from 'react';
import { RedocStandalone } from 'redoc';

class APIDocs extends Component {
    render() {
        
        return (
            <RedocStandalone 
                specUrl='/swagger.json'
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