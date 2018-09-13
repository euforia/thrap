import React, { Component } from 'react';
// import { Add, Search } from '@material-ui/icons';
import CreateProject from './project/CreateProject.js';
import ListProjects from './project/ListProjects';

import thrap from './thrap.js';
import Project from './Project.js';
import './project/Projects.css';
// import Route from 'react-router/Route';
// import { Switch } from 'react-router-dom';
// import DefaultRoute from 'react-router-dom/Route';

class Projects extends Component {
    constructor(props) {
        super(props);
    
        this.onFilterChange = this.onFilterChange.bind(this);
        
        this.onCreateProject = this.onCreateProject.bind(this);
        this.onProjectCreated = this.onProjectCreated.bind(this);
        
        this.onProjectDetails = this.onProjectDetails.bind(this);
        
        this.onCloseDialogue = this.onCloseDialogue.bind(this);


        this.state = {
            filter: '',
            createProject: false, // view 
            showProject: props.match ? props.match.params.project : '',
            projects: [],
        }

        this.getProjects();
    }

    getProjects() {
        thrap.projects()
            .then(({ data }) => {
                this.setState({ projects: data });
            }).catch(error => {
                console.log(error);
            });
    }

    onFilterChange(event) {
        this.setState({
            filter: event.target.value,
        });
    }

    onCreateProject() {
        this.setState({
            createProject: true,
        });
    }

    onProjectCreated(project) {
        this.setState({
            createProject: false,
        });

        // Re-fetch project list
        this.getProjects();
    }

    onProjectDetails(event) {
        // Get selected project
        var project = event.currentTarget.getAttribute('project');
        // Update URL
        this.props.history.push('/project/'+project);
        // Update state
        this.setState({showProject: project});
    }
    
    // Close all dialogs
    onCloseDialogue() {
        this.props.history.push('/projects')
        this.setState({
            createProject: false,
            showProject:'',
        });
    }

    render() {

        if (this.state.createProject) {
            return (
                <div style={{padding: '20px 0'}}>
                    <CreateProject onProjectCreated={this.onProjectCreated} onCloseDialogue={this.onCloseDialogue}/>
                </div>
            );
        } else if (this.state.showProject) {
            return (
                <Project 
                    projectID={this.state.showProject} 
                    onCloseDialogue={this.onCloseDialogue}
                />
            );
        }

        return (
            <ListProjects onProjectDetails={this.onProjectDetails} onCreateProject={this.onCreateProject}/>
        );
    }
}
  
export default Projects;
  