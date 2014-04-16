App.ApisNewController = Ember.ObjectController.extend({
    saveEnabled: Ember.computed.alias('apps.isFulfilled'),
    saveDisabled: Ember.computed.not('saveEnabled'),

    actions: {
        submit: function() {
            this.send('save', this.get('content'));
        }
    }
});
