/* global require, module */

var EmberApp = require('ember-cli/lib/broccoli/ember-app');

var app = new EmberApp({
    name: require('./package.json').name,

    minifyCSS: {
        enabled: true,
        options: {}
    },
    getEnvJSON: require('./config/environment'),

    trees: {
        app: 'frontend',
        styles: 'frontend/css'
    }
});

// Use this to add additional libraries to the generated output files.
// If you need to use different assets in different
// environments, specify an object as the first parameter. That
// object's keys should be the environment name and the values
// should be the asset to use in that environment.
//
// If the library that you are including contains AMD or ES6
// modules that you would like to import into your application
// please specify an object with the list of modules as keys
// along with the exports of each module as its value.

app.import('vendor/spin/index.js');
app.import('vendor/ember-easyForm/index.js');
app.import('vendor/bootstrap-sass-official/vendor/assets/javascripts/bootstrap/collapse.js');
app.import('vendor/bootstrap-sass-official/vendor/assets/javascripts/bootstrap/modal.js');
app.import('vendor/html5sortable/jquery.sortable.js');

module.exports = app.toTree();
