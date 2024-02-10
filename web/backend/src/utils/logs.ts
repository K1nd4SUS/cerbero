const internalConsoleLog = console.log
const internalConsoleError = console.error

console.log = function (...args: []) {
  internalConsoleLog(new Date(), "[LOG]", ...args)
}

console.info = function (...args: []) {
  internalConsoleLog(new Date(), "[INFO]", ...args)
}

console.debug = function (...args: []) {
  internalConsoleLog(new Date(), "[DEBUG]", ...args)
}

console.warn = function (...args: []) {
  internalConsoleError(new Date(), "[WARNING]", ...args)
}

console.error = function (...args: []) {
  internalConsoleError(new Date(), "[ERROR]", ...args)
}
