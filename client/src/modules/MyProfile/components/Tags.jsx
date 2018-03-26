import React from 'react';
import './Tags.css';

const MyTag = (props) => {
  return (
    <div key={props.index} title={'Click here to remove the tag ' + props.tag} className="profile-tag profile-matched-tag">
      <span className="position-tag" onClick={() => props.deleteTag(props.tag)}>
        #{props.tag}
      </span>
    </div>
  )
}

const NewTag = (props) => {
  var confirm = null;
  if (props.newTag) {
    confirm = (<span className="valid" title="Valid new tag" onClick={() => props.appendTag(props.newTag)}>✓</span>)
  }
  return (
    <div className="profile-tag new-tag">
      <span className="position-tag">#
      <input id="new-tag" placeholder="Add a new tag" type="text" name="newTag"
        value={props.newTag} onChange={props.handleUserInput}/>
      {confirm}
      </span>
    </div>
  )
}

const Tags = (props) => {
  var index = 0;
  var tags = [];
  props.tags.forEach((tag) => {
    tags.push(<MyTag deleteTag={props.deleteTag} key={index} index={index} tag={tag}/>);
    index += 1;
  });
  return (
    <div >
      <div id="data-profile-tags">
        {tags}
        <NewTag
          tags={props.tags}
          index={index}
          deleteTag={props.deleteTag}
          newTag={props.newTag}
          handleUserInput={props.handleUserInput}
          appendTag={props.appendTag}
        />
      </div>
    </div>
  )
}

export default Tags;
