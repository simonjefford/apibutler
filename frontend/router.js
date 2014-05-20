var Router = Ember.Router.extend({
    location: ENV.locationType
});

Router.map(function() {
    this.resource('apis', function() {
        this.route('new');
    });
    this.resource('applications', function() {
        this.route('new');
    });
    this.resource('middlewares');
    this.resource('stacks');
});

export default Router;
