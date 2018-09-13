import React, { Component } from 'react';
import { Add, Search } from '@material-ui/icons';

import thrap from '../thrap.js';
import './Projects.css';

class ListProjects extends Component {
    constructor(props) {
        super(props);
    
        this.onFilterChange = this.onFilterChange.bind(this);

        this.state = {
            filter: '',
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

    render() {

        var items = [];
        var projs = this.state.projects;
        for (var i = 0; i < projs.length; i++) {
            var proj = projs[i];

            if (this.state.filter === "" || proj.Name.includes(this.state.filter)) {
                items.push(
                    <div key={proj.Name} className="list-item" project={proj.ID} onClick={this.props.onProjectDetails}>
                        <div className="list-item-title">{proj.Name}</div>
                        <div className="list-item-desc">{proj.Source}</div>
                    </div>
                );
            }
        }

        return (
            <div id="projects"> 
                <table className="header">
                    <tbody>
                        <tr>
                            <td className="header-title">Projects</td>
                            <td className="header-body">
                                <button title="Create project" className="btn-control" onClick={this.props.onCreateProject}><Add /></button>
                            </td>
                        </tr>
                    </tbody>
                </table>
                <div id="filter">  
                    <input type="text" placeholder="Search" value={this.state.filter} onChange={this.onFilterChange} />            
                    <span><Search style={{height: 36, width: 36}}/></span>
                </div>
                <div className="list" style={{margin: "20px 0"}}>{items}</div>
            </div>
        );
    }
}
  
export default ListProjects;
  