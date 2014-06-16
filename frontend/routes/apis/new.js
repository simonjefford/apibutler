var ApisNewRoute = Ember.Route.extend({
    model: function() {
        return this.store.createRecord('api');
    },

    apps: null,

    stacks: null,

    setupController: function(controller, model) {
        this._super(controller, model);
        if(!this.get('apps')) {
            this.set('apps', this.store.find('app'));
        }

        if(!this.get('stacks')) {
            this.set('stacks', this.store.find('stack'));
        }
        controller.set('apps', this.get('apps'));
        controller.set('stacks', this.get('stacks'));
    },

    actions: {
        willTransition: function() {
            var controller = this.controllerFor('apis.new');
            if(controller.get('isNew')) {
                controller.content.deleteRecord();
            }
            controller.set('saveSucceeded', false);
        },

        cancel: function() {
            this.transitionTo('apis');
        },

        save: function(model) {
            model.save().then(function() {
                var newModel = this.model();
                var controller = this.controllerFor('apis.new');
                controller.set('saveSucceeded', true);
                this.setupController(controller, newModel);
            }.bind(this));
        }
    }
});

export default ApisNewRoute;
