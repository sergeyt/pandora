import { isFunction } from 'lodash'

export function isIterable(obj) {
  return obj && typeof obj !== 'string' && obj[Symbol.iterator]
}

export function isGenerator(obj) {
  return isFunction(obj.next) && isFunction(obj.throw)
}

export function isGeneratorFunction(obj) {
  const constructor = obj.constructor
  if (!constructor) {
    return false
  }
  if (
    'GeneratorFunction' === constructor.name ||
    'GeneratorFunction' === constructor.displayName
  ) {
    return true
  }
  return isGenerator(constructor.prototype)
}
