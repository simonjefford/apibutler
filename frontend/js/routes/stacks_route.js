App.StacksRoute = Ember.Route.extend({
    model: function() {
        return App.Middlewares;
    }
});
