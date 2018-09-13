import React, { Component } from 'react';
import { Tabs, Tab } from '@material-ui/core';

import thrap from './thrap.js';
import ProjectDeployments from './ProjectDeployments.js';
import Overview from './project/Overview.js';
import ClosableRightViewTitle from './common/ClosableRightViewTitle.js';
import './Project.css';

const viewMap = {
    0: 'Overview',
    // 1: 'Builds',
    1: 'Deployments'
}

class Project extends Component {
    constructor(props) {
        super(props);

        this.tabSelected = this.tabSelected.bind(this);

        this.state = {
            project:      {},
            environments: [],
            selectedView: 0,            
        }

    }

    componentDidMount() {
        var pid = this.props.match ? this.props.match.params.project : undefined;
        pid = pid ? pid : this.props.projectID;

        this.fetchProject(pid);
        this.fetchEnvironments()
    }

    fetchProject(id) {
        thrap.project(id)
            .then(({data}) => {
                this.setState({project: data});
            })
            .catch(error => {
                console.log(error);
            });
    }

    fetchEnvironments() {
        thrap.environments()
            .then(({data}) => {
                this.setState({environments: data});
            })
            .catch(error => {
                console.log(error);
            });
    }

    tabSelected = (event, value) => {
        this.setState({ selectedView: value });
    };

    getView() {
        var view = viewMap[this.state.selectedView];

        switch (view) {
            case 'Overview':
                return (
                    <Overview project={this.state.project} />
                );
            
            case 'Deployments':
                return (
                    <ProjectDeployments project={this.state.project} 
                        environments={this.state.environments} 
                    />
                );

            default:
                return;
        }
    }
    
    render() {
        if (this.state.project.ID === undefined) return(<div></div>);

        var view = this.getView();
        return (
            <div>
                <ClosableRightViewTitle title={this.state.project.ID} subtext={"Source: " + this.state.project.Source}
                    onCloseDialogue={this.props.onCloseDialogue} />
                <div>
                    <Tabs value={this.state.selectedView}
                        indicatorColor="primary"
                        onChange={this.tabSelected}
                        centered 
                        >
                        <Tab name="Overview" label="Overview" />
                        {/* <Tab name="Builds" label="Builds" /> */}
                        <Tab name="Deployments" label="Deployments" />
                    </Tabs>
                </div>
               {view}
            </div>
        );
    }
}

// export default withStyles(styles)(Project);
export default Project;