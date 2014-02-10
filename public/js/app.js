App = Ember.Application.create({
    LOG_TRANSITIONS: true
});


App.Router.map(function() {
    // put your routes here
});

App.IndexRoute = Ember.Route.extend({
    model: function() {
        return ['red', 'yellow', 'blue'];
    }
});
