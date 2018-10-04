import Immutable from 'immutable'
import { createSelector } from 'reselect'
import { identity } from 'lodash'

export function composeSelector(query, makeResult = identity) {
  const selectors = Immutable.OrderedMap(query).reduce(
    (res, selector) => res.concat(selector),
    []
  )
  return createSelector(selectors, (...results) =>
    makeResult(Object.assign({}, ...results))
  )
}
