var ApisRoute = Ember.Route.extend({
    model: function() {
        return this.store.find('app').then(function() {
            return this.store.find('api');
        }.bind(this));
    }
});

export default ApisRoute;
