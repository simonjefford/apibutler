export default Ember.Mixin.create({
    apps: null,

    stacks: null,

    setupController: function(controller, model) {
        this._super(controller, model);
        if(!this.get('apps')) {
            this.set('apps', this.store.find('app'));
        }

        if(!this.get('stacks')) {
            this.set('stacks', this.store.find('stack'));
        }
        controller.set('apps', this.get('apps'));
        controller.set('stacks', this.get('stacks'));
    },

    renderTemplate: function() {
        this.render('apis/api');
    },

    actions: {
        cancel: function() {
            this.transitionTo('apis');
        }
    }
});
