/*
  This component is shown when there is text in the search box. Currently, search is unimplemented.
  Neither the ConnectorDB server nor the frontend have search support right now.
  TODO: Implement search support! This should probably be done in conjunction with support
  in the server itself.
*/

import React, {Component, PropTypes} from 'react';
import {Card, CardText, CardHeader} from 'material-ui/Card';
import {connect} from 'react-redux';

class Search extends Component {
    static propTypes = {
        text: PropTypes.string
    };
    render() {
        if (this.props.text == "") {
            return null;
        }
        return (
            <Card style={{
                marginTop: "20px",
                textAlign: "left",
                marginBottom: "20px",
                backgroundColor: "#00b34a"
            }}>
                <CardHeader title={"Search"} titleColor="white" titleStyle={{
                    fontWeight: "bold"
                }} subtitle={this.props.text}/>
                <CardText style={{
                    textAlign: "center",
                    color: "white"
                }}>
                    <h3>Sorry, search is currently unimplemented.</h3>
                    <p>In future ConnectorDB versions, you would ask questions about your data here.</p>
                </CardText>
            </Card>
        );
    }
}
export default connect((state) => ({text: state.query.queryText}))(Search);
