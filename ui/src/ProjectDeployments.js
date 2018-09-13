import React, { Component } from 'react';
import { Add, SettingsOutlined } from '@material-ui/icons';
import Grid from '@material-ui/core/Grid';

import thrap from './thrap.js';

import Deployment from './deployments/Deployment.js';
import NewInstance from './deployments/NewInstance.js';
import DeploySpec from './spec/DeploySpec.js';
import { IconButton } from '@material-ui/core';

const styles = ({
    icon: {
        height: 36,
        width: 36,
    }
});

class ProjectDeployments extends Component {
    constructor(props) {
        super(props);

        this.onConfigureSpec = this.onConfigureSpec.bind(this);
        this.onCreateDeployable = this.onCreateDeployable.bind(this);
        this.onDeployableCreated = this.onDeployableCreated.bind(this);
        this.showDeployDetails = this.showDeployDetails.bind(this);
        
        this.onDeploy = this.onDeploy.bind(this);
        
        this.onCloseDialogue = this.onCloseDialogue.bind(this);
        

        this.state = {
            project:        props.project,
            environments:   props.environments,
            deployments:    {},
            // deploy:         false,
            newDeploy:      false,
            deployDetails:  false,
            configureSpec: false,

            selectedEnv:    -1,
            selectedDeploy: '',
        }

        this.getDeployments();
    }

    transformDeployment(data) {
        var dict = {};
                
        for (var i = 0; i < data.length; i++) {
            var deploy = data[i];
            if (dict[deploy.Profile.ID] === undefined) {
                dict[deploy.Profile.ID] = [];
            }
            dict[deploy.Profile.ID].push(deploy);
        }

        return dict;
    }

    getDeployments() {
        var transform = this.transformDeployment;

        thrap.deployments(this.state.project.ID)
            .then(({ data }) => {
                var d =  transform(data);
                this.setState({deployments: d});
            
            }).catch(error => {
                console.log(error);
            });
    }

    onDeploy() {
        // Refresh after creation
        this.getDeployments();
    }

    onCreateDeployable(event) {
        var envIdx = event.currentTarget.getAttribute('environment');
        this.setState({
            newDeploy: true,
            selectedEnv: envIdx,
        });
    }

    onDeployableCreated() {
        this.setState({
            newDeploy: false,
            deployDetails: false,
        });
        this.getDeployments();
    }

    onConfigureSpec() {
        this.setState({
            configureSpec: true,
        });
    }

    showDeployDetails(event) {
        var envIdx = event.currentTarget.getAttribute('environment'),
            deployIdx = event.currentTarget.getAttribute('deployment');
    
        this.setState({
            selectedEnv: envIdx,
            selectedDeploy: deployIdx,
            deployDetails: true,
        });
    }

    onCloseDialogue() {
        this.setState({
            newDeploy: false,
            deployDetails: false,
            configureSpec: false,
        });
    }

    getView() {
        var selectedEnv = this.state.selectedEnv;
        var env = this.state.environments[selectedEnv],
            project = this.state.project;

        if (this.state.newDeploy) {
            return (
                <NewInstance project={project} environment={env}
                    onCloseDialogue={this.onCloseDialogue} 
                    onDeployableCreated={this.onDeployableCreated} 
                />
            );
        
        } else if (this.state.deployDetails) {            
            
            var deployment = this.state.deployments[env.ID][this.state.selectedDeploy];
            
            return (
                <Deployment project={project} environment={env} deployment={deployment} 
                    onCloseDialogue={this.onCloseDialogue} 
                    onDeploy={this.onDeploy}
                />
            );

        }
    }

    render() {
        if (this.state.newDeploy || this.state.deployDetails) {
            return this.getView();
        } else if (this.state.configureSpec) {
            return (
                <DeploySpec project={this.state.project} onCloseDialogue={this.onCloseDialogue} />
            );
        }
        
        var items = [],
            envs = this.state.environments;

        for (var i = 0; i < envs.length; i++) {
            var denv = envs[i],
                deploys = this.state.deployments[denv.ID];
            
            if (deploys === undefined) {
                deploys = [];
            }

            items.push(
                <div key={denv.ID} className="panel">
                    <div className="panel-title">{denv.Name}</div>
                    <div className="list">
                    {deploys.map((obj, j) => 
                        <div key={obj.Name} className="list-item" environment={i} deployment={j} onClick={this.showDeployDetails}>
                            <div className="list-item-title">{obj.Name}</div>
                            <div className="list-item-desc">
                                <Grid container spacing={0}>
                                    <Grid item xs={6}>
                                        <div>Version: {obj.Version}</div>
                                    </Grid>
                                    <Grid item xs={6}>
                                        <div style={{textAlign: 'right'}}>{thrap.stateLabel(obj.State,obj.Status)}</div>
                                    </Grid>
                                </Grid>
                            </div>
                        </div>
                    )}
                    </div>
                    <div style={{padding: "40px 0px 30px 0", textAlign: "center"}}>
                        <button title="Deploy a new instance" className="btn-control" environment={i} onClick={this.onCreateDeployable}>
                            <Add />
                        </button>
                    </div>
                </div>
            );            
        }

        return (
            <div>
                <div className="header-container">
                    <IconButton aria-label="Deploy specification" onClick={this.onConfigureSpec}>
                        <SettingsOutlined style={styles.icon}/>
                    </IconButton>
                </div>
                <div style={{textAlign:"center"}}>
                    {items}
                </div>
            </div>
        );
    }
}
  
export default ProjectDeployments;
  