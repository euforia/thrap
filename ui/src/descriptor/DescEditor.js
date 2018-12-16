import React, { Component } from 'react';
import {Controlled as CodeMirror} from 'react-codemirror2';
import { Typography, Button, Grid, withStyles } from '@material-ui/core';

import {thrap} from '../api/thrap';

require('codemirror/lib/codemirror.css');
require('codemirror/theme/material.css');
require('codemirror/mode/javascript/javascript.js');



const styles = theme => ({
    header: {
        paddingTop: theme.spacing.unit*2,
        paddingBottom: theme.spacing.unit*2,
    }
});

class DescEditor extends Component {
    state = {
        rawSpec: '',
    }

    componentDidMount() {
        this.fetchSpec();
    }

    fetchSpec() {
        var {project,descriptor} = this.props.match.params;
        thrap.Spec(project, descriptor)
        .then(resp => {
            // console.log(resp.data);
            this.setState({
                rawSpec: JSON.stringify(resp.data, null, 2)
            });
        })
        .catch(err => {
            console.log(err.config);
            console.log(err.request);
            console.log(err.response);
        });
    }

    render() {
        const {classes} = this.props;
        const {descriptor} = this.props.match.params;

        return (
            <div>
                <Grid container alignItems="center" className={classes.header}>
                    <Grid item xs={11}>
                        <Typography variant="h5">{descriptor}</Typography>
                    </Grid>
                    <Grid item xs={1} style={{textAlign:'right'}}>
                        <Button 
                            variant="outlined"
                            color="primary"
                            disabled
                        >
                        Save 
                        </Button>
                    </Grid>
                </Grid>
                <CodeMirror
                    value={this.state.rawSpec}
                    options={{
                        theme: 'material',
                        mode: { name: 'javascript', json: true },
                        smartIndent: true,
                        lineNumbers: true,
                        tabSize: 2,
                    }}
                    onBeforeChange={(editor, data, value) => {
                        this.setState({rawSpec: value});
                    }}
                    onChange={(editor, data, value) => {}}
                />
            </div>
        );
    }
}

export default withStyles(styles)(DescEditor);