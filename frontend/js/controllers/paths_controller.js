/* global App */

App.PathsController = Ember.ArrayController.extend({
    pathsToShow: Ember.computed.filterBy('content', 'isNew', false)
});
