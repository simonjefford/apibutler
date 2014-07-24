var StacksController = Ember.ArrayController.extend({
    stackList: Ember.computed.filterBy('model', 'isNew', false)
});

export default StacksController;
