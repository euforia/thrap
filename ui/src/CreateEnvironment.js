import React, { Component } from 'react';
import { Clear } from '@material-ui/icons';
import './CreateEnvironment.css';

class CreateEnvironment extends Component {
    constructor(props) {
        super(props);
    
        this.onEnvCreate = this.onEnvCreate.bind(this);
        this.onInputChange = this.onInputChange.bind(this);

        this.state = {
            name: '',
            inputClass: '',
        }
    }

    onEnvCreate() {
        if (this.state.name === '') {
            this.setState({inputClass: 'invalid-input'});
            return;
        }

        this.setState({inputClass: ''});
        this.props.onEnvCreated(this.state.name);
    }
    

    onInputChange(event) {
        this.setState({
            name: event.target.value,
        });
    }
  
    render() {

        return (
            <div id="create-env">
                <table className="header">
                    <tbody>
                        <tr>
                            <td className="header-label">Create Environment</td>
                            <td className="header-body">
                                <button title="Cancel" className="btn-control" onClick={this.props.onCloseDialogue}><Clear /></button>
                            </td>
                        </tr>
                    </tbody>
                </table>
                
                <div id="create-env-result">
                </div>
                
                <div className="create-form-container">
                    <div className="section">
                        <div className="section-label">Name</div>
                        <div><input type="text" className={"create-name " + this.state.inputClass} placeholder="Environment name" value={this.state.name} onChange={this.onInputChange}/></div>
                    </div>
                    <div className="section">
                        
                    </div>

                    <div className="create-btn-container">
                        <button type="submit" className="btn-default" onClick={this.onEnvCreate}>Create Environment</button>
                    </div>
                </div>

            </div>
        );
    }
}
  
export default CreateEnvironment;
  