App.ApisNewController = Ember.ObjectController.extend({
    appsReady: Ember.computed.alias('apps.isFulfilled'),

    saveEnabled: Ember.computed.and('appsReady', 'valid'),

    saveDisabled: Ember.computed.not('saveEnabled'),

    actions: {
        submit: function() {
            if (this.get('valid')) {
                this.send('save', this.get('content'));
            }
        }
    }
});
