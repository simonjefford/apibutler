// Karma configuration
// Generated on Mon Mar 17 2014 10:55:13 GMT+0000 (GMT)

module.exports = function(config) {
    config.set({
        // base path that will be used to resolve all patterns (eg. files, exclude)
        basePath: '',

        // frameworks to use
        // available frameworks: https://npmjs.org/browse/keyword/karma-adapter
        frameworks: ['qunit'],

        // list of files / patterns to load in the browser
        files: [
            'frontend/bower_components/jquery/jquery.js',
            'frontend/bower_components/handlebars/handlebars.runtime.js',
            'frontend/bower_components/ember/ember.js',
            'frontend/bower_components/ic-ajax/dist/globals/main.js',
            'frontend/bower_components/ember-easyForm/index.js',
            'frontend/bower_components/bootstrap-sass/js/collapse.js',
            'frontend/js/app.js',
            'tests/*.js'
        ],

        // test results reporter to use
        // possible values: 'dots', 'progress'
        // available reporters: https://npmjs.org/browse/keyword/karma-reporter
        reporters: ['progress', 'growl'],

        // web server port
        port: 9876,

        // enable / disable colors in the output (reporters and logs)
        colors: true,

        // level of logging
        // possible values: config.LOG_DISABLE || config.LOG_ERROR || config.LOG_WARN || config.LOG_INFO || config.LOG_DEBUG
        logLevel: config.LOG_INFO,

        // enable / disable watching file and executing tests whenever any file changes
        autoWatch: true,

        // start these browsers
        // available browser launchers: https://npmjs.org/browse/keyword/karma-launcher
        browsers: ['PhantomJS'],

        // Continuous Integration mode
        // if true, Karma captures browsers, runs the tests and exits
        singleRun: false
    });
};
