const fs = require('fs');
const vm = require('vm');

const code = fs.readFileSync('temp.js', 'utf8');

const sandbox = {
    require: require,
    console: console,
    process: {
        stdin: process.stdin,
        stdout: process.stdout
    }
};

try {
    vm.createContext(sandbox);
    vm.runInContext(code, sandbox);
} catch (error) {
    console.error('Error during execution:', error.message);
}