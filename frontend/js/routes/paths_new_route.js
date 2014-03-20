/* global App */

App.PathsNewRoute = Ember.Route.extend({
    model: function() {
        return App.Path.create();
    },

    actions: {
        save: function(model) {
            var self = this;
            console.log('Now saving %o', JSON.stringify(model));
            model.save().then(function(result) {
                self.controllerFor('paths').addObject(model);
                self.transitionTo('paths.index');
                return result;
            }, function(arg) {
                console.log(arg);
            });
        }
    }
});
