import PropTypes from 'prop-types'
import ImmutablePropTypes from 'react-immutable-proptypes'

// from react-redux/src/utils/PropTypes
const storeShape = PropTypes.shape({
  subscribe: PropTypes.func.isRequired,
  dispatch: PropTypes.func.isRequired,
  getState: PropTypes.func.isRequired,
})

export default {
  ...ImmutablePropTypes,
  ...PropTypes,
  immutableShape: ImmutablePropTypes.shape, // fix intersection with PropTypes.shape
  storeShape,
}

export { PropTypes, ImmutablePropTypes }
