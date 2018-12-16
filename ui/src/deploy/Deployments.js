import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import { Grid, Typography, withStyles,  IconButton, Tooltip } from '@material-ui/core';

import ListIcon from '@material-ui/icons/List';
import AddToPhotosOutlinedIcon from '@material-ui/icons/AddToPhotosOutlined';
import ViewModuleIcon from '@material-ui/icons/ViewModule';

import {thrap} from '../api/thrap';
import DeploysBrief from './DeploysBrief';
import DeploysList from './DeploysList';

const styles = theme => ({
   header: {
       paddingTop: theme.spacing.unit*1.5,
       paddingBottom: theme.spacing.unit,
   },
   alignRight: {
       textAlign: 'right',
   },
   link: {
    color: 'inherit',
    '&:hover': {
        color: theme.palette.primary.main,
    }
   }
});

class Deployments extends Component {
    constructor(props) {
        super(props);

        this.state = {
            deploys:[],
            listView: false,
        }

        this.fetchDeploys()
    }

    fetchDeploys() {
        const proj = this.props.match.params.project;
        thrap.Deployments(proj).then(resp =>  {
            this.setState({deploys: resp.data});
        });
    }

    toggleView = () => {
        var l = !this.state.listView;
        this.setState({listView:l});
    }

    render() {
        const { deploys, listView} = this.state;
        const { project } = this.props;
        const { classes } = this.props;

        return (
            <div>
                <div className={classes.header}>
                    <Grid container alignItems="center" justify="space-between">
                        <Grid item xs={4}>
                            <Typography variant="h5">Deployments</Typography>
                        </Grid>
                        
                        <Grid item xs={4} style={{textAlign:'right'}}>
                            <IconButton onClick={this.toggleView}>
                                {listView 
                                ? <Tooltip title="Profile view"><ViewModuleIcon/></Tooltip>
                                : <Tooltip title="List view"><ListIcon/></Tooltip>
                                }
                            </IconButton>
                            <Tooltip title="New deployment">
                                <IconButton
                                    component={Link} 
                                    to={'/project/'+project+'/deploys/new'}
                                >
                                    <AddToPhotosOutlinedIcon/>
                                </IconButton>
                            </Tooltip>
                        </Grid>
                    </Grid>
                </div>
                {listView 
                ? <DeploysList deploys={deploys} project={project}/>
                : <DeploysBrief deploys={deploys} project={project}/>}
            </div>
        );
    }
}

export default withStyles(styles)(Deployments);