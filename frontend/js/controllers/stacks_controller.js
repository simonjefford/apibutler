App.StacksController = Ember.ArrayController.extend({
    sortProperties: ['idx'],

    updateSortOrder: function(indexes) {
        this.beginPropertyChanges();
        this.forEach(function(item) {
            var index = indexes[item.get('id')];
            item.set('idx', index);
        }, this);
        this.endPropertyChanges();
    }
});
