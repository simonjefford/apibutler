App.Api = DS.Model.extend(App.ModelMixin, {
    path: DS.attr('string'),
    needsAuth: DS.attr('boolean'),
    app: DS.belongsTo('app'),

    valid: function() {
        return !Ember.isEmpty(this.get('app')) &&
            !Ember.isEmpty(this.get('path'));
    }.property('app', 'path')
});
