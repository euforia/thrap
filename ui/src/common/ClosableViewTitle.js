import React, { Component } from 'react';
// import { Grid } from "@material-ui/core";
import { Clear } from '@material-ui/icons';

const styles = ({
    title: {
        // fontSize: '24px',
        // fontWeight: 100,
        textAlign: 'left',
    },
    btn: {
        margin: "0 10px",
    },
    icon: {
        height: 36,
        width: 36,
    }
    // subtext: {
    //     fontSize: '12px',
    //     opacity: 0.5,
    //     fontWeight: 100,
    //     paddingTop: '5px',
    // },
});

class ClosableViewTitle extends Component {
    constructor(props) {
        super(props);

        this.state = {
            title: props.title === undefined ? '' : props.title,
        }
    }

    getSubText() {
        if (this.props.subtext === undefined) {
            return;
        } else if (this.props.subtext === '') {
            return;
        }

        return (
            <div className="subscript">{this.props.subtext}</div>
        )
    }

    render() {
        return (
            <table className="header-container">
                <tbody>
                    <tr>
                        <td style={{width: '90px'}}>
                            <button title="Close" style={styles.btn} className="btn-trans" onClick={this.props.onCloseDialogue}>
                                <Clear style={styles.icon}/>
                            </button>
                        </td>
                        <td>
                            <div style={styles.title} className="header-title">{this.state.title}</div> 
                            {this.getSubText()}
                        </td>
                    </tr>
                </tbody>
            </table>
        );
    }
}

export default ClosableViewTitle;