/* global require, module */

var EmberApp = require('ember-cli/lib/broccoli/ember-app');

var app = new EmberApp({
    name: require('./package.json').name,

    wrapInEval: true,

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
app.import('vendor/ember-data/ember-data.js');
app.import('vendor/spin/index.js');
app.import('vendor/jquery-ui/ui/jquery-ui.js');
app.import('vendor/ember-easyForm/index.js');
app.import('vendor/bootstrap-sass-official/vendor/assets/javascripts/bootstrap/collapse.js');
app.import('vendor/bootstrap-sass-official/vendor/assets/javascripts/bootstrap/modal.js');

// If the library that you are including contains AMD or ES6 modules that
// you would like to import into your application please specify an
// object with the list of modules as keys along with the exports of each
// module as its value.
app.import('vendor/ic-ajax/dist/named-amd/main.js', {
    'ic-ajax': [
        'default',
        'defineFixture',
        'lookupFixture',
        'raw',
        'request',
    ]
});


module.exports = app.toTree();
