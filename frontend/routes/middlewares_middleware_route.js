/* global App */

App.MiddlewaresMiddlewareRoute = Ember.Route.extend({
    model: function(params) {
        return App.Middlewares[parseInt(params.middleware_id, 10) - 1];
    }
});
