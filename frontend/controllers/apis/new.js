var ApisNewController = Ember.ObjectController.extend({
    appsReady: Ember.computed.alias('apps.isFulfilled'),

    saveEnabled: Ember.computed.and('appsReady', 'valid'),

    saveDisabled: Ember.computed.not('saveEnabled'),

    saveSucceeded: false,

    showSuccessAlert: Ember.computed.and('saveSucceeded', 'isClean'),

    actions: {
        submit: function() {
            if (this.get('valid')) {
                this.send('save', this.get('content'));
            }
        }
    }
});

export default ApisNewController;