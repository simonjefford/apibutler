/* global App */
App.PathsRoute = Ember.Route.extend({
    model: function() {
        return this.store.find('path');
    }
});
