App = Ember.Application.create({
    LOG_TRANSITIONS: true
});

App.Router.map(function() {
    this.resource("paths", function() {
        this.route('new');
    });
});

App.PathsIndexRoute = Ember.Route.extend({
    model: function() {
        return { page: "paths "};
    }
});
