import React, { Component } from 'react';
import { withStyles } from '@material-ui/core';
import { InsertDriveFileOutlined } from '@material-ui/icons';
import { Button, TextField, FormControl, InputAdornment } from '@material-ui/core';
import { CloudUploadOutlined } from '@material-ui/icons';

import thrap from '../thrap';

const styles = theme => ({
    helperText: {
        color: '#abbbc6'
    },
    fileUploadInput: {
        opacity: 0,
        position: 'absolute', 
        height: '60px', 
        width:'100%',
        cursor: 'pointer',
    },
    importIcon: {
        fontSize: 20,
        marginRight: theme.spacing.unit,
    },

});

const descContentTypes = ([
    { 
        label: 'mold+hcl', 
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
        }

        this.onImportSpec = this.onImportSpec.bind(this);
        this.onDescContentTypeChange = this.onDescContentTypeChange.bind(this);
        this.onFileSelected = this.onFileSelected.bind(this);

        this.fileInput = React.createRef();
    }

    onDescContentTypeChange(event) {
        var value = event.target.value;
        this.setState({
            selectedType: value,
            invalidContentType: (value === ''),
        });
    };

    onFileSelected() {
        var file = this.fileInput.current.files[0];
        this.setState({
            selectedSpecFile: file.name,
            invalidSpecFile: (file.name === ''),
        });
    }

    onImportSpec() {
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

        this.setState(s);

        if (s.invalidContentType || s.invalidSpecFile) {
            return;
        }

        var f = this.fileInput.current.files[0];        
        var reader = new FileReader();
        var props = this.props;

        // Closure to capture the file information.
        reader.onload = (function(theFile) {
            return function(e) {
                thrap.importSpec(props.project.ID, mimeType, e.target.result)
                    .then(({data}) => {
                        props.onImportSpec(data);
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

    render() {

        const { classes } = this.props;

        return (
            <div>
                <div>
                    <FormControl style={{width: '100%'}}>
                      <TextField
                            id="input-with-icon-textfield"
                            label="Specification"
                            value={this.state.selectedSpecFile}
                            placeholder="Choose File"
                            InputProps={{
                                readOnly: true,
                                startAdornment: (
                                    <InputAdornment position="start">
                                        <InsertDriveFileOutlined></InsertDriveFileOutlined>
                                    </InputAdornment>
                                ),
                            }}
                            FormHelperTextProps={{className:classes.helperText}}
                            helperText="deployment specification file"
                            error={this.state.invalidSpecFile}
                        />
                        <input type="file" className={classes.fileUploadInput}
                            ref={this.fileInput} onChange={this.onFileSelected} />
                        <br />
                        <TextField 
                            id="desc-content-type"
                            name="Type"
                            label="Content Type"
                            select
                            SelectProps={{
                                native: true,
                            }}
                            value={this.state.selectedType}
                            onChange={this.onDescContentTypeChange}
                            margin="normal"
                            FormHelperTextProps={{className:classes.helperText}}
                            helperText="descriptor content type"
                            required
                            fullWidth
                            error={this.state.invalidContentType}
                        >
                            <option value=""></option>
                            {descContentTypes.map(option => (
                                <option key={option.value} value={option.value}>
                                    {option.label}
                                </option>
                            ))}
                        </TextField>
                    </FormControl>
                </div>
                <div style={{textAlign: 'center', marginTop:'30px'}}>
                    <Button color="primary" variant="contained" onClick={this.onImportSpec}>
                        <CloudUploadOutlined className={classes.importIcon}/>
                        Import
                    </Button>
                </div>
            </div>
        );
    }

}

export default withStyles(styles)(ImportSpec);