var LIVERELOAD_PORT = 35729;

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
            },
            server: {
                files: [
                    {
                        dest: '.tmp/',
                        expand: true,
                        cwd: 'public',
                        src: ['js/**', 'index.html']
                    }
                ]
            }
        },
        emberTemplates: {
            options: {
                templateName: function (sourceFile) {
                    var templatePath = 'public/templates/';
                    return sourceFile.replace(templatePath, '');
                }
            },
            dist: {
                files: {
                    '.tmp/compiled-templates.js': 'public/templates/{,*/}*.hbs'
                }
            }
        },
        compass: {
            options: {
                sassDir: 'public/css',
                cssDir: '.tmp/css',
                generatedImagesDir: '.tmp/images/generated',
                imagesDir: 'public/images',
                javascriptsDir: 'public/js',
                fontsDir: 'public/css/fonts',
                importPath: 'bower_components',
                // httpImagesPath: '/images',
                // httpGeneratedImagesPath: '/images/generated',
                // httpFontsPath: '/styles/fonts',
                relativeAssets: false
            },
            dist: {}
        },
        clean: {
            dist: {
                files: [{
                    src: [
                        '.tmp',
                        'public/dist'
                    ]
                }]
            },
            server: '.tmp'
        },
        watch: {
            emberTemplates: {
                files: 'public/templates/**/*.hbs',
                tasks: ['emberTemplates']
            },
            compass: {
                files: ['public/css/{,*/}*.{scss,sass}'],
                tasks: ['compass']
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
    ]);

    grunt.registerTask('default', [
        'jshint',
        'build'
    ]);
};
