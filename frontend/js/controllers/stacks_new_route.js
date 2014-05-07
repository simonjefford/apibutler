App.StacksNewRoute = Ember.Route.extend({
    model: function() {
        return this.store.createRecord('stack');
    },

    availableMiddlewares: null,

    setupController: function(controller, model) {
        this._super(controller, model);
        if(!this.get('availableMiddlewares')) {
            this.set('availableMiddlewares', this.store.find('middleware'));
        }
        controller.set('availableMiddlewares', this.get('availableMiddlewares'));
    },

    actions: {
        willTransition: function() {
            var controller = this.controllerFor('stacks.new');
            if (controller.get('isNew')) {
                controller.content.deleteRecord();
            }
        }
    }
});
