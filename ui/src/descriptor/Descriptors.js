import React, { Component } from 'react';
import { withStyles, MenuItem, TextField, Paper, Divider } from '@material-ui/core';
import {thrap} from '../api/thrap';
import ImportSpec from './ImportSpec';
import ReactJson from 'react-json-view';
import { Modal, Typography } from '@material-ui/core';

const styles = theme => ({
    heading: {
        paddingTop: theme.spacing.unit,
        paddingBottom:theme.spacing.unit,
    },
    jsonViewer: {
        paddingTop: theme.spacing.unit*2,
        paddingBottom:theme.spacing.unit*2,
    },
    modalCenter: {
        width: 500,
        position: 'absolute',
        top: '50%',
        left: "50%",
        outline: "none",
        transform: 'translate(-50%, -50%)',
    }
})

class Descriptors extends Component {
    constructor(props) {
        super(props);
        this.state = {
            errMsg: '',
            specName: '',
            spec: {},
            specs:[],
            modalOpen: false,
        };

        this.fetchSpecs();
    }

    fetchSpecs() {
        var project = this.props.project;
        thrap.Specs(project)
        .then(resp => {
            var data = resp.data;
            var specs = [];
            for (var i=0;i<data.length;i++) {
                specs.push({name:data[i], id:data[i]});
            }
            this.setState({specs: specs});
        })
        .catch(err => {
            console.log(err.config);
            console.log(err.request);
            console.log(err.response);
        });
    }

    fetchSpec(specName) {
        var project = this.props.project;
        thrap.Spec(project, specName)
        .then(resp => {
            this.setState({spec: resp.data});
        })
        .catch(err => {
            console.log(err.config);
            console.log(err.request);
            console.log(err.response);
        });
    }

    handleModalClose = () => {
        this.setState({modalOpen:false});
    }

    handleSelect = event => {
        var val = event.target.value;
        if (val === '_add') {
            this.setState({modalOpen:true});
            return;
        }
        this.setState({specName:val});
        this.fetchSpec(val);
        this.props.onChange(val);
    }

    onImportDesc = (name, data) => {
        this.setState({
            specName: name,
            spec: data,
            modalOpen:false,
        });
        this.fetchSpecs();
    }

    onImportDescErr = () => {
        console.log('err');
    }
    
    render() {
        const { classes, project, disabled, error } = this.props;
        const { specName, specs, errMsg } = this.state;

        return (
            <div>
                {specs.length > 0
                ? <TextField label="Descriptor"
                    value={this.state.specName}
                    onChange={this.handleSelect}
                    select
                    error={error}
                    fullWidth
                    required
                    disabled={disabled}
                >
                    {specs.map(option => (
                    <MenuItem key={option.id} value={option.id} disabled={option.disabled}>
                        {option.name}
                    </MenuItem>
                    ))}
                    <Divider/>
                    <MenuItem value="_add">Add descriptor</MenuItem>
                </TextField>
                : <TextField label="Descriptor"
                    value={this.state.specName}
                    onChange={this.handleSelect}
                    select
                    error={error}
                    fullWidth
                    required
                    disabled={disabled}
                >
                    <MenuItem value="_add">Add descriptor</MenuItem>
                </TextField>
                }
                <Typography color="secondary">{errMsg}</Typography>
                <div className={classes.jsonViewer}>
                    <ReactJson  name={specName}
                        src={this.state.spec}
                        style={{background: 'none', width:'100%', minHeight: '50px'}}
                        displayObjectSize={true} 
                        displayDataTypes={false}
                        sortKeys={true} 
                        collapsed={3} 
                    />
                </div>
                <Modal open={this.state.modalOpen} onClose={this.handleModalClose}>
                    <Paper className={classes.modalCenter}>
                        <ImportSpec project={project} 
                            onImportSpec={this.onImportDesc} 
                            onError={this.onImportDescErr}/>
                    </Paper>
                </Modal>
            </div>
        );
    }
}

export default withStyles(styles)(Descriptors);