App.ApisNewRoute = Ember.Route.extend({
    model: function() {
        var api = this.store.createRecord('api');

        api.setProperties({
            isPrefix: true,
            needsAuth: true
        });

        return api;
    },

    apps: null,

    setupController: function(controller, model) {
        this._super(controller, model);
        if(!this.get('apps')) {
            this.set('apps', this.store.find('app'));
        }
        controller.set('apps', this.get('apps'));
    },

    actions: {
        willTransition: function() {
            var controller = this.controllerFor('apis.new');
            if(controller.get('isNew')) {
                controller.content.deleteRecord();
            }
        },

        cancel: function() {
            this.transitionTo('apis');
        },

        save: function(model) {
            model.save().then(function() {
                var newModel = this.model();
                var controller = this.controllerFor('apis.new');
                this.setupController(controller, newModel);
            }.bind(this));
        }
    }
});
