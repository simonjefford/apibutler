var MiddlewareItemController = Ember.ObjectController.extend({
    configurationToggled: false,

    canBeConfigured: Ember.computed.bool('underlying.schema.length'),

    actions: {
        toggleConfiguration: function() {
            this.toggleProperty('configurationToggled');
        }
    }
});

export default MiddlewareItemController;
