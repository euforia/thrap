import React, { Component } from 'react';
import { DeviceHub, Storage, Sync } from '@material-ui/icons';
import KeyValuePairs from './common/KeyValuePairs';


class StackComponent extends Component {
	constructor(props) {
		super(props);

		this.state = {
			comp: this.props.component,
		}
	}

	render() {
		return (
			<div>
				<table>
					<tbody>
						<tr>
							<td rowSpan="3">{this.state.comp.id}</td>
							<td>
								<div>Name</div>
								<div><input type="text" value={this.state.comp.name} /></div>
							</td>
							<td>
								<div>Version</div>
								<div><input type="text" value={this.state.comp.version} /></div>
							</td>
						</tr>
						<tr>
							<td>
								<div><Sync /><span> Runtime</span></div>
								<div>Environment</div>
								<div>File</div>
								<div><input type="text" value={this.state.comp.env.file}/></div>
								<KeyValuePairs title="Variables" pairs={this.state.comp.env.vars} />			
							</td>
							<td>
								<div><DeviceHub /><span> Network</span></div>
								<KeyValuePairs title="Port map" pairs={this.state.comp.ports}/>
								<div>Health Checks</div>
							</td>
						</tr>
						<tr>
							<td>
								<div>Command</div>
								<div>
									<input type="text" />
								</div>
								<div>Arguments</div>
								<div>
									<input type="text" />
								</div>
							</td>
							<td>
								<div><Storage /><span> State</span></div>
								<div>Volumes</div>
								<div>
								{this.state.comp.volumes.map((obj) => 
									<div>
										<span>{obj}</span>
									</div>
								)}
								</div>

								<div>Secrets</div>
							</td>
							</tr>
					</tbody>
				</table>
				
				

			</div>
		);
	}
}

export default StackComponent;