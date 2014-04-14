App.ApplicationsRoute = Ember.Route.extend({
    model: function() {
        return this.store.find('app');
    }
});
