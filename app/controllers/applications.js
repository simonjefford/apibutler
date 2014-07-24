var ApplicationsController = Ember.ArrayController.extend({
    renderer: 'line',

    renderers: ['area', 'line', 'bar', 'scatterplot'],

    headers: ['Name', 'Backend URL']
});

export default ApplicationsController;
