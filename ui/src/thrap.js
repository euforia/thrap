import axios from 'axios';

// This is to support passing in the backend api host at runtime
const THRAP_HOST = process.env.REACT_APP_THRAP_HOST || '';
const VERSION = 'v1';
const THRAP_BASE = `${THRAP_HOST}/${VERSION}`;

const UI_BASE_PATH = "/ui";

const AUTH_STORE_KEY ='iam';
const TOKEN_HEADER = 'X-Vault-Token';
const PROFILE_HEADER = 'Thrap-Profile';

const deployState = {
  0: {
      0: 'Unknown',
      1: 'Unknown',
      2: 'Unknown', 
  },
  1: {
      0: 'Creating',
      1: 'Created',
      2: 'Created failed',
  },
  2: {
      0: 'Preparing',
      1: 'Prepared',
      2: 'Prepare failed',
  },
  3: {
      0: 'Deploying',
      1: 'Deployed',
      2: 'Deploy failed',
  },
};

class Thrap {
  constructor() {
    this.auth = {};

    var jsonData = sessionStorage.getItem('iam')
    if (jsonData !== null) {
      var data = JSON.parse(jsonData);
      this.auth = data;
    }
  }
  
  addr() {
    return `${THRAP_HOST}`;
  }

  setURLPath(path) {
    window.history.pushState("", "", UI_BASE_PATH+path);
  }

  isAuthd() {
    return this.auth.data !== undefined;
  }

  stateLabel(state, status) {
    return deployState[state][status]
  }

  requestHeaders() {
    return {
      [PROFILE_HEADER]: this.auth.profile,
      [TOKEN_HEADER]:  this.auth.data.id,
    };
  }

  environments() {
    const path = `${THRAP_BASE}/profiles`;
    return axios.get(path);
  }

  projects() {
    const path = `${THRAP_BASE}/projects`;
    
    return axios.get(path);
  }

  createProject(payload) {
    const path = `${THRAP_BASE}/project/${payload.Project.ID}`;

    return axios({
      method: 'POST',
      url: path,
      data: payload,
      headers: this.requestHeaders(),
    });
  }

  project(project) {
    const path = `${THRAP_BASE}/project/${project}`;
    
    return axios.get(path);
  }

  deployments(project) {
    const path = `${THRAP_BASE}/project/${project}/deployments`;
    
    return axios.get(path);
  }

  createDeployment(project, environment, instance) {
    const path = `${THRAP_BASE}/project/${project}/deployment/${environment}/${instance}`;
    
    return axios({
      method: 'POST',
      url: path,
      data: {Name: instance},
      headers: this.requestHeaders(),
    });
  }

  deployment(project, environment, instance) {
    const path = `${THRAP_BASE}/project/${project}/deployment/${environment}/${instance}`;
    
    return axios.get(path);
  }

  deploymentSpec(project) {
    const path = `${THRAP_BASE}/project/${project}/deployment/spec`;
    
    return axios.get(path);
  }

  importSpec(project, specFormat, spec) {
    const path = `${THRAP_BASE}/project/${project}/deployment/spec`;

    var headers = this.requestHeaders()
    headers['Content-Type'] = specFormat;

    return axios({
      method: 'post',
      url: path,
      data: spec,
      headers: headers,
    });
  }

  deployInstance(project, environment, instance, payload) {
    const path = `${THRAP_BASE}/project/${project}/deployment/${environment}/${instance}`;
    console.log(this.requestHeaders());
    return axios({
      method: 'put',
      url: path,
      data: payload,
      headers: this.requestHeaders(),
    });
  }

  deauthenticate() {
    this.auth = {};
    this.authProfile = '';
    sessionStorage.removeItem(AUTH_STORE_KEY);
  }

  authenticate(profile, token) {
    const path = `${THRAP_BASE}/login`

    var req = axios({
      method: 'post',
      url: path,
      headers: {
        [TOKEN_HEADER]: token,
        [PROFILE_HEADER]: profile,
      },
    });

    return new Promise((resolve, reject) => {
      req
        .then(({data}) => {
          data.profile = profile;
          this.auth = data;
          
          sessionStorage.setItem(AUTH_STORE_KEY, JSON.stringify(data));
          resolve(data);
        })
        .catch(error => {
          reject(error);
        });
    });
  }

};

const thrap = new Thrap();

export default thrap;

