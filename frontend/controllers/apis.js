var ApisController = Ember.ArrayController.extend({
    apisToShow: Ember.computed.filterBy('content', 'isNew', false),

    headers: ['Path', 'App'],
});

export default ApisController;
