/* globals require, process */

var UI = require('ember-cli/lib/ui');
var ui = new UI({
    inputStream: process.stdin,
    outputStream: process.stdout
});
var Project = require('ember-cli/lib/models/project');
var project = Project.closest(process.cwd());
var path = require('path');
var project = new Project(process.cwd(), require(path.join(process.cwd(), 'package.json')));
var buildWatcher = require('ember-cli/lib/utilities/build-watcher');
var livereload = require('ember-cli/lib/tasks/server/livereload-server');

module.exports = function(grunt) {
    require('time-grunt')(grunt);
    require('matchdep').filterDev('grunt-*').forEach(grunt.loadNpmTasks);

    grunt.initConfig({
        jshint: {
            all: {
                src: ['frontend/**/*.js',
                      'tests/**/*.js']
            },
            options: {
                jshintrc: true
            }
        },
        spawn: {
            apibutler: {
                command: './apibutler',
                commandArgs: ['-frontendPath=tmp/output']
            }
        },
        concurrent: {
            server: {
                tasks: [
                    'spawn:apibutler',
                    'watchAndBuild'
                ],
                options: {
                    logConcurrentOutput: true,
                }
            }
        },
        broccoli: {
            default: {
                dest: 'dist'
            }
        }
    });

    grunt.registerTask('watchAndBuild', 'watch and build for app changes (with Live Reload)', function() {
        this.async();
        var leek = require('leek');
        var watcher = buildWatcher({
            ui: ui,
            analytics: leek
        });

        var opts = {
            environment: 'development',
            port: 4200,
            host: '0.0.0.0',
            project: project,
            watcher: watcher,
            liveReload: true,
            liveReloadPort: 35729
        };

        livereload.ui = ui;
        livereload.analytics = leek;

        livereload.start(opts).then(function(res) { console.log("Result ", res)}, function(err) { console.log(err) });
    });

    grunt.registerTask('default', [
        'jshint',
        'broccoli:default:build'
    ]);

    grunt.registerTask('serve', [
        'concurrent:server'
    ]);
};
