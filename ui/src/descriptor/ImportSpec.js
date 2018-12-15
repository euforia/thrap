import React, { Component } from 'react';
import { withStyles } from '@material-ui/core';
import { InsertDriveFileOutlined } from '@material-ui/icons';
import { Button, MenuItem, TextField, FormControl, InputAdornment } from '@material-ui/core';
import { Typography} from '@material-ui/core';
import { CloudUploadOutlined } from '@material-ui/icons';

import thrap from '../api/thrap';

const styles = theme => ({
    helperText: {
        // color: '#abbbc6'
    },
    fileUploadInput: {
        opacity: 0,
        position: 'absolute', 
        height: '60px', 
        width:'100%',
        cursor: 'pointer',
    },
    importIcon: {
        marginRight: theme.spacing.unit,
    },
    container: {
        paddingLeft: theme.spacing.unit*3,
        paddingRight: theme.spacing.unit*3,
        paddingTop: theme.spacing.unit*4,
        paddingBottom: theme.spacing.unit*4,
    }
});

const descContentTypes = ([
    { 
        label: 'tmpl+hcl', 
        value: 'application/vnd.thrap.mold.deployment.descriptor.v1+hcl' 
    },
    { 
        label: 'nomad+hcl', 
        value: 'application/vnd.thrap.nomad.deployment.descriptor.v1+hcl' 
    },
    { 
        label: 'nomad+json', 
        value: 'application/vnd.thrap.nomad.deployment.descriptor.v1+json'
     }
]);


class ImportSpec extends Component {
    constructor(props) {
        super(props);

        this.state = {
            selectedType: '',
            invalidContentType: false,
            selectedSpecFile: '',
            invalidSpecFile: false,
            specName: '',
            invalidSpecName: false,
        }

        this.fileInput = React.createRef();
    }

    onDescContentTypeChange = (event) => {
        var value = event.target.value;
        this.setState({
            selectedType: value,
            invalidContentType: (value === ''),
        });
    };

    onFileSelected = () => {
        var file = this.fileInput.current.files[0];
        this.setState({
            selectedSpecFile: file.name,
            invalidSpecFile: (file.name === ''),
        });
    }

    onImportSpec = () => {
        var mimeType = this.state.selectedType;
        
        var s = {};
        if (mimeType === '') {
            s.invalidContentType = true;
        } else {
            s.invalidContentType = false;
        }

        if (this.state.selectedSpecFile === '') {
            s.invalidSpecFile = true;
        } else {
            s.invalidSpecFile = false;
        }

        var {specName} = this.state;
        if (specName === '' || specName.includes(" ")) {
            s.invalidSpecName = true;
        } else {
            s.invalidSpecName = false;
        }

        this.setState(s);

        if (s.invalidContentType||s.invalidSpecFile||s.invalidSpecName) {
            return;
        }

        var f = this.fileInput.current.files[0];        
        var reader = new FileReader();
        var props = this.props;
        
        s = this.state;
        // Closure to capture the file information.
        reader.onload = (function(theFile) {
            var project = props.project;
            return function(e) {
                thrap.PutSpec(project, s.specName, mimeType, e.target.result)
                    .then(({data}) => {
                        props.onImportSpec(s.specName, data);
                    })
                    .catch(error => {
                        console.log(error);
                        props.onError(error);
                    });
            };
        })(f);
  
        // Read in the image file as a data URL.
        reader.readAsBinaryString(f);
    }

    handleNameChange = (event) => {
        var val = event.currentTarget.value;
        this.setState({
            specName:val,
            invalidSpecName: val==='' || val.includes(" ")
        });
    }

    render() {

        const { classes } = this.props;

        return (
            <div className={classes.container}>
                <div>
                    <Typography variant="h6">New Descriptor</Typography>
                    <TextField
                        value={this.state.specName}
                        onChange={this.handleNameChange}
                        fullWidth
                        margin="normal"
                        required
                        label="Name"
                        error={this.state.invalidSpecName}
                    />
                    <FormControl style={{width: '100%'}}>
                        
                        <TextField
                            id="input-with-icon-textfield"
                            label="Specification"
                            value={this.state.selectedSpecFile}
                            // placeholder="Choose File"
                            margin="normal"
                            InputProps={{
                                readOnly: true,
                                endAdornment: (
                                    <InputAdornment position="start">
                                        <InsertDriveFileOutlined></InsertDriveFileOutlined>
                                    </InputAdornment>
                                ),
                            }}
                            // FormHelperTextProps={{className:classes.helperText}}
                            helperText="deployment specification file"
                            error={this.state.invalidSpecFile}
                        />
                        <input type="file" className={classes.fileUploadInput}
                            ref={this.fileInput} onChange={this.onFileSelected} />
                        {/* <br /> */}
                    </FormControl>
                        <TextField 
                            id="desc-content-type"
                            name="Type"
                            label="Content Type"
                            select
                            // SelectProps={{
                            //     native: true,
                            // }}
                            value={this.state.selectedType}
                            onChange={this.onDescContentTypeChange}
                            margin="normal"
                            // FormHelperTextProps={{className:classes.helperText}}
                            helperText="descriptor content type"
                            required
                            fullWidth
                            error={this.state.invalidContentType}
                        >
                            {/* <option value=""></option> */}
                            {descContentTypes.map(option => (
                                <MenuItem key={option.value} value={option.value}>
                                    {option.label}
                                </MenuItem>
                            ))}
                        </TextField>
                    
                </div>
                <div style={{textAlign: 'center', marginTop:'30px'}}>
                    <Button color="primary" variant="outlined" onClick={this.onImportSpec}>
                        <CloudUploadOutlined fontSize="small" className={classes.importIcon}/>
                        Import
                    </Button>
                </div>
            </div>
        );
    }

}

export default withStyles(styles)(ImportSpec);