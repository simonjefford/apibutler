/* global App */

App.ApisNewRoute = Ember.Route.extend({
    model: function() {
        return this.store.createRecord('path');
    },

    actions: {
        willTransition: function() {
            var controller = this.controllerFor('apis.new');
            if(controller.get('isNew')) {
                controller.content.deleteRecord();
            }
        },

        save: function(model) {
            var self = this;
            console.log('Now saving %o', JSON.stringify(model));
            model.save().then(function() {
                self.transitionTo('apis.index');
            }, function(arg) {
                console.log(arg);
            });
        }
    }
});
