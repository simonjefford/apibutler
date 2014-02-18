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
        }
    });

    require('matchdep').filterDev('grunt-*').forEach(grunt.loadNpmTasks);

    grunt.registerTask('build', [
        'useminPrepare',
        'emberTemplates',
        'compass',
        'concat',
        'uglify',
        'cssmin',
        'copy',
        'usemin'
    ]);

    grunt.registerTask('default', [
        'jshint',
        'build'
    ]);
};
