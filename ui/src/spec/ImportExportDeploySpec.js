import React, { Component } from 'react';

import { Button, Tab, Tabs } from '@material-ui/core';
import { withStyles } from '@material-ui/core';
import { CloudDownloadOutlined } from '@material-ui/icons';
import ImportSpec from './Import.js';

const styles = theme => ({
    exportIcon: {
        fontSize: 20,
        marginRight: theme.spacing.unit,
    },
});


const style = ({
    container: {
        padding: '20px 40px', 
        textAlign: "center", 
        borderRadius: '5px',
        minWidth: '400px',
    }
});

class ImportExportDeploySpec extends Component {
    constructor(props) {
        super(props);

        var downloadable = encodeURIComponent(JSON.stringify(props.specification, null, 2));
        this.state = {
            spec: props.specification,
            downloadable: downloadable,
            status: '',
            selectedTab: 0,
        }

        this.tabSelected = this.tabSelected.bind(this);
        this.onImportExportError = this.onImportExportError.bind(this);
    }

    tabSelected = (event, value) => {
        this.setState({ selectedTab: value });
    }

    onImportExportError(error){
        var resp = error.response;
        this.setState({
            status: resp.message,
        });
    }

    render() {
        const { classes } = this.props;

        var body;
        if (this.state.selectedTab === 0) {
            body = (
                <ImportSpec project={this.props.project} 
                    onError={this.onImportExportError}
                    onImportSpec={this.props.onImportSpec}
                />
            );
        } else if (this.state.selectedTab === 1) {
            body = (
                <Button aria-label="Export JSON" color="primary" variant="contained"
                    href={'data:text/json;charset=utf-8,'+ this.state.downloadable} 
                    download={this.props.project.ID + '.json'} >
                    <CloudDownloadOutlined className={classes.exportIcon}/>
                    Export
                </Button>
            );
        }

        return (
            <div className="theme-bg theme-color" style={style.container}>
                <div>
                    <Tabs value={this.state.selectedTab}
                        indicatorColor="primary"
                        onChange={this.tabSelected}
                        centered 
                    >

                        <Tab name="Import" label="Import" />
                        <Tab name="Export" label="Export" />
                    </Tabs>
                </div>
                <div className="error-container">{this.state.status}</div>
                <div style={{padding:'20px', marginTop: '30px'}}>
                    {body}
                </div>
            </div>
        );
    }
}

export default withStyles(styles)(ImportExportDeploySpec);