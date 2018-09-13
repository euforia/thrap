import React, { Component } from 'react';
import { ExpandMore, Clear } from '@material-ui/icons';
import DatastoreSelector from './DatastoreSelector.js';
import './CreateBasicStack.css';

class CreateBasicStack extends Component {
    constructor(props) {
        super(props);
    
        this.onStackCreate = this.onStackCreate.bind(this);
        this.onInputChange = this.onInputChange.bind(this);

        this.state = {
            name: '',
            inputClass: '',
        }
    }

    onStackCreate() {
        if (this.state.name === '') {
            this.setState({inputClass: 'invalid-input'});
            return;
        }

        this.setState({inputClass: ''});
        this.props.onStackCreated(this.state.name);
    }
    

    onInputChange(event) {
        this.setState({
            name: event.target.value,
        });
    }
  
    render() {

        return (
            <div>
                <div style={{padding: "20px 0"}}>
                    <table className="table-header">
                        <tbody>
                            <tr>
                                <td className="header-label">Create Stack</td>
                                <td className="header-body">
                                    <button title="Cancel" className="btn-control" onClick={this.props.onCloseDialogue}><Clear /></button>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
                
                <div className="create-form-container">
                    <div className="section">
                        <div className="section-label">Name</div>
                        <div><input type="text" className={"create-name " + this.state.inputClass} placeholder="Stack name" value={this.state.name} onChange={this.onInputChange}/></div>
                    </div>
                    <div className="section">
                        <div className="section-label">Language</div>
                        <div className="create-proj-input select-container">
                            <select>
                                <option value="none">None</option>
                                <option value="shell">Shell</option>
                                <option value="go">Golang</option>
                                <option value="python">Python</option>
                                <option value="java">Java</option>
                                <option value="ruby">Ruby</option>
                                <option value="csharp">C#</option>
                            </select>
                            <span className="select-expand">
                                <ExpandMore/>
                            </span>
                        </div>
                    </div>
                    <div className="section">
                        <div className="section-label">Proxy</div>
                        <div className="create-proj-input select-container">
                            <select>
                                <option value="none">None</option>
                                <option value="nginx">Nginx</option>
                            </select>
                            <span className="select-expand">
                                <ExpandMore/>
                            </span>
                        </div>
                    </div>
                    <div className="section">
                        <DatastoreSelector />
                    </div>

                    <div className="create-btn-container">
                        <button type="submit" className="btn-default" onClick={this.onStackCreate}>Create Stack</button>
                    </div>
                </div>
            </div>
        );
    }
}
  
export default CreateBasicStack;
  