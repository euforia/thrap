import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import { Chip, withStyles, TableSortLabel } from '@material-ui/core';
import { Table, TableHead, TableRow, TableCell, TableBody } from '@material-ui/core';

import thrap from '../api/thrap';

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

function desc(a, b, orderBy) {
    if (b[orderBy] < a[orderBy]) {
      return -1;
    }
    if (b[orderBy] > a[orderBy]) {
      return 1;
    }
    return 0;
  }
  
function stableSort(array, cmp) {
    const stabilizedThis = array.map((el, index) => [el, index]);
    stabilizedThis.sort((a, b) => {
        const order = cmp(a[0], b[0]);
        if (order !== 0) return order;
        return a[1] - b[1];
    });
    return stabilizedThis.map(el => el[0]);
}

function getSorting(order, orderBy) {
    return order === 'desc' ? (a, b) => desc(a, b, orderBy) : (a, b) => -desc(a, b, orderBy);
}

class DeploysList extends Component {
    constructor(props) {
        super(props);

        this.state = {
            order: 'desc',
            orderBy: 'Instance',
            deploys:[],
        }

        // this.fetchDeploys()
    }

    // fetchDeploys() {
    //     const proj = this.props.match.params.project;
    //     thrap.Deployments(proj).then(resp =>  {
    //         console.log(resp.data);
    //         this.setState({deploys: resp.data});
    //     });
    // }

    getColor(s) {
        if (s === "dead") {
            return "secondary";
        }
        return "default";
    }

    createSortHandler = property => event => {
        var order = this.state.order === 'desc' ? 'asc' : 'desc';
        this.setState({
            orderBy: property,
            order: order,
        });
    };

    handleRowClick = (event, obj) => {
        // console.log(obj);
        var path = '/project/'+this.props.project+'/deploy/'+obj.profile+'/'+obj.instance;
        this.props.history.push(path);
    }

    render() {
        const {order, orderBy} = this.state;
        const { project, deploys } = this.props;
        const { classes } = this.props;

        return (
            <Table>
                <TableHead>
                    <TableRow>
                        <TableCell sortDirection={orderBy === 'Instance' ? order : false}>
                            <TableSortLabel active={orderBy === 'Instance'}
                                direction={order}
                                onClick={this.createSortHandler('Instance')}
                            >
                                Instance
                            </TableSortLabel>
                        </TableCell>
                        <TableCell sortDirection={orderBy === 'Profile' ? order : false}>
                            <TableSortLabel active={orderBy === 'Profile'}
                                direction={order}
                                onClick={this.createSortHandler('Profile')}
                            >
                                Profile
                            </TableSortLabel>
                        </TableCell>
                        <TableCell style={{textAlign:'center'}} sortDirection={orderBy === 'Status' ? order : false}>
                            <TableSortLabel active={orderBy === 'Status'}
                                direction={order}
                                onClick={this.createSortHandler('Status')}
                            >
                                Status
                            </TableSortLabel>
                        </TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {stableSort(deploys, getSorting(order, orderBy)).map((obj, i) => (                         
                    <TableRow key={i} hover={true}>
                        <TableCell>
                            <Link to={'/project/'+project+'/deploy/'+obj.Profile.ID+'/'+obj.Name}
                                className={classes.link}
                            >
                                <div>{obj.Name}</div>
                            </Link>
                        </TableCell>
                        <TableCell>
                            {/* temp */}
                            <Link to={'/project/'+project+'/deploy/'+obj.Profile.ID+'/'+obj.Name}
                                className={classes.link}
                            >
                                <div>{obj.Profile.ID}</div>
                            </Link>
                        </TableCell>
                        <TableCell style={{textAlign:'center'}}>
                            <Chip 
                                label={thrap.stateLabel(obj.State, obj.Status)} 
                                color={thrap.stateLabelColor(obj.State, obj.Status)}
                            />
                        </TableCell>
                    </TableRow>
                    ))}
                </TableBody>
            </Table>
        );
    }
}

export default withStyles(styles)(DeploysList);