var ajax = ic.ajax;

App = Ember.Application.create({
    LOG_TRANSITIONS: true
});

App.Router.map(function() {
    this.resource("paths", function() {
        this.route('new');
    });
});

App.PathsRoute = Ember.Route.extend({
    model: function() {
        return ajax('/paths');
    }
});
