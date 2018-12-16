import React, { Component } from 'react';
import { Typography, Grid, withStyles, IconButton } from '@material-ui/core';
import { List, ListItem, ListItemText, ListItemSecondaryAction } from '@material-ui/core';
import AddIcon from '@material-ui/icons/Add';
import DeleteOutlinedIcon from '@material-ui/icons/DeleteOutlined';
import {Link} from 'react-router-dom';

import ConfirmDelete from '../common/ConfirmDelete';

import { thrap } from '../api/thrap';

const styles = theme => ({
    header: {
        paddingTop: theme.spacing.unit*2,
        paddingBottom: theme.spacing.unit*2,
        paddingLeft: theme.spacing.unit,
        paddingRight: theme.spacing.unit,
    },
    list: {
        width: '100%'
    },
    modalCenter: {
        width: 500,
        position: 'absolute',
        top: '50%',
        left: "50%",
        outline: "none",
        transform: 'translate(-50%, -50%)',
    }
});

class DescList extends Component {
    constructor(props) {
        super(props);
        this.state = {
            descs: [],
            modalOpen: false,
            delItem: -1,
            toDel: '',
        }
    }

    componentDidMount() {
        this.fetchDescs()
    }

    fetchDescs() {
        const {project} = this.props.match.params;
        thrap.Specs(project)
        .then(resp => {
            this.setState({descs: resp.data});
        })
        .catch(err => {
            console.log(err);
        })
    }

    handleModalClose = () => {
        this.setState({modalOpen: false});
    }

    handleModalOpen = (i) => {
        const {project} = this.props.match.params;
        if (!thrap.isAuthd()) {
            var path = `/login#/project/${project}/deploy/descriptors`;
            this.props.history.push(path);
            return;
        }

        var descs = this.state.descs;
        this.setState({
            delItem: i,
            modalOpen: true,
            toDel: descs[i]
        });
    }

    handleDeleteSpec = () => {
        this.setState({modalOpen:false});

        const {project} = this.props.match.params;
        const i = this.state.delItem;
        var {descs} = this.state;
        
        if (i<0) {
            console.error("Delete index not set");
            return;
        }

        thrap.DeleteSpec(project, descs[i])
        .then(resp => {
            descs.splice(i,1);
            this.setState({descs:descs});
        })
        .catch(err => {
            console.log(err);
        });
    }

    render() {
        const {classes} = this.props;
        const {project} = this.props.match.params;
        const {descs, toDel, modalOpen} = this.state;

        return (
            <div>
                <Grid container alignItems="center" className={classes.header}>
                    <Grid item xs={11}>
                        <Typography variant="h5">Descriptors</Typography>
                    </Grid>
                    <Grid item xs={1} style={{textAlign:'right'}}>
                        <IconButton disabled>
                            <AddIcon/>
                        </IconButton>
                    </Grid>
                </Grid>
                <List className={classes.list}>
                    {descs.map((desc, i) => (
                    <ListItem key={i} component={Link} button
                        to={`/project/${project}/deploy/descriptor/${desc}`}
                    >
                        <ListItemText>{desc}</ListItemText>
                        <ListItemSecondaryAction>
                            <IconButton onClick={event => this.handleModalOpen(i, event)}>
                                <DeleteOutlinedIcon/>
                            </IconButton>
                        </ListItemSecondaryAction>
                    </ListItem>
                    ))}
                </List>
                <ConfirmDelete
                    entity={toDel}
                    open={modalOpen}
                    onCancel={this.handleModalClose}
                    onDelete={this.handleDeleteSpec}
                />
            </div>
        );
    }
}

export default withStyles(styles)(DescList);