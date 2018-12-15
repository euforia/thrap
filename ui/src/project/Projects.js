import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import { withStyles } from '@material-ui/core/styles';
import { List, ListItem, ListItemText, TextField, Button, Grid } from '@material-ui/core';
import thrap from '../api/thrap';

const styles = theme => ({
  search: {
    paddingLeft: theme.spacing.unit * 2,
    paddingRight: theme.spacing.unit * 2,
  }
});

class Projects extends Component {

  constructor(props) {
    super(props);

    this.state = {
      projects: [],
      filter: '',
    }
    this.fetchProjects();
  }

  fetchProjects() {
    thrap.Projects().then(projs => {
      this.setState({projects: projs.data});
    });
  }

  handleFilterChange = event => {
    this.setState({
      filter: event.target.value
    });
  }
  
  filteredProjects() {
    var out = [],
        data = this.state.projects,
        query = this.state.filter;

    for (var i = 0; i < data.length; i++) {
        var d = data[i];
        if (d.Name.includes(query)) {
            out.push(d);
        }
    }
    return out 
  }

  render() {
    const projects = this.filteredProjects();
    const { classes } = this.props;

    return (
      <div>
        <Grid container alignItems="center" justify="space-between">
            <Grid item xs={10}>
              <div className={classes.search}>
                <TextField label="Search"
                  value={this.state.filter}
                  onChange={this.handleFilterChange}
                  margin="normal"
                  fullWidth
                >
                </TextField>
              </div>
            </Grid>
            <Grid item xs={1} style={{textAlign:'right'}}>
              <Button component={Link} to="/projects/new" 
                variant="outlined" color="primary"
              >
                New
              </Button>
            </Grid>
        </Grid>
        <List>
            {projects.map(option => (
            <ListItem button key={option.ID} component={Link} to={"/project/" + option.ID}>
              <ListItemText primary={option.Name} />
            </ListItem>
            ))}
        </List>
      </div>
    );
  }
}

export default withStyles(styles)(Projects);