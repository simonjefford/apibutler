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
        useminPrepare: {
            html: 'public/index.html',
            options: {
                dest: 'public/dist'
            }
        },
        usemin: {
            html: 'public/dist/index.html',
            options: {
                dirs: ['public/dist']
            }
        },
        copy: {
            dist: {
                src: 'public/index.html',
                dest: 'public/dist/index.html'
            }
        }
    });

    require('matchdep').filterDev('grunt-*').forEach(grunt.loadNpmTasks);

    grunt.registerTask('build', [
        'useminPrepare',
        'concat',
        'uglify',
        'copy',
        'usemin'
    ]);
};
