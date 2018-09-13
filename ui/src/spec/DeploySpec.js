import React, { Component } from 'react';
import ReactJson from 'react-json-view';

import thrap from '../thrap.js';
import { Grid, Modal } from '@material-ui/core';
import SpecControls from '../deployments/SpecControls.js';
import ImportExportDeploySpec from './ImportExportDeploySpec.js';

class DeploySpec extends Component {
    constructor(props) {
        super(props);

        this.state = {
            project: props.project,
            spec: props.specification,
            showModal: false,
        }

        this.onEdit = this.onEdit.bind(this);
        this.openModal = this.openModal.bind(this);
        this.closeModal = this.closeModal.bind(this);
        this.onImportSpec = this.onImportSpec.bind(this);

        if (props.specification === undefined) {
            this.fetchSpec()
        }
    }

    fetchSpec() {
        thrap.deploymentSpec(this.state.project.ID)
            .then(({data}) => {
                this.setState({
                    spec: data,
                });
            })
            .catch(error => {
                console.log("ERROR", error.config);
            });
    }

    onEdit(event) {
        console.log(event);
        return true;
    }

    closeModal() {
        this.setState({showModal:false});
    }

    openModal() {
        this.setState({showModal:true});
    }

    onImportSpec(data) {
        this.setState({
            showModal:false,
            spec: data,
        });
    }

    render() {
        return (
            <div>
                <Grid container spacing={0}>
                    <Grid item xs={2}>
                        <SpecControls onClose={this.props.onCloseDialogue} onImportExport={this.openModal} />
                    </Grid>
                    <Grid item xs={10}>
                        <div className="header-container">
                            <div className="header-title">Specification</div>
                            <div className="subscript">Version: </div>
                        </div>
                        <div style={{padding: '20px 0 40px 0'}}>
                            <ReactJson src={this.state.spec} name="spec" style={{background: 'none'}}
                                displayObjectSize={true} 
                                sortKeys={true} 
                                collapsed={1} 
                                onEdit={this.onEdit}
                                displayDataTypes={false}
                                theme="codeschool" iconStyle="circle" 
                            />
                        </div>
                    </Grid>
                </Grid>

                <Modal open={this.state.showModal} onClose={this.closeModal}>
                    <div className="modal-center">
                        <ImportExportDeploySpec 
                            project={this.state.project}
                            specification={this.state.spec}
                            onImportSpec={this.onImportSpec}
                        />
                    </div>
                </Modal>
            </div>
        );
    }
}

export default DeploySpec;