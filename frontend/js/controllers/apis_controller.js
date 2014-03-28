/* global App */

App.ApisController = Ember.ArrayController.extend({
    apisToShow: Ember.computed.filterBy('content', 'isNew', false)
});
