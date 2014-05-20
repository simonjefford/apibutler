var StackListView = Ember.View.extend({
    templateName: 'stack-list',

    didInsertElement: function() {
        var controller = this.get('controller');

        this.$('.sortable').sortable({
            axis: 'y',

            containment: 'parent',

            update: function() {
                var indexes = {};

                $(this).find('.item').each(function(index) {
                    indexes[$(this).data('id')] = index;
                });

                $(this).sortable('cancel');

                controller.updateSortOrder(indexes);
            }
        });
    }
});

export default StackListView;
