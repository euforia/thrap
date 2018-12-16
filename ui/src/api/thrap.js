import axios from 'axios';

// This is to support passing in the backend api host at runtime
const THRAP_HOST = process.env.REACT_APP_THRAP_HOST || '';
const VERSION = 'v1';
const THRAP_BASE = `${THRAP_HOST}/${VERSION}`;

// const UI_BASE_PATH = "/ui";

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

const nameRe = /^[a-z0-9\-_]+$/i;

const deployStateColors = {
  0: {
    0: 'default',
    1: 'default',
    2: 'default', 
  },
  1: {
      0: 'default',
      1: 'default',
      2: 'secondary',
  },
  2: {
      0: 'default',
      1: 'default',
      2: 'secondary',
  },
  3: {
      0: 'default',
      1: 'primary',
      2: 'secondary',
  },
};

class Thrap {
  constructor() {
    this.auth = {};

    var jsonData = sessionStorage.getItem(AUTH_STORE_KEY)
    if (jsonData !== null) {
      var data = JSON.parse(jsonData);
      this.auth = data;
    }
  }
  
  addr() {
    return `${THRAP_HOST}`;
  }

  isAuthd(profile) {
    // Any auth'd profile if profile not supplied
    if (profile === undefined || profile === '' || profile === null) {
      var arr = Object.keys(this.auth);
      if (arr.length > 0) {
        return true
      }
      return false
    }
    var a = this.auth[profile];
    if (a === undefined || a === null) return false
    return a.data !== undefined && a.data.id !== undefined;
  }

  stateLabel(state, status) {
    if (state === undefined || status === undefined) return 'Unknown';
    return deployState[state][status];
  }
  stateLabelColor(state, status) {
    if (state === undefined || status === undefined) return 'default';
    return deployStateColors[state][status];
  }

  translateDeploys(arr) {
    return arr.map((obj)=>(
        {
            instance: obj.Name,
            status: this.stateLabel(obj.State, obj.Status),
            profile: obj.Profile.ID,
            color: this.stateLabelColor(obj.State, obj.Status),
        }
    ));
}

  requestHeaders(profile) {
    // TODO:check profile exists
    var obj = {
      [PROFILE_HEADER]: profile,
      [TOKEN_HEADER]:  this.auth[profile].data.id,
    };
    return obj;
  }

  Profiles() {
    const path = `${THRAP_BASE}/profiles`;
    return axios.get(path);
  }

  Profile(prof) {
    const path = `${THRAP_BASE}/profile/${prof}`;
    return axios.get(path);
  }

  Projects() {
    const path = `${THRAP_BASE}/projects`;
    return axios.get(path);
  }

  CreateProject(payload) {
    const path = `${THRAP_BASE}/project/${payload.Project.ID}`;
    const profiles = Object.keys(this.auth);

    return axios({
      method: 'POST',
      url: path,
      data: payload,
      headers: this.requestHeaders(profiles[0]),
    });
  }

  Project(project) {
    const path = `${THRAP_BASE}/project/${project}`;
    return axios.get(path);
  }

  Deployments(project) {
    const path = `${THRAP_BASE}/project/${project}/deployments`;
    return axios.get(path);
  }

  CreateDeployment(project, profile, instance) {
    const path = `${THRAP_BASE}/project/${project}/deployment/${profile}/${instance}`;
    
    return axios({
      method: 'put',
      url: path,
      data: {Name: instance},
      headers: this.requestHeaders(profile),
    });
  }

  Deployment(project, profile, instance) {
    const path = `${THRAP_BASE}/project/${project}/deployment/${profile}/${instance}`;
    return axios.get(path);
  }

  Spec(project, version) {
    const path = `${THRAP_BASE}/project/${project}/deployment/spec/${version}`;
    return axios.get(path);
  }

  Specs(project) {
    const path = `${THRAP_BASE}/project/${project}/deployment/specs`;
    return axios.get(path);
  }

  PutSpec(project, specName, specFormat, spec) {
    const path = `${THRAP_BASE}/project/${project}/deployment/spec/${specName}`;
    const profiles = Object.keys(this.auth);

    var headers = this.requestHeaders(profiles[0]);
    headers['Content-Type'] = specFormat;
  
    return axios({
      method: 'put',
      url: path,
      data: spec,
      headers: headers,
    });
  }

  DeployInstance(project, profile, instance, payload) {
    const path = `${THRAP_BASE}/project/${project}/deployment/${profile}/${instance}`;
    return axios({
      method: 'post',
      url: path,
      data: payload,
      headers: this.requestHeaders(profile),
    });
  }

  AuthMethods() {
    return new Promise((resolve, reject) => {
      resolve([{
        id: 'vault',
        name: 'Vault Token',
        type: 'token',
      }]);
      // resolve([{
      //     id: 'github',
      //     name: 'Github Token',
      //     type: 'token',
      // },{
      //     id: 'userpass',
      //     name: 'Username & Password',
      //     type: 'userpass',
      // },{
      //     id: 'vault',
      //     name: 'Vault Token',
      //     type: 'token',
      // }]);
    });
  }
  
  Deauthenticate() {
    this.auth = {};
    sessionStorage.removeItem(AUTH_STORE_KEY);
  }

  Authenticate(profile, token) {
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
          this.auth[profile] = data;
          sessionStorage.setItem(AUTH_STORE_KEY, JSON.stringify(this.auth));
          resolve(data);
        })
        .catch(error => {
          reject(error);
        });
    });
  }

};

const thrap = new Thrap();

function validateName(val) {
  var match = nameRe.exec(val)
  return match ? '' 
      : 'Only alpha-numeric characters, -, and _ allowed'
}

export {
  thrap,
  validateName,
};

