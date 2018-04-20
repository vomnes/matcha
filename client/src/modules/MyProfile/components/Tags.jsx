import React from 'react';
import './Tags.css';
import api from '../../../library/api'

const unlinkTag = async (args, updateState, deleteTag) => {
  try {
    let res = await api.tag(`DELETE`, args);
    if (res && res.status >= 400) {
      const response = await res.json();
      if (res.status >= 500) {
        throw new Error("Bad response from server - Tag has failed - ", response.error);
      } else if (res.status >= 400) {
        updateState('newError', response.error);
        return;
      }
    }
    updateState('newSuccess', 'Tag ' + args.tagName + ' has been deleted');
    deleteTag(args.tagName, args.tagID);
  } catch (e) {
    console.log(e.message);
  }
}

const MyTag = (props) => {
  return (
    <div key={props.index} title={'Click here to remove the tag ' + props.tagName} className="profile-tag profile-matched-tag">
      <span className="position-tag" onClick={() => unlinkTag({tagName: props.tagName, tagID: props.tagID,}, props.updateState, props.deleteTag)}>
        #{props.tagName}
      </span>
    </div>
  )
}

const addTag = async (tags, args, updateState, appendTag) => {
  const nbTags = tags.length;
  const newTag = args.tagName.toLowerCase()
  for (var i = 0; i < nbTags; i++) {
    if (tags[i].name === newTag) {
      updateState('newError', 'Tag "' + newTag + '" is already linked with your profile');
      return;
    }
  }
  try {
    let res = await api.tag(`POST`, args);
    const response = await res.json();
    if (res.status >= 500) {
      throw new Error("Bad response from server - Tag has failed - ", response.error);
    } else if (res.status >= 400) {
      updateState('newError', response.error);
    } else {
      updateState('newError', '');
      appendTag(args.tagName, response.tag_id);
    }
  } catch (e) {
    console.log(e.message);
  }
}

const NewTag = (props) => {
  var confirm;
  if (props.newTag) {
    confirm = (<span className="valid" title="Valid new tag" onClick={() => addTag(props.tags, {tagName: props.newTag}, props.updateState, props.appendTag)}>✓</span>)
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
    tags.push(<MyTag deleteTag={props.deleteTag} key={index} index={index} tagID={tag.id} tagName={tag.name} updateState={props.updateState}/>);
    index += 1;
  });
  return (
    <div id="data-profile-tags">
      {tags}
      <NewTag
        tags={props.tags}
        index={index}
        deleteTag={props.deleteTag}
        newTag={props.newTag}
        handleUserInput={props.handleUserInput}
        appendTag={props.appendTag}
        updateState={props.updateState}
      />
    </div>
  )
}

export default Tags;
