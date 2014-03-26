/* global App */

App.PathsNewRoute = Ember.Route.extend({
    model: function() {
        return this.store.createRecord('path');
    },

    actions: {
        save: function(model) {
            var self = this;
            console.log('Now saving %o', JSON.stringify(model));
            model.save().then(function() {
                self.transitionTo('paths.index');
            }, function(arg) {
                console.log(arg);
            });
        }
    }
});
