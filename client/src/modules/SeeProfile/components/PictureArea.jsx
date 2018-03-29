import React, { Component } from 'react';
import './PictureArea.css';

const IndexPictures = (props) => {
  var elements = [];
  for (var i = 0; i < props.pictureArrayLength; i++) {
    let color = props.index === i ? "black" : "grey";
    elements.push(<span key={i} style={{ color: color }}>&#8226;</span>);
  }
  return (
    <div id="picture-index">
      {elements}
    </div>
  )
}

class PictureArea extends Component {
  constructor(props) {
    super(props);
    this.state = {
      moreInformationOpen: false,
      index: 0,
      pictureArrayLength: 0,
      picture: '',
    }
    this.openInformation = this.openInformation.bind(this);
    this.updateState = this.updateState.bind(this);
  }
  openInformation() {
    this.setState({
      moreInformationOpen: !this.state.moreInformationOpen,
    });
  }
  updateState(key, content) {
    this.setState({
      [key]: content,
    });
  }
  componentWillReceiveProps(nextProps) {
    this.updateState('pictureArrayLength', nextProps.pictureArrayLength);
    this.updateState('index', nextProps.index);
    this.updateState('picture', nextProps.picture);
  }
  render() {
    const index =  this.state.index;
    var pictureURL = '';
    if (this.state.picture) {
      pictureURL = "http://localhost:8080" + this.state.picture
    }
    var length = this.state.pictureArrayLength
    return (
      <div className="picture-area">
        <div className="picture-user-background" style={{ backgroundImage: "url(" + pictureURL + ")" }}></div>
        <div className="more-information" title={ 'See ' + (this.state.moreInformationOpen ? "less" :  "more") + ' options' }>
          <span className="more" onClick={this.openInformation}>{this.state.moreInformationOpen ? '-' : '+'}</span>
        </div>
        <div className="information-area" style={{ visibility: this.state.moreInformationOpen ? "visible" :  "hidden" } }>
            {this.props.reportedAsFakeAccount ?
              <span onClick={() => this.props.updateState("reportedAsFakeAccount", "You have just invalide you fake account report.")}>Invalidate fake account report</span> :
              <span onClick={() => this.props.updateState("reportedAsFakeAccount", "This profile has been declared as fake account.")}>Report as a fake account</span>
            }<br />
            {this.props.liked ? <span onClick={() => this.props.updateState("liked", "You have just unliked this profile")}>Unlike profile</span> : null}
        </div>
        <IndexPictures pictureArrayLength={length} index={index}/>
        <div id="picture-previous" style={{ visibility: (index === 0) ? "hidden" :  "visible" }}>
          <span className="arrow" onClick={() => this.props.changePicture(0, length)}>&#x21A9;</span>
        </div>
        <div id="picture-next" style={{ visibility: (!length || index === (length - 1)) ? "hidden" :  "visible" } }>
          <span className="arrow" onClick={() => this.props.changePicture(1, length)}>&#x21AA;</span>
        </div>
        {!this.props.usersAreConnected ? (
        <div title="Like profile" className="btn-like"
          onClick={() => this.props.likeUser()}
          style={{
            background: (this.props.liked ? "white" :  "#F80759"),
            color: (this.props.liked ? "#F80759" :  "white"),
            cursor: (this.props.liked ? "default" :  "pointer") }}
          >
          <span>&#9829;</span>
        </div>
        ) : null }
        <div>
        {this.props.usersAreConnected ? (
          <div className="profiles-linked-picture">
            <span
              role="img" aria-labelledby="Connected with"
              title={'You are connected with ' + this.props.firstname + ' - Click here to take contact ;)' }
            >&#x1f525;</span>
          </div>
        ) : null }
        </div>
      </div>
    )
  }
}

export default PictureArea;
