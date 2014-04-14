/* global App */

App.MiddlewaresRoute = Ember.Route.extend({
    model: function() {
        return App.Middlewares;
    },

    actions: {
        viewMiddleware: function(middleware) {
            this.transitionTo('middlewares.middleware', middleware);
        }
    }
});
