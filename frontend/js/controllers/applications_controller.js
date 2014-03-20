/* global App */

App.ApplicationsController = Ember.ArrayController.extend({
    renderer: 'line',

    renderers: ['area', 'line', 'bar', 'scatterplot']
});
