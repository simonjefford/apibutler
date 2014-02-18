module.exports = function(grunt) {
    grunt.initConfig({
        jshint: {
            all: {
                src: ['public/js/app.js']
            },
            options: {
                jshintrc: true
            }
        },
        uglify: {
            dist: {
                options: {
                    sourceMap: true
                },
                files: {
                    'public/js/app.min.js': [
                        'public/js/libs/jquery-2.0.0.js',
                        'public/js/libs/handlebars-1.1.2.js',
                        'public/js/libs/ember-1.3.1.js',
                        'public/js/libs/ic-ajax.js',
                        'public/js/libs/ember-easyform.js',
                        'public/js/app.js'
                    ]
                }
            }
        }
    });

    require('matchdep').filterDev('grunt-*').forEach(grunt.loadNpmTasks);
};
