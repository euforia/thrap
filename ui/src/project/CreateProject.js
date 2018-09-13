import React, { Component } from 'react';
import { Button, FormControl, TextField, withStyles } from '@material-ui/core'; 

import thrap from '../thrap.js';

import ClosableViewTitle from '../common/ClosableViewTitle.js';

const styles = theme => ({
    helperText: {
        color: '#abbbc6'
    }
});

class CreateProject extends Component {
    constructor(props) {
        super(props);
    
        this.onProjectCreate = this.onProjectCreate.bind(this);
        this.onInputChange = this.onInputChange.bind(this);

        this.state = {
            data: {
                Name: '',
                Description: '',
                Source:     '',
                Owner:      '',
                Maintainer: '',
                Developers: [],
            },
            errored: false,
            errorMsg: '',
            btnDisabled: false,
        }
    }

    onProjectCreate() {
        var project = this.state.data;
 
        if (project.Name === '') {
            this.setState({errored:true});
            return;
        }
        if (project.Maintainer === '') {
            this.setState({errored:true});
            return;
        }
        if (project.Source === '') {
            this.setState({errored:true});
            return;
        }

        // disabled button
        this.setState({
            errored: false,
            btnDisabled: true
        });

        var payload = {Project: project};
        payload.Project['ID'] = project.Name.toLowerCase();
        
        thrap.createProject(payload)
            .then(({ data }) => {
                this.setState({btnDisabled: false});
                this.props.onProjectCreated(data);
            })
            .catch(error => {
                console.log(error);
               
                var resp = error.response;
                console.log("Create project error", resp);
                this.setState({
                    errorMsg: error.message,
                    btnDisabled: false
                });
            });
    }

    onInputChange(event) {
        var name = event.target.name,
            val = event.target.value;

        var proj = this.state.data;
        proj[name] = val;
        this.setState({data: proj});    
    }
  
    render() {
        const { classes } = this.props;
        var proj = this.state.data;

        return (
            <div>
                <ClosableViewTitle onCloseDialogue={this.props.onCloseDialogue}/>
                <div className="create-form-container">
                    <div className="header-container">
                        <div className="header-title" style={{textAlign:'left'}}>Create Project</div>
                    </div>
                    <div className="error-container">{this.state.errorMsg}</div>
                    <FormControl style={{width: '100%'}}>
                        <TextField
                            id="project-name"
                            name="Name"
                            label="Name"
                            // InputLabelProps={{className:classes.inputLabel}}
                            FormHelperTextProps={{className:classes.helperText}}
                            helperText="unique project name used as the id"
                            margin="normal"
                            value={proj.Name}
                            onChange={this.onInputChange}
                            required
                            fullWidth
                            error={proj.Name === '' && this.state.errored}
                        />
                        <TextField
                            id="project-source-code"
                            name="Source"
                            label="Source Code"
                            value={proj.Source}
                            onChange={this.onInputChange}
                            // InputLabelProps={{className:classes.inputLabel}}
                            FormHelperTextProps={{className:classes.helperText}}
                            helperText="source code url"
                            margin="normal"
                            required
                            fullWidth
                            error={proj.Source === '' && this.state.errored}
                        />
                        <TextField
                            id="project-owner"
                            name="Owner"
                            label="Owner"
                            value={proj.Owner}
                            // InputLabelProps={{className:classes.inputLabel}}
                            FormHelperTextProps={{className:classes.helperText}}
                            helperText="owner of the project"
                            margin="normal"
                            onChange={this.onInputChange}
                            fullWidth
                        />
                        <TextField
                            id="project-maintainer"
                            name="Maintainer"
                            label="Maintainer"
                            value={proj.Maintainer}
                            onChange={this.onInputChange}
                            // InputLabelProps={{className:classes.inputLabel}}
                            FormHelperTextProps={{className:classes.helperText}}
                            helperText="maintainer of the project"
                            margin="normal"
                            required
                            fullWidth
                            error={proj.Maintainer === '' && this.state.errored}
                        />
                    </FormControl>
                
                    <div className="create-btn-container">
                        <Button variant="contained" color="primary" onClick={this.onProjectCreate} disabled={this.state.btnDisabled}>
                            Create
                        </Button>
                    </div>
                </div>
            </div>
        );
    }
}
  
export default withStyles(styles)(CreateProject);
  