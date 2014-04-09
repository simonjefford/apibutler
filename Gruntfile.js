var LIVERELOAD_PORT = 35729;

module.exports = function(grunt) {
    require('time-grunt')(grunt);
    grunt.initConfig({
        jshint: {
            all: {
                src: ['frontend/js/**/*.js',
                      'tests/**/*.js']
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
        neuter: {
            app: {
                options: {
                    filepathTransform: function(filepath) {
                        return 'frontend/' + filepath;
                    },
                    includeSourceMap: true
                },
                src: 'frontend/js/app.js',
                dest: '.tmp/combined-app.js'
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
                              'bower_components/ember-easyForm/**',
                              'bower_components/ic-ajax/**',
                              'bower_components/spin/**',
                              'bower_components/rickshaw/**',
                              'bower_components/ember-data/**',
                              'bower_components/jquery-ui/**',
                              'index.html']
                    },
                    {
                        dest: '.tmp/',
                        src: ['frontend/js/**/*.js']
                    }
                ]
            },
            justTheAppPlease: {
                files: [
                    {
                        dest: '.tmp/',
                        expand: true,
                        cwd: 'frontend',
                        src: ['index.html']
                    },
                    {
                        dest: '.tmp/',
                        src: ['frontend/js/**/*.js']
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
                tasks: ['copy:justTheAppPlease']
            },
            neuter: {
                files: [
                    'frontend/js/**/*.js'
                ],
                tasks: ['neuter']
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
            apibutler: {
                command: './apibutler',
                commandArgs: ['-frontendPath=.tmp']
            }
        },
        concurrent: {
            server: {
                tasks: [
                    'spawn:apibutler',
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
        'neuter',
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
        'neuter',
        'concurrent:server'
    ]);

    grunt.registerTask('default', [
        'jshint',
        'build'
    ]);
};
