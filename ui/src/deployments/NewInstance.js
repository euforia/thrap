import React, { Component } from 'react';
import { Button, TextField, withStyles, FormControl } from '@material-ui/core';

import ClosableViewTitle from '../common/ClosableViewTitle.js';
import thrap from '../thrap.js';


const styles = theme => ({
    container: {
        padding: '0  10px',
    },
    helperText: {
        color: '#abbbc6'
    },
    formControl: {
        width: '100%',
    },
    buttonPanel: {
        margin: '40px 0 25px 0',
    }
});

class NewInstance extends Component {
    constructor(props) {
        super(props);

        this.createDeployable = this.createDeployable.bind(this);
        this.onInputChange = this.onInputChange.bind(this);

        this.state = {
            project: props.project,
            env: props.environment,
            name: '',                 // selected name
            errored: false,
            status: '',
            btnDisabled: false,
        }

    }

    onInputChange(event) {
        var value = event.target.value;
        this.setState({
            name: value,
            errored: (value === ''),
        });
    }

    createDeployable() {
        var state = this.state;
        if (state.name === '') {
            this.setState({errored:true});
            return 
        }

        this.setState({
            errored: false,
            btnDisabled: true,
        });
        thrap.createDeployment(state.project.ID, state.env.ID, state.name)
            .then(({ data }) => {
                this.setState({btnDisabled: false});
                this.props.onDeployableCreated(this.state.name);        
            })
            .catch(error => {
                var resp = error.response;
                this.setState({
                    status: resp.data,
                    btnDisabled: false,
                });
            });
    }

    render() {
        const { classes } = this.props;
        return(
            <div className={classes.container}>
                <ClosableViewTitle title={this.state.env.Name} onCloseDialogue={this.props.onCloseDialogue} />

                <div className="create-form-container">
                    <div className="header-container">
                        <div className="header-title" style={{textAlign: 'left'}}>New Instance</div>
                    </div>
                    <div className="error-container">{this.state.status}</div>
                    <FormControl className={classes.formControl}>
                    <TextField value={this.state.name} onChange={this.onInputChange} 
                        helperText="unique instance name id"
                        margin="normal"
                        FormHelperTextProps={{className:classes.helperText}}
                        required
                        fullWidth
                        label="Name"
                        name="Name"
                        type="text"
                        id="deployment-instance"
                        error={this.state.errored}
                        />
                    </FormControl>
                    <div className={classes.buttonPanel}>
                        <Button variant="contained" color="primary" 
                            onClick={this.createDeployable} disabled={this.state.btnDisabled}>
                            Create
                        </Button>
                    </div>
                </div>
            </div>
        );
    }
}

export default withStyles(styles)(NewInstance);