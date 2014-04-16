App.Api = DS.Model.extend({
    fragment: DS.attr('string'),
    limit: DS.attr('number'),
    seconds: DS.attr('number'),
    isPrefix: DS.attr('boolean'),
    needsAuth: DS.attr('boolean'),
    app: DS.belongsTo('app'),

    valid: function() {
        return !Ember.isEmpty(this.get('app')) &&
            !Ember.isEmpty(this.get('fragment'));
    }.property('app', 'fragment')
});
