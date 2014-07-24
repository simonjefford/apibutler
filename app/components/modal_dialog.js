var ModalDialogComponent = Ember.Component.extend({
    didInsertElement: function() {
        this.$().modal({backdrop: 'static'});
    },

    willDestroyElement: function() {
        this.$().modal('hide');
    },

    classNames: ['modal', 'fade'],

    saveDisabled: false,

    actions: {
        saveClick: function() {
            this.sendAction('save', this.get('modelToSave'));
        },

        cancelClick: function() {
            this.sendAction('cancel');
        }
    }
});

export default ModalDialogComponent;
