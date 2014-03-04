var LIVERELOAD_PORT = 35729;

module.exports = function(grunt) {
    grunt.initConfig({
        jshint: {
            all: {
                src: ['frontend/js/app.js']
            },
            options: {
                jshintrc: true
            }
        },
        useminPrepare: {
            html: 'frontend/index.html',
            options: {
                dest: 'frontend/dist'
            }
        },
        usemin: {
            html: 'frontend/dist/index.html',
            options: {
                dirs: ['frontend/dist']
            }
        },
        copy: {
            dist: {
                src: 'frontend/index.html',
                dest: 'frontend/dist/index.html'
            },
            server: {
                files: [
                    {
                        dest: '.tmp/',
                        expand: true,
                        cwd: 'frontend',
                        src: ['bower_components/bootstrap-sass/js/**',
                              'bower_components/jquery/jquery.js',
                              'bower_components/handlebars/**',
                              'bower_components/ember/**',
                              'index.html',
                              'js/**']
                    }
                ]
            }
        },
        emberTemplates: {
            options: {
                templateName: function (sourceFile) {
                    var templatePath = 'frontend/templates/';
                    return sourceFile.replace(templatePath, '');
                }
            },
            dist: {
                files: {
                    '.tmp/compiled-templates.js': 'frontend/templates/{,*/}*.hbs'
                }
            }
        },
        compass: {
            options: {
                sassDir: 'frontend/css',
                cssDir: '.tmp/css',
                generatedImagesDir: '.tmp/images/generated',
                imagesDir: 'frontend/images',
                javascriptsDir: 'frontend/js',
                fontsDir: 'frontend/css/fonts',
                importPath: 'frontend/bower_components',
                relativeAssets: false
            },
            dist: {}
        },
        clean: {
            dist: {
                files: [{
                    src: [
                        '.tmp',
                        'frontend/dist'
                    ]
                }]
            },
            server: '.tmp'
        },
        watch: {
            emberTemplates: {
                files: 'frontend/templates/**/*.hbs',
                tasks: ['emberTemplates']
            },
            compass: {
                files: ['frontend/css/{,*/}*.{scss,sass}'],
                tasks: ['compass']
            },
            copy: {
                files: [
                    'frontend/index.html',
                    'frontend/js/**/*.js'
                ],
                tasks: ['copy:server']
            },
            livereload: {
                options: {
                    livereload: LIVERELOAD_PORT
                },
                files: [
                    '.tmp/**/*.js',
                    '.tmp/*.html',
                    '.tmp/css/*.css'
                ]
            }
        },
        spawn: {
            ratelimit: {
                command: './rateLimit',
                commandArgs: ['-frontendPath=.tmp']
            }
        },
        concurrent: {
            server: {
                tasks: [
                    'spawn:ratelimit',
                    'watch'
                ],
                options: {
                    logConcurrentOutput: true,
                }
            }
        }
    });

    require('matchdep').filterDev('grunt-*').forEach(grunt.loadNpmTasks);

    grunt.registerTask('build', [
        'clean:dist',
        'useminPrepare',
        'emberTemplates',
        'compass',
        'concat',
        'uglify',
        'cssmin',
        'copy:dist',
        'usemin'
    ]);

    grunt.registerTask('serve', [
        'clean:server',
        'copy:server',
        'emberTemplates',
        'compass',
        'concurrent:server'
    ]);

    grunt.registerTask('default', [
        'jshint',
        'build'
    ]);
};
