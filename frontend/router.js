var Router = Ember.Router.extend({
    location: ApibutlerENV.locationType
});

Router.map(function() {
    this.resource('apis', function() {
        this.route('new');
    });
    this.resource('applications', function() {
        this.route('new');
    });
    this.resource('middlewares');
    this.resource('stacks', function() {
        this.route('new');
        this.resource('stack', { path: ':stack_id' });
    });
});

export default Router;
