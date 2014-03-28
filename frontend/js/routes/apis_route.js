/* global App */
App.ApisRoute = Ember.Route.extend({
    model: function() {
        return this.store.find('path');
    }
});
