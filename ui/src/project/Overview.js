import React, { Component } from 'react';

import { FormControl, TextField, withStyles } from '@material-ui/core';

const styles = theme => ({
    helperText: {
        color: '#abbbc6'
    }
});


class Overview extends Component {
    constructor(props) {
        super(props);

        this.state = {
            project: props.project,
        }

        this.onInputChange = this.onInputChange.bind(this);
    }

    onInputChange(event) {
        var name = event.target.name,
            val = event.target.value;

        var proj = this.state.project;
        proj[name] = val;
        this.setState({project:proj});    
    };
    
    render() {
        var proj = this.state.project;
        const { classes } = this.props;
        return (
            <div style={{textAlign:'center', padding: '40px'}}>
                <FormControl style={{width: '350px'}}>
                    <TextField
                        id="project-source-code"
                        label="Source Code"
                        // InputLabelProps={{className:classes.inputLabel}}
                        FormHelperTextProps={{className:classes.helperText}}
                        InputProps={{readOnly: true}}
                        helperText="source code url"
                        margin="normal"
                        value={proj.Source}
                        onChange={this.onInputChange}
                        name="Source"
                        required
                        fullWidth={true}
                        error={proj.Source === ''}
                    />
                    <TextField
                        id="project-owner"
                        label="Owner"
                        InputProps={{readOnly: true}}
                        // InputLabelProps={{className:classes.inputLabel}}
                        FormHelperTextProps={{className:classes.helperText}}
                        helperText="owner of the project"
                        margin="normal"
                        value={proj.Owner}
                        onChange={this.onInputChange}
                        name="Owner"
                        fullWidth={true}
                    />
                    <TextField
                        id="project-maintainer"
                        label="Maintainer"
                        InputProps={{readOnly: true}}
                        // InputLabelProps={{className:classes.inputLabel}}
                        FormHelperTextProps={{className:classes.helperText}}
                        helperText="maintainer of the project"
                        margin="normal"
                        value={proj.Maintainer}
                        onChange={this.onInputChange}
                        name="Maintainer"
                        required
                        fullWidth={true}
                        error={proj.Maintainer === ''}
                    />
                </FormControl>
            </div>
        );
    }
}

export default withStyles(styles)(Overview);
// export default Overview;