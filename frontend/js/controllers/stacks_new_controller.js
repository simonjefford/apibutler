App.StacksNewController = Ember.ObjectController.extend({
    availableMiddlewares: [],

    middlewareQuery: '',

    filteredMiddlewares: function() {
        console.log('filter');
        if (Ember.isEmpty(this.get('middlewareQuery'))) {
            return this.get('availableMiddlewares');
        }

        var filterExp = new RegExp(this.get('middlewareQuery'), 'i');

        return this.get('availableMiddlewares').filter(function(item) {
            return filterExp.test(item.get('friendlyName'));
        });
    }.property('middlewareQuery', 'availableMiddlewares')
});
