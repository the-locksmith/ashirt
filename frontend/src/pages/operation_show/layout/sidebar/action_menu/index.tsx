// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import {default as Menu, MenuItem, MenuSeparator} from 'src/components/menu'
const cx = classnames.bind(require('./stylesheet'))

export default (props: {
  name: string,
  query: string,
  onDelete: () => void,
  onEdit: () => void,
}) => (
  <Menu>
    <div className={cx('top')}>
      <div className={cx('name')}>{props.name}</div>
      <div className={cx('query')}>{props.query}</div>
    </div>
    <MenuSeparator />
    <MenuItem icon={require('./edit.svg')} onClick={props.onEdit}>Edit Query</MenuItem>
    <MenuItem icon={require('./delete.svg')} onClick={props.onDelete}>Delete Query</MenuItem>
  </Menu>
)
