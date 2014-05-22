var StacksNewRoute = Ember.Route.extend({
    model: function() {
        var stack = this.store.createRecord('stack');
        stack.set('middlewares', Ember.A());
        return stack;
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
            controller.resetSelected();
        }
    }
});

export default StacksNewRoute;
