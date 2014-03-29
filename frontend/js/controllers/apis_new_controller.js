/* global App */

App.ApisNewController = Ember.ObjectController.extend({
    saveDisabled: Ember.computed.not('apps.isFulfilled')
});
