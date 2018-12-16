import React, { Component } from 'react';
import {withStyles, Paper, Typography, Chip, Grid, IconButton, Tooltip} from '@material-ui/core';
import { List, ListItem, ListItemText, ListItemSecondaryAction } from '@material-ui/core';
import {Divider} from '@material-ui/core';
import AddIcon from '@material-ui/icons/Add';
import { Link } from 'react-router-dom';
import {thrap} from '../api/thrap';


const styles = theme => ({
    paper: {
        paddingTop: theme.spacing.unit*3,
        paddingBottom: theme.spacing.unit*3,
        paddingLeft: theme.spacing.unit*3,
        paddingRight: theme.spacing.unit*3,
        marginLeft: theme.spacing.unit,
        marginRight:theme.spacing.unit,
        marginTop: theme.spacing.unit,
        marginBottom:theme.spacing.unit,
    },
    status: {
        textAlign:'right',
    },
    heading: {
        paddingTop: theme.spacing.unit,
        paddingBottom: theme.spacing.unit*3,
        display: 'inline-block',
    },
    deployInfo: {
        paddingBottom: theme.spacing.unit,
    },
    listItem:{
        paddingLeft: theme.spacing.unit,
    },
    divider:{
        // paddingLeft: theme.spacing.unit,
        marginTop: theme.spacing.unit*3,
        marginBottom: theme.spacing.unit,
    }
});

function deploysByProfile(deploys) {
    var d = thrap.translateDeploys(deploys);
    var out = {};
    for (var i=0;i<d.length;i++) {
        var prof = d[i].profile
        if (!out[prof]) out[prof]= [];
        out[prof].push(d[i]);
    }
    var arr = Object.keys(out).map((k) =>{
        return {name: k, list: out[k]};
    });
    return arr;
}

class DeploysBrief extends Component {

    render() {
        const {classes, project } = this.props;
        const deploys = deploysByProfile(this.props.deploys);

        return (
            <Grid container>
                {deploys.map((obj, i) => (
                <Grid item xs={4} key={i}>
                    <Paper className={classes.paper}>
                        <Grid container alignItems="center">
                            <Grid item xs={10}>
                                <Typography variant="h5"><b>{obj.name}</b></Typography>
                            </Grid>
                            <Grid item xs={2} style={{textAlign:'right'}}>
                                <Tooltip title={`New ${obj.name} deployment`} placement="right">
                                    <IconButton
                                        component={Link} 
                                        to={`/project/${project}/deploy/${obj.name}/new`}
                                    >
                                        <AddIcon/>
                                    </IconButton>
                                </Tooltip>
                            </Grid>
                        </Grid>
                        <Divider className={classes.divider}></Divider>
                        <List>
                            {obj.list.map((item,j) =>(
                            <ListItem button 
                                className={classes.listItem} 
                                key={j}
                                component={Link} 
                                to={'/project/'+project+'/deploy/'+obj.name+'/'+item.instance}
                            >
                                <ListItemText>{item.instance}</ListItemText>
                                <ListItemSecondaryAction>
                                    <Chip label={item.status} color={item.color}/>
                                </ListItemSecondaryAction>
                            </ListItem>
                            ))}
                        </List>
                    </Paper>
                </Grid>
                ))}
            </Grid>
        );
    }
}

export default withStyles(styles)(DeploysBrief);