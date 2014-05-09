/* global App */

App.MiddlewaresRoute = Ember.Route.extend({
    model: function() {
        return this.store.find('middleware_definition');
    }
});
