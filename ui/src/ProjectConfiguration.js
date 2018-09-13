import React, { Component } from 'react';
// import { Add } from '@material-ui/icons';
import StackComponent from './StackComponent.js';
import './StackComponent.css';

class ProjectConfiguration extends Component { 
    constructor(props) {
        super(props);

        this.state = {
            id: this.props.projectID,
            components: [{
                id: 'api',
                name: 'api',
                version: 'v0.2.0-12-abc123def',
                build: {
                    dockerfile: 'api.dockerfile',
                    context: '.',
                },
                ports: [{
                    name: "http",
                    value: 8080,
                }],
                env: {
                    file: '.env',
                    vars: [{
                        name: "VAR1",
                        value: "foobar",
                    }],
                },
                volumes: []
            },{
                id: 'ui',
                name: 'ui',
                version: 'v0.2.0-12-abc123def',
                build: {
                    dockerfile: 'ui.dockerfile',
                    context: '.',
                },
                ports: [{
                    name: "http",
                    value: 8080,
                }],
                env: {
                    file: '.env',
                    vars: [{
                        name: "API_URL",
                        value: 'comp.api.container.http.addr',
                    }],
                },
                dependencies: ["api"],
                volumes: []
            }],
            dependencies: [{
                id: 'dep1',
                name: 'dep1',
                scheme: 'http',
                external: false,
            },{
                id: 'github',
                name: 'github.com',
                scheme: 'https',
                external: true,
            },{
                id: 'rds',
                name: 'postgres.rds.aws.com:5432',
                scheme: 'tcp',
                external: true,
            }],
        };
    }

    render() {
        // In progress
        return (
            <div></div>
        );

        return (
            <div>
                <div>
                    <div>
                        <div className="config-heading">Components</div>
                    </div>
                    <div className="config-body">
                    {this.state.components.map((comp) => 
                        <StackComponent key={comp.id} component={comp} />
                    )}
                    </div>
                </div>

                <div>
                    <table>
                        <thead>
                            <tr>
                                <td colSpan="4" className="config-heading">Dependencies</td>
                                <td>

                                </td>
                            </tr>
                        </thead>
                        <tbody>
                        {this.state.dependencies.map((dep) => 
                            <tr key={dep.id}>
                                <td>{dep.id}</td>
                                <td>{dep.scheme}</td>
                                <td>{dep.name}</td>
                                <td><input type='checkbox' value={dep.external}/></td>
                                <td></td>
                            </tr>
                        )}
                        </tbody>
                    </table>
                </div>

            </div>
        );
    }
}

export default ProjectConfiguration;