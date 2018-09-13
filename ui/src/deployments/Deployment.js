import React, { Component } from 'react';
import ReactJson from 'react-json-view';

import { Grid } from '@material-ui/core';
import Modal from '@material-ui/core/Modal';

import DeployControls from './DeployControls.js';
import DeploymentDetails from './DeploymentDetails.js';
import DeployInstance from './DeployInstance.js';

import ClosableRightViewTitle from '../common/ClosableRightViewTitle.js';

import thrap from '../thrap.js';

const styles = ({
    container: {
        fontWeight: 100,
        fontSize: '13px',
        paddingTop: '40px',
    },
});

class Deployment extends Component {
    constructor(props) {
        super(props);

        this.closeDeployModal = this.closeDeployModal.bind(this);
        this.openDeployModal = this.openDeployModal.bind(this);
        
        this.onDeploy = this.onDeploy.bind(this);
        this.onDeployError = this.onDeployError.bind(this);

        this.onCloseSpec = this.onCloseSpec.bind(this);
        this.onShowSpec = this.onShowSpec.bind(this);

        this.state = {
            deployModal: false,
            project: props.project,
            environment: props.environment,
            deployment: props.deployment,
            spec: {},
            showSpec: false,
        }

        if ((props.deployment.Spec !== undefined) && (props.deployment.Spec !== null)) {
            var s = atob(props.deployment.Spec);
            this.state.spec = JSON.parse(s);
        }

    }

    onCloseSpec() {
        this.setState({
            showSpec: false,
        });
    }

    onShowSpec() {
        this.setState({
            showSpec: true,
        });
    }

    getDeployment() {
        thrap.deployment(this.state.project.ID, this.state.environment.ID, this.state.deployment.Name)
            .then(({data}) => {
                console.log("DATA",data);
                this.setState({
                    deployment: data,
                });
            })
            .catch(error => {
                var resp = error.response;
                console.log(resp);
            });
    }

    closeDeployModal() {
        this.setState({
            deployModal:false,
        });
    }

    openDeployModal() {
        this.setState({
            deployModal:true,
        });
    }

    // Called after a deploy has occured
    onDeploy() {
        // Close modal
        this.setState({
            deployModal:false,
        });

        this.refreshDeploy();
    }

    refreshDeploy() {
        // Update this one.
        this.getDeployment();
        // Upadte list
        this.props.onDeploy();
    }

    onDeployError(errStr) {
        this.refreshDeploy();
    }

    render() {
        var body;
        if (this.state.showSpec) {
            body = (
                <div>
                    <ClosableRightViewTitle title="Specification" onCloseDialogue={this.onCloseSpec} />
                    <ReactJson src={this.state.spec} name="spec" style={{background: 'none'}}
                        displayObjectSize={true} sortKeys={true} collapsed={1} displayDataTypes={false}
                        theme="codeschool" iconStyle="circle" 
                    />
                </div>
            );
        } else {
            body = (
                <DeploymentDetails
                    environment={this.state.environment} 
                    deployment={this.state.deployment} 
                />
            );
        }


        return (
            <div style={styles.container}>
                <Grid container spacing={0}>
                    <Grid item xs={2}>
                        <DeployControls 
                            onClose={this.props.onCloseDialogue} 
                            onDeploy={this.openDeployModal}
                            onSettings={this.onShowSpec}
                        />
                    </Grid>
                    <Grid item xs={10}>{body}</Grid>
                </Grid>

                <Modal open={this.state.deployModal} onClose={this.closeDeployModal}>
                    <div className="modal-center">
                        <DeployInstance 
                            project={this.state.project} 
                            environment={this.state.environment}
                            deployment={this.state.deployment}
                            onDeploy={this.onDeploy}
                            onDeployError={this.onDeployError}
                        />
                    </div>
                </Modal>
            </div>
        );
    }
}

export default Deployment;