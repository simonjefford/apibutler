/* global App */
App.PathsRoute = Ember.Route.extend({
    model: function() {
        return App.Path.findAll();
    }
});
