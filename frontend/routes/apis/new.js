import ApiRouteMixin from 'apibutler/mixins/api_route';

export default Ember.Route.extend(ApiRouteMixin, {
    model: function() {
        return this.store.createRecord('api');
    },

    actions: {
        willTransition: function() {
            var controller = this.controllerFor('apis.new');
            if(controller.get('isNew')) {
                controller.content.deleteRecord();
            }
            controller.set('saveSucceeded', false);
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
