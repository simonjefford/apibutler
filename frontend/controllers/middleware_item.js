var MiddlewareItemController = Ember.ObjectController.extend({
    configurationToggled: false,

    canBeConfigured: Ember.computed.bool('needsConfiguration'),

    actions: {
        configure: function() {
            this.toggleProperty('configurationToggled');
            return true;
        },

        finishConfiguration: function() {
            this.toggleProperty('configurationToggled');
        }
    }
});

export default MiddlewareItemController;
